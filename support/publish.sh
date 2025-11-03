#!/bin/bash
set -euo pipefail

version="v$(tr -d '[:space:]' <VERSION)"

if [[ -z "$version" ]]; then
	echo "VERSION file is empty." >&2
	exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
	echo "Found pending changes. Preparing to commit before tagging."
	read -rp "Commit message: " commit_message
	if [[ -z "${commit_message// }" ]]; then
		echo "Commit message cannot be empty when committing pending changes." >&2
		exit 1
	fi
	git add -A
	git commit -m "$commit_message"
fi

if git rev-parse "$version" >/dev/null 2>&1; then
	echo "Tag $version already exists." >&2
	exit 1
fi

git tag "$version"
git push origin "$version"
