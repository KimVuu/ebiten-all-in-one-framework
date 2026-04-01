package ebitenui

import "image/color"

type PageRoute struct {
	ID             string
	Title          string
	ParentID       string
	Children       []PageRoute
	Hidden         bool
	DefaultChildID string
}

type PageRouterConfig struct {
	Routes        []PageRoute
	InitialPageID string
}

type PageRouterState struct {
	CurrentPageID string
}

type PageRouter struct {
	routes       []PageRoute
	routeIndex   map[string]*PageRoute
	childIndex   map[string][]PageRoute
	defaultChild map[string]string
	state        PageRouterState
}

func NewPageRouter(config PageRouterConfig) *PageRouter {
	router := &PageRouter{
		routeIndex:   map[string]*PageRoute{},
		childIndex:   map[string][]PageRoute{},
		defaultChild: map[string]string{},
	}
	router.routes = normalizeRoutes(config.Routes, "", router.routeIndex, router.childIndex, router.defaultChild)

	initial := config.InitialPageID
	if initial == "" && len(router.routes) > 0 {
		initial = router.routes[0].ID
	}
	router.state.CurrentPageID = router.resolvePageID(initial)
	return router
}

func (router *PageRouter) CurrentPageID() string {
	if router == nil {
		return ""
	}
	return router.state.CurrentPageID
}

func (router *PageRouter) CurrentPage() *PageRoute {
	if router == nil {
		return nil
	}
	return router.routeIndex[router.state.CurrentPageID]
}

func (router *PageRouter) Navigate(pageID string) bool {
	if router == nil {
		return false
	}
	resolved := router.resolvePageID(pageID)
	if resolved == "" {
		return false
	}
	router.state.CurrentPageID = resolved
	return true
}

func (router *PageRouter) ParentPageID(pageID string) string {
	if router == nil {
		return ""
	}
	route := router.routeIndex[router.resolvePageID(pageID)]
	if route == nil {
		return ""
	}
	return route.ParentID
}

func (router *PageRouter) ChildrenOf(pageID string) []PageRoute {
	if router == nil {
		return nil
	}
	children := router.childIndex[pageID]
	return cloneRoutes(children)
}

func (router *PageRouter) Breadcrumb(pageID string) []PageRoute {
	if router == nil {
		return nil
	}
	current := router.routeIndex[router.resolvePageID(pageID)]
	if current == nil {
		return nil
	}

	chain := make([]PageRoute, 0, 4)
	for current != nil {
		chain = append([]PageRoute{cloneRoute(*current, false)}, chain...)
		if current.ParentID == "" {
			break
		}
		current = router.routeIndex[current.ParentID]
	}
	return chain
}

func (router *PageRouter) VisibleNavTree() []PageRoute {
	if router == nil {
		return nil
	}
	return cloneVisibleRoutes(router.routes)
}

func (router *PageRouter) resolvePageID(pageID string) string {
	if router == nil || pageID == "" {
		return ""
	}
	route := router.routeIndex[pageID]
	if route == nil {
		return ""
	}
	if defaultChild := router.defaultChild[pageID]; defaultChild != "" {
		return router.resolvePageID(defaultChild)
	}
	return route.ID
}

type PageScreenConfig struct {
	ID              string
	Sidebar         *Node
	Content         *Node
	SidebarWidth    float64
	Gap             float64
	Padding         float64
	BackgroundColor color.Color
}

func PageScreen(cfg PageScreenConfig) *Node {
	rootID := cfg.ID
	if rootID == "" {
		rootID = "page-screen"
	}
	sidebarWidth := cfg.SidebarWidth
	if sidebarWidth <= 0 {
		sidebarWidth = 280
	}
	gap := cfg.Gap
	if gap <= 0 {
		gap = 16
	}

	sidebar := cfg.Sidebar
	if sidebar == nil {
		sidebar = Div(Props{ID: rootID + "-sidebar"})
	}
	sidebar.Props.Style.Width = Px(sidebarWidth)
	sidebar.Props.Style.Height = Fill()

	content := cfg.Content
	if content == nil {
		content = Div(Props{ID: rootID + "-content"})
	}
	content.Props.Style.Width = Fill()
	content.Props.Style.Height = Fill()

	return Div(Props{
		ID: rootID,
		Style: Style{
			Width:           Fill(),
			Height:          Fill(),
			Direction:       Row,
			Padding:         All(cfg.Padding),
			Gap:             gap,
			BackgroundColor: cfg.BackgroundColor,
		},
	}, sidebar, content)
}

func normalizeRoutes(routes []PageRoute, parentID string, routeIndex map[string]*PageRoute, childIndex map[string][]PageRoute, defaultChild map[string]string) []PageRoute {
	normalized := make([]PageRoute, 0, len(routes))
	for _, route := range routes {
		route.ParentID = parentID
		route.Children = normalizeRoutes(route.Children, route.ID, routeIndex, childIndex, defaultChild)
		normalized = append(normalized, route)

		copied := route
		routeIndex[copied.ID] = &copied
		childIndex[copied.ID] = cloneRoutes(copied.Children)
		if copied.DefaultChildID != "" {
			defaultChild[copied.ID] = copied.DefaultChildID
		}
	}
	return normalized
}

func cloneVisibleRoutes(routes []PageRoute) []PageRoute {
	cloned := make([]PageRoute, 0, len(routes))
	for _, route := range routes {
		if route.Hidden {
			continue
		}
		cloned = append(cloned, cloneRoute(route, true))
	}
	return cloned
}

func cloneRoutes(routes []PageRoute) []PageRoute {
	cloned := make([]PageRoute, 0, len(routes))
	for _, route := range routes {
		cloned = append(cloned, cloneRoute(route, false))
	}
	return cloned
}

func cloneRoute(route PageRoute, visibleOnly bool) PageRoute {
	cloned := route
	if visibleOnly {
		cloned.Children = cloneVisibleRoutes(route.Children)
	} else {
		cloned.Children = cloneRoutes(route.Children)
	}
	return cloned
}
