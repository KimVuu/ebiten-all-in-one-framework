package main

import (
	"fmt"
	"image/color"
	"strings"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
	"github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui/prefabs"
)

type showcaseLayoutState struct {
	CurrentPageID string
	SidebarScroll float64
	DetailScroll  float64
}

type showcaseCallbacks struct {
	OnNavigate            func(string)
	OnSidebarScrollChange func(float64)
	OnDetailScrollChange  func(float64)
}

type ShowcasePageSpec struct {
	ID          string
	Title       string
	Group       string
	Description string
	UsageNotes  string
	CodeExample string
	DemoBuilder func(ShowcaseDemoContext) *ebitenui.Node
}

type ShowcasePageRegistry struct {
	Routes []ebitenui.PageRoute
	Pages  map[string]ShowcasePageSpec
}

type ShowcaseDemoContext struct {
	Runtime       *ebitenui.Runtime
	Registry      ShowcasePageRegistry
	CurrentPageID string
}

func buildShowcasePageRegistry() ShowcasePageRegistry {
	pages := map[string]ShowcasePageSpec{}
	add := func(spec ShowcasePageSpec) {
		pages[spec.ID] = spec
	}

	add(ShowcasePageSpec{
		ID:          "overview",
		Title:       "Overview",
		Description: "Page-based showcase for ebiten-ui. Use the left navigation to switch between foundations, inputs, layout helpers, overlays, data widgets, status controls, and prefab UI.",
		UsageNotes:  "Start with a group page for broad context, then move into a leaf page when you want the focused demo and code sample.",
		CodeExample: "router := ebitenui.NewPageRouter(ebitenui.PageRouterConfig{ /* routes */ })\ncurrent := router.CurrentPageID()\nroot := ebitenui.PageScreen(ebitenui.PageScreenConfig{Sidebar: nav, Content: detail})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
			textMuted := color.RGBA{R: 178, G: 190, B: 207, A: 255}
			return ebitenui.Div(ebitenui.Props{
				ID:    "overview-demo",
				Style: detailSectionStyle(),
			},
				ebitenui.Text("Learning Paths", ebitenui.Props{ID: "overview-demo-title", Style: detailTitleStyle()}),
				ebitenui.Div(ebitenui.Props{
					ID:    "overview-demo-cards",
					Style: ebitenui.Style{Width: ebitenui.Fill(), Direction: ebitenui.Row, Gap: 12},
				},
					infoCard("overview-card-pages", color.RGBA{R: 80, G: 160, B: 255, A: 255}, textStrong, textMuted, "Pages", "Every major UI concept now has a dedicated page."),
					infoCard("overview-card-code", color.RGBA{R: 72, G: 211, B: 161, A: 255}, textStrong, textMuted, "Code", "Each page shows canonical usage notes and code."),
					infoCard("overview-card-debug", color.RGBA{R: 255, G: 180, B: 72, A: 255}, textStrong, textMuted, "Debug", "MCP inspect/capture follows the current page state."),
				),
			)
		},
	})

	groupOverview := func(groupID, title, description string) ShowcasePageSpec {
		return ShowcasePageSpec{
			ID:          groupID,
			Title:       title,
			Description: description,
			UsageNotes:  "Use this group page to understand the category, then select a leaf page on the left for a focused demo and code sample.",
			CodeExample: fmt.Sprintf("router.Navigate(%q)\nchildren := router.ChildrenOf(%q)\nbreadcrumb := router.Breadcrumb(%q)", groupID, groupID, groupID),
			DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
				textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
				textMuted := color.RGBA{R: 178, G: 190, B: 207, A: 255}
				children := ctx.Registry.Pages[groupID]
				_ = children
				groupCards := make([]*ebitenui.Node, 0)
				for _, route := range ctx.Registry.Routes {
					if route.ID != groupID {
						continue
					}
					for _, child := range route.Children {
						childSpec := ctx.Registry.Pages[child.ID]
						groupCards = append(groupCards, infoCard("group-card-"+sanitizeID(child.ID), color.RGBA{R: 80, G: 160, B: 255, A: 255}, textStrong, textMuted, childSpec.Title, childSpec.Description))
					}
				}
				if len(groupCards) == 0 {
					groupCards = append(groupCards, infoCard("group-card-empty", color.RGBA{R: 80, G: 160, B: 255, A: 255}, textStrong, textMuted, title, description))
				}
				return ebitenui.Div(ebitenui.Props{
					ID:    "group-demo-" + sanitizeID(groupID),
					Style: detailSectionStyle(),
				},
					ebitenui.Text(title+" Pages", ebitenui.Props{ID: "group-demo-title-" + sanitizeID(groupID), Style: detailTitleStyle()}),
					ebitenui.Div(ebitenui.Props{
						ID:    "group-demo-cards-" + sanitizeID(groupID),
						Style: ebitenui.Style{Width: ebitenui.Fill(), Direction: ebitenui.Column, Gap: 12},
					}, groupCards...),
				)
			},
		}
	}

	add(groupOverview("foundations", "Foundations", "Low-level primitives that everything else in ebiten-ui builds on."))
	add(groupOverview("tags", "Tags", "HTML-like tag primitives and the basic DOM vocabulary."))
	add(groupOverview("inputs", "Inputs", "Stateful, focus-aware interactive form controls."))
	add(groupOverview("layout", "Layout", "Containers and helpers for arranging UI structures."))
	add(groupOverview("overlay", "Overlay", "Floating or layered UI such as modals and tooltips."))
	add(groupOverview("data", "Data", "Page-like data organization widgets such as tabs."))
	add(groupOverview("status", "Status", "Progress, toggles, and state indicators for gameplay UI."))
	add(groupOverview("prefabs", "Prefabs", "Reusable game UI assemblies built on top of ebiten-ui primitives."))

	add(ShowcasePageSpec{
		ID:          "foundations/image",
		Title:       "Image",
		Group:       "foundations",
		Description: "Use `Image` when you need a visual block, icon surface, or sprite region inside the same DOM tree.",
		UsageNotes:  "Images behave like nodes. You can size them with `Px` or `Fill` and combine them with stacks, cards, or grids.",
		CodeExample: "ebitenui.Image(ebitenui.Props{\n  ID: \"hero-icon\",\n  Style: ebitenui.Style{Width: ebitenui.Px(72), Height: ebitenui.Px(72)},\n  Image: ebitenui.SolidImage(72, 72, color.RGBA{R: 80, G: 160, B: 255, A: 255}),\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return previewImageCard("page-image-demo", color.RGBA{R: 80, G: 160, B: 255, A: 255})
		},
	})
	add(ShowcasePageSpec{
		ID:          "foundations/text-block",
		Title:       "TextBlock",
		Group:       "foundations",
		Description: "Use `TextBlock` for multiline descriptions and code-like wrapped copy inside constrained layouts.",
		UsageNotes:  "Set width and line-height so the block wraps predictably in cards, tooltips, and settings panels.",
		CodeExample: "ebitenui.TextBlock(\"Wrapped text\", ebitenui.Props{\n  ID: \"copy\",\n  Style: ebitenui.Style{Width: ebitenui.Px(320), LineHeight: 16},\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return ebitenui.Div(ebitenui.Props{ID: "page-text-block-demo", Style: detailSectionStyle()},
				ebitenui.Text("Wrapped Copy", ebitenui.Props{ID: "page-text-block-title", Style: detailTitleStyle()}),
				ebitenui.TextBlock("TextBlock can explain a UI feature, keep documentation inside the tree, and stay stable under responsive widths.", ebitenui.Props{
					ID: "page-text-block-copy",
					Style: ebitenui.Style{
						Width:      ebitenui.Px(360),
						Color:      color.RGBA{R: 178, G: 190, B: 207, A: 255},
						LineHeight: 16,
					},
				}),
			)
		},
	})
	add(ShowcasePageSpec{
		ID:          "foundations/spacer",
		Title:       "Spacer",
		Group:       "foundations",
		Description: "Use `Spacer` to push siblings apart in rows or columns without introducing extra layout logic.",
		UsageNotes:  "A fill-width spacer is the simplest way to keep content left/right aligned inside a row.",
		CodeExample: "ebitenui.Div(ebitenui.Props{Style: ebitenui.Style{Direction: ebitenui.Row}},\n  ebitenui.Text(\"Left\", ebitenui.Props{}),\n  ebitenui.Spacer(ebitenui.Props{Style: ebitenui.Style{Width: ebitenui.Fill(), Height: ebitenui.Px(1)}}),\n  ebitenui.Text(\"Right\", ebitenui.Props{}),\n)",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return ebitenui.Div(ebitenui.Props{ID: "page-spacer-demo", Style: detailSectionStyle()},
				ebitenui.Text("Spacer Alignment", ebitenui.Props{ID: "page-spacer-title", Style: detailTitleStyle()}),
				ebitenui.Div(ebitenui.Props{
					ID:    "page-spacer-row",
					Style: ebitenui.Style{Width: ebitenui.Fill(), Direction: ebitenui.Row, Gap: 8},
				},
					ebitenui.Text("Inventory", ebitenui.Props{ID: "page-spacer-left", Style: ebitenui.Style{Color: color.RGBA{R: 242, G: 246, B: 252, A: 255}}}),
					ebitenui.Spacer(ebitenui.Props{ID: "page-spacer-fill", Style: ebitenui.Style{Width: ebitenui.Fill(), Height: ebitenui.Px(1)}}),
					ebitenui.Text("24 slots", ebitenui.Props{ID: "page-spacer-right", Style: ebitenui.Style{Color: color.RGBA{R: 178, G: 190, B: 207, A: 255}}}),
				),
			)
		},
	})
	add(ShowcasePageSpec{
		ID:          "foundations/stack",
		Title:       "Stack",
		Group:       "foundations",
		Description: "Use `Stack` for layered UI such as icon + badge, modal backdrops, and HUD overlays.",
		UsageNotes:  "Keep the base panel and the overlay badge in the same stack so relative layering stays local.",
		CodeExample: "ebitenui.Stack(ebitenui.Props{ID: \"stack-demo\"},\n  ebitenui.Div(...base panel...),\n  ebitenui.Div(ebitenui.Props{Layout: ebitenui.LayoutSpec{Mode: ebitenui.LayoutModeAnchored, Anchor: ebitenui.AnchorTopRight}}, ...badge...),\n)",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return stackPreview("page-stack-demo", color.RGBA{R: 72, G: 211, B: 161, A: 255})
		},
	})
	add(ShowcasePageSpec{
		ID:          "foundations/scroll-view",
		Title:       "ScrollView",
		Group:       "foundations",
		Description: "Use `ScrollView` when content height can exceed the available panel or page area.",
		UsageNotes:  "Drive `OffsetY` from external state and update it through `OnScrollChange` callbacks for deterministic behavior.",
		CodeExample: "ebitenui.ScrollView(ebitenui.Props{\n  ID: \"detail-scroll\",\n  Scroll: ebitenui.ScrollState{OffsetY: offset},\n  Handlers: ebitenui.EventHandlers{OnScroll: func(ctx ebitenui.EventContext) { /* set next offset */ }},\n}, children...)",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return scrollPreview("page-scroll-view-demo")
		},
	})
	add(ShowcasePageSpec{
		ID:          "tags/basic-tags",
		Title:       "DOM Tags",
		Group:       "tags",
		Description: "The HTML-like tag layer gives ebiten-ui a readable structure: `div`, `header`, `main`, `section`, `footer`, and inline primitives.",
		UsageNotes:  "Use tags for intent, then use style and layout specs for actual placement and constraints.",
		CodeExample: "ebitenui.Div(...,\n  ebitenui.Header(...),\n  ebitenui.Main(...,\n    ebitenui.Section(...),\n  ),\n  ebitenui.Footer(...),\n)",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return ebitenui.Div(ebitenui.Props{ID: "page-tags-demo", Style: detailSectionStyle()},
				ebitenui.Text("Supported Tags", ebitenui.Props{ID: "page-tags-title", Style: detailTitleStyle()}),
				ebitenui.Div(ebitenui.Props{ID: "page-tags-row", Style: ebitenui.Style{Width: ebitenui.Fill(), Direction: ebitenui.Row, Gap: 8}},
					tagChip("page-tag-div", "div"),
					tagChip("page-tag-header", "header"),
					tagChip("page-tag-main", "main"),
					tagChip("page-tag-section", "section"),
					tagChip("page-tag-footer", "footer"),
					tagChip("page-tag-button", "button"),
					tagChip("page-tag-text", "text"),
				),
			)
		},
	})

	addInputPages(add)
	addLayoutPages(add)
	addOverlayPages(add)
	addDataPages(add)
	addStatusPages(add)
	addPrefabPages(add)

	routes := []ebitenui.PageRoute{
		{ID: "overview", Title: "Overview"},
		groupRoute("foundations", "Foundations", "foundations/image", "foundations/text-block", "foundations/spacer", "foundations/stack", "foundations/scroll-view"),
		groupRoute("tags", "Tags", "tags/basic-tags"),
		groupRoute("inputs", "Inputs", "inputs/input-field", "inputs/dropdown", "inputs/textarea", "inputs/radio-group", "inputs/stepper"),
		groupRoute("layout", "Layout", "layout/grid"),
		groupRoute("overlay", "Overlay", "overlay/modal", "overlay/tooltip"),
		groupRoute("data", "Data", "data/tabs"),
		groupRoute("status", "Status", "status/toggle", "status/slider"),
		groupRoute("prefabs", "Prefabs", "prefabs/dialog", "prefabs/inventory-grid"),
	}

	return ShowcasePageRegistry{
		Routes: routes,
		Pages:  pages,
	}
}

