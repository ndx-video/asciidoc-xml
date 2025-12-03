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
		// Check if first child is a preamble (Paragraph with role="preamble")
		hasPreamble := false
		if len(node.Children) > 0 {
			firstChild := node.Children[0]
			if firstChild.Type == Paragraph && (firstChild.GetAttribute("role") == "preamble" || firstChild.GetAttribute("data-role") == "preamble") {
				hasPreamble = true
				// Output preamble wrapper
				indentStr := strings.Repeat("    ", indent)
				var attrParts []string
				attrParts = append(attrParts, `data-role="preamble"`)
				otherAttrs := buildHTMLAttributes(firstChild, []string{"role", "id"})
				if otherAttrs != "" {
					attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
				}
				attrs := buildAttrsString(attrParts...)
				fmt.Fprintf(buf, "%s<div%s>\n", indentStr, attrs)
				// Output preamble content (children of the paragraph)
				for _, grandchild := range firstChild.Children {
					toHTML(grandchild, buf, xhtml, indent+1)
				}
				fmt.Fprintf(buf, "%s</div>\n", indentStr)
			}
		}
		
		// Output remaining children (skip first if it was preamble)
		startIdx := 0
		if hasPreamble {
			startIdx = 1
		}
		for i := startIdx; i < len(node.Children); i++ {
			toHTML(node.Children[i], buf, xhtml, indent)
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
		
		var attrParts []string
		// Handle appendix and discrete attributes
		if appendix := node.GetAttribute("appendix"); appendix != "" {
			attrParts = append(attrParts, fmt.Sprintf(`data-asciidoc-appendix="%s"`, html.EscapeString(appendix)))
		}
		if discrete := node.GetAttribute("discrete"); discrete != "" {
			attrParts = append(attrParts, fmt.Sprintf(`data-asciidoc-discrete="%s"`, html.EscapeString(discrete)))
		}
		// Add other attributes (role for ARIA, others as data-asciidoc-*)
		// Exclude id, level, appendix, discrete, title, and marker
		otherAttrs := buildHTMLAttributes(node, []string{"id", "level", "appendix", "discrete", "title", "marker"})
		if otherAttrs != "" {
			attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
		}
		attrs := buildAttrsString(attrParts...)
		
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
		// Skip preamble paragraphs - they're handled in Document case
		if node.GetAttribute("role") == "preamble" || node.GetAttribute("data-role") == "preamble" {
			// Just output children directly without <p> wrapper
			for _, child := range node.Children {
				toHTML(child, buf, xhtml, indent)
			}
			return
		}
		
		var attrParts []string
		if id := node.GetAttribute("id"); id != "" {
			attrParts = append(attrParts, fmt.Sprintf(`id="%s"`, html.EscapeString(id)))
		}
		otherAttrs := buildHTMLAttributes(node, []string{"id"})
		if otherAttrs != "" {
			attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
		}
		attrs := buildAttrsString(attrParts...)
		
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
			attrs := buildHTMLAttributes(node, []string{"src", "alt"})
			buf.WriteString(attrs)
			if xhtml {
				buf.WriteString("/>\n")
			} else {
				buf.WriteString(">\n")
			}
		} else if node.Name == "anchor" {
			// Block anchor
			id := node.GetAttribute("id")
			if id == "" {
				id = node.GetAttribute("target")
			}
			if id != "" {
				fmt.Fprintf(buf, "%s<a id=\"%s\"></a>\n", indentStr, html.EscapeString(id))
			}
		} else if node.Name == "include" {
			if len(node.Children) > 0 {
				// Included content - render children
				for _, child := range node.Children {
					toHTML(child, buf, xhtml, indent)
				}
			} else {
				// Placeholder
				file := node.GetAttribute("src")
				if file == "" {
					file = node.GetAttribute("target")
				}
				var attrParts []string
				attrParts = append(attrParts, `data-role="include"`)
				if file != "" {
					attrParts = append(attrParts, fmt.Sprintf(`data-asciidoc-file="%s"`, html.EscapeString(file)))
				}
				otherAttrs := buildHTMLAttributes(node, []string{"src", "target"})
				if otherAttrs != "" {
					attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
				}
				attrs := buildAttrsString(attrParts...)
				fmt.Fprintf(buf, "%s<div%s>[Include: %s]</div>\n", indentStr, attrs, html.EscapeString(file))
			}
		} else if node.Name == "toc" {
			// TOC generation - collect sections from document
			// For now, output placeholder - full TOC generation would require document traversal
			var attrParts []string
			attrParts = append(attrParts, `data-role="toc"`)
			levels := node.GetAttribute("levels")
			if levels != "" {
				attrParts = append(attrParts, fmt.Sprintf(`data-asciidoc-levels="%s"`, html.EscapeString(levels)))
			}
			otherAttrs := buildHTMLAttributes(node, []string{"levels"})
			if otherAttrs != "" {
				attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
			}
			attrs := buildAttrsString(attrParts...)
			fmt.Fprintf(buf, "%s<nav%s>\n", indentStr, attrs)
			fmt.Fprintf(buf, "%s    <!-- TOC would be generated here -->\n", indentStr)
			fmt.Fprintf(buf, "%s</nav>\n", indentStr)
		} else if node.Name == "video" {
			src := node.GetAttribute("src")
			if src == "" {
				src = node.GetAttribute("target")
			}
			var attrParts []string
			attrParts = append(attrParts, fmt.Sprintf(`src="%s"`, html.EscapeString(src)))
			if controls := node.GetAttribute("controls"); controls != "" && controls != "false" {
				attrParts = append(attrParts, "controls")
			}
			if autoplay := node.GetAttribute("autoplay"); autoplay == "true" {
				attrParts = append(attrParts, "autoplay")
			}
			if loop := node.GetAttribute("loop"); loop == "true" {
				attrParts = append(attrParts, "loop")
			}
			if poster := node.GetAttribute("poster"); poster != "" {
				attrParts = append(attrParts, fmt.Sprintf(`poster="%s"`, html.EscapeString(poster)))
			}
			if width := node.GetAttribute("width"); width != "" {
				attrParts = append(attrParts, fmt.Sprintf(`width="%s"`, html.EscapeString(width)))
			}
			if height := node.GetAttribute("height"); height != "" {
				attrParts = append(attrParts, fmt.Sprintf(`height="%s"`, html.EscapeString(height)))
			}
			otherAttrs := buildHTMLAttributes(node, []string{"src", "target", "controls", "autoplay", "loop", "poster", "width", "height"})
			if otherAttrs != "" {
				attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
			}
			attrs := buildAttrsString(attrParts...)
			if xhtml {
				fmt.Fprintf(buf, "%s<video%s/>\n", indentStr, attrs)
			} else {
				fmt.Fprintf(buf, "%s<video%s></video>\n", indentStr, attrs)
			}
		} else if node.Name == "audio" {
			src := node.GetAttribute("src")
			if src == "" {
				src = node.GetAttribute("target")
			}
			var attrParts []string
			attrParts = append(attrParts, fmt.Sprintf(`src="%s"`, html.EscapeString(src)))
			if controls := node.GetAttribute("controls"); controls != "" && controls != "false" {
				attrParts = append(attrParts, "controls")
			}
			if autoplay := node.GetAttribute("autoplay"); autoplay == "true" {
				attrParts = append(attrParts, "autoplay")
			}
			if loop := node.GetAttribute("loop"); loop == "true" {
				attrParts = append(attrParts, "loop")
			}
			otherAttrs := buildHTMLAttributes(node, []string{"src", "target", "controls", "autoplay", "loop"})
			if otherAttrs != "" {
				attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
			}
			attrs := buildAttrsString(attrParts...)
			if xhtml {
				fmt.Fprintf(buf, "%s<audio%s/>\n", indentStr, attrs)
			} else {
				fmt.Fprintf(buf, "%s<audio%s></audio>\n", indentStr, attrs)
			}
		} else {
			// Generic block macro
			var attrParts []string
			attrParts = append(attrParts, `data-role="macro"`)
			attrParts = append(attrParts, fmt.Sprintf(`data-asciidoc-macro="%s"`, html.EscapeString(node.Name)))
			otherAttrs := buildHTMLAttributes(node, []string{})
			if otherAttrs != "" {
				attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
			}
			attrs := buildAttrsString(attrParts...)
			fmt.Fprintf(buf, "%s<div%s>\n", indentStr, attrs)
			for _, child := range node.Children {
				toHTML(child, buf, xhtml, indent+1)
			}
			fmt.Fprintf(buf, "%s</div>\n", indentStr)
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
		
		var attrParts []string
		if id := node.GetAttribute("id"); id != "" {
			attrParts = append(attrParts, fmt.Sprintf(`id="%s"`, html.EscapeString(id)))
		}
		otherAttrs := buildHTMLAttributes(node, []string{"id", "style"})
		if otherAttrs != "" {
			attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
		}
		attrs := buildAttrsString(attrParts...)
		
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
						var ddAttrParts []string
						if id := item.GetAttribute("id"); id != "" {
							ddAttrParts = append(ddAttrParts, fmt.Sprintf(`id="%s"`, html.EscapeString(id)))
						}
						otherAttrs := buildHTMLAttributes(item, []string{"id", "term"})
						if otherAttrs != "" {
							ddAttrParts = append(ddAttrParts, strings.TrimSpace(otherAttrs))
						}
						ddAttrs := buildAttrsString(ddAttrParts...)
						if ddAttrs != "" {
							fmt.Fprintf(buf, "%s    <dd%s>", indentStr, ddAttrs)
						} else {
							fmt.Fprintf(buf, "%s    <dd>", indentStr)
						}
						toHTML(item.Children[1], buf, xhtml, 0)
						buf.WriteString("</dd>\n")
					}
				} else {
					var liAttrParts []string
					if id := item.GetAttribute("id"); id != "" {
						liAttrParts = append(liAttrParts, fmt.Sprintf(`id="%s"`, html.EscapeString(id)))
					}
					// Handle callout attribute
					if callout := item.GetAttribute("callout"); callout != "" {
						liAttrParts = append(liAttrParts, fmt.Sprintf(`data-asciidoc-callout="%s"`, html.EscapeString(callout)))
					}
					otherAttrs := buildHTMLAttributes(item, []string{"id", "callout", "term"})
					if otherAttrs != "" {
						liAttrParts = append(liAttrParts, strings.TrimSpace(otherAttrs))
					}
					liAttrs := buildAttrsString(liAttrParts...)
					
					if liAttrs != "" {
						fmt.Fprintf(buf, "%s    <li%s>", indentStr, liAttrs)
					} else {
						fmt.Fprintf(buf, "%s    <li>", indentStr)
					}
					// Add callout marker if present
					if callout := item.GetAttribute("callout"); callout != "" {
						fmt.Fprintf(buf, `<span data-role="callout-marker">%s</span> `, html.EscapeString(callout))
					}
					toHTMLInlineContent(item, buf, xhtml)
					buf.WriteString("</li>\n")
				}
			}
		}
		
		fmt.Fprintf(buf, "%s</%s>\n", indentStr, tagName)

	case CodeBlock:
		indentStr := strings.Repeat("    ", indent)
		// Check for mermaid role
		role := node.GetAttribute("role")
		if role == "mermaid" {
			// Output as cms-mermaid web component, preserve whitespace
			fmt.Fprintf(buf, "%s<cms-mermaid>", indentStr)
			for _, child := range node.Children {
				if child.Type == Text {
					buf.WriteString(child.Content) // Don't escape for mermaid
				}
			}
			buf.WriteString("</cms-mermaid>\n")
		} else {
			// Regular code block
			if title := node.GetAttribute("title"); title != "" {
				fmt.Fprintf(buf, "%s<p data-role=\"code-title\">%s</p>\n", indentStr, html.EscapeString(title))
			}
			
			var attrParts []string
			if lang := node.GetAttribute("language"); lang != "" {
				attrParts = append(attrParts, fmt.Sprintf(`data-asciidoc-language="%s"`, html.EscapeString(lang)))
			}
			otherAttrs := buildHTMLAttributes(node, []string{"language", "title", "id", "role"})
			if otherAttrs != "" {
				attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
			}
			attrs := buildAttrsString(attrParts...)
			
			if attrs != "" {
				fmt.Fprintf(buf, "%s<pre><code%s>", indentStr, attrs)
			} else {
				fmt.Fprintf(buf, "%s<pre><code>", indentStr)
			}
			for _, child := range node.Children {
				if child.Type == Text {
					buf.WriteString(html.EscapeString(child.Content))
				}
			}
			buf.WriteString("</code></pre>\n")
		}

	case LiteralBlock:
		indentStr := strings.Repeat("    ", indent)
		// Check for mermaid role
		role := node.GetAttribute("role")
		if role == "mermaid" {
			// Output as cms-mermaid web component, preserve whitespace
			fmt.Fprintf(buf, "%s<cms-mermaid>", indentStr)
			for _, child := range node.Children {
				if child.Type == Text {
					buf.WriteString(child.Content) // Don't escape for mermaid
				}
			}
			buf.WriteString("</cms-mermaid>\n")
		} else {
			// Regular literal block
			var attrParts []string
			attrParts = append(attrParts, `data-role="literal-block"`)
			if id := node.GetAttribute("id"); id != "" {
				attrParts = append(attrParts, fmt.Sprintf(`id="%s"`, html.EscapeString(id)))
			}
			otherAttrs := buildHTMLAttributes(node, []string{"id"})
			if otherAttrs != "" {
				attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
			}
			attrs := buildAttrsString(attrParts...)
			fmt.Fprintf(buf, "%s<pre%s>", indentStr, attrs)
			for _, child := range node.Children {
				if child.Type == Text {
					buf.WriteString(html.EscapeString(child.Content))
				}
			}
			buf.WriteString("</pre>\n")
		}

	case Example:
		indentStr := strings.Repeat("    ", indent)
		var attrParts []string
		attrParts = append(attrParts, `data-role="example"`)
		if id := node.GetAttribute("id"); id != "" {
			attrParts = append(attrParts, fmt.Sprintf(`id="%s"`, html.EscapeString(id)))
		}
		otherAttrs := buildHTMLAttributes(node, []string{"id", "title"})
		if otherAttrs != "" {
			attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
		}
		attrs := buildAttrsString(attrParts...)
		fmt.Fprintf(buf, "%s<div%s>\n", indentStr, attrs)
		if title := node.GetAttribute("title"); title != "" {
			fmt.Fprintf(buf, "%s    <p data-role=\"example-title\">%s</p>\n", indentStr, html.EscapeString(title))
		}
		for _, child := range node.Children {
			toHTML(child, buf, xhtml, indent+1)
		}
		fmt.Fprintf(buf, "%s</div>\n", indentStr)

	case Sidebar:
		indentStr := strings.Repeat("    ", indent)
		var attrParts []string
		attrParts = append(attrParts, `data-role="sidebar"`)
		if id := node.GetAttribute("id"); id != "" {
			attrParts = append(attrParts, fmt.Sprintf(`id="%s"`, html.EscapeString(id)))
		}
		otherAttrs := buildHTMLAttributes(node, []string{"id", "title"})
		if otherAttrs != "" {
			attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
		}
		attrs := buildAttrsString(attrParts...)
		fmt.Fprintf(buf, "%s<aside%s>\n", indentStr, attrs)
		if title := node.GetAttribute("title"); title != "" {
			fmt.Fprintf(buf, "%s    <p data-role=\"sidebar-title\">%s</p>\n", indentStr, html.EscapeString(title))
		}
		for _, child := range node.Children {
			toHTML(child, buf, xhtml, indent+1)
		}
		fmt.Fprintf(buf, "%s</aside>\n", indentStr)

	case Quote:
		indentStr := strings.Repeat("    ", indent)
		var attrParts []string
		if id := node.GetAttribute("id"); id != "" {
			attrParts = append(attrParts, fmt.Sprintf(`id="%s"`, html.EscapeString(id)))
		}
		otherAttrs := buildHTMLAttributes(node, []string{"id", "attribution", "citation"})
		if otherAttrs != "" {
			attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
		}
		attrs := buildAttrsString(attrParts...)
		if attrs != "" {
			fmt.Fprintf(buf, "%s<blockquote%s>\n", indentStr, attrs)
		} else {
			fmt.Fprintf(buf, "%s<blockquote>\n", indentStr)
		}
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
		var attrParts []string
		if id := node.GetAttribute("id"); id != "" {
			attrParts = append(attrParts, fmt.Sprintf(`id="%s"`, html.EscapeString(id)))
		}
		otherAttrs := buildHTMLAttributes(node, []string{"id"})
		if otherAttrs != "" {
			attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
		}
		attrs := buildAttrsString(attrParts...)
		
		if attrs != "" {
			fmt.Fprintf(buf, "%s<table%s>\n", indentStr, attrs)
		} else {
			fmt.Fprintf(buf, "%s<table>\n", indentStr)
		}
		
		var theadStarted, tbodyStarted, tfootStarted bool
		for _, child := range node.Children {
			if child.Type == TableRow {
				rowRole := child.GetAttribute("role")
				if rowRole == "header" && !theadStarted {
					if tbodyStarted {
						// Close tbody before thead
						fmt.Fprintf(buf, "%s    </tbody>\n", indentStr)
						tbodyStarted = false
					}
					fmt.Fprintf(buf, "%s    <thead>\n", indentStr)
					theadStarted = true
				} else if rowRole == "footer" && !tfootStarted {
					if tbodyStarted {
						fmt.Fprintf(buf, "%s    </tbody>\n", indentStr)
						tbodyStarted = false
					}
					if theadStarted {
						fmt.Fprintf(buf, "%s    </thead>\n", indentStr)
						theadStarted = false
					}
					fmt.Fprintf(buf, "%s    <tfoot>\n", indentStr)
					tfootStarted = true
				} else if !theadStarted && !tfootStarted && !tbodyStarted {
					// First row without explicit role - assume header if first row
					if !theadStarted {
						fmt.Fprintf(buf, "%s    <thead>\n", indentStr)
						theadStarted = true
					}
				} else if !tbodyStarted && !tfootStarted && theadStarted && rowRole != "header" {
					// Transition from header to body
					fmt.Fprintf(buf, "%s    </thead>\n", indentStr)
					theadStarted = false
					fmt.Fprintf(buf, "%s    <tbody>\n", indentStr)
					tbodyStarted = true
				} else if !tbodyStarted && !theadStarted && !tfootStarted {
					// Regular body row
					fmt.Fprintf(buf, "%s    <tbody>\n", indentStr)
					tbodyStarted = true
				}
				
				var trAttrParts []string
				if rowRole != "" && rowRole != "header" && rowRole != "footer" {
					trAttrParts = append(trAttrParts, fmt.Sprintf(`data-role="%s"`, html.EscapeString(rowRole)))
				}
				otherAttrs := buildHTMLAttributes(child, []string{"role"})
				if otherAttrs != "" {
					trAttrParts = append(trAttrParts, strings.TrimSpace(otherAttrs))
				}
				trAttrs := buildAttrsString(trAttrParts...)
				
				if trAttrs != "" {
					fmt.Fprintf(buf, "%s        <tr%s>\n", indentStr, trAttrs)
				} else {
					fmt.Fprintf(buf, "%s        <tr>\n", indentStr)
				}
				
				for _, cell := range child.Children {
					if cell.Type == TableCell {
						cellTag := "td"
						if rowRole == "header" || (theadStarted && !tbodyStarted) {
							cellTag = "th"
						}
						
						var cellAttrParts []string
						if align := cell.GetAttribute("align"); align != "" {
							cellAttrParts = append(cellAttrParts, fmt.Sprintf(`data-align="%s"`, html.EscapeString(align)))
						}
						if colspan := cell.GetAttribute("colspan"); colspan != "" {
							cellAttrParts = append(cellAttrParts, fmt.Sprintf(`colspan="%s"`, html.EscapeString(colspan)))
						}
						if rowspan := cell.GetAttribute("rowspan"); rowspan != "" {
							cellAttrParts = append(cellAttrParts, fmt.Sprintf(`rowspan="%s"`, html.EscapeString(rowspan)))
						}
						otherAttrs := buildHTMLAttributes(cell, []string{"align", "colspan", "rowspan"})
						if otherAttrs != "" {
							cellAttrParts = append(cellAttrParts, strings.TrimSpace(otherAttrs))
						}
						cellAttrs := buildAttrsString(cellAttrParts...)
						
						if cellAttrs != "" {
							fmt.Fprintf(buf, "%s            <%s%s>", indentStr, cellTag, cellAttrs)
						} else {
							fmt.Fprintf(buf, "%s            <%s>", indentStr, cellTag)
						}
						toHTMLInlineContent(cell, buf, xhtml)
						fmt.Fprintf(buf, "</%s>\n", cellTag)
					}
				}
				fmt.Fprintf(buf, "%s        </tr>\n", indentStr)
			}
		}
		if theadStarted {
			fmt.Fprintf(buf, "%s    </thead>\n", indentStr)
		}
		if tbodyStarted {
			fmt.Fprintf(buf, "%s    </tbody>\n", indentStr)
		}
		if tfootStarted {
			fmt.Fprintf(buf, "%s    </tfoot>\n", indentStr)
		}
		fmt.Fprintf(buf, "%s</table>\n", indentStr)

	case Admonition:
		indentStr := strings.Repeat("    ", indent)
		admType := node.GetAttribute("type")
		var attrParts []string
		attrParts = append(attrParts, `data-role="admonition"`)
		if admType != "" {
			attrParts = append(attrParts, fmt.Sprintf(`data-asciidoc-variant="%s"`, html.EscapeString(admType)))
		}
		if id := node.GetAttribute("id"); id != "" {
			attrParts = append(attrParts, fmt.Sprintf(`id="%s"`, html.EscapeString(id)))
		}
		otherAttrs := buildHTMLAttributes(node, []string{"id", "type"})
		if otherAttrs != "" {
			attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
		}
		attrs := buildAttrsString(attrParts...)
		fmt.Fprintf(buf, "%s<div%s>\n", indentStr, attrs)
		if admType != "" {
			fmt.Fprintf(buf, "%s    <p data-role=\"admonition-title\">%s</p>\n", indentStr, html.EscapeString(strings.ToUpper(admType)))
		}
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
		buf.WriteString(indentStr + `<div data-role="page-break"></div>` + "\n")

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
		// Check if this is a cross-reference (has target attribute)
		target := node.GetAttribute("target")
		href := node.GetAttribute("href")
		if target != "" {
			// Cross-reference: use target as anchor link
			href = "#" + target
		}
		
		attrs := fmt.Sprintf(`href="%s"`, html.EscapeString(href))
		if title := node.GetAttribute("title"); title != "" {
			attrs += fmt.Sprintf(` title="%s"`, html.EscapeString(title))
		}
		if window := node.GetAttribute("window"); window != "" {
			attrs += fmt.Sprintf(` target="%s"`, html.EscapeString(window))
		}
		// Add other attributes as data-asciidoc-*
		attrs += buildHTMLAttributes(node, []string{"href", "title", "window", "target"})
		
		fmt.Fprintf(buf, "<a %s>", attrs)
		toHTMLInlineContent(node, buf, xhtml)
		buf.WriteString("</a>")

	case Passthrough:
		// Write content directly without escaping (for CMS-injected HTML)
		buf.WriteString(node.Content)

	case Superscript:
		buf.WriteString(`<sup data-asciidoc="superscript">`)
		toHTMLInlineContent(node, buf, xhtml)
		buf.WriteString("</sup>")

	case Subscript:
		buf.WriteString(`<sub data-asciidoc="subscript">`)
		toHTMLInlineContent(node, buf, xhtml)
		buf.WriteString("</sub>")

	case Highlight:
		buf.WriteString(`<mark data-asciidoc="highlight">`)
		toHTMLInlineContent(node, buf, xhtml)
		buf.WriteString("</mark>")

	case VerseBlock:
		indentStr := strings.Repeat("    ", indent)
		var attrParts []string
		attrParts = append(attrParts, `data-role="verse"`)
		if id := node.GetAttribute("id"); id != "" {
			attrParts = append(attrParts, fmt.Sprintf(`id="%s"`, html.EscapeString(id)))
		}
		otherAttrs := buildHTMLAttributes(node, []string{"id"})
		if otherAttrs != "" {
			attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
		}
		attrs := buildAttrsString(attrParts...)
		
		fmt.Fprintf(buf, "%s<div%s>\n", indentStr, attrs)
		// Preserve line breaks - convert \n to <br> or wrap in <p>
		content := getTextContent(node)
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if line != "" || i < len(lines)-1 {
				if line != "" {
					fmt.Fprintf(buf, "%s    <p>%s</p>\n", indentStr, html.EscapeString(line))
				} else {
					buf.WriteString(fmt.Sprintf("%s    <br>\n", indentStr))
				}
			}
		}
		if attribution := node.GetAttribute("attribution"); attribution != "" {
			fmt.Fprintf(buf, "%s    <footer><cite>%s</cite></footer>\n", indentStr, html.EscapeString(attribution))
		}
		fmt.Fprintf(buf, "%s</div>\n", indentStr)

	case OpenBlock:
		indentStr := strings.Repeat("    ", indent)
		var attrParts []string
		attrParts = append(attrParts, `data-role="open-block"`)
		if id := node.GetAttribute("id"); id != "" {
			attrParts = append(attrParts, fmt.Sprintf(`id="%s"`, html.EscapeString(id)))
		}
		otherAttrs := buildHTMLAttributes(node, []string{"id"})
		if otherAttrs != "" {
			attrParts = append(attrParts, strings.TrimSpace(otherAttrs))
		}
		attrs := buildAttrsString(attrParts...)
		
		fmt.Fprintf(buf, "%s<div%s>\n", indentStr, attrs)
		for _, child := range node.Children {
			toHTML(child, buf, xhtml, indent+1)
		}
		fmt.Fprintf(buf, "%s</div>\n", indentStr)

	case PassthroughBlock:
		// Output raw HTML without escaping (disable-output-escaping equivalent)
		buf.WriteString(node.Content)
		buf.WriteString("\n")

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
		} else if child.Type == InlineMacro {
			// Handle inline macros
			if child.Name == "anchor" {
				id := child.GetAttribute("id")
				if id == "" {
					id = child.GetAttribute("target")
				}
				if id != "" {
					fmt.Fprintf(buf, `<a id="%s"></a>`, html.EscapeString(id))
				}
			} else if child.Name == "kbd" {
				buf.WriteString(`<kbd data-role="keyboard">`)
				toHTMLInlineContent(child, buf, xhtml)
				buf.WriteString("</kbd>")
			} else if child.Name == "btn" {
				buf.WriteString(`<span data-role="button">`)
				toHTMLInlineContent(child, buf, xhtml)
				buf.WriteString("</span>")
			} else if child.Name == "menu" {
				// Parse menu path from target attribute
				target := child.GetAttribute("target")
				buf.WriteString(`<span data-role="menu-path">`)
				if target != "" {
					// Parse menu hierarchy (e.g., "File[New]")
					parts := strings.Split(target, "[")
					if len(parts) > 1 {
						menuItems := []string{parts[0]}
						rest := strings.TrimSuffix(strings.Join(parts[1:], "["), "]")
						if rest != "" {
							menuItems = append(menuItems, rest)
						}
						for i, item := range menuItems {
							if i > 0 {
								buf.WriteString(" → ")
							}
							fmt.Fprintf(buf, `<span data-role="menu">%s</span>`, html.EscapeString(strings.TrimSpace(item)))
						}
					} else {
						fmt.Fprintf(buf, `<span data-role="menu">%s</span>`, html.EscapeString(target))
					}
				} else {
					toHTMLInlineContent(child, buf, xhtml)
				}
				buf.WriteString("</span>")
			} else if child.Name == "footnote" {
				// Footnote reference - will be collected and rendered at document end
				ref := child.GetAttribute("ref")
				fnID := "fn1"
				if ref != "" {
					fnID = "fn" + ref
				}
				// For now, output simple reference - full footnote system would need collection
				fmt.Fprintf(buf, `<sup><a href="#%s" id="fnref%s" data-role="footnote-ref">1</a></sup>`, fnID, ref)
			} else if child.Name == "footnoteref" {
				ref := child.GetAttribute("ref")
				if ref == "" {
					ref = child.GetAttribute("target")
				}
				fnID := "fn" + ref
				fmt.Fprintf(buf, `<sup><a href="#%s" id="fnref%s" data-role="footnote-ref">1</a></sup>`, fnID, ref)
			} else {
				// Generic inline macro
				attrs := fmt.Sprintf(`data-role="macro" data-asciidoc-macro="%s"`, html.EscapeString(child.Name))
				attrs += buildHTMLAttributes(child, []string{})
				fmt.Fprintf(buf, `<span %s>`, attrs)
				toHTMLInlineContent(child, buf, xhtml)
				buf.WriteString("</span>")
			}
		} else {
			toHTML(child, buf, xhtml, 0)
		}
	}
}

