package lib

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseInlineFormatting_Superscript(t *testing.T) {
	input := `This is ^superscript^ text.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find superscript nodes
	var superscriptNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == Superscript {
			superscriptNodes = append(superscriptNodes, n)
		}
	})

	if len(superscriptNodes) == 0 {
		t.Fatal("Expected at least one superscript node")
	}

	if len(superscriptNodes[0].Children) == 0 {
		t.Fatal("Expected superscript node to have children")
	}

	if superscriptNodes[0].Children[0].Content != "superscript" {
		t.Errorf("Expected superscript content 'superscript', got '%s'", superscriptNodes[0].Children[0].Content)
	}
}

func TestParseInlineFormatting_Subscript(t *testing.T) {
	input := `This is ~subscript~ text.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find subscript nodes
	var subscriptNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == Subscript {
			subscriptNodes = append(subscriptNodes, n)
		}
	})

	if len(subscriptNodes) == 0 {
		t.Fatal("Expected at least one subscript node")
	}

	if len(subscriptNodes[0].Children) == 0 {
		t.Fatal("Expected subscript node to have children")
	}

	if subscriptNodes[0].Children[0].Content != "subscript" {
		t.Errorf("Expected subscript content 'subscript', got '%s'", subscriptNodes[0].Children[0].Content)
	}
}

func TestParseInlineFormatting_Highlight(t *testing.T) {
	input := `This is #highlighted# text.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find highlight nodes
	var highlightNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == Highlight {
			highlightNodes = append(highlightNodes, n)
		}
	})

	if len(highlightNodes) == 0 {
		t.Fatal("Expected at least one highlight node")
	}

	if len(highlightNodes[0].Children) == 0 {
		t.Fatal("Expected highlight node to have children")
	}

	if highlightNodes[0].Children[0].Content != "highlighted" {
		t.Errorf("Expected highlight content 'highlighted', got '%s'", highlightNodes[0].Children[0].Content)
	}
}

func TestParseInlineFormatting_InlinePassthrough(t *testing.T) {
	input := `This is +passthrough+ text.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find passthrough nodes
	var passthroughNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == Passthrough {
			passthroughNodes = append(passthroughNodes, n)
		}
	})

	if len(passthroughNodes) == 0 {
		t.Fatal("Expected at least one passthrough node")
	}

	if passthroughNodes[0].Content != "passthrough" {
		t.Errorf("Expected passthrough content 'passthrough', got '%s'", passthroughNodes[0].Content)
	}
}

func TestParseInlineFormatting_Nested(t *testing.T) {
	input := `This is *bold with ^superscript^ inside*.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find bold nodes
	var boldNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == Bold {
			boldNodes = append(boldNodes, n)
		}
	})

	if len(boldNodes) == 0 {
		t.Fatal("Expected at least one bold node")
	}

	// Check for superscript inside bold
	var superscriptFound bool
	boldNodes[0].Traverse(func(n *Node) {
		if n.Type == Superscript {
			superscriptFound = true
		}
	})

	if !superscriptFound {
		t.Fatal("Expected superscript node inside bold node")
	}
}

func TestParseInlineMacros_Kbd(t *testing.T) {
	input := `Press kbd:[Ctrl+C] to copy.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find inline macro nodes
	var macroNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == InlineMacro && n.Name == "kbd" {
			macroNodes = append(macroNodes, n)
		}
	})

	if len(macroNodes) == 0 {
		t.Fatal("Expected at least one kbd macro node")
	}

	if len(macroNodes[0].Children) == 0 {
		t.Fatal("Expected kbd macro node to have children")
	}

	if macroNodes[0].Children[0].Content != "Ctrl+C" {
		t.Errorf("Expected kbd macro content 'Ctrl+C', got '%s'", macroNodes[0].Children[0].Content)
	}
}

