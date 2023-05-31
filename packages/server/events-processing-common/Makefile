all: clean build

generate: $(wildcard grpc/v1/*.proto)
	buf generate

build: generate
	#go build -v

test:
	#go test -v ./... -cover

install:
	go install gotest.tools/gotestsum@latest

clean:
	rm -rf gen

