# Ebiten UI Project Integration Spec

## 목적

이 문서는 실제 프로젝트에서 `ebiten-ui`, `ebiten-debug`, `ebiten-ui-debug`, `ebiten-mcp`를 어떻게 조립할지 구현 직전 수준으로 정리한다.  
목표는 showcase 코드를 복사하지 않고도 새 프로젝트가 같은 UI debug surface를 바로 붙일 수 있게 하는 것이다.

## 대상

- `projects/<name>` 아래의 실제 앱
- `examples`가 아니라 제품 조립 코드

## 기본 조립 구조

### 코어 조합

- `ebiten-ui`
- `ebiten-debug`
- `ebiten-ui-debug`

### 선택 조합

- `ebiten-mcp`
- `tools/ebiten-mcp-server`

## 프로젝트 구조 기준

최소 권장 구조:

- `projects/<name>/cmd/<app>`
- `projects/<name>/internal/app`
- `projects/<name>/internal/ui`
- `projects/<name>/internal/screens`
- `projects/<name>/README.md`

`internal/ui`는 product-specific UI builder와 theme binding을 둔다.  
`ebiten-ui` 자체를 프로젝트 내부로 복제하지 않는다.

## 런타임 조립 기준

프로젝트 app struct는 최소 아래를 가진다.

```go
type App struct {
  runtime     *ebitenui.Runtime
  router      *ebitenui.PageRouter
  uiDebug     *ebitenuidebug.Adapter
  debugBridge *ebitendebug.Bridge
  renderer    *renderer.Renderer
}
```

## 책임 분리

### 프로젝트가 소유하는 것

- 현재 screen/page 조립
- project-specific state
- theme 선택
- debug bridge 활성화 여부

### 라이브러리가 소유하는 것

- 노드 트리와 레이아웃
- 입력 런타임
- debug snapshot/query/capture/input injection

## 디버그 모드 규칙

- core UI library는 debug/product 모드에 따라 달라지지 않는다.
- 바뀌는 것은 `debug attach 여부`뿐이다.

즉:

- `debug build`
  - `ebiten-debug` attach
  - `ebiten-ui-debug` attach
  - overlay optional
- `product build`
  - attach 없음 또는 비활성

## attach API 사용 기준

프로젝트는 아래 callback을 제공해야 한다.

- `CurrentLayout`
- `CurrentViewport`
- `CurrentRuntime`
- `CurrentInput`
- `CurrentFrame`
- `OverlayEnabled`
- `SetOverlay`

이 callback은 앱 struct에 강결합된 helper가 아니라, 현재 상태를 읽는 최소 surface여야 한다.

## screenshot 경로 규칙

- 프로젝트는 `ScreenshotsDir`를 명시적으로 준다.
- 기본 규칙:
  - 저장소 루트 `screenshots/<game-id>/`
- 라이브러리 경로 추론에 의존하지 않는다.

## scene/world/ui summary 기준

프로젝트는 최소 아래 provider를 등록한다.

- frame summary
- scene summary
- world summary
- ui summary

UI 앱이면 scene/world summary에도 최소한 다음은 드러나야 한다.

- current page/screen id
- 주요 navigation/content panel
- 필요 최소 상호작용 상태

## 디버그 테스트 루프

프로젝트 UI 검증 기본 흐름:

1. 앱 실행
2. debug bridge attach
3. `ebiten-mcp-server` attach
4. `get_ui_overview`
5. `query_ui_nodes`
6. `inspect_ui_node`
7. `run_command(ui_click/ui_scroll/ui_type_text)`
8. `capture_ui_screenshot`

## 프로젝트 README 계약

각 프로젝트 README는 최소 아래를 가져야 한다.

- 실행 방법
- debug mode 실행 방법
- screenshot 저장 위치
- MCP 연결 방법
- current UI architecture 요약

## starter 완료 기준

새 프로젝트가 아래를 만족해야 한다.

- 최소 wiring만으로 UI render 가능
- debug bridge attach 가능
- `get_ui_overview` 동작
- `run_command(ui_click)` 동작
- `capture_ui_screenshot` 동작

## 마이그레이션 순서

1. `projects/<name>` 최소 앱 추가
2. `ebiten-ui` page router 또는 screen 구조 연결
3. `ebiten-ui-debug` attach
4. `ebiten-debug` provider 등록
5. MCP inspect/capture 검증
6. project README 정리
