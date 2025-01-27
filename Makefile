APP_NAME = blocktime-node
SRC = cmd/server/main.go
BUILD_DIR = bin
GO = go
GOFMT = gofmt

all: build

build: $(SRC)
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(SRC)

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

fmt:
	$(GOFMT) -s -w .

clean:
	rm $(BUILD_DIR)/$(APP_NAME)

install:
	cp bin/blocktime-node /usr/bin

help:
	@echo "Makefile for blocktime-node"
	@echo "Usage:"
	@echo "  make            Build the application"
	@echo "  make run        Build and run the application"
	@echo "  make test       Run tests"
	@echo "  make vet        Vet code"
	@echo "  make fmt        Format code"
	@echo "  make clean      Clean up build artifacts"
	@echo "  make help       Show this help message"

.PHONY: all build run test clean help
