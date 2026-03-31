# Ebiten MCP Server

`tools/ebiten-mcp-server`는 `libs/go/ebiten-mcp`를 사용해 실행 중인 `ebitendebug` 브리지에 attach 하는 stdio MCP 서버다.

내부 MCP 서버는 `github.com/modelcontextprotocol/go-sdk/mcp` 기반이다. 외부에 보이는 tool 표면은 그대로 유지하고, 공식 SDK transport와 session lifecycle을 사용한다.

## 제공 tool

- `game_health`
- `get_frame_state`
- `get_scene_state`
- `get_world_state`
- `get_ui_state`
- `list_commands`
- `run_command`

## 실행 방법

기본 주소는 `127.0.0.1:47831`이다.

```bash
go run . --addr 127.0.0.1:47831
```

또는 환경변수를 사용할 수 있다.

```bash
EBITEN_DEBUG_ADDR=127.0.0.1:47831 go run .
```

루트 스크립트도 사용할 수 있다.

```bash
./scripts/run-ebiten-mcp-server.sh
```

## 의존 라이브러리

- `libs/go/ebiten-mcp`
- `libs/go/ebitendebug`
- `github.com/modelcontextprotocol/go-sdk/mcp`

## 검증

- `go test ./...`
