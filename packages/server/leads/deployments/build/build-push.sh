#!/bin/bash
set -eo pipefail

REGISTRY=$1
IMAGE_NAME=$2
SHA=$3
PLATFORM=$4
ARCH=$5

echo "Building and pushing ${PLATFORM} image..."

# Use Docker Buildx to build and push
docker buildx build \
  --platform ${PLATFORM} \
  --push \
  --tag ${REGISTRY}/${IMAGE_NAME}:${SHA}-${ARCH} \
  --provenance=false \
  --file ./packages/server/leads/deployments/build/Containerfile \
  .

# Return success regardless of digest parsing
echo "Successfully built and pushed ${REGISTRY}/${IMAGE_NAME}:${VERSION}-${ARCH}"
exit 0
