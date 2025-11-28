package lib

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

var componentMacroRegex = regexp.MustCompile(`^component::\w+\[.*\]$`)

// Parse parses AsciiDoc content from a reader and returns a Document Node
func Parse(reader io.Reader) (*Node, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	parser := newParser(string(content))
	return parser.parse()
}

// parser is a text-based AsciiDoc parser
type parser struct {
	lines      []string
	lineNum    int
	doc        *Node
	header     *Node
	attributes map[string]string
}

func newParser(content string) *parser {
	lines := strings.Split(content, "\n")
	return &parser{
		lines:      lines,
		lineNum:    0,
		attributes: make(map[string]string),
	}
}

func (p *parser) parse() (*Node, error) {
	p.doc = NewElementNode("asciidoc")
	p.doc.SetAttribute("doctype", "article") // Default

	// Parse header and attributes
	p.parseHeader()

	// Parse content
	p.parseContent(p.doc, nil)

	return p.doc, nil
}

func (p *parser) parseHeader() {
	headerNode := NewElementNode("header")
	hasHeader := false

	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])

		// Document title
		if strings.HasPrefix(line, "=") && !strings.HasPrefix(line, "==") {
			if !hasHeader {
				hasHeader = true
				p.doc.Children = append([]*Node{headerNode}, p.doc.Children...) // Prepend header
			}
			
			titleNode := NewElementNode("h1")
			titleText := strings.TrimSpace(strings.TrimPrefix(line, "="))
			titleNode.AddChild(NewTextNode(titleText))
			headerNode.AddChild(titleNode)
			
			p.lineNum++
			continue
		}

		// Attributes
		if strings.HasPrefix(line, ":") {
			parts := strings.SplitN(line[1:], ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				p.attributes[key] = value

				if !hasHeader {
					hasHeader = true
					p.doc.Children = append([]*Node{headerNode}, p.doc.Children...)
				}

				switch key {
				case "author":
					authorNode := NewElementNode("address")
					authorNode.SetAttribute("class", "author")
					nameNode := NewElementNode("span")
					nameNode.SetAttribute("class", "author-name")
					nameNode.AddChild(NewTextNode(value))
					authorNode.AddChild(nameNode)
					headerNode.AddChild(authorNode)
				case "email":
					// Find the last author to add email
					var lastAuthor *Node
					for i := len(headerNode.Children) - 1; i >= 0; i-- {
						if headerNode.Children[i].Data == "address" {
							lastAuthor = headerNode.Children[i]
							break
						}
					}
					if lastAuthor != nil {
						emailNode := NewElementNode("a")
						emailNode.SetAttribute("class", "email")
						emailNode.SetAttribute("href", "mailto:"+value)
						emailNode.AddChild(NewTextNode(value))
						lastAuthor.AddChild(emailNode)
					}
				case "revnumber":
					revNode := p.getOrAddRevision(headerNode)
					numberNode := NewElementNode("span")
					numberNode.SetAttribute("class", "revnumber")
					numberNode.AddChild(NewTextNode(value))
					revNode.AddChild(numberNode)
				case "revdate":
					revNode := p.getOrAddRevision(headerNode)
					dateNode := NewElementNode("span")
					dateNode.SetAttribute("class", "revdate")
					dateNode.AddChild(NewTextNode(value))
					revNode.AddChild(dateNode)
				case "revremark":
					revNode := p.getOrAddRevision(headerNode)
					remarkNode := NewElementNode("span")
					remarkNode.SetAttribute("class", "revremark")
					remarkNode.AddChild(NewTextNode(value))
					revNode.AddChild(remarkNode)
				case "doctype":
					p.doc.SetAttribute("doctype", value)
				default:
					attrNode := NewElementNode("attribute")
					attrNode.SetAttribute("name", key)
					attrNode.SetAttribute("value", value)
					headerNode.AddChild(attrNode)
				}
			}
			p.lineNum++
			continue
		}

		// Empty line after header ends header parsing
		if line == "" {
			p.lineNum++
			// Only break if we actually started a header
			if hasHeader {
				break
			}
			continue
		}

		// If we hit content, stop parsing header
		if !strings.HasPrefix(line, "=") && line != "" {
			break
		}

		p.lineNum++
	}
}

