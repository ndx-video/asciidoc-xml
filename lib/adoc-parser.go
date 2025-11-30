package lib

import (
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
	anchors    map[string]*Node // Registry for anchors
}

func newParser(content string) *parser {
	lines := strings.Split(content, "\n")
	return &parser{
		lines:      lines,
		lineNum:    0,
		attributes: make(map[string]string),
		anchors:    make(map[string]*Node),
	}
}

// newSubParser creates a sub-parser with inherited attributes
func (p *parser) newSubParser(content string) *parser {
	subParser := newParser(content)
	// Inherit attributes from parent
	subParser.attributes = make(map[string]string)
	for k, v := range p.attributes {
		subParser.attributes[k] = v
	}
	subParser.doc = p.doc // Share document for attributes
	return subParser
}

func (p *parser) parse() (*Node, error) {
	p.doc = NewDocumentNode()
	p.doc.SetAttribute("doctype", "article") // Default

	// Parse header and attributes
	p.parseHeader()

	// Parse preamble (content before first section)
	// Only parse preamble if there's actually a section later
	hasSection := false
	for i := p.lineNum; i < len(p.lines); i++ {
		line := strings.TrimSpace(p.lines[i])
		if strings.HasPrefix(line, "==") {
			hasSection = true
			break
		}
	}
	
	if hasSection {
		preamble := p.parsePreamble()
		if preamble != nil && len(preamble.Children) > 0 {
			p.doc.AddChild(preamble)
		}
	}

	// Parse content (sections and remaining blocks)
	p.parseContent(p.doc, nil)

	return p.doc, nil
}

// parsePreamble parses content before the first section
func (p *parser) parsePreamble() *Node {
	preamble := NewParagraphNode()
	preamble.SetAttribute("role", "preamble")
	
	// Look ahead to find first section
	firstSectionLine := -1
	for i := p.lineNum; i < len(p.lines); i++ {
		line := strings.TrimSpace(p.lines[i])
		if strings.HasPrefix(line, "=") && !strings.HasPrefix(line, "==") {
			// This is document title, skip
			continue
		}
		if strings.HasPrefix(line, "==") {
			firstSectionLine = i
			break
		}
	}
	
	// If no section found, all content is preamble
	if firstSectionLine == -1 {
		firstSectionLine = len(p.lines)
	}
	
	// Parse content up to first section
	if firstSectionLine > p.lineNum {
		// Temporarily set lineNum to first section
		tempLineNum := p.lineNum
		p.lineNum = firstSectionLine
		
		// Parse preamble content
		subLines := p.lines[tempLineNum:firstSectionLine]
		subContent := strings.Join(subLines, "\n")
		subParser := p.newSubParser(subContent)
		subParser.lineNum = 0
		subParser.parseContent(preamble, nil)
		
		// Restore line number
		p.lineNum = firstSectionLine
	}
	
	// If preamble has no children, return nil
	if len(preamble.Children) == 0 {
		return nil
	}
	
	return preamble
}

