#!/bin/bash
set -eo pipefail

EVENT_NAME=$1
REF_NAME=$2
PR_NUMBER=$3

# For PR events, use pr-NUMBER format instead of the branch name
if [[ "${EVENT_NAME}" == "pull_request" ]]; then
  echo "pr-${PR_NUMBER}"
else
  # For normal branches, replace / with - to make it safe for tags
  echo "${REF_NAME}" | sed 's/\//-/g'
fi
