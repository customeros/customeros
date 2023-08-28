run: build
	go run github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module

build:
	go build -v

.PHONY: test
test:
	go test -v ./... -cover

install:
	go install gotest.tools/gotestsum@latest

clean:
	rm -rf github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module

all: install build
