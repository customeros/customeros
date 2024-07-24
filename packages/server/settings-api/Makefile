
build:
	go build -v

test:
	go test -v ./... -cover

install:
	go install gotest.tools/gotestsum@latest

swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --parseDependency --parseInternal -g settings_api_main.go

clean:

all: install swagger build
