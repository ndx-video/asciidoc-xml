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

	titleNode := findChild(headerNode, "h1")
	if titleNode == nil {
		t.Fatal("Title (h1) is nil")
	}

	titleText := getTextContent(titleNode)
	if titleText != "Test Document" {
		t.Errorf("Expected title 'Test Document', got '%s'", titleText)
	}

	// Find author
	var authorNode *Node
	for _, child := range headerNode.Children {
		if child.Data == "address" {
			authorNode = child
			break
		}
	}

	if authorNode == nil {
		t.Fatal("Author (address) is nil")
	}

	nameNode := findChild(authorNode, "span")
	if nameNode == nil || nameNode.GetAttribute("class") != "author-name" {
		t.Fatal("Author name (span) is nil or wrong class")
	}

	nameText := getTextContent(nameNode)
	if nameText != "John Doe" {
		t.Errorf("Expected author name 'John Doe', got '%s'", nameText)
	}

	emailNode := findChild(authorNode, "a")
	if emailNode == nil || emailNode.GetAttribute("class") != "email" {
		t.Fatal("Author email (a) is nil or wrong class")
	}

	emailText := getTextContent(emailNode)
	if emailText != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", emailText)
	}

	// Find paragraph in content
	var paraNode *Node
	for _, child := range doc.Children {
		if child.Data == "p" {
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
		if child.Data == "p" {
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

    // Assuming attributes still set (backward compat or implementation detail)
    // Or check children for heading
	titleAttr := firstSection.GetAttribute("title")
	if titleAttr == "" {
		t.Fatal("Expected section to have a title attribute")
	}
    
    // Verify it has H2 child
    foundH2 := false
    for _, child := range firstSection.Children {
        if child.Data == "h2" {
            foundH2 = true
            break
        }
    }
    if !foundH2 {
        t.Error("Expected section to have h2 child")
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
		if child.Data == "pre" {
            // Check for code child
            var codeChild *Node
            for _, c := range child.Children {
                if c.Data == "code" {
                    codeChild = c
                    break
                }
            }
            if codeChild == nil {
                continue // Might be literal block
            }
            
			foundCodeBlock = true
            // Language can be on pre or code class
			lang := codeChild.GetAttribute("class")
			if !strings.Contains(lang, "language-go") {
				t.Errorf("Expected class to contain 'language-go', got '%s'", lang)
			}
            
			content := getTextContent(codeChild)
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
		if child.Data == "pre" {
            // Check if it has code child to be sure
            hasCode := false
            for _, c := range child.Children {
                if c.Data == "code" { hasCode = true; break }
            }
            if !hasCode { continue }
            
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
		if child.Data == "pre" && child.GetAttribute("class") == "literal-block" {
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
		if child.Data == "p" {
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
			content := getTextContent(child)
			if !strings.Contains(content, "bold text") {
				t.Error("Expected bold text to contain 'bold text'")
			}
		}
		if child.Data == "em" {
			foundItalic = true
			content := getTextContent(child)
			if !strings.Contains(content, "italic text") {
				t.Error("Expected italic text to contain 'italic text'")
			}
		}
		if child.Data == "code" {
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
				if child.Data == tt.tag {
					foundList = true
					
					items := 0
					for _, item := range child.Children {
                        if tt.tag == "dl" {
                            // For dl, we used div wrappers in parseListItem
                            if item.Data == "div" {
                                items++
                            }
                        } else {
						    if item.Data == "li" {
							    items++
						    }
                        }
					}
					if items != tt.itemCount {
						t.Errorf("Expected %d items, got %d", tt.itemCount, items)
					}
					break
				}
			}

			if !foundList {
				t.Errorf("Expected to find a list with tag %s", tt.tag)
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
		if child.Data == "div" && strings.Contains(child.GetAttribute("class"), "admonition") {
			admonitions++
			class := child.GetAttribute("class")
            // Extract type from class
            parts := strings.Split(class, " ")
            for _, p := range parts {
                if strings.HasPrefix(p, "admonition-") {
                    types[strings.TrimPrefix(p, "admonition-")] = true
                }
            }
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
		if child.Data == "p" {
			for _, inlineChild := range child.Children {
				if inlineChild.Data == "a" {
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
		if child.Data == "p" {
			for _, inlineChild := range child.Children {
				if inlineChild.Data == "a" {
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
		if child.Data == "img" {
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

    // Revision is div.revision
	var revNode *Node
    for _, c := range headerNode.Children {
        if c.Data == "div" && c.GetAttribute("class") == "revision" {
            revNode = c
            break
        }
    }
    
	if revNode == nil {
		t.Fatal("Expected revision to be set (div.revision)")
	}

	numberNode := findChild(revNode, "span") // class=revnumber
	if numberNode == nil || numberNode.GetAttribute("class") != "revnumber" {
		t.Fatal("Expected revision number (span.revnumber)")
	}
	numberText := getTextContent(numberNode)
	if numberText != "1.0" {
		t.Errorf("Expected revision number '1.0', got '%s'", numberText)
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
		if child.Data == "div" && child.GetAttribute("class") == "example" {
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
		if child.Data == "aside" && child.GetAttribute("class") == "sidebar" {
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
		if child.Data == "blockquote" && child.GetAttribute("class") == "quote" {
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
		if child.Data == "hr" {
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
		if child.Data == "div" && child.GetAttribute("class") == "page-break" {
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

	authors := 0
	for _, child := range headerNode.Children {
		if child.Data == "address" && child.GetAttribute("class") == "author" {
			authors++
		}
	}

	if authors != 2 {
		t.Errorf("Expected 2 authors, got %d", authors)
	}
}
