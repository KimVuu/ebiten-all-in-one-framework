package main

import (
	"image/color"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
	"github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui/prefabs"
)

func prefabsSection() *ebitenui.Node {
	return ebitenui.Section(ebitenui.Props{
		ID:    "prefabs-section",
		Style: showcaseGroupStyle(),
	},
		ebitenui.Text("Game UI Prefabs", ebitenui.Props{
			ID:    "prefabs-section-title",
			Style: showcaseGroupTitleStyle(),
		}),
		ebitenui.TextBlock("Dialog, HUD, inventory, pause, settings, and tooltip prefabs sit on top of the same DOM tree primitives.", ebitenui.Props{
			ID:    "prefabs-section-copy",
			Style: showcaseGroupCopyStyle(),
		}),
		prefabs.Dialog(prefabs.DialogConfig{
			ID:    "dialog-demo",
			Title: "Return to title?",
			Body:  "You will lose unsaved progress from this dungeon run.",
			Width: 360,
			Actions: []prefabs.DialogAction{
				{ID: "dialog-cancel", Label: "Cancel"},
				{ID: "dialog-confirm", Label: "Confirm", State: ebitenui.InteractionState{Selected: true}},
			},
		}),
		prefabs.HUDBar(prefabs.HUDBarConfig{
			ID:      "hud-demo",
			Label:   "HP",
			Current: 84,
			Max:     120,
			Width:   280,
			Tint:    color.RGBA{R: 234, G: 94, B: 110, A: 255},
		}),
		prefabs.InventoryGrid(prefabs.InventoryGridConfig{
			ID:       "inventory-demo",
			Title:    "Inventory",
			Columns:  4,
			CellSize: 64,
			Slots: []prefabs.InventorySlot{
				{ID: "inv-slot-1", Label: "Potion", Quantity: 3, State: ebitenui.InteractionState{Selected: true}},
				{ID: "inv-slot-2", Label: "Ether", Quantity: 1},
				{ID: "inv-slot-3", Label: "Key", Quantity: 1},
				{ID: "inv-slot-4", Label: "Gem", Quantity: 2},
				{ID: "inv-slot-5", Label: "Map", Quantity: 1},
				{ID: "inv-slot-6", Label: "Bomb", Quantity: 5},
			},
		}),
		prefabs.PauseMenu(prefabs.PauseMenuConfig{
			ID:    "pause-demo",
			Title: "Paused",
			Width: 320,
			Items: []prefabs.MenuItem{
				{ID: "pause-resume", Label: "Resume", State: ebitenui.InteractionState{Focused: true}},
				{ID: "pause-settings", Label: "Settings"},
				{ID: "pause-quit", Label: "Quit", State: ebitenui.InteractionState{Disabled: true}},
			},
			Footer: "Gamepad: B to close",
		}),
		prefabs.SettingsPanel(prefabs.SettingsPanelConfig{
			ID:     "settings-demo",
			Title:  "Settings",
			Width:  420,
			Height: 240,
			Sections: []prefabs.SettingSection{
				{
					Title:       "Audio",
					Description: "Balance master, music, and effects.",
					Options: []prefabs.SettingOption{
						{Label: "Master", Value: "80%"},
						{Label: "Music", Value: "65%"},
					},
				},
				{
					Title:       "Controls",
					Description: "Mixed mouse, keyboard, and gamepad navigation.",
					Options: []prefabs.SettingOption{
						{Label: "Move", Value: "WASD / Left Stick"},
						{Label: "Confirm", Value: "Enter / A"},
					},
				},
			},
		}),
		prefabs.Tooltip(prefabs.TooltipConfig{
			ID:          "tooltip-demo",
			Title:       "Crystal Bow",
			Description: "A refined ranged weapon with high precision and low draw delay.",
			Width:       260,
			Stats: []prefabs.TooltipStat{
				{Label: "ATK", Value: "+18"},
				{Label: "Range", Value: "Long"},
				{Label: "Weight", Value: "Light"},
			},
		}),
	)
}
