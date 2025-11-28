package lib

import (
	"bytes"
	"strings"
	"testing"
)

func TestParse_BasicDocument(t *testing.T) {
	input := `= Test Document
:author: John Doe
:email: john@example.com

This is a simple paragraph.`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
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

	// Check email from document attributes
	emailText := doc.GetAttribute("email")
	if emailText != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", emailText)
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

func TestParse_EmptyDocument(t *testing.T) {
	doc, err := Parse(bytes.NewReader([]byte("")))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if doc == nil {
		t.Fatal("Document is nil")
	}

	doctype := doc.GetAttribute("doctype")
	if doctype != "article" {
		t.Errorf("Expected default doctype 'article', got '%s'", doctype)
	}
}

func TestParse_DocumentWithoutHeader(t *testing.T) {
	input := `This is a paragraph without a header.

This is another paragraph.`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
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

func TestParse_Sections(t *testing.T) {
	input := `= Document Title

== Section 1

Content of section 1.

=== Subsection 1.1

Content of subsection.

== Section 2

Content of section 2.`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
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

    // Assuming attributes still set (backward compat or implementation detail)
    // Or check children for heading
	titleAttr := firstSection.GetAttribute("title")
	if titleAttr == "" {
		t.Fatal("Expected section to have a title attribute")
	}
    
    // Verify it has title text or title attribute
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

func TestParse_CodeBlock(t *testing.T) {
	input := `= Test

[source,go]
----
package main

func main() {
    println("Hello")
}
----`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
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
			if !strings.Contains(content, "func main()") {
				t.Error("Expected code block to contain 'func main()'")
			}
			break
		}
	}

	if !foundCodeBlock {
		t.Error("Expected to find a code block (pre > code)")
	}
}

func TestParse_CodeBlockWithTitle(t *testing.T) {
	input := `= Test

.My Code Example
[source,go]
----
func main() {}
----`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	foundCodeBlock := false
	for _, child := range doc.Children {
		if child.Type == CodeBlock {
			foundCodeBlock = true
			title := child.GetAttribute("title")
			if title != "My Code Example" {
				t.Errorf("Expected title 'My Code Example', got '%s'", title)
			}
			break
		}
	}

	if !foundCodeBlock {
		t.Error("Expected to find a code block")
	}
}

func TestParse_LiteralBlock(t *testing.T) {
	input := `= Test

....
This is a literal block.
It preserves whitespace.
    And indentation.
....`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	foundLiteralBlock := false
	for _, child := range doc.Children {
		if child.Type == LiteralBlock {
			foundLiteralBlock = true
			content := getTextContent(child)
			if !strings.Contains(content, "preserves whitespace") {
				t.Error("Expected literal block to contain 'preserves whitespace'")
			}
			if !strings.Contains(content, "    And indentation") {
				t.Error("Expected literal block to preserve indentation")
			}
			break
		}
	}

	if !foundLiteralBlock {
		t.Error("Expected to find a literal block (pre.literal-block)")
	}
}

func TestParse_InlineFormatting(t *testing.T) {
	input := `= Test

This has *bold text* and _italic text_ and ` + "`monospace`" + ` text.`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
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
			content := getTextContent(child)
			if !strings.Contains(content, "bold text") {
				t.Error("Expected bold text to contain 'bold text'")
			}
		}
		if child.Type == Italic {
			foundItalic = true
			content := getTextContent(child)
			if !strings.Contains(content, "italic text") {
				t.Error("Expected italic text to contain 'italic text'")
			}
		}
		if child.Type == Monospace {
			foundMono = true
			content := getTextContent(child)
			if !strings.Contains(content, "monospace") {
				t.Error("Expected monospace text to contain 'monospace'")
			}
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

func TestParse_Lists(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		tag       string // ul, ol, dl
		itemCount int
	}{
		{
			name:      "unordered list",
			input:     "= Test\n\n* Item one\n* Item two\n* Item three",
			tag:       "ul",
			itemCount: 3,
		},
		{
			name:      "ordered list",
			input:     "= Test\n\n. First item\n. Second item\n. Third item",
			tag:       "ol",
			itemCount: 3,
		},
		{
			name:      "labeled list",
			input:     "= Test\n\nterm1:: definition1\nterm2:: definition2",
			tag:       "dl",
			itemCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(bytes.NewReader([]byte(tt.input)))
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			foundList := false
			for _, child := range doc.Children {
				if child.Type == List {
					style := child.GetAttribute("style")
					expectedStyle := "unordered"
					if tt.tag == "ol" {
						expectedStyle = "ordered"
					} else if tt.tag == "dl" {
						expectedStyle = "labeled"
					}
					if style == expectedStyle {
						foundList = true
						
						items := 0
						for _, item := range child.Children {
							if item.Type == ListItem {
								items++
							}
						}
						if items != tt.itemCount {
							t.Errorf("Expected %d items, got %d", tt.itemCount, items)
						}
						break
					}
				}
			}

			if !foundList {
				t.Errorf("Expected to find a list with style matching %s", tt.tag)
			}
		})
	}
}

