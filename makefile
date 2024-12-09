# Variables
BIN_NAME ?= some
TARGET_DIR ?= /usr/local/bin

# Default target
all: build

# Build the Go program
build:
	go build -o $(BIN_NAME)

# Copy the binary to the target directory
install: build
	cp ./$(BIN_NAME) $(TARGET_DIR)

# Clean up build artifacts
clean:
	rm -f $(BIN_NAME)

config: 
	@echo "template_root=$(shell pwd)" > $(HOME)/.$(BIN_NAME)rc
