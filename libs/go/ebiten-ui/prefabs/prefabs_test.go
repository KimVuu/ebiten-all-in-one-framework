package prefabs_test

import (
	"image/color"
	"testing"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
	"github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui/prefabs"
)

func TestDialogBuildsTitleBodyAndActions(t *testing.T) {
	node := prefabs.Dialog(prefabs.DialogConfig{
		ID:    "dialog",
		Title: "Exit game?",
		Body:  "Unsaved progress will be lost.",
		Width: 280,
		Actions: []prefabs.DialogAction{
			{ID: "cancel", Label: "Cancel"},
			{ID: "confirm", Label: "Confirm", State: ebitenui.InteractionState{Selected: true}},
		},
	})

	dom := ebitenui.New(node)
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

	layout := ebitenui.New(node).Layout(ebitenui.Viewport{Width: 220, Height: 120})
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
			{ID: "slot-1", Label: "Potion", Quantity: 3, State: ebitenui.InteractionState{Selected: true}},
			{ID: "slot-2", Label: "Ether", Quantity: 1},
			{ID: "slot-3", Label: "Key", Quantity: 1},
			{ID: "slot-4", Label: "Gem", Quantity: 2},
		},
	})

	dom := ebitenui.New(node)
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
			{ID: "resume", Label: "Resume", State: ebitenui.InteractionState{Focused: true}},
			{ID: "settings", Label: "Settings"},
			{ID: "quit", Label: "Quit", State: ebitenui.InteractionState{Disabled: true}},
		},
		Footer: "Press B to close",
	})

	dom := ebitenui.New(node)
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

	dom := ebitenui.New(node)
	scroll, ok := dom.FindByID("settings-scroll")
	if !ok || scroll.Tag != ebitenui.TagScrollView {
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

	dom := ebitenui.New(node)
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

func TestPanelUsesThemeSurfaceAndBorder(t *testing.T) {
	panelBackground := color.RGBA{R: 40, G: 46, B: 60, A: 255}
	panelBorder := color.RGBA{R: 123, G: 132, B: 170, A: 255}
	titleColor := color.RGBA{R: 250, G: 245, B: 236, A: 255}

	theme := ebitenui.DefaultTheme()
	theme.Components.Panel.Background = panelBackground
	theme.Components.Panel.Border = panelBorder
	theme.Components.Panel.TitleText = titleColor

	node := prefabs.Panel(prefabs.PanelConfig{
		ID:    "panel",
		Title: "Status",
		Theme: &theme,
	})

	dom := ebitenui.New(node)
	panel, ok := dom.FindByID("panel")
	if !ok {
		t.Fatalf("expected panel node")
	}
	if !samePrefabColor(panel.Props.Style.BackgroundColor, panelBackground) {
		t.Fatalf("expected themed panel background")
	}
	if !samePrefabColor(panel.Props.Style.BorderColor, panelBorder) {
		t.Fatalf("expected themed panel border")
	}
	title, ok := dom.FindByID("panel-title")
	if !ok {
		t.Fatalf("expected panel title")
	}
	if !samePrefabColor(title.Props.Style.Color, titleColor) {
		t.Fatalf("expected themed panel title color")
	}
}

func TestDialogUsesThemeForSelectedAction(t *testing.T) {
	selectedBackground := color.RGBA{R: 255, G: 205, B: 98, A: 255}
	selectedText := color.RGBA{R: 18, G: 22, B: 32, A: 255}

	theme := ebitenui.DefaultTheme()
	theme.Components.MenuButton.Selected.Background = selectedBackground
	theme.Components.MenuButton.Selected.Text = selectedText

	node := prefabs.Dialog(prefabs.DialogConfig{
		ID:    "dialog-theme",
		Title: "Save changes?",
		Body:  "Progress will be written to disk.",
		Theme: &theme,
		Actions: []prefabs.DialogAction{
			{ID: "dialog-theme-cancel", Label: "Cancel"},
			{ID: "dialog-theme-confirm", Label: "Confirm", State: ebitenui.InteractionState{Selected: true}},
		},
	})

	dom := ebitenui.New(node)
	action, ok := dom.FindByID("dialog-theme-confirm")
	if !ok {
		t.Fatalf("expected selected action")
	}
	if !samePrefabColor(action.Props.Style.BackgroundColor, selectedBackground) {
		t.Fatalf("expected themed selected action background")
	}
	label, ok := dom.FindByID("dialog-theme-confirm-label")
	if !ok {
		t.Fatalf("expected selected action label")
	}
	if !samePrefabColor(label.Props.Style.Color, selectedText) {
		t.Fatalf("expected themed selected action text")
	}
}

func TestHUDBarUsesThemeTintWhenTintUnset(t *testing.T) {
	fill := color.RGBA{R: 116, G: 213, B: 166, A: 255}
	track := color.RGBA{R: 30, G: 38, B: 50, A: 255}

	theme := ebitenui.DefaultTheme()
	theme.Components.HUDBar.Fill = fill
	theme.Components.HUDBar.Track = track

	node := prefabs.HUDBar(prefabs.HUDBarConfig{
		ID:      "hud-theme",
		Label:   "Shield",
		Current: 40,
		Max:     100,
		Width:   180,
		Theme:   &theme,
	})

	dom := ebitenui.New(node)
	trackNode, ok := dom.FindByID("hud-theme-track")
	if !ok {
		t.Fatalf("expected hud track")
	}
	if !samePrefabColor(trackNode.Props.Style.BackgroundColor, track) {
		t.Fatalf("expected themed hud track")
	}
	fillNode, ok := dom.FindByID("hud-theme-fill")
	if !ok {
		t.Fatalf("expected hud fill")
	}
	if !samePrefabColor(fillNode.Props.Style.BackgroundColor, fill) {
		t.Fatalf("expected themed hud fill")
	}
}

func countNodesByID(node *ebitenui.Node, id string) int {
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

func samePrefabColor(got color.Color, want color.Color) bool {
	if got == nil || want == nil {
		return got == want
	}
	gotRGBA := color.NRGBAModel.Convert(got).(color.NRGBA)
	wantRGBA := color.NRGBAModel.Convert(want).(color.NRGBA)
	return gotRGBA == wantRGBA
}
