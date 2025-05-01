#!/usr/bin/env elvish

# Strict error handling (similar to bash's set -e)
fn check-err [err]{
 if (not-eq $err $nil) {
   fail $err
 }
}

# Step into leads directory
cd ./packages/server/leads/

echo "Cleaning and updating dependencies..."
try {
 go mod tidy
} catch e {
 fail "Failed to run go mod tidy: "$e
}

echo "Generating GraphQL code..."
try {
 go run github.com/99designs/gqlgen generate --config ./api/graphql/gqlgen.yml
} catch e {
 fail "Failed to generate GraphQL code: "$e
}

echo "Generating protobuf code..."
try {
 var proto-files = (find ./proto -name "*.proto" -type f)
 for proto-file $proto-files {
   protoc \
     --proto_path=./proto \
     --go_out=./proto/pb \
     --go_opt=module=github.com/customeros/customeros/packages/server/leads/proto/pb \
     --go-grpc_out=./proto/pb \
     --go-grpc_opt=module=github.com/customeros/customeros/packages/server/leads/proto/pb \
     $proto-file
 }
} catch e {
 fail "Failed to generate protobuf code: "$e
}

echo "Building application..."
try {
 go build -v .
} catch e {
 fail "Failed to build application: "$e
}

echo "Build completed successfully."
