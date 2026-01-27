package email

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"path/filepath"
	"strings"
)

// Stats represents statistics for a single email file
type Stats struct {
	FileName  string
	LineCount int
	ByteCount int
	Err       error
}

// ProcessFile parses an email file and extracts text part statistics
func ProcessFile(filePath string) Stats {
	stats := Stats{
		FileName: filepath.Base(filePath),
	}

	file, err := openFile(filePath)
	if err != nil {
		stats.Err = fmt.Errorf("failed to open file: %w", err)
		return stats
	}
	defer file.Close()

	// Parse the email using net/mail
	msg, err := mail.ReadMessage(file)
	if err != nil {
		stats.Err = fmt.Errorf("failed to parse email: %w", err)
		return stats
	}

	// Extract text/plain part
	textContent, err := extractTextPart(msg)
	if err != nil {
		stats.Err = err
		return stats
	}

	// Calculate statistics
	stats.ByteCount = len(textContent)
	stats.LineCount = countLines(textContent)

	return stats
}

// extractTextPart extracts the text/plain part from an email message
func extractTextPart(msg *mail.Message) ([]byte, error) {
	contentType := msg.Header.Get("Content-Type")

	// Parse the content type
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		// If no content-type, assume it's plain text
		return io.ReadAll(msg.Body)
	}

	// Check if it's multipart
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

	// If it's already text/plain, return the body
	if mediaType == "text/plain" {
		return io.ReadAll(msg.Body)
	}

	return nil, fmt.Errorf("unsupported content type: %s", mediaType)
}

// countLines counts the number of lines in a byte slice
func countLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}

	count := 1 // Start with 1 since a non-empty file has at least one line
	for _, b := range data {
		if b == '\n' {
			count++
		}
	}

	// Adjust if the last character is a newline
	if data[len(data)-1] == '\n' {
		count--
	}

	return count
}
