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

# Inspect the image but handle potential errors with the digest format
digest_output=$(docker buildx imagetools inspect ${REGISTRY}/${IMAGE_NAME}:${VERSION}-${ARCH} --format "{{json .Manifest}}" 2>/dev/null || echo '{"digest":"unknown"}')
digest=$(echo $digest_output | jq -r '.digest' 2>/dev/null || echo "unknown")

# Check if we got a valid digest
if [[ "$digest" == "unknown" || "$digest" == "null" ]]; then
  echo "Could not get digest, but image was pushed successfully"
else
  echo "Image digest: $digest"
fi

# Return success regardless of digest parsing
echo "Successfully built and pushed ${REGISTRY}/${IMAGE_NAME}:${VERSION}-${ARCH}"
exit 0
