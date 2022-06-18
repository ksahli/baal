clean:
	go clean -testcache -cache

format:
	go fmt ./...

verify:
	go fmt ./...
	go vet ./...

test:
	go fmt ./...
	go vet ./...
	go test ./... -cover -covermode=atomic

install:
	go install

import:
	goimports -w .

build:
	go build -o baal
