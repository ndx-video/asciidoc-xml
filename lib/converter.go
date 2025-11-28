package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"
)

// ConvertOptions configures HTML conversion behavior
type ConvertOptions struct {
	UsePicoCSS bool // If true, includes PicoCSS styling
	
	Standalone bool // If true, generate full <html> doc. If false, generate only the <body> fragment.
	
	Title  string // Optional override for document title
	Author string // Optional override for document author
	
	XHTML         bool   // If true, outputs well-formed XHTML5
	PicoCSSPath   string // Path to PicoCSS (used for <link> tag if provided)
	PicoCSSContent string // PicoCSS content to embed inline (if PicoCSSPath is empty and UsePicoCSS is true)
}

// Metadata contains parsed document metadata
type Metadata struct {
	Title      string
	Author     string
	Attributes map[string]string // e.g. map[":toc": "true"]
}

// Result contains the converted HTML and extracted metadata
type Result struct {
	HTML string
	Meta Metadata
}

// ParseDocument converts AsciiDoc content to our generic DOM Node.
// This is used internally for XML conversion and direct DOM manipulation.
// For HTML conversion, use Convert() instead.
func ParseDocument(reader io.Reader) (*Node, error) {
	return Parse(reader)
}

// ToHTML converts an AST node to HTML string
func ToHTML(node *Node) string {
	var buf bytes.Buffer
	toHTML(node, &buf, false, 0)
	return buf.String()
}

