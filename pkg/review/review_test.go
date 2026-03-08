package review

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"review/pkg/diff"
)

type capturingReviewer struct {
	prompt string
	err    error
}

func (c *capturingReviewer) Review(prompt string, stdout, stderr io.Writer) error {
	c.prompt = prompt
	return c.err
}

func TestNewRunnerSetsDefaultLimit(t *testing.T) {
	runner := NewRunner(nil)

	if runner.MaxDiffLines != diff.DefaultMaxLines {
		t.Fatalf("NewRunner().MaxDiffLines = %d, want %d", runner.MaxDiffLines, diff.DefaultMaxLines)
	}
}

func TestBuildPromptAddsCommitMessageNewlineAndFiles(t *testing.T) {
	files := []diff.FileDiff{{FileName: "pkg/review/review.go"}}

	prompt := BuildPrompt("feat: test", files, "diff --git a b")

	if !strings.Contains(prompt, "コミットメッセージ:\nfeat: test\nファイル一覧:\n- pkg/review/review.go\n") {
		t.Fatalf("BuildPrompt() prompt missing expected sections:\n%s", prompt)
	}
	if !strings.HasSuffix(prompt, "\n差分:\ndiff --git a b") {
		t.Fatalf("BuildPrompt() prompt suffix = %q", prompt)
	}
}

func TestRunNoDiffWritesMessage(t *testing.T) {
	runner := NewRunner(nil)
	var out bytes.Buffer
	var errOut bytes.Buffer

	if err := runner.Run("msg", "", &out, &errOut); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}
	if out.String() != "No diff found\n" {
		t.Fatalf("Run() out = %q, want %q", out.String(), "No diff found\n")
	}
}

func TestRunVerbosePrintsPromptBeforeReview(t *testing.T) {
	reviewer := &capturingReviewer{}
	runner := NewRunner(reviewer)
	runner.Verbose = true
	var out bytes.Buffer
	var errOut bytes.Buffer
	rawDiff := strings.Join([]string{
		"diff --git a/a.go b/a.go",
		"@@ -1 +1 @@",
		"+line",
	}, "\n")

	if err := runner.Run("feat: msg", rawDiff, &out, &errOut); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}
	if !strings.Contains(out.String(), "Prompt:\n") {
		t.Fatalf("Run() out = %q, want prompt header", out.String())
	}
	if !strings.Contains(out.String(), "Analyzing diff...\n") {
		t.Fatalf("Run() out = %q, want analyzing message", out.String())
	}
	if reviewer.prompt == "" {
		t.Fatal("Run() did not send prompt to reviewer")
	}
}

func TestRunReturnsTooLargeError(t *testing.T) {
	runner := NewRunner(&capturingReviewer{})
	runner.MaxDiffLines = 1
	var out bytes.Buffer
	var errOut bytes.Buffer
	rawDiff := strings.Join([]string{
		"diff --git a/a.go b/a.go",
		"@@ -1 +1 @@",
	}, "\n")

	err := runner.Run("feat: msg", rawDiff, &out, &errOut)
	if err == nil {
		t.Fatal("Run() error = nil, want error")
	}
	if !errors.Is(err, diff.ErrTooLarge) {
		t.Fatalf("Run() error = %v, want ErrTooLarge", err)
	}
}