// isStandardHTMLAttribute checks if an attribute name is a standard HTML5 attribute
func isStandardHTMLAttribute(name string) bool {
	standardAttrs := map[string]bool{
		"id": true, "href": true, "src": true, "alt": true, "title": true,
		"width": true, "height": true, "colspan": true, "rowspan": true,
		"type": true, "controls": true, "autoplay": true, "loop": true, "poster": true,
		"target": true, "rel": true, "lang": true, "charset": true,
	}
	return standardAttrs[name]
}

// buildHTMLAttributes builds HTML attributes string from node attributes
// Maps non-standard attributes to data-asciidoc-* format
// role is preserved for ARIA semantics
func buildHTMLAttributes(node *Node, excludeAttrs []string) string {
	excludeMap := make(map[string]bool)
	for _, attr := range excludeAttrs {
		excludeMap[attr] = true
	}
	
	var attrs []string
	for k, v := range node.Attributes {
		if excludeMap[k] {
			continue
		}
		
		// role is preserved for ARIA semantics
		if k == "role" {
			attrs = append(attrs, fmt.Sprintf(`role="%s"`, html.EscapeString(v)))
			continue
		}
		
		// Standard HTML attributes are used as-is
		if isStandardHTMLAttribute(k) {
			attrs = append(attrs, fmt.Sprintf(`%s="%s"`, k, html.EscapeString(v)))
		} else {
			// Non-standard attributes go to data-asciidoc-*
			attrs = append(attrs, fmt.Sprintf(`data-asciidoc-%s="%s"`, k, html.EscapeString(v)))
		}
	}
	
	if len(attrs) == 0 {
		return ""
	}
	return " " + strings.Join(attrs, " ")
}

