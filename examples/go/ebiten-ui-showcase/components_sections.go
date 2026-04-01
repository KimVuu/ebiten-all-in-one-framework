package main

import (
	"image/color"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func componentsSection() *ebitenui.Node {
	return ebitenui.Section(ebitenui.Props{
		ID:    "components-section",
		Style: showcaseGroupStyle(),
	},
		ebitenui.Text("Components", ebitenui.Props{
			ID:    "components-section-title",
			Style: showcaseGroupTitleStyle(),
		}),
		ebitenui.TextBlock("Input controls, layout helpers, overlays, data widgets, and state chips are all shown in one place.", ebitenui.Props{
			ID:    "components-section-copy",
			Style: showcaseGroupCopyStyle(),
		}),
		formSection(),
		layoutSection(),
		overlaySection(),
		dataSection(),
		statusSection(),
	)
}

func formSection() *ebitenui.Node {
	return ebitenui.Section(ebitenui.Props{
		ID:    "form-section",
		Style: detailSectionStyle(),
	},
		ebitenui.Text("Form Components", ebitenui.Props{
			ID:    "form-title",
			Style: detailTitleStyle(),
		}),
		ebitenui.InputField(ebitenui.InputFieldConfig{
			ID:    "name-input",
			Label: "Player Name",
			Value: "Kim",
			Width: 260,
			State: ebitenui.InteractionState{Focused: true},
		}),
		ebitenui.Dropdown(ebitenui.DropdownConfig{
			ID:           "resolution-dropdown",
			Label:        "Resolution",
			SelectedText: "1280x720",
			Width:        260,
			Open:         true,
			Options: []ebitenui.DropdownOption{
				{ID: "resolution-720", Label: "1280x720"},
				{ID: "resolution-1080", Label: "1920x1080", State: ebitenui.InteractionState{Focused: true}},
			},
		}),
		ebitenui.Textarea(ebitenui.TextareaConfig{
			ID:    "bio-textarea",
			Label: "Profile",
			Value: "Explorer of the ember valley.\nSpecializes in bows and traps.",
			Width: 320,
			State: ebitenui.InteractionState{Focused: true},
		}),
		ebitenui.RadioGroup(ebitenui.RadioGroupConfig{
			ID:          "mode-radio",
			Label:       "Input Mode",
			Orientation: ebitenui.Row,
			Options: []ebitenui.RadioOption{
				{ID: "mode-kbm", Label: "Keyboard/Mouse"},
				{ID: "mode-pad", Label: "Gamepad", Selected: true},
			},
		}),
		ebitenui.Stepper(ebitenui.StepperConfig{
			ID:    "party-stepper",
			Label: "Party Size",
			Value: 3,
			Min:   1,
			Max:   4,
			Width: 220,
		}),
	)
}

func layoutSection() *ebitenui.Node {
	return ebitenui.Section(ebitenui.Props{
		ID:    "layout-section",
		Style: detailSectionStyle(),
	},
		ebitenui.Text("Layout and Data", ebitenui.Props{
			ID:    "layout-title",
			Style: detailTitleStyle(),
		}),
		ebitenui.Grid(ebitenui.GridConfig{
			ID:      "content-grid",
			Columns: 3,
			Gap:     10,
			Children: []*ebitenui.Node{
				ebitenui.Icon(ebitenui.IconConfig{ID: "grid-icon-0", Size: 20, Image: ebitenui.SolidImage(20, 20, color.RGBA{R: 92, G: 162, B: 255, A: 255})}),
				ebitenui.Badge(ebitenui.BadgeConfig{ID: "grid-badge-1", Label: "Rare"}),
				ebitenui.Chip(ebitenui.ChipConfig{ID: "grid-chip-2", Label: "Fire", Dismissible: true}),
				ebitenui.Text("Cell 4", ebitenui.Props{ID: "grid-text-3", Style: ebitenui.Style{Color: color.RGBA{R: 239, G: 244, B: 250, A: 255}}}),
			},
		}),
		ebitenui.List(ebitenui.ListConfig{
			ID:          "stat-list",
			Orientation: ebitenui.Column,
			Gap:         6,
			Items: []*ebitenui.Node{
				ebitenui.Text("ATK +18", ebitenui.Props{ID: "stat-item-0", Style: ebitenui.Style{Color: color.RGBA{R: 239, G: 244, B: 250, A: 255}}}),
				ebitenui.Text("DEF +10", ebitenui.Props{ID: "stat-item-1", Style: ebitenui.Style{Color: color.RGBA{R: 176, G: 188, B: 204, A: 255}}}),
			},
		}),
		ebitenui.VirtualList(ebitenui.VirtualListConfig{
			ID:           "virtual-items",
			StartIndex:   2,
			VisibleCount: 4,
			TotalCount:   12,
			Orientation:  ebitenui.Column,
			Gap:          6,
			ItemBuilder: func(index int) *ebitenui.Node {
				return ebitenui.Text(
					ebitenui.ComponentID("Loot", index),
					ebitenui.Props{
						ID: ebitenui.ComponentID("virtual-item", index),
						Style: ebitenui.Style{
							Color: color.RGBA{R: 239, G: 244, B: 250, A: 255},
						},
					},
				)
			},
		}),
		ebitenui.Scrollbar(ebitenui.ScrollbarConfig{
			ID:           "inventory-scrollbar",
			Orientation:  ebitenui.Vertical,
			Length:       120,
			Thickness:    12,
			ViewportSize: 30,
			ContentSize:  90,
			ScrollOffset: 30,
		}),
		ebitenui.Divider(ebitenui.DividerConfig{
			ID:          "layout-divider",
			Orientation: ebitenui.Horizontal,
			Length:      320,
			Thickness:   2,
		}),
	)
}

func overlaySection() *ebitenui.Node {
	return ebitenui.Section(ebitenui.Props{
		ID:    "overlay-section",
		Style: detailSectionStyle(),
	},
		ebitenui.Text("Overlay Components", ebitenui.Props{
			ID:    "overlay-title",
			Style: detailTitleStyle(),
		}),
		ebitenui.Modal(ebitenui.ModalConfig{
			ID:     "settings-modal",
			Open:   true,
			Title:  "Settings",
			Width:  280,
			Height: 160,
			Content: ebitenui.TextBlock("Audio, video, and input settings live inside the modal container.", ebitenui.Props{
				ID: "settings-modal-copy",
				Style: ebitenui.Style{
					Width:      ebitenui.Fill(),
					Color:      color.RGBA{R: 176, G: 188, B: 204, A: 255},
					LineHeight: 16,
				},
			}),
		}),
		ebitenui.Tooltip(ebitenui.TooltipConfig{
			ID:          "loot-tooltip",
			Title:       "Crystal Bow",
			Description: "A precise ranged weapon with low draw delay and high crit chance.",
			Width:       260,
		}),
		ebitenui.ContextMenu(ebitenui.ContextMenuConfig{
			ID:    "slot-context-menu",
			Width: 220,
			Items: []ebitenui.ContextMenuItem{
				{ID: "slot-use", Label: "Use"},
				{ID: "slot-equip", Label: "Equip", State: ebitenui.InteractionState{Focused: true}},
				{ID: "slot-drop", Label: "Drop"},
			},
		}),
	)
}

func dataSection() *ebitenui.Node {
	return ebitenui.Section(ebitenui.Props{
		ID:    "data-section",
		Style: detailSectionStyle(),
	},
		ebitenui.Text("Tabs and Accordion", ebitenui.Props{
			ID:    "data-title",
			Style: detailTitleStyle(),
		}),
		ebitenui.Tabs(ebitenui.TabsConfig{
			ID:            "tabs-demo",
			SelectedIndex: 1,
			Tabs: []ebitenui.TabConfig{
				{
					ID:    "tab-stats",
					Label: "Stats",
					Content: ebitenui.Text("Stats panel", ebitenui.Props{
						ID:    "tab-stats-panel",
						Style: ebitenui.Style{Color: color.RGBA{R: 176, G: 188, B: 204, A: 255}},
					}),
				},
				{
					ID:    "tab-skills",
					Label: "Skills",
					Content: ebitenui.Text("Skills panel", ebitenui.Props{
						ID:    "tab-skills-panel",
						Style: ebitenui.Style{Color: color.RGBA{R: 239, G: 244, B: 250, A: 255}},
					}),
				},
			},
		}),
		ebitenui.Accordion(ebitenui.AccordionConfig{
			ID: "accordion-demo",
			Sections: []ebitenui.AccordionSection{
				{
					ID:       "accordion-a",
					Label:    "Attack Combo",
					Expanded: false,
					Content:  ebitenui.Text("Combo details", ebitenui.Props{ID: "accordion-a-body"}),
				},
				{
					ID:       "accordion-b",
					Label:    "Passive Skills",
					Expanded: true,
					Content:  ebitenui.Text("Passive skill details", ebitenui.Props{ID: "accordion-b-body"}),
				},
			},
		}),
	)
}

func statusSection() *ebitenui.Node {
	return ebitenui.Section(ebitenui.Props{
		ID:    "status-section",
		Style: detailSectionStyle(),
	},
		ebitenui.Text("Status and Toggles", ebitenui.Props{
			ID:    "status-title",
			Style: detailTitleStyle(),
		}),
		ebitenui.Toggle(ebitenui.ToggleConfig{
			ID:      "difficulty-toggle",
			Label:   "Hardcore Mode",
			Checked: true,
		}),
		ebitenui.Slider(ebitenui.SliderConfig{
			ID:    "music-slider",
			Label: "Music",
			Min:   0,
			Max:   100,
			Value: 72,
			Width: 260,
		}),
		ebitenui.ProgressBar(ebitenui.ProgressBarConfig{
			ID:      "exp-progress",
			Label:   "EXP",
			Current: 54,
			Max:     100,
			Width:   260,
			Tint:    color.RGBA{R: 82, G: 205, B: 150, A: 255},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "status-tags",
			Style: ebitenui.Style{
				Width:     ebitenui.Fill(),
				Direction: ebitenui.Row,
				Gap:       8,
			},
		},
			ebitenui.Badge(ebitenui.BadgeConfig{
				ID:    "elite-badge",
				Label: "ELITE",
			}),
			ebitenui.Chip(ebitenui.ChipConfig{
				ID:          "fire-chip",
				Label:       "Fire",
				Dismissible: true,
			}),
			ebitenui.Checkbox(ebitenui.CheckboxConfig{
				ID:      "autosave-check",
				Label:   "Autosave",
				Checked: true,
			}),
		),
	)
}
