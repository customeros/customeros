#!/bin/sh
# Exit immediately if a command exits with a non-zero status
set -e

echo "Starting Core application..."

# Check if we need to run migrations
if [ "${RUN_MIGRATIONS:-true}" = "true" ]; then
  echo "Running migrations..."
  bin/core eval "Core.Release.migrate()" || {
    echo "Migration failed!"
    exit 1
  }
  echo "Migrations completed successfully!"
fi

# Parse command arguments
case "$1" in
  start)
    echo "Starting application..."
    exec bin/core start
    ;;
  start_iex)
    echo "Starting application with IEx attached..."
    exec bin/core start_iex
    ;;
  remote)
    echo "Starting remote console..."
    exec bin/core remote
    ;;
  eval)
    shift
    echo "Evaluating command..."
    exec bin/core eval "$@"
    ;;
  *)
    # Default behavior - start the application
    if [ -z "$1" ]; then
      echo "Starting application (default)..."
      exec bin/core start
    else
      # Pass all arguments to the release command
      echo "Running custom command..."
      exec bin/core "$@"
    fi
    ;;
esac
