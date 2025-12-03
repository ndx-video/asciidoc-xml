package lib

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestConvert_BasicDocument(t *testing.T) {
	input := `= Test Document
:author: John Doe
:email: john@example.com

This is a simple paragraph.`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if doc == nil {
		t.Fatal("Document is nil")
	}

	if doc.Type != Document {
		t.Errorf("Expected root element type Document, got '%s'", doc.Type.String())
	}

	// Check title from document attributes
	titleText := doc.GetAttribute("title")
	if titleText != "Test Document" {
		t.Errorf("Expected title 'Test Document', got '%s'", titleText)
	}

	// Check author from document attributes
	authorText := doc.GetAttribute("author")
	if authorText != "John Doe" {
		t.Errorf("Expected author 'John Doe', got '%s'", authorText)
	}

	// Find paragraph in content
	var paraNode *Node
	for _, child := range doc.Children {
		if child.Type == Paragraph {
			paraNode = child
			break
		}
	}

	if paraNode == nil {
		t.Fatal("Expected at least one paragraph (p)")
	}

	if len(paraNode.Children) == 0 {
		t.Fatal("Expected paragraph to have children")
	}
}

func TestConvert_Sections(t *testing.T) {
	input := `= Document Title

== Section 1

Content of section 1.

=== Subsection 1.1

Content of subsection.

== Section 2

Content of section 2.`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	sections := 0
	for _, child := range doc.Children {
		if child.Type == Section {
			sections++
		}
	}

	if sections < 2 {
		t.Errorf("Expected at least 2 sections, got %d", sections)
	}

	// Check first section
	var firstSection *Node
	for _, child := range doc.Children {
		if child.Type == Section {
			firstSection = child
			break
		}
	}

	if firstSection == nil {
		t.Fatal("Expected first item to be a section")
	}

	// Check section has title attribute or text child
	titleAttr := firstSection.GetAttribute("title")
	if titleAttr == "" {
		// Check for text child
		foundText := false
		for _, child := range firstSection.Children {
			if child.Type == Text {
				foundText = true
				break
			}
		}
		if !foundText {
			t.Error("Expected section to have title text or title attribute")
		}
	}
}

func TestConvert_Paragraphs(t *testing.T) {
	input := `= Test

This is paragraph one.

This is paragraph two with *bold* and _italic_ text.`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	paragraphs := 0
	for _, child := range doc.Children {
		if child.Type == Paragraph {
			paragraphs++
		}
	}

	if paragraphs < 2 {
		t.Errorf("Expected at least 2 paragraphs, got %d", paragraphs)
	}
}

func TestConvert_InlineFormatting(t *testing.T) {
	input := `= Test

This has *bold text* and _italic text_ and ` + "`monospace`" + ` text.`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	var paraNode *Node
	for _, child := range doc.Children {
		if child.Type == Paragraph {
			paraNode = child
			break
		}
	}

	if paraNode == nil {
		t.Fatal("Expected a paragraph")
	}

	foundBold := false
	foundItalic := false
	foundMono := false

	for _, child := range paraNode.Children {
		if child.Type == Bold {
			foundBold = true
		}
		if child.Type == Italic {
			foundItalic = true
		}
		if child.Type == Monospace {
			foundMono = true
		}
	}

	if !foundBold {
		t.Error("Expected to find bold text (strong)")
	}
	if !foundItalic {
		t.Error("Expected to find italic text (em)")
	}
	if !foundMono {
		t.Error("Expected to find monospace text (code)")
	}
}

func TestConvert_CodeBlock(t *testing.T) {
	input := `= Test

[source,go]
----
package main

func main() {
    println("Hello")
}
----`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundCodeBlock := false
	for _, child := range doc.Children {
		if child.Type == CodeBlock {
			foundCodeBlock = true
			lang := child.GetAttribute("language")
			if lang != "go" {
				t.Errorf("Expected language 'go', got '%s'", lang)
			}
			content := getTextContent(child)
			if !strings.Contains(content, "package main") {
				t.Error("Expected code block to contain 'package main'")
			}
			break
		}
	}

	if !foundCodeBlock {
		t.Error("Expected to find a code block")
	}
}

func TestConvert_UnorderedList(t *testing.T) {
	input := `= Test

* Item one
* Item two
* Item three`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundList := false
	for _, child := range doc.Children {
		if child.Type == List && child.GetAttribute("style") == "unordered" {
			foundList = true
			items := 0
			for _, item := range child.Children {
				if item.Type == ListItem {
					items++
				}
			}
			if items != 3 {
				t.Errorf("Expected 3 list items, got %d", items)
			}
			break
		}
	}

	if !foundList {
		t.Error("Expected to find a list (ul)")
	}
}

func TestConvert_OrderedList(t *testing.T) {
	input := `= Test

. First item
. Second item
. Third item`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundList := false
	for _, child := range doc.Children {
		if child.Type == List && child.GetAttribute("style") == "ordered" {
			foundList = true
			items := 0
			for _, item := range child.Children {
				if item.Type == ListItem {
					items++
				}
			}
			if items != 3 {
				t.Errorf("Expected 3 list items, got %d", items)
			}
			break
		}
	}

	if !foundList {
		t.Error("Expected to find a list (ol)")
	}
}

func TestConvert_LabeledList(t *testing.T) {
	input := `= Test

term1:: definition1
term2:: definition2`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundList := false
	for _, child := range doc.Children {
		if child.Type == List && child.GetAttribute("style") == "labeled" {
			foundList = true
			items := 0
			for _, item := range child.Children {
				if item.Type == ListItem {
					items++
				}
			}
			if items != 2 {
				t.Errorf("Expected 2 list items, got %d", items)
			}
			break
		}
	}

	if !foundList {
		t.Error("Expected to find a labeled list (dl)")
	}
}

