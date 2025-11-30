package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// LinkReference represents a link reference definition for CommonMark reference-style links
type LinkReference struct {
	URL   string
	Title string
}

// ConvertMarkdownToAsciiDocStreaming converts Markdown content to AsciiDoc format
// and writes directly to the provided writer, enabling streaming for large files.
// This enhanced version supports CommonMark specification and GitHub Flavored Markdown (GFM).
// It uses a two-pass approach: first pass collects link references, second pass converts content.
func ConvertMarkdownToAsciiDocStreaming(reader io.Reader, writer io.Writer) error {
	// Read all lines for two-pass processing (needed for link reference definitions)
	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading markdown: %w", err)
	}

	// First pass: collect link reference definitions (CommonMark)
	linkRefs := make(map[string]LinkReference)
	linkRefRegex := regexp.MustCompile(`^\[([^\]]+)\]:\s*(.+)$`)
	
	for _, line := range lines {
		if matches := linkRefRegex.FindStringSubmatch(line); matches != nil {
			refID := strings.ToLower(strings.TrimSpace(matches[1]))
			rest := strings.TrimSpace(matches[2])
			
			// Parse URL and optional title
			// Format: <url> or <url> "title" or <url> 'title' or <url> (title) or just url
			urlRegex := regexp.MustCompile(`^<?([^>\s]+)>?`)
			if urlMatch := urlRegex.FindStringSubmatch(rest); urlMatch != nil {
				url := urlMatch[1]
				title := ""
				
				// Check for title in quotes or parentheses
				titleRegex := regexp.MustCompile(`["']([^"']+)["']|\(([^)]+)\)`)
				if titleMatch := titleRegex.FindStringSubmatch(rest); titleMatch != nil {
					if titleMatch[1] != "" {
						title = titleMatch[1]
					} else if titleMatch[2] != "" {
						title = titleMatch[2]
					}
				}
				
				linkRefs[refID] = LinkReference{URL: url, Title: title}
			}
		}
	}

	// Second pass: convert with reference resolution
	return convertMarkdownLinesStreaming(lines, linkRefs, writer)
}

