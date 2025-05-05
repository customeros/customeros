#!/bin/bash
set -eo pipefail

# Determine the version bump type from commit message
COMMIT_MSG=$(git log -1 --pretty=format:"%s")
VERSION_BUMP="patch"  # Default to patch

if [[ "$COMMIT_MSG" =~ ^major: ]]; then
  VERSION_BUMP="major"
elif [[ "$COMMIT_MSG" =~ ^minor: ]]; then
  VERSION_BUMP="minor"
elif [[ "$COMMIT_MSG" =~ ^patch: ]]; then
  VERSION_BUMP="patch"
elif [[ "$COMMIT_MSG" == "Merge pull request"* ]] && git rev-parse HEAD^2 &>/dev/null; then
  # Check PR commits for version prefix
  PR_COMMITS=$(git log HEAD^2 --not HEAD^ --pretty=format:"%s")
  if echo "$PR_COMMITS" | grep -q "^major:"; then
    VERSION_BUMP="major"
  elif echo "$PR_COMMITS" | grep -q "^minor:"; then
    VERSION_BUMP="minor"
  fi
fi

echo "Version bump type: $VERSION_BUMP"

# Find the latest SemVer tag
tag_prefix="v"  # Simplified to just use v prefix for all tags
latest_tag=$(git tag -l "${tag_prefix}[0-9]*" | grep -E "^${tag_prefix}[0-9]+\.[0-9]+\.[0-9]+$" | sort -V | tail -n 1 || echo "${tag_prefix}0.0.0")
echo "Latest tag found: $latest_tag"

# Parse the latest tag into major, minor, patch numbers
if [ "$latest_tag" == "${tag_prefix}0.0.0" ]; then
  major=0
  minor=0
  patch=0
else
  version_part=${latest_tag#$tag_prefix}
  major=$(echo "$version_part" | cut -d. -f1)
  minor=$(echo "$version_part" | cut -d. -f2)
  patch=$(echo "$version_part" | cut -d. -f3)
fi

echo "Current version: v$major.$minor.$patch"

# Increment version based on bump type
case "$VERSION_BUMP" in
  "major")
    major=$((major + 1))
    minor=0
    patch=0
    ;;
  "minor")
    minor=$((minor + 1))
    patch=0
    ;;
  "patch")
    patch=$((patch + 1))
    ;;
esac

# Create new version tag
VERSION="v${major}.${minor}.${patch}"
echo "New version: $VERSION"

# Output the values for the GitHub Actions workflow
echo "VERSION=$VERSION" >> $GITHUB_ENV
