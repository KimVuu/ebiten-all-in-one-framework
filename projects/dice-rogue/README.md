# Dice Rogue

`projects/dice-rogue` is a first-act vertical slice for a dice-driven roguelike battle project. It assembles `ebiten-ui`, `ebiten-debug`, and `ebiten-ui-debug` without changing the shared library APIs.

## Included Flow

- party selection from 6 starter characters
- fixed act map with normal fights, rest nodes, elite, and boss
- shared team defense, random packet targeting, graveyard refill dice loop
- debug bridge, compact UI inspect/query/capture, and screenshots under `screenshots/dice-rogue`

## Run

```bash
cd projects/dice-rogue
go run ./cmd/dice-rogue
```

Debug bridge enabled:

```bash
cd projects/dice-rogue
EBITEN_DEBUG_MODE=1 EBITEN_DEBUG_ADDR=127.0.0.1:47833 go run ./cmd/dice-rogue
```

## Test

```bash
cd projects/dice-rogue
go test ./...
```

## Debug / MCP

1. Run the game with `EBITEN_DEBUG_MODE=1`.
2. Start the existing MCP server from the repository root if needed.
3. Use `/debug/ui/overview`, `/debug/ui/query`, `/debug/ui/node/{id}`, `run_command(ui_click)`, and `/debug/ui/capture`.
4. Captured PNG artifacts are stored under [screenshots/dice-rogue](/Users/kimyechan/Develop/Game/Ebiten/ebtien-aio-framework/screenshots/dice-rogue).
