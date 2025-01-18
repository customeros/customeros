run: build
	go run github.com/customeros/customeros/packages/server/customer-os-common-module

build:
	go build -v

.PHONY: test
test:
	go test -v ./... -cover

install:
	go install gotest.tools/gotestsum@latest

clean:
	rm -rf github.com/customeros/customeros/packages/server/customer-os-common-module

all: install build