func addInputPages(add func(ShowcasePageSpec)) {
	add(ShowcasePageSpec{
		ID:          "inputs/input-field",
		Title:       "InputField",
		Group:       "inputs",
		Description: "Single-line text input with focus, caret, and runtime-backed value storage.",
		UsageNotes:  "Keep the `ID` stable so runtime text values survive rerenders and page switches.",
		CodeExample: "ebitenui.InputField(ebitenui.InputFieldConfig{\n  ID: \"name-input\",\n  Label: \"Player Name\",\n  Value: \"Kim\",\n  Width: 260,\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return ebitenui.InputField(ebitenui.InputFieldConfig{ID: "name-input", Label: "Player Name", Value: "Kim", Width: 280, State: ebitenui.InteractionState{Focused: true}})
		},
	})
	add(ShowcasePageSpec{
		ID:          "inputs/dropdown",
		Title:       "Dropdown",
		Group:       "inputs",
		Description: "Selectable option list for compact choice selection.",
		UsageNotes:  "Use stable option IDs and read the selected text ID from runtime state when the page rerenders.",
		CodeExample: "ebitenui.Dropdown(ebitenui.DropdownConfig{\n  ID: \"resolution-dropdown\",\n  Label: \"Resolution\",\n  SelectedText: \"1280x720\",\n  Open: true,\n  Options: []ebitenui.DropdownOption{{ID: \"resolution-720\", Label: \"1280x720\"}},\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return ebitenui.Dropdown(ebitenui.DropdownConfig{
				ID: "resolution-dropdown", Label: "Resolution", SelectedText: "1280x720", Width: 280, Open: true,
				Options: []ebitenui.DropdownOption{
					{ID: "resolution-720", Label: "1280x720"},
					{ID: "resolution-1080", Label: "1920x1080", State: ebitenui.InteractionState{Focused: true}},
				},
			})
		},
	})
	add(ShowcasePageSpec{
		ID:          "inputs/textarea",
		Title:       "Textarea",
		Group:       "inputs",
		Description: "Multiline editable text surface for profile, note, and editor-like UI.",
		UsageNotes:  "Pair `Textarea` with a fixed width panel and line-height so long text stays predictable.",
		CodeExample: "ebitenui.Textarea(ebitenui.TextareaConfig{\n  ID: \"bio-textarea\",\n  Label: \"Profile\",\n  Value: \"Explorer of the ember valley.\",\n  Width: 320,\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return ebitenui.Textarea(ebitenui.TextareaConfig{
				ID: "bio-textarea", Label: "Profile", Value: "Explorer of the ember valley.\nSpecializes in bows and traps.", Width: 320, State: ebitenui.InteractionState{Focused: true},
			})
		},
	})
	add(ShowcasePageSpec{
		ID:          "inputs/radio-group",
		Title:       "RadioGroup",
		Group:       "inputs",
		Description: "Mutually exclusive option selection for mode and category pickers.",
		UsageNotes:  "Use row orientation for compact mode switching and column orientation for settings lists.",
		CodeExample: "ebitenui.RadioGroup(ebitenui.RadioGroupConfig{\n  ID: \"mode-radio\",\n  Orientation: ebitenui.Row,\n  Options: []ebitenui.RadioOption{{ID: \"mode-pad\", Label: \"Gamepad\", Selected: true}},\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return ebitenui.RadioGroup(ebitenui.RadioGroupConfig{
				ID: "mode-radio", Label: "Input Mode", Orientation: ebitenui.Row,
				Options: []ebitenui.RadioOption{
					{ID: "mode-kbm", Label: "Keyboard/Mouse"},
					{ID: "mode-pad", Label: "Gamepad", Selected: true},
				},
			})
		},
	})
	add(ShowcasePageSpec{
		ID:          "inputs/stepper",
		Title:       "Stepper",
		Group:       "inputs",
		Description: "Bounded numeric input with explicit increment and decrement controls.",
		UsageNotes:  "Use it for party size, stack count, quantity, and option ranges with a small footprint.",
		CodeExample: "ebitenui.Stepper(ebitenui.StepperConfig{\n  ID: \"party-stepper\",\n  Label: \"Party Size\",\n  Min: 1,\n  Max: 4,\n  Value: 3,\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return ebitenui.Stepper(ebitenui.StepperConfig{ID: "party-stepper", Label: "Party Size", Value: 3, Min: 1, Max: 4, Width: 240})
		},
	})
}

