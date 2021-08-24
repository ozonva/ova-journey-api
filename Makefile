BIN_DIR = bin

.PHONY: run, test, generate, build, clean, lint, all

run:
	go run cmd/ova-journey-api/main.go

test:
	go test -race ./...
	go test ./...

generate:
	go generate ./...

build:
	go build -o $(BIN_DIR)/main cmd/ova-journey-api/main.go

clean:
	rm -rf $(BIN_DIR)

lint:
	golangci-lint run

all: clean lint test build