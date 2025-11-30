package lib

import (
	"testing"
)

func TestNodeType_String(t *testing.T) {
	tests := []struct {
		nodeType NodeType
		expected string
	}{
		{Document, "Document"},
		{Section, "Section"},
		{Paragraph, "Paragraph"},
		{BlockMacro, "BlockMacro"},
		{InlineMacro, "InlineMacro"},
		{Text, "Text"},
		{List, "List"},
		{ListItem, "ListItem"},
		{CodeBlock, "CodeBlock"},
		{LiteralBlock, "LiteralBlock"},
		{Example, "Example"},
		{Sidebar, "Sidebar"},
		{Quote, "Quote"},
		{Table, "Table"},
		{TableRow, "TableRow"},
		{TableCell, "TableCell"},
		{Admonition, "Admonition"},
		{ThematicBreak, "ThematicBreak"},
		{PageBreak, "PageBreak"},
		{Bold, "Bold"},
		{Italic, "Italic"},
		{Monospace, "Monospace"},
		{Link, "Link"},
		{Passthrough, "Passthrough"},
	}

	for _, tt := range tests {
		result := tt.nodeType.String()
		if result != tt.expected {
			t.Errorf("NodeType(%d).String() = %q, want %q", tt.nodeType, result, tt.expected)
		}
	}
}

func TestNewDocumentNode(t *testing.T) {
	node := NewDocumentNode()
	if node.Type != Document {
		t.Errorf("Expected Document type, got %v", node.Type)
	}
	if node.Attributes == nil {
		t.Error("Attributes map should be initialized")
	}
	if node.Children == nil {
		t.Error("Children slice should be initialized")
	}
}

func TestNewSectionNode(t *testing.T) {
	node := NewSectionNode(2)
	if node.Type != Section {
		t.Errorf("Expected Section type, got %v", node.Type)
	}
	level := node.GetAttribute("level")
	if level != "2" {
		t.Errorf("Expected level '2', got %q", level)
	}
}

func TestNewTextNode(t *testing.T) {
	content := "test content"
	node := NewTextNode(content)
	if node.Type != Text {
		t.Errorf("Expected Text type, got %v", node.Type)
	}
	if node.Content != content {
		t.Errorf("Expected content %q, got %q", content, node.Content)
	}
}

func TestNewBlockMacroNode(t *testing.T) {
	name := "component"
	node := NewBlockMacroNode(name)
	if node.Type != BlockMacro {
		t.Errorf("Expected BlockMacro type, got %v", node.Type)
	}
	if node.Name != name {
		t.Errorf("Expected name %q, got %q", name, node.Name)
	}
}

func TestNewInlineMacroNode(t *testing.T) {
	name := "link"
	node := NewInlineMacroNode(name)
	if node.Type != InlineMacro {
		t.Errorf("Expected InlineMacro type, got %v", node.Type)
	}
	if node.Name != name {
		t.Errorf("Expected name %q, got %q", name, node.Name)
	}
}

func TestNode_AddChild(t *testing.T) {
	parent := NewDocumentNode()
	child1 := NewParagraphNode()
	child2 := NewTextNode("text")

	parent.AddChild(child1)
	if len(parent.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(parent.Children))
	}

	parent.AddChild(child2)
	if len(parent.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(parent.Children))
	}

	if parent.Children[0] != child1 {
		t.Error("First child should be child1")
	}
	if parent.Children[1] != child2 {
		t.Error("Second child should be child2")
	}
}

func TestNode_SetAttribute(t *testing.T) {
	node := NewDocumentNode()

	// Test setting attribute
	node.SetAttribute("title", "Test Title")
	if node.GetAttribute("title") != "Test Title" {
		t.Errorf("Expected 'Test Title', got %q", node.GetAttribute("title"))
	}

	// Test overwriting attribute
	node.SetAttribute("title", "New Title")
	if node.GetAttribute("title") != "New Title" {
		t.Errorf("Expected 'New Title', got %q", node.GetAttribute("title"))
	}

	// Test setting multiple attributes
	node.SetAttribute("author", "John Doe")
	node.SetAttribute("version", "1.0")

	if node.GetAttribute("author") != "John Doe" {
		t.Errorf("Expected 'John Doe', got %q", node.GetAttribute("author"))
	}
	if node.GetAttribute("version") != "1.0" {
		t.Errorf("Expected '1.0', got %q", node.GetAttribute("version"))
	}
}

