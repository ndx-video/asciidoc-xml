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

// Convert converts AsciiDoc content to our generic DOM Node
func Convert(reader io.Reader) (*Node, error) {
	return Parse(reader)
}

// ConvertToXML converts AsciiDoc to XML string using the DOM
func ConvertToXML(reader io.Reader) (string, error) {
	doc, err := Convert(reader)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	xmlContent, err := doc.ToXML()
	if err != nil {
		return "", err
	}
	buf.WriteString(xmlContent)

	return buf.String(), nil
}

// ConvertToHTML converts AsciiDoc to HTML5 string
// If xhtml is true, outputs well-formed XHTML5
// If usePicoCSS is true, includes PicoCSS styling
// picoCSSPath is used for <link> tag (if empty and usePicoCSS is true, picoCSSContent is embedded inline)
// picoCSSContent is embedded as <style> tag when usePicoCSS is true and picoCSSPath is empty
func ConvertToHTML(reader io.Reader, xhtml bool, usePicoCSS bool, picoCSSPath string, picoCSSContent string) (string, error) {
	doc, err := Convert(reader)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if xhtml {
		buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
		buf.WriteString(`<!DOCTYPE html>` + "\n")
	} else {
		buf.WriteString(`<!DOCTYPE html>` + "\n")
	}

	// Get lang from header attributes
	lang := "en"
	var headerNode *Node
	for _, child := range doc.Children {
		if child.Data == "header" {
			headerNode = child
			break
		}
	}
	if headerNode != nil {
		for _, attr := range headerNode.Children {
			if attr.Data == "attribute" && attr.GetAttribute("name") == "lang" {
				lang = attr.GetAttribute("value")
				break
			}
		}
	}

	if xhtml {
		fmt.Fprintf(&buf, `<html xmlns="http://www.w3.org/1999/xhtml" lang="%s">`+"\n", html.EscapeString(lang))
	} else {
		fmt.Fprintf(&buf, `<html lang="%s">`+"\n", html.EscapeString(lang))
	}

	// Head section
	buf.WriteString("  <head>\n")
	if xhtml {
		buf.WriteString("    <meta charset=\"UTF-8\"/>\n")
	} else {
		buf.WriteString("    <meta charset=\"UTF-8\">\n")
	}
	
	// Add PicoCSS if enabled
	if usePicoCSS {
		if picoCSSContent != "" {
			// Embed CSS inline
			buf.WriteString("    <style>\n")
			buf.WriteString(picoCSSContent)
			buf.WriteString("\n    </style>\n")
		} else if picoCSSPath != "" {
			// Use link tag
			if xhtml {
				fmt.Fprintf(&buf, `    <link rel="stylesheet" href="%s"/>`+"\n", html.EscapeString(picoCSSPath))
			} else {
				fmt.Fprintf(&buf, `    <link rel="stylesheet" href="%s">`+"\n", html.EscapeString(picoCSSPath))
			}
		}
	}
	
	if headerNode != nil {
		titleNode := findChild(headerNode, "title")
		if titleNode != nil {
			titleText := getTextContent(titleNode)
			fmt.Fprintf(&buf, "    <title>%s</title>\n", html.EscapeString(titleText))
		}
	}
	buf.WriteString("  </head>\n")

	// Body section
	buf.WriteString("  <body>\n")

	// Document header (title, authors, etc.)
	if headerNode != nil {
		writeHTMLHeader(&buf, headerNode, xhtml, 2)
	}

	// Main content
	buf.WriteString("    <main>\n")
	for _, child := range doc.Children {
		if child.Data != "header" {
			writeHTMLNode(&buf, child, xhtml, 2)
		}
	}
	buf.WriteString("    </main>\n")

	buf.WriteString("  </body>\n")
	buf.WriteString("</html>\n")

	return buf.String(), nil
}

func findChild(parent *Node, name string) *Node {
	for _, child := range parent.Children {
		if child.Data == name {
			return child
		}
	}
	return nil
}

func getTextContent(node *Node) string {
	var buf bytes.Buffer
	for _, child := range node.Children {
		if child.Type == TextNode {
			buf.WriteString(child.Data)
		} else {
			buf.WriteString(getTextContent(child))
		}
	}
	return buf.String()
}

