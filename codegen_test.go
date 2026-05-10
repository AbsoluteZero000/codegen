package main

import (
	"testing"

	openrouter "github.com/revrost/go-openrouter"
)

func TestParseToolCall_valid(t *testing.T) {
	content := `{"tool": "readFile", "arguments": {"filename": "main.go"}}`
	tc := parseToolCall(content)
	if tc == nil {
		t.Fatal("expected tool call, got nil")
	}
	if tc.Tool != "readFile" {
		t.Fatalf("expected readFile, got %s", tc.Tool)
	}
	if tc.Arguments.Filename != "main.go" {
		t.Fatalf("expected main.go, got %s", tc.Arguments.Filename)
	}
}

func TestParseToolCall_embedded(t *testing.T) {
	content := `some text {"tool": "bash", "arguments": {"command": "ls"}} trailing`
	tc := parseToolCall(content)
	if tc == nil {
		t.Fatal("expected tool call, got nil")
	}
	if tc.Tool != "bash" {
		t.Fatalf("expected bash, got %s", tc.Tool)
	}
}

func TestParseToolCall_invalid(t *testing.T) {
	tests := []string{
		"",
		"just some text",
		`{"tool": "", "arguments": {}}`,
	}
	for _, c := range tests {
		tc := parseToolCall(c)
		if tc != nil {
			t.Fatalf("expected nil for %q, got %+v", c, tc)
		}
	}
}

func TestSanitize_dangerous(t *testing.T) {
	tests := []string{
		"rm -rf /",
		"rm -rf /*",
		"dd if=/dev/zero of=/dev/sda",
		"mkfs.ext4 /dev/sda",
	}
	for _, c := range tests {
		result := sanitize(c)
		if result != "" {
			t.Fatalf("expected empty for %q, got %q", c, result)
		}
	}
}

func TestSanitize_safe(t *testing.T) {
	tests := []string{
		"ls -la",
		"go build .",
		"echo hello",
	}
	for _, c := range tests {
		result := sanitize(c)
		if result == "" {
			t.Fatalf("expected non-empty for %q", c)
		}
	}
}

func TestSanitize_trim(t *testing.T) {
	result := sanitize("  ls  ")
	if result != "ls" {
		t.Fatalf("expected ls, got %q", result)
	}
}

func TestIsFree(t *testing.T) {
	tests := []struct {
		model openrouter.Model
		free  bool
	}{
		{openrouter.Model{ID: "model:free", Pricing: openrouter.ModelPricing{Prompt: "0", Completion: "0"}}, true},
		{openrouter.Model{ID: "model:free", Pricing: openrouter.ModelPricing{Prompt: "1", Completion: "1"}}, true},
		{openrouter.Model{ID: "model", Pricing: openrouter.ModelPricing{Prompt: "0", Completion: "0"}}, true},
		{openrouter.Model{ID: "model", Pricing: openrouter.ModelPricing{Prompt: "1", Completion: "0.5"}}, false},
	}
	for _, tt := range tests {
		got := isFree(tt.model)
		if got != tt.free {
			t.Errorf("isFree(%+v) = %v, want %v", tt.model, got, tt.free)
		}
	}
}

func TestNewToolCallID(t *testing.T) {
	id1 := newToolCallID()
	id2 := newToolCallID()
	if id1 == id2 {
		t.Fatal("expected unique IDs")
	}
	if len(id1) < 10 {
		t.Fatalf("expected reasonably long ID, got %q", id1)
	}
}

func TestLsDir(t *testing.T) {
	result := lsDir(t.TempDir())
	if result != "" {
		t.Fatal("expected empty for empty dir")
	}
}

func TestLsDir_default(t *testing.T) {
	result := lsDir("")
	if result == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestGlob(t *testing.T) {
	dir := t.TempDir()
	result := glob(dir + "/*")
	if result != "(no matches)" {
		t.Fatalf("expected no matches, got %q", result)
	}
}

func TestGlob_invalid(t *testing.T) {
	result := glob("[")
	if result == "" {
		t.Fatal("expected error message")
	}
}