// convertMarkdownLinesStreaming performs the actual conversion with all enhanced features
func convertMarkdownLinesStreaming(lines []string, linkRefs map[string]LinkReference, writer io.Writer) error {
	var inCodeBlock bool
	var codeBlockLang string
	var inIndentedCodeBlock bool
	var codeBlockLines []string
	var inHTMLBlock bool
	var htmlBlockLines []string
	var htmlBlockTag string
	var inTable bool
	var tableAlignments []string
	var frontmatterProcessed bool
	var frontmatterLines []string
	var inFrontmatter bool
	var prevLine string
	var prevLineEmpty bool

	// Regex patterns
	headerRegex := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	setextHeaderRegex := regexp.MustCompile(`^(=+|-+)\s*$`) // Setext headers (CommonMark)
	codeBlockStartRegex := regexp.MustCompile("^```(\\w*)$")
	codeBlockEndRegex := regexp.MustCompile("^```\\s*$")
	indentedCodeBlockRegex := regexp.MustCompile(`^[ ]{4,}`) // 4+ spaces (CommonMark)
	htmlBlockStartRegex := regexp.MustCompile(`^<([a-zA-Z][a-zA-Z0-9]*)(?:\s[^>]*)?>`) // HTML block start (GFM)
	htmlBlockEndRegex := regexp.MustCompile(`^</([a-zA-Z][a-zA-Z0-9]*)>`) // HTML block end
	htmlSelfClosingRegex := regexp.MustCompile(`^<([a-zA-Z][a-zA-Z0-9]*)(?:\s[^>]*)?\s*/>`) // Self-closing tags
	htmlCommentRegex := regexp.MustCompile(`^<!--.*?-->`) // HTML comments
	tableRowRegex := regexp.MustCompile(`^\|(.+)\|$`)
	tableSeparatorRegex := regexp.MustCompile(`^\|?[\s\-:]+\|[\s\-:]+\|?$`)
	horizontalRuleRegex := regexp.MustCompile(`^[-*_]{3,}\s*$`)
	blockquoteRegex := regexp.MustCompile(`^>\s*(.*)$`)
	orderedListRegex := regexp.MustCompile(`^\s*(\d+)\.\s+(.+)$`)
	unorderedListRegex := regexp.MustCompile(`^\s*[-*+]\s+(.+)$`)
	taskListRegex := regexp.MustCompile(`^\s*[-*+]\s+\[([ xX])\]\s+(.+)$`) // GFM task lists
	frontmatterStartRegex := regexp.MustCompile(`^---\s*$`)
	frontmatterEndRegex := regexp.MustCompile(`^---\s*$|^\.\.\.\s*$`)
	linkRefDefRegex := regexp.MustCompile(`^\[([^\]]+)\]:`) // Skip link reference definitions

	lineNum := 0
	for i, line := range lines {
		lineNum++
		trimmedLine := strings.TrimSpace(line)
		prevLineEmpty = (i > 0 && strings.TrimSpace(lines[i-1]) == "")

		// Skip link reference definitions (already processed in first pass)
		if linkRefDefRegex.MatchString(trimmedLine) {
			continue
		}

		// Handle HTML blocks (GFM feature - out of spec for standard Markdown)
		if inHTMLBlock {
			// Check for closing tag matching the opening tag
			if htmlBlockEndRegex.MatchString(trimmedLine) {
				matches := htmlBlockEndRegex.FindStringSubmatch(trimmedLine)
				closingTag := strings.ToLower(matches[1])
				if closingTag == htmlBlockTag {
					// End of HTML block
					htmlBlockLines = append(htmlBlockLines, line)
					inHTMLBlock = false
					// Write HTML block as passthrough with comment
					if _, err := writer.Write([]byte("// GFM-specific HTML block (out-of-spec for standard Markdown)\n")); err != nil {
						return err
					}
					if _, err := writer.Write([]byte("++++\n")); err != nil {
						return err
					}
					for _, htmlLine := range htmlBlockLines {
						if _, err := fmt.Fprintf(writer, "%s\n", htmlLine); err != nil {
							return err
						}
					}
					if _, err := writer.Write([]byte("++++\n")); err != nil {
						return err
					}
					htmlBlockLines = nil
					htmlBlockTag = ""
					prevLine = ""
					continue
				}
			}
			// Continue collecting HTML block content
			htmlBlockLines = append(htmlBlockLines, line)
			prevLine = ""
			continue
		}

		// Check for HTML block start (must be before other block checks)
		// Only check if not in code blocks and line starts with HTML tag
		if !inCodeBlock && !inIndentedCodeBlock && strings.HasPrefix(trimmedLine, "<") {
			// Check for HTML comment (single line)
			if htmlCommentRegex.MatchString(trimmedLine) {
				if _, err := writer.Write([]byte("// GFM-specific HTML comment (out-of-spec for standard Markdown)\n")); err != nil {
					return err
				}
				if _, err := writer.Write([]byte("++++\n")); err != nil {
					return err
				}
				if _, err := fmt.Fprintf(writer, "%s\n", line); err != nil {
					return err
				}
				if _, err := writer.Write([]byte("++++\n")); err != nil {
					return err
				}
				prevLine = ""
				continue
			}
			// Check for self-closing HTML tag (single line block)
			if htmlSelfClosingRegex.MatchString(trimmedLine) {
				if _, err := writer.Write([]byte("// GFM-specific HTML block (out-of-spec for standard Markdown)\n")); err != nil {
					return err
				}
				if _, err := writer.Write([]byte("++++\n")); err != nil {
					return err
				}
				if _, err := fmt.Fprintf(writer, "%s\n", line); err != nil {
					return err
				}
				if _, err := writer.Write([]byte("++++\n")); err != nil {
					return err
				}
				prevLine = ""
				continue
			}
			// Check for opening HTML block tag
			if matches := htmlBlockStartRegex.FindStringSubmatch(trimmedLine); matches != nil {
				tagName := strings.ToLower(matches[1])
				// Common block-level HTML tags
				blockTags := map[string]bool{
					"div": true, "p": true, "table": true, "thead": true, "tbody": true, "tfoot": true,
					"tr": true, "th": true, "td": true, "ul": true, "ol": true, "li": true,
					"dl": true, "dt": true, "dd": true, "blockquote": true, "pre": true,
					"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
					"section": true, "article": true, "aside": true, "header": true, "footer": true,
					"nav": true, "main": true, "figure": true, "figcaption": true, "details": true,
					"summary": true, "form": true, "fieldset": true, "legend": true,
				}
				if blockTags[tagName] {
					inHTMLBlock = true
					htmlBlockTag = tagName
					htmlBlockLines = []string{line}
					prevLine = ""
					continue
				}
			}
		}

		// Handle YAML Frontmatter at the start of the file
		if !frontmatterProcessed && lineNum == 1 {
			if frontmatterStartRegex.MatchString(line) {
				inFrontmatter = true
				continue
			}
		}

		if inFrontmatter {
			if frontmatterEndRegex.MatchString(line) {
				frontmatterProcessed = true
				inFrontmatter = false
				if err := processFrontmatterStreaming(writer, frontmatterLines); err != nil {
					return err
				}
				frontmatterLines = nil
				continue
			}
			frontmatterLines = append(frontmatterLines, line)
			continue
		}

		// Handle indented code blocks (4+ spaces, CommonMark)
		if indentedCodeBlockRegex.MatchString(line) && !inCodeBlock {
			if !inIndentedCodeBlock {
				inIndentedCodeBlock = true
				codeBlockLines = []string{}
			}
			// Remove leading 4 spaces
			codeLine := line[4:]
			codeBlockLines = append(codeBlockLines, codeLine)
			prevLine = ""
			continue
		} else if inIndentedCodeBlock {
			// End indented code block if we hit a non-indented, non-empty line
			if trimmedLine != "" && !indentedCodeBlockRegex.MatchString(line) {
				inIndentedCodeBlock = false
				if _, err := writer.Write([]byte("----\n")); err != nil {
					return err
				}
				for _, codeLine := range codeBlockLines {
					if _, err := fmt.Fprintf(writer, "%s\n", codeLine); err != nil {
						return err
					}
				}
				if _, err := writer.Write([]byte("----\n")); err != nil {
					return err
				}
				codeBlockLines = nil
				// Continue processing this line
			} else if trimmedLine == "" {
				// Empty line in code block - check if next line continues it
				if i+1 < len(lines) {
					nextLine := lines[i+1]
					if indentedCodeBlockRegex.MatchString(nextLine) || strings.TrimSpace(nextLine) == "" {
						// Continue code block
						codeBlockLines = append(codeBlockLines, "")
						prevLine = ""
						continue
					} else {
						// End code block
						inIndentedCodeBlock = false
						if _, err := writer.Write([]byte("----\n")); err != nil {
							return err
						}
						for _, codeLine := range codeBlockLines {
							if _, err := fmt.Fprintf(writer, "%s\n", codeLine); err != nil {
								return err
							}
						}
						if _, err := writer.Write([]byte("----\n")); err != nil {
							return err
						}
						codeBlockLines = nil
						// Continue processing this line
					}
				} else {
					// End of file, close code block
					inIndentedCodeBlock = false
					if _, err := writer.Write([]byte("----\n")); err != nil {
						return err
					}
					for _, codeLine := range codeBlockLines {
						if _, err := fmt.Fprintf(writer, "%s\n", codeLine); err != nil {
							return err
						}
					}
					if _, err := writer.Write([]byte("----\n")); err != nil {
						return err
					}
					codeBlockLines = nil
				}
			} else {
				prevLine = ""
				continue
			}
		}

		// Handle fenced code blocks
		if inCodeBlock {
			if codeBlockEndRegex.MatchString(line) {
				if codeBlockLang != "" {
					if _, err := fmt.Fprintf(writer, "[source,%s]\n", codeBlockLang); err != nil {
						return err
					}
				}
				if _, err := writer.Write([]byte("----\n")); err != nil {
					return err
				}
				for _, codeLine := range codeBlockLines {
					if _, err := fmt.Fprintf(writer, "%s\n", codeLine); err != nil {
						return err
					}
				}
				if _, err := writer.Write([]byte("----\n")); err != nil {
					return err
				}
				inCodeBlock = false
				codeBlockLang = ""
				codeBlockLines = nil
				prevLine = ""
				continue
			}
			codeBlockLines = append(codeBlockLines, line)
			prevLine = ""
			continue
		}

		// Check for fenced code block start
		if codeBlockStartRegex.MatchString(line) {
			matches := codeBlockStartRegex.FindStringSubmatch(line)
			codeBlockLang = matches[1]
			inCodeBlock = true
			codeBlockLines = []string{}
			prevLine = ""
			continue
		}

		// Handle Setext-style headers (=== and ---, CommonMark)
		// Must check after code blocks but before regular headers
		if setextHeaderRegex.MatchString(trimmedLine) && prevLine != "" && !prevLineEmpty && trimmedLine != "" {
			// Determine level: === is level 1, --- is level 2
			level := 1
			if strings.HasPrefix(trimmedLine, "---") {
				level = 2
			}
			// Use previous line as title
			title := strings.TrimSpace(prevLine)
			equals := strings.Repeat("=", level)
			if _, err := fmt.Fprintf(writer, "%s %s\n", equals, title); err != nil {
				return err
			}
			prevLine = ""
			continue
		}

		// Handle ATX-style headers
		if matches := headerRegex.FindStringSubmatch(line); matches != nil {
			level := len(matches[1])
			title := strings.TrimSpace(matches[2])
			equals := strings.Repeat("=", level)
			if _, err := fmt.Fprintf(writer, "%s %s\n", equals, title); err != nil {
				return err
			}
			prevLine = ""
			continue
		}

		// Handle horizontal rules
		if horizontalRuleRegex.MatchString(trimmedLine) {
			if _, err := writer.Write([]byte("'''\n")); err != nil {
				return err
			}
			prevLine = ""
			continue
		}

		// Handle blockquotes
		if matches := blockquoteRegex.FindStringSubmatch(line); matches != nil {
			content := matches[1]
			admonitionType := ""
			boldAdmonitionRegex := regexp.MustCompile(`^\*\*(\w+)\*\*\s*`)
			if admMatches := boldAdmonitionRegex.FindStringSubmatch(content); admMatches != nil {
				admonitionType = strings.ToUpper(admMatches[1])
				content = strings.TrimSpace(boldAdmonitionRegex.ReplaceAllString(content, ""))
			} else {
				bracketAdmonitionRegex := regexp.MustCompile(`^\[!(\w+)\]\s*`)
				if admMatches := bracketAdmonitionRegex.FindStringSubmatch(content); admMatches != nil {
					admonitionType = strings.ToUpper(admMatches[1])
					content = strings.TrimSpace(bracketAdmonitionRegex.ReplaceAllString(content, ""))
				}
			}

			if admonitionType != "" {
				if _, err := fmt.Fprintf(writer, "%s: %s\n", admonitionType, content); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(writer, "[quote]\n____\n%s\n____\n", content); err != nil {
					return err
				}
			}
			prevLine = ""
			continue
		}

		// Handle tables with alignment detection (GFM)
		if tableRowRegex.MatchString(line) {
			if tableSeparatorRegex.MatchString(line) {
				// This is a separator row - extract alignments
				tableAlignments = extractTableAlignments(line)
				continue
			}
			if !inTable {
				cells := strings.Split(line, "|")
				if len(cells) > 0 && strings.TrimSpace(cells[0]) == "" {
					cells = cells[1:]
				}
				if len(cells) > 0 && strings.TrimSpace(cells[len(cells)-1]) == "" {
					cells = cells[:len(cells)-1]
				}
				colCount := len(cells)
				// Build column specification with alignments
				colSpec := buildColumnSpec(colCount, tableAlignments)
				if _, err := fmt.Fprintf(writer, "[cols=\"%s\"]\n", colSpec); err != nil {
					return err
				}
				if _, err := writer.Write([]byte("|===\n")); err != nil {
					return err
				}
				inTable = true
			}
			cells := strings.Split(line, "|")
			if len(cells) > 0 && strings.TrimSpace(cells[0]) == "" {
				cells = cells[1:]
			}
			if len(cells) > 0 && strings.TrimSpace(cells[len(cells)-1]) == "" {
				cells = cells[:len(cells)-1]
			}
			processedCells := make([]string, len(cells))
			for i, cell := range cells {
				cellContent := strings.TrimSpace(cell)
				processedCells[i] = convertInlineMarkdownEnhanced(cellContent, false, linkRefs)
			}
			if _, err := fmt.Fprintf(writer, "|%s\n", strings.Join(processedCells, " |")); err != nil {
				return err
			}
			prevLine = ""
			continue
		} else if inTable {
			if _, err := writer.Write([]byte("|===\n")); err != nil {
				return err
			}
			inTable = false
			tableAlignments = nil
		}

		// Handle task lists (GFM) - must check before regular unordered lists
		if matches := taskListRegex.FindStringSubmatch(line); matches != nil {
			checked := strings.ToLower(matches[1]) == "x"
			content := matches[2]
			leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
			content = convertInlineMarkdownEnhanced(content, false, linkRefs)
			indent := strings.Repeat(" ", leadingSpaces)
			checkbox := "[ ]"
			if checked {
				checkbox = "[x]"
			}
			if _, err := fmt.Fprintf(writer, "%s* %s %s\n", indent, checkbox, content); err != nil {
				return err
			}
			prevLine = ""
			continue
		}

		// Handle ordered lists
		if matches := orderedListRegex.FindStringSubmatch(line); matches != nil {
			content := matches[2]
			leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
			content = convertInlineMarkdownEnhanced(content, false, linkRefs)
			indent := strings.Repeat(" ", leadingSpaces)
			if _, err := fmt.Fprintf(writer, "%s. %s\n", indent, content); err != nil {
				return err
			}
			prevLine = ""
			continue
		}

		// Handle unordered lists
		if matches := unorderedListRegex.FindStringSubmatch(line); matches != nil {
			content := matches[1]
			leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
			content = convertInlineMarkdownEnhanced(content, false, linkRefs)
			indent := strings.Repeat(" ", leadingSpaces)
			if _, err := fmt.Fprintf(writer, "%s* %s\n", indent, content); err != nil {
				return err
			}
			prevLine = ""
			continue
		}

		// Handle empty lines
		if trimmedLine == "" {
			if _, err := writer.Write([]byte("\n")); err != nil {
				return err
			}
			prevLine = ""
			continue
		}

		// Check for standalone image
		imageRegex := regexp.MustCompile(`^!\[([^\]]*)\]\(([^)]+)\)\s*$`)
		if imageRegex.MatchString(trimmedLine) {
			matches := imageRegex.FindStringSubmatch(trimmedLine)
			alt := matches[1]
			src := matches[2]
			if _, err := fmt.Fprintf(writer, "image::%s[%s]\n", src, alt); err != nil {
				return err
			}
			prevLine = ""
			continue
		}

		// Check for reference-style image ![alt][ref]
		refImageRegex := regexp.MustCompile(`^!\[([^\]]*)\]\[([^\]]*)\]\s*$`)
		if matches := refImageRegex.FindStringSubmatch(trimmedLine); matches != nil {
			alt := matches[1]
			refID := strings.ToLower(strings.TrimSpace(matches[2]))
			if refID == "" {
				// Implicit reference: use alt text as ref ID
				refID = strings.ToLower(alt)
			}
			if ref, ok := linkRefs[refID]; ok {
				if _, err := fmt.Fprintf(writer, "image::%s[%s]\n", ref.URL, alt); err != nil {
					return err
				}
			}
			prevLine = ""
			continue
		}

		// Convert inline Markdown with enhanced features (includes HTML span handling)
		convertedLine := convertInlineMarkdownEnhanced(line, false, linkRefs)
		
		// Handle hard line breaks (two spaces at end of line, CommonMark)
		if strings.HasSuffix(line, "  ") && !strings.HasSuffix(convertedLine, "  ") {
			// Remove trailing spaces from converted line and add line break
			convertedLine = strings.TrimRight(convertedLine, " ")
			convertedLine += " +\n"
		} else {
			convertedLine += "\n"
		}
		
		if _, err := writer.Write([]byte(convertedLine)); err != nil {
			return err
		}
		prevLine = line
	}

	// Close table if still open
	if inTable {
		if _, err := writer.Write([]byte("|===\n")); err != nil {
			return err
		}
	}

	// Close indented code block if still open
	if inIndentedCodeBlock && len(codeBlockLines) > 0 {
		if _, err := writer.Write([]byte("----\n")); err != nil {
			return err
		}
		for _, codeLine := range codeBlockLines {
			if _, err := fmt.Fprintf(writer, "%s\n", codeLine); err != nil {
				return err
			}
		}
		if _, err := writer.Write([]byte("----\n")); err != nil {
			return err
		}
	}

	// Close HTML block if still open
	if inHTMLBlock && len(htmlBlockLines) > 0 {
		if _, err := writer.Write([]byte("// GFM-specific HTML block (out-of-spec for standard Markdown)\n")); err != nil {
			return err
		}
		if _, err := writer.Write([]byte("++++\n")); err != nil {
			return err
		}
		for _, htmlLine := range htmlBlockLines {
			if _, err := fmt.Fprintf(writer, "%s\n", htmlLine); err != nil {
				return err
			}
		}
		if _, err := writer.Write([]byte("++++\n")); err != nil {
			return err
		}
	}

	return nil
}

