# Debug Bridge and MCP Adapter

## 목적

이 문서는 Ebiten 디버그 브리지와 MCP 어댑터의 책임 경계를 고정한다.  
핵심은 게임 내부 디버그 표면과 외부 자동화 표면을 분리해서, 게임은 하네스 계약만 알고 MCP는 그 위에 얇게 올라가게 만드는 것이다.

## 구조

- `libs/go/ebitendebug`
  - 게임 프로세스 안에서 실행된다.
  - 프레임, 씬, 월드, UI tree, 디버그 커맨드를 HTTP JSON으로 노출한다.
  - provider 등록과 수명주기만 책임진다.
- `libs/go/ebiten-mcp`
  - 별도 stdio runner나 호스트 프로세스에서 소비된다.
  - 실행 중인 디버그 브리지 주소에 attach 한다.
  - HTTP endpoint를 MCP tool로 변환한다.
- `tools/ebiten-mcp-server`
  - `libs/go/ebiten-mcp`를 사용해 실제 stdio MCP 서버 프로세스를 실행한다.
  - 배포와 로컬 실행 엔트리를 담당한다.
- `examples/go/debug-bridge`
  - 라이브러리 통합과 명령 흐름을 검증한다.

## 하네스 경계

- 게임 진입점은 브리지 활성화 여부와 주소만 조립한다.
- 게임 로직은 자신의 씬/엔티티/UI 상태를 스냅샷으로 변환해 provider로 등록한다.
- 브리지는 게임 내부 구조를 직접 탐색하지 않는다.
- MCP 어댑터는 게임 엔진 타입에 의존하지 않고 HTTP 계약만 소비한다.

## v1 운영 규칙

- loopback 전용 바인딩
- 인증 없음
- 단일 게임 인스턴스 attach
- 읽기 endpoint는 스냅샷만 반환
- side effect는 등록된 command 실행으로만 허용
