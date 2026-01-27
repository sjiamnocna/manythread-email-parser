package pipeline

type fileJob struct {
	FilePath string
}

type result struct {
	FileName  string
	LineCount int
	ByteCount int
	Err       error
}

// ResultMap is a map of filename to result (exported public API)
type ResultMap map[string]*result
