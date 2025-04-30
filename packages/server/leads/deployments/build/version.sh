#!/bin/bash
# Exit immediately if a command exits with a non-zero status.
set -e

# Check if required arguments are provided
if [ -z "$1" ]; then
  echo "Error: No version bump type provided."
  exit 1
fi

if [ -z "$2" ]; then
  echo "Error: No app name provided."
  exit 1
fi

version_bump="$1"
app_name="$2"

echo "--- Running version.sh ---"
echo "Received version bump type: $version_bump"
echo "App name: $app_name"

# Find the latest SemVer tag for this specific app (app-name-vX.Y.Z)
echo "Finding latest SemVer tag for $app_name..."
tag_prefix="${app_name}-v"
latest_tag=$(git tag -l "${tag_prefix}[0-9]*" | grep -E "^${tag_prefix}[0-9]+\.[0-9]+\.[0-9]+$" | sort -V | tail -n 1 || echo "${tag_prefix}0.0.0") # sort -V for version sorting
echo "Latest tag found: $latest_tag"

# Parse the latest tag into major, minor, patch numbers
# Handle the v0.0.0 case explicitly during parsing
if [ "$latest_tag" == "${tag_prefix}0.0.0" ]; then
  major=0
  minor=0
  patch=0
else
  # Use grep with -oP for more robust parsing (requires GNU grep or compatible)
  # Fallback to sed if grep -P is not available (e.g., on some older systems)
  if grep -Pq '' /dev/null 2>/dev/null; then # Check if -P is supported
    # Extract version without the app-name prefix
    version_part=$(echo "$latest_tag" | grep -oP "${tag_prefix}\K.*")
    major=$(echo "$version_part" | grep -oP '^\K[0-9]+')
    minor=$(echo "$version_part" | grep -oP '^[0-9]+\.\K[0-9]+')
    patch=$(echo "$version_part" | grep -oP '^[0-9]+\.[0-9]+\.\K[0-9]+')
  else
    # Fallback sed parsing - less strict about format but works on more systems
    version_part=$(echo "$latest_tag" | sed "s/${tag_prefix}//")
    major=$(echo "$version_part" | sed 's/\([0-9]*\).\([0-9]*\).\([0-9]*\)/\1/')
    minor=$(echo "$version_part" | sed 's/\([0-9]*\).\([0-9]*\).\([0-9]*\)/\2/')
    patch=$(echo "$version_part" | sed 's/\([0-9]*\).\([0-9]*\).\([0-9]*\)/\3/')
  fi
fi

echo "Parsed current version: major=$major, minor=$minor, patch=$patch"

next_major="$major"
next_minor="$minor"
next_patch="$patch"
next_version=""

# Increment version based on bump type
case "$version_bump" in
  "major")
    echo "Performing major version bump..."
    next_major=$((major + 1))
    next_minor=0
    next_patch=0
    ;;
  "minor")
    echo "Performing minor version bump..."
    next_minor=$((minor + 1))
    next_patch=0
    ;;
  "patch")
    echo "Performing patch version bump..."
    next_patch=$((patch + 1))
    ;;
  *)
    echo "Error: Invalid version bump type '$version_bump'. Cannot generate next version."
    exit 1
    ;;
esac

# Create app-specific version tag
next_version="${app_name}-v${next_major}.${next_minor}.${next_patch}"
echo "Calculated next version: $next_version"

# Output the values for the GitHub Actions workflow
echo "latest_tag=$latest_tag" >> "$GITHUB_OUTPUT"
echo "next_version=$next_version" >> "$GITHUB_OUTPUT"

echo "--- version.sh finished ---"
exit 0
