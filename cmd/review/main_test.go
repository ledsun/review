package main

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"review/pkg/diff"
	"review/pkg/git"
	"review/pkg/review"
)

type stubReviewer struct{}

func (s *stubReviewer) Review(prompt string, stdout, stderr io.Writer) error {
	return nil
}

func TestRunWritesNoDiffFound(t *testing.T) {
	restore := stubMainDeps()
	defer restore()

	checkGitInstalled = func() error { return nil }
	latestCommitDiff = func() (string, error) { return "", nil }
	latestCommitMessage = func() (string, error) { return "feat: msg", nil }

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if err := run(nil, &stdout, &stderr); err != nil {
		t.Fatalf("run() error = %v, want nil", err)
	}
	if stdout.String() != "No diff found\n" {
		t.Fatalf("run() stdout = %q, want %q", stdout.String(), "No diff found\n")
	}
}

func TestRunAppliesLimitBreakFlag(t *testing.T) {
	restore := stubMainDeps()
	defer restore()

	checkGitInstalled = func() error { return nil }
	latestCommitDiff = func() (string, error) {
		return "diff --git a/a.go b/a.go\n@@ -1 +1 @@\n+line", nil
	}
	latestCommitMessage = func() (string, error) { return "feat: msg", nil }

	runner := review.NewRunner(&stubReviewer{})
	newRunner = func() *review.Runner { return runner }

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if err := run([]string{"--limit-break"}, &stdout, &stderr); err != nil {
		t.Fatalf("run() error = %v, want nil", err)
	}
	if runner.MaxDiffLines != diff.LimitBreakMaxLines {
		t.Fatalf("run() MaxDiffLines = %d, want %d", runner.MaxDiffLines, diff.LimitBreakMaxLines)
	}
}

func TestRunReturnsGitNotInstalled(t *testing.T) {
	restore := stubMainDeps()
	defer restore()

	checkGitInstalled = func() error { return git.ErrNotInstalled }

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if err := run(nil, &stdout, &stderr); !errors.Is(err, git.ErrNotInstalled) {
		t.Fatalf("run() error = %v, want ErrNotInstalled", err)
	}
}

func stubMainDeps() func() {
	originalCheckGitInstalled := checkGitInstalled
	originalLatestCommitDiff := latestCommitDiff
	originalLatestCommitMessage := latestCommitMessage
	originalNewRunner := newRunner

	return func() {
		checkGitInstalled = originalCheckGitInstalled
		latestCommitDiff = originalLatestCommitDiff
		latestCommitMessage = originalLatestCommitMessage
		newRunner = originalNewRunner
	}
}
