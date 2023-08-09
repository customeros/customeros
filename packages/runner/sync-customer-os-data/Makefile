run: build
	go run sync_customer_os_data_main.go

generate:
	export PATH="${PATH}:$(shell go env GOPATH)/bin";

build: generate
	go build -v -o bin/sync-customer-os-data

.PHONY: test
test:
	go test -v ./...

install:
	go install gotest.tools/gotestsum@latest

clean:
	find gen -mindepth 1 -maxdepth 1 -exec rm -r {} \;
	rm -rf bin

all: install build
