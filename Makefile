BIN_DIR = bin

.PHONY: run, test, build, clean, lint, all

run:
	go run cmd/ova-journey-api/main.go

test:
	go test ./...

build:
	go build -o $(BIN_DIR)/main cmd/ova-journey-api/main.go

clean:
	rm -rf $(BIN_DIR)

lint:
	golangci-lint run

all: test build