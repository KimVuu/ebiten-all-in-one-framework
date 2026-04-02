package ebitenui_test

import (
	"testing"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func TestRuntimeDispatchesButtonClickAndFocus(t *testing.T) {
	clicks := 0
	dom := ebitenui.New(
		ebitenui.Div(ebitenui.Props{
			ID: "root",
			Style: ebitenui.Style{
				Width:   ebitenui.Px(200),
				Height:  ebitenui.Px(80),
				Padding: ebitenui.All(8),
			},
		},
			ebitenui.InteractiveButton(ebitenui.Props{
				ID: "play-button",
				Style: ebitenui.Style{
					Width:  ebitenui.Px(120),
					Height: ebitenui.Px(40),
				},
				Handlers: ebitenui.EventHandlers{
					OnClick: func(ctx ebitenui.EventContext) {
						clicks++
					},
				},
			},
				ebitenui.Text("Play", ebitenui.Props{ID: "play-button-label"}),
			),
		),
	)

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 200, Height: 80}

	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: 20, PointerY: 20})
	button, ok := dom.FindByID("play-button")
	if !ok || !button.Props.State.Hovered {
		t.Fatalf("expected button hover state")
	}

	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: 20, PointerY: 20, PointerDown: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: 20, PointerY: 20})

	if clicks != 1 {
		t.Fatalf("expected click handler once, got %d", clicks)
	}
	if got, want := runtime.FocusedID(), "play-button"; got != want {
		t.Fatalf("focus mismatch: got %q want %q", got, want)
	}
}

func TestRuntimeCheckboxOnChangeUsesRuntimeValue(t *testing.T) {
	var values []bool
	dom := ebitenui.New(ebitenui.Checkbox(ebitenui.CheckboxConfig{
		ID:      "remember",
		Label:   "Remember",
		Checked: false,
		OnChange: func(value bool) {
			values = append(values, value)
		},
	}))

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 220, Height: 48}

	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: 12, PointerY: 12, PointerDown: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: 12, PointerY: 12})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: 12, PointerY: 12, PointerDown: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: 12, PointerY: 12})

	if len(values) != 2 {
		t.Fatalf("expected two toggle values, got %d", len(values))
	}
	if values[0] != true || values[1] != false {
		t.Fatalf("unexpected toggle values: %#v", values)
	}
}

func TestRuntimeSliderOnChangeTracksPointerPosition(t *testing.T) {
	var changed []float64
	dom := ebitenui.New(ebitenui.Slider(ebitenui.SliderConfig{
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

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 220, Height: 80}
	layout := runtime.Update(dom, viewport, ebitenui.InputSnapshot{})
	track, ok := layout.FindByID("volume-track")
	if !ok {
		t.Fatalf("expected slider track")
	}

	targetX := track.Frame.X + track.Frame.Width*0.75
	targetY := track.Frame.Y + track.Frame.Height*0.5
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: targetX, PointerY: targetY, PointerDown: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: targetX, PointerY: targetY})

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
	dom := ebitenui.New(ebitenui.InputField(ebitenui.InputFieldConfig{
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

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 220, Height: 80}
	layout := runtime.Update(dom, viewport, ebitenui.InputSnapshot{})
	field, ok := layout.FindByID("name")
	if !ok {
		t.Fatalf("expected input field")
	}

	x := field.Frame.X + 12
	y := field.Frame.Y + 12
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: x, PointerY: y, PointerDown: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: x, PointerY: y})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Text: "c"})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Backspace: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Submit: true})

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
	dom := ebitenui.New(ebitenui.InputField(ebitenui.InputFieldConfig{
		ID:    "name",
		Label: "Name",
		Value: "abc def",
		Width: 180,
		OnChange: func(value string) {
			changes = append(changes, value)
		},
	}))

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 220, Height: 80}
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{})
	runtime.FocusNodeByID(dom, "name", ebitenui.InputSnapshot{})

	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Control: true, Backspace: true})
	if got, want := runtime.TextValueOrDefault("name", ""), "abc"; got != want {
		t.Fatalf("expected shortcut deletion, got %q want %q", got, want)
	}

	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Home: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Text: "X"})
	if got, want := runtime.TextValueOrDefault("name", ""), "Xabc"; got != want {
		t.Fatalf("expected home insertion, got %q want %q", got, want)
	}

	runtime.Update(dom, viewport, ebitenui.InputSnapshot{End: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{ArrowLeft: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Delete: true})
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
	dom := ebitenui.New(ebitenui.Textarea(ebitenui.TextareaConfig{
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

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 260, Height: 140}
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{})
	runtime.FocusNodeByID(dom, "bio", ebitenui.InputSnapshot{})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{})

	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Submit: true})
	if got, want := runtime.TextValueOrDefault("bio", ""), "hi\n"; got != want {
		t.Fatalf("expected newline on enter, got %q want %q", got, want)
	}

	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Text: "there"})
	if got, want := runtime.TextValueOrDefault("bio", ""), "hi\nthere"; got != want {
		t.Fatalf("expected inserted text, got %q want %q", got, want)
	}

	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Control: true, Submit: true})
	if submits != 1 {
		t.Fatalf("expected ctrl+enter submit once, got %d", submits)
	}
	if got := changes[len(changes)-1]; got != "submit:hi\nthere" {
		t.Fatalf("expected submit payload, got %q", got)
	}
}