func (p *parser) getOrAddRevision(header *Node) *Node {
	for _, child := range header.Children {
		if child.Data == "div" && child.GetAttribute("class") == "revision" {
			return child
		}
	}
	revNode := NewElementNode("div")
	revNode.SetAttribute("class", "revision")
	header.AddChild(revNode)
	return revNode
}

// parseContent parses content items, optionally stopping at sections at or above maxLevel
// If maxLevel is nil, it will continue until end of document
func (p *parser) parseContent(parent *Node, maxLevel *int) {
	for p.lineNum < len(p.lines) {
		line := p.lines[p.lineNum]
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			p.lineNum++
			continue
		}

		// Code block
		if strings.HasPrefix(trimmed, "----") || strings.HasPrefix(trimmed, "```") {
			codeBlock := p.parseCodeBlock()
			if codeBlock != nil {
				parent.AddChild(codeBlock)
			}
			continue
		}

		// Literal block
		if strings.HasPrefix(trimmed, "....") {
			literalBlock := p.parseLiteralBlock()
			if literalBlock != nil {
				parent.AddChild(literalBlock)
			}
			continue
		}

		// Example block (check before sections to avoid matching ==== as a section)
		if strings.HasPrefix(trimmed, "====") {
			example := p.parseExampleBlock()
			if example != nil {
				parent.AddChild(example)
			}
			continue
		}

		// Check if we hit a section at same or higher level (if maxLevel is set)
		if strings.HasPrefix(trimmed, "=") {
			if maxLevel != nil {
				level := 0
				line := trimmed
				for strings.HasPrefix(line, "=") {
					level++
					line = line[1:]
				}
				// Convert to section level (subtract 1)
				sectionLevel := level - 1
				if sectionLevel <= *maxLevel {
					break
				}
			}
			// Parse section
			section := p.parseSection()
			if section != nil {
				parent.AddChild(section)
			}
			continue
		}

		// Sidebar
		if strings.HasPrefix(trimmed, "****") {
			sidebar := p.parseSidebar()
			if sidebar != nil {
				parent.AddChild(sidebar)
			}
			continue
		}

		// Quote
		if strings.HasPrefix(trimmed, "____") {
			quote := p.parseQuote()
			if quote != nil {
				parent.AddChild(quote)
			}
			continue
		}

		// Table
		if strings.HasPrefix(trimmed, "|===") {
			table := p.parseTable()
			if table != nil {
				parent.AddChild(table)
			}
			continue
		}

		// Component macro (check before list to avoid matching component:: as labeled list)
		// Strictly check for Block Macro syntax: component::name[attrs]
		if componentMacroRegex.MatchString(trimmed) {
			component := p.parseComponentMacro()
			if component != nil {
				parent.AddChild(component)
			}
			continue
		}

		// List (but check if it's actually a title for a block)
		if p.isListItem(trimmed) {
			if p.isBlockTitle(trimmed) {
				// This is a block title, skip this line and continue
				p.lineNum++
				continue
			}
			list := p.parseList()
			if list != nil {
				parent.AddChild(list)
			}
			continue
		}

		// Admonition
		if p.isAdmonition(trimmed) {
			admonition := p.parseAdmonition()
			if admonition != nil {
				parent.AddChild(admonition)
			}
			continue
		}

		// Image
		if strings.HasPrefix(trimmed, "image::") || strings.HasPrefix(trimmed, "image:") {
			image := p.parseImage()
			if image != nil {
				parent.AddChild(image)
			}
			continue
		}

		// Thematic break
		if trimmed == "'''" {
			tb := NewElementNode("hr")
			parent.AddChild(tb)
			p.lineNum++
			continue
		}

		// Page break
		if trimmed == "<<<" {
			pb := NewElementNode("div")
			pb.SetAttribute("class", "page-break")
			parent.AddChild(pb)
			p.lineNum++
			continue
		}

		// Paragraph
		para := p.parseParagraph()
		if para != nil {
			parent.AddChild(para)
		} else {
			// If paragraph parsing fails or returns nil, increment to avoid infinite loop
			p.lineNum++
		}
	}
}

