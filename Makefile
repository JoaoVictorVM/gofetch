BINARY := gofetch
ifeq ($(OS),Windows_NT)
BINARY := $(BINARY).exe
endif

.PHONY: fmt vet lint test build run

fmt:
	gofmt -s -w .

vet:
	go vet ./...

lint:
	golangci-lint run

test:
	go test ./...

build:
	go build -o bin/$(BINARY) ./cmd/gofetch

run:
	go run ./cmd/gofetch
