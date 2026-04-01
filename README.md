# ebiten AIO Framework

`ebiten`을 활용하기 쉬운 올인원 모노레포를 목표로 하는 저장소다. 이 저장소는 `Ebiten/Ebitengine` 기반 개발을 중심에 두되, Go와 TypeScript(Bun)를 함께 운영할 수 있는 구조와 재사용 라이브러리를 함께 쌓아가는 저장소다.

## 현재 상태

- 설계 문서와 작업 규칙을 먼저 고정했다.
- 현재 Go 라이브러리는 `libs/go/ebiten-ui`, `libs/go/ebiten-ui-debug`, `libs/go/ebiten-debug`, `libs/go/ebiten-mcp`까지 정리되어 있다.
- `ebiten-ui`는 선언형 UI, 입력 런타임, 제약 레이아웃, 프리셋을 담당한다.
- `ebiten-ui-debug`는 UI snapshot, compact inspect/query, overlay, capture, debug input queue를 재사용 가능한 어댑터로 제공한다.
- 모노레포의 최상위 인터페이스는 `libs`, `examples`, `projects`, `tools`, `scripts`, `docs` 여섯 가지로 제한한다.
- `AGENTS.md`는 에이전트 작업 허브이고, 상세 규칙은 `docs/agents/*`에서 관리한다.
- `tools/ebiten-mcp-server`와 `scripts/*`는 현재 라이브러리와 예제를 실행하고 검증하는 실제 엔트리로 사용한다.

## 디렉터리 기준

- `libs/`: 재사용 가능한 라이브러리
- `examples/`: 테스트, 실험, 검증용 샘플 프로젝트
- `projects/`: 실제 개발 중인 프로젝트
- `tools/`: 독립 실행형 개발 도구와 로컬 서비스
- `scripts/`: 반복 실행용 얇은 래퍼 스크립트
- `docs/`: 아키텍처, 규칙, 런북, 에이전트 문서

## 기본 운영 방향

- TypeScript 계열은 `Bun workspaces` 기준으로 운영한다.
- Go 계열은 `go.work` 기준으로 묶는 것을 기본안으로 삼는다.
- 루트 집계 도구는 향후 `package.json`, `bunfig.toml`, `go.work`, `Makefile` 조합을 기본안으로 사용한다.

## 문서 시작점

- 아키텍처: [docs/architecture/monorepo.md](docs/architecture/monorepo.md)
- 하네스 12원칙 평가 기준: [docs/architecture/harness-principles.md](docs/architecture/harness-principles.md)
- 디버그 브리지와 MCP: [docs/architecture/debug-bridge-mcp.md](docs/architecture/debug-bridge-mcp.md)
- 워크스페이스 계약: [docs/architecture/workspace-contract.md](docs/architecture/workspace-contract.md)
- 명명 규칙: [docs/conventions/naming.md](docs/conventions/naming.md)
- 에이전트 허브: [AGENTS.md](AGENTS.md)

## 현재 라이브러리

- Go Ebiten UI: [libs/go/ebiten-ui/README.md](libs/go/ebiten-ui/README.md)
- Go Ebiten UI Prefabs: [libs/go/ebiten-ui/prefabs/README.md](libs/go/ebiten-ui/prefabs/README.md)
- Go Ebiten UI Debug Adapter: [libs/go/ebiten-ui-debug/README.md](libs/go/ebiten-ui-debug/README.md)
- Ebiten Debug Bridge: [libs/go/ebiten-debug/README.md](libs/go/ebiten-debug/README.md)
- Ebiten MCP Adapter: [libs/go/ebiten-mcp/README.md](libs/go/ebiten-mcp/README.md)

## 현재 도구

- Ebiten MCP Server: [tools/ebiten-mcp-server/README.md](tools/ebiten-mcp-server/README.md)

## 현재 예제

- Go Ebiten UI Showcase: [examples/go/ebiten-ui-showcase/README.md](examples/go/ebiten-ui-showcase/README.md)
- Debug Bridge Example: [examples/go/ebiten-debug-bridge/README.md](examples/go/ebiten-debug-bridge/README.md)

## 전제 환경

- Go `1.25.x`
- Bun `1.2.x`

## 비고

루트 워크스페이스 집계 파일은 아직 만들지 않았다. 현재는 문서 체계, `ebiten-ui`, `ebiten-ui-debug`, `ebiten-debug`, `ebiten-mcp` 라이브러리, `ebiten-mcp-server` 도구, 그리고 이를 검증하는 예제를 기준점으로 잡는 단계다.
