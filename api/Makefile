.PHONY: fmt lint test build run

fmt:
	gofmt -s -w .

lint:
	golangci-lint run

test:
	go test ./...

build:
	go build

run:
	go run main.go
