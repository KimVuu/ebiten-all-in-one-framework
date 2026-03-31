package uidom_test

import (
	"image/color"
	"testing"

	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
)

func TestIconBuildsImageNode(t *testing.T) {
	node := uidom.Icon(uidom.IconConfig{
		ID:   "sword-icon",
		Size: 20,
		Image: uidom.SolidImage(20, 20, color.RGBA{
			R: 90, G: 140, B: 220, A: 255,
		}),
	})

	if got, want := node.Tag, uidom.TagImage; got != want {
		t.Fatalf("tag mismatch: got %q want %q", got, want)
	}
	if got, want := node.Props.Style.Width.Value, float64(20); got != want {
		t.Fatalf("width mismatch: got %v want %v", got, want)
	}
}

func TestTextareaBuildsWrappedBody(t *testing.T) {
	node := uidom.Textarea(uidom.TextareaConfig{
		ID:          "bio",
		Value:       "Line one\nLine two is longer",
		Width:       180,
		Placeholder: "Write here",
		State:       uidom.InteractionState{Focused: true},
	})

	dom := uidom.New(node)
	body, ok := dom.FindByID("bio-body")
	if !ok || body.Tag != uidom.TagTextBlock {
		t.Fatalf("expected textarea body text block")
	}
	if !node.Props.State.Focused {
		t.Fatalf("expected textarea focus state")
	}
}

func TestCheckboxPreservesCheckedState(t *testing.T) {
	node := uidom.Checkbox(uidom.CheckboxConfig{
		ID:      "remember",
		Label:   "Remember device",
		Checked: true,
	})

	dom := uidom.New(node)
	box, ok := dom.FindByID("remember-box")
	if !ok {
		t.Fatalf("expected checkbox box")
	}
	if !box.Props.State.Selected {
		t.Fatalf("expected checkbox selected state")
	}
}

func TestToggleBuildsTrackAndThumb(t *testing.T) {
	node := uidom.Toggle(uidom.ToggleConfig{
		ID:      "music-toggle",
		Label:   "Music",
		Checked: true,
	})

	dom := uidom.New(node)
	if _, ok := dom.FindByID("music-toggle-track"); !ok {
		t.Fatalf("expected toggle track")
	}
	thumb, ok := dom.FindByID("music-toggle-thumb")
	if !ok || !thumb.Props.State.Selected {
		t.Fatalf("expected selected toggle thumb")
	}
}

func TestSliderCreatesTrackAndThumb(t *testing.T) {
	node := uidom.Slider(uidom.SliderConfig{
		ID:    "volume",
		Label: "Volume",
		Min:   0,
		Max:   100,
		Value: 40,
		Width: 200,
	})

	layout := uidom.New(node).Layout(uidom.Viewport{Width: 240, Height: 100})
	fill, ok := layout.FindByID("volume-fill")
	if !ok {
		t.Fatalf("expected slider fill")
	}
	if fill.Frame.Width != 80 {
		t.Fatalf("expected 40%% slider fill, got %v", fill.Frame.Width)
	}
}

func TestScrollbarCreatesThumb(t *testing.T) {
	node := uidom.Scrollbar(uidom.ScrollbarConfig{
		ID:           "list-scrollbar",
		Orientation:  uidom.Vertical,
		Length:       120,
		ViewportSize: 30,
		ContentSize:  90,
		ScrollOffset: 30,
		Thickness:    12,
	})

	layout := uidom.New(node).Layout(uidom.Viewport{Width: 40, Height: 160})
	thumb, ok := layout.FindByID("list-scrollbar-thumb")
	if !ok {
		t.Fatalf("expected scrollbar thumb")
	}
	if thumb.Frame.Height <= 0 {
		t.Fatalf("expected positive thumb height")
	}
}

func TestDropdownBuildsOptionsWhenOpen(t *testing.T) {
	node := uidom.Dropdown(uidom.DropdownConfig{
		ID:           "resolution",
		Label:        "Resolution",
		SelectedText: "1280x720",
		Open:         true,
		Options: []uidom.DropdownOption{
			{ID: "res-720", Label: "1280x720"},
			{ID: "res-1080", Label: "1920x1080", State: uidom.InteractionState{Focused: true}},
		},
	})

	dom := uidom.New(node)
	if _, ok := dom.FindByID("resolution-options"); !ok {
		t.Fatalf("expected open dropdown options")
	}
	option, ok := dom.FindByID("res-1080")
	if !ok || !option.Props.State.Focused {
		t.Fatalf("expected focused dropdown option")
	}
}

