package pipeline

import (
	"runtime"
	"testing"
)

func TestNewPipeline(t *testing.T) {
	files := []string{"file1.eml", "file2.eml"}
	pipeline := NewPipeline(files)
	if pipeline.reader == nil {
		t.Error("expected reader to be initialized")
	}
	if pipeline.parser == nil {
		t.Error("expected parser to be initialized")
	}
	if pipeline.collector == nil {
		t.Error("expected collector to be initialized")
	}
	if pipeline.parser.numWorkers != runtime.NumCPU() {
		t.Errorf("expected %d workers, got %d", runtime.NumCPU(), pipeline.parser.numWorkers)
	}
}

func TestPipeline_SetNumberOfWorkers(t *testing.T) {
	files := []string{"file1.eml"}
	pipeline := NewPipeline(files)
	updatedPipeline := pipeline.SetNumberOfWorkers(8)
	if updatedPipeline != pipeline {
		t.Error("expected SetNumberOfWorkers to return the same pipeline for chaining")
	}
	if pipeline.parser.numWorkers != 8 {
		t.Errorf("expected 8 workers, got %d", pipeline.parser.numWorkers)
	}
}

func TestPipeline_WithoutErrorLogging(t *testing.T) {
	files := []string{"file1.eml"}
	pipeline := NewPipeline(files)
	updatedPipeline := pipeline.WithoutErrorLogging()
	if updatedPipeline != pipeline {
		t.Error("expected WithoutErrorLogging to return the same pipeline for chaining")
	}
	if pipeline.collector.logErrors {
		t.Error("expected logErrors to be false after WithoutErrorLogging()")
	}
}

func TestPipeline_Chaining(t *testing.T) {
	files := []string{"file1.eml", "file2.eml"}
	pipeline := NewPipeline(files).SetNumberOfWorkers(2).WithoutErrorLogging()
	if pipeline.parser.numWorkers != 2 {
		t.Errorf("chaining: expected 2 workers, got %d", pipeline.parser.numWorkers)
	}
	if pipeline.collector.logErrors {
		t.Error("chaining: expected logErrors to be false")
	}
}