func (p *parser) parseSection() *Node {
	line := strings.TrimSpace(p.lines[p.lineNum])
	level := 0
	marker := ""
	for strings.HasPrefix(line, "=") {
		level++
		marker += "="
		line = line[1:]
	}

	// In AsciiDoc, = is document title, == is level 1 section, === is level 2, etc.
	// So we subtract 1 from the count
	sectionLevel := level - 1

	titleText := strings.TrimSpace(line)
	section := NewElementNode("section")
	section.SetAttribute("level", fmt.Sprintf("%d", sectionLevel))
	section.SetAttribute("marker", marker)
	section.SetAttribute("title", titleText)
    
    // Add title as heading element
	hLevel := sectionLevel + 1
	if hLevel > 6 {
		hLevel = 6
	}
	hTag := fmt.Sprintf("h%d", hLevel)
	hNode := NewElementNode(hTag)
	hNode.AddChild(NewTextNode(titleText))
	section.AddChild(hNode)

	p.lineNum++

	// Parse section content, stopping at sections at same or higher level
	p.parseContent(section, &sectionLevel)

	return section
}

func (p *parser) parseParagraph() *Node {
	var lines []string

	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])

		if line == "" {
			break
		}

		// Stop at block delimiters
		if strings.HasPrefix(line, "=") ||
			strings.HasPrefix(line, "----") ||
			strings.HasPrefix(line, "....") ||
			strings.HasPrefix(line, "====") ||
			strings.HasPrefix(line, "****") ||
			strings.HasPrefix(line, "____") ||
			strings.HasPrefix(line, "|===") ||
			componentMacroRegex.MatchString(line) ||
			strings.HasPrefix(line, "image::") ||
			strings.HasPrefix(line, "image:") ||
			p.isListItem(line) ||
			p.isAdmonition(line) {
			break
		}

		lines = append(lines, line)
		p.lineNum++
	}

	if len(lines) == 0 {
		return nil
	}

	text := strings.Join(lines, " ")
	para := NewElementNode("p")
	p.parseInlineContent(para, text)
	return para
}

func (p *parser) parseCodeBlock() *Node {
	// Parse attributes from previous line(s) if present
	var language, title string
	// Look back for title and attributes (skip empty lines)
	for i := p.lineNum - 1; i >= 0 && i >= p.lineNum-3; i-- {
		prevLine := strings.TrimSpace(p.lines[i])
		if prevLine == "" {
			continue
		}
		if strings.HasPrefix(prevLine, "[") && strings.HasSuffix(prevLine, "]") {
			attrs := prevLine[1 : len(prevLine)-1]
			parts := strings.Split(attrs, ",")
			// Format is [source,language] or [language]
			if len(parts) > 1 {
				// [source,go] format
				language = strings.TrimSpace(parts[1])
			} else if len(parts) > 0 {
				// [go] format
				firstPart := strings.TrimSpace(parts[0])
				if firstPart != "source" {
					language = firstPart
				}
			}
		}
		if strings.HasPrefix(prevLine, ".") {
			title = strings.TrimPrefix(prevLine, ".")
			break // Title found, stop looking
		}
	}

	p.lineNum++ // Skip opening delimiter

	var content []string
	for p.lineNum < len(p.lines) {
		line := p.lines[p.lineNum]
		if strings.TrimSpace(line) == "----" || strings.TrimSpace(line) == "```" {
			p.lineNum++
			break
		}
		content = append(content, line)
		p.lineNum++
	}

	pre := NewElementNode("pre")
	code := NewElementNode("code")

	if language != "" {
		code.SetAttribute("class", "language-"+language)
		pre.SetAttribute("data-language", language)
	}
	if title != "" {
		pre.SetAttribute("title", title)
	}
	code.AddChild(NewTextNode(strings.Join(content, "\n")))
	pre.AddChild(code)
	return pre
}

