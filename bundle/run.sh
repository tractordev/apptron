#!/bin/bash

# Ensure script is run from its containing directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ "$PWD" != "$SCRIPT_DIR" ]; then
    echo "Error: This script must be run from its containing directory: $SCRIPT_DIR" >&2
    exit 1
fi

# Ensure kernel image exists
if [ ! -f "$SCRIPT_DIR/../kernel/bzImage" ]; then
    echo "Error: Kernel image not found at $SCRIPT_DIR/../kernel/bzImage" >&2
    exit 1
fi

# Ensure bundle exists
if [ ! -d "$SCRIPT_DIR/bundle" ]; then
    make
fi


# Function to cleanup background process on exit
cleanup() {
    echo "Cleaning up..."
    if [ ! -z "$BG_PID" ] && kill -0 $BG_PID 2>/dev/null; then
        echo "Killing background process (PID: $BG_PID)"
        kill $BG_PID
        # Give it time to terminate gracefully, then force kill if needed
        sleep 1
        if kill -0 $BG_PID 2>/dev/null; then
            kill -9 $BG_PID
        fi
    fi
    exit
}

# Set up trap to catch various exit signals
# This ensures cleanup happens on normal exit, Ctrl+C, terminal close, etc.
trap cleanup EXIT INT TERM HUP


echo "Starting wanix serve..."
wanix serve &
BG_PID=$!

# Optional: Wait a moment for background program to initialize
sleep 1

# Check if background process started successfully
if ! kill -0 $BG_PID 2>/dev/null; then
    echo "Error: Background process failed to start"
    exit 1
fi

echo "Background wanix serve started with PID: $BG_PID"

echo "Starting v86-system..."
v86-system -kernel ../kernel/bzImage -netdev user,type=virtio,relay_url=ws://localhost:7654/.well-known/ethernet -virtfs proxy,ws://localhost:7654 -append "console=ttyS0 init=/bin/init rw root=host9p rootfstype=9p rootflags=trans=virtio,version=9p2000.L,aname=bundle/rootfs"

# The cleanup function will be called automatically when the script exits