func TestRuntimeInputFieldUsesValueBinding(t *testing.T) {
	name := ebitenui.NewRef("Kim")
	dom := ebitenui.New(ebitenui.InputField(ebitenui.InputFieldConfig{
		ID:           "bound-name",
		Label:        "Name",
		ValueBinding: name,
		Width:        180,
	}))

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 220, Height: 80}
	layout := runtime.Update(dom, viewport, ebitenui.InputSnapshot{})
	field, ok := layout.FindByID("bound-name")
	if !ok {
		t.Fatalf("expected input field")
	}

	x := field.Frame.X + 12
	y := field.Frame.Y + 12
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: x, PointerY: y, PointerDown: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: x, PointerY: y})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Text: "A"})

	if got, want := name.Get(), "KimA"; got != want {
		t.Fatalf("binding mismatch: got %q want %q", got, want)
	}
}

func TestRuntimeToggleUsesCheckedBinding(t *testing.T) {
	enabled := ebitenui.NewRef(false)
	dom := ebitenui.New(ebitenui.Toggle(ebitenui.ToggleConfig{
		ID:             "music-toggle",
		Label:          "Music",
		CheckedBinding: enabled,
	}))

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 220, Height: 60}

	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: 12, PointerY: 12, PointerDown: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: 12, PointerY: 12})

	if !enabled.Get() {
		t.Fatalf("expected toggle binding to update")
	}
}

func TestRuntimeSliderUsesValueBinding(t *testing.T) {
	volume := ebitenui.NewRef(20.0)
	dom := ebitenui.New(ebitenui.Slider(ebitenui.SliderConfig{
		ID:           "bound-volume",
		Label:        "Volume",
		Min:          0,
		Max:          100,
		ValueBinding: volume,
		Width:        200,
	}))

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 220, Height: 80}
	layout := runtime.Update(dom, viewport, ebitenui.InputSnapshot{})
	track, ok := layout.FindByID("bound-volume-track")
	if !ok {
		t.Fatalf("expected slider track")
	}

	targetX := track.Frame.X + track.Frame.Width*0.75
	targetY := track.Frame.Y + track.Frame.Height*0.5
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: targetX, PointerY: targetY, PointerDown: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: targetX, PointerY: targetY})

	if got := volume.Get(); got < 70 || got > 80 {
		t.Fatalf("expected bound slider value near 75, got %v", got)
	}
}

func TestRuntimeDropdownUsesBindings(t *testing.T) {
	selected := ebitenui.NewRef("resolution-720")
	open := ebitenui.NewRef(false)
	dom := ebitenui.New(ebitenui.Dropdown(ebitenui.DropdownConfig{
		ID:              "resolution-dropdown",
		Label:           "Resolution",
		SelectedBinding: selected,
		OpenBinding:     open,
		Width:           240,
		Options: []ebitenui.DropdownOption{
			{ID: "resolution-720", Label: "1280x720"},
			{ID: "resolution-1080", Label: "1920x1080"},
		},
	}))

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 280, Height: 120}

	layout := runtime.Update(dom, viewport, ebitenui.InputSnapshot{})
	trigger, ok := layout.FindByID("resolution-dropdown-trigger")
	if !ok {
		t.Fatalf("expected dropdown trigger")
	}

	x := trigger.Frame.X + 12
	y := trigger.Frame.Y + 12
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: x, PointerY: y, PointerDown: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: x, PointerY: y})
	if !open.Get() {
		t.Fatalf("expected dropdown open binding to update")
	}

	openDOM := ebitenui.New(ebitenui.Dropdown(ebitenui.DropdownConfig{
		ID:              "resolution-dropdown",
		Label:           "Resolution",
		SelectedBinding: selected,
		OpenBinding:     open,
		Width:           240,
		Options: []ebitenui.DropdownOption{
			{ID: "resolution-720", Label: "1280x720"},
			{ID: "resolution-1080", Label: "1920x1080"},
		},
	}))
	openLayout := runtime.Update(openDOM, viewport, ebitenui.InputSnapshot{})
	option, ok := openLayout.FindByID("resolution-1080")
	if !ok {
		t.Fatalf("expected dropdown option")
	}
	optionX := option.Frame.X + 12
	optionY := option.Frame.Y + 12
	runtime.Update(openDOM, viewport, ebitenui.InputSnapshot{PointerX: optionX, PointerY: optionY, PointerDown: true})
	runtime.Update(openDOM, viewport, ebitenui.InputSnapshot{PointerX: optionX, PointerY: optionY})

	if got, want := selected.Get(), "resolution-1080"; got != want {
		t.Fatalf("selected binding mismatch: got %q want %q", got, want)
	}
}

