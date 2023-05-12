clean:
	rm -rf bin

test:
	go test -v ./... -cover

build: bin/user-admin-api

bin/user-admin-api:
	go build -v -o bin/user-admin-api

all: clean build