func addLayoutPages(add func(ShowcasePageSpec)) {
	add(ShowcasePageSpec{
		ID:          "layout/grid",
		Title:       "Grid",
		Group:       "layout",
		Description: "Structured multi-column layout for icon boards, inventory previews, and dashboard cells.",
		UsageNotes:  "Use `Grid` when row/column flow is not enough and you need stable cell placement with gaps.",
		CodeExample: "ebitenui.Grid(ebitenui.GridConfig{\n  ID: \"content-grid\",\n  Columns: 3,\n  Gap: 10,\n  Children: []*ebitenui.Node{ /* cells */ },\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return ebitenui.Grid(ebitenui.GridConfig{
				ID:      "content-grid",
				Columns: 3,
				Gap:     10,
				Children: []*ebitenui.Node{
					ebitenui.Icon(ebitenui.IconConfig{ID: "grid-icon-0", Size: 20, Image: ebitenui.SolidImage(20, 20, color.RGBA{R: 92, G: 162, B: 255, A: 255})}),
					ebitenui.Badge(ebitenui.BadgeConfig{ID: "grid-badge-1", Label: "Rare"}),
					ebitenui.Chip(ebitenui.ChipConfig{ID: "grid-chip-2", Label: "Fire", Dismissible: true}),
					ebitenui.Text("Cell 4", ebitenui.Props{ID: "grid-text-3", Style: ebitenui.Style{Color: color.RGBA{R: 239, G: 244, B: 250, A: 255}}}),
				},
			})
		},
	})
}

