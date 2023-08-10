run: build
	go run events-processing-platform.go

build:
	go build -v

.PHONY: test
test:
	go test -v ./... -cover

install:
	go install gotest.tools/gotestsum@latest

clean:
	rm -rf events-processing-platform

all: install build
