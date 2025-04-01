.PHONY: build run clean

# Default Go binary output
BINARY_NAME=fishing-game

# Default build directory
BUILD_DIR=./bin

build:
	@echo "Building fishing game..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/fishing
	@echo "Build complete! Binary is at $(BUILD_DIR)/$(BINARY_NAME)"

run: build
	@echo "Starting fishing game..."
	@$(BUILD_DIR)/$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Cleanup complete!"

test:
	@echo "Running tests..."
	@go test -v ./...

install-deps:
	@echo "Installing dependencies..."
	@go get -u github.com/charmbracelet/bubbletea
	@go get -u github.com/charmbracelet/lipgloss
	@echo "Dependencies installed!"

help:
	@echo "Fishing Game Makefile"
	@echo "---------------------"
	@echo "make build      - Build the fishing game"
	@echo "make run        - Build and run the fishing game"
	@echo "make clean      - Remove build artifacts"
	@echo "make test       - Run tests"
	@echo "make install-deps - Install dependencies"
	@echo "make help       - Show this help message" 