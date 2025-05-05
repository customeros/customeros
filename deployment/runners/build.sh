#!/bin/bash
set -eo pipefail

# Default values
GENERATE_GRAPHQL=true
GENERATE_PROTOBUF=true

# Display usage information
function show_usage {
  echo "Usage: $0 [options] APP_PATH"
  echo "Options:"
  echo "  -g  Skip GraphQL code generation"
  echo "  -p  Skip protobuf code generation"
  echo "  -h  Show this help message"
  echo "Arguments:"
  echo "  APP_PATH  Path to the application directory"
  exit 1
}

# Parse command-line options
while getopts ":gph" opt; do
  case ${opt} in
    g)
      GENERATE_GRAPHQL=false
      ;;
    p)
      GENERATE_PROTOBUF=false
      ;;
    h)
      show_usage
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      show_usage
      ;;
  esac
done

# Shift the options so $1 becomes the first non-option argument
shift $((OPTIND - 1))

# Check if APP_PATH is provided
if [ -z "$1" ]; then
  echo "Error: APP_PATH is required"
  show_usage
fi

APP_PATH=$1
GQLGEN_PATH=${2:-"./api/graphql/gqlgen.yml"}

cd ./${APP_PATH}

echo "Updating dependencies..."
go mod tidy

if [ "$GENERATE_GRAPHQL" = true ]; then
  echo "Generating GraphQL code..."
  go run github.com/99designs/gqlgen generate --config ${GQLGEN_PATH}
else
  echo "Skipping GraphQL code generation..."
fi

if [ "$GENERATE_PROTOBUF" = true ]; then
  echo "Generating protobuf code..."
  find ./proto -name "*.proto" -type f -exec \
    protoc \
    --proto_path=./proto \
    --go_out=./proto/pb \
    --go_opt=module=github.com/customeros/customeros/${APP_PATH}/proto/pb \
    --go-grpc_out=./proto/pb \
    --go-grpc_opt=module=github.com/customeros/${APP_PATH}/proto/pb \
    {} \;
else
  echo "Skipping protobuf code generation..."
fi

echo "Building application..."
go build -v .

echo "Build completed successfully."

