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

# Get and output the image digest
digest=$(docker buildx imagetools inspect ${REGISTRY}/${IMAGE_NAME}:${SHA}-${ARCH} --format "{{json .Manifest}}" | jq -r '.digest')
echo "${digest}"
