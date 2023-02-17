
run:

generate:
	go run github.com/99designs/gqlgen generate

build: generate
	go build -v

test:
	go test -v ./... -cover

install:
	go install gotest.tools/gotestsum@v1.9.0

clean:

all: install build
