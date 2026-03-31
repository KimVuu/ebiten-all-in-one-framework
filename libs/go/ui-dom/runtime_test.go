package uidom_test

import (
	"testing"

	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
)

func TestRuntimeDispatchesButtonClickAndFocus(t *testing.T) {
	clicks := 0
	dom := uidom.New(
		uidom.Div(uidom.Props{
			ID: "root",
			Style: uidom.Style{
				Width:   uidom.Px(200),
				Height:  uidom.Px(80),
				Padding: uidom.All(8),
			},
		},
			uidom.InteractiveButton(uidom.Props{
				ID: "play-button",
				Style: uidom.Style{
					Width:  uidom.Px(120),
					Height: uidom.Px(40),
				},
				Handlers: uidom.EventHandlers{
					OnClick: func(ctx uidom.EventContext) {
						clicks++
					},
				},
			},
				uidom.Text("Play", uidom.Props{ID: "play-button-label"}),
			),
		),
	)

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 200, Height: 80}

	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: 20, PointerY: 20})
	button, ok := dom.FindByID("play-button")
	if !ok || !button.Props.State.Hovered {
		t.Fatalf("expected button hover state")
	}

	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: 20, PointerY: 20, PointerDown: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: 20, PointerY: 20})

	if clicks != 1 {
		t.Fatalf("expected click handler once, got %d", clicks)
	}
	if got, want := runtime.FocusedID(), "play-button"; got != want {
		t.Fatalf("focus mismatch: got %q want %q", got, want)
	}
}

func TestRuntimeCheckboxOnChangeUsesRuntimeValue(t *testing.T) {
	var values []bool
	dom := uidom.New(uidom.Checkbox(uidom.CheckboxConfig{
		ID:      "remember",
		Label:   "Remember",
		Checked: false,
		OnChange: func(value bool) {
			values = append(values, value)
		},
	}))

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 220, Height: 48}

	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: 12, PointerY: 12, PointerDown: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: 12, PointerY: 12})
	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: 12, PointerY: 12, PointerDown: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: 12, PointerY: 12})

	if len(values) != 2 {
		t.Fatalf("expected two toggle values, got %d", len(values))
	}
	if values[0] != true || values[1] != false {
		t.Fatalf("unexpected toggle values: %#v", values)
	}
}

func TestRuntimeSliderOnChangeTracksPointerPosition(t *testing.T) {
	var changed []float64
	dom := uidom.New(uidom.Slider(uidom.SliderConfig{
		ID:    "volume",
		Label: "Volume",
		Min:   0,
		Max:   100,
		Value: 20,
		Width: 200,
		OnChange: func(value float64) {
			changed = append(changed, value)
		},
	}))

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 220, Height: 80}
	layout := runtime.Update(dom, viewport, uidom.InputSnapshot{})
	track, ok := layout.FindByID("volume-track")
	if !ok {
		t.Fatalf("expected slider track")
	}

	targetX := track.Frame.X + track.Frame.Width*0.75
	targetY := track.Frame.Y + track.Frame.Height*0.5
	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: targetX, PointerY: targetY, PointerDown: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: targetX, PointerY: targetY})

	if len(changed) == 0 {
		t.Fatalf("expected slider change event")
	}
	got := changed[len(changed)-1]
	if got < 70 || got > 80 {
		t.Fatalf("expected slider value near 75, got %v", got)
	}
}

func TestRuntimeInputFieldDispatchesTextBackspaceAndSubmit(t *testing.T) {
	var changes []string
	submits := 0
	dom := uidom.New(uidom.InputField(uidom.InputFieldConfig{
		ID:    "name",
		Label: "Name",
		Value: "ab",
		Width: 180,
		OnChange: func(value string) {
			changes = append(changes, value)
		},
		OnSubmit: func(value string) {
			submits++
			changes = append(changes, "submit:"+value)
		},
	}))

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 220, Height: 80}
	layout := runtime.Update(dom, viewport, uidom.InputSnapshot{})
	field, ok := layout.FindByID("name")
	if !ok {
		t.Fatalf("expected input field")
	}

	x := field.Frame.X + 12
	y := field.Frame.Y + 12
	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: x, PointerY: y, PointerDown: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: x, PointerY: y})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Text: "c"})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Backspace: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Submit: true})

	if got, want := runtime.FocusedID(), "name"; got != want {
		t.Fatalf("focus mismatch: got %q want %q", got, want)
	}
	if len(changes) < 3 {
		t.Fatalf("expected multiple input change events, got %#v", changes)
	}
	if changes[0] != "abc" {
		t.Fatalf("expected appended text, got %#v", changes)
	}
	if changes[1] != "ab" {
		t.Fatalf("expected backspace result, got %#v", changes)
	}
	if changes[2] != "submit:ab" || submits != 1 {
		t.Fatalf("expected submit callback, got %#v submits=%d", changes, submits)
	}
}

