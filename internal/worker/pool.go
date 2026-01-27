package worker

import (
	"sync"

	"email-analyzer/internal/email"
)

// Pool manages concurrent email processing
type Pool struct {
	NumWorkers int
}

// New creates a new worker pool
func New(numWorkers int) *Pool {
	if numWorkers < 1 {
		numWorkers = WorkersNumberDefault
	}
	return &Pool{
		NumWorkers: numWorkers,
	}
}

// Process processes email files concurrently by sending them one by one through the jobs channel and returns results
func (p *Pool) Process(emailFiles []string) []email.Stats {
	jobs := make(chan string, len(emailFiles))
	results := make(chan email.Stats, len(emailFiles))
	var wg sync.WaitGroup

	for i := 0; i < p.NumWorkers; i++ {
		wg.Add(1)
		go p.worker(jobs, results, &wg)
	}

	for _, filePath := range emailFiles {
		jobs <- filePath
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var allStats []email.Stats
	for stat := range results {
		allStats = append(allStats, stat)
	}

	return allStats
}

// worker processes email files from the jobs channel
func (p *Pool) worker(jobs <-chan string, results chan<- email.Stats, wg *sync.WaitGroup) {
	defer wg.Done()
	for filePath := range jobs {
		results <- email.ProcessFile(filePath)
	}
}
