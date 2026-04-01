# ebiten-ui-debug

`libs/go/ebiten-ui-debug`는 `ebiten-ui`와 `ebiten-debug`를 연결하는 재사용 가능한 UI 디버그 어댑터다. 앱이 현재 layout, viewport, runtime, input, frame, overlay state를 callback으로 제공하면, 이 라이브러리가 compact UI inspect surface와 debug input 흐름을 공통으로 붙여준다.

## 제공 기능

- `LayoutNode -> UISnapshot` 변환
- 저토큰 UI 표면
- `overview`, `query`, `inspect`, `issues`, `capture`
- UI screenshot artifact 저장
- debug overlay 그리기
- reusable debug input queue
- 기본 UI command 등록
- `screenshots/<game-id>/` 기준 artifact 경로 관리

## 핵심 API

- `NewAdapter(config, callbacks)`
- `Attach(bridge)`
- `ApplyQueuedInput(frame, dom, runtime, layout, input)`
- `DrawOverlay(screen, layout, enabled)`
- `UISnapshot()`
- `UIOverview()`
- `UIQuery(request)`
- `UINodeDetail(request)`
- `UIIssues(request)`
- `UICapture(request)`
- `Artifact(id)`

## 사용 방향

앱은 코어 UI를 `ebiten-ui`로 만들고, 디버그 브리지는 `ebiten-debug`로 열고, UI 디버그 표면은 `ebiten-ui-debug`로 붙인다.

```go
adapter := ebitenuidebug.NewAdapter(ebitenuidebug.Config{
	GameID: "my-game",
}, ebitenuidebug.Callbacks{
	CurrentLayout:   currentLayout,
	CurrentViewport: currentViewport,
	CurrentRuntime:  currentRuntime,
	CurrentInput:    currentInput,
	CurrentFrame:    currentFrame,
	OverlayEnabled:  overlayEnabled,
	SetOverlay:      setOverlay,
})

bridge := ebitendebug.New(ebitendebug.Config{
	Enabled: true,
	Addr:    "127.0.0.1:47831",
	GameID:  "my-game",
	Version: "v1",
})

adapter.Attach(bridge)
```

## 명령 표면

기본 등록 명령은 아래와 같다.

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

입력은 즉시 핸들러 호출이 아니라 프레임 큐로 들어가며, 다음 update에서 적용된다.

## 예제

- 통합 UI 예제: [examples/go/ebiten-ui-showcase/README.md](../../../examples/go/ebiten-ui-showcase/README.md)
- 브리지 기본 예제: [examples/go/ebiten-debug-bridge/README.md](../../../examples/go/ebiten-debug-bridge/README.md)
