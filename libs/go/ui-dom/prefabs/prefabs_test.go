package prefabs_test

import (
	"image/color"
	"testing"

	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
	"github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom/prefabs"
)

func TestDialogBuildsTitleBodyAndActions(t *testing.T) {
	node := prefabs.Dialog(prefabs.DialogConfig{
		ID:    "dialog",
		Title: "Exit game?",
		Body:  "Unsaved progress will be lost.",
		Width: 280,
		Actions: []prefabs.DialogAction{
			{ID: "cancel", Label: "Cancel"},
			{ID: "confirm", Label: "Confirm", State: uidom.InteractionState{Selected: true}},
		},
	})

	dom := uidom.New(node)
	if _, ok := dom.FindByID("dialog-title"); !ok {
		t.Fatalf("expected dialog title")
	}
	if _, ok := dom.FindByID("confirm"); !ok {
		t.Fatalf("expected confirm button")
	}
	confirm, ok := dom.FindByID("confirm")
	if !ok || !confirm.Props.State.Selected {
		t.Fatalf("expected selected confirm action")
	}
}

func TestHUDBarReflectsRatioInLayout(t *testing.T) {
	node := prefabs.HUDBar(prefabs.HUDBarConfig{
		ID:      "hp-bar",
		Label:   "HP",
		Current: 75,
		Max:     100,
		Width:   200,
		Tint:    color.RGBA{R: 230, G: 70, B: 90, A: 255},
	})

	layout := uidom.New(node).Layout(uidom.Viewport{Width: 220, Height: 120})
	fill, ok := layout.FindByID("hp-bar-fill")
	if !ok {
		t.Fatalf("expected fill node")
	}
	if fill.Frame.Width != 150 {
		t.Fatalf("expected 75%% fill width, got %v", fill.Frame.Width)
	}
}

func TestInventoryGridBuildsRowsAndSlots(t *testing.T) {
	node := prefabs.InventoryGrid(prefabs.InventoryGridConfig{
		ID:       "inventory",
		Title:    "Inventory",
		Columns:  3,
		CellSize: 48,
		Slots: []prefabs.InventorySlot{
			{ID: "slot-1", Label: "Potion", Quantity: 3, State: uidom.InteractionState{Selected: true}},
			{ID: "slot-2", Label: "Ether", Quantity: 1},
			{ID: "slot-3", Label: "Key", Quantity: 1},
			{ID: "slot-4", Label: "Gem", Quantity: 2},
		},
	})

	dom := uidom.New(node)
	if _, ok := dom.FindByID("inventory-row-0"); !ok {
		t.Fatalf("expected first inventory row")
	}
	slot, ok := dom.FindByID("slot-1")
	if !ok || !slot.Props.State.Selected {
		t.Fatalf("expected selected inventory slot")
	}
}

func TestPauseMenuBuildsSelectableMenu(t *testing.T) {
	node := prefabs.PauseMenu(prefabs.PauseMenuConfig{
		ID:    "pause-menu",
		Title: "Paused",
		Width: 280,
		Items: []prefabs.MenuItem{
			{ID: "resume", Label: "Resume", State: uidom.InteractionState{Focused: true}},
			{ID: "settings", Label: "Settings"},
			{ID: "quit", Label: "Quit", State: uidom.InteractionState{Disabled: true}},
		},
		Footer: "Press B to close",
	})

	dom := uidom.New(node)
	item, ok := dom.FindByID("resume")
	if !ok || !item.Props.State.Focused {
		t.Fatalf("expected focused resume item")
	}
	quit, ok := dom.FindByID("quit")
	if !ok || !quit.Props.State.Disabled {
		t.Fatalf("expected disabled quit item")
	}
}

func TestSettingsPanelUsesScrollViewAndSections(t *testing.T) {
	node := prefabs.SettingsPanel(prefabs.SettingsPanelConfig{
		ID:     "settings",
		Title:  "Settings",
		Width:  320,
		Height: 220,
		Sections: []prefabs.SettingSection{
			{
				Title:       "Audio",
				Description: "Master volume and effect levels.",
				Options: []prefabs.SettingOption{
					{Label: "Master", Value: "80%"},
					{Label: "SFX", Value: "65%"},
				},
			},
		},
	})

	dom := uidom.New(node)
	scroll, ok := dom.FindByID("settings-scroll")
	if !ok || scroll.Tag != uidom.TagScrollView {
		t.Fatalf("expected settings scroll view")
	}
	if _, ok := dom.FindByID("settings-section-0"); !ok {
		t.Fatalf("expected first settings section")
	}
}

func TestTooltipBuildsTitleDescriptionAndStats(t *testing.T) {
	node := prefabs.Tooltip(prefabs.TooltipConfig{
		ID:          "tooltip",
		Title:       "Iron Sword",
		Description: "A steady starter blade for close combat.",
		Width:       220,
		Stats: []prefabs.TooltipStat{
			{Label: "ATK", Value: "+12"},
			{Label: "Weight", Value: "Normal"},
		},
	})

	dom := uidom.New(node)
	if _, ok := dom.FindByID("tooltip-title"); !ok {
		t.Fatalf("expected tooltip title")
	}
	if _, ok := dom.FindByID("tooltip-stat-0"); !ok {
		t.Fatalf("expected tooltip stat row")
	}
	if got := countNodesByID(node, "tooltip-title"); got != 1 {
		t.Fatalf("expected tooltip title to appear once, got %d", got)
	}
}

func countNodesByID(node *uidom.Node, id string) int {
	if node == nil {
		return 0
	}
	count := 0
	if node.Props.ID == id {
		count++
	}
	for _, child := range node.Children {
		count += countNodesByID(child, id)
	}
	return count
}
