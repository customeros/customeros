#!/bin/bash
set -eo pipefail

cd ./packages/server/leads/

echo "Cleaning and updating dependencies..."
go mod tidy

echo "Generating GraphQL code..."
go run github.com/99designs/gqlgen generate --config ./api/graphql/gqlgen.yml

echo "Generating protobuf code..."
find ./proto -name "*.proto" -type f -exec \
  protoc \
  --proto_path=./proto \
  --go_out=./proto/pb \
  --go_opt=module=github.com/customeros/customeros/packages/server/leads/proto/pb \
  --go-grpc_out=./proto/pb \
  --go-grpc_opt=module=github.com/customeros/customeros/packages/server/leads/proto/pb \
  {} \;

echo "Building application..."
go build -v .

echo "Build completed successfully."
