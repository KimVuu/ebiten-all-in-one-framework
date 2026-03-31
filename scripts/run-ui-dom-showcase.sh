#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ADDR="${EBITEN_DEBUG_ADDR:-127.0.0.1:47831}"

cd "${ROOT_DIR}/examples/go/ui-dom-showcase"
exec env EBITEN_DEBUG_MODE=1 EBITEN_DEBUG_ADDR="${ADDR}" go run . "$@"
