# Email Analyzer

A high-performance, multithreaded Go application that analyzes email files and generates CSV statistics.

This tool processes email files (.eml format) from a directory, extracts the text content from each email, and generates a CSV report with statistics about lines and bytes in the text portion of each email.

## Run

> Prepare test data in the `*.eml` format inside the `test_data` directory or provide your own folder with `.eml` files as argument or modify prepared Makefile variable.

> Test run using Makefile, it alrady includes the `go mod download` step.

```bash
make run
```

### Default target builds and runs the analyzer on the included test data:

```bash
make
```

## Build

```bash
# Simple build
go build -o analyzer ./cmd/analyzer

# Or simple using make
make build
```

### Usage

```bash
./email-analyzer <input_folder> <output_csv>
```

### Examples

```bash
# Process emails from a folder
./email-analyzer ./emails results.csv

# Process with absolute paths
./email-analyzer /home/user/emails /tmp/report.csv

# Process test emails from the included test data
./email-analyzer "test_data" results.csv
```

## Output Format

The generated CSV file contains the following columns:

| Column | Description |
|--------|-------------|
| `filename` | Name of the email file |
| `lines` | Number of lines in the text/plain part of the email |
| `bytes` | Number of bytes in the text/plain part of the email |

### Example Output

```csv
filename,lines,bytes
test_email_001.eml,42,1320
test_email_002.eml,67,2105
test_email_003.eml,31,950
```

### Processing Notes

- Only the `text/plain` part of emails is analyzed (HTML parts are ignored)
- Empty lines and trailing newlines are handled correctly
- If an email has no text/plain part, it's logged as an error but processing continues
- Lines are counted by newline characters; an email without trailing newline counts as one line

## How It Works

The analyzer uses a three-stage **pipeline pattern** for efficient concurrent processing:

```
Reader → Parser → Collector
```

1. **Reader Stage**: Reads the list of email file paths and sends them to a channel
2. **Parser Stage**: Multiple workers (one per CPU core) concurrently read and analyze each email file
3. **Collector Stage**: Gathers all results into a map for easy access

This design allows the system to:
- Start processing files immediately without waiting for the full file list to be scanned
- Efficiently use all available CPU cores for parallel processing
- Handle errors gracefully while continuing with other files
- Keep memory usage constant regardless of the number of files

## Development

### Building

```bash
make build          # Build the binary
```

### Testing

```bash
make test           # Run all unit tests
go test ./...       # Alternative command
go test ./... -v    # Verbose output
go test -cover ./...  # With coverage report
```

### Project Structure

```plaintext
.
├── cmd/analyzer/        # Main application entry point
│   └── main.go          # CLI argument parsing and orchestration
├── internal/
│   ├── pipeline/        # Core pipeline stages
│   │   ├── pipeline.go  # Pipeline orchestration with builder pattern
│   │   ├── reader.go    # File path generation stage
│   │   ├── parser.go    # Email analysis and line counting
│   │   ├── collector.go # Result aggregation
│   │   ├── types.go     # Internal and public types
│   │   └── *_test.go    # Unit tests
│   ├── email/           # Email file operations
│   │   ├── file.go      # Directory scanning for .eml files
│   │   └── *_test.go    # Tests
│   └── output/          # Output formatting
│       ├── csv.go       # CSV file generation
│       └── *_test.go    # Tests
├── test_data/           # Sample email data for testing
├── Makefile             # Build automation
├── README.md            # This file
└── TESTS.md             # Detailed test documentation
```
