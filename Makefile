collector: test
	go build ./.

test: setup
	go test ./...

setup:
	go get ./...
