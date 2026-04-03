# 주사위 로그

`projects/dice-rogue`는 주사위 기반 로그라이크 전투 프로젝트의 1막 수직 슬라이스입니다. 공유 라이브러리 API는 바꾸지 않고 `ebiten-ui`, `ebiten-debug`, `ebiten-ui-debug`를 조립해 사용합니다.

## 포함된 흐름

- 6명 스타터 캐릭터 중 3명 파티 선택
- 일반전, 휴식, 엘리트, 보스로 이어지는 고정 1막 지도
- 파티 공유 방어, 무작위 피해 타깃, 묘지 재충전 주사위 순환
- UI 텍스트는 `ebitenui.SetTextFace(...)`로 `asset/fonts/NeoDunggeunmoPro-Regular.ttf`를 적용
- 디버그 브리지, UI inspect/query/capture, `screenshots/dice-rogue` 스크린샷 저장

## 실행

```bash
cd projects/dice-rogue
go run ./cmd/dice-rogue
```

디버그 브리지 활성화:

```bash
cd projects/dice-rogue
EBITEN_DEBUG_MODE=1 EBITEN_DEBUG_ADDR=127.0.0.1:47833 go run ./cmd/dice-rogue
```

## 테스트

```bash
cd projects/dice-rogue
go test ./...
```

## 디버그 / MCP

1. `EBITEN_DEBUG_MODE=1`로 게임을 실행합니다.
2. 필요하면 저장소 루트에서 기존 MCP 서버를 함께 실행합니다.
3. `/debug/ui/overview`, `/debug/ui/query`, `/debug/ui/node/{id}`, `run_command(ui_click)`, `/debug/ui/capture`를 사용합니다.
4. 캡처된 PNG는 [screenshots/dice-rogue](/Users/kimyechan/Develop/Game/Ebiten/ebtien-aio-framework/screenshots/dice-rogue)에 저장됩니다.
