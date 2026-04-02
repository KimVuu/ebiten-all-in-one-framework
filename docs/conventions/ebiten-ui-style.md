# Ebiten UI Style

## 목적

이 문서는 `ebiten-ui` 계열 코드와 문서를 어떤 스타일로 쌓을지 고정한다.  
핵심은 `문서형 showcase`, `명확한 taxonomy`, `선언형 builder`, `재사용 가능한 debug surface`를 일관되게 유지하는 것이다.

## 기본 스타일 원칙

- `ebiten-ui`는 직접 그리기보다 선언형 tree 조립을 우선한다.
- low-level node와 high-level prefab은 같은 파일에서 섞지 않는다.
- component는 `Config struct -> *Node` 패턴을 기본값으로 둔다.
- 예제는 “스크래치 패드”가 아니라 문서형 reference page여야 한다.
- debug 코드와 product 코드는 같은 core를 공유하고, debug surface만 선택적으로 붙인다.

## 컴포넌트 작성 규칙

- stable `ID`를 기준으로 runtime state를 연결한다.
- component는 가능한 한 외부에서 상태를 주입받는다.
- 상태가 필요한 경우:
  - immediate interaction state는 `Runtime`
  - navigation state는 `PageRouter`
  - 나머지 UI 값 상태는 점진적으로 binding 계층으로 분리한다.
- reusable component는 특정 showcase 페이지 구조를 알면 안 된다.

## 시각 스타일 규칙

- reusable component는 inline color를 줄이고 token/theme 중심으로 이동한다.
- spacing, stroke, background, text tone은 반복될수록 helper 또는 theme에 올린다.
- one-off demo styling은 showcase에 두고, reusable styling contract는 라이브러리에 둔다.

## showcase 작성 규칙

각 page는 아래 4개를 기본으로 가져야 한다.

1. 짧은 설명
2. 실제 interactive demo
3. 사용법 설명
4. canonical code example

group page는 개요 역할을 하고, leaf page는 단일 컴포넌트나 prefab에 집중한다.

## taxonomy 규칙

새 component나 prefab은 먼저 이 그룹 중 하나에 속해야 한다.

- `foundations`
- `tags`
- `inputs`
- `layout`
- `overlay`
- `data`
- `status`
- `prefabs`

taxonomy 밖의 새 그룹을 만들려면, 기존 그룹으로는 왜 표현이 안 되는지 먼저 설명해야 한다.

## debug 연동 규칙

- `ebiten-ui` 코어는 debug bridge를 직접 알지 않는다.
- `ebiten-ui-debug`가 bridge attach, overlay, inspect/query/capture, input injection을 담당한다.
- showcase나 실제 프로젝트는 `ebiten-ui-debug`를 wiring만 해야 한다.

## code example 규칙

- showcase code example은 런타임에 실제 파일을 읽어오지 않는다.
- 사람이 큐레이션한 canonical example string을 사용한다.
- 이유:
  - 문서 안정성
  - 길이 제어
  - 설명 목적 최적화

## 문서 우선순위

- README: 시작점
- architecture docs: 경계와 책임
- roadmap: 다음 순서
- showcase: 살아 있는 사용 문서

같은 내용을 여러 군데 장문으로 중복하지 않는다.

## 하지 않을 것

- 외부 UI 라이브러리 API를 그대로 복제하지 않는다.
- 디버그 기능 때문에 코어 UI 라이브러리의 책임을 흐리지 않는다.
- example 전용 glue를 reusable library 안으로 무비판적으로 넣지 않는다.
- product UI와 debug UI가 서로 다른 core를 보게 만들지 않는다.
