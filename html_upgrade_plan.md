# HTML Converter Enhancement Plan

## Overview

Upgrade the `toHTML` function in `lib/converter.go` to support all AST node types and features that are currently parsed but not yet converted to HTML. This includes inline formatting, macros, cross-references, footnotes, and block enhancements.

## Current State Audit

### Ensure the following is Implemented in HTML Conversion

- **Document**: Outputs children directly
- **Section**: Converts to `<h1>`-`<h6>` with id and data attributes (class reserved for user)
- **Paragraph**: Converts to `<p>` with id and data attributes (class reserved for user)
- **BlockMacro**: Handles `component`, `image`, and generic macros
- **Text**: Escaped and output
- **List**: Converts to `<ul>`, `<ol>`, or `<dl>` based on style
- **CodeBlock**: Converts to `<pre><code>` with language support and mermaid role
- **LiteralBlock**: Converts to `<pre>` with data-role="literal-block"
- **Example**: Converts to `<div data-role="example">`
- **Sidebar**: Converts to `<aside data-role="sidebar">`
- **Quote**: Converts to `<blockquote>` with attribution/citation
- **Table**: Basic table conversion with `<thead>` and `<tbody>`
- **Admonition**: Converts to `<div data-role="admonition" data-asciidoc-variant="note">` (variant like "note", "warning", "tip", etc.)
- **ThematicBreak**: Converts to `<hr>`
- **PageBreak**: Converts to `<div data-role="page-break">`
- **Bold**: Converts to `<strong>`
- **Italic**: Converts to `<em>`
- **Monospace**: Converts to `<code>`
- **Link**: Converts to `<a>` with href, title, target (class reserved for user)
- **Passthrough**: Outputs raw HTML (for CMS)

### This is what is Missing from HTML Conversion

#### Inline Formatting

- **Superscript** (`^text^`): AST exists, XML exists, HTML missing
- **Subscript** (`~text~`): AST exists, XML exists, HTML missing
- **Highlight** (`#text#`): AST exists, XML exists, HTML missing

#### Block Elements

- **VerseBlock** (`[verse]____...____`): AST exists, XML exists, HTML missing
- **OpenBlock** (`--`): AST exists, XML exists, HTML missing
- **Preamble**: Not explicitly handled in HTML (first paragraph with role="preamble")

#### Inline Macros

- **Anchor** (`[[id]]`, `[#id]`): Parsed as InlineMacro/BlockMacro, HTML missing
- **Footnote** (`footnote:[text]`): Parsed as InlineMacro, HTML missing
- **FootnoteRef** (`footnoteref:ref[]`): Parsed as InlineMacro, HTML missing
- **Kbd** (`kbd:[keys]`): Parsed as InlineMacro, HTML missing
- **Btn** (`btn:[label]`): Parsed as InlineMacro, HTML missing
- **Menu** (`menu:File[New]`): Parsed as InlineMacro, HTML missing
- **Generic inline macros**: Not handled in `toHTML` switch statement

#### Block Macros

- **Include** (`include::file[]`): Parsed but not converted
- **TOC** (`toc::[]`): Parsed but not converted
- **Video** (`video::url[]`): Parsed but not converted
- **Audio** (`audio::url[]`): Parsed but not converted

#### Cross-References

- **Xref** (`<<anchor-id>>`, `xref:anchor-id[]`): Parsed as Link with target attribute, needs HTML conversion

#### Table Enhancements

- **Cell alignment**: `align` attribute exists but not applied to `<td>`/`<th>`
- **Cell colspan/rowspan**: Attributes exist but not applied
- **Row role**: `role="header"`/`role="footer"` partially handled but needs enhancement
- **Table role/id**: Attributes exist but may not be fully applied

#### List Enhancements

- **Callout lists**: `callout` attribute exists but not rendered
- **List continuations**: Parsed but visual representation missing
- **List item attributes**: `id`, `role`, `term` exist but may not be fully applied

#### Section Enhancements

- **Appendix**: `appendix` attribute exists but not used
- **Discrete**: `discrete` attribute exists but not used

## Implementation Plan

### Critical Constraint: Clean Class Policy

**Important**: We have a strict **"Clean Class" Policy**. The converter must NOT pollute the `class` attribute with library-specific styling names. The `class` attribute is reserved for user-defined classes only.

**Strategy:**

