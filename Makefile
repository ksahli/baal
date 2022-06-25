BINARY_NAME=baal

clean:
	go clean -testcache -cache
	rm ${BINARY_NAME}-freebsd
	rm ${BINARY_NAME}-openbsd
	rm ${BINARY_NAME}-linux

format:
	go fmt ./...

verify:
	go vet ./...

test:
	go test ./... -coverprofile coverage.out

install:
	go install

build:
	GOARCH=amd64 GOOS=freebsd go build -o ${BINARY_NAME}-freebsd main.go
	GOARCH=amd64 GOOS=openbsd go build -o ${BINARY_NAME}-openbsd main.go
	GOARCH=amd64 GOOS=linux   go build -o ${BINARY_NAME}-linux   main.go


