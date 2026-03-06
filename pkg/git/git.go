package git

import (
	"errors"
	"fmt"
	"os/exec"
)

var ErrNotInstalled = errors.New("git is not installed")

func CheckInstalled() error {
	if _, err := exec.LookPath("git"); err != nil {
		return ErrNotInstalled
	}
	return nil
}

// LatestCommitDiff returns the full-function context diff for the latest commit.
func LatestCommitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "-W", "-U3", "HEAD~1", "HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get latest commit diff: %w: %s", err, string(out))
	}
	return string(out), nil
}