func writeHTMLHeader(buf *bytes.Buffer, header *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)
	
	titleNode := findChild(header, "title")
	if titleNode != nil {
		buf.WriteString(indentStr + "<header>\n")
		buf.WriteString(indentStr + "      <h1>")
		writeHTMLInlineContent(buf, titleNode, xhtml)
		buf.WriteString("</h1>\n")

		// Authors
		for _, child := range header.Children {
			if child.Data == "author" {
				buf.WriteString(indentStr + "      <address class=\"authors\">\n")
				nameNode := findChild(child, "name")
				if nameNode != nil {
					buf.WriteString(indentStr + "        <p>")
					fmt.Fprintf(buf, `<span class="author-name">%s</span>`, html.EscapeString(getTextContent(nameNode)))
					emailNode := findChild(child, "email")
					if emailNode != nil {
						email := getTextContent(emailNode)
						fmt.Fprintf(buf, ` <a href="mailto:%s" class="author-email">%s</a>`, html.EscapeString(email), html.EscapeString(email))
					}
					buf.WriteString("</p>\n")
				}
				buf.WriteString(indentStr + "      </address>\n")
			}
		}

		buf.WriteString(indentStr + "    </header>\n")
	}
}

func writeHTMLNode(buf *bytes.Buffer, node *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)

	switch node.Data {
	case "section":
		writeHTMLSection(buf, node, xhtml, indent)
	case "paragraph":
		writeHTMLParagraph(buf, node, xhtml, indent)
	case "codeblock":
		writeHTMLCodeBlock(buf, node, xhtml, indent)
	case "literalblock":
		writeHTMLLiteralBlock(buf, node, xhtml, indent)
	case "list":
		writeHTMLList(buf, node, xhtml, indent)
	case "table":
		writeHTMLTable(buf, node, xhtml, indent)
	case "image":
		writeHTMLImage(buf, node, xhtml, indent)
	case "admonition":
		writeHTMLAdmonition(buf, node, xhtml, indent)
	case "example":
		writeHTMLExample(buf, node, xhtml, indent)
	case "sidebar":
		writeHTMLSidebar(buf, node, xhtml, indent)
	case "quote":
		writeHTMLQuote(buf, node, xhtml, indent)
	case "thematicbreak":
		if xhtml {
			buf.WriteString(indentStr + "<hr/>\n")
		} else {
			buf.WriteString(indentStr + "<hr>\n")
		}
	case "pagebreak":
		buf.WriteString(indentStr + `<div class="page-break"></div>` + "\n")
	default:
		// For unknown elements, just write children
		for _, child := range node.Children {
			writeHTMLNode(buf, child, xhtml, indent)
		}
	}
}

func writeHTMLSection(buf *bytes.Buffer, section *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)
	level := 1
	if levelAttr := section.GetAttribute("level"); levelAttr != "" {
		fmt.Sscanf(levelAttr, "%d", &level)
		level++ // HTML h1-h6
	}
	if level > 6 {
		level = 6
	}

	tagName := fmt.Sprintf("h%d", level)
	attrs := ""
	if id := section.GetAttribute("id"); id != "" {
		attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
	}
	if role := section.GetAttribute("role"); role != "" {
		attrs += fmt.Sprintf(` class="%s"`, html.EscapeString(role))
	}

	titleNode := findChild(section, "title")
	if titleNode != nil {
		fmt.Fprintf(buf, "%s<%s%s>", indentStr, tagName, attrs)
		writeHTMLInlineContent(buf, titleNode, xhtml)
		fmt.Fprintf(buf, "</%s>\n", tagName)
	}

	// Write section content
	for _, child := range section.Children {
		if child.Data != "title" {
			writeHTMLNode(buf, child, xhtml, indent)
		}
	}
}

func writeHTMLParagraph(buf *bytes.Buffer, para *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)
	attrs := ""
	if id := para.GetAttribute("id"); id != "" {
		attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
	}
	if role := para.GetAttribute("role"); role != "" {
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
	writeHTMLInlineContent(buf, para, xhtml)
	buf.WriteString("</p>\n")
}

