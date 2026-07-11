.PHONY: fmt lint test build run

fmt:
	cd api && gofmt -s -w .

lint:
	cd api && golangci-lint run

test:
	cd api && go test ./...

build:
	cd api && go build

run:
	cd api && go run main.go