- Use `data-role` or `data-asciidoc-*` attributes for structural semantics
- Use semantic HTML elements (`<div>`, `<span>`, `<section>`, etc.) without classes
- Preserve all AsciiDoc attributes as `data-*` attributes
- Use `id` attributes when provided
- Use `role` attributes for ARIA semantics when appropriate
- Map non-standard AsciiDoc attributes to `data-asciidoc-{name}="value"` format

### Phase 1: Inline Formatting (Direct HTML Equivalents)

**Superscript** (`^text^`)

- **HTML**: `<sup data-asciidoc="superscript">text</sup>`
- **Implementation**: Add case for `Superscript` node type
- **XHTML**: Same as HTML (well-formed)

**Subscript** (`~text~`)

- **HTML**: `<sub data-asciidoc="subscript">text</sub>`
- **Implementation**: Add case for `Subscript` node type
- **XHTML**: Same as HTML (well-formed)

**Highlight** (`#text#`)

- **HTML**: `<mark data-asciidoc="highlight">text</mark>` (HTML5 semantic element)
- **Implementation**: Add case for `Highlight` node type, use `<mark>` for HTML5
- **XHTML**: `<mark>` is valid in XHTML5

### Phase 2: Block Elements

**Note on Comments**: For clarity in debug mode, prepend HTML comments like `<!-- verse block -->` before block elements when helpful.

**VerseBlock** (`[verse]____...____`)

- **Challenge**: Preserve line breaks and whitespace
- **HTML Solution**: 
- Use `<div data-role="verse">` with `<p>` for each line or `<br>` for line breaks
- Attribution: Use `<footer><cite>` similar to Quote
- **Implementation**: 
- Add case for `VerseBlock`
- Preserve line breaks (convert `\n` to `<br>` or wrap lines in `<p>`)
- Output attribution if present
- **Attributes**: `title`, `attribution`, `id`, `role` → map to `data-asciidoc-*` if not standard HTML

**OpenBlock** (`--`)

- **Challenge**: Generic container with no semantic HTML equivalent
- **HTML Solution**: `<div data-role="open-block">` with all attributes preserved
- **Implementation**: Add case for `OpenBlock`, output as `<div>` with data-role and data attributes
- **Attributes**: All attributes preserved as HTML attributes (if standard) or `data-asciidoc-*` attributes

**Preamble**

- **Challenge**: First paragraph with `role="preamble"` should be wrapped
- **HTML Solution**: Wrap in `<div data-role="preamble">` or use `<section data-role="preamble">`
- **Implementation**: 
- Check first child of Document for Paragraph with `role="preamble"` or `data-role="preamble"`
- Wrap in semantic container
- Similar to XML preamble handling

### Phase 3: Inline Macros

**Anchor** (`[[id]]`, `[#id]`)

- **HTML**: `<a id="anchor-id"></a>` (empty anchor for linking)
- **Implementation**: 
- Handle `InlineMacro` with `name="anchor"` in inline content
- Handle `BlockMacro` with `name="anchor"` in block context
- Output as `<a id="..."></a>` or `<span id="..."></span>` (span for inline anchors)
- **Note**: Block anchors can be `<a id="..."></a>` or `<div id="..."></div>`

**Footnote** (`footnote:[text]`)

- **Challenge**: Footnotes need to be collected and rendered at document end
- **HTML Solution**: 
- Inline: `<sup><a href="#fn1" id="fnref1" data-role="footnote-ref">1</a></sup>`
- Footer: `<div data-role="footnotes"><ol><li id="fn1">text <a href="#fnref1">↩</a></li></ol></div>`
- **Implementation**: 
- Collect footnotes during traversal
- Render inline reference with link
- Render footnote list at document end
- Handle `ref` attribute for custom IDs
- **Attributes**: `ref` (footnote reference ID) → `data-asciidoc-ref`

**FootnoteRef** (`footnoteref:ref[]`)

- **HTML**: `<sup><a href="#fn1" id="fnref1" data-role="footnote-ref">1</a></sup>`
- **Implementation**: Similar to footnote but references existing footnote
- **Attributes**: `ref` (footnote reference ID) → `data-asciidoc-ref`

**Kbd** (`kbd:[Ctrl+C]`)

- **HTML**: `<kbd data-role="keyboard">Ctrl+C</kbd>` (HTML5 semantic element)
- **Implementation**: Handle `InlineMacro` with `name="kbd"` in inline content
- **XHTML**: `<kbd>` is valid in XHTML5

