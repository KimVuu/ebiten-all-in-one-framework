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
- `Theme`, `DefaultTheme`, `NewTheme`, `ResolveTheme` 기반 theme token 시스템
- `TextFace`, `SetTextFace` 기반 전역 text face 설정
- `Value[T]`, `WritableValue[T]`, `Ref[T]`, `Computed[T]` 기반 reactive binding primitive
- ID 기반 DOM 조회
- `LayoutNode`의 `ParentID`, `ContentBounds`, `ClipRect`, `ClickableRect`, `Overflow` 계산 필드
- `ValidateLayout` 기반의 레이아웃 검증과 constraint patch 제안
- `Runtime`, `InputSnapshot`, `EventHandlers` 기반 상호작용 런타임
- `OnPointerDown`, `OnPointerHold`, `OnPointerUp`, `OnClick` 기반 버튼 포인터 lifecycle
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

버튼과 interactive node는 포인터 lifecycle을 분리해서 다룰 수 있다.

- `OnPointerDown`: 누르기 시작할 때 1회
- `OnPointerHold`: 누르고 있는 동안 지속 프레임마다
- `OnPointerUp`: 손을 뗄 때 1회
- `OnClick`: 같은 타깃에서 눌렀다가 뗐을 때 1회

기본 런타임은 `Tab` 포커스 이동, `Escape` 포커스 해제, `Enter` 기반 focused button 활성화, arrow 기반 scroll dispatch까지 처리한다.

페이지 전체 스크롤이 필요한 화면은 `PageLayout`로 고정 헤더와 `ScrollView` 본문을 함께 구성하고, 문서형 UI나 툴형 화면은 `PageRouter + PageScreen`으로 좌측 navigation과 우측 detail panel을 조립하는 방식을 기본값으로 둔다.

`ScrollView`와 `ClipChildren` 제약을 가진 노드는 렌더 단계에서도 `ClipRect` 바깥 자식을 잘라서 그린다.

레이아웃 계산 결과를 디버그하거나 AI가 검사할 때는 `ValidateLayout(layout, viewport, opts)`를 사용하면 된다. 반환값은 `LayoutIssue`와 `ConstraintPatch` 중심으로 구성되어 있어 절대좌표가 아니라 제약 수정 단위로 다룰 수 있다.

## Theme 사용

대표 컴포넌트와 프리팹은 `Theme *Theme`를 통해 token override를 받을 수 있다.

```go
theme := ebitenui.NewTheme("forest")
theme.Components.Panel.Background = color.RGBA{R: 20, G: 39, B: 34, A: 255}
theme.Components.InputField.Border = color.RGBA{R: 88, G: 187, B: 152, A: 255}

node := ebitenui.InputField(ebitenui.InputFieldConfig{
	ID:    "search",
	Label: "Search",
	Theme: &theme,
})
```

v1에서는 대표 입력계와 대표 프리팹부터 theme token을 읽고, 나머지 컴포넌트는 같은 구조로 점진 확장한다.

## Font Face 사용

텍스트 측정과 렌더링은 같은 global text face를 공유한다.

```go
face := loadYourFontFace()
ebitenui.SetTextFace(face)
```

설정하지 않으면 기본값은 `basicfont.Face7x13`이다.

## Reactive Binding 사용

대표 입력 컴포넌트는 controlled/uncontrolled 둘 다 지원한다.

```go
name := ebitenui.NewRef("Kim")
enabled := ebitenui.NewRef(true)

field := ebitenui.InputField(ebitenui.InputFieldConfig{
	ID:           "name-input",
	Label:        "Player Name",
	ValueBinding: name,
})

toggle := ebitenui.Toggle(ebitenui.ToggleConfig{
	ID:             "music-toggle",
	Label:          "Music",
	CheckedBinding: enabled,
})
```

현재 v1 binding 적용 범위:

- `InputField`
- `Textarea`
- `Checkbox`
- `Toggle`
- `Slider`
- `Dropdown`

`Computed`는 derived value를 읽기 전용으로 노출하는 작은 helper이고, hook runtime이나 global store는 아직 범위 밖이다.

## 패키지 구성

- `ebitenui`: 노드, 스타일, 레이아웃, DOM 조회
- `ebitenui` router surface: `PageRoute`, `PageRouter`, `PageScreen`
- `ebitenui/renderer`: Ebiten 이미지 렌더링
- `ebitenui/prefabs`: 게임 UI 프리셋

## 예제

- 통합 쇼케이스: [examples/go/ebiten-ui-showcase/README.md](../../../examples/go/ebiten-ui-showcase/README.md)
- 재사용 가능한 UI 디버그 어댑터: [../ebiten-ui-debug/README.md](../ebiten-ui-debug/README.md)
- 아키텍처 기준: [../../../docs/architecture/ebiten-ui-architecture.md](../../../docs/architecture/ebiten-ui-architecture.md)
- 로드맵: [../../../docs/architecture/ebiten-ui-roadmap.md](../../../docs/architecture/ebiten-ui-roadmap.md)
- Theme spec: [../../../docs/architecture/ebiten-ui-theme-system.md](../../../docs/architecture/ebiten-ui-theme-system.md)
- Reactive binding spec: [../../../docs/architecture/ebiten-ui-reactive-binding.md](../../../docs/architecture/ebiten-ui-reactive-binding.md)
- Project integration spec: [../../../docs/architecture/ebiten-ui-project-integration.md](../../../docs/architecture/ebiten-ui-project-integration.md)
- 스타일 규칙: [../../../docs/conventions/ebiten-ui-style.md](../../../docs/conventions/ebiten-ui-style.md)
- 리뷰 체크리스트: [../../../docs/conventions/ebiten-ui-review-checklist.md](../../../docs/conventions/ebiten-ui-review-checklist.md)

## 현재 제약

- 레이아웃은 현재 `row`, `column`, `stack`, `scroll-view`, `anchored`, `grid` 중심의 v1 구현이다.
- `grid`는 track 기반 span/auto-placement와 alignment 분배를 지원한다.
- 텍스트 렌더링은 기본 비트맵 폰트를 사용한다.
- 런타임은 콜백 기반 1차 구현이며, React 스타일 훅 계층은 아직 없다.
- 상태 변경 후 화면 값 반영은 호출자가 DOM을 다시 빌드하는 흐름을 기본값으로 둔다.
