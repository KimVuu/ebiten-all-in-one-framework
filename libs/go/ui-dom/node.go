package uidom

type Tag string

const (
	TagDiv        Tag = "div"
	TagHeader     Tag = "header"
	TagMain       Tag = "main"
	TagSection    Tag = "section"
	TagFooter     Tag = "footer"
	TagButton     Tag = "button"
	TagSpan       Tag = "span"
	TagText       Tag = "#text"
	TagImage      Tag = "img"
	TagTextBlock  Tag = "text-block"
	TagSpacer     Tag = "spacer"
	TagStack      Tag = "stack"
	TagScrollView Tag = "scroll-view"
)

type Props struct {
	ID        string
	ClassName string
	Semantic  SemanticSpec
	Layout    LayoutSpec
	Style     Style
	State     InteractionState
	Image     ImageSource
	Scroll    ScrollState
	Focusable bool
	Handlers  EventHandlers
}

type Node struct {
	Tag      Tag
	Props    Props
	Text     string
	Children []*Node
}

func Element(tag Tag, props Props, children ...*Node) *Node {
	node := &Node{
		Tag:      tag,
		Props:    props,
		Children: make([]*Node, 0, len(children)),
	}
	for _, child := range children {
		if child == nil {
			continue
		}
		node.Children = append(node.Children, child)
	}
	return node
}

func Div(props Props, children ...*Node) *Node {
	return Element(TagDiv, props, children...)
}

func Header(props Props, children ...*Node) *Node {
	return Element(TagHeader, props, children...)
}

func Main(props Props, children ...*Node) *Node {
	return Element(TagMain, props, children...)
}

func Section(props Props, children ...*Node) *Node {
	return Element(TagSection, props, children...)
}

func Footer(props Props, children ...*Node) *Node {
	return Element(TagFooter, props, children...)
}

func Button(props Props, children ...*Node) *Node {
	return Element(TagButton, props, children...)
}

func Span(props Props, children ...*Node) *Node {
	return Element(TagSpan, props, children...)
}

func Text(content string, props Props) *Node {
	node := Element(TagText, props)
	node.Text = content
	return node
}

func TextBlock(content string, props Props) *Node {
	node := Element(TagTextBlock, props)
	node.Text = content
	return node
}

func Image(props Props) *Node {
	return Element(TagImage, props)
}

func Spacer(props Props) *Node {
	return Element(TagSpacer, props)
}

func Stack(props Props, children ...*Node) *Node {
	return Element(TagStack, props, children...)
}

func ScrollView(props Props, children ...*Node) *Node {
	return Element(TagScrollView, props, children...)
}

func InteractiveButton(props Props, children ...*Node) *Node {
	props.Focusable = true
	return Element(TagButton, props, children...)
}

func (n *Node) FindByID(id string) (*Node, bool) {
	if n == nil {
		return nil, false
	}
	if n.Props.ID == id {
		return n, true
	}
	for _, child := range n.Children {
		if found, ok := child.FindByID(id); ok {
			return found, true
		}
	}
	return nil, false
}

type DOM struct {
	Root *Node
}

func New(root *Node) *DOM {
	return &DOM{Root: root}
}

func (d *DOM) FindByID(id string) (*Node, bool) {
	if d == nil {
		return nil, false
	}
	return d.Root.FindByID(id)
}
