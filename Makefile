# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary name
BINARY_NAME=re

# Build target
all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)_x86
	rm -f $(BINARY_NAME)_win_x86.exe
	rm -f $(BINARY_NAME)_darwin_x86
	rm -f $(BINARY_NAME)_arm
	rm -f $(BINARY_NAME)_darwin_arm

run: build
	./$(BINARY_NAME)

deps:
	$(GOMOD) download

# Cross compilation
build-x86:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_x86 -v
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_win_x86.exe -v
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_darwin_x86 -v

build-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)_arm -v
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)_darwin_arm -v

build-all:
	make build-x86
	make build-arm

.PHONY: all build test clean run deps build-x86 build-arm build-all
