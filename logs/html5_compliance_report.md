# HTML5 Standards Compliance Report

**Date:** Generated from testbed/*.html files  
**Files Analyzed:** 16 HTML files in ./testbed/  
**Files with Issues:** 15 (93.75%)

## Executive Summary

All HTML files in the testbed directory were analyzed for HTML5 standards compliance. **Critical issues were found in 15 out of 16 files**, indicating systematic problems in the HTML generation code.

## Critical Issues Identified

### Issue #1: Missing Space Before Attributes (1,648+ occurrences)

**Severity:** ERROR  
**Affected Tags:** `<p>`, `<table>`, `<div>`, `<blockquote>`

**Pattern Found:**
- `<pclass="preamble">` instead of `<p class="preamble">` or `<p data-role="preamble">`
- `<tableclass="table">` instead of `<table class="table">` or proper data attributes
- `<divclass="example">` instead of `<div class="example">` or `<div data-role="example">`
- `<blockquoteclass="quote">` instead of `<blockquote class="quote">` or proper data attributes

**HTML5 Violation:** HTML5 requires whitespace between tag name and attributes.

**Example from 15-complex-documents.html:**
```html
Line 13: <pclass="preamble"><h1 id="complex_realworld_documents">Complex Real-World Documents</h1>
```

**Source (15-complex-documents.adoc):**
```asciidoc
Line 1: = Complex Real-World Documents
```

**Root Cause:** In `lib/converter.go`, the attribute string building is missing proper spacing. The format string `%s<p%s>` expects `attrs` to start with a space, but it appears attributes are being concatenated without the required leading space.

**Code Location:** `lib/converter.go` - Multiple cases (Paragraph, Table, Example, Quote, etc.)

---

### Issue #2: Invalid Nesting - Block Elements Inside `<p>` Tags (15 occurrences)

**Severity:** ERROR  
**HTML5 Violation:** The HTML5 content model prohibits block-level elements inside `<p>` tags.

**Pattern Found:**
```html
<pclass="preamble"><h1>...</h1></p>
```

**Example from 01-headers-and-structure.html:**
```html
Line 16: <pclass="preamble"><h1 id="level_1_header_atx_style">Level 1 Header (ATX style)</h1>
        </p>
```

**Source (01-headers-and-structure.adoc):**
```asciidoc
Line 1: = Headers and Document Structure
Line 8: = Level 1 Header (ATX style)
```

**Root Cause:** The preamble handling in the Document case wraps content in a `<div>`, but the paragraph itself is also being processed and output as a `<p>` tag containing block-level elements like `<h1>`.

**Code Location:** `lib/converter.go` - Document case (line ~59) and Paragraph case (line ~133)

**Fix Required:** Preamble paragraphs should not be output as `<p>` tags when they contain block-level content. They should either:
1. Use `<div>` instead of `<p>` for block content
2. Not wrap block elements in paragraph tags

---

### Issue #3: Mismatched Paragraph Tags (Multiple files)

**Severity:** WARNING  
**Pattern:** More opening `<p>` tags than closing `</p>` tags

**Affected Files:**
- 03-emphasis-and-formatting.html: 41 open, 40 close
- 04-lists.html: 11 open, 8 close
- 05-links.html: 81 open, 80 close
- 07-code-blocks.html: 86 open, 35 close (severe)
- 09-blockquotes.html: 332 open, 326 close
- 10-frontmatter.html: 56 open, 3 close (severe)
- And others...

**Root Cause:** Likely related to Issue #2 - when block elements are nested inside `<p>` tags, the paragraph structure becomes invalid and tags may not be properly closed.

---

## Detailed File-by-File Analysis

### 01-headers-and-structure.html
- **Issues:** 12 missing space errors, 1 invalid nesting
- **Source:** First heading `= Headers and Document Structure` creates preamble paragraph

### 08-tables.html  
- **Issues:** 20 missing space errors in table tags
- **Source:** Multiple `[cols="1,1,1"]` table definitions

### 09-blockquotes.html
- **Issues:** 194 missing space errors in blockquote tags
- **Most affected file** - nearly every blockquote has the spacing issue

### 15-complex-documents.html
- **Issues:** 9 missing space errors, 1 invalid nesting
- **Source:** Complex document with multiple sections, tables, and blockquotes

---

## Recommendations

1. **Fix attribute spacing:** Ensure all attribute strings start with a space when used in format strings
2. **Fix preamble handling:** Use `<div>` for preamble instead of `<p>` when block content is present
3. **Fix paragraph nesting:** Ensure `<p>` tags only contain inline content
4. **Add HTML5 validation:** Consider adding automated HTML5 validation in the test suite

## Files Analyzed

1. 01-headers-and-structure.html
2. 02-paragraphs-and-text.html
3. 03-emphasis-and-formatting.html
4. 04-lists.html
5. 05-links.html
6. 06-images.html
7. 07-code-blocks.html
8. 08-tables.html
9. 09-blockquotes.html
10. 10-frontmatter.html
11. 11-horizontal-rules.html
12. 12-html-in-markdown.html
13. 13-escaped-characters.html
14. 14-edge-cases.html
15. 15-complex-documents.html
16. README.html (compliant)

