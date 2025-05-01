#!/usr/bin/env elvish

use str
use path
use os

# Configuration (would be passed as environment variables from GitHub Actions)
var registry = $E:REGISTRY  # ghcr.io
var image_name = $E:IMAGE_NAME  # github.repository_owner/customer-os-leads
var github_ref = $E:GITHUB_REF  # refs/heads/otter or PR ref
var github_repository = $E:GITHUB_REPOSITORY
var github_token = $E:GITHUB_TOKEN  # Injected secret
var github_workspace = $E:GITHUB_WORKSPACE

# Path configuration
var leads_dir = (path:join $github_workspace "packages" "server" "leads")
var build_scripts_dir = (path:join $leads_dir "deployments" "build")

# Utility functions
fn log [message]{
  echo "::group::$message"
}

fn log-end {
  echo "::endgroup::"
}

fn set-output [name value]{
  echo "$name=$value" >> $E:GITHUB_OUTPUT
}

fn is-otter-branch {
  eq $github_ref "refs/heads/otter"
}

# Main workflow logic
fn main {
  # Step 1: Build (runs on both PRs and pushes to otter)
  log "Building Go Binary"
  
  # Setup Go (simplified, assuming Go is installed)
  set E:CGO_ENABLED = 0
  set E:LANG = "C.UTF-8"
  set E:LC_ALL = "C.UTF-8"
  
  # Install dependencies
  cd $leads_dir
  external $build_scripts_dir/install-dependencies.sh
  
  # Build binary
  external $build_scripts_dir/build.sh
  
  log-end
  
  # Step 2: Check if release is needed (only on otter branch)
  if (not (is-otter-branch)) {
    echo "Not on otter branch, skipping release steps"
    return
  }
  
  log "Checking if release is needed"
  
  var check-result = (external $build_scripts_dir/check-release-needed.sh "leads")
  var version-bump = (echo $check-result | grep -o 'version_bump=.*' | str:trim | str:split = | last)
  var app-name = (echo $check-result | grep -o 'app_name=.*' | str:trim | str:split = | last)
  
  echo "Version bump: $version-bump"
  echo "App name: $app-name"
  
  if (eq $version-bump "skip") {
    echo "No release needed, skipping"
    return
  }
  
  log-end
  
  # Step 3: Build and release
  log "Building and releasing"
  
  # Determine version
  var version-result = (external $build_scripts_dir/version.sh $version-bump $app-name)
  var next-version = (echo $version-result | grep -o 'next_version=.*' | str:trim | str:split = | last)
  
  echo "Next version: $next-version"
  
  # Build and push Docker image
  var digest = (external $build_scripts_dir/build-push.sh $registry $image_name $next-version "linux/arm64" "arm64")
  
  # Create latest tag
  external docker buildx imagetools create \
    --tag $registry/$image_name:latest \
    $registry/$image_name:$next-version-arm64
  
  # Create Git tag
  external git config --local user.email "action@github.com"
  external git config --local user.name "GitHub Action"
  external git tag -a $next-version -m "Release $next-version"
  external git push origin $next-version
  
  # Create GitHub release (using GitHub CLI for simplicity)
  external gh release create $next-version \
    --title "Leads Release $next-version" \
    --generate-notes \
    --repo $github_repository
  
  log-end
}

# Execute main function
main
