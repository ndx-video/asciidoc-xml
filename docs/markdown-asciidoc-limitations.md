# Markdown to AsciiDoc Conversion Limitations

## Overview
This document outlines Markdown features (CommonMark and GFM) that either:
1. Don't have direct AsciiDoc equivalents
2. Require workarounds or approximations
3. Are partially supported

## Features with Direct Support ✅

All of these have been implemented:
- Setext headers → AsciiDoc headers
- Reference links → AsciiDoc links
- Autolinks (`<url>`) → AsciiDoc links
- Hard line breaks → AsciiDoc `+\n`
- Indented code blocks → AsciiDoc code blocks
- Task lists → AsciiDoc checklists
- Strikethrough → AsciiDoc `[.line-through]#text#`
- Table alignment → AsciiDoc column specs

## Features Requiring Workarounds ⚠️

### 1. HTML Blocks and Spans
**Status**: Can be supported via passthrough blocks

**Markdown**:
```markdown
<div class="custom">
  Content
</div>
```

**AsciiDoc Equivalent**:
```asciidoc
++++
<div class="custom">
  Content
</div>
++++
```

**Implementation**: AsciiDoc supports passthrough blocks (`++++`) that allow raw HTML. This can be implemented but requires detecting HTML blocks in Markdown and converting them appropriately.

**Note**: Currently marked as "pending" in implementation plan.

### 2. GFM Autolinks Without Angle Brackets
**Status**: Not automatically supported by AsciiDoc

**Markdown (GFM)**:
```markdown
Visit https://example.com or email user@example.com
```

**Issue**: GFM automatically converts URLs and emails to links even without `< >` brackets. AsciiDoc doesn't do this automatically.

**Workaround**: Could detect URLs/emails in text and convert them, but this is complex and error-prone (might convert URLs in code blocks, etc.).

**Current Status**: Only autolinks with angle brackets (`<url>`) are supported.

### 3. GFM Line Breaks in Tables
**Status**: Requires passthrough

**Markdown (GFM)**:
```markdown
| Column 1 | Column 2 |
|----------|----------|
| Line 1<br>Line 2 | Content |
```

**AsciiDoc**: Can support via inline passthrough, but requires explicit conversion:
```asciidoc
|Line 1 pass:[<br>]Line 2 |Content|
```

**Implementation**: Would need to detect `<br>` tags in table cells and convert to passthrough.

### 4. GFM Mentions and Issue References
**Status**: No direct equivalent

**Markdown (GFM)**:
```markdown
@username mentioned in #123
```

**Issue**: These are GitHub-specific features. AsciiDoc has no equivalent.

**Workaround**: Could convert to plain text or custom macros, but loses the GitHub-specific functionality.

### 5. GFM Emoji Shortcodes
**Status**: No direct equivalent

**Markdown (GFM)**:
```markdown
:smile: :heart: :rocket:
```

**Issue**: GitHub-specific emoji shortcodes. AsciiDoc doesn't support these.

**Workaround**: Could convert to Unicode emoji or leave as-is, but loses GitHub rendering.

### 6. GFM Disallowed Raw HTML
**Status**: Security/sanitization feature, not conversion

**Issue**: GFM has a list of HTML tags that are disallowed for security. This is about sanitization, not conversion.

**Note**: This is more of a security concern than a conversion limitation.

### 7. Multiple Consecutive Blank Lines
**Status**: Different handling

**Markdown**: Collapses multiple blank lines to a single paragraph break.

**AsciiDoc**: May preserve blank lines differently depending on context.

**Impact**: Minor - usually doesn't affect output significantly.

### 8. GFM Table Cell Alignment with Colons
**Status**: ✅ Implemented

**Note**: This was successfully implemented using AsciiDoc column specifications.

## Features That Could Be Added

### HTML Block Support
- **Feasibility**: High
- **Method**: Detect HTML blocks in Markdown, convert to AsciiDoc passthrough blocks (`++++`)
- **Complexity**: Medium (need to detect block vs inline HTML, handle edge cases)

### GFM Autolinks (without brackets)
- **Feasibility**: Medium
- **Method**: Regex-based URL/email detection in paragraphs
- **Complexity**: High (need to avoid converting in code blocks, code spans, etc.)
- **Risk**: False positives (converting non-URLs)

### GFM Mentions/Issues
- **Feasibility**: Low
- **Method**: Convert to plain text or custom AsciiDoc macros
- **Complexity**: Low
- **Value**: Low (GitHub-specific, loses functionality)

## Summary

### Fully Supported ✅
- All CommonMark core features
- Most GFM extensions (tables, task lists, strikethrough, table alignment)
- Reference-style links and images
- Autolinks with angle brackets

### Can Be Supported (Not Yet Implemented) ⚠️
- HTML blocks/spans (via passthrough)
- GFM autolinks without brackets (complex, error-prone)

### No Direct Equivalent ❌
- GFM mentions (`@username`)
- GFM issue references (`#123`)
- GFM emoji shortcodes (`:smile:`)
- GFM disallowed HTML (security feature, not conversion)

### Conclusion

The converter successfully handles **all major CommonMark features** and **most GFM extensions**. The remaining gaps are primarily:
1. **HTML passthrough** - Can be added if needed
2. **GitHub-specific features** - No equivalent in AsciiDoc (mentions, issues, emoji)
3. **GFM autolinks without brackets** - Complex to implement correctly

For most use cases, the converter provides comprehensive coverage of standard Markdown and GFM features.

