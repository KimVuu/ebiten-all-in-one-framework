package main

import (
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	debugMode := envEnabled("EBITEN_DEBUG_MODE")
	game := newGame(debugMode)
	if debugMode {
		if err := game.startDebugBridge(resolveDebugAddr()); err != nil {
			log.Fatal(err)
		}
		defer func() {
			_ = game.stopDebugBridge()
		}()
	}

	ebiten.SetWindowTitle("ebiten debug bridge example")
	ebiten.SetWindowSize(960, 540)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func envEnabled(name string) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(name)))
	switch value {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func resolveDebugAddr() string {
	if value := os.Getenv("EBITEN_DEBUG_ADDR"); value != "" {
		return value
	}
	return ebitendebugDefaultAddr
}

const ebitendebugDefaultAddr = "127.0.0.1:47831"