func (p *parser) parseHeader() {
	hasHeader := false

	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])

		// Document title
		if strings.HasPrefix(line, "=") && !strings.HasPrefix(line, "==") {
			hasHeader = true
			titleText := strings.TrimSpace(strings.TrimPrefix(line, "="))
			p.doc.SetAttribute("title", titleText)
			
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

				hasHeader = true

				switch key {
				case "author":
					p.doc.SetAttribute("author", value)
				case "email":
					p.doc.SetAttribute("email", value)
				case "revnumber":
					p.doc.SetAttribute("revnumber", value)
				case "revdate":
					p.doc.SetAttribute("revdate", value)
				case "revremark":
					p.doc.SetAttribute("revremark", value)
				case "doctype":
					p.doc.SetAttribute("doctype", value)
				default:
					// Store custom attributes with : prefix to distinguish from standard ones
					p.doc.SetAttribute(":"+key, value)
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
			// Check if previous line is an attribute line and skip it
			if p.lineNum > 0 {
				prevLine := strings.TrimSpace(p.lines[p.lineNum-1])
				if strings.HasPrefix(prevLine, "[") && strings.HasSuffix(prevLine, "]") {
					// This is an attribute line for the code block, it will be consumed by parseCodeBlock
					// We don't need to do anything here as parseCodeBlock looks backwards
				}
			}
			codeBlock := p.parseCodeBlock()
			if codeBlock != nil {
				parent.AddChild(codeBlock)
			}
			continue
		}
		
		// Skip attribute lines that are immediately before code blocks
		// Check if this line is an attribute and next line is a code block delimiter
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			if p.lineNum+1 < len(p.lines) {
				nextLine := strings.TrimSpace(p.lines[p.lineNum+1])
				if strings.HasPrefix(nextLine, "----") || strings.HasPrefix(nextLine, "```") {
					// This is an attribute for a code block, skip it (parseCodeBlock will handle it)
					p.lineNum++
					continue
				}
			}
		}

		// Literal block
		if strings.HasPrefix(trimmed, "....") {
			literalBlock := p.parseLiteralBlock()
			if literalBlock != nil {
				parent.AddChild(literalBlock)
			}
			continue
		}
		
		// Skip attribute lines that are immediately before literal blocks
		// Check if this line is an attribute and next line is a literal block delimiter
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			if p.lineNum+1 < len(p.lines) {
				nextLine := strings.TrimSpace(p.lines[p.lineNum+1])
				if strings.HasPrefix(nextLine, "....") {
					// This is an attribute for a literal block, skip it (parseLiteralBlock will handle it)
					p.lineNum++
					continue
				}
			}
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

		// Block macros: include::, toc::, video::, audio::, etc.
		if strings.Contains(trimmed, "::") && strings.HasSuffix(trimmed, "]") {
			// Check for standard block macros
			if strings.HasPrefix(trimmed, "include::") {
				include := p.parseIncludeMacro()
				if include != nil {
					parent.AddChild(include)
				}
				continue
			} else if strings.HasPrefix(trimmed, "toc::") {
				toc := p.parseTOCMacro()
				if toc != nil {
					parent.AddChild(toc)
				}
				continue
			} else if strings.HasPrefix(trimmed, "video::") {
				video := p.parseVideoMacro()
				if video != nil {
					parent.AddChild(video)
				}
				continue
			} else if strings.HasPrefix(trimmed, "audio::") {
				audio := p.parseAudioMacro()
				if audio != nil {
					parent.AddChild(audio)
				}
				continue
			}
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

		// Block anchor: [[anchor-id]] or [#anchor-id]
		if strings.HasPrefix(trimmed, "[[") && strings.HasSuffix(trimmed, "]]") {
			anchorID := strings.TrimPrefix(strings.TrimSuffix(trimmed, "]]"), "[[")
			anchor := NewBlockMacroNode("anchor")
			anchor.SetAttribute("id", anchorID)
			parent.AddChild(anchor)
			p.anchors[anchorID] = anchor
			p.lineNum++
			continue
		}

		// Thematic break
		if trimmed == "'''" {
			tb := NewThematicBreakNode()
			parent.AddChild(tb)
			p.lineNum++
			continue
		}

		// Page break
		if trimmed == "<<<" {
			pb := NewPageBreakNode()
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

	// Check for section ID: [id] or [#id] before the title
	var sectionID string
	if p.lineNum > 0 {
		prevLine := strings.TrimSpace(p.lines[p.lineNum-1])
		if strings.HasPrefix(prevLine, "[") && strings.HasSuffix(prevLine, "]") {
			// Could be [id] or [#id] or [.role] or [id.role]
			idContent := prevLine[1 : len(prevLine)-1]
			if strings.HasPrefix(idContent, "#") {
				sectionID = strings.TrimPrefix(idContent, "#")
			} else if !strings.HasPrefix(idContent, ".") {
				// If it starts with a dot, it's a role, not an ID
				// Otherwise, it could be an ID
				parts := strings.Split(idContent, ".")
				if len(parts) > 0 {
					sectionID = parts[0]
				}
			}
		}
	}

	titleText := strings.TrimSpace(line)
	section := NewSectionNode(sectionLevel)
	section.SetAttribute("title", titleText)
	section.SetAttribute("marker", marker)
	
	if sectionID != "" {
		section.SetAttribute("id", sectionID)
		// Register anchor for cross-references
		p.anchors[sectionID] = section
		// Also register with underscore prefix for section references
		p.anchors["_"+sectionID] = section
	}
	
	// Generate automatic ID from title if no ID specified
	if sectionID == "" {
		autoID := p.generateSectionID(titleText)
		section.SetAttribute("id", autoID)
		p.anchors[autoID] = section
		p.anchors["_"+autoID] = section
	}
	
	// Add title as text content (converters will handle rendering)
	section.AddChild(NewTextNode(titleText))

	p.lineNum++

	// Parse section content, stopping at sections at same or higher level
	p.parseContent(section, &sectionLevel)

	return section
}

// generateSectionID creates an ID from a section title
func (p *parser) generateSectionID(title string) string {
	// Convert to lowercase, replace spaces with underscores, remove special chars
	id := strings.ToLower(title)
	id = strings.ReplaceAll(id, " ", "_")
	// Remove special characters, keep only alphanumeric and underscores
	var result strings.Builder
	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// getAllAttributes returns all attributes including built-in ones
func (p *parser) getAllAttributes() map[string]string {
	attrs := make(map[string]string)
	
	// Copy document attributes
	for k, v := range p.doc.Attributes {
		attrs[k] = v
	}
	
	// Copy parser attributes (may override document attributes)
	for k, v := range p.attributes {
		attrs[k] = v
		// Also add with : prefix for custom attributes
		if k != "author" && k != "email" && k != "revnumber" && k != "revdate" && k != "revremark" && k != "doctype" && k != "title" {
			attrs[":"+k] = v
		}
	}
	
	return attrs
}

func (p *parser) parseParagraph() *Node {
	var lines []string

	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])

		if line == "" {
			break
		}

		// Check for attribute assignment in body: :attr-name: value
		if strings.HasPrefix(line, ":") && strings.Contains(line, ":") {
			parts := strings.SplitN(line[1:], ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				p.attributes[key] = value
				// Also set on document for built-in attributes
				if key == "author" || key == "email" || key == "revnumber" || key == "revdate" || key == "revremark" || key == "doctype" {
					p.doc.SetAttribute(key, value)
				} else {
					p.doc.SetAttribute(":"+key, value)
				}
				p.lineNum++
				continue
			}
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
	// Substitute attributes in text
	text = SubstituteAttributes(text, p.getAllAttributes())
	para := NewParagraphNode()
	p.parseInlineContent(para, text)
	return para
}

func (p *parser) parseCodeBlock() *Node {
	// Parse attributes from previous line(s) if present
	var language, title, role string
	// Look back for title and attributes (skip empty lines)
	for i := p.lineNum - 1; i >= 0 && i >= p.lineNum-3; i-- {
		prevLine := strings.TrimSpace(p.lines[i])
		if prevLine == "" {
			continue
		}
		if strings.HasPrefix(prevLine, "[") && strings.HasSuffix(prevLine, "]") {
			attrs := prevLine[1 : len(prevLine)-1]
			parts := strings.Split(attrs, ",")
			// Format is [source,language] or [language] or [mermaid] or [role="mermaid"]
			// Check for mermaid role first
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part == "mermaid" {
					role = "mermaid"
				} else if strings.HasPrefix(part, "role=") {
					// Handle role="mermaid" or role='mermaid'
					roleValue := strings.TrimSpace(strings.TrimPrefix(part, "role="))
					roleValue = strings.Trim(roleValue, "\"'")
					if roleValue == "mermaid" {
						role = "mermaid"
					}
				}
			}
			if len(parts) > 1 {
				// [source,go] format
				language = strings.TrimSpace(parts[1])
			} else if len(parts) > 0 {
				// [go] format or [mermaid]
				firstPart := strings.TrimSpace(parts[0])
				if firstPart != "source" && firstPart != "mermaid" {
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

	codeBlock := NewCodeBlockNode()
	if role != "" {
		codeBlock.SetAttribute("role", role)
	}
	if language != "" {
		codeBlock.SetAttribute("language", language)
	}
	if title != "" {
		codeBlock.SetAttribute("title", title)
	}
	codeBlock.AddChild(NewTextNode(strings.Join(content, "\n")))
	return codeBlock
}

func (p *parser) parseLiteralBlock() *Node {
	// Parse attributes from previous line(s) if present
	var role string
	// Look back for attributes (skip empty lines)
	for i := p.lineNum - 1; i >= 0 && i >= p.lineNum-3; i-- {
		prevLine := strings.TrimSpace(p.lines[i])
		if prevLine == "" {
			continue
		}
		if strings.HasPrefix(prevLine, "[") && strings.HasSuffix(prevLine, "]") {
			attrs := prevLine[1 : len(prevLine)-1]
			parts := strings.Split(attrs, ",")
			// Check for mermaid role
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part == "mermaid" {
					role = "mermaid"
					break
				} else if strings.HasPrefix(part, "role=") {
					// Handle role="mermaid" or role='mermaid'
					roleValue := strings.TrimSpace(strings.TrimPrefix(part, "role="))
					roleValue = strings.Trim(roleValue, "\"'")
					if roleValue == "mermaid" {
						role = "mermaid"
						break
					}
				}
			}
			break
		}
		if strings.HasPrefix(prevLine, ".") {
			break // Title found, stop looking
		}
	}

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
	literalBlock := NewLiteralBlockNode()
	if role != "" {
		literalBlock.SetAttribute("role", role)
	}
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
	subParser := p.newSubParser(subContent)
	example := NewExampleNode()
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

	sidebar := NewSidebarNode()
	if title != "" {
		sidebar.SetAttribute("title", title)
	}
	
	subContent := strings.Join(contentLines, "\n")
	subParser := p.newSubParser(subContent)
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

	quote := NewQuoteNode()
	if attribution != "" {
		quote.SetAttribute("attribution", attribution)
	}
	if citation != "" {
		quote.SetAttribute("citation", citation)
	}

	subContent := strings.Join(contentLines, "\n")
	subParser := p.newSubParser(subContent)
	subParser.parseContent(quote, nil)

	return quote
}

func (p *parser) parseVerseBlock() *Node {
	// Check for [verse] attribute on previous line
	var title, attribution string
	if p.lineNum > 0 {
		prevLine := strings.TrimSpace(p.lines[p.lineNum-1])
		if strings.HasPrefix(prevLine, "[verse") {
			// Parse attributes if present
			if strings.Contains(prevLine, ",") {
				parts := strings.SplitN(prevLine[1:len(prevLine)-1], ",", 2)
				if len(parts) > 1 {
					attribution = strings.TrimSpace(parts[1])
				}
			}
		}
		// Check for title before that
		if p.lineNum > 1 {
			titleLine := strings.TrimSpace(p.lines[p.lineNum-2])
			if strings.HasPrefix(titleLine, ".") {
				title = strings.TrimPrefix(titleLine, ".")
			}
		}
	}

	p.lineNum++ // Skip opening
	var contentLines []string

	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])
		if line == "____" {
			p.lineNum++
			break
		}
		contentLines = append(contentLines, p.lines[p.lineNum])
		p.lineNum++
	}

	verse := NewVerseBlockNode()
	if title != "" {
		verse.SetAttribute("title", title)
	}
	if attribution != "" {
		verse.SetAttribute("attribution", attribution)
	}

	// Verse blocks preserve line breaks, so we parse content but preserve structure
	subContent := strings.Join(contentLines, "\n")
	subParser := p.newSubParser(subContent)
	subParser.parseContent(verse, nil)

	// Check for attribution after closing delimiter
	if p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])
		if strings.HasPrefix(line, "--") {
			attribution = strings.TrimPrefix(line, "--")
			attribution = strings.TrimSpace(attribution)
			if attribution != "" {
				verse.SetAttribute("attribution", attribution)
			}
			p.lineNum++
		}
	}

	return verse
}

