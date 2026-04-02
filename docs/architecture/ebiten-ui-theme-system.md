# Ebiten UI Theme System Spec

## 목적

이 문서는 `ebiten-ui`의 `Theme system`을 구현 직전 수준으로 고정한다.  
목표는 현재 `components.go`, `prefabs.go`, showcase 내부에 흩어져 있는 색상과 spacing 결정을 token 기반 구조로 승격하는 것이다.

## 현재 문제

- reusable component 내부에 hard-coded color가 많다.
- `prefabs`도 자체 색상 세트를 별도로 가진다.
- showcase는 별도 데모 색상을 또 가진다.
- 결과적으로:
  - 시각 규칙 변경 비용이 크고
  - 일관성 검증이 어렵고
  - 프로젝트별 테마 교체가 사실상 불가능하다

## 설계 원칙

- `Theme`는 reusable visual decision의 canonical source다.
- component는 가능한 한 literal color 대신 theme token을 읽는다.
- theme 미지정 시 기본 `DefaultTheme()`를 사용한다.
- theme override는 `전역 -> prefab group -> component instance` 순으로 덮는다.
- debug surface는 theme 자체를 몰라도 되며, 최종 computed 결과만 다루면 된다.

## 범위

이번 spec이 다루는 것:

- 색상 token
- spacing token
- radius/stroke token
- component state tone
- theme 적용 경로

이번 spec이 다루지 않는 것:

- dynamic font loading 시스템
- rich text theme
- XML/template theme injection

## 타입 구조

### Root type

```go
type Theme struct {
  Name       string
  Palette    ThemePalette
  Spacing    ThemeSpacing
  Radius     ThemeRadius
  Stroke     ThemeStroke
  Typography ThemeTypography
  Components ThemeComponents
}
```

### Palette

```go
type ThemePalette struct {
  Background ThemeSurfacePalette
  Text       ThemeTextPalette
  Accent     ThemeAccentPalette
  Status     ThemeStatusPalette
  Overlay    color.Color
}
```

### Component themes

```go
type ThemeComponents struct {
  InputField   InputFieldTheme
  Textarea     TextareaTheme
  Dropdown     DropdownTheme
  Checkbox     CheckboxTheme
  Toggle       ToggleTheme
  Slider       SliderTheme
  Tabs         TabsTheme
  Modal        ModalTheme
  Tooltip      TooltipTheme
  ProgressBar  ProgressBarTheme
  Card         CardTheme
  Panel        PanelTheme
  Dialog       DialogTheme
  InventoryGrid InventoryGridTheme
}
```

필드명은 v1에서 필요한 컴포넌트만 포함하고, 나머지는 점진적으로 추가한다.

## 토큰 분류

### 색상 토큰

- `surface.base`
- `surface.elevated`
- `surface.sunken`
- `surface.panel`
- `surface.card`
- `text.strong`
- `text.muted`
- `text.inverse`
- `accent.primary`
- `accent.secondary`
- `accent.warning`
- `status.success`
- `status.warning`
- `status.error`
- `overlay.backdrop`

### spacing 토큰

- `xs`
- `sm`
- `md`
- `lg`
- `xl`

### radius 토큰

- `sm`
- `md`
- `lg`

### stroke 토큰

- `thin`
- `medium`
- `focus`

## 상태 모델

모든 interactive component는 최소 아래 상태를 가진다고 가정한다.

- `default`
- `hovered`
- `focused`
- `pressed`
- `selected`
- `disabled`

theme는 이 상태별 차이를 표현할 수 있어야 한다.

예:

```go
type StateColors struct {
  Background color.Color
  Border     color.Color
  Text       color.Color
}

type ButtonTheme struct {
  Default  StateColors
  Hovered  StateColors
  Focused  StateColors
  Pressed  StateColors
  Selected StateColors
  Disabled StateColors
}
```

## 적용 방식

### 1. Runtime theme source

v1 기본안:

- `Theme`는 `Runtime`가 아니라 build context 또는 helper 경로로 주입한다.
- simplest path:
  - `ThemeContext`
  - `WithTheme(theme Theme, root *Node) *Node`
  - 또는 component config에 `ThemeOverride *ComponentTheme`

### 2. Component read path

component는 아래 순서로 style을 계산한다.

1. instance override
2. component theme
3. palette/shared token
4. default fallback

### 3. Prefab read path

prefab은 자신만의 literal color를 들고 있지 않고:

1. `Theme.Components.Panel/Card/...`
2. 필요하면 `Theme.Components.Dialog/InventoryGrid/...`

를 읽는다.

## API 방향

v1 공개 표면 제안:

```go
func DefaultTheme() Theme
func NewTheme(name string) Theme
func ResolveTheme(theme *Theme) Theme
```

component config 확장 예:

```go
type InputFieldConfig struct {
  ...
  Theme *Theme
}
```

v1에서는 전역 mutable singleton을 두지 않는다.

## 마이그레이션 순서

1. `Theme` 타입과 `DefaultTheme()` 추가
2. `components.go` 내부 `componentTextStrong`, `componentPanel` 계열을 theme로 이동
3. `prefabs.go` 내부 `panelBG`, `accentBlue` 계열을 theme로 이동
4. showcase의 one-off 색상은 demo 전용 값만 남기고, reusable token은 theme로 이동
5. README와 showcase에 theme page 추가

## 테스트 기준

- `DefaultTheme()`가 비어 있지 않은 유효 token set을 반환하는지
- component가 theme override를 적용하는지
- prefab이 theme token을 읽는지
- 같은 component라도 state에 따라 색상이 달라지는지
- showcase에서 theme 교체 시 최소 대표 페이지의 visual capture가 달라지는지

## 현재 v1 반영 범위

- `Theme`, `DefaultTheme()`, `NewTheme()`, `ResolveTheme()` 구현
- 대표 입력계:
  - `InputField`
  - `Textarea`
  - `Checkbox`
  - `Toggle`
  - `Slider`
  - `ProgressBar`
- 대표 프리팹:
  - `Panel`
  - `Card`
  - `MenuList`
  - `Dialog`
  - `HUDBar`
  - `InventoryGrid`
- showcase에 `foundations/theme` 페이지 추가

## 완료 기준

- reusable component에서 color literal 비중이 크게 줄어든다.
- `prefabs`가 독자적인 색상 상수를 거의 갖지 않는다.
- showcase에 theme 설명 페이지가 생긴다.
- 실제 프로젝트가 `Theme` 하나로 기본 look을 바꿀 수 있다.
