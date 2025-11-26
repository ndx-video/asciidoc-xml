package converter

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	asciidocxml "asciidoc-xml"
)

func init() {
	// Silence libasciidoc logging if used
	logrus.SetLevel(logrus.ErrorLevel)
}

// Convert converts AsciiDoc content to our custom XML format
// This is a text-based parser that parses AsciiDoc syntax directly
func Convert(reader io.Reader) (*asciidocxml.Document, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %w", err)
	}

	parser := newParser(string(content))
	return parser.parse()
}

// ConvertToXML converts AsciiDoc to XML string
func ConvertToXML(reader io.Reader) (string, error) {
	doc, err := Convert(reader)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return "", fmt.Errorf("failed to encode XML: %w", err)
	}

	return buf.String(), nil
}

// parser is a text-based AsciiDoc parser
type parser struct {
	lines     []string
	lineNum   int
	doc       *asciidocxml.Document
	header    *asciidocxml.Header
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

func (p *parser) parse() (*asciidocxml.Document, error) {
	p.doc = &asciidocxml.Document{
		DocType: "article",
		Content: asciidocxml.Content{},
	}

	// Parse header and attributes
	p.parseHeader()

	// Parse content
	p.doc.Content = p.parseContent()

	if p.header != nil {
		p.doc.Header = p.header
	}

	return p.doc, nil
}

func (p *parser) parseHeader() {
	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])
		
		// Document title
		if strings.HasPrefix(line, "=") && !strings.HasPrefix(line, "==") {
			if p.header == nil {
				p.header = &asciidocxml.Header{}
			}
			p.header.Title = strings.TrimSpace(strings.TrimPrefix(line, "="))
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
				
				if p.header == nil {
					p.header = &asciidocxml.Header{}
				}
				
				switch key {
				case "author":
					p.header.Authors = append(p.header.Authors, asciidocxml.Author{
						Name: value,
					})
				case "email":
					if len(p.header.Authors) > 0 {
						p.header.Authors[len(p.header.Authors)-1].Email = value
					}
				case "revnumber":
					if p.header.Revision == nil {
						p.header.Revision = &asciidocxml.Revision{}
					}
					p.header.Revision.Number = value
				case "revdate":
					if p.header.Revision == nil {
						p.header.Revision = &asciidocxml.Revision{}
					}
					p.header.Revision.Date = value
				case "revremark":
					if p.header.Revision == nil {
						p.header.Revision = &asciidocxml.Revision{}
					}
					p.header.Revision.Remark = value
				case "doctype":
					p.doc.DocType = value
				default:
					p.header.Attributes = append(p.header.Attributes, asciidocxml.Attribute{
						Name:  key,
						Value: value,
					})
				}
			}
			p.lineNum++
			continue
		}

		// Empty line after header ends header parsing
		if line == "" {
			p.lineNum++
			if p.header != nil && p.header.Title != "" {
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

func (p *parser) parseContent() asciidocxml.Content {
	content := asciidocxml.Content{}
	
	for p.lineNum < len(p.lines) {
		line := p.lines[p.lineNum]
		trimmed := strings.TrimSpace(line)
		
		if trimmed == "" {
			p.lineNum++
			continue
		}

		// Section
		if strings.HasPrefix(trimmed, "=") {
			section := p.parseSection()
			if section != nil {
				content.Items = append(content.Items, asciidocxml.ContentItem{
					Section: section,
				})
			}
			continue
		}

		// Code block
		if strings.HasPrefix(trimmed, "----") || strings.HasPrefix(trimmed, "```") {
			codeBlock := p.parseCodeBlock()
			if codeBlock != nil {
				content.Items = append(content.Items, asciidocxml.ContentItem{
					CodeBlock: codeBlock,
				})
			}
			continue
		}

		// Literal block
		if strings.HasPrefix(trimmed, "....") {
			literalBlock := p.parseLiteralBlock()
			if literalBlock != nil {
				content.Items = append(content.Items, asciidocxml.ContentItem{
					LiteralBlock: literalBlock,
				})
			}
			continue
		}

		// Example block
		if strings.HasPrefix(trimmed, "====") {
			example := p.parseExampleBlock()
			if example != nil {
				content.Items = append(content.Items, asciidocxml.ContentItem{
					Example: example,
				})
			}
			continue
		}

		// Sidebar
		if strings.HasPrefix(trimmed, "****") {
			sidebar := p.parseSidebar()
			if sidebar != nil {
				content.Items = append(content.Items, asciidocxml.ContentItem{
					Sidebar: sidebar,
				})
			}
			continue
		}

		// Quote
		if strings.HasPrefix(trimmed, "____") {
			quote := p.parseQuote()
			if quote != nil {
				content.Items = append(content.Items, asciidocxml.ContentItem{
					Quote: quote,
				})
			}
			continue
		}

		// Table
		if strings.HasPrefix(trimmed, "|===") {
			table := p.parseTable()
			if table != nil {
				content.Items = append(content.Items, asciidocxml.ContentItem{
					Table: table,
				})
			}
			continue
		}

		// List
		if p.isListItem(trimmed) {
			list := p.parseList()
			if list != nil {
				content.Items = append(content.Items, asciidocxml.ContentItem{
					List: list,
				})
			}
			continue
		}

		// Admonition
		if p.isAdmonition(trimmed) {
			admonition := p.parseAdmonition()
			if admonition != nil {
				content.Items = append(content.Items, asciidocxml.ContentItem{
					Admonition: admonition,
				})
			}
			continue
		}

		// Image
		if strings.HasPrefix(trimmed, "image::") || strings.HasPrefix(trimmed, "image:") {
			image := p.parseImage()
			if image != nil {
				content.Items = append(content.Items, asciidocxml.ContentItem{
					Image: image,
				})
			}
			continue
		}

		// Thematic break
		if trimmed == "'''" {
			content.Items = append(content.Items, asciidocxml.ContentItem{
				ThematicBreak: &asciidocxml.ThematicBreak{},
			})
			p.lineNum++
			continue
		}

		// Page break
		if trimmed == "<<<" {
			content.Items = append(content.Items, asciidocxml.ContentItem{
				PageBreak: &asciidocxml.PageBreak{},
			})
			p.lineNum++
			continue
		}

		// Paragraph
		para := p.parseParagraph()
		if para != nil {
			content.Items = append(content.Items, asciidocxml.ContentItem{
				Paragraph: para,
			})
		} else {
			// If paragraph parsing fails or returns nil, increment to avoid infinite loop
			p.lineNum++
		}
	}

	return content
}

func (p *parser) parseSection() *asciidocxml.Section {
	line := strings.TrimSpace(p.lines[p.lineNum])
	level := 0
	for strings.HasPrefix(line, "=") {
		level++
		line = line[1:]
	}
	
	title := strings.TrimSpace(line)
	section := &asciidocxml.Section{
		Level: level,
		Title: p.parseInlineContent(title),
	}
	
	p.lineNum++
	
	// Parse section content
	section.Content = p.parseContentUntilSection(level)
	
	return section
}

func (p *parser) parseContentUntilSection(maxLevel int) asciidocxml.Content {
	content := asciidocxml.Content{}
	
	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])
		
		// Check if we hit a section at same or higher level
		if strings.HasPrefix(line, "=") {
			level := 0
			for strings.HasPrefix(line, "=") {
				level++
				line = line[1:]
			}
			if level <= maxLevel {
				break
			}
		}
		
		// Parse content (similar to parseContent but stop at sections)
		trimmed := strings.TrimSpace(p.lines[p.lineNum])
		
		if trimmed == "" {
			p.lineNum++
			continue
		}

		// Code block
		if strings.HasPrefix(trimmed, "----") {
			codeBlock := p.parseCodeBlock()
			if codeBlock != nil {
				content.Items = append(content.Items, asciidocxml.ContentItem{
					CodeBlock: codeBlock,
				})
			}
			continue
		}

		// Paragraph
		para := p.parseParagraph()
		if para != nil {
			content.Items = append(content.Items, asciidocxml.ContentItem{
				Paragraph: para,
			})
		} else {
			// If paragraph parsing fails or returns nil, increment to avoid infinite loop
			p.lineNum++
		}
	}
	
	return content
}

