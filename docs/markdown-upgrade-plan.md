# Markdown Converter Upgrade Plan

## Overview
Comprehensive upgrade to support the official CommonMark specification and GitHub Flavored Markdown (GFM) extensions.

## Current Feature Status

### ✅ Currently Supported
- ATX-style headers (# through ######)
- Bold text (**text** and __text__)
- Italic text (*text* and _text_)
- Unordered lists (*, -, +)
- Ordered lists (1., 2., 3.)
- Nested lists
- Inline links [text](url)
- Inline images ![alt](src)
- Fenced code blocks ```
- Basic tables
- Blockquotes (>)
- Horizontal rules (---, ***, ___)
- YAML frontmatter
- Inline code `code`

### ❌ Missing CommonMark Features
1. **Setext-style headers** (=== and ---)
2. **Reference-style links** [text][ref] and images ![alt][ref]
3. **Link reference definitions** [ref]: url "title"
4. **Autolinks** <url> and <email@example.com>
5. **Hard line breaks** (two spaces + newline)
6. **Indented code blocks** (4 spaces)
7. **Escaped characters** (backslash escaping)
8. **HTML blocks and spans**
9. **Better emphasis word boundaries** (underscore handling)

### ❌ Missing GFM Features
1. **Task lists** (- [ ] and - [x])
2. **Strikethrough** (~~text~~)
3. **Table alignment** (left, center, right)
4. **Autolinks** (URLs and emails without angle brackets)
5. **Disallowed raw HTML** handling

## Implementation Plan

### Phase 1: CommonMark Core Features
1. Setext-style headers
2. Reference-style links and images
3. Link reference definitions collection
4. Autolinks
5. Hard line breaks
6. Indented code blocks
7. Escaped character handling

### Phase 2: GFM Extensions
1. Task lists
2. Strikethrough
3. Table alignment
4. Enhanced autolinks

### Phase 3: Advanced Features
1. HTML block and span support
2. Better emphasis word boundaries
3. Enhanced escaping

## Implementation Details

### Setext Headers
- Detect lines with `===` or `---` underlines
- Match to previous non-empty line
- Convert to appropriate AsciiDoc header level

### Reference Links
- Collect link reference definitions during first pass
- Resolve references during inline conversion
- Support implicit reference links [text][]

### Autolinks
- Detect URLs (http://, https://, ftp://)
- Detect email addresses
- Convert to AsciiDoc links

### Hard Line Breaks
- Detect two spaces at end of line
- Convert to AsciiDoc line break (+\n)

### Indented Code Blocks
- Detect 4+ spaces at start of line
- Treat as code block until non-indented line
- Convert to AsciiDoc code block

### Task Lists
- Detect `- [ ]` and `- [x]` in lists
- Convert to AsciiDoc checklist format

### Strikethrough
- Detect `~~text~~`
- Convert to AsciiDoc [.line-through]#text#

### Table Alignment
- Parse alignment from separator row (|:---|:---:|---:|)
- Convert to AsciiDoc column specifications