func TestRuntimeInputFieldSupportsCursorNavigationAndShortcutDeletion(t *testing.T) {
	var changes []string
	dom := uidom.New(uidom.InputField(uidom.InputFieldConfig{
		ID:    "name",
		Label: "Name",
		Value: "abc def",
		Width: 180,
		OnChange: func(value string) {
			changes = append(changes, value)
		},
	}))

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 220, Height: 80}
	runtime.Update(dom, viewport, uidom.InputSnapshot{})
	runtime.FocusNodeByID(dom, "name", uidom.InputSnapshot{})

	runtime.Update(dom, viewport, uidom.InputSnapshot{Control: true, Backspace: true})
	if got, want := runtime.TextValueOrDefault("name", ""), "abc"; got != want {
		t.Fatalf("expected shortcut deletion, got %q want %q", got, want)
	}

	runtime.Update(dom, viewport, uidom.InputSnapshot{Home: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Text: "X"})
	if got, want := runtime.TextValueOrDefault("name", ""), "Xabc"; got != want {
		t.Fatalf("expected home insertion, got %q want %q", got, want)
	}

	runtime.Update(dom, viewport, uidom.InputSnapshot{End: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{ArrowLeft: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Delete: true})
	if got, want := runtime.TextValueOrDefault("name", ""), "Xab"; got != want {
		t.Fatalf("expected delete after cursor move, got %q want %q", got, want)
	}

	if len(changes) < 3 {
		t.Fatalf("expected multiple changes, got %#v", changes)
	}
}

func TestRuntimeTextareaUsesEnterForNewlineAndCtrlEnterForSubmit(t *testing.T) {
	var changes []string
	submits := 0
	dom := uidom.New(uidom.Textarea(uidom.TextareaConfig{
		ID:     "bio",
		Label:  "Bio",
		Value:  "hi",
		Width:  220,
		Height: 80,
		OnChange: func(value string) {
			changes = append(changes, value)
		},
		OnSubmit: func(value string) {
			submits++
			changes = append(changes, "submit:"+value)
		},
	}))

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 260, Height: 140}
	runtime.Update(dom, viewport, uidom.InputSnapshot{})
	runtime.FocusNodeByID(dom, "bio", uidom.InputSnapshot{})
	runtime.Update(dom, viewport, uidom.InputSnapshot{})

	runtime.Update(dom, viewport, uidom.InputSnapshot{Submit: true})
	if got, want := runtime.TextValueOrDefault("bio", ""), "hi\n"; got != want {
		t.Fatalf("expected newline on enter, got %q want %q", got, want)
	}

	runtime.Update(dom, viewport, uidom.InputSnapshot{Text: "there"})
	if got, want := runtime.TextValueOrDefault("bio", ""), "hi\nthere"; got != want {
		t.Fatalf("expected inserted text, got %q want %q", got, want)
	}

	runtime.Update(dom, viewport, uidom.InputSnapshot{Control: true, Submit: true})
	if submits != 1 {
		t.Fatalf("expected ctrl+enter submit once, got %d", submits)
	}
	if got := changes[len(changes)-1]; got != "submit:hi\nthere" {
		t.Fatalf("expected submit payload, got %q", got)
	}
}

func TestRuntimeDropdownDispatchesOpenAndSelect(t *testing.T) {
	openStates := []bool{}
	selected := []string{}
	dom := uidom.New(uidom.Dropdown(uidom.DropdownConfig{
		ID:           "resolution",
		Label:        "Resolution",
		SelectedText: "1280x720",
		Open:         true,
		Width:        220,
		Options: []uidom.DropdownOption{
			{ID: "res-720", Label: "1280x720"},
			{ID: "res-1080", Label: "1920x1080"},
		},
		OnOpenChange: func(value bool) {
			openStates = append(openStates, value)
		},
		OnSelect: func(id string) {
			selected = append(selected, id)
		},
	}))

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 260, Height: 180}
	layout := runtime.Update(dom, viewport, uidom.InputSnapshot{})

	trigger, ok := layout.FindByID("resolution-trigger")
	if !ok {
		t.Fatalf("expected dropdown trigger")
	}
	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: trigger.Frame.X + 10, PointerY: trigger.Frame.Y + 10, PointerDown: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: trigger.Frame.X + 10, PointerY: trigger.Frame.Y + 10})

	option, ok := layout.FindByID("res-1080")
	if !ok {
		t.Fatalf("expected dropdown option")
	}
	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: option.Frame.X + 10, PointerY: option.Frame.Y + 10, PointerDown: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: option.Frame.X + 10, PointerY: option.Frame.Y + 10})

	if len(openStates) == 0 || openStates[0] != false {
		t.Fatalf("expected trigger to request close, got %#v", openStates)
	}
	if len(selected) == 0 || selected[0] != "res-1080" {
		t.Fatalf("expected selected option, got %#v", selected)
	}
}