func TestNode_GetAttribute(t *testing.T) {
	node := NewDocumentNode()

	// Test getting non-existent attribute
	if node.GetAttribute("nonexistent") != "" {
		t.Error("Expected empty string for non-existent attribute")
	}

	// Test getting existing attribute
	node.SetAttribute("test", "value")
	if node.GetAttribute("test") != "value" {
		t.Errorf("Expected 'value', got %q", node.GetAttribute("test"))
	}

	// Test with nil attributes map (should not panic)
	node2 := &Node{Type: Document}
	if node2.GetAttribute("test") != "" {
		t.Error("Expected empty string when attributes is nil")
	}
}

func TestNode_Traverse(t *testing.T) {
	// Create a tree structure
	// Document
	//   - Section
	//     - Paragraph
	//       - Text
	//   - Paragraph
	//     - Text

	doc := NewDocumentNode()
	section := NewSectionNode(1)
	para1 := NewParagraphNode()
	text1 := NewTextNode("text1")
	para2 := NewParagraphNode()
	text2 := NewTextNode("text2")

	doc.AddChild(section)
	section.AddChild(para1)
	para1.AddChild(text1)
	doc.AddChild(para2)
	para2.AddChild(text2)

	// Collect visited nodes
	var visited []*Node
	doc.Traverse(func(n *Node) {
		visited = append(visited, n)
	})

	// Should visit all 6 nodes (doc, section, para1, text1, para2, text2)
	if len(visited) != 6 {
		t.Errorf("Expected 6 nodes visited, got %d", len(visited))
	}

	// Verify order (depth-first)
	// We have 5 nodes: doc, section, para1, text1, para2, text2
	// Visited will be: doc, section, para1, text1, para2, text2
	if visited[0].Type != Document {
		t.Error("First node should be Document")
	}
	if visited[1].Type != Section {
		t.Error("Second node should be Section")
	}
	if visited[2].Type != Paragraph {
		t.Error("Third node should be Paragraph")
	}
	if visited[3].Type != Text {
		t.Error("Fourth node should be Text")
	}
	if visited[4].Type != Paragraph {
		t.Error("Fifth node should be Paragraph")
	}
}

func TestNode_Traverse_EmptyTree(t *testing.T) {
	node := NewDocumentNode()
	count := 0
	node.Traverse(func(n *Node) {
		count++
	})
	if count != 1 {
		t.Errorf("Expected 1 node visited, got %d", count)
	}
}

func TestNode_FindElementsByTag(t *testing.T) {
	// Create tree with various node types
	doc := NewDocumentNode()
	section1 := NewSectionNode(1)
	section2 := NewSectionNode(2)
	para1 := NewParagraphNode()
	para2 := NewParagraphNode()

	doc.AddChild(section1)
	doc.AddChild(section2)
	section1.AddChild(para1)
	section2.AddChild(para2)

	// Find all sections
	sections := doc.FindElementsByTag("Section")
	if len(sections) != 2 {
		t.Errorf("Expected 2 sections, got %d", len(sections))
	}

	// Find all paragraphs
	paragraphs := doc.FindElementsByTag("Paragraph")
	if len(paragraphs) != 2 {
		t.Errorf("Expected 2 paragraphs, got %d", len(paragraphs))
	}

	// Find non-existent type
	none := doc.FindElementsByTag("NonExistent")
	if len(none) != 0 {
		t.Errorf("Expected 0 results, got %d", len(none))
	}
}

func TestNode_FindElementsByTag_BlockMacro(t *testing.T) {
	doc := NewDocumentNode()
	component1 := NewBlockMacroNode("component")
	component2 := NewBlockMacroNode("component")
	image := NewBlockMacroNode("image")

	doc.AddChild(component1)
	doc.AddChild(component2)
	doc.AddChild(image)

	// Find component macros
	components := doc.FindElementsByTag("component")
	if len(components) != 2 {
		t.Errorf("Expected 2 component macros, got %d", len(components))
	}

	// Find image macro
	images := doc.FindElementsByTag("image")
	if len(images) != 1 {
		t.Errorf("Expected 1 image macro, got %d", len(images))
	}
}

func TestNode_FindElementsByTag_InlineMacro(t *testing.T) {
	para := NewParagraphNode()
	link1 := NewInlineMacroNode("link")
	link2 := NewInlineMacroNode("link")
	button := NewInlineMacroNode("button")

	para.AddChild(link1)
	para.AddChild(link2)
	para.AddChild(button)

	// Find link macros
	links := para.FindElementsByTag("link")
	if len(links) != 2 {
		t.Errorf("Expected 2 link macros, got %d", len(links))
	}

	// Find button macro
	buttons := para.FindElementsByTag("button")
	if len(buttons) != 1 {
		t.Errorf("Expected 1 button macro, got %d", len(buttons))
	}
}