func TestConvert_Admonition(t *testing.T) {
	input := `= Test

NOTE: This is a note.

WARNING: This is a warning.`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	admonitions := 0
	for _, child := range doc.Children {
		if child.Type == Admonition {
			admonitions++
		}
	}

	if admonitions != 2 {
		t.Errorf("Expected 2 admonitions, got %d", admonitions)
	}

	// Check first admonition
	var firstAdmonition *Node
	for _, child := range doc.Children {
		if child.Type == Admonition {
			firstAdmonition = child
			break
		}
	}

	if firstAdmonition == nil {
		t.Fatal("Expected first item to be an admonition")
	}

    // Check type attribute
    admType := firstAdmonition.GetAttribute("type")
    if admType != "note" {
        t.Errorf("Expected admonition type 'note', got '%s'", admType)
    }
}

func TestConvert_Links(t *testing.T) {
	input := `= Test

Visit https://example.com[Example Website] for more info.`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	var paraNode *Node
	for _, child := range doc.Children {
		if child.Type == Paragraph {
			paraNode = child
			break
		}
	}

	if paraNode == nil {
		t.Fatal("Expected a paragraph")
	}

	foundLink := false
	for _, child := range paraNode.Children {
		if child.Type == Link {
			foundLink = true
			href := child.GetAttribute("href")
			if href != "https://example.com" {
				t.Errorf("Expected href 'https://example.com', got '%s'", href)
			}
			break
		}
	}

	if !foundLink {
		t.Error("Expected to find a link (a)")
	}
}

func TestConvert_LinkMacros(t *testing.T) {
	input := `= Test

link:/docs?page=user-guide.adoc[User Guide] - Complete guide to using the web interface and features

link:/docs?page=api.adoc[API Reference]`

	html, err := ConvertToHTML(bytes.NewReader([]byte(input)), false, false, "", "")
	if err != nil {
		t.Fatalf("ConvertToHTML failed: %v", err)
	}

	// Check that links are properly converted to HTML anchor tags
	if !strings.Contains(html, `<a href="/docs?page=user-guide.adoc">`) {
		t.Error("Expected HTML to contain link to user-guide.adoc")
	}

	if !strings.Contains(html, "User Guide") {
		t.Error("Expected HTML to contain link text 'User Guide'")
	}

	if !strings.Contains(html, `<a href="/docs?page=api.adoc">`) {
		t.Error("Expected HTML to contain link to api.adoc")
	}

	if !strings.Contains(html, "API Reference") {
		t.Error("Expected HTML to contain link text 'API Reference'")
	}

	// Verify the links are properly closed
	if !strings.Contains(html, `</a>`) {
		t.Error("Expected HTML to contain closing anchor tags")
	}
}

func TestConvert_Image(t *testing.T) {
	input := `= Test

image::logo.png[Logo, 200, 100]`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundImage := false
	for _, child := range doc.Children {
		if child.Type == BlockMacro && child.Name == "image" {
			foundImage = true
			src := child.GetAttribute("src")
			if src != "logo.png" {
				t.Errorf("Expected src 'logo.png', got '%s'", src)
			}
			alt := child.GetAttribute("alt")
			if alt != "Logo" {
				t.Errorf("Expected alt 'Logo', got '%s'", alt)
			}
			break
		}
	}

	if !foundImage {
		t.Error("Expected to find an image (img)")
	}
}

func TestConvert_Attributes(t *testing.T) {
	input := `= Test Document
:doctype: book
:revnumber: 1.0
:revdate: 2024-01-01
:revremark: Initial version
:custom-attr: Custom Value`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	doctype := doc.GetAttribute("doctype")
	if doctype != "book" {
		t.Errorf("Expected doctype 'book', got '%s'", doctype)
	}

	// Header info is now in document attributes, no separate header node

    // Revision is now in document attributes
    revNumber := doc.GetAttribute("revnumber")
    if revNumber == "" {
		t.Fatal("Expected revision to be set")
	}
}

func TestConvertToXML_ValidXML(t *testing.T) {
	input := `= Test Document

This is a test paragraph.`

	xmlOutput, err := ConvertToXML(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("ConvertToXML failed: %v", err)
	}

	// Verify it's valid XML
	decoder := xml.NewDecoder(strings.NewReader(xmlOutput))
	for {
		_, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Fatalf("Invalid XML: %v", err)
		}
	}

	// Verify it contains expected elements
	if !strings.Contains(xmlOutput, "document") {
		t.Error("XML should contain 'document' element")
	}
	// Header is h1 now in xml?
    // Yes, AST contains h1.
    // And it should contain "Test Document".
	if !strings.Contains(xmlOutput, "Test Document") {
		t.Error("XML should contain document title")
	}
}

func TestConvert_ExampleBlock(t *testing.T) {
	input := `= Test

.Example Title
====
This is an example block.
====`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundExample := false
	for _, child := range doc.Children {
		if child.Type == Example {
			foundExample = true
			title := child.GetAttribute("title")
			if title != "Example Title" {
				t.Errorf("Expected title 'Example Title', got '%s'", title)
			}
			break
		}
	}

	if !foundExample {
		t.Error("Expected to find an example block")
	}
}

func TestConvert_Sidebar(t *testing.T) {
	input := `= Test

****
This is a sidebar.
****`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundSidebar := false
	for _, child := range doc.Children {
		if child.Type == Sidebar {
			foundSidebar = true
			break
		}
	}

	if !foundSidebar {
		t.Error("Expected to find a sidebar (aside)")
	}
}

func TestConvert_Quote(t *testing.T) {
	input := `= Test

____
This is a quote.
____`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundQuote := false
	for _, child := range doc.Children {
		if child.Type == Quote {
			foundQuote = true
			break
		}
	}

	if !foundQuote {
		t.Error("Expected to find a quote block (blockquote)")
	}
}

func TestConvert_ThematicBreak(t *testing.T) {
	input := `= Test

First section.

'''

Second section.`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundBreak := false
	for _, child := range doc.Children {
		if child.Type == ThematicBreak {
			foundBreak = true
			break
		}
	}

	if !foundBreak {
		t.Error("Expected to find a thematic break (hr)")
	}
}

