#!/bin/bash

# Simple script to run the fishing game
# Makes it easier for users who don't have make installed

echo "Building and running the fishing game..."

# Ensure bin directory exists
mkdir -p bin

# Build the game
go build -o bin/fishing-game ./cmd/fishing

# Check if build was successful
if [ $? -eq 0 ]; then
    echo "Build successful, starting the game..."
    # Run the game
    ./bin/fishing-game
else
    echo "Build failed. Please check for errors and try again."
    exit 1
fi 