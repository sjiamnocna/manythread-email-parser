package pipeline

// FileReader produces work items (file paths) for the pipeline to process.
// It reads a list of file paths and yields them one-by-one to a channel,
// allowing downstream workers to begin processing immediately without waiting
// for the full list to be loaded into memory.
type FileReader struct {
	files []string
}

// NewFileReader creates a reader for the given list of file paths.
func NewFileReader(files []string) *FileReader {
	return &FileReader{
		files: files,
	}
}

// Read sends file jobs to the output channel as they're requested.
// This non-blocking design means files can be processed as they're yielded,
// rather than waiting for the entire batch to be loaded first.
func (fr *FileReader) Read() <-chan *fileJob {
	out := make(chan *fileJob)

	go func() {
		defer close(out)
		for _, filePath := range fr.files {
			out <- &fileJob{
				FilePath: filePath,
			}
		}
	}()

	return out
}