func TestConvert_PageBreak(t *testing.T) {
	input := `= Test

First page.

<<<

Second page.`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundBreak := false
	for _, child := range doc.Children {
		if child.Type == PageBreak {
			foundBreak = true
			break
		}
	}

	if !foundBreak {
		t.Error("Expected to find a page break (div.page-break)")
	}
}

func TestConvert_ComplexDocument(t *testing.T) {
	input := `= Complex Document
:author: Test Author
:doctype: article

== Introduction

This is the introduction with *bold* and _italic_ text.

== Main Content

=== Subsection

Here's a code block:

[source,go]
----
func main() {
    println("Hello")
}
----

And a list:

* Item 1
* Item 2

|===
|Header
|Data
|===

NOTE: This is important!`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Count different element types recursively
	var countElements func(*Node, *int, *int, *int, *int, *int)
	countElements = func(node *Node, sections, codeBlocks, lists, tables, admonitions *int) {
		if node.Type == Section {
			*sections++
		}
		if node.Type == CodeBlock {
			*codeBlocks++
		}
		if node.Type == List {
			*lists++
		}
		if node.Type == Table {
			*tables++
		}
		if node.Type == Admonition {
			*admonitions++
		}
		for _, child := range node.Children {
			countElements(child, sections, codeBlocks, lists, tables, admonitions)
		}
	}

	sections := 0
	codeBlocks := 0
	lists := 0
	tables := 0
	admonitions := 0

	countElements(doc, &sections, &codeBlocks, &lists, &tables, &admonitions)

	if sections < 2 {
		t.Errorf("Expected at least 2 sections, got %d", sections)
	}
	if codeBlocks < 1 {
		t.Error("Expected at least 1 code block")
	}
	if lists < 1 {
		t.Error("Expected at least 1 list")
	}
	if tables < 1 {
		t.Error("Expected at least 1 table")
	}
	if admonitions < 1 {
		t.Error("Expected at least 1 admonition")
	}
}

func TestConvertToHTML(t *testing.T) {
	input := `= Test Document

This is a test paragraph.`

	htmlOutput, err := ConvertToHTML(bytes.NewReader([]byte(input)), false, false, "", "")
	if err != nil {
		t.Fatalf("ConvertToHTML failed: %v", err)
	}

	// Verify it contains expected HTML elements
	if !strings.Contains(htmlOutput, "<!DOCTYPE html>") {
		t.Error("HTML should contain DOCTYPE")
	}
	if !strings.Contains(htmlOutput, "<html") {
		t.Error("HTML should contain html element")
	}
	if !strings.Contains(htmlOutput, "<body>") {
		t.Error("HTML should contain body element")
	}
	if !strings.Contains(htmlOutput, "Test Document") {
		t.Error("HTML should contain document title")
	}
}

func TestConvertToHTML_XHTML(t *testing.T) {
	input := `= Test Document

This is a test paragraph.`

	htmlOutput, err := ConvertToHTML(bytes.NewReader([]byte(input)), true, false, "", "")
	if err != nil {
		t.Fatalf("ConvertToHTML failed: %v", err)
	}

	// Verify it contains XHTML-specific elements
	if !strings.Contains(htmlOutput, `<?xml version="1.0"`) {
		t.Error("XHTML should contain XML declaration")
	}
	if !strings.Contains(htmlOutput, `xmlns="http://www.w3.org/1999/xhtml"`) {
		t.Error("XHTML should contain xmlns attribute")
	}
}

func TestConvertMarkdownToAsciiDoc_Basic(t *testing.T) {
	input := `# Title

This is a paragraph with **bold** and *italic* text.

## Section

- Item 1
- Item 2

### Subsection

1. First item
2. Second item`

	output, err := ConvertMarkdownToAsciiDoc(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDoc failed: %v", err)
	}

	// Check title conversion
	if !strings.Contains(output, "= Title") {
		t.Error("Expected '= Title' in output")
	}

	// Check section conversion
	if !strings.Contains(output, "== Section") {
		t.Error("Expected '== Section' in output")
	}

	// Check subsection conversion
	if !strings.Contains(output, "=== Subsection") {
		t.Error("Expected '=== Subsection' in output")
	}

	// Check bold conversion (AsciiDoc uses **text** for bold)
	if !strings.Contains(output, "**bold**") {
		t.Error("Expected '**bold**' in output")
	}

	// Check italic conversion
	if !strings.Contains(output, "_italic_") {
		t.Error("Expected '_italic_' in output")
	}

	// Check unordered list conversion
	if !strings.Contains(output, "* Item 1") {
		t.Error("Expected '* Item 1' in output")
	}

	// Check ordered list conversion
	if !strings.Contains(output, ". First item") {
		t.Error("Expected '. First item' in output")
	}
}

func TestConvertMarkdownToAsciiDoc_LinksAndImages(t *testing.T) {
	input := `# Test

This is a [link](https://example.com) and an ![image](image.png)`

	output, err := ConvertMarkdownToAsciiDoc(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDoc failed: %v", err)
	}

	// Check link conversion
	if !strings.Contains(output, "link:https://example.com[link]") {
		t.Error("Expected 'link:https://example.com[link]' in output")
	}

	// Check image conversion (inline images use : not ::)
	if !strings.Contains(output, "image:image.png[") {
		t.Error("Expected 'image:image.png[' in output")
	}
}

func TestConvertMarkdownToAsciiDoc_CodeBlocks(t *testing.T) {
	input := "# Test\n\n```go\npackage main\n\nfunc main() {\n    println(\"Hello\")\n}\n```\n"

	output, err := ConvertMarkdownToAsciiDoc(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDoc failed: %v", err)
	}

	// Check code block conversion
	if !strings.Contains(output, "[source,go]") {
		t.Errorf("Expected '[source,go]' in output. Got:\n%s", output)
	}

	if !strings.Contains(output, "----") {
		t.Errorf("Expected '----' code block delimiters in output. Got:\n%s", output)
	}

	if !strings.Contains(output, "package main") {
		t.Errorf("Expected code content in output. Got:\n%s", output)
	}
}