func (p *parser) parseParagraph() *asciidocxml.Paragraph {
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
	return &asciidocxml.Paragraph{
		Items: p.parseInlineContent(text).Items,
	}
}

func (p *parser) parseCodeBlock() *asciidocxml.CodeBlock {
	p.lineNum++ // Skip opening delimiter
	
	// Parse attributes from previous line if present
	var language, title string
	if p.lineNum > 1 {
		prevLine := strings.TrimSpace(p.lines[p.lineNum-2])
		if strings.HasPrefix(prevLine, "[") && strings.HasSuffix(prevLine, "]") {
			attrs := prevLine[1 : len(prevLine)-1]
			parts := strings.Split(attrs, ",")
			if len(parts) > 0 {
				language = strings.TrimSpace(parts[0])
			}
		}
		if strings.HasPrefix(prevLine, ".") {
			title = strings.TrimPrefix(prevLine, ".")
		}
	}
	
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
	
	return &asciidocxml.CodeBlock{
		Language: language,
		Title:    title,
		Content:  strings.Join(content, "\n"),
	}
}

func (p *parser) parseLiteralBlock() *asciidocxml.LiteralBlock {
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
	return &asciidocxml.LiteralBlock{
		Content: strings.Join(content, "\n"),
	}
}

func (p *parser) parseExampleBlock() *asciidocxml.Example {
	p.lineNum++ // Skip opening
	var title string
	if p.lineNum > 1 {
		prevLine := strings.TrimSpace(p.lines[p.lineNum-2])
		if strings.HasPrefix(prevLine, ".") {
			title = strings.TrimPrefix(prevLine, ".")
		}
	}
	
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
	
	// Simple content parsing
	content := asciidocxml.Content{}
	for _, cl := range contentLines {
		if strings.TrimSpace(cl) != "" {
			content.Items = append(content.Items, asciidocxml.ContentItem{
				Paragraph: &asciidocxml.Paragraph{
					Items: p.parseInlineContent(cl).Items,
				},
			})
		}
	}
	
	return &asciidocxml.Example{
		Title:   title,
		Content: content,
	}
}