func (p *parser) parseLiteralBlock() *Node {
	p.lineNum++ // Skip opening
	var content []string
	for p.lineNum < len(p.lines) {
		line := p.lines[p.lineNum]
		if strings.TrimSpace(line) == "...." {
			p.lineNum++
			break
		}
		content = append(content, line)
		p.lineNum++
	}
	literalBlock := NewElementNode("pre")
	literalBlock.SetAttribute("class", "literal-block")
	literalBlock.AddChild(NewTextNode(strings.Join(content, "\n")))
	return literalBlock
}

func (p *parser) parseExampleBlock() *Node {
	// Check for title on previous non-empty line(s)
	var title string
	for i := p.lineNum - 1; i >= 0; i-- {
		prevLine := strings.TrimSpace(p.lines[i])
		if prevLine == "" {
			continue
		}
		if strings.HasPrefix(prevLine, ".") {
			title = strings.TrimPrefix(prevLine, ".")
			break
		}
		break
	}

	p.lineNum++ // Skip opening

	var contentLines []string
	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])
		if line == "====" {
			p.lineNum++
			break
		}
		contentLines = append(contentLines, p.lines[p.lineNum])
		p.lineNum++
	}

	// Re-parse content within the example block
	// For simplicity in this refactor, we can create a new parser or recursively call parseContent logic
	// But we need to pass the lines. 
	// Let's just use a sub-parser approach for simplicity since we have the lines extracted.
	
	subContent := strings.Join(contentLines, "\n")
	subParser := newParser(subContent)
	example := NewElementNode("div")
	example.SetAttribute("class", "example")
	if title != "" {
		example.SetAttribute("title", title)
	}
	subParser.parseContent(example, nil)
	
	return example
}

func (p *parser) parseSidebar() *Node {
	p.lineNum++ // Skip opening
	var title string
	var contentLines []string
	
	// Check for title in previous lines logic is tricky if we already advanced. 
	// Assuming title was checked before or is inside.
	// Actually the previous parser checked lines *before* the block start.
	// Let's check if there was a title before.
	for i := p.lineNum - 2; i >= 0; i-- {
		prevLine := strings.TrimSpace(p.lines[i])
		if prevLine == "" {
			continue
		}
		if strings.HasPrefix(prevLine, ".") {
			title = strings.TrimPrefix(prevLine, ".")
		}
		break
	}

	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])
		if line == "****" {
			p.lineNum++
			break
		}
		contentLines = append(contentLines, p.lines[p.lineNum])
		p.lineNum++
	}

	sidebar := NewElementNode("aside")
	sidebar.SetAttribute("class", "sidebar")
	if title != "" {
		sidebar.SetAttribute("title", title)
	}
	
	subContent := strings.Join(contentLines, "\n")
	subParser := newParser(subContent)
	subParser.parseContent(sidebar, nil)
	
	return sidebar
}

