package git

import (
	"errors"
	"fmt"
	"os/exec"
)

var ErrNotInstalled = errors.New("git is not installed")

var (
	lookPath = exec.LookPath
	runGit   = func(args ...string) ([]byte, error) {
		return exec.Command("git", args...).CombinedOutput()
	}
)

func CheckInstalled() error {
	if _, err := lookPath("git"); err != nil {
		return ErrNotInstalled
	}
	return nil
}

// LatestCommitDiff returns the full-function context diff for the latest commit.
func LatestCommitDiff() (string, error) {
	out, err := runGit("diff", "-W", "-U3", "HEAD~1", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to get latest commit diff: %w: %s", err, string(out))
	}
	return string(out), nil
}

func LatestCommitMessage() (string, error) {
	out, err := runGit("log", "-1", "--pretty=%B")
	if err != nil {
		return "", fmt.Errorf("failed to get latest commit message: %w: %s", err, string(out))
	}
	return string(out), nil
}