func (p *parser) parseSidebar() *asciidocxml.Sidebar {
	p.lineNum++ // Skip opening
	var title string
	var contentLines []string
	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])
		if line == "****" {
			p.lineNum++
			break
		}
		if strings.HasPrefix(line, ".") {
			title = strings.TrimPrefix(line, ".")
		} else {
			contentLines = append(contentLines, p.lines[p.lineNum])
		}
		p.lineNum++
	}
	
	content := asciidocxml.Content{}
	for _, cl := range contentLines {
		if strings.TrimSpace(cl) != "" {
			content.Items = append(content.Items, asciidocxml.ContentItem{
				Paragraph: &asciidocxml.Paragraph{
					Items: p.parseInlineContent(cl).Items,
				},
			})
		}
	}
	
	return &asciidocxml.Sidebar{
		Title:   title,
		Content: content,
	}
}

func (p *parser) parseQuote() *asciidocxml.Quote {
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
	
	content := asciidocxml.Content{}
	for _, cl := range contentLines {
		if strings.TrimSpace(cl) != "" {
			content.Items = append(content.Items, asciidocxml.ContentItem{
				Paragraph: &asciidocxml.Paragraph{
					Items: p.parseInlineContent(cl).Items,
				},
			})
		}
	}
	
	return &asciidocxml.Quote{
		Content: content,
	}
}

func (p *parser) parseTable() *asciidocxml.Table {
	p.lineNum++ // Skip opening
	table := &asciidocxml.Table{
		Rows: []asciidocxml.TableRow{},
	}
	
	var headerRow *asciidocxml.TableRow
	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])
		if line == "|===" {
			p.lineNum++
			break
		}
		
		if strings.HasPrefix(line, "|") {
			cells := strings.Split(line, "|")
			row := asciidocxml.TableRow{
				Cells: []asciidocxml.TableCell{},
			}
			for i := 1; i < len(cells); i++ {
				cellText := strings.TrimSpace(cells[i])
				row.Cells = append(row.Cells, asciidocxml.TableCell{
					Items: p.parseInlineContent(cellText).Items,
				})
			}
			
			// First row is typically header
			if headerRow == nil {
				headerRow = &row
			} else {
				table.Rows = append(table.Rows, row)
			}
		}
		p.lineNum++
	}
	
	if headerRow != nil {
		table.Header = headerRow
	}
	
	return table
}

func (p *parser) parseList() *asciidocxml.List {
	var items []asciidocxml.ListItem
	style := "unordered"
	
	for p.lineNum < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.lineNum])
		
		if !p.isListItem(line) {
			break
		}
		
		// Determine style
		if strings.HasPrefix(line, ".") {
			style = "ordered"
		} else if strings.HasPrefix(line, "*") || strings.HasPrefix(line, "-") {
			style = "unordered"
		} else if strings.Contains(line, "::") {
			style = "labeled"
		}
		
		item := p.parseListItem(line)
		if item != nil {
			items = append(items, *item)
		}
		p.lineNum++
	}
	
	return &asciidocxml.List{
		Style: style,
		Items: items,
	}
}