func (p *parser) parseQuote() *Node {
	p.lineNum++ // Skip opening
	var contentLines []string
	
	var attribution, citation string
	// Check for attribution line at bottom usually? 
	// AsciiDoc quotes: 
	// [quote, attribution, citation]
	// ____
	// content
	// ____
	
	// Check for attributes block above
	for i := p.lineNum - 2; i >= 0 && i >= p.lineNum-4; i-- {
		prevLine := strings.TrimSpace(p.lines[i])
		if prevLine == "" {
			continue
		}
		if strings.HasPrefix(prevLine, "[") && strings.HasSuffix(prevLine, "]") {
			attrs := prevLine[1 : len(prevLine)-1]
			parts := strings.Split(attrs, ",")
			// quote, attribution, citation
			if len(parts) > 1 {
				attribution = strings.TrimSpace(parts[1])
			}
			if len(parts) > 2 {
				citation = strings.TrimSpace(parts[2])
			}
		}
		break
	}

	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])
		if line == "____" {
			p.lineNum++
			break
		}
		contentLines = append(contentLines, p.lines[p.lineNum])
		p.lineNum++
	}

	quote := NewElementNode("blockquote")
	quote.SetAttribute("class", "quote")
	if attribution != "" {
		quote.SetAttribute("attribution", attribution)
	}
	if citation != "" {
		quote.SetAttribute("citation", citation)
	}

	subContent := strings.Join(contentLines, "\n")
	subParser := newParser(subContent)
	subParser.parseContent(quote, nil)

	return quote
}

func (p *parser) parseTable() *Node {
	// Parse table attributes
	tableAttrs := make(map[string]string)
	
	// Look back for attributes [cols="...", options="..."]
	for i := p.lineNum - 1; i >= 0; i-- {
		prevLine := strings.TrimSpace(p.lines[i])
		if prevLine == "" {
			continue
		}
		if strings.HasPrefix(prevLine, "[") && strings.HasSuffix(prevLine, "]") {
			// Parse attributes: key=val, key=val or just val
			content := prevLine[1 : len(prevLine)-1]
			
			// Simple parser for key="val" or key=val
			// Note: This is a basic implementation and might need a proper scanner for complex attributes
			parts := strings.Split(content, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.Contains(part, "=") {
					kv := strings.SplitN(part, "=", 2)
					key := strings.TrimSpace(kv[0])
					val := strings.TrimSpace(kv[1])
					val = strings.Trim(val, "\"")
					tableAttrs[key] = val
				} else {
					// Positional or boolean attributes
					// Currently just storing as is if needed, or specific handling
					// For now, handle 'header', 'footer' which are often in options
					if part == "header" || part == "footer" {
						// Append to options
						if opts, ok := tableAttrs["options"]; ok {
							tableAttrs["options"] = opts + "," + part
						} else {
							tableAttrs["options"] = part
						}
					}
				}
			}
			break // Attributes found
		}
		// If we hit a title (starts with .) continue looking before it
		if strings.HasPrefix(prevLine, ".") {
			continue
		}
		break // Hit something else
	}

	p.lineNum++ // Skip opening
	table := NewElementNode("table")
	
	// Apply attributes
	for k, v := range tableAttrs {
		table.SetAttribute(k, v)
	}

	var headerRow *Node
	var rows []*Node
	
	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])
		if line == "|===" {
			p.lineNum++
			break
		}

		if strings.HasPrefix(line, "|") {
			cells := strings.Split(line, "|")
			row := NewElementNode("tr")

			for i := 1; i < len(cells); i++ {
				cellText := strings.TrimSpace(cells[i])
				cell := NewElementNode("td")
				p.parseInlineContent(cell, cellText)
				row.AddChild(cell)
			}

			// First row is typically header logic needs refinement, but for now stick to simple logic
			// Check 'options' attribute for 'header'
			if headerRow == nil {
				// Logic for implicit header could go here, or respect options
				// For now just add as row
			}
			rows = append(rows, row)
		}
		p.lineNum++
	}

	// Determine if first row is header (often done via [options="header"] or empty line between)
	// For now, just append all as rows.
	for _, row := range rows {
		table.AddChild(row)
	}

	return table
}