func TestInputFieldUsesValueAndPlaceholderNodes(t *testing.T) {
	withValue := uidom.InputField(uidom.InputFieldConfig{
		ID:    "player-name",
		Label: "Name",
		Value: "Kim",
		State: uidom.InteractionState{Focused: true},
	})
	domWithValue := uidom.New(withValue)
	if _, ok := domWithValue.FindByID("player-name-value"); !ok {
		t.Fatalf("expected value node")
	}
	if _, ok := domWithValue.FindByID("player-name-caret"); !ok {
		t.Fatalf("expected caret node when focused")
	}

	withPlaceholder := uidom.InputField(uidom.InputFieldConfig{
		ID:          "email",
		Label:       "Email",
		Placeholder: "name@example.com",
	})
	domWithPlaceholder := uidom.New(withPlaceholder)
	if _, ok := domWithPlaceholder.FindByID("email-placeholder"); !ok {
		t.Fatalf("expected placeholder node")
	}
}

func TestRadioGroupMarksSelectedOption(t *testing.T) {
	node := uidom.RadioGroup(uidom.RadioGroupConfig{
		ID:          "difficulty",
		Label:       "Difficulty",
		Orientation: uidom.Row,
		Options: []uidom.RadioOption{
			{ID: "easy", Label: "Easy"},
			{ID: "normal", Label: "Normal", Selected: true},
			{ID: "hard", Label: "Hard"},
		},
	})

	dom := uidom.New(node)
	option, ok := dom.FindByID("normal")
	if !ok || !option.Props.State.Selected {
		t.Fatalf("expected selected radio option")
	}
}

func TestStepperBuildsDecrementValueIncrement(t *testing.T) {
	node := uidom.Stepper(uidom.StepperConfig{
		ID:    "party-size",
		Label: "Party Size",
		Value: 3,
		Min:   1,
		Max:   4,
	})

	dom := uidom.New(node)
	for _, id := range []string{"party-size-decrement", "party-size-value", "party-size-increment"} {
		if _, ok := dom.FindByID(id); !ok {
			t.Fatalf("expected stepper child %q", id)
		}
	}
}

func TestProgressBarUsesRatioWidth(t *testing.T) {
	node := uidom.ProgressBar(uidom.ProgressBarConfig{
		ID:      "exp-bar",
		Label:   "EXP",
		Current: 25,
		Max:     100,
		Width:   160,
	})

	layout := uidom.New(node).Layout(uidom.Viewport{Width: 180, Height: 80})
	fill, ok := layout.FindByID("exp-bar-fill")
	if !ok {
		t.Fatalf("expected progress fill")
	}
	if fill.Frame.Width != 40 {
		t.Fatalf("expected 25%% fill width, got %v", fill.Frame.Width)
	}
}

func TestDividerRespectsOrientation(t *testing.T) {
	node := uidom.Divider(uidom.DividerConfig{
		ID:          "divider",
		Orientation: uidom.Horizontal,
		Length:      90,
		Thickness:   2,
	})

	layout := uidom.New(node).Layout(uidom.Viewport{Width: 120, Height: 40})
	if got, want := layout.Frame.Height, float64(2); got != want {
		t.Fatalf("divider thickness mismatch: got %v want %v", got, want)
	}
}

func TestGridBuildsRowsAndCells(t *testing.T) {
	node := uidom.Grid(uidom.GridConfig{
		ID:      "grid",
		Columns: 2,
		Gap:     8,
		Children: []*uidom.Node{
			uidom.Text("A", uidom.Props{ID: "grid-a"}),
			uidom.Text("B", uidom.Props{ID: "grid-b"}),
			uidom.Text("C", uidom.Props{ID: "grid-c"}),
		},
	})

	dom := uidom.New(node)
	if _, ok := dom.FindByID("grid-row-0"); !ok {
		t.Fatalf("expected grid row")
	}
	if _, ok := dom.FindByID("grid-c"); !ok {
		t.Fatalf("expected grid cell")
	}
}

