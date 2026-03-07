package review

import (
	"fmt"
	"io"
	"strings"

	"review/pkg/copilot"
	"review/pkg/diff"
)

type Runner struct {
	Copilot *copilot.Client
	Verbose bool
}

func NewRunner(client *copilot.Client) *Runner {
	return &Runner{Copilot: client}
}

func BuildPrompt(files []diff.FileDiff, rawDiff string) string {
	var b strings.Builder
	b.WriteString("あなたは実践的なコードレビューアです。\n")
	b.WriteString("以下の最新コミット差分をレビューし、問題点と改善提案を日本語で簡潔に出力してください。\n")
	b.WriteString("出力形式:\n")
	b.WriteString("問題点\n- ...\n\n改善提案\n- ...\n\n")
	b.WriteString("対象ファイル一覧:\n")
	for _, f := range files {
		b.WriteString("- ")
		b.WriteString(f.FileName)
		b.WriteString("\n")
	}
	b.WriteString("\n差分:\n")
	b.WriteString(rawDiff)
	return b.String()
}

func (r *Runner) Run(rawDiff string, out, errOut io.Writer) error {
	if err := diff.ValidateSize(rawDiff); err != nil {
		return err
	}
	files := diff.Parse(rawDiff)
	if len(files) == 0 {
		_, _ = fmt.Fprintln(out, "No diff found")
		return nil
	}

	_, _ = fmt.Fprintln(out, "Analyzing diff...")
	prompt := BuildPrompt(files, rawDiff)
	if r.Verbose {
		_, _ = fmt.Fprintln(out, "Prompt:")
		_, _ = fmt.Fprintln(out, prompt)
	}
	return r.Copilot.Review(prompt, out, errOut)
}
