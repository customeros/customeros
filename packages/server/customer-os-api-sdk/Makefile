run: build
	go run main.go

generate:
	go get github.com/99designs/gqlgen@v0.17.47
	go run github.com/99designs/gqlgen generate

clean:
	rm -rf bin
	rm -rf graph/generated
	rm -rf graph/model/models_gen.go

test:
	go test -v ./... -cover

graph/generated/generated.go: generate

build: bin/file-store-api graph/generated/generated.go

bin/file-store-api:
	go build -v -o bin/file-store-api


all: clean  build
