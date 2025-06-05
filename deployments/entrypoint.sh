#!/bin/sh

# Exit immediately if a command exits with a non-zero status.
set -e

# Here you can add any commands that need to run before starting the main application.
# For example:
# - Running database migrations
# - Waiting for other services to be ready
# - Setting up configurations

echo "Entrypoint script started..."
echo "Executing command: $@"

# Execute the command passed as arguments to the script (which is CMD from Dockerfile).
exec "$@"
