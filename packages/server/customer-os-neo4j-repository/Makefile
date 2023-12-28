run: build
	go run github.com/openline-ai/customer-os-neo4j-repository

build:
	go build -v

.PHONY: test
test:
	go test -v ./... -cover

install:
	go install gotest.tools/gotestsum@latest

clean:
	rm -rf github.com/openline-ai/customer-os-neo4j-repository

all: install build
