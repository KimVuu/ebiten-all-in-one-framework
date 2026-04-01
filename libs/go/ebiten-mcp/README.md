# ebiten-mcp

`ebiten-mcp`는 `ebitendebug` loopback HTTP 브리지를 MCP stdio 서버 형태로 감싸는 Go 라이브러리다. 게임 상태 조회와 디버그 명령 실행을 MCP 툴 호출로 연결한다.

내부 서버 구현은 `github.com/modelcontextprotocol/go-sdk/mcp`를 사용한다. tool 이름과 bridge surface는 유지하고, stdio transport와 tool registration만 공식 SDK 위로 올렸다.

지금은 두 transport를 지원한다.

- `ServeStdio(...)`: 로컬 MCP host용 stdio 서버
- `StreamableHTTPHandler(...)`: HTTP MCP client용 streamable HTTP handler

## 제공 툴

- `game_health`
- `get_frame_state`
- `get_scene_state`
- `get_world_state`
- `get_ui_state`
- `get_ui_overview`
- `query_ui_nodes`
- `inspect_ui_node`
- `list_ui_issues`
- `capture_ui_screenshot`
- `list_commands`
- `run_command`

## UI 디버그 표면

- `get_ui_state`는 `semantic`, `layout`, `computed`, `issues`, `inputState`를 포함한 full UI 스냅샷을 그대로 전달한다. legacy/high token cost 경로다.
- `get_ui_overview`, `query_ui_nodes`, `inspect_ui_node`, `list_ui_issues`는 design/layout testing용 compact 표면이다.
- `capture_ui_screenshot`는 PNG artifact metadata와 절대 경로만 반환하고, 이미지 bytes는 inline으로 싣지 않는다.
- `run_command`는 `validate_ui_layout`, `inspect_ui_node`, `suggest_ui_constraint_fixes`, `set_ui_debug_overlay`, `ui_click`, `ui_scroll`, `ui_type_text`, `ui_key_event` 같은 명령을 전달하는 용도다.

권장 루프는 `get_ui_overview -> query_ui_nodes -> inspect_ui_node -> run_command(ui_*) -> capture_ui_screenshot` 순서다.

## 사용 방향

- 게임 앱은 `libs/go/ebitendebug`로 디버그 브리지를 연다.
- MCP 호스트와 연결할 때 이 라이브러리를 tool runner가 소비한다.
- stdio 연결은 `go-sdk`의 `mcp.Server`와 `mcp.IOTransport`로 처리한다.
- HTTP 연결은 `go-sdk`의 `mcp.NewStreamableHTTPHandler(...)`를 사용한다.
- 실제 실행 엔트리는 `tools/ebiten-mcp-server`에 둔다.

## 검증

```bash
go test ./...
```
