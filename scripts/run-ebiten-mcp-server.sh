#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ADDR="${EBITEN_DEBUG_ADDR:-127.0.0.1:47831}"
TRANSPORT="${EBITEN_MCP_TRANSPORT:-stdio}"
LISTEN_ADDR="${EBITEN_MCP_LISTEN_ADDR:-127.0.0.1:47830}"
HTTP_PATH="${EBITEN_MCP_HTTP_PATH:-/mcp}"

cd "${ROOT_DIR}/tools/ebiten-mcp-server"
exec env \
  EBITEN_DEBUG_ADDR="${ADDR}" \
  EBITEN_MCP_TRANSPORT="${TRANSPORT}" \
  EBITEN_MCP_LISTEN_ADDR="${LISTEN_ADDR}" \
  EBITEN_MCP_HTTP_PATH="${HTTP_PATH}" \
  go run . \
    --addr "${ADDR}" \
    --transport "${TRANSPORT}" \
    --listen "${LISTEN_ADDR}" \
    --path "${HTTP_PATH}" \
    "$@"