// toHTML is the internal recursive function for HTML conversion
func toHTML(node *Node, buf *bytes.Buffer, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)

	switch node.Type {
	case Document:
		// Document has no wrapper, just output children
		for _, child := range node.Children {
			toHTML(child, buf, xhtml, indent)
		}

	case Section:
		level := 1
		if levelAttr := node.GetAttribute("level"); levelAttr != "" {
			fmt.Sscanf(levelAttr, "%d", &level)
		}
		hLevel := level + 1
		if hLevel > 6 {
			hLevel = 6
		}
		tagName := fmt.Sprintf("h%d", hLevel)
		
		attrs := ""
		if id := node.GetAttribute("id"); id != "" {
			attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
		}
		if role := node.GetAttribute("role"); role != "" {
			attrs += fmt.Sprintf(` class="%s"`, html.EscapeString(role))
		}
		
		// Section title (first text child or title attribute)
		titleText := node.GetAttribute("title")
		if titleText == "" && len(node.Children) > 0 && node.Children[0].Type == Text {
			titleText = node.Children[0].Content
		}
		
		if titleText != "" {
			fmt.Fprintf(buf, "%s<%s%s>%s</%s>\n", indentStr, tagName, attrs, html.EscapeString(titleText), tagName)
		}
		
		// Section content (skip first child if it was the title text)
		startIdx := 0
		if len(node.Children) > 0 && node.Children[0].Type == Text {
			startIdx = 1
		}
		for i := startIdx; i < len(node.Children); i++ {
			toHTML(node.Children[i], buf, xhtml, indent)
		}

	case Paragraph:
		attrs := ""
		if id := node.GetAttribute("id"); id != "" {
			attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
		}
		if role := node.GetAttribute("role"); role != "" {
			if attrs != "" {
				attrs += " "
			}
			attrs += fmt.Sprintf(`class="%s"`, html.EscapeString(role))
		}
		
		if attrs != "" {
			fmt.Fprintf(buf, "%s<p%s>", indentStr, attrs)
		} else {
			fmt.Fprintf(buf, "%s<p>", indentStr)
		}
		toHTMLInlineContent(node, buf, xhtml)
		buf.WriteString("</p>\n")

	case BlockMacro:
		if node.Name == "component" {
			componentName := node.GetAttribute("component-name")
			buf.WriteString(indentStr + "<cms-component")
			for k, v := range node.Attributes {
				if k != "component-name" {
					buf.WriteString(fmt.Sprintf(` %s="%s"`, k, html.EscapeString(v)))
				}
			}
			if componentName != "" {
				buf.WriteString(fmt.Sprintf(` component-name="%s"`, html.EscapeString(componentName)))
			}
			if len(node.Children) == 0 {
				if xhtml {
					buf.WriteString("/>\n")
				} else {
					buf.WriteString(">\n")
				}
			} else {
				buf.WriteString(">\n")
				for _, child := range node.Children {
					toHTML(child, buf, xhtml, indent+1)
				}
				buf.WriteString(indentStr + "</cms-component>\n")
			}
		} else if node.Name == "image" {
			src := node.GetAttribute("src")
			alt := node.GetAttribute("alt")
			buf.WriteString(indentStr + `<img style="display:block" src="` + html.EscapeString(src) + `"`)
			if alt != "" {
				buf.WriteString(` alt="` + html.EscapeString(alt) + `"`)
			} else {
				buf.WriteString(` alt=""`)
			}
			if xhtml {
				buf.WriteString("/>\n")
			} else {
				buf.WriteString(">\n")
			}
		} else {
			// Generic block macro
			buf.WriteString(indentStr + `<div class="macro macro-` + html.EscapeString(node.Name) + `"`)
			for k, v := range node.Attributes {
				buf.WriteString(fmt.Sprintf(` %s="%s"`, k, html.EscapeString(v)))
			}
			buf.WriteString(">\n")
			for _, child := range node.Children {
				toHTML(child, buf, xhtml, indent+1)
			}
			buf.WriteString(indentStr + "</div>\n")
		}

	case Text:
		buf.WriteString(html.EscapeString(node.Content))

	case List:
		style := node.GetAttribute("style")
		tagName := "ul"
		if style == "ordered" {
			tagName = "ol"
		} else if style == "labeled" {
			tagName = "dl"
		}
		
		attrs := ""
		if id := node.GetAttribute("id"); id != "" {
			attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
		}
		if role := node.GetAttribute("role"); role != "" {
			if attrs != "" {
				attrs += " "
			}
			attrs += fmt.Sprintf(`class="%s"`, html.EscapeString(role))
		}
		
		if attrs != "" {
			fmt.Fprintf(buf, "%s<%s%s>\n", indentStr, tagName, attrs)
		} else {
			fmt.Fprintf(buf, "%s<%s>\n", indentStr, tagName)
		}
		
		for _, item := range node.Children {
			if item.Type == ListItem {
				if style == "labeled" {
					// For labeled lists, item has term and description as children
					term := item.GetAttribute("term")
					if term != "" {
						fmt.Fprintf(buf, "%s    <dt>", indentStr)
						toHTMLInlineContent(item, buf, xhtml)
						buf.WriteString("</dt>\n")
					}
					// Description is the second child paragraph
					if len(item.Children) > 1 {
						fmt.Fprintf(buf, "%s    <dd>", indentStr)
						toHTML(item.Children[1], buf, xhtml, 0)
						buf.WriteString("</dd>\n")
					}
				} else {
					fmt.Fprintf(buf, "%s    <li>", indentStr)
					toHTMLInlineContent(item, buf, xhtml)
					buf.WriteString("</li>\n")
				}
			}
		}
		
		fmt.Fprintf(buf, "%s</%s>\n", indentStr, tagName)

	case CodeBlock:
		indentStr := strings.Repeat("    ", indent)
		if title := node.GetAttribute("title"); title != "" {
			fmt.Fprintf(buf, "%s<p class=\"code-title\">%s</p>\n", indentStr, html.EscapeString(title))
		}
		
		attrs := ""
		if lang := node.GetAttribute("language"); lang != "" {
			attrs += fmt.Sprintf(` class="language-%s" data-language="%s"`, html.EscapeString(lang), html.EscapeString(lang))
		}
		
		fmt.Fprintf(buf, "%s<pre><code%s>", indentStr, attrs)
		for _, child := range node.Children {
			if child.Type == Text {
				buf.WriteString(html.EscapeString(child.Content))
			}
		}
		buf.WriteString("</code></pre>\n")

	case LiteralBlock:
		indentStr := strings.Repeat("    ", indent)
		attrs := `class="literal-block"`
		if id := node.GetAttribute("id"); id != "" {
			attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
		}
		fmt.Fprintf(buf, "%s<pre%s>", indentStr, attrs)
		for _, child := range node.Children {
			if child.Type == Text {
				buf.WriteString(html.EscapeString(child.Content))
			}
		}
		buf.WriteString("</pre>\n")

	case Example:
		indentStr := strings.Repeat("    ", indent)
		attrs := `class="example"`
		if id := node.GetAttribute("id"); id != "" {
			attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
		}
		fmt.Fprintf(buf, "%s<div%s>\n", indentStr, attrs)
		if title := node.GetAttribute("title"); title != "" {
			fmt.Fprintf(buf, "%s    <p class=\"example-title\">%s</p>\n", indentStr, html.EscapeString(title))
		}
		for _, child := range node.Children {
			toHTML(child, buf, xhtml, indent+1)
		}
		fmt.Fprintf(buf, "%s</div>\n", indentStr)

	case Sidebar:
		indentStr := strings.Repeat("    ", indent)
		attrs := `class="sidebar"`
		if id := node.GetAttribute("id"); id != "" {
			attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
		}
		fmt.Fprintf(buf, "%s<aside%s>\n", indentStr, attrs)
		if title := node.GetAttribute("title"); title != "" {
			fmt.Fprintf(buf, "%s    <p class=\"sidebar-title\">%s</p>\n", indentStr, html.EscapeString(title))
		}
		for _, child := range node.Children {
			toHTML(child, buf, xhtml, indent+1)
		}
		fmt.Fprintf(buf, "%s</aside>\n", indentStr)

	case Quote:
		indentStr := strings.Repeat("    ", indent)
		attrs := `class="quote"`
		if id := node.GetAttribute("id"); id != "" {
			attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
		}
		fmt.Fprintf(buf, "%s<blockquote%s>\n", indentStr, attrs)
		for _, child := range node.Children {
			toHTML(child, buf, xhtml, indent+1)
		}
		if attribution := node.GetAttribute("attribution"); attribution != "" {
			fmt.Fprintf(buf, "%s    <footer><cite>%s</cite></footer>\n", indentStr, html.EscapeString(attribution))
		}
		if citation := node.GetAttribute("citation"); citation != "" {
			fmt.Fprintf(buf, "%s    <cite>%s</cite>\n", indentStr, html.EscapeString(citation))
		}
		fmt.Fprintf(buf, "%s</blockquote>\n", indentStr)

	case Table:
		indentStr := strings.Repeat("    ", indent)
		attrs := `class="table"`
		if id := node.GetAttribute("id"); id != "" {
			attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
		}
		fmt.Fprintf(buf, "%s<table%s>\n", indentStr, attrs)
		
		firstRow := true
		for _, child := range node.Children {
			if child.Type == TableRow {
				if firstRow {
					fmt.Fprintf(buf, "%s    <thead>\n", indentStr)
					fmt.Fprintf(buf, "%s        <tr>\n", indentStr)
					for _, cell := range child.Children {
						if cell.Type == TableCell {
							fmt.Fprintf(buf, "%s            <th>", indentStr)
							toHTMLInlineContent(cell, buf, xhtml)
							buf.WriteString("</th>\n")
						}
					}
					fmt.Fprintf(buf, "%s        </tr>\n", indentStr)
					fmt.Fprintf(buf, "%s    </thead>\n", indentStr)
					fmt.Fprintf(buf, "%s    <tbody>\n", indentStr)
					firstRow = false
				} else {
					fmt.Fprintf(buf, "%s        <tr>\n", indentStr)
					for _, cell := range child.Children {
						if cell.Type == TableCell {
							fmt.Fprintf(buf, "%s            <td>", indentStr)
							toHTMLInlineContent(cell, buf, xhtml)
							buf.WriteString("</td>\n")
						}
					}
					fmt.Fprintf(buf, "%s        </tr>\n", indentStr)
				}
			}
		}
		if !firstRow {
			fmt.Fprintf(buf, "%s    </tbody>\n", indentStr)
		}
		fmt.Fprintf(buf, "%s</table>\n", indentStr)

	case Admonition:
		indentStr := strings.Repeat("    ", indent)
		admType := node.GetAttribute("type")
		attrs := fmt.Sprintf(`class="admonition admonition-%s"`, html.EscapeString(admType))
		if id := node.GetAttribute("id"); id != "" {
			attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
		}
		fmt.Fprintf(buf, "%s<div%s>\n", indentStr, attrs)
		fmt.Fprintf(buf, "%s    <p class=\"admonition-title\">%s</p>\n", indentStr, html.EscapeString(strings.ToUpper(admType)))
		for _, child := range node.Children {
			toHTML(child, buf, xhtml, indent+1)
		}
		fmt.Fprintf(buf, "%s</div>\n", indentStr)

	case ThematicBreak:
		indentStr := strings.Repeat("    ", indent)
		if xhtml {
			buf.WriteString(indentStr + "<hr/>\n")
		} else {
			buf.WriteString(indentStr + "<hr>\n")
		}

	case PageBreak:
		indentStr := strings.Repeat("    ", indent)
		buf.WriteString(indentStr + `<div class="page-break"></div>` + "\n")

	case Bold:
		buf.WriteString("<strong>")
		toHTMLInlineContent(node, buf, xhtml)
		buf.WriteString("</strong>")

	case Italic:
		buf.WriteString("<em>")
		toHTMLInlineContent(node, buf, xhtml)
		buf.WriteString("</em>")

	case Monospace:
		buf.WriteString("<code>")
		toHTMLInlineContent(node, buf, xhtml)
		buf.WriteString("</code>")

	case Link:
		href := node.GetAttribute("href")
		attrs := fmt.Sprintf(`href="%s"`, html.EscapeString(href))
		if title := node.GetAttribute("title"); title != "" {
			attrs += fmt.Sprintf(` title="%s"`, html.EscapeString(title))
		}
		if class := node.GetAttribute("class"); class != "" {
			attrs += fmt.Sprintf(` class="%s"`, html.EscapeString(class))
		}
		if window := node.GetAttribute("window"); window != "" {
			attrs += fmt.Sprintf(` target="%s"`, html.EscapeString(window))
		}
		fmt.Fprintf(buf, "<a %s>", attrs)
		toHTMLInlineContent(node, buf, xhtml)
		buf.WriteString("</a>")

	case Passthrough:
		// Write content directly without escaping (for CMS-injected HTML)
		buf.WriteString(node.Content)

	default:
		// Unknown type, just output children
		for _, child := range node.Children {
			toHTML(child, buf, xhtml, indent)
		}
	}
}