func TestNode_FindElementsByTag_Nested(t *testing.T) {
	doc := NewDocumentNode()
	section := NewSectionNode(1)
	para1 := NewParagraphNode()
	para2 := NewParagraphNode()

	doc.AddChild(section)
	section.AddChild(para1)
	section.AddChild(para2)

	// Find paragraphs (should find nested ones)
	paragraphs := doc.FindElementsByTag("Paragraph")
	if len(paragraphs) != 2 {
		t.Errorf("Expected 2 paragraphs (including nested), got %d", len(paragraphs))
	}
}

func TestNode_FindElementsByTag_SelfMatch(t *testing.T) {
	para := NewParagraphNode()
	results := para.FindElementsByTag("Paragraph")
	if len(results) != 1 {
		t.Errorf("Expected node to match itself, got %d results", len(results))
	}
	if results[0] != para {
		t.Error("Result should be the node itself")
	}
}

func TestNode_NodeCreationFunctions(t *testing.T) {
	// Test all node creation functions
	tests := []struct {
		name     string
		createFn func() *Node
		expected NodeType
	}{
		{"NewParagraphNode", NewParagraphNode, Paragraph},
		{"NewListNode", NewListNode, List},
		{"NewListItemNode", NewListItemNode, ListItem},
		{"NewCodeBlockNode", NewCodeBlockNode, CodeBlock},
		{"NewLiteralBlockNode", NewLiteralBlockNode, LiteralBlock},
		{"NewExampleNode", NewExampleNode, Example},
		{"NewSidebarNode", NewSidebarNode, Sidebar},
		{"NewQuoteNode", NewQuoteNode, Quote},
		{"NewTableNode", NewTableNode, Table},
		{"NewTableRowNode", NewTableRowNode, TableRow},
		{"NewTableCellNode", NewTableCellNode, TableCell},
		{"NewAdmonitionNode", NewAdmonitionNode, Admonition},
		{"NewThematicBreakNode", NewThematicBreakNode, ThematicBreak},
		{"NewPageBreakNode", NewPageBreakNode, PageBreak},
		{"NewBoldNode", NewBoldNode, Bold},
		{"NewItalicNode", NewItalicNode, Italic},
		{"NewMonospaceNode", NewMonospaceNode, Monospace},
		{"NewLinkNode", NewLinkNode, Link},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.createFn()
			if node.Type != tt.expected {
				t.Errorf("Expected type %v, got %v", tt.expected, node.Type)
			}
			if node.Attributes == nil {
				t.Error("Attributes should be initialized")
			}
			if node.Children == nil {
				t.Error("Children should be initialized")
			}
		})
	}
}

func TestNode_PassthroughNode(t *testing.T) {
	content := "passthrough content"
	node := NewPassthroughNode(content)
	if node.Type != Passthrough {
		t.Errorf("Expected Passthrough type, got %v", node.Type)
	}
	if node.Content != content {
		t.Errorf("Expected content %q, got %q", content, node.Content)
	}
}

func TestNode_ComplexTree(t *testing.T) {
	// Create a complex tree structure
	doc := NewDocumentNode()
	doc.SetAttribute("title", "Test Document")

	section1 := NewSectionNode(1)
	section1.SetAttribute("id", "section1")
	para1 := NewParagraphNode()
	text1 := NewTextNode("First paragraph")
	bold1 := NewBoldNode()
	boldText := NewTextNode("bold text")
	bold1.AddChild(boldText)
	para1.AddChild(text1)
	para1.AddChild(bold1)

	section2 := NewSectionNode(2)
	section2.SetAttribute("id", "section2")
	list := NewListNode()
	item1 := NewListItemNode()
	item1Text := NewTextNode("Item 1")
	item1.AddChild(item1Text)
	list.AddChild(item1)

	doc.AddChild(section1)
	section1.AddChild(para1)
	doc.AddChild(section2)
	section2.AddChild(list)

	// Verify structure
	if len(doc.Children) != 2 {
		t.Errorf("Expected 2 top-level children, got %d", len(doc.Children))
	}

	// Verify attributes
	if doc.GetAttribute("title") != "Test Document" {
		t.Error("Document title attribute not set correctly")
	}
	if section1.GetAttribute("id") != "section1" {
		t.Error("Section1 id attribute not set correctly")
	}

	// Find all text nodes
	var textNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == Text {
			textNodes = append(textNodes, n)
		}
	})

	// We have: text1, boldText, item1Text = 3 text nodes
	if len(textNodes) != 3 {
		t.Errorf("Expected 3 text nodes, got %d", len(textNodes))
	}
}

func TestNode_SetAttribute_NilMap(t *testing.T) {
	// Test SetAttribute when Attributes is nil
	node := &Node{Type: Document}
	node.SetAttribute("test", "value")
	if node.GetAttribute("test") != "value" {
		t.Error("SetAttribute should initialize map if nil")
	}
}