func (p *parser) parseOpenBlock() *Node {
	// Parse attributes from previous line(s)
	attrs := make(map[string]string)
	if p.lineNum > 0 {
		prevLine := strings.TrimSpace(p.lines[p.lineNum-1])
		if strings.HasPrefix(prevLine, "[") && strings.HasSuffix(prevLine, "]") {
			attrContent := prevLine[1 : len(prevLine)-1]
			parts := strings.Split(attrContent, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.Contains(part, "=") {
					kv := strings.SplitN(part, "=", 2)
					key := strings.TrimSpace(kv[0])
					val := strings.TrimSpace(kv[1])
					if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
						val = val[1 : len(val)-1]
					}
					attrs[key] = val
				} else if strings.HasPrefix(part, "#") {
					attrs["id"] = strings.TrimPrefix(part, "#")
				} else if strings.HasPrefix(part, ".") {
					attrs["role"] = strings.TrimPrefix(part, ".")
				}
			}
		}
	}

	p.lineNum++ // Skip opening
	var contentLines []string

	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])
		if line == "--" {
			p.lineNum++
			break
		}
		contentLines = append(contentLines, p.lines[p.lineNum])
		p.lineNum++
	}

	openBlock := NewOpenBlockNode()
	for k, v := range attrs {
		openBlock.SetAttribute(k, v)
	}

	subContent := strings.Join(contentLines, "\n")
	subParser := p.newSubParser(subContent)
	subParser.parseContent(openBlock, nil)

	return openBlock
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
					// Support modern HTML attributes: keys can contain :, @, ., _, -
					// Values can contain @, $, {, }, (, ) and other special chars
					kv := strings.SplitN(part, "=", 2)
					key := strings.TrimSpace(kv[0])
					val := strings.TrimSpace(kv[1])
					// Remove quotes only from the outer edges (handles both "value" and 'value')
					if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
						val = val[1 : len(val)-1]
					}
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
	table := NewTableNode()
	
	// Apply attributes
	for k, v := range tableAttrs {
		table.SetAttribute(k, v)
	}

	var rows []*Node
	
	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])
		if line == "|===" {
			p.lineNum++
			break
		}

		if strings.HasPrefix(line, "|") {
			cells := strings.Split(line, "|")
			row := NewTableRowNode()
			
			// Check if this is a header row (first non-empty row or explicitly marked)
			isHeaderRow := len(rows) == 0
			if opts, ok := tableAttrs["options"]; ok && strings.Contains(opts, "header") {
				// Header rows are explicitly marked
				isHeaderRow = len(rows) == 0 || (len(rows) == 1 && strings.Contains(tableAttrs["options"], "header"))
			}
			if isHeaderRow {
				row.SetAttribute("role", "header")
			}

			for i := 1; i < len(cells); i++ {
				cellText := strings.TrimSpace(cells[i])
				
				// Parse cell attributes and alignment: [.class]#text# or [align=left]#text# or |^text| or |vtext|
				var cellAttrs map[string]string
				var alignment string
				var actualText string
				
				// Check for alignment markers: ^ (top), v (bottom), or default (middle)
				if strings.HasPrefix(cellText, "^") {
					alignment = "top"
					actualText = cellText[1:]
				} else if strings.HasPrefix(cellText, "v") {
					alignment = "bottom"
					actualText = cellText[1:]
				} else {
					actualText = cellText
				}
				
				// Check for cell attributes: [.class]#text# or [colspan=2]#text#
				if strings.Contains(actualText, "#") && strings.Contains(actualText, "[") {
					attrStart := strings.Index(actualText, "[")
					attrEnd := strings.Index(actualText[attrStart:], "]")
					if attrEnd > 0 {
						attrEnd += attrStart
						attrStr := actualText[attrStart+1 : attrEnd]
						contentStart := strings.Index(actualText[attrEnd:], "#")
						contentEnd := strings.LastIndex(actualText, "#")
						if contentStart >= 0 && contentEnd > attrEnd {
							contentStart += attrEnd
							cellAttrs = make(map[string]string)
							// Parse attributes
							if strings.HasPrefix(attrStr, ".") {
								cellAttrs["role"] = strings.TrimPrefix(attrStr, ".")
							} else if strings.HasPrefix(attrStr, "#") {
								cellAttrs["id"] = strings.TrimPrefix(attrStr, "#")
							} else if strings.Contains(attrStr, "=") {
								// Key-value attributes like colspan=2
								parts := strings.Split(attrStr, ",")
								for _, part := range parts {
									part = strings.TrimSpace(part)
									if strings.Contains(part, "=") {
										kv := strings.SplitN(part, "=", 2)
										key := strings.TrimSpace(kv[0])
										val := strings.TrimSpace(kv[1])
										if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
											val = val[1 : len(val)-1]
										}
										cellAttrs[key] = val
									}
								}
							} else {
								// Could be id.role format
								parts := strings.Split(attrStr, ".")
								if len(parts) > 0 {
									cellAttrs["id"] = parts[0]
								}
								if len(parts) > 1 {
									cellAttrs["role"] = strings.Join(parts[1:], ".")
								}
							}
							// Extract content between # markers
							actualText = actualText[contentStart+1:contentEnd]
						}
					}
				}
				
				// Check for horizontal alignment in cell text: left, center, right
				// AsciiDoc uses <, ^, > for left, center, right alignment
				if strings.HasPrefix(actualText, "<") {
					alignment = "left"
					actualText = actualText[1:]
				} else if strings.HasPrefix(actualText, "^") {
					if alignment == "" {
						alignment = "center"
					}
					actualText = actualText[1:]
				} else if strings.HasPrefix(actualText, ">") {
					alignment = "right"
					actualText = actualText[1:]
				}
				
				// Substitute attributes in cell text
				actualText = SubstituteAttributes(actualText, p.getAllAttributes())
				cell := NewTableCellNode()
				
				// Set cell attributes
				if alignment != "" {
					cell.SetAttribute("align", alignment)
				}
				if cellAttrs != nil {
					for k, v := range cellAttrs {
						cell.SetAttribute(k, v)
					}
				}
				
				p.parseInlineContent(cell, actualText)
				row.AddChild(cell)
			}

			rows = append(rows, row)
		}
		p.lineNum++
	}

	// Append all rows
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

		// Check for list continuation: + on separate line
		if line == "+" {
			// This is a continuation marker - the next content should be added to the last item
			p.lineNum++
			if len(items) > 0 {
				// Parse the next block and add it to the last item
				lastItem := items[len(items)-1]
				// The next line should be a block (code, table, etc.)
				// We'll handle this by parsing content into the last item
				if p.lineNum < len(p.lines) {
					nextLine := strings.TrimSpace(p.lines[p.lineNum])
					// Check what type of block follows
					if strings.HasPrefix(nextLine, "----") || strings.HasPrefix(nextLine, "```") {
						codeBlock := p.parseCodeBlock()
						if codeBlock != nil {
							lastItem.AddChild(codeBlock)
						}
						continue
					} else if strings.HasPrefix(nextLine, "|===") {
						table := p.parseTable()
						if table != nil {
							lastItem.AddChild(table)
						}
						continue
					} else if strings.HasPrefix(nextLine, "====") {
						example := p.parseExampleBlock()
						if example != nil {
							lastItem.AddChild(example)
						}
						continue
					} else if strings.HasPrefix(nextLine, "....") {
						literalBlock := p.parseLiteralBlock()
						if literalBlock != nil {
							lastItem.AddChild(literalBlock)
						}
						continue
					} else {
						// Regular paragraph continuation
						para := p.parseParagraph()
						if para != nil {
							lastItem.AddChild(para)
						}
						continue
					}
				}
			}
			continue
		}

		// Check for callout list: <1>, <2>, etc.
		if strings.HasPrefix(line, "<") && strings.HasSuffix(line, ">") {
			calloutNum := strings.Trim(line, "<>")
			item := NewListItemNode()
			item.SetAttribute("callout", calloutNum)
			// Callout lists don't have content, just the number
			items = append(items, item)
			p.lineNum++
			continue
		}

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

	list := NewListNode()
	list.SetAttribute("style", style)

	for _, item := range items {
		list.AddChild(item)
	}

	return list
}

