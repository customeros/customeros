
generate:
	export PATH="${PATH}:$(shell go env GOPATH)/bin";go generate ./...
	$(shell cd proto;./generate.sh)

build: generate
	go build -v -o bin/message-store-api

test:
	go test -v ./...

install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install entgo.io/contrib/entproto/cmd/protoc-gen-entgrpc@latest
	go install gotest.tools/gotestsum@latest

clean:
	rm -rf bin
	rm -rf test/graph/generated
	rm -rf proto/generated
	rm test/graph/model/models_gen.go

all: install build