func (p *parser) parseList() *Node {
	var items []*Node
	style := "unordered"

	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])

		if !p.isListItem(line) {
			break
		}

		// Determine style and marker
		if strings.HasPrefix(line, ".") {
			style = "ordered"
		} else if strings.HasPrefix(line, "*") {
			style = "unordered"
		} else if strings.HasPrefix(line, "-") {
			style = "unordered"
		} else if strings.Contains(line, "::") {
			style = "labeled"
		}

		item := p.parseListItem(line)
		if item != nil {
			items = append(items, item)
		}
		p.lineNum++
	}

	tagName := "ul"
	if style == "ordered" {
		tagName = "ol"
	} else if style == "labeled" {
		tagName = "dl"
	}

	list := NewElementNode(tagName)
	if style != "unordered" && style != "ordered" && style != "labeled" {
		list.SetAttribute("class", "list-"+style)
	}

	for _, item := range items {
		list.AddChild(item)
	}

	return list
}

func (p *parser) parseListItem(line string) *Node {
	// Remove list marker
	content := strings.TrimLeft(line, ".*-")
	content = strings.TrimSpace(content)

	// Handle labeled lists
	if strings.Contains(content, "::") {
		parts := strings.SplitN(content, "::", 2)
		item := NewElementNode("div")
		item.SetAttribute("class", "labeled-item")

		term := strings.TrimSpace(parts[0])
		dt := NewElementNode("dt")
		p.parseInlineContent(dt, term)
		item.AddChild(dt)

		desc := strings.TrimSpace(parts[1])
		dd := NewElementNode("dd")
		p.parseInlineContent(dd, desc)
		item.AddChild(dd)

		return item
	}

	item := NewElementNode("li")
	// Implicit paragraph for list item content? Or just mixed content?
	// Let's use mixed content for list item for now.
	p.parseInlineContent(item, content)
	return item
}

func (p *parser) parseAdmonition() *Node {
	line := strings.TrimSpace(p.lines[p.lineNum])
	p.lineNum++

	// Extract type
	admonitionType := "note"
	if strings.HasPrefix(line, "NOTE:") {
		admonitionType = "note"
	} else if strings.HasPrefix(line, "TIP:") {
		admonitionType = "tip"
	} else if strings.HasPrefix(line, "IMPORTANT:") {
		admonitionType = "important"
	} else if strings.HasPrefix(line, "WARNING:") {
		admonitionType = "warning"
	} else if strings.HasPrefix(line, "CAUTION:") {
		admonitionType = "caution"
	}

	content := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])

	admonition := NewElementNode("div")
	admonition.SetAttribute("class", "admonition admonition-"+admonitionType)
	// admonition.SetAttribute("type", admonitionType) // Removed, using class
	
	// Content is typically a paragraph
	para := NewElementNode("p")
	p.parseInlineContent(para, content)
	admonition.AddChild(para)
	return admonition
}

func (p *parser) parseImage() *Node {
	line := strings.TrimSpace(p.lines[p.lineNum])
	p.lineNum++

	// Parse image::path[alt,width,height] or image:path[alt]
	var src, alt string
	isBlock := strings.HasPrefix(line, "image::")
	
	if isBlock {
		line = strings.TrimPrefix(line, "image::")
	} else {
		line = strings.TrimPrefix(line, "image:")
	}

	// Extract path and attributes
	parts := strings.SplitN(line, "[", 2)
	src = strings.TrimSpace(parts[0])
	if len(parts) > 1 {
		attrs := strings.TrimSuffix(parts[1], "]")
		attrParts := strings.Split(attrs, ",")
		if len(attrParts) > 0 {
			alt = strings.TrimSpace(attrParts[0])
		}
	}

	image := NewElementNode("img")
	image.SetAttribute("src", src)
	if alt != "" {
		image.SetAttribute("alt", alt)
	}
	
	return image
}

