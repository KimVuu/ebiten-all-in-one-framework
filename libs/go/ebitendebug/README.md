# ebitendebug

`ebitendebug`는 Ebiten/Ebitengine 앱의 디버그 상태를 loopback HTTP로 노출하는 Go 라이브러리다. 프레임, 씬, 월드, UI 스냅샷과 디버그 명령 실행을 공통 계약으로 제공한다.

## 제공 기능

- `/health`: 게임 식별자와 연결 상태
- `/debug/frame`: 프레임 상태
- `/debug/scene`: 현재 씬 상태
- `/debug/world`: 엔티티 월드 상태
- `/debug/ui`: `semantic`, `layout`, `computed`, `issues`, `inputState`를 포함한 UI 상태
- `/debug/commands`: 등록된 디버그 명령 목록
- `/debug/commands/{name}`: 디버그 명령 실행

## UI 스냅샷

`/debug/ui`는 화면별 `semantic` 정보, 레이아웃 제약, 계산된 bounds, 진단 이슈, 입력 상태를 함께 노출한다. 새 명령 표면은 `validate_ui_layout`, `inspect_ui_node`, `suggest_ui_constraint_fixes`, `set_ui_debug_overlay`, `ui_click`, `ui_scroll`, `ui_type_text`, `ui_key_event` 같은 작업을 main app이 직접 등록해 처리하는 전제를 따른다.

명령 응답은 `success`, `status`, `resolvedTarget`, `queuedFrame`, `reason`, `payload`를 담을 수 있다.

## 사용 방향

- 게임 앱은 `libs/go/ebitendebug`를 직접 붙여 디버그 브리지를 연다.
- MCP 어댑터가 필요하면 [../ebiten-mcp/README.md](../ebiten-mcp/README.md)를 사용한다.

## 검증

```bash
go test ./...
```