func addOverlayPages(add func(ShowcasePageSpec)) {
	add(ShowcasePageSpec{
		ID:          "overlay/modal",
		Title:       "Modal",
		Group:       "overlay",
		Description: "Layered panel for blocking tasks such as settings, confirmation, or onboarding.",
		UsageNotes:  "Keep modal content short and focused. Use stack or overlay layouts behind it when you need backdrop context.",
		CodeExample: "ebitenui.Modal(ebitenui.ModalConfig{\n  ID: \"settings-modal\",\n  Open: true,\n  Title: \"Settings\",\n  Width: 280,\n  Height: 160,\n  Content: body,\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return ebitenui.Modal(ebitenui.ModalConfig{
				ID: "settings-modal", Open: true, Title: "Settings", Width: 280, Height: 160,
				Content: ebitenui.TextBlock("Audio, video, and input settings live inside the modal container.", ebitenui.Props{
					ID:    "settings-modal-copy",
					Style: ebitenui.Style{Width: ebitenui.Fill(), Color: color.RGBA{R: 176, G: 188, B: 204, A: 255}, LineHeight: 16},
				}),
			})
		},
	})
	add(ShowcasePageSpec{
		ID:          "overlay/tooltip",
		Title:       "Tooltip",
		Group:       "overlay",
		Description: "Small contextual information surface for items, stats, and icon explanations.",
		UsageNotes:  "Keep the title short and stats scannable. Tooltips should be readable at a glance.",
		CodeExample: "ebitenui.Tooltip(ebitenui.TooltipConfig{\n  ID: \"loot-tooltip\",\n  Title: \"Crystal Bow\",\n  Description: \"Precise ranged weapon.\",\n  Width: 260,\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return ebitenui.Tooltip(ebitenui.TooltipConfig{
				ID: "loot-tooltip", Title: "Crystal Bow", Description: "A precise ranged weapon with low draw delay and high crit chance.", Width: 260,
			})
		},
	})
}