func TestParse_Admonitions(t *testing.T) {
	input := `= Test

NOTE: This is a note.

WARNING: This is a warning.

TIP: This is a tip.

IMPORTANT: This is important.

CAUTION: This is a caution.`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	admonitions := 0
	types := make(map[string]bool)
	for _, child := range doc.Children {
		if child.Type == Admonition {
			admonitions++
			admType := child.GetAttribute("type")
			types[admType] = true
		}
	}

	if admonitions != 5 {
		t.Errorf("Expected 5 admonitions, got %d", admonitions)
	}

	expectedTypes := []string{"note", "warning", "tip", "important", "caution"}
	for _, expectedType := range expectedTypes {
		if !types[expectedType] {
			t.Errorf("Expected to find admonition type '%s'", expectedType)
		}
	}
}

func TestParse_Links(t *testing.T) {
	input := `= Test

Visit https://example.com[Example Website] for more info.

Or just https://example.org`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	foundLinks := 0
	for _, child := range doc.Children {
		if child.Type == Paragraph {
			for _, inlineChild := range child.Children {
				if inlineChild.Type == Link {
					foundLinks++
					href := inlineChild.GetAttribute("href")
					if href == "https://example.com" {
						content := getTextContent(inlineChild)
						if !strings.Contains(content, "Example Website") {
							t.Error("Expected link text to be 'Example Website'")
						}
					}
				}
			}
		}
	}

	if foundLinks < 1 {
		t.Error("Expected to find at least one link (a)")
	}
}

func TestParse_LinkMacros(t *testing.T) {
	input := `= Test

link:/docs?page=user-guide.adoc[User Guide] - Complete guide to using the web interface and features

link:/docs?page=user-testing.adoc[UAT Plan] - User Acceptance Testing scenarios and test cases

link:/docs?page=api.adoc[API Reference] - Complete REST API documentation for programmatic access

link:/simple/path - Simple link without text`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	foundLinks := 0
	linkMap := make(map[string]string) // href -> text

	for _, child := range doc.Children {
		if child.Type == Paragraph {
			for _, inlineChild := range child.Children {
				if inlineChild.Type == Link {
					foundLinks++
					href := inlineChild.GetAttribute("href")
					text := getTextContent(inlineChild)
					linkMap[href] = text
				}
			}
		}
	}

	if foundLinks < 4 {
		t.Errorf("Expected at least 4 links, found %d", foundLinks)
	}

	// Check specific links
	expectedLinks := map[string]string{
		"/docs?page=user-guide.adoc":   "User Guide",
		"/docs?page=user-testing.adoc": "UAT Plan",
		"/docs?page=api.adoc":           "API Reference",
		"/simple/path":                 "/simple/path", // When no text, href is used
	}

	for expectedHref, expectedText := range expectedLinks {
		if text, found := linkMap[expectedHref]; found {
			if text != expectedText {
				t.Errorf("Expected link '%s' to have text '%s', got '%s'", expectedHref, expectedText, text)
			}
		} else {
			t.Errorf("Expected to find link with href '%s'", expectedHref)
		}
	}
}