func TestParseInlineMacros_Btn(t *testing.T) {
	input := `Click btn:[Submit] to continue.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find inline macro nodes
	var macroNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == InlineMacro && n.Name == "btn" {
			macroNodes = append(macroNodes, n)
		}
	})

	if len(macroNodes) == 0 {
		t.Fatal("Expected at least one btn macro node")
	}

	if len(macroNodes[0].Children) == 0 {
		t.Fatal("Expected btn macro node to have children")
	}

	if macroNodes[0].Children[0].Content != "Submit" {
		t.Errorf("Expected btn macro content 'Submit', got '%s'", macroNodes[0].Children[0].Content)
	}
}

func TestParseInlineMacros_Menu(t *testing.T) {
	input := `Go to menu:File[New Document].`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find inline macro nodes
	var macroNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == InlineMacro && n.Name == "menu" {
			macroNodes = append(macroNodes, n)
		}
	})

	if len(macroNodes) == 0 {
		t.Fatal("Expected at least one menu macro node")
	}

	if macroNodes[0].GetAttribute("target") != "File" {
		t.Errorf("Expected menu macro target 'File', got '%s'", macroNodes[0].GetAttribute("target"))
	}

	if len(macroNodes[0].Children) == 0 {
		t.Fatal("Expected menu macro node to have children")
	}

	if macroNodes[0].Children[0].Content != "New Document" {
		t.Errorf("Expected menu macro content 'New Document', got '%s'", macroNodes[0].Children[0].Content)
	}
}

func TestParseInlineMacros_Generic(t *testing.T) {
	input := `Use custom:macro[text] here.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find inline macro nodes
	var macroNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == InlineMacro && n.Name == "custom" {
			macroNodes = append(macroNodes, n)
		}
	})

	if len(macroNodes) == 0 {
		t.Fatal("Expected at least one custom macro node")
	}

	if macroNodes[0].GetAttribute("target") != "macro" {
		t.Errorf("Expected custom macro target 'macro', got '%s'", macroNodes[0].GetAttribute("target"))
	}

	if len(macroNodes[0].Children) == 0 {
		t.Fatal("Expected custom macro node to have children")
	}

	if macroNodes[0].Children[0].Content != "text" {
		t.Errorf("Expected custom macro content 'text', got '%s'", macroNodes[0].Children[0].Content)
	}
}

func TestParseInlineFormatting_Complex(t *testing.T) {
	input := `This has *bold*, _italic_, ` + "`monospace`" + `, ^sup^, ~sub~, #highlight#, and +pass+ all together.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Count different node types
	counts := make(map[NodeType]int)
	doc.Traverse(func(n *Node) {
		counts[n.Type]++
	})

	if counts[Bold] == 0 {
		t.Error("Expected bold node")
	}
	if counts[Italic] == 0 {
		t.Error("Expected italic node")
	}
	if counts[Monospace] == 0 {
		t.Error("Expected monospace node")
	}
	if counts[Superscript] == 0 {
		t.Error("Expected superscript node")
	}
	if counts[Subscript] == 0 {
		t.Error("Expected subscript node")
	}
	if counts[Highlight] == 0 {
		t.Error("Expected highlight node")
	}
	if counts[Passthrough] == 0 {
		t.Error("Expected passthrough node")
	}
}

func TestParseAnchors_BlockAnchor(t *testing.T) {
	input := `[[my-anchor]]
This is a paragraph with an anchor above it.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find anchor nodes
	var anchorNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == BlockMacro && n.Name == "anchor" {
			anchorNodes = append(anchorNodes, n)
		}
	})

	if len(anchorNodes) == 0 {
		t.Fatal("Expected at least one anchor node")
	}

	if anchorNodes[0].GetAttribute("id") != "my-anchor" {
		t.Errorf("Expected anchor ID 'my-anchor', got '%s'", anchorNodes[0].GetAttribute("id"))
	}
}