func writeHTMLCodeBlock(buf *bytes.Buffer, code *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)
	attrs := ""
	if id := code.GetAttribute("id"); id != "" {
		attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
	}
	classParts := []string{"code-block"}
	if lang := code.GetAttribute("language"); lang != "" {
		classParts = append(classParts, "language-"+lang)
		attrs += fmt.Sprintf(` data-language="%s"`, html.EscapeString(lang))
	}
	if role := code.GetAttribute("role"); role != "" {
		classParts = append(classParts, role)
	}
	if len(classParts) > 0 {
		if attrs != "" {
			attrs += " "
		}
		attrs += fmt.Sprintf(`class="%s"`, html.EscapeString(strings.Join(classParts, " ")))
	}

	if title := code.GetAttribute("title"); title != "" {
		fmt.Fprintf(buf, "%s<p class=\"code-title\">%s</p>\n", indentStr, html.EscapeString(title))
	}

	if attrs != "" {
		fmt.Fprintf(buf, "%s<pre><code%s>", indentStr, attrs)
	} else {
		fmt.Fprintf(buf, "%s<pre><code>", indentStr)
	}
	
	// Write code content
	for _, child := range code.Children {
		if child.Type == TextNode {
			buf.WriteString(html.EscapeString(child.Data))
		}
	}
	
	buf.WriteString("</code></pre>\n")
}

func writeHTMLLiteralBlock(buf *bytes.Buffer, literal *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)
	attrs := `class="literal-block"`
	if id := literal.GetAttribute("id"); id != "" {
		attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
	}

	fmt.Fprintf(buf, "%s<pre%s>", indentStr, attrs)
	for _, child := range literal.Children {
		if child.Type == TextNode {
			buf.WriteString(html.EscapeString(child.Data))
		}
	}
	buf.WriteString("</pre>\n")
}

func writeHTMLList(buf *bytes.Buffer, list *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)
	style := list.GetAttribute("style")
	tagName := "ul"
	if style == "ordered" {
		tagName = "ol"
	} else if style == "labeled" {
		tagName = "dl"
	}

	attrs := ""
	if id := list.GetAttribute("id"); id != "" {
		attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
	}
	if role := list.GetAttribute("role"); role != "" {
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

	for _, item := range list.Children {
		if item.Data == "item" {
			if style == "labeled" {
				term := findChild(item, "term")
				if term != nil {
					fmt.Fprintf(buf, "%s    <dt>", indentStr)
					writeHTMLInlineContent(buf, term, xhtml)
					buf.WriteString("</dt>\n")
				}
				desc := findChild(item, "description")
				if desc != nil {
					fmt.Fprintf(buf, "%s    <dd>", indentStr)
					writeHTMLInlineContent(buf, desc, xhtml)
					buf.WriteString("</dd>\n")
				}
			} else {
				fmt.Fprintf(buf, "%s    <li>", indentStr)
				writeHTMLInlineContent(buf, item, xhtml)
				buf.WriteString("</li>\n")
			}
		}
	}

	fmt.Fprintf(buf, "%s</%s>\n", indentStr, tagName)
}

