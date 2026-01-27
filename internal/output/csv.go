package output

import (
	"encoding/csv"
	"fmt"
	"os"

	"email-analyzer/internal/pipeline"
)

// WriteCSV writes the email statistics to a CSV file from a ResultMap
func WriteCSV(outputPath string, results pipeline.ResultMap) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	if err := writer.Write([]string{"filename", "lines", "bytes"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, result := range results {
		if result.Err != nil {
			// Skip errors - they were already logged by collector
			continue
		}

		row := []string{
			result.FileName,
			fmt.Sprintf("%d", result.LineCount),
			fmt.Sprintf("%d", result.ByteCount),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	if err := writer.Error(); err != nil {
		return fmt.Errorf("CSV writer error: %w", err)
	}

	return nil
}
