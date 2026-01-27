package output

import (
	"os"
	"strings"
	"testing"

	"email-analyzer/internal/pipeline"
)

func TestWriteCSV_Basic(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	results := make(pipeline.ResultMap)

	err = WriteCSV(tmpFile.Name(), results)
	if err != nil {
		t.Fatalf("WriteCSV failed: %v", err)
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read CSV file: %v", err)
	}

	csvContent := string(content)
	if !strings.Contains(csvContent, "filename,lines,bytes") {
		t.Errorf("CSV header not found in output")
	}
}

func TestWriteCSV_Empty(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	results := make(pipeline.ResultMap)

	err = WriteCSV(tmpFile.Name(), results)
	if err != nil {
		t.Fatalf("WriteCSV failed: %v", err)
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read CSV file: %v", err)
	}

	csvContent := string(content)
	lines := strings.Split(strings.TrimSpace(csvContent), "\n")

	if len(lines) != 1 || !strings.Contains(lines[0], "filename") {
		t.Errorf("expected only header line, got %d lines", len(lines))
	}
}
