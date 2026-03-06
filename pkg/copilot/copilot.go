package copilot

import (
	"fmt"
	"io"
	"os/exec"
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

	cmd := exec.Command("copilot", "--model", c.Model, "-p", prompt)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("copilot execution failed: %w", err)
	}
	return nil
}
