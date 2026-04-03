package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kimyechan/ebiten-aio-framework/projects/dice-rogue/internal/app"
)

func main() {
	debugMode := debugModeEnabled()
	game := app.NewGame(app.GameConfig{
		DebugEnabled: debugMode,
		Seed:         7,
	})
	if debugMode {
		if err := game.StartDebugBridge(resolveDebugAddr()); err != nil {
			log.Fatal(err)
		}
		defer func() {
			_ = game.StopDebugBridge()
		}()
	}

	ebiten.SetWindowTitle("주사위 로그")
	ebiten.SetWindowSize(app.DefaultWindowWidth, app.DefaultWindowHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func debugModeEnabled() bool {
	value := os.Getenv("EBITEN_DEBUG_MODE")
	return value == "1" || value == "true" || value == "yes"
}

func resolveDebugAddr() string {
	value := os.Getenv("EBITEN_DEBUG_ADDR")
	if value == "" {
		return "127.0.0.1:47833"
	}
	return value
}
