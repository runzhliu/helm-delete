#!/usr/bin/env bash

# Copied and adapted from https://github.com/chartmuseum/helm-push

set -e

BINARY="helm-cm-delete"
PLUGIN_NAME="cm-delete"

# Skip download if running in development mode.
if [ "${HELM_CM_DELETE_PLUGIN_NO_INSTALL_HOOK}" = "1" ]; then
    echo "Development mode: skipping binary download"
    exit 0
fi

# Extract the plugin version from plugin.yaml.
VERSION=$(grep "^version:" "${HELM_PLUGIN_DIR}/plugin.yaml" | awk '{print $2}')

# Detect OS.
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "${OS}" in
    darwin)  OS="darwin" ;;
    linux)   OS="linux" ;;
    mingw*|msys*|cygwin*) OS="windows" ;;
    *)
        echo "Unsupported OS: ${OS}"
        exit 1
        ;;
esac

# Detect architecture.
ARCH=$(uname -m)
case "${ARCH}" in
    x86_64|amd64)   ARCH="amd64" ;;
    aarch64|arm64)  ARCH="arm64" ;;
    *)
        echo "Unsupported architecture: ${ARCH}"
        exit 1
        ;;
esac

FILENAME="${BINARY}_${OS}_${ARCH}"
if [ "${OS}" = "windows" ]; then
    FILENAME="${FILENAME}.exe"
fi

DOWNLOAD_URL="https://github.com/chartmuseum/helm-cm-delete/releases/download/v${VERSION}/${FILENAME}.tar.gz"

mkdir -p "${HELM_PLUGIN_DIR}/bin"
cd "${HELM_PLUGIN_DIR}/bin"

echo "Downloading ${BINARY} v${VERSION} for ${OS}/${ARCH}..."

if command -v curl > /dev/null 2>&1; then
    curl -sSfL "${DOWNLOAD_URL}" | tar xz
elif command -v wget > /dev/null 2>&1; then
    wget -qO- "${DOWNLOAD_URL}" | tar xz
else
    echo "Error: curl or wget is required to download the plugin binary."
    exit 1
fi

# Ensure the binary is executable.
chmod +x "${HELM_PLUGIN_DIR}/bin/${BINARY}"

echo "Installed ${PLUGIN_NAME} plugin successfully."
