#!/bin/bash
set -euo pipefail

version=$(tr -d '[:space:]' <VERSION)

if [[ -z "$version" ]]; then
	echo "VERSION file is empty." >&2
	exit 1
fi

if git rev-parse "$version" >/dev/null 2>&1; then
	echo "Tag $version already exists." >&2
	exit 1
fi

git tag "$version"
