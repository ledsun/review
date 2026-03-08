package diff

import (
	"errors"
	"fmt"
	"strings"
)

const (
	DefaultMaxLines    = 300
	LimitBreakMaxLines = 3000
)

var ErrTooLarge = errors.New("diff too large")

type TooLargeError struct {
	Lines    int
	MaxLines int
}

func (e *TooLargeError) Error() string {
	return fmt.Sprintf("Diff too large (%d lines, max %d lines)", e.Lines, e.MaxLines)
}

func (e *TooLargeError) Is(target error) bool {
	return target == ErrTooLarge
}

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

func ValidateSize(raw string, maxLines int) error {
	lines := CountLines(raw)
	if lines > maxLines {
		return &TooLargeError{Lines: lines, MaxLines: maxLines}
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
