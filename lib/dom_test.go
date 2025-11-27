package lib

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewElementNode(t *testing.T) {
	node := NewElementNode("test")
	
	if node == nil {
		t.Fatal("NewElementNode should not return nil")
	}
	
	if node.Type != ElementNode {
		t.Errorf("Expected ElementNode, got %v", node.Type)
	}
	
	if node.Data != "test" {
		t.Errorf("Expected Data 'test', got '%s'", node.Data)
	}
	
	if node.Attributes == nil {
		t.Error("Attributes map should be initialized")
	}
	
	if node.Children == nil {
		t.Error("Children slice should be initialized")
	}
	
	if len(node.Children) != 0 {
		t.Errorf("Expected empty children slice, got %d", len(node.Children))
	}
}

func TestNewTextNode(t *testing.T) {
	text := "Hello, World!"
	node := NewTextNode(text)
	
	if node == nil {
		t.Fatal("NewTextNode should not return nil")
	}
	
	if node.Type != TextNode {
		t.Errorf("Expected TextNode, got %v", node.Type)
	}
	
	if node.Data != text {
		t.Errorf("Expected Data '%s', got '%s'", text, node.Data)
	}
	
	if node.Children == nil {
		t.Error("Children slice should be initialized")
	}
}

func TestNode_AddChild(t *testing.T) {
	parent := NewElementNode("parent")
	child := NewElementNode("child")
	
	parent.AddChild(child)
	
	if len(parent.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(parent.Children))
	}
	
	if parent.Children[0] != child {
		t.Error("Child should be added to parent")
	}
	
	if child.Parent != parent {
		t.Error("Child's parent should be set")
	}
	
	// Test adding multiple children
	child2 := NewTextNode("text")
	parent.AddChild(child2)
	
	if len(parent.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(parent.Children))
	}
}

func TestNode_SetAttribute(t *testing.T) {
	node := NewElementNode("test")
	
	// Test setting first attribute
	node.SetAttribute("id", "test-id")
	
	if node.Attributes["id"] != "test-id" {
		t.Errorf("Expected attribute 'id' to be 'test-id', got '%s'", node.Attributes["id"])
	}
	
	// Test setting multiple attributes
	node.SetAttribute("class", "test-class")
	node.SetAttribute("role", "test-role")
	
	if len(node.Attributes) != 3 {
		t.Errorf("Expected 3 attributes, got %d", len(node.Attributes))
	}
	
	// Test overwriting attribute
	node.SetAttribute("id", "new-id")
	if node.Attributes["id"] != "new-id" {
		t.Errorf("Expected attribute 'id' to be 'new-id', got '%s'", node.Attributes["id"])
	}
}

func TestNode_GetAttribute(t *testing.T) {
	node := NewElementNode("test")
	
	// Test getting non-existent attribute
	value := node.GetAttribute("nonexistent")
	if value != "" {
		t.Errorf("Expected empty string for non-existent attribute, got '%s'", value)
	}
	
	// Test getting existing attribute
	node.SetAttribute("id", "test-id")
	value = node.GetAttribute("id")
	if value != "test-id" {
		t.Errorf("Expected 'test-id', got '%s'", value)
	}
	
	// Test with nil attributes map (should not panic)
	node2 := &Node{Type: ElementNode, Data: "test"}
	value = node2.GetAttribute("test")
	if value != "" {
		t.Errorf("Expected empty string for node with nil attributes, got '%s'", value)
	}
}

func TestNode_Traverse(t *testing.T) {
	root := NewElementNode("root")
	child1 := NewElementNode("child1")
	child2 := NewElementNode("child2")
	grandchild := NewElementNode("grandchild")
	
	root.AddChild(child1)
	root.AddChild(child2)
	child1.AddChild(grandchild)
	
	visited := []string{}
	root.Traverse(func(n *Node) {
		visited = append(visited, n.Data)
	})
	
	expected := []string{"root", "child1", "grandchild", "child2"}
	if len(visited) != len(expected) {
		t.Errorf("Expected %d nodes visited, got %d", len(expected), len(visited))
	}
	
	for i, v := range expected {
		if i >= len(visited) || visited[i] != v {
			t.Errorf("Expected visit order %v, got %v", expected, visited)
			break
		}
	}
}

func TestNode_ToXML_SimpleElement(t *testing.T) {
	node := NewElementNode("test")
	
	xml, err := node.ToXML()
	if err != nil {
		t.Fatalf("ToXML failed: %v", err)
	}
	
	if !strings.Contains(xml, "<test") {
		t.Errorf("Expected XML to contain '<test', got: %s", xml)
	}
}

func TestNode_ToXML_ElementWithAttributes(t *testing.T) {
	node := NewElementNode("test")
	node.SetAttribute("id", "test-id")
	node.SetAttribute("class", "test-class")
	
	xml, err := node.ToXML()
	if err != nil {
		t.Fatalf("ToXML failed: %v", err)
	}
	
	if !strings.Contains(xml, "id=\"test-id\"") {
		t.Errorf("Expected XML to contain 'id=\"test-id\"', got: %s", xml)
	}
	
	if !strings.Contains(xml, "class=\"test-class\"") {
		t.Errorf("Expected XML to contain 'class=\"test-class\"', got: %s", xml)
	}
}