// toHTMLInlineContent writes inline content (text and inline nodes)
func toHTMLInlineContent(node *Node, buf *bytes.Buffer, xhtml bool) {
	for _, child := range node.Children {
		if child.Type == Text {
			buf.WriteString(html.EscapeString(child.Content))
		} else {
			toHTML(child, buf, xhtml, 0)
		}
	}
}

// ConvertToXML converts AsciiDoc to XML string using the AST
func ConvertToXML(reader io.Reader) (string, error) {
	doc, err := ParseDocument(reader)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	xmlContent := ToXML(doc)
	buf.WriteString(xmlContent)

	return buf.String(), nil
}

// ToXML converts an AST node to XML string
func ToXML(node *Node) string {
	var buf bytes.Buffer
	toXML(node, &buf, 0)
	return buf.String()
}

// toXML is the internal recursive function for XML conversion
func toXML(node *Node, buf *bytes.Buffer, indentLevel int) {
	indent := strings.Repeat("  ", indentLevel)

	switch node.Type {
	case Document:
		buf.WriteString("<document")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>")
		} else {
			buf.WriteString(">\n")
			for _, child := range node.Children {
				toXML(child, buf, indentLevel+1)
			}
			buf.WriteString(indent + "</document>\n")
		}

	case Section:
		buf.WriteString(indent + "<section")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">\n")
			for _, child := range node.Children {
				toXML(child, buf, indentLevel+1)
			}
			buf.WriteString(indent + "</section>\n")
		}

	case Paragraph:
		buf.WriteString(indent + "<paragraph")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">")
			toXMLInlineContent(node, buf)
			buf.WriteString("</paragraph>\n")
		}

	case BlockMacro:
		buf.WriteString(indent + `<macro type="block" name="` + escapeXML(node.Name) + `"`)
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">\n")
			for _, child := range node.Children {
				toXML(child, buf, indentLevel+1)
			}
			buf.WriteString(indent + "</macro>\n")
		}

	case InlineMacro:
		buf.WriteString(`<macro type="inline" name="` + escapeXML(node.Name) + `"`)
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>")
		} else {
			buf.WriteString(">")
			toXMLInlineContent(node, buf)
			buf.WriteString("</macro>")
		}

	case Text:
		buf.WriteString(escapeXML(node.Content))

	case List:
		buf.WriteString(indent + "<list")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">\n")
			for _, child := range node.Children {
				toXML(child, buf, indentLevel+1)
			}
			buf.WriteString(indent + "</list>\n")
		}

	case ListItem:
		buf.WriteString(indent + "<listitem")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">")
			toXMLInlineContent(node, buf)
			buf.WriteString("</listitem>\n")
		}

	case CodeBlock:
		buf.WriteString(indent + "<codeblock")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">")
			toXMLInlineContent(node, buf)
			buf.WriteString("</codeblock>\n")
		}

	case LiteralBlock:
		buf.WriteString(indent + "<literalblock")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">")
			toXMLInlineContent(node, buf)
			buf.WriteString("</literalblock>\n")
		}

	case Example:
		buf.WriteString(indent + "<example")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">\n")
			for _, child := range node.Children {
				toXML(child, buf, indentLevel+1)
			}
			buf.WriteString(indent + "</example>\n")
		}

	case Sidebar:
		buf.WriteString(indent + "<sidebar")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">\n")
			for _, child := range node.Children {
				toXML(child, buf, indentLevel+1)
			}
			buf.WriteString(indent + "</sidebar>\n")
		}

	case Quote:
		buf.WriteString(indent + "<quote")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">\n")
			for _, child := range node.Children {
				toXML(child, buf, indentLevel+1)
			}
			buf.WriteString(indent + "</quote>\n")
		}

	case Table:
		buf.WriteString(indent + "<table")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">\n")
			for _, child := range node.Children {
				toXML(child, buf, indentLevel+1)
			}
			buf.WriteString(indent + "</table>\n")
		}

	case TableRow:
		buf.WriteString(indent + "<row")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">\n")
			for _, child := range node.Children {
				toXML(child, buf, indentLevel+1)
			}
			buf.WriteString(indent + "</row>\n")
		}

	case TableCell:
		buf.WriteString(indent + "<cell")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">")
			toXMLInlineContent(node, buf)
			buf.WriteString("</cell>\n")
		}

	case Admonition:
		buf.WriteString(indent + "<admonition")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">\n")
			for _, child := range node.Children {
				toXML(child, buf, indentLevel+1)
			}
			buf.WriteString(indent + "</admonition>\n")
		}

	case ThematicBreak:
		buf.WriteString(indent + "<thematicbreak/>\n")

	case PageBreak:
		buf.WriteString(indent + "<pagebreak/>\n")

	case Bold:
		buf.WriteString("<strong>")
		toXMLInlineContent(node, buf)
		buf.WriteString("</strong>")

	case Italic:
		buf.WriteString("<emphasis>")
		toXMLInlineContent(node, buf)
		buf.WriteString("</emphasis>")

	case Monospace:
		buf.WriteString("<monospace>")
		toXMLInlineContent(node, buf)
		buf.WriteString("</monospace>")

	case Link:
		buf.WriteString("<link")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>")
		} else {
			buf.WriteString(">")
			toXMLInlineContent(node, buf)
			buf.WriteString("</link>")
		}

	case Passthrough:
		// Wrap content in CDATA for XML
		buf.WriteString("<passthrough><![CDATA[")
		buf.WriteString(node.Content)
		buf.WriteString("]]></passthrough>")

	default:
		// Unknown type
		for _, child := range node.Children {
			toXML(child, buf, indentLevel)
		}
	}
}

