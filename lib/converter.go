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

// ConvertToXML converts AsciiDoc to XML string using the DOM
func ConvertToXML(reader io.Reader) (string, error) {
	doc, err := ParseDocument(reader)
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

// extractMetadata extracts metadata from a parsed document
func extractMetadata(doc *Node) Metadata {
	meta := Metadata{
		Attributes: make(map[string]string),
	}

	var headerNode *Node
	for _, child := range doc.Children {
		if child.Data == "header" {
			headerNode = child
			break
		}
	}

	if headerNode != nil {
		// Extract title
		titleNode := findChild(headerNode, "h1")
		if titleNode != nil {
			meta.Title = getTextContent(titleNode)
		}

		// Extract author
		for _, child := range headerNode.Children {
			if child.Data == "address" {
				nameNode := findChild(child, "span")
				if nameNode != nil && nameNode.GetAttribute("class") == "author-name" {
					meta.Author = getTextContent(nameNode)
					break // Take first author
				}
			}
		}

		// Extract all attributes
		for _, child := range headerNode.Children {
			if child.Data == "attribute" {
				name := child.GetAttribute("name")
				value := child.GetAttribute("value")
				if name != "" {
					meta.Attributes[":"+name] = value
				}
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

	// Document header (title, authors, etc.) - only if standalone
	var headerNode *Node
	for _, child := range doc.Children {
		if child.Data == "header" {
			headerNode = child
			break
		}
	}

	if opts.Standalone {
		// Write document header in standalone mode
		if headerNode != nil {
			writeHTMLHeaderWithOverride(&buf, headerNode, opts.XHTML, 2, opts.Title, opts.Author)
		}
		buf.WriteString("    <main>\n")
	}

	// Write content (excluding header node in both modes)
	indent := 2
	if !opts.Standalone {
		indent = 0
	}
	for _, child := range doc.Children {
		if child.Data != "header" {
			writeHTMLNode(&buf, child, opts.XHTML, indent)
		}
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
	writeHTMLHeaderWithOverride(buf, header, xhtml, indent, "", "")
}

func writeHTMLHeaderWithOverride(buf *bytes.Buffer, header *Node, xhtml bool, indent int, titleOverride, authorOverride string) {
	indentStr := strings.Repeat("    ", indent)
	
	titleNode := findChild(header, "title")
	titleText := ""
	if titleNode != nil {
		titleText = getTextContent(titleNode)
	}
	
	// Use override if provided
	if titleOverride != "" {
		titleText = titleOverride
	}
	
	if titleText != "" || titleNode != nil {
		buf.WriteString(indentStr + "<header>\n")
		buf.WriteString(indentStr + "      <h1>")
		if titleOverride != "" {
			// Write override title directly (escaped)
			buf.WriteString(html.EscapeString(titleOverride))
		} else if titleNode != nil {
			writeHTMLInlineContent(buf, titleNode, xhtml)
		}
		buf.WriteString("</h1>\n")

		// Authors - use override if provided, otherwise from header
		if authorOverride != "" {
			buf.WriteString(indentStr + "      <address class=\"authors\">\n")
			buf.WriteString(indentStr + "        <p>")
			fmt.Fprintf(buf, `<span class="author-name">%s</span>`, html.EscapeString(authorOverride))
			buf.WriteString("</p>\n")
			buf.WriteString(indentStr + "      </address>\n")
		} else {
			// Authors from header
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
		}

		buf.WriteString(indentStr + "    </header>\n")
	}
}

func writeHTMLNode(buf *bytes.Buffer, node *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)

	switch node.Data {
	case "section":
		writeHTMLSection(buf, node, xhtml, indent)
	case "p":
		writeHTMLParagraph(buf, node, xhtml, indent)
	case "pre":
		writeHTMLPre(buf, node, xhtml, indent)
	case "div":
		class := node.GetAttribute("class")
		if strings.Contains(class, "admonition") {
			writeHTMLAdmonition(buf, node, xhtml, indent)
		} else if strings.Contains(class, "example") {
			writeHTMLExample(buf, node, xhtml, indent)
		} else if strings.Contains(class, "page-break") {
			buf.WriteString(indentStr + `<div class="page-break"></div>` + "\n")
		} else if strings.Contains(class, "labeled-item") {
            // Should be handled inside writeHTMLList, but if encountered here?
            // Just generic block.
            writeHTMLGenericBlock(buf, node, xhtml, indent)
        } else {
			writeHTMLGenericBlock(buf, node, xhtml, indent)
		}
	case "aside":
		writeHTMLSidebar(buf, node, xhtml, indent)
	case "blockquote":
		writeHTMLQuote(buf, node, xhtml, indent)
	case "hr":
		if xhtml {
			buf.WriteString(indentStr + "<hr/>\n")
		} else {
			buf.WriteString(indentStr + "<hr>\n")
		}
	default:
		// Generic handling for other elements (including custom ones like cms-component)
		writeHTMLGenericBlock(buf, node, xhtml, indent)
	}
}

func writeHTMLGenericBlock(buf *bytes.Buffer, node *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)
	
	buf.WriteString(indentStr + "<" + node.Data)
	
	// Sort attributes for stability? Or just iterate
	for k, v := range node.Attributes {
		buf.WriteString(fmt.Sprintf(` %s="%s"`, k, html.EscapeString(v)))
	}

	if len(node.Children) == 0 {
		if isVoidElement(node.Data) {
			if xhtml {
				buf.WriteString("/>\n")
			} else {
				buf.WriteString(">\n")
			}
		} else {
			buf.WriteString("></" + node.Data + ">\n")
		}
		return
	}

	buf.WriteString(">\n")
	for _, child := range node.Children {
		if child.Type == TextNode {
             // If strictly inline text, maybe no newline?
             // But if block contains text, usually indentation is tricky.
             // We'll just write it.
			buf.WriteString(indentStr + "    " + html.EscapeString(child.Data) + "\n")
		} else {
			writeHTMLNode(buf, child, xhtml, indent+1)
		}
	}
	buf.WriteString(indentStr + "</" + node.Data + ">\n")
}

func isVoidElement(tag string) bool {
	switch tag {
	case "area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "param", "source", "track", "wbr":
		return true
	}
	return false
}

func writeHTMLSection(buf *bytes.Buffer, section *Node, xhtml bool, indent int) {
	writeHTMLGenericBlock(buf, section, xhtml, indent)
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

func writeHTMLPre(buf *bytes.Buffer, pre *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)

	// Handle title (outside pre)
	if title := pre.GetAttribute("title"); title != "" {
		fmt.Fprintf(buf, "%s<p class=\"code-title\">%s</p>\n", indentStr, html.EscapeString(title))
	}

	buf.WriteString(indentStr + "<pre")
	for k, v := range pre.Attributes {
		if k == "title" || k == "data-language" { continue } // Skip title, and skip data-language if we want cleaner pre? NO, keep data-language.
		// Actually, let's keep data-language.
        if k == "title" { continue }
		buf.WriteString(fmt.Sprintf(` %s="%s"`, k, html.EscapeString(v)))
	}
    // Add data-language if present and skipped above? No loop handles it.
    
	buf.WriteString(">")

	// Check for code child
	var codeChild *Node
	for _, child := range pre.Children {
		if child.Data == "code" {
			codeChild = child
			break
		}
	}

	if codeChild != nil {
		buf.WriteString("<code")
		for k, v := range codeChild.Attributes {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, k, html.EscapeString(v)))
		}
		buf.WriteString(">")

		for _, grandChild := range codeChild.Children {
			if grandChild.Type == TextNode {
				buf.WriteString(html.EscapeString(grandChild.Data))
			}
		}

		buf.WriteString("</code>")
	} else {
		for _, child := range pre.Children {
			if child.Type == TextNode {
				buf.WriteString(html.EscapeString(child.Data))
			}
		}
	}

	buf.WriteString("</pre>\n")
}

func writeHTMLAdmonition(buf *bytes.Buffer, adm *Node, xhtml bool, indent int) {
	indentStr := strings.Repeat("    ", indent)
	class := adm.GetAttribute("class")
	admType := "NOTE"
	parts := strings.Split(class, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "admonition-") {
			admType = strings.TrimPrefix(part, "admonition-")
			break
		}
	}

	attrs := fmt.Sprintf(`class="%s"`, html.EscapeString(class))
	if id := adm.GetAttribute("id"); id != "" {
		attrs += fmt.Sprintf(` id="%s"`, html.EscapeString(id))
	}

	fmt.Fprintf(buf, "%s<div%s>\n", indentStr, attrs)
	fmt.Fprintf(buf, "%s    <p class=\"admonition-title\">%s</p>\n", indentStr, html.EscapeString(strings.ToUpper(admType)))
	for _, child := range adm.Children {
		writeHTMLNode(buf, child, xhtml, indent+1)
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
			case "em":
				buf.WriteString("<em>")
				writeHTMLInlineContent(buf, child, xhtml)
				buf.WriteString("</em>")
			case "code":
				buf.WriteString("<code>")
				writeHTMLInlineContent(buf, child, xhtml)
				buf.WriteString("</code>")
			case "a":
				href := child.GetAttribute("href")
				attrs := fmt.Sprintf(`href="%s"`, html.EscapeString(href))
				if title := child.GetAttribute("title"); title != "" {
					attrs += fmt.Sprintf(` title="%s"`, html.EscapeString(title))
				}
				if class := child.GetAttribute("class"); class != "" {
					attrs += fmt.Sprintf(` class="%s"`, html.EscapeString(class))
				}
				if window := child.GetAttribute("window"); window != "" {
					attrs += fmt.Sprintf(` target="%s"`, html.EscapeString(window))
				}
				fmt.Fprintf(buf, "<a %s>", attrs)
				writeHTMLInlineContent(buf, child, xhtml)
				buf.WriteString("</a>")
			default:
				if child.Type == ElementNode {
					buf.WriteString("<" + child.Data)
					for k, v := range child.Attributes {
						buf.WriteString(fmt.Sprintf(` %s="%s"`, k, html.EscapeString(v)))
					}
					buf.WriteString(">")
					writeHTMLInlineContent(buf, child, xhtml)
					buf.WriteString("</" + child.Data + ">")
				} else {
					writeHTMLInlineContent(buf, child, xhtml)
				}
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

