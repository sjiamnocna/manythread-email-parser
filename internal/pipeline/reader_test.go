package pipeline

import "testing"

func TestFileReader_Read(t *testing.T) {
	files := []string{"file1.eml", "file2.eml", "file3.eml"}
	reader := NewFileReader(files)
	jobs := reader.Read()
	jobCount := 0
	for job := range jobs {
		if job == nil {
			t.Fatal("received nil job")
		}
		jobCount++
	}
	if jobCount != len(files) {
		t.Errorf("expected %d jobs, got %d", len(files), jobCount)
	}
}

func TestFileReader_Read_Empty(t *testing.T) {
	files := []string{}
	reader := NewFileReader(files)
	jobs := reader.Read()
	jobCount := 0
	for range jobs {
		jobCount++
	}
	if jobCount != 0 {
		t.Errorf("expected 0 jobs for empty input, got %d", jobCount)
	}
}
