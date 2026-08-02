#!/bin/sh

set -eu

ROOT_DIR=$(CDPATH= cd "$(dirname "$0")/.." && pwd)
TEST_DIR=$(mktemp -d "${TMPDIR:-/tmp}/helm-cm-delete-install-test.XXXXXX")
cleanup() {
	rm -rf "${TEST_DIR}"
}
trap cleanup 0 HUP INT TERM

PLUGIN_DIR="${TEST_DIR}/plugin"
mkdir -p "${PLUGIN_DIR}/scripts" "${TEST_DIR}/archive" "${TEST_DIR}/fake-bin"
cp "${ROOT_DIR}/plugin.yaml" "${PLUGIN_DIR}/plugin.yaml"
cp "${ROOT_DIR}/scripts/install_plugin.sh" "${PLUGIN_DIR}/scripts/install_plugin.sh"
cp "${ROOT_DIR}/scripts/plugin-helm4.yaml.tpl" "${PLUGIN_DIR}/scripts/plugin-helm4.yaml.tpl"

# Build a release-shaped fixture without making a network request.
printf '#!/bin/sh\necho helm-cm-delete test binary\n' > "${TEST_DIR}/archive/helm-cm-delete"
chmod +x "${TEST_DIR}/archive/helm-cm-delete"
tar -czf "${TEST_DIR}/release.tar.gz" -C "${TEST_DIR}/archive" helm-cm-delete

cat > "${TEST_DIR}/fake-bin/curl" <<'EOF'
#!/bin/sh
set -eu

url=""
output=""
while [ "$#" -gt 0 ]; do
	case "$1" in
		-o)
			shift
			output="$1"
			;;
		-*) ;;
		*) url="$1" ;;
	esac
	shift
done

if [ "${url}" != "${TEST_DOWNLOAD_URL}" ]; then
	echo "unexpected download URL: ${url}" >&2
	exit 1
fi
cp "${TEST_ARCHIVE}" "${output}"
EOF
chmod +x "${TEST_DIR}/fake-bin/curl"

case "$(uname -s)" in
	Darwin) EXPECTED_OS="darwin" ;;
	Linux) EXPECTED_OS="linux" ;;
	*)
		echo "unsupported test OS" >&2
		exit 1
		;;
esac

case "$(uname -m)" in
	x86_64|amd64) EXPECTED_ARCH="amd64" ;;
	aarch64|arm64) EXPECTED_ARCH="arm64" ;;
	*)
		echo "unsupported test architecture" >&2
		exit 1
		;;
esac

export TEST_ARCHIVE="${TEST_DIR}/release.tar.gz"
export TEST_DOWNLOAD_URL="https://github.com/runzhliu/helm-delete/releases/download/v0.0.3/helm-cm-delete_${EXPECTED_OS}_${EXPECTED_ARCH}.tar.gz"

PATH="${TEST_DIR}/fake-bin:${PATH}" \
	HELM_PLUGIN_DIR="${PLUGIN_DIR}" \
	HELM_MAJOR_VERSION=4 \
	sh "${PLUGIN_DIR}/scripts/install_plugin.sh"

test -x "${PLUGIN_DIR}/bin/helm-cm-delete"
cmp "${TEST_DIR}/archive/helm-cm-delete" "${PLUGIN_DIR}/bin/helm-cm-delete"
grep -q '^apiVersion: v1$' "${PLUGIN_DIR}/plugin.yaml"
grep -q '^version: "0.0.3"$' "${PLUGIN_DIR}/plugin.yaml"

if command -v helm > /dev/null 2>&1; then
	HELM_MAJOR=$(helm version --short 2>/dev/null | sed -n 's/^v\([0-9][0-9]*\).*/\1/p')
	if [ "${HELM_MAJOR}" = "4" ]; then
		mkdir -p "${TEST_DIR}/helm-data/plugins"
		ln -s "${PLUGIN_DIR}" "${TEST_DIR}/helm-data/plugins/helm-delete"
		HELM_DATA_HOME="${TEST_DIR}/helm-data" helm plugin list | grep -q 'cm-delete.*v1'
		HELM_DATA_HOME="${TEST_DIR}/helm-data" helm cm-delete --help | grep -q 'helm-cm-delete test binary'
	else
		echo "Skipping Helm 4 runtime test (found Helm ${HELM_MAJOR:-unknown})."
	fi
fi

echo "Install script test passed."
