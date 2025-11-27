#!/bin/bash

# Build script for IaC Security Scanner
# Builds standalone binaries for multiple platforms

set -e

VERSION=${1:-"1.0.0"}
BUILD_DIR="dist"
BINARY_NAME="iac-audit"

echo "🔨 Building IaC Security Scanner v${VERSION}"
echo "=========================================="

# Create dist directory
mkdir -p ${BUILD_DIR}

# Build for Linux (amd64)
echo "📦 Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 .

# Build for Linux (arm64)
echo "📦 Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.version=${VERSION}" -o ${BUILD_DIR}/${BINARY_NAME}-linux-arm64 .

# Build for macOS (amd64)
echo "📦 Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64 .

# Build for macOS (arm64 - Apple Silicon)
echo "📦 Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.version=${VERSION}" -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64 .

# Build for Windows (amd64)
echo "📦 Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe .

echo ""
echo "✅ Build complete!"
echo "📁 Binaries are in: ${BUILD_DIR}/"
echo ""
echo "Files created:"
ls -lh ${BUILD_DIR}/

