.PHONY: build run test cover vet lint clean fmt

# Build the game binary
build:
	go build -o gozork .

# Run the game
run: build
	./gozork

# Run all tests (verbose, no cache)
test:
	go test -v -count=1 -timeout 300s ./...

# Run tests with coverage summary
cover:
	@go test -count=1 -timeout 300s -coverprofile=coverage.out ./... 2>&1 | \
		awk '/^ok/ { split($$0,a,"coverage: "); split(a[2],b," "); printf "  %-50s %s\n", $$2, b[1] } \
		     /^[^o]/ && /coverage:/ { for(i=1;i<=NF;i++) if($$i ~ /^github/) pkg=$$i; split($$0,a,"coverage: "); split(a[2],b," "); printf "  %-50s %s (no tests)\n", pkg, b[1] }'
	@echo ""
	@go tool cover -func=coverage.out | awk '/^total:/ { printf "Total coverage: %s\n", $$NF }'
	@rm -f coverage.out

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
	rm -f gozork coverage.out

# Run all checks (format, vet, lint, test)
check: fmt vet lint test
