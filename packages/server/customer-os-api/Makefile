PYTHON ?= python3
OPENAPI_INPUT = docs/openapi.yaml
OPENAPI_OUTPUT_DIR = docs/

run: build
	go run customer_os_api_main.go

generate:
	go run github.com/99designs/gqlgen generate

build: generate
	go build -v

.PHONY: test doc-deps doc-clean docs
test:
	go test -v ./... -cover

install:
	go install gotest.tools/gotestsum@latest

clean:
	rm -rf customer-os-api

all: install build

doc-deps:
	go install github.com/swaggo/swag/cmd/swag@latest
	$(PYTHON) -m pip install pyyaml

doc-clean:
	rm -rf $(OPENAPI_OUTPUT_DIR)

docs: doc-clean
	swag init -g customer_os_api_main.go
	$(PYTHON) swag_convert.py docs/swagger.yaml docs/openapi.yaml 
	$(PYTHON) swag_split.py docs/openapi.yaml docs/
	rm docs/swagger.yaml
	rm docs/openapi.yaml