func addDataPages(add func(ShowcasePageSpec)) {
	add(ShowcasePageSpec{
		ID:          "data/tabs",
		Title:       "Tabs",
		Group:       "data",
		Description: "Page-like data partitioning inside a single panel for stats, skills, logs, and settings groups.",
		UsageNotes:  "Tabs are useful when the user should stay within the same screen context while switching sections.",
		CodeExample: "ebitenui.Tabs(ebitenui.TabsConfig{\n  ID: \"tabs-demo\",\n  SelectedIndex: 1,\n  Tabs: []ebitenui.TabConfig{{ID: \"tab-skills\", Label: \"Skills\", Content: panel}},\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return ebitenui.Tabs(ebitenui.TabsConfig{
				ID: "tabs-demo", SelectedIndex: 1,
				Tabs: []ebitenui.TabConfig{
					{ID: "tab-stats", Label: "Stats", Content: ebitenui.Text("Stats panel", ebitenui.Props{ID: "tab-stats-panel", Style: ebitenui.Style{Color: color.RGBA{R: 176, G: 188, B: 204, A: 255}}})},
					{ID: "tab-skills", Label: "Skills", Content: ebitenui.Text("Skills panel", ebitenui.Props{ID: "tab-skills-panel", Style: ebitenui.Style{Color: color.RGBA{R: 239, G: 244, B: 250, A: 255}}})},
				},
			})
		},
	})
}

