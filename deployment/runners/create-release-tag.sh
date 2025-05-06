#!/bin/bash
set -e

echo "Creating and pushing tag $VERSION..."
git config --local user.email "action@github.com"
git config --local user.name "GitHub Action"
git tag -a $VERSION -m "Release $VERSION"
git push origin $VERSION --force