// extractTableAlignments extracts column alignments from table separator row (GFM)
func extractTableAlignments(separatorRow string) []string {
	alignments := []string{}
	cells := strings.Split(separatorRow, "|")
	
	for _, cell := range cells {
		cell = strings.TrimSpace(cell)
		if cell == "" {
			continue
		}
		
		// Check alignment: :--- = left, ---: = right, :---: = center, --- = default
		leftAlign := strings.HasPrefix(cell, ":")
		rightAlign := strings.HasSuffix(cell, ":")
		
		if leftAlign && rightAlign {
			alignments = append(alignments, "^") // Center
		} else if rightAlign {
			alignments = append(alignments, ">") // Right
		} else if leftAlign {
			alignments = append(alignments, "<") // Left
		} else {
			alignments = append(alignments, "") // Default
		}
	}
	
	return alignments
}

// buildColumnSpec builds AsciiDoc column specification with alignments
func buildColumnSpec(colCount int, alignments []string) string {
	if len(alignments) == 0 {
		// No alignments specified, use default
		colSpec := strings.Repeat("1,", colCount)
		if len(colSpec) > 0 {
			colSpec = colSpec[:len(colSpec)-1]
		}
		return colSpec
	}
	
	// Build spec with alignments: <1,^1,>1,1
	var parts []string
	for i := 0; i < colCount; i++ {
		align := ""
		if i < len(alignments) {
			align = alignments[i]
		}
		parts = append(parts, align+"1")
	}
	return strings.Join(parts, ",")
}

