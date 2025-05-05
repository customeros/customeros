#!/bin/bash
set -eo pipefail

# Use a known valid protoc version (e.g., 27.0 or find the latest valid one)
PROTOC_VERSION=30.2
PROTOC_GEN_VERSION=v1.36.6
PROTOC_GRPC_VERSION=v1.5.1
GQLGEN_VERSION=0.17.72 

sudo apt install unzip

# Setup directories
TEMP_DIR=$(mktemp -d)
# Use trap for cleanup to ensure it runs even on errors before the end
trap 'cd /; rm -rf "$TEMP_DIR"' EXIT
cd "$TEMP_DIR"

echo "Installing protoc v${PROTOC_VERSION}..."
PROTOC_ZIP=protoc-${PROTOC_VERSION}-linux-x86_64.zip
curl -sSL -o "$PROTOC_ZIP" "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}"
# Check if download was successful
if [ ! -f "$PROTOC_ZIP" ]; then
  echo "Error: Failed to download protoc zip."
  exit 1
fi
sudo unzip -qo "$PROTOC_ZIP" -d /usr/local bin/protoc
sudo unzip -qo "$PROTOC_ZIP" -d /usr/local 'include/*' -x readme.txt # Be more specific, quiet (-q), overwrite (-o)
rm -f "$PROTOC_ZIP"
echo "protoc version: $(protoc --version)"

echo "Installing Go tools..."
# Ensure Go is set up (usually done by actions/setup-go before this script runs)
# Determine Go bin path reliably
GOBIN=$(go env GOBIN)
if [ -z "$GOBIN" ]; then
  GOBIN=$(go env GOPATH)/bin
fi
# Fallback if GOPATH is also not set (common case in simple setups)
if [ ! -d "$GOBIN" ]; then
    GOBIN="$HOME/go/bin"
fi

echo "Attempting to install Go tools to: $GOBIN"
go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_VERSION}
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GRPC_VERSION}
go install github.com/99designs/gqlgen@v${GQLGEN_VERSION}

# Add the Go bin directory to the PATH for this script's environment
echo "Adding $GOBIN to PATH"
export PATH="$PATH:$GOBIN"
echo "Current PATH: $PATH" # For debugging

# Verify installations
echo "Verifying installations..."
which protoc
which protoc-gen-go
which protoc-gen-go-grpc
which gqlgen

echo "All dependencies installed successfully."

