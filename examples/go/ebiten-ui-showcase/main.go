package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	debugMode := debugModeEnabled()
	game := newGame(debugMode)
	if debugMode {
		if err := game.startDebugBridge(resolveShowcaseDebugAddr()); err != nil {
			log.Fatal(err)
		}
		defer func() {
			_ = game.stopDebugBridge()
		}()
	}

	ebiten.SetWindowTitle("ebiten ebiten-ui showcase")
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func getenv(name string) string {
	return os.Getenv(name)
}