**Btn** (`btn:[Save]`)

- **Challenge**: No semantic HTML element for buttons in text
- **HTML Solution**: `<span data-role="button">Save</span>` or `<button type="button" data-role="button-inline">Save</button>`
- **Implementation**: Use `<span>` with `data-role="button"` for inline context
- **Note**: Could use `<button>` but may have styling implications

**Menu** (`menu:File[New]`)

- **Challenge**: No semantic HTML for menu paths
- **HTML Solution**: `<span data-role="menu-path"><span data-role="menu">File</span> → <span data-role="menu">New</span></span>`
- **Implementation**: Parse target (e.g., "File[New]") and render with separators
- **Attributes**: Parse menu hierarchy from target → `data-asciidoc-menu-path`

**Generic Inline Macros**

- **HTML Solution**: `<span data-role="macro" data-asciidoc-macro="{name}">content</span>`
- **Implementation**: Handle all `InlineMacro` nodes not specifically handled
- **Attributes**: Preserve all attributes as `data-asciidoc-*` attributes

### Phase 4: Block Macros

**Include** (`include::file[]`)

- **Challenge**: File inclusion happens at parse time or needs placeholder
- **HTML Solution**: 
- If included: Render included content
- If placeholder: `<div data-role="include" data-asciidoc-file="path">[Include: path]</div>`
- **Implementation**: 
- Check if macro has children (included content)
- If yes, render children
- If no, render placeholder
- **Attributes**: `src` or `target` (file path) → `data-asciidoc-file` or `data-asciidoc-target`

**TOC** (`toc::[]`)

- **Challenge**: Table of contents generation
- **HTML Solution**: 
- Generate `<nav data-role="toc"><ul><li><a href="#section-id">Section Title</a></li></ul></nav>`
- Requires document traversal to collect sections
- **Implementation**: 
- Traverse document to collect sections with IDs
- Generate nested list structure
- Handle `levels` attribute for depth
- **Attributes**: `levels` (max depth) → `data-asciidoc-levels`, `title` → standard HTML `title`

**Video** (`video::url[]`)

- **HTML**: `<video src="url" controls></video>`
- **Implementation**: Handle `BlockMacro` with `name="video"`
- **Attributes**: `src`, `width`, `height`, `autoplay`, `loop`, `controls`, `poster`

**Audio** (`audio::url[]`)

- **HTML**: `<audio src="url" controls></audio>`
- **Implementation**: Handle `BlockMacro` with `name="audio"`
- **Attributes**: `src`, `autoplay`, `loop`, `controls`

### Phase 5: Cross-References

**Xref** (`<<anchor-id>>`, `xref:anchor-id[]`)

- **Current**: Parsed as `Link` with `target` attribute pointing to anchor
- **HTML Solution**: `<a href="#anchor-id">text</a>`
- **Implementation**: 
- Check if `Link` has `target` attribute (cross-reference)
- Use `href="#{target}"` instead of `href="{href}"`
- Use link text or anchor ID as display text
- **Attributes**: `target` (anchor ID), `href` (may be empty for xref)

### Phase 6: Table Enhancements

**Cell Alignment**

- **HTML**: Use `data-align="{align}"` attribute on `<td>`/`<th>` (pure data attribute approach)
- **Implementation**: Apply `data-align` attribute to table cells
- **Note**: Keep it pure - use `data-align` instead of inline styles or classes

**Cell Spanning**

- **HTML**: `colspan` and `rowspan` attributes on `<td>`/`<th>`
- **Implementation**: Read `colspan` and `rowspan` attributes and apply to cells

**Row Roles**

- **HTML**: `<thead>`, `<tbody>`, `<tfoot>` sections
- **Implementation**: 
- Check `data-role="header"` on rows for `<thead>` (or `role="header"` for ARIA)
- Check `data-role="footer"` on rows for `<tfoot>` (or `role="footer"` for ARIA)
- Current implementation assumes first row is header; enhance to check role attribute
- Use `<tr data-role="header">` for semantic markup

**Table Attributes**

- **HTML**: Apply `id`, `role` (for ARIA), and `data-asciidoc-*` attributes to `<table>` element
- **Implementation**: Ensure all table attributes are output as standard HTML attributes or `data-asciidoc-*` (no `class`)

### Phase 7: List Enhancements

**Callout Lists**

