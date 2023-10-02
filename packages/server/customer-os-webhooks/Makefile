run: build
	go run customer_os_webhooks_main.go

generate:

build: generate
	go build -v

.PHONY: test
test:
	go test -v ./... -cover

install:
	go install gotest.tools/gotestsum@latest

clean:
	rm -rf customer-os-webhooks

all: install build
