# Project variables
BINARY_NAME=vps-setup
BUILD_DIR=bin
VERSION=1.0.0

# Build configuration
# CGO_ENABLED=0 ensures a static binary with no dependencies on host C libraries
GOOS=linux
GOARCH=amd64
CGO_ENABLED=0

# LDFLAGS: 
# -s: Omit the symbol table and debug information
# -w: Omit the DWARF symbol table
# This makes the binary smaller for faster uploading to your VPS
LDFLAGS=-ldflags="-s -w"

all: build

## build: Compiles the project for the target VPS architecture
build:
	@echo "Building static binary for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## clean: Removes the build directory
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

## deploy: Builds and then attempts to SCP the binary to a test VPS
# Usage: make deploy USER=root IP=1.2.3.4
deploy: build
	@echo "Deploying to $(USER)@$(IP)..."
	scp $(BUILD_DIR)/$(BINARY_NAME) $(USER)@$(IP):/tmp/$(BINARY_NAME)
	@echo "Upload finished. Binary is at /tmp/$(BINARY_NAME)"

## help: Shows this help menu
help:
	@echo "Usage: make [target]"
	@echo ""
	@grep -E '^##' $(MAKEFILE_LIST) | sed -e 's/## //g' | column -t -s ':'


# All action-based targets
.PHONY: all build clean deploy help