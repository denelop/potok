#!/usr/bin/env bash
set -e

SCRIPT_DIR=$(cd -P -- "$(dirname -- "$0")" && pwd -P)

echo "Deploying server..."
scp "$SCRIPT_DIR/../potok-server" "$1"
