# Detect the operating system
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

# Set installation directories based on the operating system
ifeq ($(UNAME_S), Linux)
	INSTALL_DIR := /usr/local/bin
else ifeq ($(UNAME_S), Darwin)
	ifeq ($(UNAME_M), arm64)
		INSTALL_DIR := /opt/homebrew/bin
	else
		INSTALL_DIR := /usr/local/bin
	endif
endif

# Project name
PROJECT_NAME := ehost

# Go compiler and flags
GO := go
GOFLAGS := -ldflags "-s -w"

# Versioning
VERSION := $(shell git describe --tags --always --dirty)

# Default target
all: build

# Build target
build:
	$(GO) build $(GOFLAGS) -o $(PROJECT_NAME)-$(VERSION) main.go

# Install target
install: build
	install -m 0755 $(PROJECT_NAME)-$(VERSION) $(INSTALL_DIR)/$(PROJECT_NAME)

# Clean target
clean:
	rm -f $(PROJECT_NAME)

# Uninstall target
uninstall:
	rm -f $(INSTALL_DIR)/$(PROJECT_NAME)

# Phony targets
.PHONY: all build install clean uninstall
