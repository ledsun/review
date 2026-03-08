package review

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"strings"
	"text/template"

	"review/pkg/diff"
)

//go:embed prompt.tmpl
var promptTemplateText string

var promptTemplate = template.Must(template.New("prompt").Parse(promptTemplateText))

type Reviewer interface {
	Review(prompt string, stdout, stderr io.Writer) error
}

type Runner struct {
	Copilot      Reviewer
	Verbose      bool
	MaxDiffLines int
}

func NewRunner(client Reviewer) *Runner {
	return &Runner{
		Copilot:      client,
		MaxDiffLines: diff.DefaultMaxLines,
	}
}

type PromptData struct {
	CommitMessage string
	Files         []diff.FileDiff
	RawDiff       string
}

func BuildPrompt(commitMessage string, files []diff.FileDiff, rawDiff string) (string, error) {
	data := PromptData{
		CommitMessage: ensureTrailingNewline(commitMessage),
		Files:         files,
		RawDiff:       rawDiff,
	}

	var b bytes.Buffer
	if err := promptTemplate.Execute(&b, data); err != nil {
		return "", fmt.Errorf("build prompt: %w", err)
	}
	return strings.TrimSuffix(b.String(), "\n"), nil
}

func ensureTrailingNewline(s string) string {
	if strings.HasSuffix(s, "\n") {
		return s
	}
	return s + "\n"
}

func (r *Runner) Run(commitMessage, rawDiff string, out, errOut io.Writer) error {
	files := diff.Parse(rawDiff)
	if len(files) == 0 {
		_, _ = fmt.Fprintln(out, "No diff found")
		return nil
	}

	prompt, err := BuildPrompt(commitMessage, files, rawDiff)
	if err != nil {
		return err
	}
	if r.Verbose {
		_, _ = fmt.Fprintln(out, "Prompt:")
		_, _ = fmt.Fprintln(out, prompt)
	}
	if err := diff.ValidateSize(rawDiff, r.MaxDiffLines); err != nil {
		return err
	}
	_, _ = fmt.Fprintln(out, "Analyzing diff...")
	return r.Copilot.Review(prompt, out, errOut)
}