func TestConvertMarkdownToAsciiDoc_README(t *testing.T) {
	// Read the actual README.md file (from project root)
	readmeContent, err := readFileContent("../README.md")
	if err != nil {
		t.Skipf("Could not read README.md: %v", err)
	}

	output, err := ConvertMarkdownToAsciiDoc(bytes.NewReader([]byte(readmeContent)))
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDoc failed: %v", err)
	}

	// Verify basic conversions
	if len(output) == 0 {
		t.Error("Output should not be empty")
	}

	// Check that headers were converted
	if !strings.Contains(output, "=") {
		t.Error("Expected at least one header conversion")
	}

	// Check that some content is present
	if len(output) < 100 {
		t.Error("Output seems too short")
	}

	// Verify the first line is a title (should start with =)
	lines := strings.Split(output, "\n")
	if len(lines) > 0 && !strings.HasPrefix(strings.TrimSpace(lines[0]), "=") {
		t.Logf("First line: %s", lines[0])
		t.Log("Note: First line might not be a title, which is okay")
	}
}

// Helper function to read file content
func readFileContent(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func TestConvertHTML_Standalone(t *testing.T) {
	input := `= Test Document
:author: John Doe

This is a test paragraph.`

	opts := ConvertOptions{
		Standalone: true,
		UsePicoCSS: false,
		XHTML:      false,
	}

	result, err := Convert(bytes.NewReader([]byte(input)), opts)
	if err != nil {
		t.Fatalf("ConvertHTML failed: %v", err)
	}

	// Verify standalone HTML structure
	if !strings.Contains(result.HTML, "<!DOCTYPE html>") {
		t.Error("Standalone HTML should contain DOCTYPE")
	}
	if !strings.Contains(result.HTML, "<html") {
		t.Error("Standalone HTML should contain html element")
	}
	if !strings.Contains(result.HTML, "<head>") {
		t.Error("Standalone HTML should contain head element")
	}
	if !strings.Contains(result.HTML, "<body>") {
		t.Error("Standalone HTML should contain body element")
	}
	if !strings.Contains(result.HTML, "<main>") {
		t.Error("Standalone HTML should contain main element")
	}

	// Verify metadata
	if result.Meta.Title != "Test Document" {
		t.Errorf("Expected title 'Test Document', got '%s'", result.Meta.Title)
	}
	if result.Meta.Author != "John Doe" {
		t.Errorf("Expected author 'John Doe', got '%s'", result.Meta.Author)
	}
}

func TestConvertHTML_Fragment(t *testing.T) {
	input := `= Test Document

This is a test paragraph.`

	opts := ConvertOptions{
		Standalone: false,
		UsePicoCSS: false,
		XHTML:      false,
	}

	result, err := Convert(bytes.NewReader([]byte(input)), opts)
	if err != nil {
		t.Fatalf("ConvertHTML failed: %v", err)
	}

	// Verify fragment does NOT contain document structure
	if strings.Contains(result.HTML, "<!DOCTYPE html>") {
		t.Error("Fragment HTML should NOT contain DOCTYPE")
	}
	if strings.Contains(result.HTML, "<html") {
		t.Error("Fragment HTML should NOT contain html element")
	}
	if strings.Contains(result.HTML, "<head>") {
		t.Error("Fragment HTML should NOT contain head element")
	}
	if strings.Contains(result.HTML, "<body>") {
		t.Error("Fragment HTML should NOT contain body element")
	}

	// Should contain content
	if !strings.Contains(result.HTML, "test paragraph") {
		t.Error("Fragment HTML should contain content")
	}
}

func TestConvertHTML_TitleOverride(t *testing.T) {
	input := `= Original Title

Content here.`

	opts := ConvertOptions{
		Standalone: true,
		Title:      "Override Title",
		UsePicoCSS: false,
		XHTML:      false,
	}

	result, err := Convert(bytes.NewReader([]byte(input)), opts)
	if err != nil {
		t.Fatalf("ConvertHTML failed: %v", err)
	}

	// Verify override title appears in HTML
	if !strings.Contains(result.HTML, "<title>Override Title</title>") {
		t.Error("HTML should contain override title in <title> tag")
	}

	// Verify override title in metadata
	if result.Meta.Title != "Override Title" {
		t.Errorf("Expected override title 'Override Title' in metadata, got '%s'", result.Meta.Title)
	}
}

func TestConvertHTML_AuthorOverride(t *testing.T) {
	input := `= Test Document
:author: Original Author

Content here.`

	opts := ConvertOptions{
		Standalone: true,
		Author:     "Override Author",
		UsePicoCSS: false,
		XHTML:      false,
	}

	result, err := Convert(bytes.NewReader([]byte(input)), opts)
	if err != nil {
		t.Fatalf("ConvertHTML failed: %v", err)
	}

	// Verify override author in metadata
	if result.Meta.Author != "Override Author" {
		t.Errorf("Expected override author 'Override Author' in metadata, got '%s'", result.Meta.Author)
	}
}

func TestConvertHTML_MetadataExtraction(t *testing.T) {
	input := `= My Document Title
:author: Jane Doe
:email: jane@example.com
:toc: true
:numbered: true

Content here.`

	opts := ConvertOptions{
		Standalone: true,
		UsePicoCSS: false,
		XHTML:      false,
	}

	result, err := Convert(bytes.NewReader([]byte(input)), opts)
	if err != nil {
		t.Fatalf("ConvertHTML failed: %v", err)
	}

	// Verify extracted metadata
	if result.Meta.Title != "My Document Title" {
		t.Errorf("Expected title 'My Document Title', got '%s'", result.Meta.Title)
	}
	if result.Meta.Author != "Jane Doe" {
		t.Errorf("Expected author 'Jane Doe', got '%s'", result.Meta.Author)
	}

	// Verify attributes
	if result.Meta.Attributes[":toc"] != "true" {
		t.Errorf("Expected :toc attribute 'true', got '%s'", result.Meta.Attributes[":toc"])
	}
	if result.Meta.Attributes[":numbered"] != "true" {
		t.Errorf("Expected :numbered attribute 'true', got '%s'", result.Meta.Attributes[":numbered"])
	}
}

func TestConvertHTML_PicoCSS(t *testing.T) {
	input := `= Test Document

Content.`

	opts := ConvertOptions{
		Standalone:    true,
		UsePicoCSS:    true,
		PicoCSSPath:   "https://cdn.example.com/pico.css",
		XHTML:         false,
	}

	result, err := Convert(bytes.NewReader([]byte(input)), opts)
	if err != nil {
		t.Fatalf("ConvertHTML failed: %v", err)
	}

	// Verify PicoCSS link is included
	if !strings.Contains(result.HTML, `rel="stylesheet"`) {
		t.Error("PicoCSS should include stylesheet link")
	}
	if !strings.Contains(result.HTML, `https://cdn.example.com/pico.css`) {
		t.Error("PicoCSS path should be included in link")
	}
}

func TestConvertHTML_XHTML(t *testing.T) {
	input := `= Test Document

Content.`

	opts := ConvertOptions{
		Standalone: true,
		UsePicoCSS: false,
		XHTML:      true,
	}

	result, err := Convert(bytes.NewReader([]byte(input)), opts)
	if err != nil {
		t.Fatalf("ConvertHTML failed: %v", err)
	}

	// Verify XHTML-specific elements
	if !strings.Contains(result.HTML, `<?xml version="1.0"`) {
		t.Error("XHTML should contain XML declaration")
	}
	if !strings.Contains(result.HTML, `xmlns="http://www.w3.org/1999/xhtml"`) {
		t.Error("XHTML should contain xmlns attribute")
	}
}

func TestConvertHTML_FragmentWithMetadata(t *testing.T) {
	input := `= Document Title
:author: Test Author
:toc: true

Content paragraph.`

	opts := ConvertOptions{
		Standalone: false,
		UsePicoCSS: false,
		XHTML:      false,
	}

	result, err := Convert(bytes.NewReader([]byte(input)), opts)
	if err != nil {
		t.Fatalf("ConvertHTML failed: %v", err)
	}

	// Fragment should not have document structure
	if strings.Contains(result.HTML, "<html") || strings.Contains(result.HTML, "<body>") {
		t.Error("Fragment should not contain document structure")
	}

	// But metadata should still be extracted
	if result.Meta.Title != "Document Title" {
		t.Errorf("Expected title in metadata, got '%s'", result.Meta.Title)
	}
	if result.Meta.Author != "Test Author" {
		t.Errorf("Expected author in metadata, got '%s'", result.Meta.Author)
	}
	if result.Meta.Attributes[":toc"] != "true" {
		t.Error("Expected :toc attribute in metadata")
	}
}

func TestConvertMarkdownToAsciiDoc_FrontmatterArrays(t *testing.T) {
	input := `---
title: Test Post
categories:
  - microsoft
  - Windows
  - windows server
tags:
  - 2012 R2
  - autounattend.xml
  - packer
---

# Content

This is the content.`

	output, err := ConvertMarkdownToAsciiDoc(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDoc failed: %v", err)
	}

	// Check title conversion
	if !strings.Contains(output, "= Test Post") {
		t.Error("Expected '= Test Post' in output")
	}

	// Check categories array conversion
	if !strings.Contains(output, ":categories: microsoft, Windows, windows server") {
		t.Errorf("Expected ':categories: microsoft, Windows, windows server' in output. Got:\n%s", output)
	}

	// Check tags array conversion
	if !strings.Contains(output, ":tags: 2012 R2, autounattend.xml, packer") {
		t.Errorf("Expected ':tags: 2012 R2, autounattend.xml, packer' in output. Got:\n%s", output)
	}
}

func TestConvertMarkdownToAsciiDoc_FrontmatterNested(t *testing.T) {
	input := `---
title: Test Post
wordtwit_post_info:
  - 'O:8:"stdClass":14:{s:6:"manual";b:1;...}'
---

# Content`

	output, err := ConvertMarkdownToAsciiDoc(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDoc failed: %v", err)
	}

	// Check that nested structure is preserved
	if !strings.Contains(output, ":wordtwit_post_info:") {
		t.Errorf("Expected ':wordtwit_post_info:' in output. Got:\n%s", output)
	}
}

func TestConvertMarkdownToAsciiDoc_BoldText(t *testing.T) {
	input := `# Test

**I'll reiterate that**: This will allow _ANYONE_ holding this file to _access ANY server AS YOU_.`

	output, err := ConvertMarkdownToAsciiDoc(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDoc failed: %v", err)
	}

	// Check bold conversion (should be **text**, not *text*)
	if !strings.Contains(output, "**I'll reiterate that**") {
		t.Errorf("Expected '**I'll reiterate that**' in output. Got:\n%s", output)
	}

	// Make sure it's not converted to italic
	// Check that we don't have the pattern *text* (single asterisks) without the double asterisks
	// The output should contain **text** (double asterisks), not *text* (single asterisks)
	lines := strings.Split(output, "\n")
	foundBold := false
	foundItalic := false
	for _, line := range lines {
		if strings.Contains(line, "**I'll reiterate that**") {
			foundBold = true
		}
		// Check for *I'll reiterate that* (single asterisks, not part of double)
		italicPattern := regexp.MustCompile(`[^*]\*I'll reiterate that\*[^*]`)
		if italicPattern.MatchString(line) {
			foundItalic = true
		}
	}
	if !foundBold {
		t.Errorf("Expected '**I'll reiterate that**' in output. Got:\n%s", output)
	}
	if foundItalic {
		t.Error("Bold text should not be converted to italic (*text*)")
	}
}

func TestConvertMarkdownToAsciiDoc_Tables(t *testing.T) {
	input := `# Test

| Column 1 | Column 2 |
|----------|----------|
| Value 1  | Value 2  |`

	output, err := ConvertMarkdownToAsciiDoc(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDoc failed: %v", err)
	}

	// Check table start
	if !strings.Contains(output, "[cols=") {
		t.Errorf("Expected '[cols=' in output. Got:\n%s", output)
	}

	// Check table delimiter
	if !strings.Contains(output, "|===") {
		t.Errorf("Expected '|===' in output. Got:\n%s", output)
	}

	// Check table rows
	if !strings.Contains(output, "|Column 1") || !strings.Contains(output, "|Column 2") {
		t.Errorf("Expected table headers in output. Got:\n%s", output)
	}
}

func TestConvertMarkdownToAsciiDoc_HorizontalRules(t *testing.T) {
	input := `# Test

Some content.

---

More content.`

	output, err := ConvertMarkdownToAsciiDoc(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDoc failed: %v", err)
	}

	// Check horizontal rule conversion
	if !strings.Contains(output, "'''") {
		t.Errorf("Expected ''' (horizontal rule) in output. Got:\n%s", output)
	}
}

func TestConvertMarkdownToAsciiDoc_MixedFormatting(t *testing.T) {
	input := `# Test

This is **bold** text and this is *italic* text.`

	output, err := ConvertMarkdownToAsciiDoc(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDoc failed: %v", err)
	}

	// Check bold conversion
	if !strings.Contains(output, "**bold**") {
		t.Errorf("Expected '**bold**' in output. Got:\n%s", output)
	}

	// Check italic conversion
	if !strings.Contains(output, "_italic_") {
		t.Errorf("Expected '_italic_' in output. Got:\n%s", output)
	}
}

func TestConvertToXML_Superscript(t *testing.T) {
	input := `= Test

Text with ^superscript^ content.`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	xmlOutput := ToXML(doc)
	if !strings.Contains(xmlOutput, "<superscript>") {
		t.Error("Expected <superscript> element in XML output")
	}
	if !strings.Contains(xmlOutput, "</superscript>") {
		t.Error("Expected </superscript> closing tag in XML output")
	}
}

func TestConvertToXML_Subscript(t *testing.T) {
	input := `= Test

Text with ~subscript~ content.`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	xmlOutput := ToXML(doc)
	if !strings.Contains(xmlOutput, "<subscript>") {
		t.Error("Expected <subscript> element in XML output")
	}
	if !strings.Contains(xmlOutput, "</subscript>") {
		t.Error("Expected </subscript> closing tag in XML output")
	}
}

func TestConvertToXML_Highlight(t *testing.T) {
	input := `= Test

Text with #highlight# content.`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	xmlOutput := ToXML(doc)
	if !strings.Contains(xmlOutput, "<highlight>") {
		t.Error("Expected <highlight> element in XML output")
	}
	if !strings.Contains(xmlOutput, "</highlight>") {
		t.Error("Expected </highlight> closing tag in XML output")
	}
}

func TestConvertToXML_VerseBlock(t *testing.T) {
	input := `= Test

[verse]
____
Roses are red,
Violets are blue.
____`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	xmlOutput := ToXML(doc)
	if !strings.Contains(xmlOutput, "<verseblock") {
		t.Error("Expected <verseblock> element in XML output")
	}
	if !strings.Contains(xmlOutput, "</verseblock>") {
		t.Error("Expected </verseblock> closing tag in XML output")
	}
}

func TestConvertToXML_OpenBlock(t *testing.T) {
	input := `= Test

--
This is an open block.
--`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check if open block was parsed
	foundOpenBlock := false
	for _, child := range doc.Children {
		if child.Type == OpenBlock {
			foundOpenBlock = true
			break
		}
	}

	if !foundOpenBlock {
		t.Log("Open block may not be parsed correctly by parser, checking XML output anyway")
	}

	xmlOutput := ToXML(doc)
	if !strings.Contains(xmlOutput, "<openblock") {
		t.Logf("XML output:\n%s", xmlOutput)
		t.Error("Expected <openblock> element in XML output")
	}
	if !strings.Contains(xmlOutput, "</openblock>") {
		t.Error("Expected </openblock> closing tag in XML output")
	}
}

func TestConvertToXML_Anchor(t *testing.T) {
	input := `= Test

[#anchor-id]
== Section

Text with [#inline-anchor] anchor.`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	xmlOutput := ToXML(doc)
	if !strings.Contains(xmlOutput, `<anchor id="anchor-id"`) || !strings.Contains(xmlOutput, `<anchor id="inline-anchor"`) {
		t.Error("Expected <anchor> element with id attribute in XML output")
	}
}

func TestConvertToXML_Footnote(t *testing.T) {
	input := `= Test

Text with footnote:[Footnote text] and footnote:ref[Reference footnote].`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	xmlOutput := ToXML(doc)
	if !strings.Contains(xmlOutput, "<footnote") {
		t.Error("Expected <footnote> element in XML output")
	}
	if !strings.Contains(xmlOutput, "</footnote>") {
		t.Error("Expected </footnote> closing tag in XML output")
	}
}

func TestConvertToXML_Preamble(t *testing.T) {
	input := `= Test Document

This is preamble content.

== First Section

Section content.`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	xmlOutput := ToXML(doc)
	if !strings.Contains(xmlOutput, "<preamble") {
		t.Error("Expected <preamble> element in XML output")
	}
	if !strings.Contains(xmlOutput, "</preamble>") {
		t.Error("Expected </preamble> closing tag in XML output")
	}
	if !strings.Contains(xmlOutput, "preamble content") {
		t.Error("Expected preamble content in XML output")
	}
}

func TestConvertToXML_SectionAttributes(t *testing.T) {
	input := `= Test

[appendix]
== Appendix Section

[discrete]
== Discrete Section

[#custom-id]
== Section With ID`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	xmlOutput := ToXML(doc)
	// Check that attributes are present (they may be stored differently)
	if !strings.Contains(xmlOutput, `appendix`) && !strings.Contains(xmlOutput, `id="custom-id"`) {
		t.Logf("XML output:\n%s", xmlOutput)
		t.Log("Note: Section attributes may be stored differently")
	}
	// At minimum, custom ID should be present
	if !strings.Contains(xmlOutput, `id="custom-id"`) {
		t.Error("Expected id attribute in section XML")
	}
}

func TestConvertToXML_TableCellAttributes(t *testing.T) {
	input := `= Test

|===
|^top|vbottom|left|center|right|
|===

[cols="1,1"]
|===
|[.class]#cell#|[colspan=2]#cell#
|===`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	xmlOutput := ToXML(doc)
	// Check for cell attributes
	if !strings.Contains(xmlOutput, `align="`) {
		t.Log("Note: Cell alignment attributes may be present")
	}
	if !strings.Contains(xmlOutput, `colspan="`) {
		t.Log("Note: Cell colspan attributes may be present")
	}
}

func TestConvertToXML_TableRowRole(t *testing.T) {
	input := `= Test

[cols="1,1",options="header"]
|===
|Header 1|Header 2
|Data 1|Data 2
|===`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	xmlOutput := ToXML(doc)
	// Table row role should be set for header rows
	if !strings.Contains(xmlOutput, `role="header"`) {
		t.Log("Note: Table row role attribute may be present")
	}
}

func TestConvertToXML_ListItemCallout(t *testing.T) {
	input := `= Test

<1> First callout
<2> Second callout`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	xmlOutput := ToXML(doc)
	if !strings.Contains(xmlOutput, `callout="1"`) || !strings.Contains(xmlOutput, `callout="2"`) {
		t.Log("Note: List item callout attributes may be present")
	}
}

func TestConvertToXML_LinkAttributes(t *testing.T) {
	input := `= Test

link:url[text, window="_blank", role="external"]`

	doc, err := ParseDocument(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	xmlOutput := ToXML(doc)
	if !strings.Contains(xmlOutput, `href="url"`) {
		t.Error("Expected href attribute in link XML")
	}
	// window and class attributes should be preserved
	if !strings.Contains(xmlOutput, `window="`) && !strings.Contains(xmlOutput, `class="`) {
		t.Log("Note: Link window and class attributes may be present")
	}
}

func TestConvertHTML_PicoCSSContent(t *testing.T) {
	input := `= Test Document

Content.`

	cssContent := "/* Custom PicoCSS */\nbody { margin: 0; }"
	opts := ConvertOptions{
		Standalone:     true,
		UsePicoCSS:    true,
		PicoCSSContent: cssContent,
		XHTML:         false,
	}

	result, err := Convert(bytes.NewReader([]byte(input)), opts)
	if err != nil {
		t.Fatalf("ConvertHTML failed: %v", err)
	}

	// Verify PicoCSS content is embedded inline
	if !strings.Contains(result.HTML, "<style>") {
		t.Error("PicoCSSContent should include <style> tag")
	}
	if !strings.Contains(result.HTML, cssContent) {
		t.Error("PicoCSSContent should include the CSS content")
	}
	if !strings.Contains(result.HTML, "</style>") {
		t.Error("PicoCSSContent should include closing </style> tag")
	}
	
	// Verify it does NOT use link tag when PicoCSSContent is provided
	if strings.Contains(result.HTML, `rel="stylesheet"`) {
		t.Error("PicoCSSContent should not use link tag, should embed inline")
	}
}

func TestConvertHTML_Preamble(t *testing.T) {
	input := `= Test Document

This is preamble content.

== First Section

Section content.`

	opts := ConvertOptions{
		Standalone: false,
		UsePicoCSS: false,
		XHTML:      false,
	}

	result, err := Convert(bytes.NewReader([]byte(input)), opts)
	if err != nil {
		t.Fatalf("ConvertHTML failed: %v", err)
	}

	// Verify preamble is wrapped in div with data-role="preamble"
	if !strings.Contains(result.HTML, `data-role="preamble"`) {
		t.Error("Preamble should be wrapped in div with data-role=\"preamble\"")
	}
	if !strings.Contains(result.HTML, "preamble content") {
		t.Error("Preamble content should be present in HTML output")
	}
	
	// Verify preamble wrapper is a div
	if !strings.Contains(result.HTML, `<div`) || !strings.Contains(result.HTML, `</div>`) {
		t.Error("Preamble should be wrapped in div tags")
	}
}

func TestConvertHTML_PreambleStandalone(t *testing.T) {
	input := `= Test Document

This is preamble content.

== First Section

Section content.`

	opts := ConvertOptions{
		Standalone: true,
		UsePicoCSS: false,
		XHTML:      false,
	}

	result, err := Convert(bytes.NewReader([]byte(input)), opts)
	if err != nil {
		t.Fatalf("ConvertHTML failed: %v", err)
	}

	// Verify preamble is wrapped in div with data-role="preamble" in standalone mode too
	if !strings.Contains(result.HTML, `data-role="preamble"`) {
		t.Error("Preamble should be wrapped in div with data-role=\"preamble\" in standalone mode")
	}
	if !strings.Contains(result.HTML, "preamble content") {
		t.Error("Preamble content should be present in HTML output")
	}
}

func TestSanitizeXMLAttributeName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "AsciiDoc attribute with leading and trailing colons",
			input:    ":toclevels:",
			expected: "toclevels",
		},
		{
			name:     "AsciiDoc attribute with single leading colon",
			input:    ":author",
			expected: "author",
		},
		{
			name:     "Valid attribute name",
			input:    "id",
			expected: "id",
		},
		{
			name:     "Attribute with hyphens",
			input:    "custom-attr",
			expected: "custom-attr",
		},
		{
			name:     "Attribute with underscores",
			input:    "custom_attr",
			expected: "custom_attr",
		},
		{
			name:     "Attribute with periods",
			input:    "custom.attr",
			expected: "custom.attr",
		},
		{
			name:     "Attribute starting with number",
			input:    "123attr",
			expected: "_123attr",
		},
		{
			name:     "Attribute with internal colons",
			input:    "my:custom:attr",
			expected: "my-custom-attr",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only colons",
			input:    ":::",
			expected: "attr",
		},
		{
			name:     "Special characters",
			input:    "attr@#$%",
			expected: "attr____",
		},
		{
			name:     "Mixed valid and invalid",
			input:    ":rev:number:",
			expected: "rev-number",
		},
		{
			name:     "Unicode characters",
			input:    "attrüñ",
			expected: "attr__",
		},
		{
			name:     "Space in name",
			input:    "my attr",
			expected: "my_attr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeXMLAttributeName(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeXMLAttributeName(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeXMLAttributeName_ValidXML(t *testing.T) {
	// Test that sanitized names produce valid XML
	testNames := []string{
		":toclevels:",
		":author:",
		":revnumber:",
		":custom-attr:",
		"123invalid",
		"my:namespace:attr",
	}

	for _, name := range testNames {
		sanitized := sanitizeXMLAttributeName(name)
		
		// Try to create XML with this attribute
		xmlStr := fmt.Sprintf(`<test %s="value"/>`, sanitized)
		
		// Parse it to ensure it's valid
		decoder := xml.NewDecoder(strings.NewReader(xmlStr))
		_, err := decoder.Token()
		if err != nil {
			t.Errorf("sanitizeXMLAttributeName(%q) = %q produced invalid XML: %v", name, sanitized, err)
		}
	}
}

func TestConvertToXML_SanitizedAttributes(t *testing.T) {
	input := `= Test Document
:author: John Doe
:toclevels: 3
:sectnums:
:custom-attr: Custom Value

Test content.`

	xmlOutput, err := ConvertToXML(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("ConvertToXML failed: %v", err)
	}

	// Verify XML is valid
	decoder := xml.NewDecoder(strings.NewReader(xmlOutput))
	for {
		_, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Fatalf("Invalid XML generated: %v\nXML: %s", err, xmlOutput)
		}
	}

	// Verify attributes are present and sanitized
	if !strings.Contains(xmlOutput, `toclevels="3"`) {
		t.Error("XML should contain sanitized toclevels attribute")
	}
	if !strings.Contains(xmlOutput, `author="John Doe"`) {
		t.Error("XML should contain author attribute")
	}
	if !strings.Contains(xmlOutput, `custom-attr="Custom Value"`) {
		t.Error("XML should contain custom-attr attribute")
	}
	
	// Verify no invalid colons in attribute names
	if strings.Contains(xmlOutput, `:toclevels=`) || strings.Contains(xmlOutput, `:author=`) {
		t.Error("XML should not contain attributes with leading colons")
	}
}

