#!/bin/bash
# Exit immediately if a command exits with a non-zero status.
set -e
echo "--- Running check_release_needed.sh ---"

# App name is now the primary parameter
APP_NAME=${1:-$APP_NAME}

if [ -z "$APP_NAME" ]; then
  echo "Error: No app name specified. Please provide APP_NAME as first argument or environment variable."
  exit 1
fi

echo "Checking for releases for app: $APP_NAME"

# Map of app names to their directories in the monorepo
# This can be moved to a config file if you prefer
declare -A APP_PATHS
APP_PATHS["leads"]="packages/server/leads"
APP_PATHS["users"]="packages/server/users"
APP_PATHS["admin"]="packages/client/admin"
# Add other apps as needed

# Get the directory path for the app
APP_DIR=${APP_PATHS[$APP_NAME]}

if [ -z "$APP_DIR" ]; then
  echo "Error: No path configured for app '$APP_NAME'. Please add it to the APP_PATHS map."
  exit 1
fi

echo "App directory path: $APP_DIR"

# Check if a tag already exists for the current commit (HEAD)
echo "Checking for existing tags on commit $(git rev-parse HEAD)..."
existing_tags=$(git tag --points-at HEAD)
if [ -n "$existing_tags" ]; then
  # Check if there's an app-specific tag in the format app-name-vX.Y.Z
  app_tag_prefix="${APP_NAME}-v"
  
  if echo "$existing_tags" | grep -q "$app_tag_prefix"; then
    echo "App-specific tag already exists for this commit: $(echo "$existing_tags" | grep "$app_tag_prefix")"
    echo "version_bump=skip" >> "$GITHUB_OUTPUT"
    echo "app_name=$APP_NAME" >> "$GITHUB_OUTPUT"
    echo "--- check_release_needed.sh finished (skipped due to existing app-specific tag) ---"
    exit 0 # Exit successfully, signaling skip
  else
    echo "No app-specific tag found for $APP_NAME, continuing..."
  fi
else
  echo "No existing tag found on this commit."
fi

# Get the last commit message
echo "Getting last commit message..."
merge_msg=$(git log -1 --pretty=format:"%s")
echo "Last commit message: $merge_msg"

# Initialize variables
version_bump="skip" # Default to skip
echo "app_name=$APP_NAME" >> "$GITHUB_OUTPUT"

# Check if this commit modified files in the specific app directory
if git diff --name-only HEAD HEAD~1 | grep -q "^${APP_DIR}/"; then
  echo "Changes detected in $APP_DIR directory."
  
  # Check if it's a standard GitHub merge commit message
  if [[ "$merge_msg" == "Merge pull request"* ]]; then
    echo "Identified as a PR merge commit."
    # Get the commits from the PR (between merge base and PR head)
    # Assuming standard merge: HEAD^ is merge base, HEAD^2 is PR branch tip
    echo "Getting PR commits between $(git rev-parse HEAD^) and $(git rev-parse HEAD^2)..."
    
    # Check if HEAD^2 exists (it might not for squash merges or initial commits)
    if git rev-parse HEAD^2 &>/dev/null; then
      # Check if any of the changed files in the PR are in our app directory
      changed_files=$(git diff --name-only HEAD^ HEAD^2)
      if echo "$changed_files" | grep -q "^${APP_DIR}/"; then
        echo "PR contains changes to the $APP_DIR directory."
        
        pr_commits=$(git log HEAD^2 --not HEAD^ --pretty=format:"%s")
        echo "PR commits:"
        echo "$pr_commits"
        
        # Check for version prefixes in PR commits
        if echo "$pr_commits" | grep -q "^major:"; then
          echo "Found 'major:' prefix in PR commits."
          version_bump="major"
        elif echo "$pr_commits" | grep -q "^minor:"; then
          echo "Found 'minor:' prefix in PR commits."
          version_bump="minor"
        elif echo "$pr_commits" | grep -q "^patch:"; then
          echo "Found 'patch:' prefix in PR commits."
          version_bump="patch"
        fi
      else
        echo "PR does not contain changes to the $APP_DIR directory. Skipping release check."
      fi
    else
      echo "Could not find HEAD^2 (likely a squash merge or single-commit PR). Checking the merge commit message itself."
      # Fallback for squash merges or non-standard merges - check the merge message itself
      if [[ "$merge_msg" =~ ^major: ]]; then
        echo "Found 'major:' prefix in merge commit message."
        version_bump="major"
      elif [[ "$merge_msg" =~ ^minor: ]]; then
        echo "Found 'minor:' prefix in merge commit message."
        version_bump="minor"
      elif [[ "$merge_msg" =~ ^patch: ]]; then
        echo "Found 'patch:' prefix in merge commit message."
        version_bump="patch"
      fi
    fi
  else
    echo "Not a standard PR merge commit. Checking the commit message directly."
    # Direct commit to branch or squash merge with custom message
    if [[ "$merge_msg" =~ ^major: ]]; then
      echo "Found 'major:' prefix in commit message."
      version_bump="major"
    elif [[ "$merge_msg" =~ ^minor: ]]; then
      echo "Found 'minor:' prefix in commit message."
      version_bump="minor"
    elif [[ "$merge_msg" =~ ^patch: ]]; then
      echo "Found 'patch:' prefix in commit message."
      version_bump="patch"
    fi
  fi
else
  echo "No changes detected in $APP_DIR directory. Skipping release."
fi

# Output the determined version_bump type
echo "Determined version bump type: $version_bump"
echo "version_bump=$version_bump" >> "$GITHUB_OUTPUT"

if [ "$version_bump" == "skip" ]; then
  echo "No version prefix found or no relevant changes. Skipping release."
else
  echo "Release is needed for $APP_NAME with a '$version_bump' bump."
fi

echo "--- check_release_needed.sh finished ---"
# Exit with success status regardless of bump type, as 'skip' is a valid outcome.
exit 0