func (p *parser) parseListItem(line string) *Node {
	// Remove list marker
	content := strings.TrimLeft(line, ".*-")
	content = strings.TrimSpace(content)

	item := NewListItemNode()
	
	// Handle labeled lists
	if strings.Contains(content, "::") {
		parts := strings.SplitN(content, "::", 2)
		term := strings.TrimSpace(parts[0])
		termNode := NewParagraphNode()
		p.parseInlineContent(termNode, term)
		item.SetAttribute("term", term)
		item.AddChild(termNode)

		desc := strings.TrimSpace(parts[1])
		descNode := NewParagraphNode()
		p.parseInlineContent(descNode, desc)
		item.AddChild(descNode)
	} else {
		// Regular list item content
		p.parseInlineContent(item, content)
	}
	
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

	admonition := NewAdmonitionNode()
	admonition.SetAttribute("type", admonitionType)
	
	// Content is typically a paragraph
	para := NewParagraphNode()
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

	image := NewBlockMacroNode("image")
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

	component := NewBlockMacroNode("component")
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
				// Support modern HTML attributes: keys can contain :, @, ., _, -
				// Values can contain @, $, {, }, (, ) and other special chars
				kv := strings.SplitN(part, "=", 2)
				key := strings.TrimSpace(kv[0])
				val := strings.TrimSpace(kv[1])
				// Remove quotes only from the outer edges (handles both "value" and 'value')
				if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
					val = val[1 : len(val)-1]
				}
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
	// Superscript: ^text^
	// Subscript: ~text~
	// Highlight: #text#
	// Inline passthrough: +text+
	// Link: http://...[text] or link:url[text]
	// Inline macros: macroname:[text] or macroname:text[]
	
	boldRegex := regexp.MustCompile(`\*([^*]+)\*`)
	italicRegex := regexp.MustCompile(`_([^_]+)_`)
	monoRegex := regexp.MustCompile("`([^`]+)`")
	superscriptRegex := regexp.MustCompile(`\^([^^]+)\^`)
	subscriptRegex := regexp.MustCompile(`~([^~]+)~`)
	highlightRegex := regexp.MustCompile(`#([^#]+)#`)
	passthroughInlineRegex := regexp.MustCompile(`\+([^+]+)\+`)
	// Match HTTP/HTTPS links: http://... or https://... with optional [text]
	httpLinkRegex := regexp.MustCompile(`(https?://[^\s\[\]]+)(\[([^\]]+)\])?`)
	// Match AsciiDoc link: macros: link:url[text] or link:url
	linkMacroRegex := regexp.MustCompile(`link:([^\s\[\]]+)(\[([^\]]+)\])?`)
	// Match inline macros: macroname:[text] or macroname:target[text]
	// This matches: word characters, colon, optional target (non-bracket chars), brackets with content
	// Exclude link: which is handled separately
	inlineMacroRegex := regexp.MustCompile(`(\w+):([^\[]*)\[([^\]]+)\]`)
	// Match cross-references: <<anchor-id>> or xref:anchor-id[]
	xrefRegex := regexp.MustCompile(`<<([^>>]+)>>`)
	xrefMacroRegex := regexp.MustCompile(`xref:([^\s\[\]]+)(\[([^\]]+)\])?`)
	// Match inline anchor: [#anchor-id]
	inlineAnchorRegex := regexp.MustCompile(`\[#([^\]]+)\]`)
	// Match footnotes: footnote:[text] or footnote:ref[text]
	footnoteRegex := regexp.MustCompile(`footnote:([^\s\[\]]*)(\[([^\]]+)\])?`)
	footnoterefRegex := regexp.MustCompile(`footnoteref:([^\s\[\]]+)(\[([^\]]+)\])?`)

	// Find all bold matches
	boldMatches := boldRegex.FindAllStringIndex(text, -1)
	for _, match := range boldMatches {
		content := text[match[0]+1 : match[1]-1]
		strong := NewBoldNode()
		// Recursively parse nested content
		p.parseInlineContent(strong, content)
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
		emphasis := NewItalicNode()
		// Recursively parse nested content
		p.parseInlineContent(emphasis, content)
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
		mono := NewMonospaceNode()
		// Recursively parse nested content
		p.parseInlineContent(mono, content)
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  mono,
		})
	}

	// Find all superscript matches
	superscriptMatches := superscriptRegex.FindAllStringIndex(text, -1)
	for _, match := range superscriptMatches {
		content := text[match[0]+1 : match[1]-1]
		sup := NewSuperscriptNode()
		// Recursively parse nested content
		p.parseInlineContent(sup, content)
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  sup,
		})
	}

	// Find all subscript matches
	subscriptMatches := subscriptRegex.FindAllStringIndex(text, -1)
	for _, match := range subscriptMatches {
		content := text[match[0]+1 : match[1]-1]
		sub := NewSubscriptNode()
		// Recursively parse nested content
		p.parseInlineContent(sub, content)
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  sub,
		})
	}

	// Find all highlight matches
	highlightMatches := highlightRegex.FindAllStringIndex(text, -1)
	for _, match := range highlightMatches {
		content := text[match[0]+1 : match[1]-1]
		highlight := NewHighlightNode()
		// Recursively parse nested content
		p.parseInlineContent(highlight, content)
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  highlight,
		})
	}

	// Find all inline passthrough matches
	passthroughMatches := passthroughInlineRegex.FindAllStringIndex(text, -1)
	for _, match := range passthroughMatches {
		content := text[match[0]+1 : match[1]-1]
		passthrough := NewPassthroughNode(content)
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  passthrough,
		})
	}

	// Find all inline macro matches (before HTTP links to avoid conflicts)
	inlineMacroMatches := inlineMacroRegex.FindAllStringSubmatchIndex(text, -1)
	for _, match := range inlineMacroMatches {
		macroName := text[match[2]:match[3]]
		target := text[match[4]:match[5]]
		macroText := text[match[6]:match[7]]
		
		// Skip if it's a link: macro (handled separately)
		if macroName == "link" {
			continue
		}
		
		macro := NewInlineMacroNode(macroName)
		if target != "" {
			macro.SetAttribute("target", target)
		}
		macro.AddChild(NewTextNode(macroText))
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  macro,
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
		link := NewLinkNode()
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
		link := NewLinkNode()
		link.SetAttribute("href", href)
		link.AddChild(NewTextNode(linkText))
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  link,
		})
	}

	// Find all cross-reference matches: <<anchor-id>>
	xrefMatches := xrefRegex.FindAllStringSubmatchIndex(text, -1)
	for _, match := range xrefMatches {
		anchorID := text[match[2]:match[3]]
		xref := NewInlineMacroNode("xref")
		xref.SetAttribute("target", anchorID)
		// Default text is the anchor ID, but can be overridden
		xref.AddChild(NewTextNode(anchorID))
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  xref,
		})
	}

	// Find all xref: macro matches: xref:anchor-id[] or xref:anchor-id[text]
	xrefMacroMatches := xrefMacroRegex.FindAllStringSubmatchIndex(text, -1)
	for _, match := range xrefMacroMatches {
		anchorID := text[match[2]:match[3]]
		xrefText := anchorID
		if match[6] > 0 && match[7] > 0 {
			xrefText = text[match[6]:match[7]] // Text from [brackets]
		}
		xref := NewInlineMacroNode("xref")
		xref.SetAttribute("target", anchorID)
		xref.AddChild(NewTextNode(xrefText))
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  xref,
		})
	}

	// Find all inline anchor matches: [#anchor-id]
	inlineAnchorMatches := inlineAnchorRegex.FindAllStringSubmatchIndex(text, -1)
	for _, match := range inlineAnchorMatches {
		anchorID := text[match[2]:match[3]]
		anchor := NewInlineMacroNode("anchor")
		anchor.SetAttribute("id", anchorID)
		p.anchors[anchorID] = anchor
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  anchor,
		})
	}

	// Find all footnote matches: footnote:[text] or footnote:ref[text]
	footnoteMatches := footnoteRegex.FindAllStringSubmatchIndex(text, -1)
	for _, match := range footnoteMatches {
		ref := text[match[2]:match[3]]
		footnoteText := ""
		if match[6] > 0 && match[7] > 0 {
			footnoteText = text[match[6]:match[7]]
		}
		footnote := NewInlineMacroNode("footnote")
		if ref != "" {
			footnote.SetAttribute("ref", ref)
		}
		if footnoteText != "" {
			footnote.AddChild(NewTextNode(footnoteText))
		}
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  footnote,
		})
	}

	// Find all footnoteref matches: footnoteref:ref[]
	footnoterefMatches := footnoterefRegex.FindAllStringSubmatchIndex(text, -1)
	for _, match := range footnoterefMatches {
		ref := text[match[2]:match[3]]
		footnoterefText := ref
		if match[6] > 0 && match[7] > 0 {
			footnoterefText = text[match[6]:match[7]]
		}
		footnoteref := NewInlineMacroNode("footnoteref")
		footnoteref.SetAttribute("ref", ref)
		footnoteref.AddChild(NewTextNode(footnoterefText))
		allMatches = append(allMatches, matchInfo{
			start: match[0],
			end:   match[1],
			node:  footnoteref,
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

		// Add text before match (with attribute substitution)
		if match.start > lastPos {
			textBefore := text[lastPos:match.start]
			textBefore = SubstituteAttributes(textBefore, p.getAllAttributes())
			parent.AddChild(NewTextNode(textBefore))
		}
		
		// Add match
		parent.AddChild(match.node)
		lastPos = match.end
	}

	// Add remaining text (with attribute substitution)
	if lastPos < len(text) {
		textRemaining := text[lastPos:]
		textRemaining = SubstituteAttributes(textRemaining, p.getAllAttributes())
		parent.AddChild(NewTextNode(textRemaining))
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

// parseIncludeMacro parses include::file[] directive
func (p *parser) parseIncludeMacro() *Node {
	line := strings.TrimSpace(p.lines[p.lineNum])
	p.lineNum++
	
	// Parse include::file[lines=1..5] or include::file[]
	line = strings.TrimPrefix(line, "include::")
	parts := strings.SplitN(line, "[", 2)
	filePath := strings.TrimSpace(parts[0])
	
	include := NewBlockMacroNode("include")
	include.SetAttribute("target", filePath)
	
	if len(parts) > 1 {
		attrs := strings.TrimSuffix(parts[1], "]")
		// Parse attributes like lines=1..5, tags=tag1;tag2
		attrParts := strings.Split(attrs, ",")
		for _, part := range attrParts {
			part = strings.TrimSpace(part)
			if strings.Contains(part, "=") {
				kv := strings.SplitN(part, "=", 2)
				key := strings.TrimSpace(kv[0])
				val := strings.TrimSpace(kv[1])
				include.SetAttribute(key, val)
			}
		}
	}
	
	return include
}

// parseTOCMacro parses toc::[] directive
func (p *parser) parseTOCMacro() *Node {
	line := strings.TrimSpace(p.lines[p.lineNum])
	p.lineNum++
	
	// Parse toc::[] or toc::[levels=2]
	line = strings.TrimPrefix(line, "toc::")
	toc := NewBlockMacroNode("toc")
	
	if strings.HasPrefix(line, "[") {
		attrs := strings.TrimSuffix(strings.TrimPrefix(line, "["), "]")
		attrParts := strings.Split(attrs, ",")
		for _, part := range attrParts {
			part = strings.TrimSpace(part)
			if strings.Contains(part, "=") {
				kv := strings.SplitN(part, "=", 2)
				key := strings.TrimSpace(kv[0])
				val := strings.TrimSpace(kv[1])
				toc.SetAttribute(key, val)
			}
		}
	}
	
	return toc
}

// parseVideoMacro parses video::url[] directive
func (p *parser) parseVideoMacro() *Node {
	line := strings.TrimSpace(p.lines[p.lineNum])
	p.lineNum++
	
	// Parse video::url[] or video::url[width=640,height=360]
	line = strings.TrimPrefix(line, "video::")
	parts := strings.SplitN(line, "[", 2)
	url := strings.TrimSpace(parts[0])
	
	video := NewBlockMacroNode("video")
	video.SetAttribute("target", url)
	
	if len(parts) > 1 {
		attrs := strings.TrimSuffix(parts[1], "]")
		attrParts := strings.Split(attrs, ",")
		for _, part := range attrParts {
			part = strings.TrimSpace(part)
			if strings.Contains(part, "=") {
				kv := strings.SplitN(part, "=", 2)
				key := strings.TrimSpace(kv[0])
				val := strings.TrimSpace(kv[1])
				video.SetAttribute(key, val)
			}
		}
	}
	
	return video
}

// parseAudioMacro parses audio::url[] directive
func (p *parser) parseAudioMacro() *Node {
	line := strings.TrimSpace(p.lines[p.lineNum])
	p.lineNum++
	
	// Parse audio::url[] or audio::url[autoplay]
	line = strings.TrimPrefix(line, "audio::")
	parts := strings.SplitN(line, "[", 2)
	url := strings.TrimSpace(parts[0])
	
	audio := NewBlockMacroNode("audio")
	audio.SetAttribute("target", url)
	
	if len(parts) > 1 {
		attrs := strings.TrimSuffix(parts[1], "]")
		attrParts := strings.Split(attrs, ",")
		for _, part := range attrParts {
			part = strings.TrimSpace(part)
			if strings.Contains(part, "=") {
				kv := strings.SplitN(part, "=", 2)
				key := strings.TrimSpace(kv[0])
				val := strings.TrimSpace(kv[1])
				audio.SetAttribute(key, val)
			} else {
				// Boolean attribute
				audio.SetAttribute(part, "true")
			}
		}
	}
	
	return audio
}

// Validate attempts to parse the AsciiDoc content and returns an error if invalid
func Validate(reader io.Reader) error {
	_, err := Parse(reader)
	return err
}