// convertInlineMarkdownEnhanced converts inline Markdown with enhanced features:
// - Reference-style links and images
// - Autolinks
// - Strikethrough (GFM)
// - HTML spans (GFM - out of spec for standard Markdown)
// - Better escaped character handling
func convertInlineMarkdownEnhanced(text string, isStandaloneLine bool, linkRefs map[string]LinkReference) string {
	// Handle escaped characters first (backslash escaping, CommonMark)
	text = handleEscapedCharacters(text)
	
	// Convert HTML spans to AsciiDoc inline passthrough (GFM feature)
	// Common inline HTML tags: strong, em, code, a, img, br, span, mark, del, ins, sub, sup, etc.
	// Note: Go regexp doesn't support backreferences, so we use a simpler approach
	htmlSelfClosingSpanRegex := regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*)(?:\s[^>]*)?\s*/>`)
	
	// Common inline HTML tags
	inlineTags := []string{"strong", "b", "em", "i", "code", "a", "img", "br", "span", "mark",
		"del", "ins", "sub", "sup", "small", "kbd", "samp", "var", "time", "abbr",
		"cite", "q", "dfn", "u", "bdi", "bdo"}
	
	// Process self-closing inline tags first (like <br/>, <img/>, etc.)
	text = htmlSelfClosingSpanRegex.ReplaceAllStringFunc(text, func(match string) string {
		// Convert to AsciiDoc inline passthrough
		// Note: GFM-specific, out of spec for standard Markdown
		return fmt.Sprintf("pass:[%s]", match)
	})
	
	// Process paired HTML tags (like <strong>text</strong>)
	// We need to match each tag type individually since Go doesn't support backreferences
	for _, tagName := range inlineTags {
		// Match opening tag, content, and closing tag for this specific tag
		pattern := fmt.Sprintf(`<(%s)(?:\s[^>]*)?>(.*?)</%s>`, tagName, tagName)
		htmlTagRegex := regexp.MustCompile(pattern)
		text = htmlTagRegex.ReplaceAllStringFunc(text, func(match string) string {
			// Convert to AsciiDoc inline passthrough
			// Note: GFM-specific, out of spec for standard Markdown
			return fmt.Sprintf("pass:[%s]", match)
		})
	}
	
	// Convert images - check for reference style first
	refImageRegex := regexp.MustCompile(`!\[([^\]]*)\]\[([^\]]*)\]`)
	text = refImageRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := refImageRegex.FindStringSubmatch(match)
		alt := matches[1]
		refID := strings.ToLower(strings.TrimSpace(matches[2]))
		if refID == "" {
			// Implicit reference: use alt text as ref ID
			refID = strings.ToLower(alt)
		}
		if ref, ok := linkRefs[refID]; ok {
			if isStandaloneLine {
				return fmt.Sprintf("image::%s[%s]", ref.URL, alt)
			}
			return fmt.Sprintf("image:%s[%s]", ref.URL, alt)
		}
		return match // Return original if reference not found
	})
	
	// Convert inline images
	imageRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	trimmedText := strings.TrimSpace(text)
	imageOnly := imageRegex.MatchString(trimmedText) && strings.TrimSpace(imageRegex.ReplaceAllString(trimmedText, "")) == ""
	
	text = imageRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := imageRegex.FindStringSubmatch(match)
		alt := matches[1]
		src := matches[2]
		// Check for title in src (src "title" or src 'title')
		titleRegex := regexp.MustCompile(`^(.+?)\s+["']([^"']+)["']$`)
		if titleMatch := titleRegex.FindStringSubmatch(src); titleMatch != nil {
			actualSrc := titleMatch[1]
			title := titleMatch[2]
			if imageOnly || isStandaloneLine {
				return fmt.Sprintf("image::%s[%s, title=\"%s\"]", actualSrc, alt, title)
			}
			return fmt.Sprintf("image:%s[%s, title=\"%s\"]", actualSrc, alt, title)
		}
		if imageOnly || isStandaloneLine {
			return fmt.Sprintf("image::%s[%s]", src, alt)
		}
		return fmt.Sprintf("image:%s[%s]", src, alt)
	})

	// Convert autolinks (<url> and <email@example.com>, CommonMark/GFM)
	autolinkRegex := regexp.MustCompile(`<([^>]+)>`)
	text = autolinkRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := autolinkRegex.FindStringSubmatch(match)
		url := matches[1]
		// Check if it's an email (contains @ but not ://)
		if strings.Contains(url, "@") && !strings.Contains(url, "://") {
			return fmt.Sprintf("link:mailto:%s[%s]", url, url)
		}
		// Regular URL - check if it already has a scheme
		if !strings.Contains(url, "://") {
			url = "http://" + url
		}
		return fmt.Sprintf("link:%s[%s]", url, url)
	})

	// Convert reference-style links [text][ref] (CommonMark)
	refLinkRegex := regexp.MustCompile(`\[([^\]]+)\]\[([^\]]*)\]`)
	text = refLinkRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := refLinkRegex.FindStringSubmatch(match)
		linkText := matches[1]
		refID := strings.ToLower(strings.TrimSpace(matches[2]))
		if refID == "" {
			// Implicit reference: use link text as ref ID
			refID = strings.ToLower(linkText)
		}
		if ref, ok := linkRefs[refID]; ok {
			if ref.Title != "" {
				return fmt.Sprintf("link:%s[%s, title=\"%s\"]", ref.URL, linkText, ref.Title)
			}
			return fmt.Sprintf("link:%s[%s]", ref.URL, linkText)
		}
		return match // Return original if reference not found
	})

	// Convert inline links [text](url)
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = linkRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := linkRegex.FindStringSubmatch(match)
		linkText := matches[1]
		url := matches[2]
		// Check for title in URL (url "title" or url 'title')
		titleRegex := regexp.MustCompile(`^(.+?)\s+["']([^"']+)["']$`)
		if titleMatch := titleRegex.FindStringSubmatch(url); titleMatch != nil {
			actualURL := titleMatch[1]
			title := titleMatch[2]
			return fmt.Sprintf("link:%s[%s, title=\"%s\"]", actualURL, linkText, title)
		}
		return fmt.Sprintf("link:%s[%s]", url, linkText)
	})

	// Convert strikethrough ~~text~~ (GFM)
	strikethroughRegex := regexp.MustCompile(`~~([^~]+)~~`)
	text = strikethroughRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := strikethroughRegex.FindStringSubmatch(match)
		content := matches[1]
		// AsciiDoc uses [.line-through]#text# for strikethrough
		return fmt.Sprintf("[.line-through]#%s#", content)
	})

	// Convert bold **text** and __text__ (using placeholder approach)
	boldDoubleStarRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	placeholderMap := make(map[string]string)
	placeholderCounter := 0
	
	text = boldDoubleStarRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholderCounter++
		placeholder := fmt.Sprintf("__BOLD_PLACEHOLDER_%d__", placeholderCounter)
		matches := boldDoubleStarRegex.FindStringSubmatch(match)
		placeholderMap[placeholder] = "**" + matches[1] + "**"
		return placeholder
	})
	
	boldDoubleUnderscoreRegex := regexp.MustCompile(`__([^_]+)__`)
	text = boldDoubleUnderscoreRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholderCounter++
		placeholder := fmt.Sprintf("__BOLD_PLACEHOLDER_%d__", placeholderCounter)
		matches := boldDoubleUnderscoreRegex.FindStringSubmatch(match)
		placeholderMap[placeholder] = "**" + matches[1] + "**"
		return placeholder
	})

	// Convert italic *text* (now safe since bold is in placeholders)
	// Important: Don't match *text* if it's part of **text** (double asterisks for bold)
	// Match *text* but ensure it's not preceded or followed by another *
	italicAsteriskRegex := regexp.MustCompile(`([^*]|^)\*([^*\n]+?)\*([^*]|$)`)
	text = italicAsteriskRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := italicAsteriskRegex.FindStringSubmatch(match)
		if len(matches) >= 4 {
			before := matches[1]
			content := matches[2]
			after := matches[3]
			if strings.TrimSpace(content) != "" {
				return before + "_" + content + "_" + after
			}
		}
		return match
	})

	// Convert italic _text_ (with word boundary awareness for better CommonMark compliance)
	// Only convert if it's not part of a word (word boundary)
	italicUnderscoreRegex := regexp.MustCompile(`(^|[^_a-zA-Z0-9])_([^_\n]+?)_([^_a-zA-Z0-9]|$)`)
	text = italicUnderscoreRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := italicUnderscoreRegex.FindStringSubmatch(match)
		if len(matches) >= 4 {
			before := matches[1]
			content := matches[2]
			after := matches[3]
			if strings.TrimSpace(content) != "" {
				return before + "_" + content + "_" + after
			}
		}
		return match
	})

	// Restore bold placeholders
	for placeholder, replacement := range placeholderMap {
		text = strings.ReplaceAll(text, placeholder, replacement)
	}

	// Inline code `code` is already in correct format for AsciiDoc
	// No conversion needed

	return text
}

// handleEscapedCharacters processes backslash-escaped characters (CommonMark)
func handleEscapedCharacters(text string) string {
	// CommonMark escaping: backslash escapes special characters
	// We need to remove the backslash and keep the character
	// Handle escaped backslash first: \\ -> \
	result := strings.ReplaceAll(text, `\\`, `\ESCAPED_BACKSLASH\`)
	
	// Replace other escaped characters
	escapedChars := map[string]string{
		`\*`: `*`,
		`\_`: `_`,
		`\#`: `#`,
		`\+`: `+`,
		`\-`: `-`,
		`\.`: `.`,
		`\!`: `!`,
		`\[`: `[`,
		`\]`: `]`,
		`\(`: `(`,
		`\)`: `)`,
		`\{`: `{`,
		`\}`: `}`,
		`\|`: `|`,
		`\~`: `~`,
		`\>`: `>`,
		`\<`: `<`,
		`\\`: `\`,
	}
	
	for esc, char := range escapedChars {
		result = strings.ReplaceAll(result, esc, char)
	}
	
	// Restore escaped backslashes
	result = strings.ReplaceAll(result, `\ESCAPED_BACKSLASH\`, `\`)
	
	return result
}

// processFrontmatterStreaming converts YAML frontmatter to AsciiDoc header attributes
// and writes directly to the writer. This is a copy of processFrontmatter logic adapted for streaming.
func processFrontmatterStreaming(writer io.Writer, lines []string) error {
	// Use the same logic as processFrontmatter but write to writer instead of buffer
	var buf bytes.Buffer
	processFrontmatter(&buf, lines)
	_, err := buf.WriteTo(writer)
	return err
}