// buildAttrsString builds a properly spaced attribute string from multiple attribute strings
// Ensures consistent spacing and returns a string that starts with a space
// for use in format strings like %s<tag%s>
func buildAttrsString(initialAttrs ...string) string {
	var parts []string
	for _, attr := range initialAttrs {
		trimmed := strings.TrimSpace(attr)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return " " + strings.Join(parts, " ")
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
		buf.WriteString(`<document xmlns="https://github.com/ndx-video/asciidoc-xml"`)
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>")
		} else {
			buf.WriteString(">\n")
			
			// Check if first child is a preamble (Paragraph with role="preamble")
			hasPreamble := false
			if len(node.Children) > 0 {
				firstChild := node.Children[0]
				if firstChild.Type == Paragraph && firstChild.GetAttribute("role") == "preamble" {
					hasPreamble = true
					// Output preamble wrapper
					buf.WriteString(indent + "  <preamble")
					if firstChild.GetAttribute("role") != "" {
						buf.WriteString(` role="` + escapeXML(firstChild.GetAttribute("role")) + `"`)
					}
					buf.WriteString(">\n")
					// Output preamble content (children of the paragraph)
					for _, grandchild := range firstChild.Children {
						toXML(grandchild, buf, indentLevel+2)
					}
					buf.WriteString(indent + "  </preamble>\n")
				}
			}
			
			// Output remaining children (skip first if it was preamble)
			startIdx := 0
			if hasPreamble {
				startIdx = 1
			}
			for i := startIdx; i < len(node.Children); i++ {
				toXML(node.Children[i], buf, indentLevel+1)
			}
			buf.WriteString(indent + "</document>\n")
		}

	case Section:
		buf.WriteString(indent + "<section")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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
		// Special handling for anchor and footnote macros
		if node.Name == "anchor" {
			// Output as dedicated <anchor> element
			id := node.GetAttribute("id")
			if id != "" {
				buf.WriteString(`<anchor id="` + escapeXML(id) + `"/>`)
			} else {
				// Fallback to macro if no id
				buf.WriteString(`<macro type="inline" name="anchor"`)
				for k, v := range node.Attributes {
					buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
				}
				buf.WriteString("/>")
			}
		} else if node.Name == "footnote" {
			// Output as dedicated <footnote> element
			buf.WriteString("<footnote")
			ref := node.GetAttribute("ref")
			if ref != "" {
				buf.WriteString(` ref="` + escapeXML(ref) + `"`)
			}
			buf.WriteString(">")
			toXMLInlineContent(node, buf)
			buf.WriteString("</footnote>")
		} else {
			// Generic inline macro
			buf.WriteString(`<macro type="inline" name="` + escapeXML(node.Name) + `"`)
			for k, v := range node.Attributes {
				buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
			}
			if len(node.Children) == 0 {
				buf.WriteString("/>")
			} else {
				buf.WriteString(">")
				toXMLInlineContent(node, buf)
				buf.WriteString("</macro>")
			}
		}

	case Text:
		buf.WriteString(escapeXML(node.Content))

	case List:
		buf.WriteString(indent + "<list")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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

	case VerseBlock:
		buf.WriteString(indent + "<verseblock")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">\n")
			for _, child := range node.Children {
				toXML(child, buf, indentLevel+1)
			}
			buf.WriteString(indent + "</verseblock>\n")
		}

	case OpenBlock:
		buf.WriteString(indent + "<openblock")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
		}
		if len(node.Children) == 0 {
			buf.WriteString("/>\n")
		} else {
			buf.WriteString(">\n")
			for _, child := range node.Children {
				toXML(child, buf, indentLevel+1)
		}
		buf.WriteString(indent + "</openblock>\n")
	}

	case PassthroughBlock:
		buf.WriteString(indent + "<passthrough")
		for k, v := range node.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
		}
		buf.WriteString(">")
		buf.WriteString(escapeXML(node.Content))
		buf.WriteString("</passthrough>\n")

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
			buf.WriteString(fmt.Sprintf(` %s="%s"`, sanitizeXMLAttributeName(k), escapeXML(v)))
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

	case Superscript:
		buf.WriteString("<superscript>")
		toXMLInlineContent(node, buf)
		buf.WriteString("</superscript>")

	case Subscript:
		buf.WriteString("<subscript>")
		toXMLInlineContent(node, buf)
		buf.WriteString("</subscript>")

	case Highlight:
		buf.WriteString("<highlight>")
		toXMLInlineContent(node, buf)
		buf.WriteString("</highlight>")

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

