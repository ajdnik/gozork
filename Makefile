.PHONY: build run test vet lint clean fmt

# Build the game binary
build:
	go build -o gozork .

# Run the game
run: build
	./gozork

# Run all tests (verbose, no cache)
test:
	go test -v -count=1 -timeout 300s ./...

# Run go vet
vet:
	go vet ./...

# Run golangci-lint (install: https://golangci-lint.run/welcome/install/)
lint:
	golangci-lint run ./...

# Format all Go source files
fmt:
	gofmt -w .

# Remove build artifacts
clean:
	rm -f gozork

# Run all checks (format, vet, lint, test)
check: fmt vet lint test
