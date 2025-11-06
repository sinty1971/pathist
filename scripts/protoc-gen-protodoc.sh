#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
REPO_ROOT=$(cd -- "$SCRIPT_DIR/.." && pwd)
GOCACHE_DIR="$REPO_ROOT/.gocache"
mkdir -p "$GOCACHE_DIR"

exec env GOCACHE="$GOCACHE_DIR" go run "$REPO_ROOT/backend-grpc/cmd/protodoc" "$@"