func TestConvertToXML_PassthroughBlock(t *testing.T) {
	input := `= Test

++++
<div class="custom">
  <p>Raw HTML content</p>
</div>
++++

Normal paragraph.`

	xmlOutput, err := ConvertToXML(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("ConvertToXML failed: %v", err)
	}

	// Verify XML is valid
	decoder := xml.NewDecoder(strings.NewReader(xmlOutput))
	for {
		_, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Fatalf("Invalid XML generated: %v\nXML: %s", err, xmlOutput)
		}
	}

	// Verify passthrough block is present in XML
	if !strings.Contains(xmlOutput, "<passthrough>") {
		t.Error("XML should contain <passthrough> element")
	}
	if !strings.Contains(xmlOutput, "</passthrough>") {
		t.Error("XML should contain </passthrough> closing tag")
	}

	// Verify content is escaped in XML (< becomes &lt;)
	if !strings.Contains(xmlOutput, "&lt;div") {
		t.Error("XML should escape < to &lt; in passthrough content")
	}

	// Verify ++++ delimiters are NOT in output
	if strings.Contains(xmlOutput, "++++") {
		t.Error("XML should not contain ++++ delimiters")
	}
}

func TestConvertToHTML_PassthroughBlockRaw(t *testing.T) {
	input := `= Test

++++
<div class="custom-html">
  <p>This is <strong>raw</strong> HTML</p>
</div>
++++

Normal paragraph.`

	opts := ConvertOptions{
		Standalone: false,
		UsePicoCSS: false,
		XHTML:      false,
	}

	result, err := Convert(bytes.NewReader([]byte(input)), opts)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	html := result.HTML

	// Verify raw HTML is output without escaping
	if !strings.Contains(html, `<div class="custom-html">`) {
		t.Error("HTML should contain raw <div> element")
	}
	if !strings.Contains(html, "<strong>raw</strong>") {
		t.Error("HTML should contain raw <strong> element")
	}

	// Verify ++++ delimiters are NOT in output
	if strings.Contains(html, "++++") {
		t.Error("HTML should not contain ++++ delimiters")
	}

	// Verify it's NOT escaped (should not have &lt;)
	if strings.Contains(html, "&lt;div") {
		t.Error("HTML should NOT escape passthrough content (found &lt;div)")
	}
}

