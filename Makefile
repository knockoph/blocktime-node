APP_NAME_SERVER = blocktime-node
APP_NAME_NOTIFY = blocktime-node-notify
SRC_SERVER = cmd/server/main.go
SRC_NOTIFY = cmd/notify/main.go
BUILD_DIR = bin
GO = go
GOFMT = gofmt

all: build

build: $(SRC)
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME_SERVER) $(SRC_SERVER)
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME_NOTIFY) $(SRC_NOTIFY)

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

fmt:
	$(GOFMT) -s -w .

clean:
	rm $(BUILD_DIR)/$(APP_NAME_SERVER)

install:
	cp bin/blocktime-node /usr/bin
	cp bin/blocktime-node-notify /usr/bin

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
