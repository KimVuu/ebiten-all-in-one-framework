# ebiten Ebiten UI Showcase

`examples/go/ebiten-ui-showcase`는 `libs/go/ebiten-ui`의 코어 태그, 확장 컴포넌트, 게임 UI 프리셋을 한 화면에 합쳐 보여주는 단일 Ebiten 예제다. 동시에 `libs/go/ebiten-ui-debug`를 소비하는 기준 통합 예제이기도 하다.

## 포함 내용

- 코어 DOM 태그
- `Image`, `TextBlock`, `Spacer`, `Stack`, `ScrollView`
- 입력/상태 컴포넌트
- 데이터/오버레이 컴포넌트
- `ebitenui/prefabs` 게임 UI 프리셋

## 실행

```bash
cd examples/go/ebiten-ui-showcase
go run .
```

디버그 브리지 활성화 실행:

```bash
cd examples/go/ebiten-ui-showcase
EBITEN_DEBUG_MODE=1 EBITEN_DEBUG_ADDR=127.0.0.1:47832 go run .
```

## 테스트

```bash
cd examples/go/ebiten-ui-showcase
go test ./...
```

## 목적

- 선언형 DOM 트리 조립 예시 제공
- `ebiten-ui` 레이아웃과 Ebiten 렌더러 연결 예시 제공
- 태그, 컴포넌트, 프리셋을 한 번에 검증하는 단일 기준 예제 제공
- `ebiten-ui-debug`로 UI tree와 주요 섹션 레이아웃을 디버그 스냅샷으로 노출하는 기준 예제 제공
- `ebitenui.PageLayout`와 런타임 기반 wheel scroll을 함께 검증하는 예제 제공
- nested layout constraint, keyboard focus/navigation, 디버그 입력 주입을 함께 검증하는 예제 제공

## 디버그 브리지

디버그 모드일 때 loopback HTTP 브리지를 열고, `ebiten-ui-debug` 어댑터가 아래 UI 표면과 명령을 등록한다.

- `/health`
- `/debug/frame`
- `/debug/scene`
- `/debug/world`
- `/debug/ui`
- `/debug/ui/overview`
- `/debug/ui/query`
- `/debug/ui/node/{id}`
- `/debug/ui/issues`
- `/debug/ui/capture`
- `/debug/ui/artifacts/{artifactId}`

디자인 테스트 기본 경로는 full tree dump 대신 아래 순서를 쓴다.

1. `/debug/ui/overview`
2. `/debug/ui/query`
3. `/debug/ui/node/{id}`
4. `run_command`로 `ui_click`, `ui_scroll`, `ui_type_text`
5. `/debug/ui/capture`

`/debug/ui/capture`는 PNG artifact metadata만 반환하고, 이미지 bytes는 inline으로 싣지 않는다.
생성된 PNG는 저장소 루트의 [screenshots](/Users/kimyechan/Develop/Game/Ebiten/ebtien-aio-framework/screenshots) 아래에 쌓인다.

추가로 `run_command`를 통해 아래 UI 디버그 명령을 사용할 수 있다.

- `validate_ui_layout`
- `inspect_ui_node`
- `suggest_ui_constraint_fixes`
- `set_ui_debug_overlay`
- `ui_click`
- `ui_pointer_move`
- `ui_pointer_down`
- `ui_pointer_up`
- `ui_scroll`
- `ui_focus_node`
- `ui_type_text`
- `ui_key_event`
- `ui_clear_input_queue`

입력은 프레임 큐로 들어가며, 클릭/스크롤/텍스트/키 입력은 `node_id` 또는 좌표 기준으로 재현된다. 오버레이는 bounds, anchor/pivot, invalid state, clickable rect, focus/hover 상태를 표시한다.
