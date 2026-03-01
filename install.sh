#!/usr/bin/env bash

set -e

REPO="isayme/port-detector"
INSTALL_DIR="/usr/local/bin"
BIN_NAME="port-detector"

echo "Detecting OS and architecture..."

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    linux*)   OS="linux" ;;
    darwin*)  OS="darwin" ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Detect ARCH
ARCH=$(uname -m)
case "$ARCH" in
    x86_64 | amd64)
        ARCH="amd64"
        ;;
    aarch64 | arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "OS: $OS"
echo "ARCH: $ARCH"

# Get latest version from GitHub
echo "Fetching latest release info..."
LATEST_VERSION=$(curl -s https://api.github.com/repos/${REPO}/releases/latest | grep tag_name | cut -d '"' -f4)

if [ -z "$LATEST_VERSION" ]; then
    echo "Failed to fetch latest version."
    exit 1
fi

echo "Latest version: $LATEST_VERSION"

# Construct download URL
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BIN_NAME}-${OS}-${ARCH}"

echo "Downloading from:"
echo "$DOWNLOAD_URL"

TMP_FILE="/tmp/${BIN_NAME}"

curl -L -o "$TMP_FILE" "$DOWNLOAD_URL"

chmod +x "$TMP_FILE"

echo "Installing to ${INSTALL_DIR} (may require sudo)..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_FILE" "${INSTALL_DIR}/${BIN_NAME}"
else
    sudo mv "$TMP_FILE" "${INSTALL_DIR}/${BIN_NAME}"
fi

echo "Installation complete."
echo "Run:"
echo "  ${BIN_NAME} -h"