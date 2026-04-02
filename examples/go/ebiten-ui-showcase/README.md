# ebiten Ebiten UI Showcase

`examples/go/ebiten-ui-showcase`는 `libs/go/ebiten-ui`의 코어 태그, 확장 컴포넌트, 게임 UI 프리셋을 문서형 다중 페이지 구조로 보여주는 Ebiten 예제다. 동시에 `libs/go/ebiten-ui-debug`를 소비하는 기준 통합 예제이기도 하다.

## 화면 구조

- 상단 헤더
- 좌측 사이드바 navigation
- 우측 상세 패널
- 그룹 페이지와 leaf 페이지
- 사이드바는 현재 그룹만 펼치는 collapsed navigation 구조
- 각 페이지별 설명, 실제 데모, 사용법, 코드 예제 문자열
- 헤더의 global theme preset switcher
- detail 영역의 live state 패널
- `foundations/theme` 페이지에서 default theme와 override theme 비교
- `reactive/*` 페이지에서 `Ref`와 `Computed` 기반 상태 흐름 설명
- 대표 입력/상태 페이지에서 `Ref` 기반 binding demo 사용

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
- `PageRouter`와 `PageScreen` 기반 문서형 UI 조립 예시 제공
- 태그, 컴포넌트, 프리셋을 페이지별로 검증하는 기준 예제 제공
- `ebiten-ui-debug`로 UI tree와 주요 섹션 레이아웃을 디버그 스냅샷으로 노출하는 기준 예제 제공
- sidebar/detail 별도 scroll과 페이지 전환을 함께 검증하는 예제 제공
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

현재 선택된 페이지 정보는 scene summary와 UI snapshot root props의 `currentPageID`에서도 확인할 수 있다.
현재 theme preset은 scene summary와 UI snapshot root props의 `themePreset`에서도 확인할 수 있다.

`/debug/ui/capture`는 PNG artifact metadata만 반환하고, 이미지 bytes는 inline으로 싣지 않는다.
이 예제는 `Config.ScreenshotsDir`를 명시해서 생성된 PNG를 저장소 루트의 [screenshots](/Users/kimyechan/Develop/Game/Ebiten/ebtien-aio-framework/screenshots) 아래에 고정 저장한다.

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

입력은 프레임 큐로 들어가며, 클릭/스크롤/텍스트/키 입력은 `node_id` 또는 좌표 기준으로 재현된다. 오버레이는 bounds, anchor/pivot, invalid state, clickable rect, focus/hover 상태를 표시한다. 기본값은 `off`이고, 필요할 때만 `set_ui_debug_overlay`로 켠다.
