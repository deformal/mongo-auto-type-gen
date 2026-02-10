#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="${ROOT_DIR}/dist"

# Clean dist directory to ensure fresh build
echo "Cleaning ${OUT_DIR}..."
rm -rf "${OUT_DIR}"
mkdir -p "${OUT_DIR}"

# Clear Go build cache to prevent stale builds
echo "Clearing Go build cache..."
go clean -cache

build() {
  local goos="$1"
  local goarch="$2"
  local ext="$3"
  local tag="${4:-${goos}}"
  local out="${OUT_DIR}/mongots-${tag}-${goarch}${ext}"

  echo "Building ${out}"
  GOOS="${goos}" GOARCH="${goarch}" go build -o "${out}" ./core/cmd/main.go
}

build "darwin" "arm64" ""
build "darwin" "amd64" ""
build "linux" "amd64" ""
build "linux" "arm64" ""
build "windows" "amd64" ".exe" "win32"

# Verify all binaries were built
echo ""
echo "Build complete! Verifying binaries..."
ls -lh "${OUT_DIR}"/mongots-*
echo ""
echo "Total binaries: $(ls -1 "${OUT_DIR}"/mongots-* | wc -l)"
