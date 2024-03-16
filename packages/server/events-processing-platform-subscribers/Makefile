run: build
	go run events_processing_platform_subscribers.go

generate:

build: generate
	go build -v

.PHONY: test
test:
	go test -v ./... -cover

install:
	go install gotest.tools/gotestsum@latest

clean:

all: install build