// toXMLInlineContent writes inline content for XML
func toXMLInlineContent(node *Node, buf *bytes.Buffer) {
	for _, child := range node.Children {
		if child.Type == Text {
			buf.WriteString(escapeXML(child.Content))
		} else {
			toXML(child, buf, 0)
		}
	}
}

// escapeXML escapes XML special characters
func escapeXML(s string) string {
	var result strings.Builder
	for _, c := range s {
		switch c {
		case '<':
			result.WriteString("&lt;")
		case '>':
			result.WriteString("&gt;")
		case '&':
			result.WriteString("&amp;")
		case '"':
			result.WriteString("&quot;")
		case '\'':
			result.WriteString("&apos;")
		default:
			result.WriteRune(c)
		}
	}
	return result.String()
}

// extractMetadata extracts metadata from a parsed document
func extractMetadata(doc *Node) Metadata {
	meta := Metadata{
		Attributes: make(map[string]string),
	}

	if doc.Type == Document {
		// Extract title from document attributes
		meta.Title = doc.GetAttribute("title")
		meta.Author = doc.GetAttribute("author")

		// Extract all attributes (including custom ones with : prefix)
		for k, v := range doc.Attributes {
			if strings.HasPrefix(k, ":") {
				meta.Attributes[k] = v
			}
		}
	}

	return meta
}

