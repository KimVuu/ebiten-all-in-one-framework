# ebiten Ebiten UI

`libs/go/ebiten-ui`은 Ebiten에서 선언형으로 UI 트리를 만들기 위한 Go 라이브러리다. HTML 태그처럼 `Div`, `Section`, `Text` 같은 노드를 조합하고, React처럼 함수가 노드를 반환하는 방식으로 화면 구조를 만들 수 있게 두는 것을 목표로 한다.

## 현재 제공 범위

- 선언형 노드 트리
- `div`, `section`, `header`, `main`, `footer`, `button`, `span`, `text`, `img`, `text-block`, `spacer`, `stack`, `scroll-view` 노드 함수
- `Style` 기반의 `row`/`column` 레이아웃
- `Semantic`/`Layout` 메타 기반의 제약 레이아웃
- nested `anchored`/`grid` 제약과 parent-bounds clamp
- `grid`의 `row/column span`, auto-placement, `justify/align` 계열 track 규칙
- `px`와 `fill` 길이 단위
- `padding`, `gap`, 배경색, 테두리, 텍스트 색상, 텍스트 정렬
- `Image`, `Icon`, `TextBlock`, `Textarea`, `Spacer`, `Stack`, `ScrollView`, `InteractiveButton`
- `Checkbox`, `Toggle`, `Slider`, `Scrollbar`, `Dropdown`, `InputField`, `RadioGroup`, `Stepper`
- `ProgressBar`, `Divider`, `Grid`, `List`, `VirtualList`, `Modal`, `Tooltip`, `ContextMenu`, `Tabs`, `Accordion`, `Badge`, `Chip`
- `Dialog`, `HUDBar`, `InventoryGrid`, `PauseMenu`, `SettingsPanel`, `Tooltip` 프리셋
- `PageLayout` 기반의 고정 헤더 + 스크롤 본문 레이아웃 헬퍼
- `PageRoute`, `PageRouter`, `PageScreen` 기반의 중첩 페이지 탐색 조립
- ID 기반 DOM 조회
- `LayoutNode`의 `ParentID`, `ContentBounds`, `ClipRect`, `ClickableRect`, `Overflow` 계산 필드
- `ValidateLayout` 기반의 레이아웃 검증과 constraint patch 제안
- `Runtime`, `InputSnapshot`, `EventHandlers` 기반 상호작용 런타임
- `Tab`, `Escape`, arrow key, `Home/End`, `SelectAll`, shortcut 문자열 기반 포커스/텍스트 편집
- `OnClick`, `OnChange`, `OnSubmit` 콜백이 연결된 주요 입력 컴포넌트
- Ebiten 렌더러 서브패키지

`ebiten-ui`는 코어 라이브러리다. 디버그 브리지 연결, compact UI inspect/query, debug overlay, screenshot capture, MCP 친화 command surface는 [../ebiten-ui-debug/README.md](../ebiten-ui-debug/README.md)에서 담당한다.

## 기본 사용 예시

```go
package main

import (
	"image/color"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
	"github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui/renderer"
)

func buildUI() *ebitenui.DOM {
	return ebitenui.New(
		ebitenui.Div(ebitenui.Props{
			ID: "app",
			Style: ebitenui.Style{
				Width:           ebitenui.Fill(),
				Height:          ebitenui.Fill(),
				Direction:       ebitenui.Column,
				Padding:         ebitenui.All(16),
				Gap:             12,
				BackgroundColor: color.RGBA{R: 18, G: 24, B: 33, A: 255},
			},
		},
			ebitenui.Header(ebitenui.Props{},
				ebitenui.Text("Ebiten UI", ebitenui.Props{
					ID: "title",
					Style: ebitenui.Style{
						Color: color.RGBA{R: 245, G: 247, B: 250, A: 255},
					},
				}),
			),
			ebitenui.Section(ebitenui.Props{
				Style: ebitenui.Style{
					Direction:       ebitenui.Row,
					Gap:             8,
					BackgroundColor: color.RGBA{R: 32, G: 40, B: 54, A: 255},
					Padding:         ebitenui.All(12),
				},
			},
				ebitenui.Div(ebitenui.Props{
					Style: ebitenui.Style{
						Width:           ebitenui.Px(48),
						Height:          ebitenui.Px(48),
						BackgroundColor: color.RGBA{R: 84, G: 136, B: 255, A: 255},
					},
				}),
				ebitenui.Text("Build UI like a DOM tree.", ebitenui.Props{
					Style: ebitenui.Style{
						Color: color.RGBA{R: 220, G: 226, B: 235, A: 255},
					},
				}),
			),
		),
	)
}

func drawUI(screenWidth, screenHeight float64) {
	dom := buildUI()
	renderer := renderer.New()
	_ = renderer

	_ = dom.Layout(ebitenui.Viewport{
		Width:  screenWidth,
		Height: screenHeight,
	})
}
```