func writeHTMLTable(buf *bytes.Buffer, table *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)
	attrs := `class="table"`
	if id := table.GetAttribute("id"); id != "" {
		attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
	}

	fmt.Fprintf(buf, "%s<table%s>\n", indentStr, attrs)

	// First row is typically header
	firstRow := true
	for _, child := range table.Children {
		if child.Data == "row" {
			if firstRow {
				fmt.Fprintf(buf, "%s    <thead>\n", indentStr)
				fmt.Fprintf(buf, "%s        <tr>\n", indentStr)
				for _, cell := range child.Children {
					if cell.Data == "cell" {
						fmt.Fprintf(buf, "%s            <th>", indentStr)
						writeHTMLInlineContent(buf, cell, xhtml)
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
					if cell.Data == "cell" {
						fmt.Fprintf(buf, "%s            <td>", indentStr)
						writeHTMLInlineContent(buf, cell, xhtml)
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
}

func writeHTMLImage(buf *bytes.Buffer, img *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)
	attrs := fmt.Sprintf(`src="%s"`, html.EscapeString(img.GetAttribute("src")))
	if alt := img.GetAttribute("alt"); alt != "" {
		attrs += fmt.Sprintf(` alt="%s"`, html.EscapeString(alt))
	} else {
		attrs += ` alt=""`
	}

	if xhtml {
		fmt.Fprintf(buf, "%s<img %s/>\n", indentStr, attrs)
	} else {
		fmt.Fprintf(buf, "%s<img %s>\n", indentStr, attrs)
	}
}

func writeHTMLAdmonition(buf *bytes.Buffer, adm *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)
	admType := adm.GetAttribute("type")
	attrs := fmt.Sprintf(`class="admonition admonition-%s"`, html.EscapeString(admType))
	if id := adm.GetAttribute("id"); id != "" {
		attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
	}

	fmt.Fprintf(buf, "%s<div%s>\n", indentStr, attrs)
	fmt.Fprintf(buf, "%s    <p class=\"admonition-title\">%s</p>\n", indentStr, html.EscapeString(strings.ToUpper(admType)))
	for _, child := range adm.Children {
		writeHTMLNode(buf, child, xhtml, indent)
	}
	fmt.Fprintf(buf, "%s</div>\n", indentStr)
}

func writeHTMLExample(buf *bytes.Buffer, ex *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)
	attrs := `class="example"`
	if id := ex.GetAttribute("id"); id != "" {
		attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
	}

	fmt.Fprintf(buf, "%s<div%s>\n", indentStr, attrs)
	if title := ex.GetAttribute("title"); title != "" {
		fmt.Fprintf(buf, "%s    <p class=\"example-title\">%s</p>\n", indentStr, html.EscapeString(title))
	}
	for _, child := range ex.Children {
		writeHTMLNode(buf, child, xhtml, indent)
	}
	fmt.Fprintf(buf, "%s</div>\n", indentStr)
}

func writeHTMLSidebar(buf *bytes.Buffer, sidebar *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)
	attrs := `class="sidebar"`
	if id := sidebar.GetAttribute("id"); id != "" {
		attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
	}

	fmt.Fprintf(buf, "%s<aside%s>\n", indentStr, attrs)
	if title := sidebar.GetAttribute("title"); title != "" {
		fmt.Fprintf(buf, "%s    <p class=\"sidebar-title\">%s</p>\n", indentStr, html.EscapeString(title))
	}
	for _, child := range sidebar.Children {
		writeHTMLNode(buf, child, xhtml, indent)
	}
	fmt.Fprintf(buf, "%s</aside>\n", indentStr)
}

func writeHTMLQuote(buf *bytes.Buffer, quote *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)
	attrs := `class="quote"`
	if id := quote.GetAttribute("id"); id != "" {
		attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
	}

	fmt.Fprintf(buf, "%s<blockquote%s>\n", indentStr, attrs)
	for _, child := range quote.Children {
		writeHTMLNode(buf, child, xhtml, indent)
	}
	if attribution := quote.GetAttribute("attribution"); attribution != "" {
		fmt.Fprintf(buf, "%s    <footer>", indentStr)
		fmt.Fprintf(buf, "<cite>%s</cite>", html.EscapeString(attribution))
		buf.WriteString("</footer>\n")
	}
	if citation := quote.GetAttribute("citation"); citation != "" {
		fmt.Fprintf(buf, "%s    <cite>%s</cite>\n", indentStr, html.EscapeString(citation))
	}
	fmt.Fprintf(buf, "%s</blockquote>\n", indentStr)
}

func writeHTMLInlineContent(buf *bytes.Buffer, node *Node, xhtml bool) {
	for _, child := range node.Children {
		if child.Type == TextNode {
			buf.WriteString(html.EscapeString(child.Data))
		} else {
			switch child.Data {
			case "strong":
				buf.WriteString("<strong>")
				writeHTMLInlineContent(buf, child, xhtml)
				buf.WriteString("</strong>")
			case "emphasis":
				buf.WriteString("<em>")
				writeHTMLInlineContent(buf, child, xhtml)
				buf.WriteString("</em>")
			case "monospace":
				buf.WriteString("<code>")
				writeHTMLInlineContent(buf, child, xhtml)
				buf.WriteString("</code>")
			case "link":
				href := child.GetAttribute("href")
				attrs := fmt.Sprintf(`href="%s"`, html.EscapeString(href))
				if title := child.GetAttribute("title"); title != "" {
					attrs += fmt.Sprintf(` title="%s"`, html.EscapeString(title))
				}
				if window := child.GetAttribute("window"); window != "" {
					attrs += fmt.Sprintf(` target="%s"`, html.EscapeString(window))
				}
				fmt.Fprintf(buf, "<a %s>", attrs)
				writeHTMLInlineContent(buf, child, xhtml)
				buf.WriteString("</a>")
			default:
				writeHTMLInlineContent(buf, child, xhtml)
			}
		}
	}
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
