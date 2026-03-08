package diff

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateSizeWithinLimit(t *testing.T) {
	raw := strings.Repeat("+x\n", DefaultMaxLines)

	if err := ValidateSize(raw, DefaultMaxLines); err != nil {
		t.Fatalf("ValidateSize() error = %v, want nil", err)
	}
}

func TestValidateSizeTooLarge(t *testing.T) {
	raw := strings.Repeat("+x\n", LimitBreakMaxLines+1)

	err := ValidateSize(raw, LimitBreakMaxLines)
	if err == nil {
		t.Fatal("ValidateSize() error = nil, want error")
	}
	if !errors.Is(err, ErrTooLarge) {
		t.Fatalf("ValidateSize() error = %v, want ErrTooLarge", err)
	}

	var tooLarge *TooLargeError
	if !errors.As(err, &tooLarge) {
		t.Fatalf("ValidateSize() error = %T, want *TooLargeError", err)
	}
	if tooLarge.Lines != LimitBreakMaxLines+1 {
		t.Fatalf("TooLargeError.Lines = %d, want %d", tooLarge.Lines, LimitBreakMaxLines+1)
	}
	if tooLarge.MaxLines != LimitBreakMaxLines {
		t.Fatalf("TooLargeError.MaxLines = %d, want %d", tooLarge.MaxLines, LimitBreakMaxLines)
	}
}
