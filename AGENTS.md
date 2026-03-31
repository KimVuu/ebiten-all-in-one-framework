# AGENTS

이 문서는 에이전트 작업의 진입점이다. 상세 정책은 `docs/agents/*`와 `docs/architecture/*`에서 관리한다.

## 목차

- [목적](#목적)
- [에이전트 역할](#에이전트-역할)
- [작업 분담 원칙](#작업-분담-원칙)
- [쓰기 범위](#쓰기-범위)
- [작업 프로세스](#작업-프로세스)
- [하네스 엔지니어링](#하네스-엔지니어링)
- [응답 원칙](#응답-원칙)
- [요청 템플릿](#요청-템플릿)
- [검증 규칙](#검증-규칙)
- [상세 문서 링크](#상세-문서-링크)

## 목적

에이전트가 저장소 구조를 이해하고, 충돌 없이 역할을 나누며, 같은 형식으로 결과를 보고하도록 돕는다.  
상세 문서: [docs/agents/roles.md](docs/agents/roles.md), [docs/agents/task-prompts.md](docs/agents/task-prompts.md)

## 에이전트 역할

`architect`, `bootstrapper`, `library-author`, `example-author`, `project-integrator`, `tool-author`, `docs-curator`를 기본 역할로 사용한다.  
상세 문서: [docs/agents/roles.md](docs/agents/roles.md)

## 작업 분담 원칙

역할별 쓰기 범위를 분리하고, 공통 규칙은 문서 계약을 우선한다. 구현 전에는 자신이 수정할 경로와 검증 범위를 먼저 고정한다.  
상세 문서: [docs/agents/roles.md](docs/agents/roles.md)

## 쓰기 범위

`libs`, `examples`, `projects`, `tools`, `scripts`, `docs`는 책임이 다르며, 에이전트는 맡은 영역 바깥의 구조를 임의로 재설계하지 않는다.  
상세 문서: [docs/architecture/monorepo.md](docs/architecture/monorepo.md), [docs/agents/roles.md](docs/agents/roles.md)

## 작업 프로세스

모든 구현 작업은 같은 순서를 따른다. 먼저 현재 테스트와 검증 명령을 확인하고, 필요한 테스트를 먼저 추가한 뒤 구현을 진행한다. 구현이 끝나면 테스트를 다시 돌리고, 마지막에 관련 문서를 갱신한다.  
상세 문서: [docs/agents/workflow.md](docs/agents/workflow.md), [docs/agents/task-prompts.md](docs/agents/task-prompts.md)

## 하네스 엔지니어링

이 저장소의 Ebiten/Ebitengine 앱은 하네스 중심 조립을 기본값으로 삼는다. 앱 진입점은 조립만 하고, 루프와 씬 생명주기, 공통 컨텍스트, 디버그 기능은 하네스가 소유하는 방향으로 설계한다.  
상세 문서: [docs/architecture/harness-engineering.md](docs/architecture/harness-engineering.md)

## 응답 원칙

컨텍스트를 아끼기 위해 응답은 짧고 직접적으로 유지한다. 현재 단계, 다음 행동, 막힌 점만 말하고, 불필요한 배경 설명과 반복 요약은 줄인다.  
상세 문서: [docs/agents/workflow.md](docs/agents/workflow.md), [docs/agents/task-prompts.md](docs/agents/task-prompts.md)

## 요청 템플릿

모노레포 부트스트랩, Go 라이브러리 추가, Bun 라이브러리 추가, 예제 추가, 프로젝트 추가, 도구 추가 요청은 표준 템플릿을 사용한다.  
상세 문서: [docs/agents/task-prompts.md](docs/agents/task-prompts.md)

## 검증 규칙

에이전트 결과는 항상 `변경 경로`, `실행 명령`, `검증 결과`, `남은 리스크` 형식으로 보고한다.  
상세 문서: [docs/agents/task-prompts.md](docs/agents/task-prompts.md)

## 상세 문서 링크

- 역할과 책임: [docs/agents/roles.md](docs/agents/roles.md)
- 작업 프로세스: [docs/agents/workflow.md](docs/agents/workflow.md)
- 작업 요청 템플릿: [docs/agents/task-prompts.md](docs/agents/task-prompts.md)
- 모노레포 구조: [docs/architecture/monorepo.md](docs/architecture/monorepo.md)
- 워크스페이스 계약: [docs/architecture/workspace-contract.md](docs/architecture/workspace-contract.md)
- 하네스 엔지니어링: [docs/architecture/harness-engineering.md](docs/architecture/harness-engineering.md)
- 명명 규칙: [docs/conventions/naming.md](docs/conventions/naming.md)
