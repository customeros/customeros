#!/bin/bash
set -eo pipefail

# Arguments
CLOUD_REPO_URL=$1   
KUSTOMIZATION_PATH=$2 
NEW_TAG=$3           
GITHUB_TOKEN=$4      
COMMIT_MESSAGE="Update leads image to ${NEW_TAG}"

echo "Updating kustomization file for tag: ${NEW_TAG}"

# Set up Git user
git config --global user.name "GitHub Actions"
git config --global user.email "actions@github.com"

# Clone the cloud repository
TEMP_DIR=$(mktemp -d)
echo "Cloning repository to ${TEMP_DIR}..."
git clone https://${GITHUB_TOKEN}@${CLOUD_REPO_URL#https://} ${TEMP_DIR}

# Navigate to the repository
cd ${TEMP_DIR}

# Update the kustomization.yaml file
echo "Updating ${KUSTOMIZATION_PATH}..."
# Use yq if available, otherwise use sed
if command -v yq &> /dev/null; then
  yq -i '.images[0].newTag = "'${NEW_TAG}'"' ${KUSTOMIZATION_PATH}
else
  # Backup approach using sed 
  sed -i 's/newTag: .*/newTag: '${NEW_TAG}'/' ${KUSTOMIZATION_PATH}
fi

# Check if there are changes
if git diff --quiet; then
  echo "No changes to commit."
  exit 0
fi

# Commit and push changes
git add ${KUSTOMIZATION_PATH}
git commit -m "${COMMIT_MESSAGE}"
git push

echo "Successfully updated kustomization file and pushed changes."

# Clean up
rm -rf ${TEMP_DIR}