// sanitizeXMLAttributeName converts an attribute name to be XML-compliant
// XML attribute names must:
// - Start with a letter or underscore
// - Contain only letters, digits, hyphens, underscores, and periods
// - Not contain colons (reserved for namespaces)
func sanitizeXMLAttributeName(name string) string {
	if name == "" {
		return ""
	}
	
	// Trim leading/trailing colons (common in AsciiDoc attributes like :toclevels:)
	name = strings.Trim(name, ":")
	
	// If still empty after trimming, return a valid default
	if name == "" {
		return "attr"
	}
	
	var result strings.Builder
	
	for i, c := range name {
		// First character must be letter or underscore
		if i == 0 {
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' {
				result.WriteRune(c)
			} else if c >= '0' && c <= '9' {
				// Numbers not allowed at start, prefix with underscore
				result.WriteRune('_')
				result.WriteRune(c)
			} else {
				// Invalid character at start, replace with underscore
				result.WriteRune('_')
			}
		} else {
			// Subsequent characters can be letters, digits, hyphens, underscores, or periods
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' {
				result.WriteRune(c)
			} else if c == ':' {
				// Replace colons with hyphens
				result.WriteRune('-')
			} else {
				// Replace other invalid characters with underscores
				result.WriteRune('_')
			}
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


// ConvertMarkdownToAsciiDoc converts Markdown content to AsciiDoc format.
// For large files, consider using ConvertMarkdownToAsciiDocStreaming instead
// to avoid loading the entire result into memory.
func ConvertMarkdownToAsciiDoc(reader io.Reader) (string, error) {
	var result bytes.Buffer
	if err := ConvertMarkdownToAsciiDocStreaming(reader, &result); err != nil {
		return "", err
	}
	return result.String(), nil
}

// convertMarkdownToAsciiDocLegacy is the original implementation kept for reference.
// It's been replaced by ConvertMarkdownToAsciiDocStreaming for better memory efficiency.
func convertMarkdownToAsciiDocLegacy(reader io.Reader) (string, error) {
	scanner := bufio.NewScanner(reader)
	var result bytes.Buffer
	var inCodeBlock bool
	var codeBlockLang string
	var inTable bool
	var codeBlockLines []string
	var frontmatterProcessed bool
	var frontmatterLines []string
	var inFrontmatter bool

	// Regex patterns
	headerRegex := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	codeBlockStartRegex := regexp.MustCompile("^```(\\w*)$")
	codeBlockEndRegex := regexp.MustCompile("^```$")
	tableRowRegex := regexp.MustCompile(`^\|(.+)\|$`)
	tableSeparatorRegex := regexp.MustCompile(`^\|?[\s\-:]+\|[\s\-:]+\|?$`)
	horizontalRuleRegex := regexp.MustCompile(`^[-*_]{3,}$`)
	blockquoteRegex := regexp.MustCompile(`^>\s*(.*)$`)
	orderedListRegex := regexp.MustCompile(`^\s*(\d+)\.\s+(.+)$`)
	unorderedListRegex := regexp.MustCompile(`^\s*[-*+]\s+(.+)$`)
	frontmatterStartRegex := regexp.MustCompile(`^---\s*$`)
	frontmatterEndRegex := regexp.MustCompile(`^---\s*$|^\.\.\.\s*$`)

	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Task 2: Handle YAML Frontmatter at the start of the file
		if !frontmatterProcessed && lineNum == 1 {
			if frontmatterStartRegex.MatchString(line) {
				inFrontmatter = true
				continue
			}
		}

		if inFrontmatter {
			if frontmatterEndRegex.MatchString(line) {
				// Process frontmatter
				frontmatterProcessed = true
				inFrontmatter = false
				processFrontmatter(&result, frontmatterLines)
				continue
			}
			frontmatterLines = append(frontmatterLines, line)
			continue
		}

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

		// Task 5: Handle blockquotes (including admonitions)
		if matches := blockquoteRegex.FindStringSubmatch(line); matches != nil {
			content := matches[1]
			// Check for admonition markers
			admonitionType := ""
			// Try **Note** or **WARNING** format first
			boldAdmonitionRegex := regexp.MustCompile(`^\*\*(\w+)\*\*\s*`)
			if admMatches := boldAdmonitionRegex.FindStringSubmatch(content); admMatches != nil {
				admonitionType = strings.ToUpper(admMatches[1])
				content = strings.TrimSpace(boldAdmonitionRegex.ReplaceAllString(content, ""))
			} else {
				// Try [!NOTE] or [!WARNING] format
				bracketAdmonitionRegex := regexp.MustCompile(`^\[!(\w+)\]\s*`)
				if admMatches := bracketAdmonitionRegex.FindStringSubmatch(content); admMatches != nil {
					admonitionType = strings.ToUpper(admMatches[1])
					content = strings.TrimSpace(bracketAdmonitionRegex.ReplaceAllString(content, ""))
				}
			}
			
			if admonitionType != "" {
				// Convert to AsciiDoc admonition
				result.WriteString(fmt.Sprintf("%s: %s\n", admonitionType, content))
			} else {
				// Regular blockquote
				result.WriteString(fmt.Sprintf("[quote]\n____\n%s\n____\n", content))
			}
			continue
		}

		// Task 3: Handle tables (skip separator rows)
		if tableRowRegex.MatchString(line) {
			// Task 3: Skip table separator rows
			if tableSeparatorRegex.MatchString(line) {
				continue
			}
			if !inTable {
				// Count columns from first row
				cells := strings.Split(line, "|")
				// Remove empty first/last cells from split
				if len(cells) > 0 && strings.TrimSpace(cells[0]) == "" {
					cells = cells[1:]
				}
				if len(cells) > 0 && strings.TrimSpace(cells[len(cells)-1]) == "" {
					cells = cells[:len(cells)-1]
				}
				colCount := len(cells)
				// Write table start with column count
				result.WriteString(fmt.Sprintf("[cols=\"%s\"]\n", strings.Repeat("1,", colCount)[:len(strings.Repeat("1,", colCount))-1]))
				result.WriteString("|===\n")
				inTable = true
			}
			// Convert table row - process inline markdown in cells
			cells := strings.Split(line, "|")
			// Remove empty first/last cells from split
			if len(cells) > 0 && strings.TrimSpace(cells[0]) == "" {
				cells = cells[1:]
			}
			if len(cells) > 0 && strings.TrimSpace(cells[len(cells)-1]) == "" {
				cells = cells[:len(cells)-1]
			}
			// Process each cell and convert inline markdown
			processedCells := make([]string, len(cells))
			for i, cell := range cells {
				cellContent := strings.TrimSpace(cell)
				processedCells[i] = convertInlineMarkdown(cellContent, false)
			}
			// Write row with proper AsciiDoc table format
			result.WriteString("|" + strings.Join(processedCells, " |") + "\n")
			continue
		} else if inTable {
			// End table if we hit a non-table line
			result.WriteString("|===\n")
			inTable = false
		}

		// Task 4: Handle ordered lists (with nesting support)
		if matches := orderedListRegex.FindStringSubmatch(line); matches != nil {
			content := matches[2]
			// Calculate indentation level
			leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
			// Convert inline formatting in list items
			content = convertInlineMarkdown(content, false)
			// AsciiDoc uses . for ordered lists, nesting with additional spaces
			indent := strings.Repeat(" ", leadingSpaces)
			result.WriteString(fmt.Sprintf("%s. %s\n", indent, content))
			continue
		}

		// Task 4: Handle unordered lists (with nesting support)
		if matches := unorderedListRegex.FindStringSubmatch(line); matches != nil {
			content := matches[1]
			// Calculate indentation level
			leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
			// Convert inline formatting in list items
			content = convertInlineMarkdown(content, false)
			// AsciiDoc uses * for unordered lists, nesting with additional spaces or *
			indent := strings.Repeat(" ", leadingSpaces)
			result.WriteString(fmt.Sprintf("%s* %s\n", indent, content))
			continue
		}

		// Handle regular lines
		if strings.TrimSpace(line) == "" {
			result.WriteString("\n")
			continue
		}

		// Check if line is a standalone image
		trimmedLine := strings.TrimSpace(line)
		imageRegex := regexp.MustCompile(`^!\[([^\]]*)\]\(([^)]+)\)\s*$`)
		if imageRegex.MatchString(trimmedLine) {
			matches := imageRegex.FindStringSubmatch(trimmedLine)
			alt := matches[1]
			src := matches[2]
			result.WriteString(fmt.Sprintf("image::%s[%s]\n", src, alt))
			continue
		}

		// Convert inline Markdown in regular text (handles images, links, formatting)
		convertedLine := convertInlineMarkdown(line, false)
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
// isStandaloneLine indicates if this is a line containing only the image (for block vs inline)
func convertInlineMarkdown(text string, isStandaloneLine bool) string {
	// Task 1: Convert images - determine if block or inline
	imageRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	// Check if the line is just an image (possibly with whitespace)
	trimmedText := strings.TrimSpace(text)
	imageOnly := imageRegex.MatchString(trimmedText) && strings.TrimSpace(imageRegex.ReplaceAllString(trimmedText, "")) == ""
	
	text = imageRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := imageRegex.FindStringSubmatch(match)
		alt := matches[1]
		src := matches[2]
		// Task 1: Use :: for block images, : for inline
		if imageOnly || isStandaloneLine {
			return fmt.Sprintf("image::%s[%s]", src, alt)
		}
		return fmt.Sprintf("image:%s[%s]", src, alt)
	})

	// Task 1: Convert links [text](url) to link:url[text]
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = linkRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := linkRegex.FindStringSubmatch(match)
		linkText := matches[1]
		url := matches[2]
		return fmt.Sprintf("link:%s[%s]", url, linkText)
	})

	// Convert bold **text** to **text** (AsciiDoc bold) using placeholder to avoid conflicts with italic
	boldDoubleStarRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	placeholderMap := make(map[string]string)
	placeholderCounter := 0
	
	// Replace bold with placeholder first
	text = boldDoubleStarRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholderCounter++
		placeholder := fmt.Sprintf("__BOLD_PLACEHOLDER_%d__", placeholderCounter)
		matches := boldDoubleStarRegex.FindStringSubmatch(match)
		// AsciiDoc uses **text** for bold, same as Markdown
		placeholderMap[placeholder] = "**" + matches[1] + "**"
		return placeholder
	})
	
	// Convert bold __text__ to **text** (AsciiDoc bold, also use placeholder)
	boldDoubleUnderscoreRegex := regexp.MustCompile(`__([^_]+)__`)
	text = boldDoubleUnderscoreRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholderCounter++
		placeholder := fmt.Sprintf("__BOLD_PLACEHOLDER_%d__", placeholderCounter)
		matches := boldDoubleUnderscoreRegex.FindStringSubmatch(match)
		// Convert Markdown __text__ to AsciiDoc **text**
		placeholderMap[placeholder] = "**" + matches[1] + "**"
		return placeholder
	})

	// Convert italic *text* to _text_ (now safe since bold is in placeholders)
	// Use a more robust regex that ensures we have word boundaries or non-word characters
	// Important: Don't match *text* if it's part of **text** (double asterisks for bold)
	// Match *text* but ensure it's not preceded or followed by another *
	italicAsteriskRegex := regexp.MustCompile(`([^*]|^)\*([^*\n]+?)\*([^*]|$)`)
	text = italicAsteriskRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := italicAsteriskRegex.FindStringSubmatch(match)
		if len(matches) >= 4 {
			before := matches[1]
			content := matches[2]
			after := matches[3]
			// Only convert if the content is not empty
			if strings.TrimSpace(content) != "" {
				return before + "_" + content + "_" + after
			}
		}
		return match
	})

	// Convert italic _text_ to _text_ (already correct format, no change needed)
	// Note: Markdown uses _text_ for italic, AsciiDoc also uses _text_ for italic, so no conversion needed
	// However, we should ensure that _text_ is preserved correctly

	// Restore bold placeholders
	for placeholder, replacement := range placeholderMap {
		text = strings.ReplaceAll(text, placeholder, replacement)
	}

	// Inline code `code` is already in correct format for AsciiDoc
	// No conversion needed

	return text
}

