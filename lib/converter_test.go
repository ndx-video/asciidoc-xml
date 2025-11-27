package lib

import (
	"bytes"
	"encoding/xml"
	"os"
	"strings"
	"testing"
)

func TestConvert_BasicDocument(t *testing.T) {
	input := `= Test Document
:author: John Doe
:email: john@example.com

This is a simple paragraph.`

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if doc == nil {
		t.Fatal("Document is nil")
	}

	if doc.Data != "asciidoc" {
		t.Errorf("Expected root element 'asciidoc', got '%s'", doc.Data)
	}

	// Find header
	var headerNode *Node
	for _, child := range doc.Children {
		if child.Data == "header" {
			headerNode = child
			break
		}
	}

	if headerNode == nil {
		t.Fatal("Header is nil")
	}

	titleNode := findChild(headerNode, "title")
	if titleNode == nil {
		t.Fatal("Title is nil")
	}

	titleText := getTextContent(titleNode)
	if titleText != "Test Document" {
		t.Errorf("Expected title 'Test Document', got '%s'", titleText)
	}

	// Find author
	var authorNode *Node
	for _, child := range headerNode.Children {
		if child.Data == "author" {
			authorNode = child
			break
		}
	}

	if authorNode == nil {
		t.Fatal("Author is nil")
	}

	nameNode := findChild(authorNode, "name")
	if nameNode == nil {
		t.Fatal("Author name is nil")
	}

	nameText := getTextContent(nameNode)
	if nameText != "John Doe" {
		t.Errorf("Expected author name 'John Doe', got '%s'", nameText)
	}

	// Find paragraph in content
	var paraNode *Node
	for _, child := range doc.Children {
		if child.Data == "paragraph" {
			paraNode = child
			break
		}
	}

	if paraNode == nil {
		t.Fatal("Expected at least one paragraph")
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

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	sections := 0
	for _, child := range doc.Children {
		if child.Data == "section" {
			sections++
		}
	}

	if sections < 2 {
		t.Errorf("Expected at least 2 sections, got %d", sections)
	}

	// Check first section
	var firstSection *Node
	for _, child := range doc.Children {
		if child.Data == "section" {
			firstSection = child
			break
		}
	}

	if firstSection == nil {
		t.Fatal("Expected first item to be a section")
	}

	level := firstSection.GetAttribute("level")
	if level != "1" {
		t.Errorf("Expected section level 1, got %s", level)
	}

	titleAttr := firstSection.GetAttribute("title")
	if titleAttr == "" {
		t.Fatal("Expected section to have a title attribute")
	}
}

func TestConvert_Paragraphs(t *testing.T) {
	input := `= Test

This is paragraph one.

This is paragraph two with *bold* and _italic_ text.`

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	paragraphs := 0
	for _, child := range doc.Children {
		if child.Data == "paragraph" {
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

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	var paraNode *Node
	for _, child := range doc.Children {
		if child.Data == "paragraph" {
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
		if child.Data == "strong" {
			foundBold = true
		}
		if child.Data == "emphasis" {
			foundItalic = true
		}
		if child.Data == "monospace" {
			foundMono = true
		}
	}

	if !foundBold {
		t.Error("Expected to find bold text")
	}
	if !foundItalic {
		t.Error("Expected to find italic text")
	}
	if !foundMono {
		t.Error("Expected to find monospace text")
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

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundCodeBlock := false
	for _, child := range doc.Children {
		if child.Data == "codeblock" {
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

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundList := false
	for _, child := range doc.Children {
		if child.Data == "list" {
			foundList = true
			style := child.GetAttribute("style")
			if style != "unordered" {
				t.Errorf("Expected unordered list, got '%s'", style)
			}
			items := 0
			for _, item := range child.Children {
				if item.Data == "item" {
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
		t.Error("Expected to find a list")
	}
}

func TestConvert_OrderedList(t *testing.T) {
	input := `= Test

. First item
. Second item
. Third item`

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundList := false
	for _, child := range doc.Children {
		if child.Data == "list" {
			foundList = true
			style := child.GetAttribute("style")
			if style != "ordered" {
				t.Errorf("Expected ordered list, got '%s'", style)
			}
			items := 0
			for _, item := range child.Children {
				if item.Data == "item" {
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
		t.Error("Expected to find a list")
	}
}

func TestConvert_LabeledList(t *testing.T) {
	input := `= Test

term1:: definition1
term2:: definition2`

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundList := false
	for _, child := range doc.Children {
		if child.Data == "list" {
			foundList = true
			style := child.GetAttribute("style")
			if style != "labeled" {
				t.Errorf("Expected labeled list, got '%s'", style)
			}
			items := 0
			for _, item := range child.Children {
				if item.Data == "item" {
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
		t.Error("Expected to find a labeled list")
	}
}

func TestConvert_Admonition(t *testing.T) {
	input := `= Test

NOTE: This is a note.

WARNING: This is a warning.`

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	admonitions := 0
	for _, child := range doc.Children {
		if child.Data == "admonition" {
			admonitions++
		}
	}

	if admonitions != 2 {
		t.Errorf("Expected 2 admonitions, got %d", admonitions)
	}

	// Check first admonition
	var firstAdmonition *Node
	for _, child := range doc.Children {
		if child.Data == "admonition" {
			firstAdmonition = child
			break
		}
	}

	if firstAdmonition == nil {
		t.Fatal("Expected first item to be an admonition")
	}

	admType := firstAdmonition.GetAttribute("type")
	if admType != "note" {
		t.Errorf("Expected admonition type 'note', got '%s'", admType)
	}
}

func TestConvert_Links(t *testing.T) {
	input := `= Test

Visit https://example.com[Example Website] for more info.`

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	var paraNode *Node
	for _, child := range doc.Children {
		if child.Data == "paragraph" {
			paraNode = child
			break
		}
	}

	if paraNode == nil {
		t.Fatal("Expected a paragraph")
	}

	foundLink := false
	for _, child := range paraNode.Children {
		if child.Data == "link" {
			foundLink = true
			href := child.GetAttribute("href")
			if href != "https://example.com" {
				t.Errorf("Expected href 'https://example.com', got '%s'", href)
			}
			break
		}
	}

	if !foundLink {
		t.Error("Expected to find a link")
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

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundImage := false
	for _, child := range doc.Children {
		if child.Data == "image" {
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
		t.Error("Expected to find an image")
	}
}

func TestConvert_Attributes(t *testing.T) {
	input := `= Test Document
:doctype: book
:revnumber: 1.0
:revdate: 2024-01-01
:revremark: Initial version
:custom-attr: Custom Value`

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	doctype := doc.GetAttribute("doctype")
	if doctype != "book" {
		t.Errorf("Expected doctype 'book', got '%s'", doctype)
	}

	var headerNode *Node
	for _, child := range doc.Children {
		if child.Data == "header" {
			headerNode = child
			break
		}
	}

	if headerNode == nil {
		t.Fatal("Expected header")
	}

	revNode := findChild(headerNode, "revision")
	if revNode == nil {
		t.Fatal("Expected revision to be set")
	}

	foundCustomAttr := false
	for _, attr := range headerNode.Children {
		if attr.Data == "attribute" && attr.GetAttribute("name") == "custom-attr" {
			if attr.GetAttribute("value") == "Custom Value" {
				foundCustomAttr = true
				break
			}
		}
	}

	if !foundCustomAttr {
		t.Error("Expected to find custom-attr attribute")
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
	if !strings.Contains(xmlOutput, "asciidoc") {
		t.Error("XML should contain 'asciidoc' element")
	}
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

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundExample := false
	for _, child := range doc.Children {
		if child.Data == "example" {
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

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundSidebar := false
	for _, child := range doc.Children {
		if child.Data == "sidebar" {
			foundSidebar = true
			break
		}
	}

	if !foundSidebar {
		t.Error("Expected to find a sidebar")
	}
}

func TestConvert_Quote(t *testing.T) {
	input := `= Test

____
This is a quote.
____`

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundQuote := false
	for _, child := range doc.Children {
		if child.Data == "quote" {
			foundQuote = true
			break
		}
	}

	if !foundQuote {
		t.Error("Expected to find a quote block")
	}
}

func TestConvert_ThematicBreak(t *testing.T) {
	input := `= Test

First section.

'''

Second section.`

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundBreak := false
	for _, child := range doc.Children {
		if child.Data == "thematicbreak" {
			foundBreak = true
			break
		}
	}

	if !foundBreak {
		t.Error("Expected to find a thematic break")
	}
}

func TestConvert_PageBreak(t *testing.T) {
	input := `= Test

First page.

<<<

Second page.`

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundBreak := false
	for _, child := range doc.Children {
		if child.Data == "pagebreak" {
			foundBreak = true
			break
		}
	}

	if !foundBreak {
		t.Error("Expected to find a page break")
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

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Count different element types recursively
	var countElements func(*Node, *int, *int, *int, *int, *int)
	countElements = func(node *Node, sections, codeBlocks, lists, tables, admonitions *int) {
		if node.Data == "section" {
			*sections++
		}
		if node.Data == "codeblock" {
			*codeBlocks++
		}
		if node.Data == "list" {
			*lists++
		}
		if node.Data == "table" {
			*tables++
		}
		if node.Data == "admonition" {
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

	// Check bold conversion
	if !strings.Contains(output, "*bold*") {
		t.Error("Expected '*bold*' in output")
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

	// Check image conversion
	if !strings.Contains(output, "image::image.png[") {
		t.Error("Expected 'image::image.png[' in output")
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
