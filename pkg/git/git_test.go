package git

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestCheckInstalledReturnsErrNotInstalled(t *testing.T) {
	originalLookPath := lookPath
	t.Cleanup(func() {
		lookPath = originalLookPath
	})

	lookPath = func(file string) (string, error) {
		if file != "git" {
			t.Fatalf("lookPath() file = %q, want git", file)
		}
		return "", errors.New("missing")
	}

	if err := CheckInstalled(); !errors.Is(err, ErrNotInstalled) {
		t.Fatalf("CheckInstalled() error = %v, want ErrNotInstalled", err)
	}
}

func TestLatestCommitDiffRunsExpectedCommand(t *testing.T) {
	originalRunGit := runGit
	t.Cleanup(func() {
		runGit = originalRunGit
	})

	runGit = func(args ...string) ([]byte, error) {
		got := strings.Join(args, " ")
		want := "diff -W -U3 HEAD~1 HEAD"
		if got != want {
			t.Fatalf("runGit() args = %q, want %q", got, want)
		}
		return []byte("diff output"), nil
	}

	got, err := LatestCommitDiff()
	if err != nil {
		t.Fatalf("LatestCommitDiff() error = %v, want nil", err)
	}
	if got != "diff output" {
		t.Fatalf("LatestCommitDiff() = %q, want diff output", got)
	}
}

func TestLatestCommitMessageIncludesCommandOutputInError(t *testing.T) {
	originalRunGit := runGit
	t.Cleanup(func() {
		runGit = originalRunGit
	})

	runGit = func(args ...string) ([]byte, error) {
		return []byte("fatal: bad revision"), exec.ErrNotFound
	}

	_, err := LatestCommitMessage()
	if err == nil {
		t.Fatal("LatestCommitMessage() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "failed to get latest commit message") {
		t.Fatalf("LatestCommitMessage() error = %q, want command context", err)
	}
	if !strings.Contains(err.Error(), "fatal: bad revision") {
		t.Fatalf("LatestCommitMessage() error = %q, want command output", err)
	}
}
