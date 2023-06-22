run: build
	go run customer-os-api.go

generate:
	go run github.com/99designs/gqlgen generate

build: generate
	go build -v

test:
	go test -v ./... -cover

install:
	go install gotest.tools/gotestsum@latest

clean:
	rm -rf customer-os-api

all: install build test
