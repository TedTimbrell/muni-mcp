#!/bin/bash

# Run all tests with coverage
go test -cover ./...

# If you want to see detailed coverage report, uncomment the following lines
# go test -coverprofile=coverage.out ./...
# go tool cover -html=coverage.out -o coverage.html
# open coverage.html 