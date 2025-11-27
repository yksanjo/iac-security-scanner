.PHONY: build clean test install run-example

# Build the binary
build:
	@echo "Building iac-audit..."
	@go build -o iac-audit .
	@echo "✅ Build complete: ./iac-audit"

# Build for multiple platforms
build-all:
	@echo "Building for all platforms..."
	@GOOS=linux GOARCH=amd64 go build -o dist/iac-audit-linux-amd64 .
	@GOOS=darwin GOARCH=amd64 go build -o dist/iac-audit-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -o dist/iac-audit-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 go build -o dist/iac-audit-windows-amd64.exe .
	@echo "✅ Builds complete in ./dist/"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f iac-audit iac-audit.exe
	@rm -rf dist/
	@rm -f *.json *.pdf security-report*
	@echo "✅ Clean complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Run example scan
run-example:
	@echo "Running example scan..."
	@./iac-audit scan ./examples --format=json --output=example-report

# Format code
fmt:
	@go fmt ./...

# Lint code
lint:
	@golangci-lint run || echo "Install golangci-lint: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

