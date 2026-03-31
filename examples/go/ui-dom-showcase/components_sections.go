package main

import (
	"image/color"

	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
)

func componentsSection() *uidom.Node {
	return uidom.Section(uidom.Props{
		ID:    "components-section",
		Style: showcaseGroupStyle(),
	},
		uidom.Text("Components", uidom.Props{
			ID:    "components-section-title",
			Style: showcaseGroupTitleStyle(),
		}),
		uidom.TextBlock("Input controls, layout helpers, overlays, data widgets, and state chips are all shown in one place.", uidom.Props{
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

func formSection() *uidom.Node {
	return uidom.Section(uidom.Props{
		ID:    "form-section",
		Style: detailSectionStyle(),
	},
		uidom.Text("Form Components", uidom.Props{
			ID:    "form-title",
			Style: detailTitleStyle(),
		}),
		uidom.InputField(uidom.InputFieldConfig{
			ID:    "name-input",
			Label: "Player Name",
			Value: "Kim",
			Width: 260,
			State: uidom.InteractionState{Focused: true},
		}),
		uidom.Dropdown(uidom.DropdownConfig{
			ID:           "resolution-dropdown",
			Label:        "Resolution",
			SelectedText: "1280x720",
			Width:        260,
			Open:         true,
			Options: []uidom.DropdownOption{
				{ID: "resolution-720", Label: "1280x720"},
				{ID: "resolution-1080", Label: "1920x1080", State: uidom.InteractionState{Focused: true}},
			},
		}),
		uidom.Textarea(uidom.TextareaConfig{
			ID:    "bio-textarea",
			Label: "Profile",
			Value: "Explorer of the ember valley.\nSpecializes in bows and traps.",
			Width: 320,
			State: uidom.InteractionState{Focused: true},
		}),
		uidom.RadioGroup(uidom.RadioGroupConfig{
			ID:          "mode-radio",
			Label:       "Input Mode",
			Orientation: uidom.Row,
			Options: []uidom.RadioOption{
				{ID: "mode-kbm", Label: "Keyboard/Mouse"},
				{ID: "mode-pad", Label: "Gamepad", Selected: true},
			},
		}),
		uidom.Stepper(uidom.StepperConfig{
			ID:    "party-stepper",
			Label: "Party Size",
			Value: 3,
			Min:   1,
			Max:   4,
			Width: 220,
		}),
	)
}

func layoutSection() *uidom.Node {
	return uidom.Section(uidom.Props{
		ID:    "layout-section",
		Style: detailSectionStyle(),
	},
		uidom.Text("Layout and Data", uidom.Props{
			ID:    "layout-title",
			Style: detailTitleStyle(),
		}),
		uidom.Grid(uidom.GridConfig{
			ID:      "content-grid",
			Columns: 3,
			Gap:     10,
			Children: []*uidom.Node{
				uidom.Icon(uidom.IconConfig{ID: "grid-icon-0", Size: 20, Image: uidom.SolidImage(20, 20, color.RGBA{R: 92, G: 162, B: 255, A: 255})}),
				uidom.Badge(uidom.BadgeConfig{ID: "grid-badge-1", Label: "Rare"}),
				uidom.Chip(uidom.ChipConfig{ID: "grid-chip-2", Label: "Fire", Dismissible: true}),
				uidom.Text("Cell 4", uidom.Props{ID: "grid-text-3", Style: uidom.Style{Color: color.RGBA{R: 239, G: 244, B: 250, A: 255}}}),
			},
		}),
		uidom.List(uidom.ListConfig{
			ID:          "stat-list",
			Orientation: uidom.Column,
			Gap:         6,
			Items: []*uidom.Node{
				uidom.Text("ATK +18", uidom.Props{ID: "stat-item-0", Style: uidom.Style{Color: color.RGBA{R: 239, G: 244, B: 250, A: 255}}}),
				uidom.Text("DEF +10", uidom.Props{ID: "stat-item-1", Style: uidom.Style{Color: color.RGBA{R: 176, G: 188, B: 204, A: 255}}}),
			},
		}),
		uidom.VirtualList(uidom.VirtualListConfig{
			ID:           "virtual-items",
			StartIndex:   2,
			VisibleCount: 4,
			TotalCount:   12,
			Orientation:  uidom.Column,
			Gap:          6,
			ItemBuilder: func(index int) *uidom.Node {
				return uidom.Text(
					uidom.ComponentID("Loot", index),
					uidom.Props{
						ID: uidom.ComponentID("virtual-item", index),
						Style: uidom.Style{
							Color: color.RGBA{R: 239, G: 244, B: 250, A: 255},
						},
					},
				)
			},
		}),
		uidom.Scrollbar(uidom.ScrollbarConfig{
			ID:           "inventory-scrollbar",
			Orientation:  uidom.Vertical,
			Length:       120,
			Thickness:    12,
			ViewportSize: 30,
			ContentSize:  90,
			ScrollOffset: 30,
		}),
		uidom.Divider(uidom.DividerConfig{
			ID:          "layout-divider",
			Orientation: uidom.Horizontal,
			Length:      320,
			Thickness:   2,
		}),
	)
}

func overlaySection() *uidom.Node {
	return uidom.Section(uidom.Props{
		ID:    "overlay-section",
		Style: detailSectionStyle(),
	},
		uidom.Text("Overlay Components", uidom.Props{
			ID:    "overlay-title",
			Style: detailTitleStyle(),
		}),
		uidom.Modal(uidom.ModalConfig{
			ID:     "settings-modal",
			Open:   true,
			Title:  "Settings",
			Width:  280,
			Height: 160,
			Content: uidom.TextBlock("Audio, video, and input settings live inside the modal container.", uidom.Props{
				ID: "settings-modal-copy",
				Style: uidom.Style{
					Width:      uidom.Fill(),
					Color:      color.RGBA{R: 176, G: 188, B: 204, A: 255},
					LineHeight: 16,
				},
			}),
		}),
		uidom.Tooltip(uidom.TooltipConfig{
			ID:          "loot-tooltip",
			Title:       "Crystal Bow",
			Description: "A precise ranged weapon with low draw delay and high crit chance.",
			Width:       260,
		}),
		uidom.ContextMenu(uidom.ContextMenuConfig{
			ID:    "slot-context-menu",
			Width: 220,
			Items: []uidom.ContextMenuItem{
				{ID: "slot-use", Label: "Use"},
				{ID: "slot-equip", Label: "Equip", State: uidom.InteractionState{Focused: true}},
				{ID: "slot-drop", Label: "Drop"},
			},
		}),
	)
}

func dataSection() *uidom.Node {
	return uidom.Section(uidom.Props{
		ID:    "data-section",
		Style: detailSectionStyle(),
	},
		uidom.Text("Tabs and Accordion", uidom.Props{
			ID:    "data-title",
			Style: detailTitleStyle(),
		}),
		uidom.Tabs(uidom.TabsConfig{
			ID:            "tabs-demo",
			SelectedIndex: 1,
			Tabs: []uidom.TabConfig{
				{
					ID:    "tab-stats",
					Label: "Stats",
					Content: uidom.Text("Stats panel", uidom.Props{
						ID:    "tab-stats-panel",
						Style: uidom.Style{Color: color.RGBA{R: 176, G: 188, B: 204, A: 255}},
					}),
				},
				{
					ID:    "tab-skills",
					Label: "Skills",
					Content: uidom.Text("Skills panel", uidom.Props{
						ID:    "tab-skills-panel",
						Style: uidom.Style{Color: color.RGBA{R: 239, G: 244, B: 250, A: 255}},
					}),
				},
			},
		}),
		uidom.Accordion(uidom.AccordionConfig{
			ID: "accordion-demo",
			Sections: []uidom.AccordionSection{
				{
					ID:       "accordion-a",
					Label:    "Attack Combo",
					Expanded: false,
					Content:  uidom.Text("Combo details", uidom.Props{ID: "accordion-a-body"}),
				},
				{
					ID:       "accordion-b",
					Label:    "Passive Skills",
					Expanded: true,
					Content:  uidom.Text("Passive skill details", uidom.Props{ID: "accordion-b-body"}),
				},
			},
		}),
	)
}

func statusSection() *uidom.Node {
	return uidom.Section(uidom.Props{
		ID:    "status-section",
		Style: detailSectionStyle(),
	},
		uidom.Text("Status and Toggles", uidom.Props{
			ID:    "status-title",
			Style: detailTitleStyle(),
		}),
		uidom.Toggle(uidom.ToggleConfig{
			ID:      "difficulty-toggle",
			Label:   "Hardcore Mode",
			Checked: true,
		}),
		uidom.Slider(uidom.SliderConfig{
			ID:    "music-slider",
			Label: "Music",
			Min:   0,
			Max:   100,
			Value: 72,
			Width: 260,
		}),
		uidom.ProgressBar(uidom.ProgressBarConfig{
			ID:      "exp-progress",
			Label:   "EXP",
			Current: 54,
			Max:     100,
			Width:   260,
			Tint:    color.RGBA{R: 82, G: 205, B: 150, A: 255},
		}),
		uidom.Div(uidom.Props{
			ID: "status-tags",
			Style: uidom.Style{
				Width:     uidom.Fill(),
				Direction: uidom.Row,
				Gap:       8,
			},
		},
			uidom.Badge(uidom.BadgeConfig{
				ID:    "elite-badge",
				Label: "ELITE",
			}),
			uidom.Chip(uidom.ChipConfig{
				ID:          "fire-chip",
				Label:       "Fire",
				Dismissible: true,
			}),
			uidom.Checkbox(uidom.CheckboxConfig{
				ID:      "autosave-check",
				Label:   "Autosave",
				Checked: true,
			}),
		),
	)
}
