# Makefile

# Variables
BINARY_NAME=ecommerce
GO_CMD=go

# Default target
all: build

# Build the application
build:
	$(GO_CMD) build -o $(BINARY_NAME) .

# Run the application
run: build
	./$(BINARY_NAME)

# Clean the binary
clean:
	rm -f $(BINARY_NAME)

# Phony targets
.PHONY: all build run clean