package copilot

import "testing"

func TestNewSetsModel(t *testing.T) {
	client := New("gpt-5-mini")

	if client.Model != "gpt-5-mini" {
		t.Fatalf("New().Model = %q, want %q", client.Model, "gpt-5-mini")
	}
}