func TestNode_ToXML_ElementWithTextChild(t *testing.T) {
	node := NewElementNode("test")
	textNode := NewTextNode("Hello, World!")
	node.AddChild(textNode)
	
	xml, err := node.ToXML()
	if err != nil {
		t.Fatalf("ToXML failed: %v", err)
	}
	
	if !strings.Contains(xml, "Hello, World!") {
		t.Errorf("Expected XML to contain 'Hello, World!', got: %s", xml)
	}
	
	if !strings.Contains(xml, "<test>") {
		t.Errorf("Expected XML to contain '<test>', got: %s", xml)
	}
	
	if !strings.Contains(xml, "</test>") {
		t.Errorf("Expected XML to contain '</test>', got: %s", xml)
	}
}

func TestNode_ToXML_ElementWithElementChild(t *testing.T) {
	parent := NewElementNode("parent")
	child := NewElementNode("child")
	parent.AddChild(child)
	
	xml, err := parent.ToXML()
	if err != nil {
		t.Fatalf("ToXML failed: %v", err)
	}
	
	if !strings.Contains(xml, "<parent>") {
		t.Errorf("Expected XML to contain '<parent>', got: %s", xml)
	}
	
	if !strings.Contains(xml, "<child") {
		t.Errorf("Expected XML to contain '<child', got: %s", xml)
	}
}

func TestNode_ToXML_ComplexStructure(t *testing.T) {
	root := NewElementNode("root")
	root.SetAttribute("id", "root-id")
	
	header := NewElementNode("header")
	title := NewElementNode("title")
	titleText := NewTextNode("Document Title")
	title.AddChild(titleText)
	header.AddChild(title)
	root.AddChild(header)
	
	section := NewElementNode("section")
	section.SetAttribute("level", "1")
	para := NewElementNode("paragraph")
	paraText := NewTextNode("This is a paragraph.")
	para.AddChild(paraText)
	section.AddChild(para)
	root.AddChild(section)
	
	xml, err := root.ToXML()
	if err != nil {
		t.Fatalf("ToXML failed: %v", err)
	}
	
	// Verify structure
	if !strings.Contains(xml, "<root") {
		t.Error("Expected XML to contain '<root'")
	}
	
	if !strings.Contains(xml, "<header>") {
		t.Error("Expected XML to contain '<header>'")
	}
	
	if !strings.Contains(xml, "<title>") {
		t.Error("Expected XML to contain '<title>'")
	}
	
	if !strings.Contains(xml, "Document Title") {
		t.Error("Expected XML to contain 'Document Title'")
	}
	
	if !strings.Contains(xml, "<section") {
		t.Error("Expected XML to contain '<section'")
	}
	
	if !strings.Contains(xml, "level=\"1\"") {
		t.Error("Expected XML to contain 'level=\"1\"'")
	}
}

func TestNode_ToXML_XMLEscaping(t *testing.T) {
	node := NewElementNode("test")
	
	// Test various characters that need escaping
	testCases := []struct {
		name     string
		text     string
		expected string
	}{
		{"less than", "<", "&lt;"},
		{"greater than", ">", "&gt;"},
		{"ampersand", "&", "&amp;"},
		{"quote", "\"", "&quot;"},
		{"apostrophe", "'", "&apos;"},
		{"mixed", "A < B & C > D", "&lt;"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			textNode := NewTextNode(tc.text)
			node.Children = []*Node{textNode}
			
			xml, err := node.ToXML()
			if err != nil {
				t.Fatalf("ToXML failed: %v", err)
			}
			
			if !strings.Contains(xml, tc.expected) {
				t.Errorf("Expected XML to contain '%s', got: %s", tc.expected, xml)
			}
		})
	}
}

func TestNode_ToXML_AttributeEscaping(t *testing.T) {
	node := NewElementNode("test")
	
	// Test attribute value escaping
	node.SetAttribute("attr", "value with \"quotes\"")
	
	xml, err := node.ToXML()
	if err != nil {
		t.Fatalf("ToXML failed: %v", err)
	}
	
	// Quotes in attribute values should be escaped
	if !strings.Contains(xml, "&quot;") {
		t.Errorf("Expected XML to escape quotes in attributes, got: %s", xml)
	}
}

func TestNode_ToXML_MixedContent(t *testing.T) {
	para := NewElementNode("paragraph")
	para.AddChild(NewTextNode("This is "))
	
	strong := NewElementNode("strong")
	strong.AddChild(NewTextNode("bold"))
	para.AddChild(strong)
	
	para.AddChild(NewTextNode(" text."))
	
	xml, err := para.ToXML()
	if err != nil {
		t.Fatalf("ToXML failed: %v", err)
	}
	
	// Mixed content should be handled correctly
	if !strings.Contains(xml, "This is ") {
		t.Error("Expected XML to contain 'This is '")
	}
	
	if !strings.Contains(xml, "<strong>") {
		t.Error("Expected XML to contain '<strong>'")
	}
	
	if !strings.Contains(xml, "bold") {
		t.Error("Expected XML to contain 'bold'")
	}
	
	if !strings.Contains(xml, " text.") {
		t.Error("Expected XML to contain ' text.'")
	}
}

