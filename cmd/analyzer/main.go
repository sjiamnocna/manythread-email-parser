package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"email-analyzer/internal/email"
	"email-analyzer/internal/output"
	"email-analyzer/internal/pipeline"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input_folder> <output_csv>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExample: %s './emails' 'results.csv'\n", os.Args[0])
		os.Exit(1)
	}

	inputDir := os.Args[1]
	outputCSV := os.Args[2]

	if info, err := os.Stat(inputDir); err != nil || !info.IsDir() {
		log.Fatalf("Invalid input directory: %s (directory must exist)", inputDir)
	}

	// Recursively scan the directory and find all .eml email files
	emailFiles, err := email.CollectFiles(inputDir)
	if err != nil {
		log.Fatalf("Failed to collect email files: %v", err)
	}

	if len(emailFiles) == 0 {
		log.Fatalf("No .eml email files found in: %s", inputDir)
	}

	log.Printf("Found %d email files to process", len(emailFiles))
	log.Printf("Using %d worker threads (one per CPU core)", runtime.NumCPU())

	results := pipeline.NewPipeline(emailFiles).
		SetNumberOfWorkers(runtime.NumCPU()).
		Execute()

	if err := output.WriteCSV(outputCSV, results); err != nil {
		log.Fatalf("Failed to write results to CSV file: %v", err)
	}

	successCount := 0
	for _, result := range results {
		if result.Err == nil {
			successCount++
		}
	}

	log.Printf("✓ Successfully processed %d/%d emails. Results saved to: %s",
		successCount, len(results), outputCSV)
}