func (p *parser) parseComponentMacro() *Node {
	line := strings.TrimSpace(p.lines[p.lineNum])
	p.lineNum++

	// Parse component::name[attrs]
	// Remove the "component::" prefix
	line = strings.TrimPrefix(line, "component::")

	// Extract component name and attributes
	parts := strings.SplitN(line, "[", 2)
	componentName := strings.TrimSpace(parts[0])

	component := NewElementNode("cms-component")
	component.SetAttribute("component-name", componentName)

	// Parse attributes if present
	if len(parts) > 1 {
		attrsStr := strings.TrimSuffix(parts[1], "]")
		
		// Parse attributes: support both key="value" and comma-separated formats
		// Split by comma, but need to be careful about commas inside quoted values
		// Simple approach: split and parse each part
		attrParts := splitAttributes(attrsStr)
		
		for _, part := range attrParts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			
			if strings.Contains(part, "=") {
				// Key-value format: key="value" or key=value
				kv := strings.SplitN(part, "=", 2)
				key := strings.TrimSpace(kv[0])
				val := strings.TrimSpace(kv[1])
				val = strings.Trim(val, "\"'") // Remove quotes
				component.SetAttribute(key, val)
			} else {
				// Positional attribute - use as both key and value, or store with generic key
				// For now, we'll skip positional-only attributes or could store with index
				// Following AsciiDoc convention, positional attributes might be role or id
				// For simplicity, we'll just store with the value as key and empty value
				// Or better: store positional attributes with numeric keys
				// Actually, let's store them with the attribute name being the value itself
				// and value being empty, or we could use a convention like "attr0", "attr1"
				// For CMS components, it's more likely to be key="value" format
			}
		}
	}

	return component
}

// splitAttributes splits attribute string by comma, respecting quoted values
func splitAttributes(attrsStr string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(attrsStr); i++ {
		c := attrsStr[i]
		
		if (c == '"' || c == '\'') && (i == 0 || attrsStr[i-1] != '\\') {
			if !inQuotes {
				inQuotes = true
				quoteChar = c
			} else if c == quoteChar {
				inQuotes = false
				quoteChar = 0
			}
			current.WriteByte(c)
		} else if c == ',' && !inQuotes {
			parts = append(parts, current.String())
			current.Reset()
		} else {
			current.WriteByte(c)
		}
	}
	
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func (p *parser) parseInlineContent(parent *Node, text string) {
	// Collect all matches with their positions
	type matchInfo struct {
		start int
		end   int
		node  *Node
	}
	var allMatches []matchInfo

	// Regex for inline formatting
	// Bold: *text*
	// Italic: _text_
	// Monospace: `text`
	// Link: http://...[text] or link:url[text]
	
	boldRegex := regexp.MustCompile(`\*([^*]+)\*`)
	italicRegex := regexp.MustCompile(`_([^_]+)_`)
	monoRegex := regexp.MustCompile("`([^`]+)`")
	// Match HTTP/HTTPS links: http://... or https://... with optional [text]
	httpLinkRegex := regexp.MustCompile(`(https?://[^\s\[\]]+)(\[([^\]]+)\])?`)
	// Match AsciiDoc link: macros: link:url[text] or link:url
	linkMacroRegex := regexp.MustCompile(`link:([^\s\[\]]+)(\[([^\]]+)\])?`)

	// Find all bold matches
	boldMatches := boldRegex.FindAllStringIndex(text, -1)
	for _, match := range boldMatches {
		content := text[match[0]+1 : match[1]-1]
		strong := NewElementNode("strong")
		strong.SetAttribute("marker", "*")
		strong.AddChild(NewTextNode(content))
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  strong,
		})
	}

	// Find all italic matches
	italicMatches := italicRegex.FindAllStringIndex(text, -1)
	for _, match := range italicMatches {
		content := text[match[0]+1 : match[1]-1]
		emphasis := NewElementNode("em")
		emphasis.SetAttribute("marker", "_")
		emphasis.AddChild(NewTextNode(content))
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  emphasis,
		})
	}

	// Find all monospace matches
	monoMatches := monoRegex.FindAllStringIndex(text, -1)
	for _, match := range monoMatches {
		content := text[match[0]+1 : match[1]-1]
		mono := NewElementNode("code")
		mono.SetAttribute("marker", "`")
		mono.AddChild(NewTextNode(content))
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  mono,
		})
	}

	// Find all HTTP/HTTPS link matches
	httpLinkMatches := httpLinkRegex.FindAllStringSubmatchIndex(text, -1)
	for _, match := range httpLinkMatches {
		href := text[match[2]:match[3]]
		linkText := href
		if match[6] > 0 && match[7] > 0 {
			linkText = text[match[6]:match[7]]
		}
		link := NewElementNode("a")
		link.SetAttribute("href", href)
		link.AddChild(NewTextNode(linkText))
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  link,
		})
	}

	// Find all link: macro matches
	linkMacroMatches := linkMacroRegex.FindAllStringSubmatchIndex(text, -1)
	for _, match := range linkMacroMatches {
		href := text[match[2]:match[3]] // URL without "link:" prefix
		linkText := href
		if match[6] > 0 && match[7] > 0 {
			linkText = text[match[6]:match[7]] // Text from [brackets]
		}
		link := NewElementNode("a")
		link.SetAttribute("href", href)
		link.AddChild(NewTextNode(linkText))
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  link,
		})
	}

	// Sort by start position
	for i := 0; i < len(allMatches)-1; i++ {
		for j := i + 1; j < len(allMatches); j++ {
			if allMatches[i].start > allMatches[j].start {
				allMatches[i], allMatches[j] = allMatches[j], allMatches[i]
			}
		}
	}

	// Build result, handling overlaps (prefer earlier matches, skip overlapped)
	lastPos := 0

	for _, match := range allMatches {
		// Overlap check
		if match.start < lastPos {
			continue 
		}

		// Add text before match
		if match.start > lastPos {
			parent.AddChild(NewTextNode(text[lastPos:match.start]))
		}
		
		// Add match
		parent.AddChild(match.node)
		lastPos = match.end
	}

	// Add remaining text
	if lastPos < len(text) {
		parent.AddChild(NewTextNode(text[lastPos:]))
	}
}

