
generate:
	export PATH="${PATH}:$(shell go env GOPATH)/bin";go generate ./...

build: generate
	go build -v -o bin/sync-tracked-data

test:
	go test -v ./...

install:

clean:
	find gen -mindepth 1 -maxdepth 1 -exec rm -r {} \;
	rm -rf bin

all: install build
