#!/bin/bash
set -eo pipefail

APP_PATH=$1

cd ./${APP_PATH}

echo "Updating dependencies..."
go mod tidy

echo "Generating GraphQL code..."
go run github.com/99designs/gqlgen generate --config ./api/graphql/gqlgen.yml

echo "Generating protobuf code..."
find ./proto -name "*.proto" -type f -exec \
  protoc \
  --proto_path=./proto \
  --go_out=./proto/pb \
  --go_opt=module=github.com/customeros/customeros/${APP_PATH}/proto/pb \
  --go-grpc_out=./proto/pb \
  --go-grpc_opt=module=github.com/customeros/${APP_PATH}/proto/pb \
  {} \;

echo "Building application..."
go build -v .

echo "Build completed successfully."