func (p *parser) isListItem(line string) bool {
	trimmed := strings.TrimSpace(line)
	// Don't match image lines
	if strings.HasPrefix(trimmed, "image::") || strings.HasPrefix(trimmed, "image:") {
		return false
	}
	return strings.HasPrefix(trimmed, "*") ||
		strings.HasPrefix(trimmed, ".") ||
		strings.HasPrefix(trimmed, "-") ||
		strings.Contains(trimmed, "::")
}

func (p *parser) isBlockTitle(line string) bool {
	// Check if a line starting with '.' is actually a title for a block
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, ".") {
		return false
	}

	// Look ahead for block delimiters
	for i := p.lineNum + 1; i < len(p.lines); i++ {
		nextLine := strings.TrimSpace(p.lines[i])
		if nextLine == "" {
			continue
		}
		// Check if it's a block delimiter
		if strings.HasPrefix(nextLine, "====") ||
			strings.HasPrefix(nextLine, "----") ||
			strings.HasPrefix(nextLine, "....") ||
			strings.HasPrefix(nextLine, "****") ||
			strings.HasPrefix(nextLine, "____") {
			return true
		}
		// If we hit something else, it's not a title
		break
	}
	return false
}

func (p *parser) isAdmonition(line string) bool {
	trimmed := strings.ToUpper(strings.TrimSpace(line))
	return strings.HasPrefix(trimmed, "NOTE:") ||
		strings.HasPrefix(trimmed, "TIP:") ||
		strings.HasPrefix(trimmed, "IMPORTANT:") ||
		strings.HasPrefix(trimmed, "WARNING:") ||
		strings.HasPrefix(trimmed, "CAUTION:")
}

// Validate attempts to parse the AsciiDoc content and returns an error if invalid
func Validate(reader io.Reader) error {
	_, err := Parse(reader)
	return err
}
