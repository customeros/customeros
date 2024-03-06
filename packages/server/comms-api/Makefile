run: build
	go run main.go

test:
	go test -v ./... -cover

build: bin/comms-api

bin/comms-api:
	go build -v -o bin/comms-api

install:
	go install gotest.tools/gotestsum@latest

all: build
