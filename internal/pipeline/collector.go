package pipeline

import (
	"log"
)

// Collector gathers individual email analysis results into a single map.
// It also handles logging of any errors encountered during processing.
// Results are keyed by filename for easy lookup.
type Collector struct {
	logErrors bool
}

// NewCollector creates a new result collector.
// If logErrors is true, any processing errors will be logged as they're received.
func NewCollector(logErrors bool) *Collector {
	return &Collector{
		logErrors: logErrors,
	}
}

// Collect reads results from the input channel and aggregates them into a map.
// This blocks until all results have been received (the channel is closed).
// Returns a ResultMap with the filename as key and analysis result as value.
func (c *Collector) Collect(in <-chan *result) ResultMap {
	results := make(ResultMap)

	for res := range in {
		// Log errors as they occur (if configured to do so)
		if res.Err != nil && c.logErrors {
			log.Printf("Error processing %s: %v", res.FileName, res.Err)
		}
		// Store the result regardless of whether there was an error
		// (errors are captured in the result.Err field)
		results[res.FileName] = res
	}

	return results
}
