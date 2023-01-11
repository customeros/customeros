#!/bin/bash

mkdir -p ../proto/generated
protoc --go_out=../proto/generated --go_opt=paths=source_relative --go-grpc_out=../proto/generated --go-grpc_opt=paths=source_relative ./*.proto
sed -i "" -e "s/,omitempty//g" ../proto/generated/*.go