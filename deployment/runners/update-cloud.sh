#!/bin/bash
set -eo pipefail

APP=$1
VERSION=$2

# Navigate to the cloud repo directory
cd cloud-repo

# Set Git identity for the commit
git config user.name "GitHub Actions Bot"
git config user.email "actions@github.com"

# Get latest changes
git pull origin main  

# Path to the file that needs to be updated
FILE_PATH="./deployments/${APP}/kustomization.yaml"

# Check if file exists
if [ ! -f "${FILE_PATH}" ]; then
  echo "Error: File ${FILE_PATH} not found!"
  exit 1
fi

# Update the version in the specified file
sed -i "s|newTag: .*|newTag: ${VERSION}|" "${FILE_PATH}"

# Check if there are changes to commit
if git diff --quiet "${FILE_PATH}"; then
  echo "No changes detected in ${FILE_PATH}. Version might already be ${VERSION}."
  exit 0
fi

# Commit and push changes
git add "${FILE_PATH}"
git commit -m "Update ${APP} version to ${VERSION}"
git push origin HEAD

echo "Successfully updated ${APP} to version ${VERSION}"
