#!/bin/bash
set -eo pipefail

APP=$1
VERSION=$2

# Set Git identity for the commit
cd cloud-repo
git config user.name "GitHub Actions Bot"
git config user.email "actions@github.com"

# Path to the file that needs to be updated
FILE_PATH="./deployments/${APP}/kustomization.yaml"

# Update the version in the specified file
sed -i "s|newTag: .*|newTag: ${VERSION}|" "${FILE_PATH}"

# Commit and push changes
git add $FILE_PATH
git commit -m "Update leads version to ${VERSION}"
git push
