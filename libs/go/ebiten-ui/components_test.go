package ebitenui_test

import (
	"image/color"
	"testing"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func TestIconBuildsImageNode(t *testing.T) {
	node := ebitenui.Icon(ebitenui.IconConfig{
		ID:   "sword-icon",
		Size: 20,
		Image: ebitenui.SolidImage(20, 20, color.RGBA{
			R: 90, G: 140, B: 220, A: 255,
		}),
	})

	if got, want := node.Tag, ebitenui.TagImage; got != want {
		t.Fatalf("tag mismatch: got %q want %q", got, want)
	}
	if got, want := node.Props.Style.Width.Value, float64(20); got != want {
		t.Fatalf("width mismatch: got %v want %v", got, want)
	}
}

func TestTextareaBuildsWrappedBody(t *testing.T) {
	node := ebitenui.Textarea(ebitenui.TextareaConfig{
		ID:          "bio",
		Value:       "Line one\nLine two is longer",
		Width:       180,
		Placeholder: "Write here",
		State:       ebitenui.InteractionState{Focused: true},
	})

	dom := ebitenui.New(node)
	body, ok := dom.FindByID("bio-body")
	if !ok || body.Tag != ebitenui.TagTextBlock {
		t.Fatalf("expected textarea body text block")
	}
	if !node.Props.State.Focused {
		t.Fatalf("expected textarea focus state")
	}
}

func TestCheckboxPreservesCheckedState(t *testing.T) {
	node := ebitenui.Checkbox(ebitenui.CheckboxConfig{
		ID:      "remember",
		Label:   "Remember device",
		Checked: true,
	})

	dom := ebitenui.New(node)
	box, ok := dom.FindByID("remember-box")
	if !ok {
		t.Fatalf("expected checkbox box")
	}
	if !box.Props.State.Selected {
		t.Fatalf("expected checkbox selected state")
	}
}

func TestToggleBuildsTrackAndThumb(t *testing.T) {
	node := ebitenui.Toggle(ebitenui.ToggleConfig{
		ID:      "music-toggle",
		Label:   "Music",
		Checked: true,
	})

	dom := ebitenui.New(node)
	if _, ok := dom.FindByID("music-toggle-track"); !ok {
		t.Fatalf("expected toggle track")
	}
	thumb, ok := dom.FindByID("music-toggle-thumb")
	if !ok || !thumb.Props.State.Selected {
		t.Fatalf("expected selected toggle thumb")
	}
}

func TestSliderCreatesTrackAndThumb(t *testing.T) {
	node := ebitenui.Slider(ebitenui.SliderConfig{
		ID:    "volume",
		Label: "Volume",
		Min:   0,
		Max:   100,
		Value: 40,
		Width: 200,
	})

	layout := ebitenui.New(node).Layout(ebitenui.Viewport{Width: 240, Height: 100})
	fill, ok := layout.FindByID("volume-fill")
	if !ok {
		t.Fatalf("expected slider fill")
	}
	if fill.Frame.Width != 80 {
		t.Fatalf("expected 40%% slider fill, got %v", fill.Frame.Width)
	}
}

func TestScrollbarCreatesThumb(t *testing.T) {
	node := ebitenui.Scrollbar(ebitenui.ScrollbarConfig{
		ID:           "list-scrollbar",
		Orientation:  ebitenui.Vertical,
		Length:       120,
		ViewportSize: 30,
		ContentSize:  90,
		ScrollOffset: 30,
		Thickness:    12,
	})

	layout := ebitenui.New(node).Layout(ebitenui.Viewport{Width: 40, Height: 160})
	thumb, ok := layout.FindByID("list-scrollbar-thumb")
	if !ok {
		t.Fatalf("expected scrollbar thumb")
	}
	if thumb.Frame.Height <= 0 {
		t.Fatalf("expected positive thumb height")
	}
}

func TestDropdownBuildsOptionsWhenOpen(t *testing.T) {
	node := ebitenui.Dropdown(ebitenui.DropdownConfig{
		ID:           "resolution",
		Label:        "Resolution",
		SelectedText: "1280x720",
		Open:         true,
		Options: []ebitenui.DropdownOption{
			{ID: "res-720", Label: "1280x720"},
			{ID: "res-1080", Label: "1920x1080", State: ebitenui.InteractionState{Focused: true}},
		},
	})

	dom := ebitenui.New(node)
	if _, ok := dom.FindByID("resolution-options"); !ok {
		t.Fatalf("expected open dropdown options")
	}
	option, ok := dom.FindByID("res-1080")
	if !ok || !option.Props.State.Focused {
		t.Fatalf("expected focused dropdown option")
	}
}

func TestInputFieldUsesValueAndPlaceholderNodes(t *testing.T) {
	withValue := ebitenui.InputField(ebitenui.InputFieldConfig{
		ID:    "player-name",
		Label: "Name",
		Value: "Kim",
		State: ebitenui.InteractionState{Focused: true},
	})
	domWithValue := ebitenui.New(withValue)
	if _, ok := domWithValue.FindByID("player-name-value"); !ok {
		t.Fatalf("expected value node")
	}
	if _, ok := domWithValue.FindByID("player-name-caret"); !ok {
		t.Fatalf("expected caret node when focused")
	}

	withPlaceholder := ebitenui.InputField(ebitenui.InputFieldConfig{
		ID:          "email",
		Label:       "Email",
		Placeholder: "name@example.com",
	})
	domWithPlaceholder := ebitenui.New(withPlaceholder)
	if _, ok := domWithPlaceholder.FindByID("email-placeholder"); !ok {
		t.Fatalf("expected placeholder node")
	}
}

func TestRadioGroupMarksSelectedOption(t *testing.T) {
	node := ebitenui.RadioGroup(ebitenui.RadioGroupConfig{
		ID:          "difficulty",
		Label:       "Difficulty",
		Orientation: ebitenui.Row,
		Options: []ebitenui.RadioOption{
			{ID: "easy", Label: "Easy"},
			{ID: "normal", Label: "Normal", Selected: true},
			{ID: "hard", Label: "Hard"},
		},
	})

	dom := ebitenui.New(node)
	option, ok := dom.FindByID("normal")
	if !ok || !option.Props.State.Selected {
		t.Fatalf("expected selected radio option")
	}
}

func TestStepperBuildsDecrementValueIncrement(t *testing.T) {
	node := ebitenui.Stepper(ebitenui.StepperConfig{
		ID:    "party-size",
		Label: "Party Size",
		Value: 3,
		Min:   1,
		Max:   4,
	})

	dom := ebitenui.New(node)
	for _, id := range []string{"party-size-decrement", "party-size-value", "party-size-increment"} {
		if _, ok := dom.FindByID(id); !ok {
			t.Fatalf("expected stepper child %q", id)
		}
	}
}

func TestProgressBarUsesRatioWidth(t *testing.T) {
	node := ebitenui.ProgressBar(ebitenui.ProgressBarConfig{
		ID:      "exp-bar",
		Label:   "EXP",
		Current: 25,
		Max:     100,
		Width:   160,
	})

	layout := ebitenui.New(node).Layout(ebitenui.Viewport{Width: 180, Height: 80})
	fill, ok := layout.FindByID("exp-bar-fill")
	if !ok {
		t.Fatalf("expected progress fill")
	}
	if fill.Frame.Width != 40 {
		t.Fatalf("expected 25%% fill width, got %v", fill.Frame.Width)
	}
}

func TestDividerRespectsOrientation(t *testing.T) {
	node := ebitenui.Divider(ebitenui.DividerConfig{
		ID:          "divider",
		Orientation: ebitenui.Horizontal,
		Length:      90,
		Thickness:   2,
	})

	layout := ebitenui.New(node).Layout(ebitenui.Viewport{Width: 120, Height: 40})
	if got, want := layout.Frame.Height, float64(2); got != want {
		t.Fatalf("divider thickness mismatch: got %v want %v", got, want)
	}
}

func TestGridBuildsRowsAndCells(t *testing.T) {
	node := ebitenui.Grid(ebitenui.GridConfig{
		ID:      "grid",
		Columns: 2,
		Gap:     8,
		Children: []*ebitenui.Node{
			ebitenui.Text("A", ebitenui.Props{ID: "grid-a"}),
			ebitenui.Text("B", ebitenui.Props{ID: "grid-b"}),
			ebitenui.Text("C", ebitenui.Props{ID: "grid-c"}),
		},
	})

	dom := ebitenui.New(node)
	if _, ok := dom.FindByID("grid-row-0"); !ok {
		t.Fatalf("expected grid row")
	}
	if _, ok := dom.FindByID("grid-c"); !ok {
		t.Fatalf("expected grid cell")
	}
}

func TestListUsesConfiguredOrientation(t *testing.T) {
	node := ebitenui.List(ebitenui.ListConfig{
		ID:          "menu-list",
		Orientation: ebitenui.Column,
		Gap:         6,
		Items: []*ebitenui.Node{
			ebitenui.Text("One", ebitenui.Props{ID: "list-1"}),
			ebitenui.Text("Two", ebitenui.Props{ID: "list-2"}),
		},
	})

	layout := ebitenui.New(node).Layout(ebitenui.Viewport{Width: 200, Height: 80})
	first, ok := layout.FindByID("list-1")
	if !ok {
		t.Fatalf("expected list item")
	}
	second, ok := layout.FindByID("list-2")
	if !ok {
		t.Fatalf("expected second list item")
	}
	if !(second.Frame.Y > first.Frame.Y) {
		t.Fatalf("expected vertical list flow")
	}
}

func TestVirtualListRendersVisibleWindow(t *testing.T) {
	node := ebitenui.VirtualList(ebitenui.VirtualListConfig{
		ID:           "virtual",
		StartIndex:   2,
		VisibleCount: 3,
		Orientation:  ebitenui.Column,
		ItemBuilder: func(index int) *ebitenui.Node {
			return ebitenui.Text("item", ebitenui.Props{ID: ebitenui.ComponentID("virtual-item", index)})
		},
		TotalCount: 10,
	})

	dom := ebitenui.New(node)
	if _, ok := dom.FindByID("virtual-item-2"); !ok {
		t.Fatalf("expected visible item")
	}
	if _, ok := dom.FindByID("virtual-item-5"); ok {
		t.Fatalf("expected non-visible item to be absent")
	}
}

func TestModalBuildsOverlayAndContent(t *testing.T) {
	node := ebitenui.Modal(ebitenui.ModalConfig{
		ID:      "settings-modal",
		Open:    true,
		Width:   260,
		Height:  140,
		Title:   "Settings",
		Content: ebitenui.Text("Body", ebitenui.Props{ID: "settings-modal-body"}),
	})

	dom := ebitenui.New(node)
	for _, id := range []string{"settings-modal-overlay", "settings-modal-content", "settings-modal-body"} {
		if _, ok := dom.FindByID(id); !ok {
			t.Fatalf("expected modal node %q", id)
		}
	}
}

func TestTooltipBuildsTitleAndDescription(t *testing.T) {
	node := ebitenui.Tooltip(ebitenui.TooltipConfig{
		ID:          "item-tooltip",
		Title:       "Iron Bow",
		Description: "A reliable ranged weapon.",
		Width:       200,
	})

	dom := ebitenui.New(node)
	if _, ok := dom.FindByID("item-tooltip-title"); !ok {
		t.Fatalf("expected tooltip title")
	}
	if _, ok := dom.FindByID("item-tooltip-description"); !ok {
		t.Fatalf("expected tooltip description")
	}
}

func TestContextMenuBuildsSelectableItems(t *testing.T) {
	node := ebitenui.ContextMenu(ebitenui.ContextMenuConfig{
		ID: "slot-menu",
		Items: []ebitenui.ContextMenuItem{
			{ID: "use-item", Label: "Use"},
			{ID: "drop-item", Label: "Drop", State: ebitenui.InteractionState{Focused: true}},
		},
	})

	dom := ebitenui.New(node)
	item, ok := dom.FindByID("drop-item")
	if !ok || !item.Props.State.Focused {
		t.Fatalf("expected focused context menu item")
	}
}

func TestTabsBuildsOnlySelectedPanel(t *testing.T) {
	node := ebitenui.Tabs(ebitenui.TabsConfig{
		ID:            "tabs",
		SelectedIndex: 1,
		Tabs: []ebitenui.TabConfig{
			{ID: "tab-0", Label: "Stats", Content: ebitenui.Text("stats", ebitenui.Props{ID: "stats-panel"})},
			{ID: "tab-1", Label: "Skills", Content: ebitenui.Text("skills", ebitenui.Props{ID: "skills-panel"})},
		},
	})

	dom := ebitenui.New(node)
	tab, ok := dom.FindByID("tab-1")
	if !ok || !tab.Props.State.Selected {
		t.Fatalf("expected selected tab")
	}
	if _, ok := dom.FindByID("skills-panel"); !ok {
		t.Fatalf("expected selected panel")
	}
	if _, ok := dom.FindByID("stats-panel"); ok {
		t.Fatalf("expected non-selected panel to be absent")
	}
}

func TestAccordionOnlyExpandsSelectedSections(t *testing.T) {
	node := ebitenui.Accordion(ebitenui.AccordionConfig{
		ID: "accordion",
		Sections: []ebitenui.AccordionSection{
			{ID: "section-a", Label: "A", Expanded: false, Content: ebitenui.Text("A content", ebitenui.Props{ID: "section-a-content"})},
			{ID: "section-b", Label: "B", Expanded: true, Content: ebitenui.Text("B content", ebitenui.Props{ID: "section-b-content"})},
		},
	})

	dom := ebitenui.New(node)
	if _, ok := dom.FindByID("section-b-content"); !ok {
		t.Fatalf("expected expanded content")
	}
	if _, ok := dom.FindByID("section-a-content"); ok {
		t.Fatalf("expected collapsed content to be absent")
	}
}

func TestBadgeAndChipBuildExpectedNodes(t *testing.T) {
	badge := ebitenui.Badge(ebitenui.BadgeConfig{
		ID:    "new-badge",
		Label: "NEW",
	})
	chip := ebitenui.Chip(ebitenui.ChipConfig{
		ID:          "filter-chip",
		Label:       "Fire",
		Dismissible: true,
	})

	if _, ok := ebitenui.New(badge).FindByID("new-badge-label"); !ok {
		t.Fatalf("expected badge label")
	}
	if _, ok := ebitenui.New(chip).FindByID("filter-chip-dismiss"); !ok {
		t.Fatalf("expected dismissible chip")
	}
}
