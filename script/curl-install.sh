#!/bin/bash
set -e

# Detect OS and Arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map to match release filenames
if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" == "arm64" ] || [ "$ARCH" == "aarch64" ]; then
    ARCH="arm64"
fi

if [ "$OS" == "darwin" ]; then
    OS="darwin"
elif [ "$OS" == "linux" ]; then
    OS="linux"
elif [ "$OS" == "mingw64_nt" ] || [ "$OS" == "msys_nt" ]; then
    OS="windows"
fi

# We need a tag
if [ -z "$1" ]; then
    echo "Usage: curl -sL https://raw.githubusercontent.com/jefflunt/breakdown/main/script/curl-install.sh | bash -s <tag>"
    echo "Example tag: b1"
    exit 1
fi
TAG=$1

BINARY_NAME="breakdown-${OS}-${ARCH}"
if [ "$OS" == "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
fi

URL="https://github.com/jefflunt/breakdown/releases/download/${TAG}/${BINARY_NAME}"

echo "Downloading ${BINARY_NAME} from ${URL}..."

curl -sL "$URL" -o breakdown

chmod +x breakdown

echo "Installed breakdown to ./breakdown"
echo "Move it to your PATH to use it."
