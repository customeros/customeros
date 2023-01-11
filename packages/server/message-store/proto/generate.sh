#!/bin/bash

mkdir -p ../gen/proto
protoc --go_out=../gen/proto --go_opt=paths=source_relative --go-grpc_out=../gen/proto --go-grpc_opt=paths=source_relative ./*.proto
sed -i "" -e "s/,omitempty//g" ../gen/proto/*.go