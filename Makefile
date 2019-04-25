test:
	go test ./...
build:
	go build -o archiver main.go

.PHONY: test build