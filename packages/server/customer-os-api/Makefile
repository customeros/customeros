run: build
	go run customer_os_api_main.go

generate:
	go run github.com/99designs/gqlgen generate

build: generate
	go build -v

swag:
	swag init -g customer_os_api_main.go

.PHONY: test
test:
	go test -v ./... -cover

install:
	go install gotest.tools/gotestsum@latest

clean:
	rm -rf customer-os-api

all: install build