// Convert converts AsciiDoc to HTML using the new options-based API.
// This is the recommended entry point for HTML conversion with full control over options.
//
// Example:
//
//	opts := ConvertOptions{
//	    Standalone: true,
//	    UsePicoCSS: true,
//	    Title: "My Document",
//	}
//	result, err := Convert(reader, opts)
//	fmt.Println(result.HTML)
//	fmt.Println(result.Meta.Title)
func Convert(reader io.Reader, opts ConvertOptions) (Result, error) {
	doc, err := ParseDocument(reader)
	if err != nil {
		return Result{}, err
	}

	// Extract metadata first
	meta := extractMetadata(doc)
	
	// Override title/author if provided in options
	if opts.Title != "" {
		meta.Title = opts.Title
	}
	if opts.Author != "" {
		meta.Author = opts.Author
	}

	var buf bytes.Buffer

	// If standalone, write full HTML document structure
	if opts.Standalone {
		if opts.XHTML {
			buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
			buf.WriteString(`<!DOCTYPE html>` + "\n")
		} else {
			buf.WriteString(`<!DOCTYPE html>` + "\n")
		}

		// Get lang from document attributes
		lang := doc.GetAttribute(":lang")
		if lang == "" {
			lang = "en"
		}

		if opts.XHTML {
			fmt.Fprintf(&buf, `<html xmlns="http://www.w3.org/1999/xhtml" lang="%s">`+"\n", html.EscapeString(lang))
		} else {
			fmt.Fprintf(&buf, `<html lang="%s">`+"\n", html.EscapeString(lang))
		}

		// Head section
		buf.WriteString("  <head>\n")
		if opts.XHTML {
			buf.WriteString("    <meta charset=\"UTF-8\"/>\n")
		} else {
			buf.WriteString("    <meta charset=\"UTF-8\">\n")
		}

		// Add PicoCSS if enabled
		if opts.UsePicoCSS {
			if opts.PicoCSSContent != "" {
				// Embed CSS inline
				buf.WriteString("    <style>\n")
				buf.WriteString(opts.PicoCSSContent)
				buf.WriteString("\n    </style>\n")
			} else if opts.PicoCSSPath != "" {
				// Use link tag
				if opts.XHTML {
					fmt.Fprintf(&buf, `    <link rel="stylesheet" href="%s"/>`+"\n", html.EscapeString(opts.PicoCSSPath))
				} else {
					fmt.Fprintf(&buf, `    <link rel="stylesheet" href="%s">`+"\n", html.EscapeString(opts.PicoCSSPath))
				}
			}
		}

		// Title in head (use override if provided)
		title := meta.Title
		if opts.Title != "" {
			title = opts.Title
		}
		if title != "" {
			fmt.Fprintf(&buf, "    <title>%s</title>\n", html.EscapeString(title))
		}
		buf.WriteString("  </head>\n")

		// Body section
		buf.WriteString("  <body>\n")
	}

	if opts.Standalone {
		// Write document header in standalone mode
		title := meta.Title
		if opts.Title != "" {
			title = opts.Title
		}
		author := meta.Author
		if opts.Author != "" {
			author = opts.Author
		}
		
		if title != "" || author != "" {
			buf.WriteString("    <header>\n")
			if title != "" {
				fmt.Fprintf(&buf, "      <h1>%s</h1>\n", html.EscapeString(title))
			}
			if author != "" {
				buf.WriteString("      <address class=\"authors\">\n")
				fmt.Fprintf(&buf, "        <p><span class=\"author-name\">%s</span>", html.EscapeString(author))
				if email := doc.GetAttribute("email"); email != "" {
					fmt.Fprintf(&buf, ` <a href="mailto:%s" class="author-email">%s</a>`, html.EscapeString(email), html.EscapeString(email))
				}
				buf.WriteString("</p>\n")
				buf.WriteString("      </address>\n")
			}
			buf.WriteString("    </header>\n")
		}
		buf.WriteString("    <main>\n")
	}

	// Write content
	htmlContent := ToHTML(doc)
	if opts.Standalone {
		// Indent the content
		lines := strings.Split(htmlContent, "\n")
		for _, line := range lines {
			if line != "" {
				buf.WriteString("      " + line + "\n")
			} else {
				buf.WriteString("\n")
			}
		}
	} else {
		buf.WriteString(htmlContent)
	}

	if opts.Standalone {
		buf.WriteString("    </main>\n")
		buf.WriteString("  </body>\n")
		buf.WriteString("</html>\n")
	}

	return Result{
		HTML: buf.String(),
		Meta: meta,
	}, nil
}

