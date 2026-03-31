package main

import (
	"context"
	"image/color"
	"strconv"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kimyechan/ebiten-aio-framework/libs/go/ebitendebug"
)

type scene struct {
	ID   string
	Name string
}

type entity struct {
	ID      string
	Type    string
	Tags    []string
	X       float64
	Y       float64
	Width   float64
	Height  float64
	Visible bool
	Enabled bool
	Props   map[string]any
}

type game struct {
	mu sync.RWMutex

	width        int
	height       int
	frame        int
	tick         int
	paused       bool
	debugEnabled bool

	scenes       []scene
	currentScene string
	entities     []entity

	debugBridge *ebitendebug.Bridge
}

func newGame(debugMode bool) *game {
	return &game{
		width:        960,
		height:       540,
		debugEnabled: true,
		scenes: []scene{
			{ID: "menu", Name: "Menu"},
			{ID: "battle", Name: "Battle"},
		},
		currentScene: "menu",
		entities: []entity{
			{
				ID:      "player-hero",
				Type:    "player",
				Tags:    []string{"party", "selected"},
				X:       96,
				Y:       140,
				Width:   24,
				Height:  24,
				Visible: true,
				Enabled: true,
				Props: map[string]any{
					"hp": 100,
				},
			},
			{
				ID:      "npc-guide",
				Type:    "npc",
				Tags:    []string{"guide"},
				X:       180,
				Y:       160,
				Width:   20,
				Height:  28,
				Visible: true,
				Enabled: true,
				Props: map[string]any{
					"dialogue": "welcome",
				},
			},
		},
	}
}

func (game *game) Update() error {
	game.mu.Lock()
	defer game.mu.Unlock()

	game.tick++
	if !game.paused {
		game.frame++
	}

	return nil
}

func (game *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 17, G: 23, B: 34, A: 255})

	game.mu.RLock()
	defer game.mu.RUnlock()

	lines := []string{
		"Debug Bridge Example",
		"Scene: " + game.currentScene,
		"Paused: " + strconv.FormatBool(game.paused),
		"Entities: " + strconv.Itoa(len(game.entities)),
	}
	if game.debugBridge != nil {
		lines = append(lines, "Bridge: http://"+game.debugBridge.Address())
	}

	ebitenutil.DebugPrint(screen, strings.Join(lines, "\n"))
}

func (game *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return game.width, game.height
}

func (game *game) startDebugBridge(addr string) error {
	bridge := ebitendebug.New(ebitendebug.Config{
		Enabled: true,
		Addr:    addr,
		GameID:  "debug-bridge-example",
		Version: "v1",
	})
	game.registerDebugProviders(bridge)
	if err := bridge.Start(); err != nil {
		return err
	}
	game.debugBridge = bridge
	return nil
}

func (game *game) stopDebugBridge() error {
	if game.debugBridge == nil {
		return nil
	}
	return game.debugBridge.Close(context.Background())
}

func (game *game) registerDebugProviders(bridge *ebitendebug.Bridge) {
	bridge.SetFrameProvider(game.frameSnapshot)
	bridge.SetSceneProvider(game.sceneSnapshot)
	bridge.SetWorldProvider(game.worldSnapshot)
	bridge.SetUIProvider(game.uiSnapshot)
	bridge.RegisterCommand("pause.toggle", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		game.mu.Lock()
		defer game.mu.Unlock()

		game.paused = !game.paused
		return ebitendebug.CommandResult{
			Success: true,
			Message: "pause toggled",
			Payload: map[string]any{
				"paused": game.paused,
			},
		}
	})
	bridge.RegisterCommand("scene.switch", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		target, _ := request.Args["scene"].(string)
		if target == "" {
			return ebitendebug.CommandResult{
				Success: false,
				Message: "scene arg is required",
			}
		}

		game.mu.Lock()
		defer game.mu.Unlock()
		if !game.hasScene(target) {
			return ebitendebug.CommandResult{
				Success: false,
				Message: "scene not found",
			}
		}

		game.currentScene = target
		return ebitendebug.CommandResult{
			Success: true,
			Message: "scene switched",
			Payload: map[string]any{
				"scene": target,
			},
		}
	})
	bridge.RegisterCommand("entity.visibility.toggle", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		target, _ := request.Args["entity"].(string)
		if target == "" {
			return ebitendebug.CommandResult{
				Success: false,
				Message: "entity arg is required",
			}
		}

		game.mu.Lock()
		defer game.mu.Unlock()
		entity := game.entityByIDLocked(target)
		if entity == nil {
			return ebitendebug.CommandResult{
				Success: false,
				Message: "entity not found",
			}
		}

		entity.Visible = !entity.Visible
		return ebitendebug.CommandResult{
			Success: true,
			Message: "entity visibility toggled",
			Payload: map[string]any{
				"entity":  entity.ID,
				"visible": entity.Visible,
			},
		}
	})
}