func TestParse_Image(t *testing.T) {
	input := `= Test

image::logo.png[Logo, 200, 100]

image:screenshot.png[Screenshot]`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	foundImages := 0
	for _, child := range doc.Children {
		if child.Type == BlockMacro && child.Name == "image" {
			foundImages++
			src := child.GetAttribute("src")
			if src == "logo.png" {
				alt := child.GetAttribute("alt")
				if alt != "Logo" {
					t.Errorf("Expected alt 'Logo', got '%s'", alt)
				}
			}
		}
	}

	if foundImages < 2 {
		t.Errorf("Expected 2 images, found %d", foundImages)
	}
}

func TestParse_Attributes(t *testing.T) {
	input := `= Test Document
:doctype: book
:revnumber: 1.0
:revdate: 2024-01-01
:revremark: Initial version
:custom-attr: Custom Value`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	doctype := doc.GetAttribute("doctype")
	if doctype != "book" {
		t.Errorf("Expected doctype 'book', got '%s'", doctype)
	}

	// Check revision attributes on document
	revNumber := doc.GetAttribute("revnumber")
	if revNumber != "1.0" {
		t.Errorf("Expected revision number '1.0', got '%s'", revNumber)
	}

	revDate := doc.GetAttribute("revdate")
	if revDate != "2024-01-01" {
		t.Errorf("Expected revision date '2024-01-01', got '%s'", revDate)
	}

	foundCustomAttr := false
	customAttr := doc.GetAttribute(":custom-attr")
	if customAttr == "Custom Value" {
		foundCustomAttr = true
	}

	if !foundCustomAttr {
		t.Error("Expected to find custom-attr attribute")
	}
}

func TestParse_ExampleBlock(t *testing.T) {
	input := `= Test

.Example Title
====
This is an example block.

With multiple paragraphs.
====`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	foundExample := false
	for _, child := range doc.Children {
		if child.Type == Example {
			foundExample = true
			title := child.GetAttribute("title")
			if title != "Example Title" {
				t.Errorf("Expected title 'Example Title', got '%s'", title)
			}
			if len(child.Children) == 0 {
				t.Error("Expected example to have content")
			}
			break
		}
	}

	if !foundExample {
		t.Error("Expected to find an example block")
	}
}

func TestParse_Sidebar(t *testing.T) {
	input := `= Test

****
This is a sidebar.
****`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	foundSidebar := false
	for _, child := range doc.Children {
		if child.Type == Sidebar {
			foundSidebar = true
			break
		}
	}

	if !foundSidebar {
		t.Error("Expected to find a sidebar (aside.sidebar)")
	}
}

func TestParse_Quote(t *testing.T) {
	input := `= Test

____
This is a quote.
____`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	foundQuote := false
	for _, child := range doc.Children {
		if child.Type == Quote {
			foundQuote = true
			break
		}
	}

	if !foundQuote {
		t.Error("Expected to find a quote block (blockquote.quote)")
	}
}

func TestParse_ThematicBreak(t *testing.T) {
	input := `= Test

First section.

'''

Second section.`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
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

func TestParse_PageBreak(t *testing.T) {
	input := `= Test

First page.

<<<

Second page.`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
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

func TestValidate_ValidDocument(t *testing.T) {
	input := `= Test Document

This is valid content.`

	err := Validate(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Errorf("Expected valid document, got error: %v", err)
	}
}

func TestParse_MultipleAuthors(t *testing.T) {
	input := `= Test
:author: First Author
:email: first@example.com
:author: Second Author
:email: second@example.com`

	doc, err := Parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check author attribute (can have multiple authors, but we store as single attribute for now)
	author := doc.GetAttribute("author")
	if author == "" {
		t.Error("Expected author attribute on document")
	}
	
	// The parser currently stores only one author attribute
	// Multiple authors would need special handling
	if author == "" {
		t.Error("Expected at least one author")
	}
}
