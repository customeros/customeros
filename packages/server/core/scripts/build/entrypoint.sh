#!/bin/sh
set -e

# Run migrations before starting the app
bin/realtime eval "Realtime.Release.migrate()"

# Start the app
exec bin/realtime start
