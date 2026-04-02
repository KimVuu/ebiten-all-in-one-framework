# Ebiten UI Roadmap

## 목적

이 문서는 `ebiten-ui` 생태계의 다음 작업 순서를 고정한다.  
목표는 기능을 무작정 늘리는 것이 아니라, `컴포넌트`, `상태`, `테마`, `디버그`, `문서 경험`이 같이 성숙하도록 단계별로 정리하는 것이다.

## 현재 기준점

이미 갖춘 것:

- 선언형 노드 트리
- 제약 레이아웃과 validator
- 입력 런타임
- 페이지 라우터
- `ebiten-debug`, `ebiten-ui-debug`, `ebiten-mcp`
- compact UI inspect/query/capture
- 문서형 `ebiten-ui-showcase`

이제부터는 “컴포넌트를 더 만든다”보다 “라이브러리 경험을 한 단계 올린다”가 우선이다.

## Track 1. Theme System

### 목표

- reusable component가 theme token을 읽도록 바꾼다.
- showcase에서 theme variation을 문서화한다.

### 작업 항목

- `Theme` 타입 도입
- color, spacing, stroke, radius, typography token 정의
- component 기본 스타일을 theme 기반으로 이동
- dark/light 또는 game-specific theme 샘플 추가
- showcase에 theme page 추가

### 완료 기준

- reusable component에서 inline hard-coded color가 줄어든다.
- demo page에서 theme 교체 효과를 확인할 수 있다.

## Track 2. Reactive State and Binding

### 목표

- 상호작용 상태와 UI 값 상태를 더 명확히 분리한다.
- rerender를 덜 기계적으로 만들고, 선언형 사용성을 높인다.

### 작업 항목

- state binding model 설계
- controlled/uncontrolled component 기준 정리
- computed value helper
- page/router와 binding 계층의 경계 고정
- form state helper 또는 view-model 계층 도입 검토

### 완료 기준

- 입력 컴포넌트가 현재보다 덜 수동적으로 연결된다.
- showcase 예제 코드가 더 짧고 명확해진다.

## Track 3. Widget System Completion

### 목표

- taxonomy별 핵심 위젯 세트를 빠짐없이 정리한다.

### 작업 항목

- `layout`: `List`, `VirtualList`
- `overlay`: `ContextMenu`
- `data`: data browser / list navigation helpers
- `status`: `ProgressBar` 중심 usage 강화
- `prefabs`: settings, HUD, menu, inventory, shop, character sheet
- page별 usage note와 canonical code 보강

### 완료 기준

- showcase taxonomy 안에 “비어 있는 그룹”이 거의 없어야 한다.
- 실제 프로젝트가 prefab 조합으로 빠르게 UI를 만들 수 있어야 한다.

## Track 4. Navigation and App Composition

### 목표

- 페이지 라우터를 실제 앱 조립 수준으로 끌어올린다.

### 작업 항목

- nested page semantics 확장
- optional history/back-forward
- route-level metadata
- group page와 leaf page 조립 helper 강화
- project starter에서 page router 기본 wiring 제공

### 완료 기준

- settings, codex-like inspector, docs-style tool UI를 라우터로 자연스럽게 만들 수 있다.

## Track 5. Visual Testing and Tooling

### 목표

- UI 테스트를 구조 검증에서 시각 검증까지 확장한다.

### 작업 항목

- page별 screenshot baseline 규칙
- visual diff workflow
- capture naming 규칙
- issue triage flow
- script/tooling integration

### 완료 기준

- `overview -> inspect -> input -> capture -> diff` 흐름이 반복 가능한 루틴이 된다.

## Track 6. Project Adoption

### 목표

- 새 프로젝트가 showcase 코드를 복사하지 않고도 라이브러리 조합만으로 바로 UI debug surface를 쓸 수 있게 한다.

### 작업 항목

- 최소 project template
- `ebiten-ui + ebiten-debug + ebiten-ui-debug` 조립 예제
- project-level theme sample
- project-level command registration sample

### 완료 기준

- `projects/<name>` 하나를 새로 만들고 짧은 wiring만으로 MCP inspect/capture가 동작한다.

## 우선순위

### 바로 할 것

1. `Theme system`
2. `Reactive state/binding` 1차 설계
3. `showcase` 페이지별 canonical example 보강

### 그 다음

1. widget taxonomy completion
2. visual testing baseline
3. project starter/sample

### 나중

1. history/deep-link 같은 router 확장
2. template/XML 스타일 선언 계층 검토
3. richer editor/tool integration

## 운영 규칙

- 새 기능은 먼저 taxonomy와 아키텍처 문서에 위치를 정한다.
- reusable component는 showcase page 없이 추가하지 않는다.
- page 추가 시 최소 `설명`, `live demo`, `usage`, `code example`을 같이 넣는다.
- debug surface가 없는 UI 기능은 기본값으로 완료로 보지 않는다.

## 이번 문서의 의미

이 로드맵은 외부 라이브러리를 따라가기 위한 문서가 아니다.  
`ebiten-ui`를 `독립적인 Ebiten UI toolkit + debug harness + doc-driven showcase`로 키우기 위한 기준 문서다.
