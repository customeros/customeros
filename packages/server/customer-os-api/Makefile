
generate:
    go run github.com/99designs/gqlgen generate

build: generate
	go build -v

test:
	go test -v ./... -cover

install:
	go install gotest.tools/gotestsum@latest

clean:

all: install build
