package copilot

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type Client struct {
	Model string
}

func New(model string) *Client {
	return &Client{Model: model}
}

func (c *Client) Review(prompt string, stdout, stderr io.Writer) error {
	if _, err := exec.LookPath("copilot"); err != nil {
		return fmt.Errorf("copilot CLI not found: %w", err)
	}

	promptPath, err := writePrompt(prompt)
	if err != nil {
		return err
	}
	defer os.RemoveAll(filepath.Dir(promptPath))

	cmd := exec.Command("copilot", "--model", c.Model, "-p", "@"+promptPath)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("copilot execution failed: %w", err)
	}
	return nil
}

func writePrompt(prompt string) (string, error) {
	dir, err := os.MkdirTemp("", "review-copilot-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	path := filepath.Join(dir, "prompt.txt")
	if err := os.WriteFile(path, []byte(prompt), 0o600); err != nil {
		_ = os.RemoveAll(dir)
		return "", fmt.Errorf("failed to write prompt file: %w", err)
	}
	return path, nil
}
