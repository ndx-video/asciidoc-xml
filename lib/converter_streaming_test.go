package lib

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestConvertMarkdownToAsciiDocStreaming_Basic(t *testing.T) {
	input := `# Test Document

This is a paragraph.`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "= Test Document") {
		t.Errorf("Expected header conversion. Got:\n%s", result)
	}
	if !strings.Contains(result, "This is a paragraph") {
		t.Errorf("Expected paragraph content. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_Frontmatter(t *testing.T) {
	input := `---
title: Test Title
author: John Doe
categories:
  - tech
  - programming
---

# Content`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "= Test Title") {
		t.Errorf("Expected title from frontmatter. Got:\n%s", result)
	}
	if !strings.Contains(result, ":author: John Doe") {
		t.Errorf("Expected author attribute. Got:\n%s", result)
	}
	if !strings.Contains(result, ":categories: tech, programming") {
		t.Errorf("Expected categories array. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_FrontmatterNested(t *testing.T) {
	input := `---
title: Test
metadata:
  version: 1.0
  tags: [test, example]
---

# Content`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, ":metadata:") {
		t.Errorf("Expected nested metadata. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_Headers(t *testing.T) {
	input := `# H1
## H2
### H3
#### H4
##### H5
###### H6`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	expected := []string{"= H1", "== H2", "=== H3", "==== H4", "===== H5", "====== H6"}
	for _, exp := range expected {
		if !strings.Contains(result, exp) {
			t.Errorf("Expected %s in output. Got:\n%s", exp, result)
		}
	}
}

func TestConvertMarkdownToAsciiDocStreaming_CodeBlocks(t *testing.T) {
	input := "```go\nfunc main() {\n    println(\"hello\")\n}\n```"

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "[source,go]") {
		t.Errorf("Expected language specification. Got:\n%s", result)
	}
	if !strings.Contains(result, "----") {
		t.Errorf("Expected code block delimiters. Got:\n%s", result)
	}
	if !strings.Contains(result, "func main()") {
		t.Errorf("Expected code content. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_CodeBlockNoLanguage(t *testing.T) {
	input := "```\nplain text\n```"

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "----") {
		t.Errorf("Expected code block delimiters. Got:\n%s", result)
	}
	if !strings.Contains(result, "plain text") {
		t.Errorf("Expected code content. Got:\n%s", result)
	}
	// Should not have [source,lang] when no language specified
	if strings.Contains(result, "[source,") {
		t.Errorf("Should not have language spec without language. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_Tables(t *testing.T) {
	input := `| Col1 | Col2 |
|------|------|
| Val1 | Val2 |`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "[cols=") {
		t.Errorf("Expected column specification. Got:\n%s", result)
	}
	if !strings.Contains(result, "|===") {
		t.Errorf("Expected table delimiters. Got:\n%s", result)
	}
	if !strings.Contains(result, "|Col1") || !strings.Contains(result, "|Col2") {
		t.Errorf("Expected table headers. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_HorizontalRules(t *testing.T) {
	input := `Content

---

More content`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "'''") {
		t.Errorf("Expected horizontal rule. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_Lists(t *testing.T) {
	input := `- Item 1
- Item 2
* Item 3

1. Ordered 1
2. Ordered 2`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "* Item 1") || !strings.Contains(result, "* Item 2") {
		t.Errorf("Expected unordered list items. Got:\n%s", result)
	}
	if !strings.Contains(result, ". Ordered 1") || !strings.Contains(result, ". Ordered 2") {
		t.Errorf("Expected ordered list items. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_Blockquotes(t *testing.T) {
	input := `> This is a quote`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "[quote]") || !strings.Contains(result, "____") {
		t.Errorf("Expected blockquote conversion. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_Admonitions(t *testing.T) {
	input := `> **NOTE** This is a note
> [!WARNING] This is a warning`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "NOTE: This is a note") {
		t.Errorf("Expected NOTE admonition. Got:\n%s", result)
	}
	if !strings.Contains(result, "WARNING: This is a warning") {
		t.Errorf("Expected WARNING admonition. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_InlineFormatting(t *testing.T) {
	input := `This is **bold** and *italic* text`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "**bold**") {
		t.Errorf("Expected bold formatting. Got:\n%s", result)
	}
	if !strings.Contains(result, "_italic_") {
		t.Errorf("Expected italic formatting. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_Links(t *testing.T) {
	input := `[Link text](https://example.com)`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "https://example.com[Link text]") {
		t.Errorf("Expected link conversion. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_Images(t *testing.T) {
	input := `![Alt text](image.png)`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "image::image.png[Alt text]") {
		t.Errorf("Expected image conversion. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_LargeFile(t *testing.T) {
	// Create a large markdown file (simulate streaming with many lines)
	var inputBuilder strings.Builder
	inputBuilder.WriteString("# Large Document\n\n")
	for i := 0; i < 1000; i++ {
		inputBuilder.WriteString("This is paragraph ")
		inputBuilder.WriteString(strings.Repeat("content ", 50))
		inputBuilder.WriteString("\n\n")
	}

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(inputBuilder.String()), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "= Large Document") {
		t.Errorf("Expected header in large file output")
	}
	if len(result) < 1000 {
		t.Errorf("Expected substantial output from large file, got %d bytes", len(result))
	}
}

func TestConvertMarkdownToAsciiDocStreaming_EmptyInput(t *testing.T) {
	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(""), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming should handle empty input: %v", err)
	}
	// Empty input should produce empty or minimal output
}

func TestConvertMarkdownToAsciiDocStreaming_TableClosure(t *testing.T) {
	input := `| Col1 | Col2 |
|------|------|
| Val1 | Val2 |
Some text after table`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	// Table should be closed before "Some text"
	tableEnd := strings.Index(result, "|===")
	textStart := strings.Index(result, "Some text")
	if tableEnd == -1 || textStart == -1 || tableEnd >= textStart {
		t.Errorf("Table should be closed before following text. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_NestedLists(t *testing.T) {
	input := `- Item 1
  - Nested 1
  - Nested 2
- Item 2`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "* Item 1") {
		t.Errorf("Expected top-level item. Got:\n%s", result)
	}
	// Check that nested items have proper indentation
	if !strings.Contains(result, "  * Nested") {
		t.Errorf("Expected nested items with indentation. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_ComplexDocument(t *testing.T) {
	input := `---
title: Complex Test
author: Test Author
---

# Main Title

## Section 1

This is a paragraph with **bold** and *italic* text.

` + "```python" + `
def hello():
    print("world")
` + "```" + `

| Header 1 | Header 2 |
|----------|----------|
| Cell 1   | Cell 2   |

- List item 1
- List item 2

---

## Section 2

> **NOTE** This is important!`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	// Verify all major components
	checks := []string{
		"= Complex Test",
		":author: Test Author",
		"= Main Title",
		"== Section 1",
		"**bold**",
		"_italic_",
		"[source,python]",
		"[cols=",
		"|===",
		"* List item",
		"'''",
		"== Section 2",
		"NOTE: This is important",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("Expected to find '%s' in output. Got:\n%s", check, result)
		}
	}
}

func TestConvertMarkdownToAsciiDocStreaming_ErrorHandling(t *testing.T) {
	// Test with a writer that fails
	failingWriter := &failingWriter{}
	input := `# Test`

	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), failingWriter)
	if err == nil {
		t.Error("Expected error from failing writer")
	}
}

// failingWriter is a writer that always fails
type failingWriter struct{}

func (f *failingWriter) Write(p []byte) (n int, err error) {
	return 0, io.ErrClosedPipe
}

func TestConvertMarkdownToAsciiDocStreaming_FrontmatterWithEmptyValues(t *testing.T) {
	input := `---
title: Test
empty_field:
nested:
  key: value
---

# Content`

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "= Test") {
		t.Errorf("Expected title. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_MultipleCodeBlocks(t *testing.T) {
	input := "```go\ncode1\n```\n\nText\n\n```python\ncode2\n```"

	var output bytes.Buffer
	err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
	}

	result := output.String()
	// Should have both code blocks
	goIndex := strings.Index(result, "[source,go]")
	pythonIndex := strings.Index(result, "[source,python]")
	if goIndex == -1 || pythonIndex == -1 {
		t.Errorf("Expected both code blocks. Got:\n%s", result)
	}
	if goIndex >= pythonIndex {
		t.Errorf("Go block should come before Python block. Got:\n%s", result)
	}
}

func TestConvertMarkdownToAsciiDocStreaming_FrontmatterEndingVariations(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "three dashes",
			input:  "---\ntitle: Test\n---\n# Content",
			expect: "= Test",
		},
		{
			name:   "three dots",
			input:  "---\ntitle: Test\n...\n# Content",
			expect: "= Test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			err := ConvertMarkdownToAsciiDocStreaming(strings.NewReader(tt.input), &output)
			if err != nil {
				t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
			}

			result := output.String()
			if !strings.Contains(result, tt.expect) {
				t.Errorf("Expected %s. Got:\n%s", tt.expect, result)
			}
		})
	}
}

