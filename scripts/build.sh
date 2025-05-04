#!/bin/bash

# Create output directory if it doesn't exist
mkdir -p build

# Build for macOS (both Intel and Apple Silicon)
GOOS=darwin GOARCH=amd64 go build -o build/muni-mcp-darwin-amd64 ./cmd/server
GOOS=darwin GOARCH=arm64 go build -o build/muni-mcp-darwin-arm64 ./cmd/server

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o build/muni-mcp-linux-amd64 ./cmd/server
GOOS=linux GOARCH=arm64 go build -o build/muni-mcp-linux-arm64 ./cmd/server

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o build/muni-mcp-windows-amd64.exe ./cmd/server

# Make the binaries executable
chmod +x build/muni-mcp-*

echo "Build complete! Binaries are available in the build directory:"
ls -l build/ 