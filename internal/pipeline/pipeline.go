package pipeline

import (
	"runtime"
)

// Pipeline orchestrates the email processing workflow.
// It uses a three-stage design pattern for efficient concurrent processing:
// - Reader stage: Generates file paths to process
// - Parser stage: Concurrently analyzes email files
// - Collector stage: Aggregates results into a single map
//
// The pipeline supports fluent interface method chaining for easy configuration.
type Pipeline struct {
	reader    *FileReader
	parser    *Parser
	collector *Collector
}

// NewPipeline creates a new email processing pipeline with sensible defaults.
func NewPipeline(files []string) *Pipeline {
	return &Pipeline{
		reader:    NewFileReader(files),
		parser:    NewParser(runtime.NumCPU()),
		collector: NewCollector(true),
	}
}

// SetNumberOfWorkers adjusts the number of concurrent workers processing emails.
func (p *Pipeline) SetNumberOfWorkers(n int) *Pipeline {
	p.parser = NewParser(n)
	return p
}

// WithoutErrorLogging disables logging of per-file errors.
// Errors are still collected and available in results, just not logged.
// Returns the pipeline for method chaining.
func (p *Pipeline) WithoutErrorLogging() *Pipeline {
	p.collector = NewCollector(false)
	return p
}

// Execute runs the complete pipeline from start to finish.
// It processes all emails and returns results in a map keyed by filename.
func (p *Pipeline) Execute() ResultMap {
	fileJobs := p.reader.Read()
	results := p.parser.Parse(fileJobs)
	return p.collector.Collect(results)
}
