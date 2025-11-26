package converter

import (
	"bytes"
	"encoding/xml"
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

	if doc.Header == nil {
		t.Fatal("Header is nil")
	}

	if doc.Header.Title != "Test Document" {
		t.Errorf("Expected title 'Test Document', got '%s'", doc.Header.Title)
	}

	if len(doc.Header.Authors) != 1 {
		t.Fatalf("Expected 1 author, got %d", len(doc.Header.Authors))
	}

	if doc.Header.Authors[0].Name != "John Doe" {
		t.Errorf("Expected author name 'John Doe', got '%s'", doc.Header.Authors[0].Name)
	}

	if doc.Header.Authors[0].Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", doc.Header.Authors[0].Email)
	}

	if len(doc.Content.Items) == 0 {
		t.Fatal("Expected at least one content item")
	}

	para := doc.Content.Items[0].Paragraph
	if para == nil {
		t.Fatal("Expected first item to be a paragraph")
	}

	if len(para.Items) == 0 {
		t.Fatal("Expected paragraph to have items")
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
	for _, item := range doc.Content.Items {
		if item.Section != nil {
			sections++
		}
	}

	if sections < 2 {
		t.Errorf("Expected at least 2 sections, got %d", sections)
	}

	// Check first section
	firstSection := doc.Content.Items[0].Section
	if firstSection == nil {
		t.Fatal("Expected first item to be a section")
	}

	if firstSection.Level != 1 {
		t.Errorf("Expected section level 1, got %d", firstSection.Level)
	}

	if len(firstSection.Title.Items) == 0 {
		t.Fatal("Expected section to have a title")
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
	for _, item := range doc.Content.Items {
		if item.Paragraph != nil {
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

	para := doc.Content.Items[0].Paragraph
	if para == nil {
		t.Fatal("Expected a paragraph")
	}

	foundBold := false
	foundItalic := false
	foundMono := false

	for _, item := range para.Items {
		if item.Strong != nil {
			foundBold = true
		}
		if item.Emphasis != nil {
			foundItalic = true
		}
		if item.Monospace != nil {
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
	for _, item := range doc.Content.Items {
		if item.CodeBlock != nil {
			foundCodeBlock = true
			if item.CodeBlock.Language != "go" {
				t.Errorf("Expected language 'go', got '%s'", item.CodeBlock.Language)
			}
			if !strings.Contains(item.CodeBlock.Content, "package main") {
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
	for _, item := range doc.Content.Items {
		if item.List != nil {
			foundList = true
			if item.List.Style != "unordered" {
				t.Errorf("Expected unordered list, got '%s'", item.List.Style)
			}
			if len(item.List.Items) != 3 {
				t.Errorf("Expected 3 list items, got %d", len(item.List.Items))
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
	for _, item := range doc.Content.Items {
		if item.List != nil {
			foundList = true
			if item.List.Style != "ordered" {
				t.Errorf("Expected ordered list, got '%s'", item.List.Style)
			}
			if len(item.List.Items) != 3 {
				t.Errorf("Expected 3 list items, got %d", len(item.List.Items))
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
	for _, item := range doc.Content.Items {
		if item.List != nil {
			foundList = true
			if item.List.Style != "labeled" {
				t.Errorf("Expected labeled list, got '%s'", item.List.Style)
			}
			if len(item.List.Items) != 2 {
				t.Errorf("Expected 2 list items, got %d", len(item.List.Items))
			}
			break
		}
	}

	if !foundList {
		t.Error("Expected to find a labeled list")
	}
}

func TestConvert_Table(t *testing.T) {
	input := `= Test

|===
|Header 1 |Header 2
|Cell 1   |Cell 2
|===`

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	foundTable := false
	for _, item := range doc.Content.Items {
		if item.Table != nil {
			foundTable = true
			if item.Table.Header == nil {
				t.Error("Expected table to have a header")
			} else {
				if len(item.Table.Header.Cells) != 2 {
					t.Errorf("Expected 2 header cells, got %d", len(item.Table.Header.Cells))
				}
			}
			if len(item.Table.Rows) != 1 {
				t.Errorf("Expected 1 data row, got %d", len(item.Table.Rows))
			}
			break
		}
	}

	if !foundTable {
		t.Error("Expected to find a table")
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
	for _, item := range doc.Content.Items {
		if item.Admonition != nil {
			admonitions++
		}
	}

	if admonitions != 2 {
		t.Errorf("Expected 2 admonitions, got %d", admonitions)
	}

	// Check first admonition
	firstAdmonition := doc.Content.Items[0].Admonition
	if firstAdmonition == nil {
		t.Fatal("Expected first item to be an admonition")
	}

	if firstAdmonition.Type != "note" {
		t.Errorf("Expected admonition type 'note', got '%s'", firstAdmonition.Type)
	}
}

func TestConvert_Links(t *testing.T) {
	input := `= Test

Visit https://example.com[Example Website] for more info.`

	doc, err := Convert(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	para := doc.Content.Items[0].Paragraph
	if para == nil {
		t.Fatal("Expected a paragraph")
	}

	foundLink := false
	for _, item := range para.Items {
		if item.Link != nil {
			foundLink = true
			if item.Link.Href != "https://example.com" {
				t.Errorf("Expected href 'https://example.com', got '%s'", item.Link.Href)
			}
			break
		}
	}

	if !foundLink {
		t.Error("Expected to find a link")
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
	for _, item := range doc.Content.Items {
		if item.Image != nil {
			foundImage = true
			if item.Image.Src != "logo.png" {
				t.Errorf("Expected src 'logo.png', got '%s'", item.Image.Src)
			}
			if item.Image.Alt != "Logo" {
				t.Errorf("Expected alt 'Logo', got '%s'", item.Image.Alt)
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

	if doc.DocType != "book" {
		t.Errorf("Expected doctype 'book', got '%s'", doc.DocType)
	}

	if doc.Header.Revision == nil {
		t.Fatal("Expected revision to be set")
	}

	if doc.Header.Revision.Number != "1.0" {
		t.Errorf("Expected revision number '1.0', got '%s'", doc.Header.Revision.Number)
	}

	if doc.Header.Revision.Date != "2024-01-01" {
		t.Errorf("Expected revision date '2024-01-01', got '%s'", doc.Header.Revision.Date)
	}

	foundCustomAttr := false
	for _, attr := range doc.Header.Attributes {
		if attr.Name == "custom-attr" && attr.Value == "Custom Value" {
			foundCustomAttr = true
			break
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
	for _, item := range doc.Content.Items {
		if item.Example != nil {
			foundExample = true
			if item.Example.Title != "Example Title" {
				t.Errorf("Expected title 'Example Title', got '%s'", item.Example.Title)
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
	for _, item := range doc.Content.Items {
		if item.Sidebar != nil {
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
	for _, item := range doc.Content.Items {
		if item.Quote != nil {
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
	for _, item := range doc.Content.Items {
		if item.ThematicBreak != nil {
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
	for _, item := range doc.Content.Items {
		if item.PageBreak != nil {
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

	// Verify structure
	if doc.Header == nil {
		t.Fatal("Expected header")
	}

	if len(doc.Content.Items) == 0 {
		t.Fatal("Expected content items")
	}

	// Count different element types
	sections := 0
	codeBlocks := 0
	lists := 0
	tables := 0
	admonitions := 0

	for _, item := range doc.Content.Items {
		if item.Section != nil {
			sections++
		}
		if item.CodeBlock != nil {
			codeBlocks++
		}
		if item.List != nil {
			lists++
		}
		if item.Table != nil {
			tables++
		}
		if item.Admonition != nil {
			admonitions++
		}
	}

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

