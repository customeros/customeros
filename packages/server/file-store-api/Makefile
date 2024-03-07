run: build
	go run main.go

test:
	go test -v ./... -cover

build: bin/file-store-api

bin/file-store-api:
	go build -v -o bin/file-store-api

install:
	go install gotest.tools/gotestsum@latest

all: install build