// ConvertToHTML converts AsciiDoc to HTML5 string
//
// Deprecated: Use ConvertToHTMLWithOptions instead for better flexibility.
// This function is maintained for backward compatibility only.
//
// If xhtml is true, outputs well-formed XHTML5
// If usePicoCSS is true, includes PicoCSS styling
// picoCSSPath is used for <link> tag (if empty and usePicoCSS is true, picoCSSContent is embedded inline)
// picoCSSContent is embedded as <style> tag when usePicoCSS is true and picoCSSPath is empty
func ConvertToHTML(reader io.Reader, xhtml bool, usePicoCSS bool, picoCSSPath string, picoCSSContent string) (string, error) {
	opts := ConvertOptions{
		UsePicoCSS:    usePicoCSS,
		Standalone:    true,
		XHTML:         xhtml,
		PicoCSSPath:   picoCSSPath,
		PicoCSSContent: picoCSSContent,
	}
	
	result, err := Convert(reader, opts)
	if err != nil {
		return "", err
	}
	
	return result.HTML, nil
}

func getTextContent(node *Node) string {
	var buf bytes.Buffer
	for _, child := range node.Children {
		if child.Type == Text {
			buf.WriteString(child.Content)
		} else {
			buf.WriteString(getTextContent(child))
		}
	}
	return buf.String()
}


