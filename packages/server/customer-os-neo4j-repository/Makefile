build:
	go build -v ./...

.PHONY: test
test:
	go test -v ./... -cover

install:
	go install gotest.tools/gotestsum@latest

clean:
	rm -rf ./bin/*

all: install build