func TestParseAnchors_InlineAnchor(t *testing.T) {
	input := `This is a paragraph with an [#inline-anchor] inline anchor.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find anchor nodes
	var anchorNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == InlineMacro && n.Name == "anchor" {
			anchorNodes = append(anchorNodes, n)
		}
	})

	if len(anchorNodes) == 0 {
		t.Fatal("Expected at least one inline anchor node")
	}

	if anchorNodes[0].GetAttribute("id") != "inline-anchor" {
		t.Errorf("Expected anchor ID 'inline-anchor', got '%s'", anchorNodes[0].GetAttribute("id"))
	}
}

func TestParseCrossReferences_Xref(t *testing.T) {
	input := `[[my-anchor]]
This is the anchor.

See <<my-anchor>> for details.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find xref nodes
	var xrefNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == InlineMacro && n.Name == "xref" {
			xrefNodes = append(xrefNodes, n)
		}
	})

	if len(xrefNodes) == 0 {
		t.Fatal("Expected at least one xref node")
	}

	if xrefNodes[0].GetAttribute("target") != "my-anchor" {
		t.Errorf("Expected xref target 'my-anchor', got '%s'", xrefNodes[0].GetAttribute("target"))
	}
}

func TestParseCrossReferences_XrefMacro(t *testing.T) {
	input := `[[my-anchor]]
This is the anchor.

See xref:my-anchor[custom text] for details.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find xref nodes
	var xrefNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == InlineMacro && n.Name == "xref" {
			xrefNodes = append(xrefNodes, n)
		}
	})

	if len(xrefNodes) == 0 {
		t.Fatal("Expected at least one xref node")
	}

	if xrefNodes[0].GetAttribute("target") != "my-anchor" {
		t.Errorf("Expected xref target 'my-anchor', got '%s'", xrefNodes[0].GetAttribute("target"))
	}

	if len(xrefNodes[0].Children) == 0 {
		t.Fatal("Expected xref node to have children")
	}

	if xrefNodes[0].Children[0].Content != "custom text" {
		t.Errorf("Expected xref text 'custom text', got '%s'", xrefNodes[0].Children[0].Content)
	}
}

func TestParseSection_WithID(t *testing.T) {
	input := `[#section-id]
== Section Title

Content here.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find section nodes
	var sectionNodes []*Node
	doc.Traverse(func(n *Node) {
		if n.Type == Section {
			sectionNodes = append(sectionNodes, n)
		}
	})

	if len(sectionNodes) == 0 {
		t.Fatal("Expected at least one section node")
	}

	if sectionNodes[0].GetAttribute("id") != "section-id" {
		t.Errorf("Expected section ID 'section-id', got '%s'", sectionNodes[0].GetAttribute("id"))
	}
}

func TestParseAttributeSubstitution_Basic(t *testing.T) {
	input := `:custom-attr: Custom Value

This is {custom-attr} text.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find text nodes containing the substituted value
	var found bool
	doc.Traverse(func(n *Node) {
		if n.Type == Text && strings.Contains(n.Content, "Custom Value") {
			found = true
		}
	})

	if !found {
		t.Error("Expected attribute substitution to replace {custom-attr} with 'Custom Value'")
	}
}

func TestParseAttributeSubstitution_BuiltIn(t *testing.T) {
	input := `= Document Title
:author: John Doe

Author: {author}`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find text nodes containing the substituted value
	var found bool
	doc.Traverse(func(n *Node) {
		if n.Type == Text && strings.Contains(n.Content, "John Doe") {
			found = true
		}
	})

	if !found {
		t.Error("Expected attribute substitution to replace {author} with 'John Doe'")
	}
}

func TestParseAttributeSubstitution_InBody(t *testing.T) {
	input := `:attr1: Value1

Some text.

:attr2: Value2

More text with {attr2}.`
	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find text nodes containing the substituted value
	var found bool
	doc.Traverse(func(n *Node) {
		if n.Type == Text && strings.Contains(n.Content, "Value2") {
			found = true
		}
	})

	if !found {
		t.Error("Expected attribute assignment in body to work and substitution to replace {attr2} with 'Value2'")
	}
}
