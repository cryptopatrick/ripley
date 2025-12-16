#!/bin/bash
# Script to run the Ripley daemon

set -e

# Change to project root
cd "$(dirname "$0")/.."

echo "=== Ripley Daemon Launcher ==="
echo ""

# Check if Claude CLI is installed
if ! command -v claude &> /dev/null; then
    echo "Error: Claude CLI not found in PATH"
    echo "Please install the Claude CLI: https://github.com/anthropics/claude-code"
    exit 1
fi

# Build if needed
if [ ! -f "./ripleyd" ]; then
    echo "Building daemon..."
    make build
    echo ""
fi

# Check for config file
if [ ! -f "./config.yaml" ]; then
    echo "Warning: config.yaml not found, using defaults"
    echo "Tip: Copy config.yaml.example to config.yaml and customize"
    echo ""
fi

# Run daemon
echo "Starting Ripley daemon..."
echo ""
./ripleyd
