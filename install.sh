#!/bin/bash

# Installation script for the fishing game
# Sets up dependencies and ensures the game can run

echo "Setting up the fishing game..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed!"
    echo "Please install Go 1.18 or higher from https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | grep -oE '[0-9]+\.[0-9]+')
GO_VERSION_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
GO_VERSION_MINOR=$(echo $GO_VERSION | cut -d. -f2)

if [ "$GO_VERSION_MAJOR" -lt 1 ] || ([ "$GO_VERSION_MAJOR" -eq 1 ] && [ "$GO_VERSION_MINOR" -lt 18 ]); then
    echo "Error: Go version 1.18 or higher is required!"
    echo "Current version: $GO_VERSION"
    echo "Please upgrade your Go installation."
    exit 1
fi

echo "Go version $GO_VERSION detected. ✓"

# Install dependencies
echo "Installing dependencies..."
go get -u github.com/charmbracelet/bubbletea
go get -u github.com/charmbracelet/lipgloss

# Ensure the project builds
echo "Building the game to verify setup..."
mkdir -p bin
go build -o bin/fishing-game ./cmd/fishing

if [ $? -eq 0 ]; then
    echo "======================="
    echo "Setup complete! ✓"
    echo "======================="
    echo "You can now run the game using:"
    echo "  ./run.sh"
    echo "  or"
    echo "  make run (if you have make installed)"
    echo ""
    echo "Enjoy fishing!"
else
    echo "Build failed. Please check the error messages above."
    exit 1
fi 