# Define variables
GO_CMD := go
BINARY_NAME := myapp
BUILD_DIR := bin
SRC_DIR := .

# Default target
.PHONY: all
all: build

# Build target
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO_CMD) build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)

# Run target
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME)

# Test target
.PHONY: test
test:
	@echo "Running tests..."
	$(GO_CMD) test ./... -v

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