렌더링할 때는 `renderer.Renderer`로 `*ebiten.Image`에 그리면 된다.

## 런타임 계층

`ebiten-ui`은 이제 정적 노드 트리만이 아니라 입력 런타임도 제공한다.

- `Runtime`: hover, press, focus, click, text input, submit, scroll dispatch
- `InputSnapshot`: 포인터 위치, 버튼 상태, 텍스트 입력, backspace, submit, scroll 입력, tab/escape/arrow/modifier 키 입력
- `EventHandlers`: 노드 단위 이벤트 핸들러

대표 흐름:

```go
runtime := ebitenui.NewRuntime()
layout := runtime.Update(dom, ebitenui.Viewport{
	Width:  1280,
	Height: 720,
}, ebitenui.InputSnapshot{
	PointerX:    mouseX,
	PointerY:    mouseY,
	PointerDown: mouseDown,
	Text:        typedText,
	Backspace:   backspacePressed,
	Submit:      submitPressed,
})

_ = layout
```

입력 컴포넌트는 `OnChange`, `OnSubmit`, `OnOpenChange`, `OnSelect` 같은 콜백을 받아 실제 기능 연결이 가능하다.

기본 런타임은 `Tab` 포커스 이동, `Escape` 포커스 해제, `Enter` 기반 focused button 활성화, arrow 기반 scroll dispatch까지 처리한다.

페이지 전체 스크롤이 필요한 화면은 `PageLayout`로 고정 헤더와 `ScrollView` 본문을 함께 구성하고, 문서형 UI나 툴형 화면은 `PageRouter + PageScreen`으로 좌측 navigation과 우측 detail panel을 조립하는 방식을 기본값으로 둔다.

레이아웃 계산 결과를 디버그하거나 AI가 검사할 때는 `ValidateLayout(layout, viewport, opts)`를 사용하면 된다. 반환값은 `LayoutIssue`와 `ConstraintPatch` 중심으로 구성되어 있어 절대좌표가 아니라 제약 수정 단위로 다룰 수 있다.

## 패키지 구성

- `ebitenui`: 노드, 스타일, 레이아웃, DOM 조회
- `ebitenui` router surface: `PageRoute`, `PageRouter`, `PageScreen`
- `ebitenui/renderer`: Ebiten 이미지 렌더링
- `ebitenui/prefabs`: 게임 UI 프리셋

## 예제

- 통합 쇼케이스: [examples/go/ebiten-ui-showcase/README.md](../../../examples/go/ebiten-ui-showcase/README.md)
- 재사용 가능한 UI 디버그 어댑터: [../ebiten-ui-debug/README.md](../ebiten-ui-debug/README.md)

## 현재 제약

- 레이아웃은 현재 `row`, `column`, `stack`, `scroll-view`, `anchored`, `grid` 중심의 v1 구현이다.
- `grid`는 track 기반 span/auto-placement와 alignment 분배를 지원한다.
- 텍스트 렌더링은 기본 비트맵 폰트를 사용한다.
- 런타임은 콜백 기반 1차 구현이며, React 스타일 훅 계층은 아직 없다.
- 상태 변경 후 화면 값 반영은 호출자가 DOM을 다시 빌드하는 흐름을 기본값으로 둔다.
