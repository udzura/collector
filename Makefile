VERSION := $(shell go run main.go version | sed 's/version: //')

collector: test
	go build ./.

test: setup
	go test ./...

setup:
	go get ./...
