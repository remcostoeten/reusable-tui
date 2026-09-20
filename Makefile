.PHONY: build run test golden lint vet tidy

build:
	go build -trimpath -o bin/tui ./cmd/tui

run:
	go run ./cmd/tui

test:
	go test ./...

golden:
	go test ./internal/app/ -run TestFrameGolden -update

lint:
	golangci-lint run ./...

vet:
	go vet ./...

tidy:
	go mod tidy
