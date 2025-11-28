package lib

import "fmt"

// NodeType represents the semantic type of a node in the AST
type NodeType int

const (
	Document NodeType = iota
	Section
	Paragraph
	BlockMacro
	InlineMacro
	Text
	List
	ListItem
	CodeBlock
	LiteralBlock
	Example
	Sidebar
	Quote
	Table
	TableRow
	TableCell
	Admonition
	ThematicBreak
	PageBreak
	Bold
	Italic
	Monospace
	Link
	Passthrough
)

// String returns a human-readable name for the NodeType
func (t NodeType) String() string {
	switch t {
	case Document:
		return "Document"
	case Section:
		return "Section"
	case Paragraph:
		return "Paragraph"
	case BlockMacro:
		return "BlockMacro"
	case InlineMacro:
		return "InlineMacro"
	case Text:
		return "Text"
	case List:
		return "List"
	case ListItem:
		return "ListItem"
	case CodeBlock:
		return "CodeBlock"
	case LiteralBlock:
		return "LiteralBlock"
	case Example:
		return "Example"
	case Sidebar:
		return "Sidebar"
	case Quote:
		return "Quote"
	case Table:
		return "Table"
	case TableRow:
		return "TableRow"
	case TableCell:
		return "TableCell"
	case Admonition:
		return "Admonition"
	case ThematicBreak:
		return "ThematicBreak"
	case PageBreak:
		return "PageBreak"
	case Bold:
		return "Bold"
	case Italic:
		return "Italic"
	case Monospace:
		return "Monospace"
	case Link:
		return "Link"
	case Passthrough:
		return "Passthrough"
	default:
		return "Unknown"
	}
}

// Node represents a node in the Abstract Syntax Tree
type Node struct {
	Type       NodeType
	Content    string            // For Text nodes
	Name       string            // For BlockMacro/InlineMacro (macro name)
	Attributes map[string]string // Key-value pairs for attributes
	Children   []*Node
}

// NewDocumentNode creates a new Document node
func NewDocumentNode() *Node {
	return &Node{
		Type:       Document,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewSectionNode creates a new Section node with the specified level
func NewSectionNode(level int) *Node {
	node := &Node{
		Type:       Section,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
	node.SetAttribute("level", fmt.Sprintf("%d", level))
	return node
}

// NewParagraphNode creates a new Paragraph node
func NewParagraphNode() *Node {
	return &Node{
		Type:       Paragraph,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewBlockMacroNode creates a new BlockMacro node with the specified macro name
func NewBlockMacroNode(name string) *Node {
	return &Node{
		Type:       BlockMacro,
		Name:       name,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewInlineMacroNode creates a new InlineMacro node with the specified macro name
func NewInlineMacroNode(name string) *Node {
	return &Node{
		Type:       InlineMacro,
		Name:       name,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewTextNode creates a new Text node with the specified content
func NewTextNode(content string) *Node {
	return &Node{
		Type:     Text,
		Content:  content,
		Children: make([]*Node, 0),
	}
}

// NewListNode creates a new List node
func NewListNode() *Node {
	return &Node{
		Type:       List,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewListItemNode creates a new ListItem node
func NewListItemNode() *Node {
	return &Node{
		Type:       ListItem,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewCodeBlockNode creates a new CodeBlock node
func NewCodeBlockNode() *Node {
	return &Node{
		Type:       CodeBlock,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewLiteralBlockNode creates a new LiteralBlock node
func NewLiteralBlockNode() *Node {
	return &Node{
		Type:       LiteralBlock,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewExampleNode creates a new Example node
func NewExampleNode() *Node {
	return &Node{
		Type:       Example,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewSidebarNode creates a new Sidebar node
func NewSidebarNode() *Node {
	return &Node{
		Type:       Sidebar,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewQuoteNode creates a new Quote node
func NewQuoteNode() *Node {
	return &Node{
		Type:       Quote,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewTableNode creates a new Table node
func NewTableNode() *Node {
	return &Node{
		Type:       Table,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewTableRowNode creates a new TableRow node
func NewTableRowNode() *Node {
	return &Node{
		Type:       TableRow,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewTableCellNode creates a new TableCell node
func NewTableCellNode() *Node {
	return &Node{
		Type:       TableCell,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewAdmonitionNode creates a new Admonition node
func NewAdmonitionNode() *Node {
	return &Node{
		Type:       Admonition,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewThematicBreakNode creates a new ThematicBreak node
func NewThematicBreakNode() *Node {
	return &Node{
		Type:       ThematicBreak,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewPageBreakNode creates a new PageBreak node
func NewPageBreakNode() *Node {
	return &Node{
		Type:       PageBreak,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewBoldNode creates a new Bold node
func NewBoldNode() *Node {
	return &Node{
		Type:       Bold,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewItalicNode creates a new Italic node
func NewItalicNode() *Node {
	return &Node{
		Type:       Italic,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewMonospaceNode creates a new Monospace node
func NewMonospaceNode() *Node {
	return &Node{
		Type:       Monospace,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewLinkNode creates a new Link node
func NewLinkNode() *Node {
	return &Node{
		Type:       Link,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

// NewPassthroughNode creates a new Passthrough node with the specified content
func NewPassthroughNode(content string) *Node {
	return &Node{
		Type:     Passthrough,
		Content:  content,
		Children: make([]*Node, 0),
	}
}

// AddChild adds a child node to this node
func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

// SetAttribute sets an attribute on the node
func (n *Node) SetAttribute(key, value string) {
	if n.Attributes == nil {
		n.Attributes = make(map[string]string)
	}
	n.Attributes[key] = value
}

// GetAttribute retrieves an attribute value, returning empty string if not found
func (n *Node) GetAttribute(key string) string {
	if n.Attributes == nil {
		return ""
	}
	return n.Attributes[key]
}

// Traverse traverses the AST tree depth-first, calling visit for each node
func (n *Node) Traverse(visit func(*Node)) {
	visit(n)
	for _, child := range n.Children {
		child.Traverse(visit)
	}
}

// FindElementsByTag performs a recursive depth-first search and returns all nodes
// matching the tag name. For macro searches (e.g. "component"), checks node.Type == BlockMacro && node.Name == tagName
func (n *Node) FindElementsByTag(tagName string) []*Node {
	var results []*Node

	// Check if current node matches
	// For macros, check Type and Name
	if n.Type == BlockMacro && n.Name == tagName {
		results = append(results, n)
	} else if n.Type == InlineMacro && n.Name == tagName {
		results = append(results, n)
	} else {
		// For other types, check if tagName matches the type name (for backward compatibility)
		if n.Type.String() == tagName {
			results = append(results, n)
		}
	}

	// Recursively search children
	for _, child := range n.Children {
		childResults := child.FindElementsByTag(tagName)
		results = append(results, childResults...)
	}

	return results
}