// processFrontmatter converts YAML frontmatter to AsciiDoc header attributes
func processFrontmatter(result *bytes.Buffer, lines []string) {
	frontmatter := make(map[string]interface{})
	var currentKey string
	var currentArray []string
	var inArray bool
	
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Check if this is an array item (starts with -)
		if strings.HasPrefix(line, "-") {
			if inArray && currentKey != "" {
				// Extract the value after the dash
				value := strings.TrimSpace(strings.TrimPrefix(line, "-"))
				// Remove quotes if present
				if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
					value = value[1 : len(value)-1]
				}
				currentArray = append(currentArray, value)
			}
			continue
		}
		
		// If we were in an array, save it now
		if inArray && currentKey != "" && len(currentArray) > 0 {
			frontmatter[currentKey] = currentArray
			currentArray = []string{}
			inArray = false
		}
		
		// Check if this is a key:value pair
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			
			// Check if value is empty (indicating start of array or nested structure)
			if value == "" {
				// Check if next line starts with - (array) or has indentation (nested structure)
				if i+1 < len(lines) {
					nextLine := strings.TrimSpace(lines[i+1])
					if strings.HasPrefix(nextLine, "-") {
						// This is the start of an array
						currentKey = key
						currentArray = []string{}
						inArray = true
						continue
					} else {
						// Check if next line has indentation (nested structure)
						nextLineRaw := lines[i+1]
						if strings.HasPrefix(nextLineRaw, "  ") || strings.HasPrefix(nextLineRaw, "\t") {
							// This is a nested structure - collect all indented lines
							var nestedLines []string
							baseIndent := ""
							// Determine base indentation level
							for _, char := range nextLineRaw {
								if char == ' ' || char == '\t' {
									baseIndent += string(char)
								} else {
									break
								}
							}
							
							for j := i + 1; j < len(lines); j++ {
								indentedLine := lines[j]
								trimmedLine := strings.TrimSpace(indentedLine)
								if trimmedLine == "" {
									// Keep empty lines in nested structure
									nestedLines = append(nestedLines, "")
									continue
								}
								// Check if line is still part of nested structure
								if strings.HasPrefix(indentedLine, baseIndent) || strings.HasPrefix(indentedLine, "  ") || strings.HasPrefix(indentedLine, "\t") {
									nestedLines = append(nestedLines, indentedLine)
								} else {
									break
								}
							}
							// Store nested structure as a single string value (preserve formatting)
							if len(nestedLines) > 0 {
								frontmatter[key] = strings.Join(nestedLines, "\n")
								// Mark nested lines as processed by clearing them
								for k := 0; k < len(nestedLines); k++ {
									if i+1+k < len(lines) {
										lines[i+1+k] = ""
									}
								}
								continue
							}
						}
					}
				}
				// Empty value, store as empty string
				frontmatter[key] = ""
				continue
			}
			
			// Remove quotes if present
			if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
				value = value[1 : len(value)-1]
			}
			frontmatter[key] = value
		}
	}
	
	// Handle any remaining array
	if inArray && currentKey != "" && len(currentArray) > 0 {
		frontmatter[currentKey] = currentArray
	}
	
	// Special handling: title becomes document header
	if title, ok := frontmatter["title"]; ok {
		if titleStr, ok := title.(string); ok {
			result.WriteString(fmt.Sprintf("= %s\n", titleStr))
			delete(frontmatter, "title")
		}
	}
	
	// All other keys become AsciiDoc attributes
	for key, value := range frontmatter {
		var attrValue string
		
		switch v := value.(type) {
		case []string:
			// Array: convert to comma-separated values
			attrValue = strings.Join(v, ", ")
		case string:
			// String value
			attrValue = v
		default:
			// Fallback: convert to string
			attrValue = fmt.Sprintf("%v", v)
		}
		
		result.WriteString(fmt.Sprintf(":%s: %s\n", key, attrValue))
	}
	
	// Add blank line after frontmatter
	result.WriteString("\n")
}

