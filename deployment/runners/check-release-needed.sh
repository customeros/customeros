#!/bin/bash
set -e

# Check if a tag already exists for the current commit
if git tag --points-at HEAD | grep -q "v[0-9]"; then
  echo "RELEASE_NEEDED=false" >> $GITHUB_ENV
  echo "No release needed (tag exists)"
  exit 0
fi

# Get the last commit message
commit_msg=$(git log -1 --pretty=format:"%s")

# Check for version prefixes in commit message or PR commits
if [[ "$commit_msg" =~ ^major: ]] || [[ "$commit_msg" =~ ^minor: ]] || [[ "$commit_msg" =~ ^patch: ]]; then
  echo "RELEASE_NEEDED=true" >> $GITHUB_ENV
  echo "Release needed (version prefix in commit)"
  exit 0
fi

# Check PR commits if it's a merge
if [[ "$commit_msg" == "Merge pull request"* ]] && git rev-parse HEAD^2 &>/dev/null; then
  pr_commits=$(git log HEAD^2 --not HEAD^ --pretty=format:"%s")
  if echo "$pr_commits" | grep -q "^major:\|^minor:\|^patch:"; then
    echo "RELEASE_NEEDED=true" >> $GITHUB_ENV
    echo "Release needed (version prefix in PR commits)"
    exit 0
  fi
fi

# Default: no release needed
echo "RELEASE_NEEDED=false" >> $GITHUB_ENV
echo "No release needed (no version prefix found)"
exit 0
