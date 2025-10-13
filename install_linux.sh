#!/usr/bin/env bash
set -euo pipefail

REPO="arsiba/tofulint"

echo "===================================================="
echo "Fetching release version ..."

get_latest_release() {
  headers=()
  if [ -n "${GITHUB_TOKEN:-}" ]; then
      headers=(-H "Authorization: Bearer ${GITHUB_TOKEN}")
  fi
  curl --fail -sS "${headers[@]}" "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name":' \
    | sed -E 's/.*"([^"]+)".*/\1/'
}

if [ -z "${TOFULINT_VERSION:-}" ] || [ "${TOFULINT_VERSION}" = "latest" ]; then
  echo "Determining latest version ..."
  version=$(get_latest_release)
else
  version=${TOFULINT_VERSION}
fi

echo "Version: $version"

BINARY_NAME="tofulint-${version}"

echo "===================================================="
echo "Downloading ${BINARY_NAME} ..."

download_path=$(mktemp -d -t tofulint.XXXXXXXXXX)
download_executable="${download_path}/tofulint"

curl --fail -sSL -o "${download_executable}" \
  "https://github.com/${REPO}/releases/download/${version}/${BINARY_NAME}" || {
    echo "Download failed. Check if ${BINARY_NAME} exists for release ${version}."
    exit 1
  }

chmod +x "${download_executable}"

echo "===================================================="
echo "Installing ${BINARY_NAME} ..."

dest="${TOFULINT_INSTALL_PATH:-/usr/local/bin}"

if [[ -w "$dest" ]]; then
  SUDO=""
else
  SUDO="sudo"
fi

$SUDO mkdir -p "$dest"
$SUDO install -c -v "${download_executable}" "$dest/tofulint"

echo "===================================================="
echo "Cleaning up ..."
rm -rf "${download_path}"

echo "===================================================="
echo "tofulint has been installed to ${dest}"
"${dest}/tofulint" -v || echo "Version could not be displayed"

echo "===================================================="
echo "Start by calling 'tofulint' in your terminal"
