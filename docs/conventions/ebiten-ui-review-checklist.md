# Ebiten UI Review Checklist

## 목적

이 문서는 `ebiten-ui` 계열 코드 변경을 리뷰할 때 보는 체크리스트다.  
목표는 코드 스타일과 구조 기준을 구현 전후에 같은 눈으로 점검하는 것이다.

## 현재 코드 기준 체크

### 이미 좋은 점

- `Config struct -> *Node` 패턴이 넓게 적용돼 있다.
- `prefabs`와 low-level component가 패키지 수준에서 분리돼 있다.
- `ebiten-ui-debug`가 코어 UI와 분리돼 있다.
- showcase가 문서형 reference app으로 이동했다.

### 현재 보완이 필요한 점

- reusable color/stroke 값이 아직 많이 literal이다.
- binding 계층이 없어서 form state 연결이 수동적이다.
- theme token이 없어 component/prefab/showcase 간 시각 일관성 점검이 어렵다.
- 일부 code example은 canonical reference로 더 다듬을 여지가 있다.

## 리뷰 질문

### 구조

- 이 변경이 `ebiten-ui`, `ebiten-ui-debug`, `ebiten-debug`, `showcase` 중 어디 책임인가
- reusable logic가 example/project 안에 갇히지 않았는가
- debug 전용 glue가 코어 라이브러리에 섞이지 않았는가

### 컴포넌트

- `Config struct -> *Node` 패턴을 따르는가
- stable `ID`가 필요한 곳에 빠지지 않았는가
- 컴포넌트가 taxonomy 안에 자연스럽게 속하는가
- component와 prefab이 같은 책임을 중복하지 않는가

### 상태

- interaction state와 value state가 구분되는가
- `Runtime`가 값 상태를 과도하게 소유하지 않는가
- controlled/uncontrolled 기준이 모호하지 않은가

### 시각 스타일

- literal color 대신 token/theme로 올릴 수 있는 값이 아닌가
- spacing/padding/gap이 반복되면 helper나 theme로 올려야 하는 것 아닌가
- demo 전용 스타일과 reusable 스타일이 섞이지 않았는가

### 문서

- 새 reusable component면 showcase page가 있는가
- showcase page에 설명, demo, usage, code example이 모두 있는가
- README/architecture/roadmap 중 어디에 링크가 필요한지 반영됐는가

### 디버그

- compact inspect/query/capture surface에서 검증 가능한가
- debug overlay 없이도 구조 검증이 가능한가
- screenshot capture나 MCP input 흐름을 깨지 않았는가

## 리뷰 상태 기준

- `pass`
  - 구조와 스타일 기준이 명확히 지켜짐
- `warn`
  - 동작은 가능하지만 theme/binding/doc/debug 기준 중 일부가 약함
- `fail`
  - 책임 경계 붕괴, 문서 누락, 디버그 검증 불가, 재사용성 저하

## 이번 기준 문서와의 연결

- 구조 기준: `docs/architecture/ebiten-ui-architecture.md`
- 로드맵 기준: `docs/architecture/ebiten-ui-roadmap.md`
- 스타일 기준: `docs/conventions/ebiten-ui-style.md`
- 구현 직전 설계서:
  - `docs/architecture/ebiten-ui-theme-system.md`
  - `docs/architecture/ebiten-ui-reactive-binding.md`
  - `docs/architecture/ebiten-ui-project-integration.md`
