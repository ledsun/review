package diff

import (
	"errors"
	"strings"
)

const MaxLines = 300

var ErrTooLarge = errors.New("Diff too large (max 300 lines)")

type Hunk struct {
	Header string
	Lines  []string
}

type FileDiff struct {
	FileName string
	RawLines []string
	Hunks    []Hunk
}

func CountLines(raw string) int {
	if raw == "" {
		return 0
	}
	trimmed := strings.TrimSuffix(raw, "\n")
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(trimmed, "\n"))
}

func ValidateSize(raw string) error {
	if CountLines(raw) > MaxLines {
		return ErrTooLarge
	}
	return nil
}

func Parse(raw string) []FileDiff {
	lines := strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n")
	files := make([]FileDiff, 0)
	var current *FileDiff
	var currentHunk *Hunk

	flushHunk := func() {
		if current != nil && currentHunk != nil {
			current.Hunks = append(current.Hunks, *currentHunk)
			currentHunk = nil
		}
	}

	flushFile := func() {
		if current != nil {
			flushHunk()
			files = append(files, *current)
			current = nil
		}
	}

	for _, l := range lines {
		if strings.HasPrefix(l, "diff --git ") {
			flushFile()
			current = &FileDiff{}
			parts := strings.Fields(l)
			if len(parts) >= 4 {
				current.FileName = strings.TrimPrefix(parts[3], "b/")
			}
			current.RawLines = append(current.RawLines, l)
			continue
		}
		if current == nil {
			continue
		}
		current.RawLines = append(current.RawLines, l)
		if strings.HasPrefix(l, "@@") {
			flushHunk()
			currentHunk = &Hunk{Header: l}
			continue
		}
		if currentHunk != nil {
			currentHunk.Lines = append(currentHunk.Lines, l)
		}
	}
	flushFile()
	return files
}
