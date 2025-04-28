#!/bin/bash
set -eo pipefail

REGISTRY=$1
IMAGE_NAME=$2
SHA=$3
SAFE_TAG=$4
REF_NAME=$5
EVENT_NAME=$6

# Create and push the SHA manifest
echo "Creating SHA manifest..."
docker buildx imagetools create \
  --tag ${REGISTRY}/${IMAGE_NAME}:${SHA} \
  ${REGISTRY}/${IMAGE_NAME}:${SHA}-amd64 \
  ${REGISTRY}/${IMAGE_NAME}:${SHA}-arm64

# Create and push the branch/PR manifest
echo "Creating branch/PR manifest..."
docker buildx imagetools create \
  --tag ${REGISTRY}/${IMAGE_NAME}:${SAFE_TAG} \
  ${REGISTRY}/${IMAGE_NAME}:${SHA}-amd64 \
  ${REGISTRY}/${IMAGE_NAME}:${SHA}-arm64

# If this is the otter branch, also tag as latest
if [[ "${REF_NAME}" == "otter" ]]; then
  echo "Creating latest manifest..."
  docker buildx imagetools create \
    --tag ${REGISTRY}/${IMAGE_NAME}:latest \
    ${REGISTRY}/${IMAGE_NAME}:${SHA}-amd64 \
    ${REGISTRY}/${IMAGE_NAME}:${SHA}-arm64
fi

# For releases, also tag with the release tag
if [[ "${EVENT_NAME}" == "release" ]]; then
  # For releases, the tag should be the release tag
  RELEASE_TAG=$(echo "${REF_NAME}" | sed 's/\//-/g')
  echo "Creating release manifest for ${RELEASE_TAG}..."
  docker buildx imagetools create \
    --tag ${REGISTRY}/${IMAGE_NAME}:${RELEASE_TAG} \
    ${REGISTRY}/${IMAGE_NAME}:${SHA}-amd64 \
    ${REGISTRY}/${IMAGE_NAME}:${SHA}-arm64
fi

# Inspect the created manifest
echo "Inspecting manifest..."
docker buildx imagetools inspect ${REGISTRY}/${IMAGE_NAME}:${SAFE_TAG}
