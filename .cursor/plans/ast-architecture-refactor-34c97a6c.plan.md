<!-- 34c97a6c-07df-4653-8146-c7201a4a3b5b fc937fc3-a093-4843-9dbb-9105d670f81c -->
# AST Architecture Refactor (v0.3.0)

## Overview

Pivot to a strict Abstract Syntax Tree (AST) architecture that completely separates parsing from output generation. The parser will produce semantic node types with no HTML/XML bias, and converters will be responsible solely for serialization.

## Task 1: Define AST Structure (`lib/ast.go`)

Create new file with semantic node types:

**NodeType enum:**

- `Document` - Root document node
- `Section` - Document sections (with level attribute)
- `Paragraph` - Text paragraphs
- `BlockMacro` - Block-level macros (component::, image::, etc.)
- `InlineMacro` - Inline macros (if any)
- `Text` - Plain text content
- `List` - Lists (unordered, ordered, labeled)
- `ListItem` - List items
- `CodeBlock` - Code blocks with language
- `LiteralBlock` - Literal/preformatted blocks
- `Example` - Example blocks
- `Sidebar` - Sidebar blocks
- `Quote` - Quote blocks
- `Table` - Tables
- `TableRow` - Table rows
- `TableCell` - Table cells
- `Admonition` - Admonitions (NOTE, WARNING, etc.)
- `ThematicBreak` - Horizontal rules
- `PageBreak` - Page breaks
- `Bold` - Bold text
- `Italic` - Italic text
- `Monospace` - Monospace/code text
- `Link` - Links

**Node struct:**

```go
type Node struct {
    Type       NodeType
    Content    string            // For Text nodes
    Name       string            // For BlockMacro/InlineMacro (macro name)
    Attributes map[string]string
    Children   []*Node
}
```

**Helper constructors:**

- `NewDocumentNode() *Node`
- `NewSectionNode(level int) *Node`
- `NewParagraphNode() *Node`
- `NewBlockMacroNode(name string) *Node`
- `NewTextNode(content string) *Node`
- etc.

## Task 2: Replace DOM Structure

**Remove/Replace `lib/dom.go`:**

- Remove existing `NodeType` (ElementNode, TextNode, CommentNode)
- Remove existing `Node` struct with `Data` field
- Keep utility functions that are format-agnostic (Traverse, FindElementsByTag) but update them to work with new AST
- Move XML serialization logic to converter (see Task 3)

**Update all references:**

- Replace `NewElementNode(tagName)` calls with semantic constructors
- Replace `node.Data` access with `node.Type` checks
- Update `FindElementsByTag` to search by NodeType or Name attribute

## Task 3: Refactor Parser (`lib/adoc-parser.go`)

**Key changes:**

1. Parser creates AST nodes with semantic types only (no HTML tag names)
2. `component::name[]` creates `Node{Type: BlockMacro, Name: "component"}`
3. Regular text lines create `Node{Type: Paragraph}`
4. Sections create `Node{Type: Section}` with level attribute
5. All block types use semantic NodeTypes (CodeBlock, Table, Admonition, etc.)
6. Inline formatting creates separate node types (Bold, Italic, Monospace, Link)

**Parser methods to update:**

- `parse()` - Returns Document node
- `parseSection()` - Returns Section node
- `parseParagraph()` - Returns Paragraph node with inline children
- `parseComponentMacro()` - Returns BlockMacro node
- `parseImage()` - Returns BlockMacro node with Name="image"
- `parseList()` - Returns List node
- `parseCodeBlock()` - Returns CodeBlock node
- `parseInlineContent()` - Creates Bold, Italic, Monospace, Link nodes

**Critical constraint:** Parser must contain ZERO string literals like "<p>", "<paragraph>", etc. Only semantic logic.

## Task 4: Implement Converters (`lib/converter.go`)

**Create two separate converter functions:**

1. **`ToHTML(node *Node) string`**

   - `Document` -> (no wrapper, just children)
   - `Section` -> `<section>` with level attribute
   - `Paragraph` -> `<p>`
   - `BlockMacro` with Name="component" -> `<cms-component>`
   - `BlockMacro` with Name="image" -> `<img style="display:block">`
   - `List` -> `<ul>`, `<ol>`, or `<dl>` based on attributes
   - `CodeBlock` -> `<pre><code>`
   - `Bold` -> `<strong>`
   - `Italic` -> `<em>`
   - `Link` -> `<a href="...">`
   - etc.

2. **`ToXML(node *Node) string`**

   - `Document` -> `<document>`
   - `Section` -> `<section level="...">`
   - `Paragraph` -> `<paragraph>`
   - `BlockMacro` -> `<macro type="block" name="...">`
   - `List` -> `<list style="...">`
   - `CodeBlock` -> `<codeblock>`
   - `Bold` -> `<strong>`
   - `Italic` -> `<emphasis>`
   - `Link` -> `<link>`
   - etc.

**Update existing functions:**

- `ConvertToHTML()` - Calls `Parse()` then `ToHTML()`
- `ConvertToXML()` - Calls `Parse()` then `ToXML()`
- `Convert()` - Calls `Parse()` then `ToHTML()` with options

**Remove:**

- All HTML tag name logic from parser
- `node.ToXML()` method (move to converter)

## Task 5: Update Entry Points

**Update `lib/adoc-parser.go`:**

- `Parse(reader io.Reader) (*Node, error)` - Returns AST Document node

**Update `lib/converter.go`:**

- `ConvertToHTML(reader io.Reader, ...) (string, error)` - Parse then ToHTML
- `ConvertToXML(reader io.Reader) (string, error)` - Parse then ToXML
- `Convert(reader io.Reader, opts ConvertOptions) (Result, error)` - Parse then ToHTML

## Task 6: Update Tests

**Update test files:**

- `lib/adoc-parser_test.go` - Check for AST node types instead of HTML tags
- `lib/converter_test.go` - Verify HTML/XML output from AST
- `lib/dom_test.go` - Update or remove if no longer relevant

## Files to Modify

1. **Create:** `lib/ast.go` - New AST structure
2. **Replace:** `lib/dom.go` - Remove DOM structure, keep only AST-agnostic utilities
3. **Refactor:** `lib/adoc-parser.go` - Pure AST generation
4. **Refactor:** `lib/converter.go` - Separate HTML/XML serialization
5. **Update:** `lib/adoc-parser_test.go` - Test AST structure
6. **Update:** `lib/converter_test.go` - Test converter output
7. **Update:** `lib/version.go` - Bump to v0.3.0

## Constraints

- Zero dependencies (maintain current)
- Parser contains NO HTML/XML string literals
- AST is completely format-agnostic
- Converters handle all format-specific logic

### To-dos

- [ ] Create lib/ast.go with NodeType enum and Node struct with semantic types
- [ ] Replace lib/dom.go - remove DOM structure, keep only AST-agnostic utilities
- [ ] Refactor lib/adoc-parser.go to create AST nodes with semantic types only (no HTML tags)
- [ ] Implement ToHTML() function in lib/converter.go for HTML serialization
- [ ] Implement ToXML() function in lib/converter.go for XML serialization
- [ ] Update Parse(), ConvertToHTML(), ConvertToXML(), Convert() entry points
- [ ] Update all test files to work with new AST structure
- [ ] Update lib/version.go to v0.3.0