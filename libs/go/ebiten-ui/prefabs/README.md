# ebiten Ebiten UI Prefabs

`ebitenui/prefabs`는 게임에서 자주 쓰는 일회성 UI를 빠르게 조립하기 위한 상위 팩토리 집합이다. 각 프리셋은 설정 struct를 받아 `*ebitenui.Node`를 반환한다.

## 현재 프리셋

- `Panel`
- `Card`
- `StatusRow`
- `MenuList`
- `Dialog`
- `HUDBar`
- `InventoryGrid`
- `PauseMenu`
- `SettingsPanel`
- `Tooltip`

## 원칙

- 프리셋은 그리기 위젯이 아니라 DOM 팩토리다.
- 입력 상태는 config나 props로 주입한다.
- 공통 레이아웃과 스타일은 `ebiten-ui`의 저수준 노드 위에 조립한다.

## 예제

- 통합 쇼케이스: [examples/go/ebiten-ui-showcase/README.md](../../../../examples/go/ebiten-ui-showcase/README.md)
