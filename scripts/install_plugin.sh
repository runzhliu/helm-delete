#!/bin/sh

# Copied and adapted from https://github.com/chartmuseum/helm-push

set -eu

BINARY="helm-cm-delete"
PLUGIN_NAME="cm-delete"
PLUGIN_DIR="${HELM_PLUGIN_DIR:-$(CDPATH= cd "$(dirname "$0")/.." && pwd)}"

# Skip download if running in development mode.
if [ "${HELM_CM_DELETE_PLUGIN_NO_INSTALL_HOOK:-}" = "1" ]; then
    echo "Development mode: skipping binary download"
    exit 0
fi

# Extract the plugin version and remove the optional YAML quotes. Keeping the
# quotes would produce a non-existent GitHub release tag such as v"0.0.2".
VERSION=$(awk '/^version:/ { print $2; exit }' "${PLUGIN_DIR}/plugin.yaml" | tr -d "\"'")
case "${VERSION}" in
    ""|*[!0-9A-Za-z.+-]*)
        echo "Error: invalid plugin version in ${PLUGIN_DIR}/plugin.yaml: ${VERSION}" >&2
        exit 1
        ;;
esac

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

DOWNLOAD_URL="https://github.com/runzhliu/helm-delete/releases/download/v${VERSION}/${FILENAME}.tar.gz"

TMP_DIR=$(mktemp -d "${TMPDIR:-/tmp}/helm-cm-delete.XXXXXX")
cleanup() {
    rm -rf "${TMP_DIR}"
}
trap cleanup 0 HUP INT TERM

ARCHIVE_PATH="${TMP_DIR}/${FILENAME}.tar.gz"
EXTRACT_DIR="${TMP_DIR}/extract"
mkdir -p "${EXTRACT_DIR}"

echo "Downloading ${BINARY} v${VERSION} for ${OS}/${ARCH}..."

if command -v curl > /dev/null 2>&1; then
    curl -fsSL "${DOWNLOAD_URL}" -o "${ARCHIVE_PATH}"
elif command -v wget > /dev/null 2>&1; then
    wget -q "${DOWNLOAD_URL}" -O "${ARCHIVE_PATH}"
else
    echo "Error: curl or wget is required to download the plugin binary."
    exit 1
fi

tar -xzf "${ARCHIVE_PATH}" -C "${EXTRACT_DIR}"

INSTALLED_BINARY="${BINARY}"
if [ "${OS}" = "windows" ]; then
    INSTALLED_BINARY="${BINARY}.exe"
fi

if [ ! -f "${EXTRACT_DIR}/${INSTALLED_BINARY}" ]; then
    echo "Error: release archive does not contain ${INSTALLED_BINARY}." >&2
    exit 1
fi

mkdir -p "${PLUGIN_DIR}/bin"
mv "${EXTRACT_DIR}/${INSTALLED_BINARY}" "${PLUGIN_DIR}/bin/${INSTALLED_BINARY}"
chmod +x "${PLUGIN_DIR}/bin/${INSTALLED_BINARY}"

# Helm 4 understands the legacy Helm 3 manifest, but reports it as legacy.
# Once the install hook has run, switch to the versioned Helm 4 manifest.
if [ -n "${HELM_MAJOR_VERSION:-}" ]; then
    HELM_MAJOR="${HELM_MAJOR_VERSION}"
else
    HELM_MAJOR=$("${HELM_BIN:-helm}" version --short 2>/dev/null | sed -n 's/^v\([0-9][0-9]*\).*/\1/p')
fi

if [ "${HELM_MAJOR}" = "4" ]; then
    sed "s/@VERSION@/${VERSION}/g" \
        "${PLUGIN_DIR}/scripts/plugin-helm4.yaml.tpl" > "${TMP_DIR}/plugin.yaml"
    mv "${TMP_DIR}/plugin.yaml" "${PLUGIN_DIR}/plugin.yaml"
fi

echo "Installed ${PLUGIN_NAME} plugin successfully."
