#!/usr/bin/env bash
set -euo pipefail

root="$(git rev-parse --show-toplevel)"
cd "$root"

profile_bin="$root/.devenv/profile/bin"
if [[ ! -d "$profile_bin" ]]; then
  echo "❌ Missing $profile_bin — run: devenv shell (or direnv allow)" >&2
  exit 1
fi

export PATH="$root/bin:$profile_bin:$PATH"
export GOROOT="$("$profile_bin/go" env GOROOT)"
export GOTOOLCHAIN=local

exec task check