func addStatusPages(add func(ShowcasePageSpec)) {
	add(ShowcasePageSpec{
		ID:          "status/toggle",
		Title:       "Toggle",
		Group:       "status",
		Description: "Binary state switch for feature flags, gameplay options, and settings panels.",
		UsageNotes:  "Use `Toggle` when the choice is on/off and the user should see the current state at a glance.",
		CodeExample: "ebitenui.Toggle(ebitenui.ToggleConfig{\n  ID: \"difficulty-toggle\",\n  Label: \"Hardcore Mode\",\n  Checked: true,\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return ebitenui.Toggle(ebitenui.ToggleConfig{ID: "difficulty-toggle", Label: "Hardcore Mode", Checked: true})
		},
	})
	add(ShowcasePageSpec{
		ID:          "status/slider",
		Title:       "Slider",
		Group:       "status",
		Description: "Analog range control for volume, sensitivity, brightness, and progression knobs.",
		UsageNotes:  "Pair slider labels with immediate feedback values when the exact amount matters.",
		CodeExample: "ebitenui.Slider(ebitenui.SliderConfig{\n  ID: \"music-slider\",\n  Label: \"Music\",\n  Min: 0,\n  Max: 100,\n  Value: 65,\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return ebitenui.Slider(ebitenui.SliderConfig{ID: "music-slider", Label: "Music", Min: 0, Max: 100, Value: 65, Width: 300})
		},
	})
}

