package lib

import (
	"bytes"
	"fmt"
	"strings"
)

type NodeType int

const (
	ElementNode NodeType = iota
	TextNode
	CommentNode
)

type Node struct {
	Type       NodeType
	Data       string            // Tag name (for Element) or content (for Text)
	Attributes map[string]string // Key-value pairs for attributes
	Children   []*Node
	Parent     *Node
}

func NewElementNode(tagName string) *Node {
	return &Node{
		Type:       ElementNode,
		Data:       tagName,
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0),
	}
}

func NewTextNode(text string) *Node {
	return &Node{
		Type:     TextNode,
		Data:     text,
		Children: make([]*Node, 0),
	}
}

func (n *Node) AddChild(child *Node) {
	child.Parent = n
	n.Children = append(n.Children, child)
}

func (n *Node) SetAttribute(key, value string) {
	if n.Attributes == nil {
		n.Attributes = make(map[string]string)
	}
	n.Attributes[key] = value
}

func (n *Node) GetAttribute(key string) string {
	if n.Attributes == nil {
		return ""
	}
	return n.Attributes[key]
}

// Traverse traverses the DOM tree depth-first
func (n *Node) Traverse(visit func(*Node)) {
	visit(n)
	for _, child := range n.Children {
		child.Traverse(visit)
	}
}

// FindElementsByTag performs a recursive depth-first search and returns all nodes
// where node.Data == tagName, starting from this node and including it.
func (n *Node) FindElementsByTag(tagName string) []*Node {
	var results []*Node
	
	// Check if current node matches
	if n.Data == tagName {
		results = append(results, n)
	}
	
	// Recursively search children
	for _, child := range n.Children {
		childResults := child.FindElementsByTag(tagName)
		results = append(results, childResults...)
	}
	
	return results
}

// ToXML generates an XML string from the node
func (n *Node) ToXML() (string, error) {
	var buf bytes.Buffer
	if err := n.writeXML(&buf, 0); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (n *Node) writeXML(buf *bytes.Buffer, indentLevel int) error {
	indent := strings.Repeat("  ", indentLevel)

	switch n.Type {
	case TextNode:
		// For text nodes, we might want to trim space if it's just whitespace between blocks,
		// but generally we want to preserve it if it's inline content.
		// Since we are building a clean DOM, text nodes should be meaningful.
		// We'll just escape special characters.
		buf.WriteString(escapeXML(n.Data))

	case ElementNode:
		// If it's the root document, don't write indentation for it (optional)
		// but typically we want the root element.
		// If it's inline element (like strong, emphasis), we shouldn't add newlines/indent around it
		// unless it's block level.
		// For simplicity in this "bloat-free" version, let's assume pretty printing only for block elements.
		// But deciding what is block vs inline is context dependent.
		// A simple heuristic: if it has text children mixed with elements, it's mixed content (inline-ish).
		// If it only has element children, it's block content.
		
		isMixed := hasMixedContent(n)
		
		if !isMixed && n.Parent != nil { // Only indent if not mixed content
			buf.WriteString("\n" + indent)
		}

		buf.WriteString("<" + n.Data)
		
		// Attributes
		for k, v := range n.Attributes {
			buf.WriteString(fmt.Sprintf(" %s=\"%s\"", k, escapeXML(v)))
		}

		if len(n.Children) == 0 {
			buf.WriteString("/>")
		} else {
			buf.WriteString(">")
			for _, child := range n.Children {
				var err error
				if isMixed {
					// No indent for children of mixed content
					err = child.writeXML(buf, 0) // 0 or effectively ignored for indentation logic if we handled it perfectly
				} else {
					err = child.writeXML(buf, indentLevel+1)
				}
				if err != nil {
					return err
				}
			}
			if !isMixed {
				buf.WriteString("\n" + indent)
			}
			buf.WriteString("</" + n.Data + ">")
		}
	}
	return nil
}

func hasMixedContent(n *Node) bool {
	for _, child := range n.Children {
		if child.Type == TextNode {
			return true
		}
	}
	return false
}

func escapeXML(s string) string {
	var buf bytes.Buffer
	if err := xmlEscape(&buf, s); err != nil {
		return s // Fallback
	}
	return buf.String()
}

// xmlEscape is a simple helper to escape XML characters
func xmlEscape(w *bytes.Buffer, s string) error {
	for _, c := range s {
		switch c {
		case '<':
			w.WriteString("&lt;")
		case '>':
			w.WriteString("&gt;")
		case '&':
			w.WriteString("&amp;")
		case '"':
			w.WriteString("&quot;")
		case '\'':
			w.WriteString("&apos;")
		default:
			w.WriteRune(c)
		}
	}
	return nil
}

