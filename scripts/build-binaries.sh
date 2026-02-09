#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="${ROOT_DIR}/dist"

mkdir -p "${OUT_DIR}"

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
