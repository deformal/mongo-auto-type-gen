#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${1:-}" ]]; then
  echo "Usage: scripts/release.sh <version>"
  exit 1
fi

VERSION="$1"

npm version "${VERSION}" --no-git-tag-version

git add package.json package-lock.json
git commit -m "release: v${VERSION}"
git tag "v${VERSION}"
git push
git push origin "v${VERSION}"
