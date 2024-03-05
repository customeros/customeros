run: build
	go run main.go

generate:
	export PATH="${PATH}:$(shell go env GOPATH)/bin";go generate ./...

clean:
	rm -rf bin
	rm -rf test/graph/generated
	rm -rf test/graph/model/models_gen.go

test:
	go test -v ./... -cover

test/graph/generated/generated.go: generate

build: bin/file-store-api test/graph/generated/generated.go

bin/file-store-api:
	go build -v -o bin/file-store-api

install:
	go install gotest.tools/gotestsum@latest

all: clean install build