func TestRuntimeDropdownDispatchesOpenAndSelect(t *testing.T) {
	openStates := []bool{}
	selected := []string{}
	dom := ebitenui.New(ebitenui.Dropdown(ebitenui.DropdownConfig{
		ID:           "resolution",
		Label:        "Resolution",
		SelectedText: "1280x720",
		Open:         true,
		Width:        220,
		Options: []ebitenui.DropdownOption{
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

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 260, Height: 180}
	layout := runtime.Update(dom, viewport, ebitenui.InputSnapshot{})

	trigger, ok := layout.FindByID("resolution-trigger")
	if !ok {
		t.Fatalf("expected dropdown trigger")
	}
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: trigger.Frame.X + 10, PointerY: trigger.Frame.Y + 10, PointerDown: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: trigger.Frame.X + 10, PointerY: trigger.Frame.Y + 10})

	option, ok := layout.FindByID("res-1080")
	if !ok {
		t.Fatalf("expected dropdown option")
	}
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: option.Frame.X + 10, PointerY: option.Frame.Y + 10, PointerDown: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: option.Frame.X + 10, PointerY: option.Frame.Y + 10})

	if len(openStates) == 0 || openStates[0] != false {
		t.Fatalf("expected trigger to request close, got %#v", openStates)
	}
	if len(selected) == 0 || selected[0] != "res-1080" {
		t.Fatalf("expected selected option, got %#v", selected)
	}
}

func TestRuntimeTabCyclesFocusableNodesAndEscapeClearsFocus(t *testing.T) {
	dom := ebitenui.New(
		ebitenui.Div(ebitenui.Props{
			ID: "root",
			Style: ebitenui.Style{
				Width:     ebitenui.Px(320),
				Height:    ebitenui.Px(160),
				Direction: ebitenui.Column,
				Gap:       8,
			},
		},
			ebitenui.InputField(ebitenui.InputFieldConfig{
				ID:    "name",
				Label: "Name",
				Value: "ab",
				Width: 180,
			}),
			ebitenui.InteractiveButton(ebitenui.Props{
				ID: "submit",
				Style: ebitenui.Style{
					Width:  ebitenui.Px(120),
					Height: ebitenui.Px(40),
				},
			}, ebitenui.Text("Submit", ebitenui.Props{ID: "submit-label"})),
		),
	)

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 320, Height: 160}

	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Tab: true})
	if got, want := runtime.FocusedID(), "name"; got != want {
		t.Fatalf("focus mismatch after first tab: got %q want %q", got, want)
	}
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Tab: true})
	if got, want := runtime.FocusedID(), "submit"; got != want {
		t.Fatalf("focus mismatch after second tab: got %q want %q", got, want)
	}
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Escape: true})
	if got := runtime.FocusedID(); got != "" {
		t.Fatalf("expected escape to clear focus, got %q", got)
	}
}

func TestRuntimeSubmitActivatesFocusedButton(t *testing.T) {
	clicks := 0
	dom := ebitenui.New(
		ebitenui.Div(ebitenui.Props{
			ID: "root",
			Style: ebitenui.Style{
				Width:  ebitenui.Px(200),
				Height: ebitenui.Px(80),
			},
		},
			ebitenui.InteractiveButton(ebitenui.Props{
				ID: "play-button",
				Handlers: ebitenui.EventHandlers{
					OnClick: func(ctx ebitenui.EventContext) {
						clicks++
					},
				},
				Style: ebitenui.Style{
					Width:  ebitenui.Px(120),
					Height: ebitenui.Px(40),
				},
			}, ebitenui.Text("Play", ebitenui.Props{ID: "play-label"})),
		),
	)

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 200, Height: 80}
	runtime.FocusNodeByID(dom, "play-button", ebitenui.InputSnapshot{})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Submit: true})

	if clicks != 1 {
		t.Fatalf("expected keyboard submit to trigger click, got %d", clicks)
	}
}

