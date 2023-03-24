generate:
	export PATH="${PATH}:$(shell go env GOPATH)/bin";go generate ./...

clean:
	rm -rf bin
	rm -rf test/graph/generated
	rm -rf test/graph/model/models_gen.go

test:
	go test -v ./... -cover

test/graph/generated/generated.go: generate


build: bin/comms-api test/graph/generated/generated.go

bin/comms-api:
	go build -v -o bin/comms-api

all: clean build