func TestRuntimeTabCyclesFocusableNodesAndEscapeClearsFocus(t *testing.T) {
	dom := uidom.New(
		uidom.Div(uidom.Props{
			ID: "root",
			Style: uidom.Style{
				Width:     uidom.Px(320),
				Height:    uidom.Px(160),
				Direction: uidom.Column,
				Gap:       8,
			},
		},
			uidom.InputField(uidom.InputFieldConfig{
				ID:    "name",
				Label: "Name",
				Value: "ab",
				Width: 180,
			}),
			uidom.InteractiveButton(uidom.Props{
				ID: "submit",
				Style: uidom.Style{
					Width:  uidom.Px(120),
					Height: uidom.Px(40),
				},
			}, uidom.Text("Submit", uidom.Props{ID: "submit-label"})),
		),
	)

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 320, Height: 160}

	runtime.Update(dom, viewport, uidom.InputSnapshot{Tab: true})
	if got, want := runtime.FocusedID(), "name"; got != want {
		t.Fatalf("focus mismatch after first tab: got %q want %q", got, want)
	}
	runtime.Update(dom, viewport, uidom.InputSnapshot{Tab: true})
	if got, want := runtime.FocusedID(), "submit"; got != want {
		t.Fatalf("focus mismatch after second tab: got %q want %q", got, want)
	}
	runtime.Update(dom, viewport, uidom.InputSnapshot{Escape: true})
	if got := runtime.FocusedID(); got != "" {
		t.Fatalf("expected escape to clear focus, got %q", got)
	}
}

func TestRuntimeSubmitActivatesFocusedButton(t *testing.T) {
	clicks := 0
	dom := uidom.New(
		uidom.Div(uidom.Props{
			ID: "root",
			Style: uidom.Style{
				Width:  uidom.Px(200),
				Height: uidom.Px(80),
			},
		},
			uidom.InteractiveButton(uidom.Props{
				ID: "play-button",
				Handlers: uidom.EventHandlers{
					OnClick: func(ctx uidom.EventContext) {
						clicks++
					},
				},
				Style: uidom.Style{
					Width:  uidom.Px(120),
					Height: uidom.Px(40),
				},
			}, uidom.Text("Play", uidom.Props{ID: "play-label"})),
		),
	)

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 200, Height: 80}
	runtime.FocusNodeByID(dom, "play-button", uidom.InputSnapshot{})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Submit: true})

	if clicks != 1 {
		t.Fatalf("expected keyboard submit to trigger click, got %d", clicks)
	}
}

func TestRuntimeArrowKeysScrollFocusedScrollView(t *testing.T) {
	var offsets []float64
	dom := uidom.New(uidom.ScrollView(uidom.Props{
		ID:        "scroll",
		Focusable: true,
		Layout: uidom.LayoutSpec{
			Mode: uidom.LayoutModeFlowVertical,
			Size: uidom.LayoutSize{
				Width:  uidom.Px(220),
				Height: uidom.Px(120),
			},
		},
		Handlers: uidom.EventHandlers{
			OnScroll: func(ctx uidom.EventContext) {
				offsets = append(offsets, ctx.ScrollY)
			},
		},
	},
		uidom.Div(uidom.Props{ID: "row-a", Style: uidom.Style{Height: uidom.Px(120)}}),
		uidom.Div(uidom.Props{ID: "row-b", Style: uidom.Style{Height: uidom.Px(120)}}),
	))

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 220, Height: 120}
	runtime.FocusNodeByID(dom, "scroll", uidom.InputSnapshot{})
	runtime.Update(dom, viewport, uidom.InputSnapshot{ArrowDown: true})

	if len(offsets) == 0 {
		t.Fatalf("expected arrow key to dispatch scroll")
	}
	if offsets[0] <= 0 {
		t.Fatalf("expected positive scroll delta, got %#v", offsets)
	}
}