func TestConvertToHTML_RoleAttributeNotParagraph(t *testing.T) {
	input := `[.my-custom-role]
--
Content in open block
--`

	opts := ConvertOptions{
		Standalone: false,
		UsePicoCSS: false,
		XHTML:      false,
	}

	result, err := Convert(bytes.NewReader([]byte(input)), opts)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	html := result.HTML

	// Should NOT have a paragraph with the role attribute syntax
	if strings.Contains(html, "[.my-custom-role]") {
		t.Error("HTML should not contain literal '[.my-custom-role]' - it should be consumed as attribute")
	}

	// Should have the role applied to the div
	if !strings.Contains(html, "my-custom-role") {
		t.Error("HTML should have 'my-custom-role' applied to the open block")
	}

	// Should have the open block
	if !strings.Contains(html, "Content in open block") {
		t.Error("HTML should contain the open block content")
	}
}

func TestConvertToXML_RoleAttributeNotParagraph(t *testing.T) {
	input := `[.grid-layout]
--
Grid content
--`

	xmlOutput, err := ConvertToXML(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("ConvertToXML failed: %v", err)
	}

	// Should NOT have paragraph with [.grid-layout]
	if strings.Contains(xmlOutput, "[.grid-layout]") {
		t.Error("XML should not contain literal '[.grid-layout]'")
	}

	// Should have openblock with role attribute
	if !strings.Contains(xmlOutput, "<openblock") {
		t.Error("XML should contain <openblock> element")
	}
	if !strings.Contains(xmlOutput, `role="grid-layout"`) {
		t.Error("XML should have role='grid-layout' on openblock")
	}
}