// ConvertMarkdownToAsciiDoc converts Markdown content to AsciiDoc format
func ConvertMarkdownToAsciiDoc(reader io.Reader) (string, error) {
	scanner := bufio.NewScanner(reader)
	var result bytes.Buffer
	var inCodeBlock bool
	var codeBlockLang string
	var inTable bool
	var codeBlockLines []string

	// Regex patterns
	headerRegex := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	codeBlockStartRegex := regexp.MustCompile("^```(\\w*)$")
	codeBlockEndRegex := regexp.MustCompile("^```$")
	tableRowRegex := regexp.MustCompile(`^\|(.+)\|$`)
	horizontalRuleRegex := regexp.MustCompile(`^[-*_]{3,}$`)
	blockquoteRegex := regexp.MustCompile(`^>\s*(.*)$`)
	orderedListRegex := regexp.MustCompile(`^(\d+)\.\s+(.+)$`)
	unorderedListRegex := regexp.MustCompile(`^[-*+]\s+(.+)$`)

	for scanner.Scan() {
		line := scanner.Text()

		// Handle code blocks - check end first if we're in a code block
		if inCodeBlock {
			if codeBlockEndRegex.MatchString(line) {
				// End of code block
				if codeBlockLang != "" {
					result.WriteString(fmt.Sprintf("[source,%s]\n", codeBlockLang))
				}
				result.WriteString("----\n")
				for _, codeLine := range codeBlockLines {
					result.WriteString(codeLine + "\n")
				}
				result.WriteString("----\n")
				inCodeBlock = false
				codeBlockLang = ""
				codeBlockLines = []string{}
				continue
			}
			codeBlockLines = append(codeBlockLines, line)
			continue
		}

		// Check for code block start (only when not already in a code block)
		if codeBlockStartRegex.MatchString(line) {
			matches := codeBlockStartRegex.FindStringSubmatch(line)
			codeBlockLang = matches[1]
			inCodeBlock = true
			codeBlockLines = []string{}
			continue
		}

		// Handle headers
		if matches := headerRegex.FindStringSubmatch(line); matches != nil {
			level := len(matches[1])
			title := strings.TrimSpace(matches[2])
			// Convert # to =, ## to ==, etc.
			equals := strings.Repeat("=", level)
			result.WriteString(fmt.Sprintf("%s %s\n", equals, title))
			continue
		}

		// Handle horizontal rules
		if horizontalRuleRegex.MatchString(line) {
			result.WriteString("'''\n")
			continue
		}

		// Handle blockquotes
		if matches := blockquoteRegex.FindStringSubmatch(line); matches != nil {
			result.WriteString(fmt.Sprintf("[quote]\n____\n%s\n____\n", matches[1]))
			continue
		}

		// Handle tables
		if tableRowRegex.MatchString(line) {
			if !inTable {
				result.WriteString("|===\n")
				inTable = true
			}
			// Convert table row - keep as is for now, AsciiDoc uses similar syntax
			result.WriteString(line + "\n")
			continue
		} else if inTable {
			// End table if we hit a non-table line
			result.WriteString("|===\n")
			inTable = false
		}

		// Handle ordered lists
		if matches := orderedListRegex.FindStringSubmatch(line); matches != nil {
			content := matches[2]
			// Convert inline formatting in list items
			content = convertInlineMarkdown(content)
			result.WriteString(fmt.Sprintf(". %s\n", content))
			continue
		}

		// Handle unordered lists
		if matches := unorderedListRegex.FindStringSubmatch(line); matches != nil {
			content := matches[1]
			// Convert inline formatting in list items
			content = convertInlineMarkdown(content)
			result.WriteString(fmt.Sprintf("* %s\n", content))
			continue
		}

		// Handle regular lines
		if strings.TrimSpace(line) == "" {
			result.WriteString("\n")
			continue
		}

		// Convert inline Markdown in regular text
		convertedLine := convertInlineMarkdown(line)
		result.WriteString(convertedLine + "\n")
	}

	// Close table if still open
	if inTable {
		result.WriteString("|===\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading markdown: %w", err)
	}

	return result.String(), nil
}