func addPrefabPages(add func(ShowcasePageSpec)) {
	add(ShowcasePageSpec{
		ID:          "prefabs/dialog",
		Title:       "Dialog Prefab",
		Group:       "prefabs",
		Description: "Prebuilt confirmation and decision surface for quit, reset, and transactional prompts.",
		UsageNotes:  "Use dialog prefabs when you want a stable interaction pattern instead of building confirmation flows by hand each time.",
		CodeExample: "prefabs.Dialog(prefabs.DialogConfig{\n  ID: \"dialog-demo\",\n  Title: \"Return to title?\",\n  Body: \"You will lose unsaved progress.\",\n  Actions: []prefabs.DialogAction{{ID: \"dialog-confirm\", Label: \"Confirm\"}},\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return prefabs.Dialog(prefabs.DialogConfig{
				ID: "dialog-demo", Title: "Return to title?", Body: "You will lose unsaved progress from this dungeon run.", Width: 360,
				Actions: []prefabs.DialogAction{
					{ID: "dialog-cancel", Label: "Cancel"},
					{ID: "dialog-confirm", Label: "Confirm", State: ebitenui.InteractionState{Selected: true}},
				},
			})
		},
	})
	add(ShowcasePageSpec{
		ID:          "prefabs/inventory-grid",
		Title:       "InventoryGrid Prefab",
		Group:       "prefabs",
		Description: "Slot-based inventory assembly built on top of cards, labels, and selection states.",
		UsageNotes:  "Use this prefab for quick inventory screens, then layer tooltips and context menus around it.",
		CodeExample: "prefabs.InventoryGrid(prefabs.InventoryGridConfig{\n  ID: \"inventory-demo\",\n  Title: \"Inventory\",\n  Columns: 4,\n  CellSize: 64,\n  Slots: slots,\n})",
		DemoBuilder: func(ctx ShowcaseDemoContext) *ebitenui.Node {
			return prefabs.InventoryGrid(prefabs.InventoryGridConfig{
				ID: "inventory-demo", Title: "Inventory", Columns: 4, CellSize: 64,
				Slots: []prefabs.InventorySlot{
					{ID: "inv-slot-1", Label: "Potion", Quantity: 3, State: ebitenui.InteractionState{Selected: true}},
					{ID: "inv-slot-2", Label: "Ether", Quantity: 1},
					{ID: "inv-slot-3", Label: "Key", Quantity: 1},
					{ID: "inv-slot-4", Label: "Gem", Quantity: 2},
				},
			})
		},
	})
}

func groupRoute(id, title string, children ...string) ebitenui.PageRoute {
	routeChildren := make([]ebitenui.PageRoute, 0, len(children))
	for _, child := range children {
		routeChildren = append(routeChildren, ebitenui.PageRoute{ID: child, Title: pageTitleFromID(child)})
	}
	defaultChildID := ""
	if len(children) > 0 {
		defaultChildID = children[0]
	}
	return ebitenui.PageRoute{
		ID:             id,
		Title:          title,
		Children:       routeChildren,
		DefaultChildID: defaultChildID,
	}
}

func pageTitleFromID(id string) string {
	last := id
	if index := strings.LastIndex(id, "/"); index >= 0 {
		last = id[index+1:]
	}
	parts := strings.Split(last, "-")
	for i, part := range parts {
		if len(part) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}

func sanitizeID(value string) string {
	value = strings.ReplaceAll(value, "/", "-")
	value = strings.ReplaceAll(value, " ", "-")
	return value
}
