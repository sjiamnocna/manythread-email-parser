package pipeline

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Parser struct {
	numWorkers int
}

func NewParser(numWorkers int) *Parser {
	if numWorkers < 1 {
		numWorkers = 1
	}
	return &Parser{
		numWorkers: numWorkers,
	}
}

func (p *Parser) Parse(in <-chan *fileJob) <-chan *result {
	out := make(chan *result)

	var wg sync.WaitGroup
	for i := 0; i < p.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range in {
				result := p.processFile(job.FilePath)
				out <- result
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func (p *Parser) processFile(filePath string) *result {
	res := &result{
		FileName: filepath.Base(filePath),
	}

	file, err := os.Open(filePath)
	if err != nil {
		res.Err = fmt.Errorf("failed to open file: %w", err)
		return res
	}
	defer file.Close()

	msg, err := mail.ReadMessage(file)
	if err != nil {
		res.Err = fmt.Errorf("failed to parse email: %w", err)
		return res
	}

	textContent, err := extractTextPart(msg)
	if err != nil {
		res.Err = err
		return res
	}

	res.ByteCount = len(textContent)
	res.LineCount = countLines(textContent)

	return res
}

func extractTextPart(msg *mail.Message) ([]byte, error) {
	contentType := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return io.ReadAll(msg.Body)
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		boundary := params["boundary"]
		if boundary == "" {
			return nil, fmt.Errorf("multipart message missing boundary")
		}

		mr := multipart.NewReader(msg.Body, boundary)
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("failed to read multipart: %w", err)
			}

			partContentType := part.Header.Get("Content-Type")
			partMediaType, _, _ := mime.ParseMediaType(partContentType)

			if partMediaType == "text/plain" {
				content, err := io.ReadAll(part)
				if err != nil {
					return nil, fmt.Errorf("failed to read text part: %w", err)
				}
				return content, nil
			}
		}

		return nil, fmt.Errorf("no text/plain part found in email")
	}

	if mediaType == "text/plain" {
		return io.ReadAll(msg.Body)
	}

	return nil, fmt.Errorf("unsupported content type: %s", mediaType)
}

func countLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}

	count := 1
	for _, b := range data {
		if b == '\n' {
			count++
		}
	}

	if data[len(data)-1] == '\n' {
		count--
	}

	return count
}