func (game *game) frameSnapshot() ebitendebug.FrameSnapshot {
	game.mu.RLock()
	defer game.mu.RUnlock()

	return ebitendebug.FrameSnapshot{
		Frame:        game.frame,
		Tick:         game.tick,
		FPS:          ebiten.ActualFPS(),
		TPS:          ebiten.ActualTPS(),
		Paused:       game.paused,
		DebugEnabled: game.debugEnabled,
	}
}

func (game *game) sceneSnapshot() ebitendebug.SceneSnapshot {
	game.mu.RLock()
	defer game.mu.RUnlock()

	known := make([]ebitendebug.SceneRef, 0, len(game.scenes))
	current := ebitendebug.SceneRef{}
	for _, scene := range game.scenes {
		ref := ebitendebug.SceneRef{ID: scene.ID, Name: scene.Name}
		known = append(known, ref)
		if scene.ID == game.currentScene {
			current = ref
		}
	}

	return ebitendebug.SceneSnapshot{
		Current: current,
		Known:   known,
		Summary: map[string]any{
			"entityCount": len(game.entities),
			"paused":      game.paused,
		},
	}
}

func (game *game) worldSnapshot() ebitendebug.WorldSnapshot {
	game.mu.RLock()
	defer game.mu.RUnlock()

	entities := make([]ebitendebug.EntitySnapshot, 0, len(game.entities))
	for _, entity := range game.entities {
		props := map[string]any{}
		for key, value := range entity.Props {
			props[key] = value
		}

		entities = append(entities, ebitendebug.EntitySnapshot{
			ID:      entity.ID,
			Type:    entity.Type,
			Tags:    append([]string(nil), entity.Tags...),
			Visible: entity.Visible,
			Enabled: entity.Enabled,
			Position: ebitendebug.Vector2{
				X: entity.X,
				Y: entity.Y,
			},
			Size: ebitendebug.Vector2{
				X: entity.Width,
				Y: entity.Height,
			},
			Props: props,
		})
	}

	return ebitendebug.WorldSnapshot{Entities: entities}
}

func (game *game) uiSnapshot() ebitendebug.UISnapshot {
	game.mu.RLock()
	defer game.mu.RUnlock()

	return ebitendebug.UISnapshot{
		Width:  float64(game.width),
		Height: float64(game.height),
		Root: ebitendebug.UINodeSnapshot{
			ID:      "screen-root",
			Type:    "screen",
			Visible: true,
			Bounds: ebitendebug.Rect{
				X:      0,
				Y:      0,
				Width:  float64(game.width),
				Height: float64(game.height),
			},
			Children: []ebitendebug.UINodeSnapshot{
				{
					ID:      "title",
					Type:    "text",
					Text:    "Debug Bridge Example",
					Visible: true,
					Bounds: ebitendebug.Rect{
						X:      8,
						Y:      8,
						Width:  220,
						Height: 16,
					},
				},
				{
					ID:      "scene-label",
					Type:    "text",
					Text:    "Scene: " + game.currentScene,
					Visible: true,
					Bounds: ebitendebug.Rect{
						X:      8,
						Y:      24,
						Width:  180,
						Height: 16,
					},
				},
				{
					ID:      "paused-label",
					Type:    "text",
					Text:    "Paused: " + strconv.FormatBool(game.paused),
					Visible: true,
					Bounds: ebitendebug.Rect{
						X:      8,
						Y:      40,
						Width:  180,
						Height: 16,
					},
				},
				{
					ID:      "entity-list",
					Type:    "list",
					Visible: true,
					Bounds: ebitendebug.Rect{
						X:      8,
						Y:      72,
						Width:  320,
						Height: float64(len(game.entities) * 18),
					},
					Children: game.uiEntityNodesLocked(),
				},
			},
		},
	}
}

func (game *game) uiEntityNodesLocked() []ebitendebug.UINodeSnapshot {
	children := make([]ebitendebug.UINodeSnapshot, 0, len(game.entities))
	for index, entity := range game.entities {
		children = append(children, ebitendebug.UINodeSnapshot{
			ID:      entity.ID + "-row",
			Type:    "text",
			Text:    entity.ID + " visible=" + strconv.FormatBool(entity.Visible),
			Visible: true,
			Bounds: ebitendebug.Rect{
				X:      12,
				Y:      72 + float64(index*18),
				Width:  260,
				Height: 16,
			},
			Props: map[string]any{
				"entityId": entity.ID,
				"type":     entity.Type,
			},
		})
	}
	return children
}

func (game *game) entityByID(id string) *entity {
	game.mu.RLock()
	defer game.mu.RUnlock()
	return game.entityByIDLocked(id)
}

func (game *game) entityByIDLocked(id string) *entity {
	for index := range game.entities {
		if game.entities[index].ID == id {
			return &game.entities[index]
		}
	}
	return nil
}

func (game *game) hasScene(id string) bool {
	for _, scene := range game.scenes {
		if scene.ID == id {
			return true
		}
	}
	return false
}
