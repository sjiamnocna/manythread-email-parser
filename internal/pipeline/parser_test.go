package pipeline

import "testing"

func TestCountLines_Empty(t *testing.T) {
	result := countLines([]byte{})
	if result != 0 {
		t.Errorf("expected 0 lines for empty input, got %d", result)
	}
}

func TestCountLines_SingleLine(t *testing.T) {
	result := countLines([]byte("hello world"))
	if result != 1 {
		t.Errorf("expected 1 line, got %d", result)
	}
}

func TestCountLines_MultipleLines(t *testing.T) {
	result := countLines([]byte("line1\nline2\nline3"))
	if result != 3 {
		t.Errorf("expected 3 lines, got %d", result)
	}
}

func TestCountLines_WithTrailingNewline(t *testing.T) {
	result := countLines([]byte("line1\nline2\n"))
	if result != 2 {
		t.Errorf("expected 2 lines with trailing newline, got %d", result)
	}
}

func TestNewParser(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 1},
		{-5, 1},
		{1, 1},
		{4, 4},
		{100, 100},
	}
	for _, tt := range tests {
		parser := NewParser(tt.input)
		if parser.numWorkers != tt.expected {
			t.Errorf("NewParser(%d): expected %d workers, got %d", tt.input, tt.expected, parser.numWorkers)
		}
	}
}
