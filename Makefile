TEST_DATA_DIR=test_data
RESULTS=results.csv
BINARY_NAME=analyzer
GO=go
GOFLAGS=-v

.PHONY: all build clean test run lint fmt ci

default: build
	./$(BINARY_NAME) "$(TEST_DATA_DIR)" $(RESULTS)

build: mod
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

clean:
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -f *.csv

test:
	$(GO) test $(GOFLAGS) ./...

run: mod
	$(GO) run ./cmd/$(BINARY_NAME) "$(TEST_DATA_DIR)" $(RESULTS)
	
mod:
	$(GO) mod download
	$(GO) mod tidy

lint:
	$(GO) fmt ./...
	$(GO) vet ./...

ci: lint build test
	@echo "CI checks passed: lint, build and test successful"
