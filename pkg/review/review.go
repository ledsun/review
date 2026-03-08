package review

import (
	"fmt"
	"io"
	"strings"

	"review/pkg/diff"
)

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

func BuildPrompt(commitMessage string, files []diff.FileDiff, rawDiff string) string {
	var b strings.Builder
	b.WriteString("あなたは実践的なコードレビューをするお嬢様「皇戸麗風子」です。\n")
	b.WriteString("以下の最新コミットのコミットメッセージ､ファイル一覧、差分を見て、コミットメッセージと修正内容に相違はないか？タイプミスはないか？確認してください。\n")
	b.WriteString("静的に解析してください。CLI、build、format コマンドは実行しないでください。\n")
	b.WriteString("語尾は以下を参考にしてください。\n")
	b.WriteString("- ですの\n")
	b.WriteString("- ございますわ\n")
	b.WriteString("- いますの\n")
	b.WriteString("- いませんの\n")
	b.WriteString("- でしょうか\n")
	b.WriteString("- かしら\n\n")
	b.WriteString("出力形式:\n")
	b.WriteString("メッセージの齟齬\n- ...\n\nタイプミス\n- ...\n\n以上ですの。\n\n")
	b.WriteString("コミットメッセージ:\n")
	b.WriteString(commitMessage)
	if !strings.HasSuffix(commitMessage, "\n") {
		b.WriteString("\n")
	}
	b.WriteString("ファイル一覧:\n")
	for _, f := range files {
		b.WriteString("- ")
		b.WriteString(f.FileName)
		b.WriteString("\n")
	}
	b.WriteString("\n差分:\n")
	b.WriteString(rawDiff)
	return b.String()
}

func (r *Runner) Run(commitMessage, rawDiff string, out, errOut io.Writer) error {
	files := diff.Parse(rawDiff)
	if len(files) == 0 {
		_, _ = fmt.Fprintln(out, "No diff found")
		return nil
	}

	prompt := BuildPrompt(commitMessage, files, rawDiff)
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
