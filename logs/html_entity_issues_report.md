# HTML Entity Escaping Issues Report

**Date:** Generated from testbed/*.html files  
**Files Analyzed:** 15 HTML files in ./testbed/

## Executive Summary

Multiple issues with HTML entity encoding have been identified in the generated HTML files. The main problems are:

1. **HTML entities appearing in visible content** where they should render as characters
2. **Double escaping** of HTML entities
3. **Link titles appearing as text** instead of being parsed as attributes
4. **Code block content** with escaped quotes that may not render correctly

## Critical Issues

### Issue #1: HTML Entities in Visible Content (100+ occurrences)

**Severity:** HIGH  
**Pattern:** `&#34;` (double quote) and `&#39;` (single quote) appearing in visible text

**Examples:**

1. **Code Blocks** - `15-complex-documents.html:70`
   ```html
   <pre><code data-asciidoc-language="json">{
     &#34;name&#34;: &#34;John Doe&#34;,
     &#34;email&#34;: &#34;john@example.com&#34;
   }</code></pre>
   ```
   **Expected:** Quotes should render as `"` in the browser
   **Actual:** HTML entities `&#34;` are visible (if viewing source) or not rendering correctly

2. **Link Text** - `05-links.html:18`
   ```html
   <p>This is a <a href="https://example.com">link with title, title=&#34;Link Title&#34;</a>.</p>
   ```
   **Problem:** The `title="Link Title"` part should be a link attribute, not part of the link text

3. **Headers** - `01-headers-and-structure.html:39`
   ```html
   <h1>Header with &#34;quotes&#34; and &#39;single quotes&#39;</h1>
   ```
   **Expected:** Should display as: `Header with "quotes" and 'single quotes'`
   **Issue:** Entities may not render correctly in all contexts

4. **Paragraphs** - `07-code-blocks.html:24`
   ```html
   <p>This paragraph has <code>code with &#34;quotes&#34;</code> inside.</p>
   ```
   **Expected:** Should display quotes correctly in code

### Issue #2: Double Escaping (10+ occurrences)

**Severity:** MEDIUM  
**Pattern:** `&amp;lt;`, `&amp;quot;`, `&amp;amp;` instead of `&lt;`, `&quot;`, `&amp;`

**Examples:**

1. `01-headers-and-structure.html:60`
   ```html
   <h1>Header with &amp;lt;tags&amp;gt;</h1>
   ```
   **Should be:** `<h1>Header with &lt;tags&gt;</h1>`

2. `04-lists.html:204`
   ```html
   <li>Item with &amp;ampersand&amp;</li>
   ```
   **Should be:** `<li>Item with &amp; ampersand</li>`

**Root Cause:** Content is being HTML-escaped twice - once during parsing/processing and once during HTML generation.

### Issue #3: Link Titles as Text Content

**Severity:** HIGH  
**Pattern:** Link titles appearing in link text instead of as `title` attribute

**Examples:**

- `05-links.html:18-19, 48-50`
  ```html
  <a href="https://example.com">link with title, title=&#34;Link Title&#34;</a>
  ```
  **Expected:** 
  ```html
  <a href="https://example.com" title="Link Title">link with title</a>
  ```

**Root Cause:** AsciiDoc link syntax with titles is not being parsed correctly. The parser is treating the title as part of the link text rather than extracting it as an attribute.

### Issue #4: Table Attributes with Entities

**Severity:** MEDIUM  
**Pattern:** Table column specifications with escaped quotes

**Examples:**

- `08-tables.html:17-18`
  ```html
  <p>[cols=&#34;1,1,1&#34;]</p>
  <table data-asciidoc-cols="&#34;1">
  ```
  **Problem:** The `cols` attribute value contains `&#34;1` instead of `"1"`

## Files Affected

All 15 HTML files in ./testbed/ have some form of these issues:

1. **01-headers-and-structure.html** - Headers with entities, double escaping
2. **02-paragraphs-and-text.html** - Minor issues
3. **03-emphasis-and-formatting.html** - Minor issues
4. **04-lists.html** - Entities in list items, double escaping
5. **05-links.html** - **CRITICAL:** Link titles as text (26 occurrences)
6. **06-images.html** - Minor issues
7. **07-code-blocks.html** - Code with entities
8. **08-tables.html** - Table attributes with entities
9. **09-blockquotes.html** - Code in blockquotes with entities, double escaping
10. **10-frontmatter.html** - Minor issues
11. **11-horizontal-rules.html** - Minor issues
12. **12-html-in-markdown.html** - Minor issues
13. **13-escaped-characters.html** - Double escaping, entities
14. **14-edge-cases.html** - Minor issues
15. **15-complex-documents.html** - **CRITICAL:** Code blocks with entities (43 occurrences)

## Root Cause Analysis

### HTML Escaping Behavior

The Go `html.EscapeString()` function correctly escapes:
- `"` → `&#34;`
- `'` → `&#39;`
- `<` → `&lt;`
- `>` → `&gt;`
- `&` → `&amp;`

This is **correct** for HTML content. However, the issues arise from:

1. **Double Escaping:** Content that's already escaped is being escaped again
2. **Parsing Issues:** AsciiDoc syntax (like link titles) not being parsed correctly
3. **Context Issues:** Content that should be in attributes ending up in text content

### Specific Problems

1. **Code Blocks:** Quotes in code are correctly escaped as `&#34;`, which should render as `"` in browsers. If they're not rendering, it may be a display issue or the entities are being double-processed.

2. **Link Titles:** The AsciiDoc parser is not extracting link titles from syntax like:
   ```
   link:url[text, title="Title"]
   ```
   Instead, it's treating the entire `text, title="Title"` as the link text.

3. **Double Escaping:** Some content is being escaped during parsing/processing, then escaped again during HTML generation.

## Recommendations

1. **Fix Link Title Parsing:** Update the AsciiDoc parser to correctly extract link titles and place them in the `title` attribute instead of link text.

2. **Review Escaping Logic:** Ensure content is only escaped once, at the HTML generation stage.

3. **Code Block Handling:** Verify that code block content is being handled correctly - quotes should be escaped for HTML but should render correctly in browsers.

4. **Attribute Values:** Ensure attribute values are properly escaped but don't contain double-escaped entities.

5. **Test in Browsers:** Verify that `&#34;` entities actually render as `"` in browsers, as they should according to HTML spec.

## Impact

- **User Experience:** HTML entities visible in rendered content would be confusing
- **Accessibility:** Link titles not in `title` attributes reduce accessibility
- **Code Examples:** Code blocks with visible entities instead of quotes are hard to read
- **Standards Compliance:** While technically valid HTML, the output may not match expected behavior

