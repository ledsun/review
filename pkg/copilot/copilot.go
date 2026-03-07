package copilot

import (
	"context"
	"fmt"
	"io"

	sdk "github.com/github/copilot-sdk/go"
)

type Client struct {
	Model string
}

func New(model string) *Client {
	return &Client{Model: model}
}

func (c *Client) Review(prompt string, stdout, stderr io.Writer) error {
	_ = stderr

	client := sdk.NewClient(&sdk.ClientOptions{
		Cwd:      ".",
		LogLevel: "error",
	})

	ctx := context.Background()
	if err := client.Start(ctx); err != nil {
		return fmt.Errorf("copilot sdk start failed: %w", err)
	}
	defer client.Stop()

	session, err := client.CreateSession(ctx, &sdk.SessionConfig{
		Model:               c.Model,
		Streaming:           true,
		OnPermissionRequest: sdk.PermissionHandler.ApproveAll,
	})
	if err != nil {
		return fmt.Errorf("copilot session creation failed: %w", err)
	}
	defer session.Destroy()

	streamed := false
	trailingLine := true
	errCh := make(chan error, 1)
	doneCh := make(chan struct{}, 1)

	unsubscribe := session.On(func(event sdk.SessionEvent) {
		switch event.Type {
		case sdk.AssistantMessageDelta:
			if event.Data.DeltaContent == nil {
				return
			}
			streamed = true
			text := *event.Data.DeltaContent
			_, _ = fmt.Fprint(stdout, text)
			trailingLine = len(text) == 0 || text[len(text)-1] == '\n'
		case sdk.AssistantMessage:
			if streamed || event.Data.Content == nil {
				return
			}
			text := *event.Data.Content
			_, _ = fmt.Fprint(stdout, text)
			trailingLine = len(text) == 0 || text[len(text)-1] == '\n'
		case sdk.SessionError:
			msg := "session error"
			if event.Data.Message != nil {
				msg = *event.Data.Message
			}
			select {
			case errCh <- fmt.Errorf("copilot session error: %s", msg):
			default:
			}
		case sdk.SessionIdle:
			if !trailingLine {
				_, _ = fmt.Fprintln(stdout)
			}
			select {
			case doneCh <- struct{}{}:
			default:
			}
		}
	})
	defer unsubscribe()

	if _, err := session.Send(ctx, sdk.MessageOptions{Prompt: prompt}); err != nil {
		return fmt.Errorf("copilot send failed: %w", err)
	}

	select {
	case err := <-errCh:
		return err
	case <-doneCh:
		return nil
	}
}
