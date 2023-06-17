all: clean install build

generate: $(wildcard grpc/v1/*.proto)
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	buf generate

build: generate

test:

install:

clean:
	rm -rf gen