func TestRuntimeArrowKeysScrollFocusedScrollView(t *testing.T) {
	var offsets []float64
	dom := ebitenui.New(ebitenui.ScrollView(ebitenui.Props{
		ID:        "scroll",
		Focusable: true,
		Layout: ebitenui.LayoutSpec{
			Mode: ebitenui.LayoutModeFlowVertical,
			Size: ebitenui.LayoutSize{
				Width:  ebitenui.Px(220),
				Height: ebitenui.Px(120),
			},
		},
		Handlers: ebitenui.EventHandlers{
			OnScroll: func(ctx ebitenui.EventContext) {
				offsets = append(offsets, ctx.ScrollY)
			},
		},
	},
		ebitenui.Div(ebitenui.Props{ID: "row-a", Style: ebitenui.Style{Height: ebitenui.Px(120)}}),
		ebitenui.Div(ebitenui.Props{ID: "row-b", Style: ebitenui.Style{Height: ebitenui.Px(120)}}),
	))

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 220, Height: 120}
	runtime.FocusNodeByID(dom, "scroll", ebitenui.InputSnapshot{})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{ArrowDown: true})

	if len(offsets) == 0 {
		t.Fatalf("expected arrow key to dispatch scroll")
	}
	if offsets[0] <= 0 {
		t.Fatalf("expected positive scroll delta, got %#v", offsets)
	}
}

func TestRuntimeInputFieldSupportsCursorDeleteHomeEndAndReplacement(t *testing.T) {
	var changes []string
	dom := ebitenui.New(ebitenui.InputField(ebitenui.InputFieldConfig{
		ID:    "name",
		Label: "Name",
		Value: "abcd",
		Width: 180,
		OnChange: func(value string) {
			changes = append(changes, value)
		},
	}))

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 220, Height: 80}
	runtime.FocusNodeByID(dom, "name", ebitenui.InputSnapshot{})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Home: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Text: "X"})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Delete: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{End: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Backspace: true})

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
	dom := ebitenui.New(ebitenui.InputField(ebitenui.InputFieldConfig{
		ID:    "query",
		Label: "Query",
		Value: "alpha beta",
		Width: 220,
		OnChange: func(value string) {
			changes = append(changes, value)
		},
	}))

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 260, Height: 80}
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{})
	runtime.FocusNodeByID(dom, "query", ebitenui.InputSnapshot{})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Shortcut: "ctrl+a"})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Text: "Zeta beta"})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Control: true, Backspace: true})

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
	dom := ebitenui.New(ebitenui.InteractiveButton(ebitenui.Props{
		ID: "command",
		Handlers: ebitenui.EventHandlers{
			OnShortcut: func(ctx ebitenui.EventContext) {
				shortcut = ctx.Shortcut
			},
		},
		Style: ebitenui.Style{
			Width:  ebitenui.Px(120),
			Height: ebitenui.Px(40),
		},
	}, ebitenui.Text("Command", ebitenui.Props{ID: "command-label"})))

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 160, Height: 80}
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{})
	runtime.FocusNodeByID(dom, "command", ebitenui.InputSnapshot{})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Shortcut: "ctrl+k"})

	if shortcut != "ctrl+k" {
		t.Fatalf("expected shortcut routing, got %q", shortcut)
	}
}

func TestRuntimeTextEditingSupportsCursorMovementAndDeletion(t *testing.T) {
	var changes []string
	dom := ebitenui.New(ebitenui.InputField(ebitenui.InputFieldConfig{
		ID:    "nickname",
		Label: "Nickname",
		Value: "ab",
		Width: 180,
		OnChange: func(value string) {
			changes = append(changes, value)
		},
	}))

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 220, Height: 80}
	runtime.FocusNodeByID(dom, "nickname", ebitenui.InputSnapshot{})

	runtime.Update(dom, viewport, ebitenui.InputSnapshot{ArrowLeft: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Text: "c"})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Home: true})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{Delete: true})

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
	dom := ebitenui.New(ebitenui.Textarea(ebitenui.TextareaConfig{
		ID:    "notes",
		Value: "hello",
		Width: 220,
	}))

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 240, Height: 100}
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{})
	runtime.FocusNodeByID(dom, "notes", ebitenui.InputSnapshot{})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{})
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{SelectAll: true, Text: "ok"})

	if got, want := runtime.TextValueOrDefault("notes", ""), "ok"; got != want {
		t.Fatalf("select-all replacement mismatch: got %q want %q", got, want)
	}
}