// convertInlineMarkdown converts inline Markdown formatting to AsciiDoc
func convertInlineMarkdown(text string) string {
	// Convert images first (before links, as images contain link syntax)
	imageRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	text = imageRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := imageRegex.FindStringSubmatch(match)
		alt := matches[1]
		src := matches[2]
		return fmt.Sprintf("image::%s[%s]", src, alt)
	})

	// Convert links [text](url) to link:url[text]
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = linkRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := linkRegex.FindStringSubmatch(match)
		linkText := matches[1]
		url := matches[2]
		return fmt.Sprintf("link:%s[%s]", url, linkText)
	})

	// Convert bold **text** to *text* using placeholder to avoid conflicts with italic
	boldDoubleStarRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	placeholderMap := make(map[string]string)
	placeholderCounter := 0
	
	// Replace bold with placeholder first
	text = boldDoubleStarRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholderCounter++
		placeholder := fmt.Sprintf("__BOLD_PLACEHOLDER_%d__", placeholderCounter)
		matches := boldDoubleStarRegex.FindStringSubmatch(match)
		placeholderMap[placeholder] = "*" + matches[1] + "*"
		return placeholder
	})
	
	// Convert bold __text__ to *text* (also use placeholder)
	boldDoubleUnderscoreRegex := regexp.MustCompile(`__([^_]+)__`)
	text = boldDoubleUnderscoreRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholderCounter++
		placeholder := fmt.Sprintf("__BOLD_PLACEHOLDER_%d__", placeholderCounter)
		matches := boldDoubleUnderscoreRegex.FindStringSubmatch(match)
		placeholderMap[placeholder] = "*" + matches[1] + "*"
		return placeholder
	})

	// Convert italic *text* to _text_ (now safe since bold is in placeholders)
	italicAsteriskRegex := regexp.MustCompile(`\*([^*\n]+)\*`)
	text = italicAsteriskRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := italicAsteriskRegex.FindStringSubmatch(match)
		if len(matches) > 1 {
			return "_" + matches[1] + "_"
		}
		return match
	})

	// Convert italic _text_ to _text_ (already correct format, no change needed)
	// Note: Markdown uses _text_ for italic, AsciiDoc also uses _text_ for italic, so no conversion needed

	// Restore bold placeholders
	for placeholder, replacement := range placeholderMap {
		text = strings.ReplaceAll(text, placeholder, replacement)
	}

	// Inline code `code` is already in correct format for AsciiDoc
	// No conversion needed

	return text
}

