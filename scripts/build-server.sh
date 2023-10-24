#!/usr/bin/env bash
set -e

SCRIPT_DIR=$(cd -P -- "$(dirname -- "$0")" && pwd -P)

echo "Building server..."
GOOS=linux GOARCH=amd64 go build -o potok "$SCRIPT_DIR/../cmd/server"
