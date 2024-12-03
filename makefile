# Build variables
BINARY_NAME=ehost
BUILD_DIR=build

# Git version information
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "-dirty" || echo "")
VERSION=$(GIT_COMMIT)$(GIT_DIRTY)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

# Supported GOOS and GOARCH
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: all clean build build-all

all: clean build

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

build:
	$(GOBUILD) -o $(BINARY_NAME) -ldflags "-X main.Version=$(VERSION)" .

build-all:
	mkdir -p $(BUILD_DIR)
	$(foreach platform,$(PLATFORMS),\
		$(eval OS := $(word 1,$(subst /, ,$(platform))))\
		$(eval ARCH := $(word 2,$(subst /, ,$(platform))))\
		GOOS=$(OS) GOARCH=$(ARCH) $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)_$(OS)_$(ARCH)$(if $(findstring windows,$(OS)),.exe,) -ldflags "-X main.Version=$(VERSION)" . ;)

# Individual platform builds
linux-amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)_linux_amd64 -ldflags "-X main.Version=$(VERSION)" .

darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)_darwin_amd64 -ldflags "-X main.Version=$(VERSION)" .

windows-amd64:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)_windows_amd64.exe -ldflags "-X main.Version=$(VERSION)" .
