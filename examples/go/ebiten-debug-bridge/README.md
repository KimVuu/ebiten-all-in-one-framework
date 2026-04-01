# Debug Bridge Example

`examples/go/ebiten-debug-bridge`는 `ebiten-debug`를 붙인 작은 Ebiten 샘플 게임이다.  
디버그 모드일 때만 로컬 HTTP 서버를 열고, 외부 `ebiten-mcp` 소비 코드가 이 주소에 attach 할 수 있다.

## 실행 방법

일반 실행:

```bash
go run .
```

디버그 브리지 활성화 실행:

```bash
EBITEN_DEBUG_MODE=1 EBITEN_DEBUG_ADDR=127.0.0.1:47831 go run .
```

루트 스크립트도 사용할 수 있다.

```bash
./scripts/run-debug-bridge.sh
```

## 노출되는 상태

- 프레임, 틱, FPS, TPS, pause 상태
- 현재 씬과 등록된 씬 목록
- 샘플 엔티티 월드 스냅샷
- 샘플 UI tree, 텍스트, 레이아웃 스냅샷
- 디버그 커맨드
  - `pause.toggle`
  - `scene.switch`
  - `entity.visibility.toggle`

## 수동 검증

1. `./scripts/run-debug-bridge.sh`
2. 다른 터미널에서 `./scripts/run-ebiten-mcp-server.sh`
3. MCP client에서 `game_health`, `get_world_state`, `get_ui_state`, `run_command`를 호출한다.
4. `scene.switch`, `pause.toggle` 실행 결과가 게임 상태에 반영되는지 확인한다.

## 검증

- `go test ./...`