func (p *parser) parseListItem(line string) *asciidocxml.ListItem {
	// Remove list marker
	content := strings.TrimLeft(line, ".*-")
	content = strings.TrimSpace(content)
	
	// Handle labeled lists
	if strings.Contains(content, "::") {
		parts := strings.SplitN(content, "::", 2)
		return &asciidocxml.ListItem{
			Term: p.parseInlineContent(strings.TrimSpace(parts[0])),
			Items: []asciidocxml.ListItemContentItem{
				{
					InlineItems: p.parseInlineContent(strings.TrimSpace(parts[1])).Items,
				},
			},
		}
	}
	
	return &asciidocxml.ListItem{
		Items: []asciidocxml.ListItemContentItem{
			{
				InlineItems: p.parseInlineContent(content).Items,
			},
		},
	}
}

func (p *parser) parseAdmonition() *asciidocxml.Admonition {
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
	
	return &asciidocxml.Admonition{
		Type: admonitionType,
		Content: asciidocxml.Content{
			Items: []asciidocxml.ContentItem{
				{
					Paragraph: &asciidocxml.Paragraph{
						Items: p.parseInlineContent(content).Items,
					},
				},
			},
		},
	}
}

func (p *parser) parseImage() *asciidocxml.Image {
	line := strings.TrimSpace(p.lines[p.lineNum])
	p.lineNum++
	
	// Parse image::path[alt,width,height] or image:path[alt]
	var src, alt string
	if strings.HasPrefix(line, "image::") {
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
	
	return &asciidocxml.Image{
		Src: src,
		Alt: alt,
	}
}

func (p *parser) parseInlineContent(text string) asciidocxml.InlineContent {
	items := []asciidocxml.InlineItem{}
	
	// Simple inline parsing - can be enhanced
	// Handle bold
	boldRegex := regexp.MustCompile(`\*([^*]+)\*`)
	text = boldRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := match[1 : len(match)-1]
		items = append(items, asciidocxml.InlineItem{
			Strong: &asciidocxml.InlineContent{
				Items: []asciidocxml.InlineItem{{Text: content}},
			},
		})
		return ""
	})
	
	// Handle italic
	italicRegex := regexp.MustCompile(`_([^_]+)_`)
	text = italicRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := match[1 : len(match)-1]
		items = append(items, asciidocxml.InlineItem{
			Emphasis: &asciidocxml.InlineContent{
				Items: []asciidocxml.InlineItem{{Text: content}},
			},
		})
		return ""
	})
	
	// Handle monospace
	monoRegex := regexp.MustCompile("`([^`]+)`")
	text = monoRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := match[1 : len(match)-1]
		items = append(items, asciidocxml.InlineItem{
			Monospace: &asciidocxml.InlineContent{
				Items: []asciidocxml.InlineItem{{Text: content}},
			},
		})
		return ""
	})
	
	// Handle links
	linkRegex := regexp.MustCompile(`(https?://[^\s\[\]]+)(\[([^\]]+)\])?`)
	text = linkRegex.ReplaceAllStringFunc(text, func(match string) string {
		parts := linkRegex.FindStringSubmatch(match)
		href := parts[1]
		text := href
		if len(parts) > 3 && parts[3] != "" {
			text = parts[3]
		}
		items = append(items, asciidocxml.InlineItem{
			Link: &asciidocxml.Link{
				Href: href,
				Items: []asciidocxml.InlineItem{{Text: text}},
			},
		})
		return ""
	})
	
	// Add remaining text
	if text != "" {
		items = append(items, asciidocxml.InlineItem{Text: text})
	}
	
	return asciidocxml.InlineContent{Items: items}
}

func (p *parser) isListItem(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "*") ||
		strings.HasPrefix(trimmed, ".") ||
		strings.HasPrefix(trimmed, "-") ||
		strings.Contains(trimmed, "::")
}

func (p *parser) isAdmonition(line string) bool {
	trimmed := strings.ToUpper(strings.TrimSpace(line))
	return strings.HasPrefix(trimmed, "NOTE:") ||
		strings.HasPrefix(trimmed, "TIP:") ||
		strings.HasPrefix(trimmed, "IMPORTANT:") ||
		strings.HasPrefix(trimmed, "WARNING:") ||
		strings.HasPrefix(trimmed, "CAUTION:")
}