func TestRuntimeInputFieldSupportsCursorDeleteHomeEndAndReplacement(t *testing.T) {
	var changes []string
	dom := uidom.New(uidom.InputField(uidom.InputFieldConfig{
		ID:    "name",
		Label: "Name",
		Value: "abcd",
		Width: 180,
		OnChange: func(value string) {
			changes = append(changes, value)
		},
	}))

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 220, Height: 80}
	runtime.FocusNodeByID(dom, "name", uidom.InputSnapshot{})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Home: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Text: "X"})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Delete: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{End: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Backspace: true})

	if len(changes) < 3 {
		t.Fatalf("expected editing changes, got %#v", changes)
	}
	if changes[0] != "Xabcd" {
		t.Fatalf("expected insert at home, got %#v", changes)
	}
	if changes[1] != "Xbcd" {
		t.Fatalf("expected delete at cursor, got %#v", changes)
	}
	if changes[2] != "Xbc" {
		t.Fatalf("expected backspace at end, got %#v", changes)
	}
}

func TestRuntimeInputFieldSupportsSelectAllShortcutAndWordDelete(t *testing.T) {
	var changes []string
	dom := uidom.New(uidom.InputField(uidom.InputFieldConfig{
		ID:    "query",
		Label: "Query",
		Value: "alpha beta",
		Width: 220,
		OnChange: func(value string) {
			changes = append(changes, value)
		},
	}))

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 260, Height: 80}
	runtime.Update(dom, viewport, uidom.InputSnapshot{})
	runtime.FocusNodeByID(dom, "query", uidom.InputSnapshot{})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Shortcut: "ctrl+a"})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Text: "Zeta beta"})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Control: true, Backspace: true})

	if len(changes) < 2 {
		t.Fatalf("expected multiple shortcut-driven changes, got %#v", changes)
	}
	if changes[0] != "Zeta beta" {
		t.Fatalf("expected select-all replacement, got %#v", changes)
	}
	if changes[1] != "Zeta" {
		t.Fatalf("expected appended text after replacement, got %#v", changes)
	}
}

func TestRuntimeDispatchesShortcutToFocusedNode(t *testing.T) {
	var shortcut string
	dom := uidom.New(uidom.InteractiveButton(uidom.Props{
		ID: "command",
		Handlers: uidom.EventHandlers{
			OnShortcut: func(ctx uidom.EventContext) {
				shortcut = ctx.Shortcut
			},
		},
		Style: uidom.Style{
			Width:  uidom.Px(120),
			Height: uidom.Px(40),
		},
	}, uidom.Text("Command", uidom.Props{ID: "command-label"})))

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 160, Height: 80}
	runtime.Update(dom, viewport, uidom.InputSnapshot{})
	runtime.FocusNodeByID(dom, "command", uidom.InputSnapshot{})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Shortcut: "ctrl+k"})

	if shortcut != "ctrl+k" {
		t.Fatalf("expected shortcut routing, got %q", shortcut)
	}
}

func TestRuntimeTextEditingSupportsCursorMovementAndDeletion(t *testing.T) {
	var changes []string
	dom := uidom.New(uidom.InputField(uidom.InputFieldConfig{
		ID:    "nickname",
		Label: "Nickname",
		Value: "ab",
		Width: 180,
		OnChange: func(value string) {
			changes = append(changes, value)
		},
	}))

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 220, Height: 80}
	runtime.FocusNodeByID(dom, "nickname", uidom.InputSnapshot{})

	runtime.Update(dom, viewport, uidom.InputSnapshot{ArrowLeft: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Text: "c"})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Home: true})
	runtime.Update(dom, viewport, uidom.InputSnapshot{Delete: true})

	if got, want := runtime.TextValueOrDefault("nickname", ""), "cb"; got != want {
		t.Fatalf("text edit mismatch: got %q want %q", got, want)
	}
	if len(changes) < 2 {
		t.Fatalf("expected text change callbacks, got %#v", changes)
	}
	if changes[0] != "acb" {
		t.Fatalf("expected cursor insertion result, got %#v", changes)
	}
	if changes[1] != "cb" {
		t.Fatalf("expected delete result, got %#v", changes)
	}
}

func TestRuntimeSelectAllReplacesFocusedText(t *testing.T) {
	dom := uidom.New(uidom.Textarea(uidom.TextareaConfig{
		ID:    "notes",
		Value: "hello",
		Width: 220,
	}))

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 240, Height: 100}
	runtime.Update(dom, viewport, uidom.InputSnapshot{})
	runtime.FocusNodeByID(dom, "notes", uidom.InputSnapshot{})
	runtime.Update(dom, viewport, uidom.InputSnapshot{})
	runtime.Update(dom, viewport, uidom.InputSnapshot{SelectAll: true, Text: "ok"})

	if got, want := runtime.TextValueOrDefault("notes", ""), "ok"; got != want {
		t.Fatalf("select-all replacement mismatch: got %q want %q", got, want)
	}
}