func TestListUsesConfiguredOrientation(t *testing.T) {
	node := uidom.List(uidom.ListConfig{
		ID:          "menu-list",
		Orientation: uidom.Column,
		Gap:         6,
		Items: []*uidom.Node{
			uidom.Text("One", uidom.Props{ID: "list-1"}),
			uidom.Text("Two", uidom.Props{ID: "list-2"}),
		},
	})

	layout := uidom.New(node).Layout(uidom.Viewport{Width: 200, Height: 80})
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
	node := uidom.VirtualList(uidom.VirtualListConfig{
		ID:           "virtual",
		StartIndex:   2,
		VisibleCount: 3,
		Orientation:  uidom.Column,
		ItemBuilder: func(index int) *uidom.Node {
			return uidom.Text("item", uidom.Props{ID: uidom.ComponentID("virtual-item", index)})
		},
		TotalCount: 10,
	})

	dom := uidom.New(node)
	if _, ok := dom.FindByID("virtual-item-2"); !ok {
		t.Fatalf("expected visible item")
	}
	if _, ok := dom.FindByID("virtual-item-5"); ok {
		t.Fatalf("expected non-visible item to be absent")
	}
}

func TestModalBuildsOverlayAndContent(t *testing.T) {
	node := uidom.Modal(uidom.ModalConfig{
		ID:      "settings-modal",
		Open:    true,
		Width:   260,
		Height:  140,
		Title:   "Settings",
		Content: uidom.Text("Body", uidom.Props{ID: "settings-modal-body"}),
	})

	dom := uidom.New(node)
	for _, id := range []string{"settings-modal-overlay", "settings-modal-content", "settings-modal-body"} {
		if _, ok := dom.FindByID(id); !ok {
			t.Fatalf("expected modal node %q", id)
		}
	}
}

func TestTooltipBuildsTitleAndDescription(t *testing.T) {
	node := uidom.Tooltip(uidom.TooltipConfig{
		ID:          "item-tooltip",
		Title:       "Iron Bow",
		Description: "A reliable ranged weapon.",
		Width:       200,
	})

	dom := uidom.New(node)
	if _, ok := dom.FindByID("item-tooltip-title"); !ok {
		t.Fatalf("expected tooltip title")
	}
	if _, ok := dom.FindByID("item-tooltip-description"); !ok {
		t.Fatalf("expected tooltip description")
	}
}

func TestContextMenuBuildsSelectableItems(t *testing.T) {
	node := uidom.ContextMenu(uidom.ContextMenuConfig{
		ID: "slot-menu",
		Items: []uidom.ContextMenuItem{
			{ID: "use-item", Label: "Use"},
			{ID: "drop-item", Label: "Drop", State: uidom.InteractionState{Focused: true}},
		},
	})

	dom := uidom.New(node)
	item, ok := dom.FindByID("drop-item")
	if !ok || !item.Props.State.Focused {
		t.Fatalf("expected focused context menu item")
	}
}

func TestTabsBuildsOnlySelectedPanel(t *testing.T) {
	node := uidom.Tabs(uidom.TabsConfig{
		ID:            "tabs",
		SelectedIndex: 1,
		Tabs: []uidom.TabConfig{
			{ID: "tab-0", Label: "Stats", Content: uidom.Text("stats", uidom.Props{ID: "stats-panel"})},
			{ID: "tab-1", Label: "Skills", Content: uidom.Text("skills", uidom.Props{ID: "skills-panel"})},
		},
	})

	dom := uidom.New(node)
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
	node := uidom.Accordion(uidom.AccordionConfig{
		ID: "accordion",
		Sections: []uidom.AccordionSection{
			{ID: "section-a", Label: "A", Expanded: false, Content: uidom.Text("A content", uidom.Props{ID: "section-a-content"})},
			{ID: "section-b", Label: "B", Expanded: true, Content: uidom.Text("B content", uidom.Props{ID: "section-b-content"})},
		},
	})

	dom := uidom.New(node)
	if _, ok := dom.FindByID("section-b-content"); !ok {
		t.Fatalf("expected expanded content")
	}
	if _, ok := dom.FindByID("section-a-content"); ok {
		t.Fatalf("expected collapsed content to be absent")
	}
}

func TestBadgeAndChipBuildExpectedNodes(t *testing.T) {
	badge := uidom.Badge(uidom.BadgeConfig{
		ID:    "new-badge",
		Label: "NEW",
	})
	chip := uidom.Chip(uidom.ChipConfig{
		ID:          "filter-chip",
		Label:       "Fire",
		Dismissible: true,
	})

	if _, ok := uidom.New(badge).FindByID("new-badge-label"); !ok {
		t.Fatalf("expected badge label")
	}
	if _, ok := uidom.New(chip).FindByID("filter-chip-dismiss"); !ok {
		t.Fatalf("expected dismissible chip")
	}
}
