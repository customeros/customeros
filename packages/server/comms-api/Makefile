clean:
	rm -rf bin

test:
	go test -v ./... -cover

build: bin/comms-api

bin/comms-api:
	go build -v -o bin/comms-api

all: clean build
