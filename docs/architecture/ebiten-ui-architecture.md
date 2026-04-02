# Ebiten UI Architecture

## 목적

이 문서는 `ebiten-ui` 생태계의 목표 구조를 고정한다.  
핵심은 외부 UI 라이브러리를 직접 의존하지 않고도, `문서형 showcase`, `컴포넌트 분류`, `반응형 상태`, `테마`, `디버그 도구`가 함께 움직이는 독립적인 Ebiten UI 스택을 만드는 것이다.

## 외부 참고에 대한 기준

- 이 저장소는 `Willow UI`를 사용하지 않는다.
- 다만 아래 방향은 참고한다.
  - 문서와 갤러리가 같은 제품 경험을 만든다는 점
  - 위젯을 역할별 taxonomy로 정리하는 점
  - 런타임 상태와 컴포넌트 계층을 분리하는 점
  - 개발 도구와 시각 검증을 UI 라이브러리의 일부 경험으로 보는 점
- 참고하는 것은 구조와 운영 방식이지, 코드나 API를 그대로 복제하는 것이 아니다.

## 목표 계층

### 1. Core UI

- 소유 경로: `libs/go/ebiten-ui`
- 책임:
  - DOM 유사 노드 트리
  - 레이아웃과 제약 계산
  - 입력 런타임
  - 페이지 라우터
  - 기본 컴포넌트와 prefab 조립
- 비책임:
  - HTTP 디버그 브리지
  - MCP tool 표면
  - screenshot artifact 저장
  - debug overlay

### 2. UI Debug Adapter

- 소유 경로: `libs/go/ebiten-ui-debug`
- 책임:
  - `ebiten-ui` 레이아웃을 디버그 snapshot으로 변환
  - compact inspect/query/issues/capture
  - overlay
  - debug input queue
  - 기본 UI command surface
- 비책임:
  - 게임 전체 frame/scene/world contract
  - 제품 전용 UI 로직

### 3. Generic Debug Bridge

- 소유 경로: `libs/go/ebiten-debug`
- 책임:
  - frame, scene, world, ui provider를 HTTP로 노출
  - command registration
  - generic transport surface
- 비책임:
  - 특정 UI 라이브러리 해석
  - 앱별 debug glue

### 4. MCP Adapter

- 소유 경로: `libs/go/ebiten-mcp`, `tools/ebiten-mcp-server`
- 책임:
  - 디버그 브리지를 MCP tool surface로 변환
  - compact UI testing loop 지원
- 비책임:
  - UI 렌더링
  - 앱 내부 상태 관리

### 5. Doc-Driven Showcase

- 소유 경로: `examples/go/ebiten-ui-showcase`
- 책임:
  - 각 컴포넌트와 prefab을 페이지 단위로 설명
  - 데모, 사용법, 코드 예시를 한 화면에서 제공
  - 디버그/MCP 검증의 기준 앱 역할 수행
- 비책임:
  - 재사용 가능한 라이브러리 로직 보관

## 컴포넌트 taxonomy

`ebiten-ui`는 컴포넌트를 기능이 아니라 역할 기준으로 분류한다.

- `foundations`
  - `Image`, `TextBlock`, `Spacer`, `Stack`, `ScrollView`
- `tags`
  - `div`, `header`, `main`, `section`, `footer`, `button`, `span`, `text`
- `inputs`
  - `InputField`, `Textarea`, `Dropdown`, `RadioGroup`, `Stepper`
- `layout`
  - `Grid`, 향후 `List`, `VirtualList`, page composition helper
- `overlay`
  - `Modal`, `Tooltip`, 향후 `ContextMenu`
- `data`
  - `Tabs`, 향후 data-oriented navigation components
- `status`
  - `Toggle`, `Slider`, `ProgressBar`, `Badge`, `Chip`
- `prefabs`
  - `Dialog`, `InventoryGrid`, 향후 게임별 reusable assemblies

새 컴포넌트는 먼저 이 taxonomy 안에 들어가야 한다. 어느 그룹에도 자연스럽게 들어가지 않으면, 컴포넌트 추가 전에 taxonomy를 먼저 재검토한다.

## 상태 모델

### 현재 기준

- `Runtime`은 focus, hover, pointer, text, selection, scroll 같은 즉시 상호작용 상태를 소유한다.
- `PageRouter`는 현재 페이지와 계층 상태를 소유한다.
- 호출자는 DOM을 다시 빌드해서 상태를 반영한다.

### 다음 목표

- `Runtime`은 상호작용 상태만 담당한다.
- UI 값 상태는 점진적으로 `reactive binding` 계층으로 분리한다.
- `PageRouter`는 페이지 네비게이션의 canonical source로 유지한다.
- 컴포넌트는 가능한 한 `Config struct + stable ID + external state` 패턴을 유지한다.

## 테마 모델

현재는 showcase와 component 내부에 색상과 spacing이 분산되어 있다.  
다음 단계의 기본 방향은 아래와 같다.

- `Theme`는 색상, spacing, radius, stroke, text tone, state tone을 가진다.
- reusable component는 inline hard-coded color를 줄이고 theme token을 우선 사용한다.
- showcase는 theme sample을 문서 페이지로 보여줄 수 있어야 한다.
- debug surface는 theme 값 자체를 몰라도 되며, 최종 computed 결과만 보면 된다.

## 코드 스타일 기준

- 컴포넌트는 `Config struct -> *Node` 패턴을 기본값으로 둔다.
- 런타임에 필요한 상태는 stable `ID`를 기준으로 연결한다.
- low-level node helper와 high-level prefab을 섞지 않는다.
- 디버그 전용 glue는 코어 라이브러리 안으로 끌어들이지 않는다.
- 예제 페이지에 보여주는 코드는 “설명용 canonical example”을 별도로 유지한다.
- `showcase`는 문서형 reference app이고, 제품 코드를 대신하지 않는다.

## 문서 전략

문서 표면은 아래 3단계로 나눈다.

1. 루트/라이브러리 README
- 무엇을 하는지
- 어디서 시작하는지

2. 아키텍처 문서
- 왜 이렇게 나뉘는지
- 계층과 책임이 무엇인지

3. showcase 페이지
- 실제 데모
- 사용법
- 코드 예시

즉, README는 허브고, 아키텍처 문서는 기준이며, showcase는 살아 있는 reference다.

## 디버그와 검증 전략

UI 개발의 기본 검증 경로는 아래 순서를 따른다.

1. `get_ui_overview`
2. `query_ui_nodes`
3. `inspect_ui_node`
4. `run_command`로 `ui_click`, `ui_scroll`, `ui_type_text`
5. `capture_ui_screenshot`

full tree dump는 legacy/debug 용도로만 쓰고, 기본 설계 테스트 경로는 compact surface를 사용한다.

## 앞으로의 설계 방향

- `ebiten-ui`는 `core widget toolkit + page router + theme-ready component system`으로 간다.
- `ebiten-ui-debug`는 `UI testing harness`로 간다.
- `ebiten-ui-showcase`는 `문서형 gallery + 검증용 app`으로 간다.
- 실제 제품 프로젝트는 showcase 코드를 복사하지 않고, 라이브러리 조합만으로 같은 debug surface를 붙일 수 있어야 한다.