- **Challenge**: Callouts are typically rendered as numbered markers in code blocks
- **HTML Solution**: 
- For callout list items: `<li data-asciidoc-callout="1"><span data-role="callout-marker">1</span> description</li>`
- Link from code block callout markers to list items
- **Implementation**: 
- Check `callout` attribute on list items
- Render callout number as visual marker with `data-role="callout-marker"`
- Add `data-asciidoc-callout` attribute for linking

**List Continuations**

- **Challenge**: Continuation marker `+` adds content to previous item
- **HTML Solution**: Already handled structurally (content is child of list item)
- **Implementation**: No special HTML needed, structure already correct

**List Item Attributes**

- **HTML**: Apply `id`, `role` (for ARIA), and `data-asciidoc-*` to `<li>`, `<dt>`, `<dd>` elements
- **Implementation**: Ensure all list item attributes are output as standard HTML attributes or `data-asciidoc-*` (no `class`)
- **Term attribute**: For labeled lists, ensure term is in `<dt>`

### Phase 8: Section Enhancements

**Appendix**

### Phase 9: Attribute Reflection

**General Rule**: If an AsciiDoc attribute (e.g., `[foo="bar"]`) does not map to a standard HTML5 attribute (`id`, `href`, `src`, `alt`, `title`, `width`, `height`, `role` for ARIA, etc.), map it to `data-asciidoc-foo="bar"`.

**Examples**:
- `[variant="note"]` on Admonition → `data-asciidoc-variant="note"`
- `[callout="1"]` on ListItem → `data-asciidoc-callout="1"`
- `[appendix="true"]` on Section → `data-asciidoc-appendix="true"`
- `[discrete="true"]` on Section → `data-asciidoc-discrete="true"`
- `[ref="fn1"]` on Footnote → `data-asciidoc-ref="fn1"`
- Custom attributes like `[custom="value"]` → `data-asciidoc-custom="value"`

**Standard HTML Attributes** (preserve as-is):
- `id`, `href`, `src`, `alt`, `title`, `width`, `height`, `colspan`, `rowspan`
- `role` (for ARIA semantics)
- `type`, `controls`, `autoplay`, `loop`, `poster` (for media elements)

**Implementation**: 
- Check each attribute against standard HTML5 attribute list
- If standard → use as-is
- If not standard → prefix with `data-asciidoc-` and use as data attribute

- **HTML Solution**: Add `data-asciidoc-appendix="true"` to section
- **Implementation**: Check `appendix` attribute and add as data attribute

**Discrete**

- **HTML Solution**: Add `data-asciidoc-discrete="true"` to section (discrete sections have no number)
- **Implementation**: Check `discrete` attribute and add as data attribute

## Implementation Strategy

### File: `lib/converter.go`

1. **Add missing case handlers** in `toHTML` function:

- `Superscript`, `Subscript`, `Highlight`
- `VerseBlock`, `OpenBlock`
- `InlineMacro` (comprehensive handling)

2. **Enhance existing handlers**:

- `Table`: Add cell alignment, colspan, rowspan, row roles
- `List`: Add callout rendering, item attributes
- `Section`: Add appendix, discrete attributes
- `Link`: Handle cross-references (target attribute)

3. **Add footnote collection**:

- Create footnote registry during document traversal
- Render footnotes at document end
- Link footnote references

4. **Add TOC generation**:

- Traverse document to collect sections
- Generate TOC HTML structure
- Handle TOC macro rendering

5. **Enhance inline content handling**:

- Update `toHTMLInlineContent` to handle all inline macros
- Ensure proper nesting and escaping

### Testing Strategy

1. **Unit tests** for each new node type
2. **Integration tests** for complex features (footnotes, TOC, cross-refs)
3. **Visual tests** to verify HTML output renders correctly
4. **XHTML tests** to ensure well-formed output

## Success Criteria

- All AST node types have HTML output
- All attributes are preserved in HTML (as standard HTML attributes or `data-asciidoc-*` attributes)
- **No `class` attributes are used** (reserved for user-defined classes only)
- HTML output is valid and semantic
- XHTML output is well-formed
- Features without direct HTML equivalents have creative but standards-compliant solutions
- Footnotes are properly linked and rendered
- Cross-references work correctly
- Tables have full attribute support (using `data-align` for alignment)
- Lists support callouts and continuations
- All macros have appropriate HTML representation
- All structural semantics use `data-role` or `data-asciidoc-*` attributes