func TestEscapeXML(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"no escaping needed", "Hello World", "Hello World"},
		{"less than", "<", "&lt;"},
		{"greater than", ">", "&gt;"},
		{"ampersand", "&", "&amp;"},
		{"quote", "\"", "&quot;"},
		{"apostrophe", "'", "&apos;"},
		{"mixed", "A < B & C > D \"test\" 'test'", "&lt;"},
		{"empty string", "", ""},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := escapeXML(tc.input)
			
			// Verify all special characters are escaped
			if strings.Contains(result, "<") && !strings.Contains(result, "&lt;") {
				t.Errorf("Expected '<' to be escaped in '%s'", result)
			}
			if strings.Contains(result, ">") && !strings.Contains(result, "&gt;") {
				t.Errorf("Expected '>' to be escaped in '%s'", result)
			}
			if strings.Contains(result, "&") && !strings.Contains(result, "&amp;") && !strings.Contains(result, "&lt;") && !strings.Contains(result, "&gt;") && !strings.Contains(result, "&quot;") && !strings.Contains(result, "&apos;") {
				t.Errorf("Expected '&' to be escaped in '%s'", result)
			}
		})
	}
}

func TestHasMixedContent(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Node
		expected bool
	}{
		{
			name: "only element children",
			setup: func() *Node {
				node := NewElementNode("parent")
				node.AddChild(NewElementNode("child"))
				return node
			},
			expected: false,
		},
		{
			name: "only text children",
			setup: func() *Node {
				node := NewElementNode("parent")
				node.AddChild(NewTextNode("text"))
				return node
			},
			expected: true,
		},
		{
			name: "mixed children",
			setup: func() *Node {
				node := NewElementNode("parent")
				node.AddChild(NewTextNode("text"))
				node.AddChild(NewElementNode("child"))
				return node
			},
			expected: true,
		},
		{
			name: "no children",
			setup: func() *Node {
				return NewElementNode("parent")
			},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.setup()
			result := hasMixedContent(node)
			if result != tt.expected {
				t.Errorf("Expected hasMixedContent to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNode_ToXML_EmptyElement(t *testing.T) {
	node := NewElementNode("empty")
	
	xml, err := node.ToXML()
	if err != nil {
		t.Fatalf("ToXML failed: %v", err)
	}
	
	// Empty element should be self-closing or have no content
	if !strings.Contains(xml, "<empty") {
		t.Errorf("Expected XML to contain '<empty', got: %s", xml)
	}
}

func TestNode_ToXML_DeepNesting(t *testing.T) {
	// Create a deeply nested structure
	root := NewElementNode("root")
	current := root
	
	for i := 0; i < 10; i++ {
		child := NewElementNode("level")
		child.SetAttribute("depth", string(rune('0'+i)))
		current.AddChild(child)
		current = child
	}
	
	// Add text at the deepest level
	current.AddChild(NewTextNode("deep content"))
	
	xml, err := root.ToXML()
	if err != nil {
		t.Fatalf("ToXML failed: %v", err)
	}
	
	// Verify deep nesting is handled
	if !strings.Contains(xml, "deep content") {
		t.Error("Expected XML to contain 'deep content'")
	}
	
	// Verify structure by counting opening and closing tags
	openTags := strings.Count(xml, "<level")
	closeTags := strings.Count(xml, "</level")
	if openTags != closeTags {
		t.Errorf("XML structure may be unbalanced: %d opening tags, %d closing tags", openTags, closeTags)
	}
}

func TestXmlEscape(t *testing.T) {
	var buf bytes.Buffer
	
	testCases := []struct {
		input    string
		expected string
	}{
		{"normal text", "normal text"},
		{"<tag>", "&lt;tag&gt;"},
		{"A & B", "A &amp; B"},
		{`"quoted"`, "&quot;quoted&quot;"},
		{"'apostrophe'", "&apos;apostrophe&apos;"},
	}
	
	for _, tc := range testCases {
		buf.Reset()
		err := xmlEscape(&buf, tc.input)
		if err != nil {
			t.Fatalf("xmlEscape failed: %v", err)
		}
		
		result := buf.String()
		// Verify escaping worked
		if strings.Contains(tc.input, "<") && !strings.Contains(result, "&lt;") {
			t.Errorf("Expected '<' to be escaped in '%s', got '%s'", tc.input, result)
		}
		if strings.Contains(tc.input, ">") && !strings.Contains(result, "&gt;") {
			t.Errorf("Expected '>' to be escaped in '%s', got '%s'", tc.input, result)
		}
		if strings.Contains(tc.input, "&") && !strings.Contains(result, "&amp;") && !strings.Contains(result, "&lt;") && !strings.Contains(result, "&gt;") && !strings.Contains(result, "&quot;") && !strings.Contains(result, "&apos;") {
			t.Errorf("Expected '&' to be escaped in '%s', got '%s'", tc.input, result)
		}
	}
}

