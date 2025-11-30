# Corrupt/Malicious Test Files

This directory contains test files designed to test the parser's resilience against:
- Malformed Markdown syntax
- Security vulnerabilities (XSS, injection attacks)
- Edge cases and boundary conditions
- Invalid input data
- Resource exhaustion attempts

## File Categories

### Malformed Syntax
- `02-unclosed-frontmatter.md` - Frontmatter without closing delimiter
- `03-nested-frontmatter.md` - Incorrectly nested frontmatter
- `07-unclosed-code-block.md` - Code block that never closes
- `08-unclosed-html.md` - Unclosed HTML tags
- `09-malformed-table.md` - Table with inconsistent columns
- `25-malformed-yaml.md` - Invalid YAML syntax
- `41-unclosed-bold.md` - Unclosed emphasis
- `42-nested-emphasis.md` - Incorrectly nested emphasis
- `44-mismatched-fences.md` - Code block with mismatched fence characters
- `45-mismatched-pipes.md` - Table with mismatched pipe characters

### Security Tests
- `11-xss-link.md` - XSS attempt in link URL
- `12-xss-image.md` - XSS attempt in image URL
- `13-html-script.md` - Script tag injection
- `14-html-events.md` - HTML event handler injection
- `15-path-traversal.md` - Path traversal in links
- `16-path-traversal-image.md` - Path traversal in images
- `52-sql-injection.md` - SQL injection attempt
- `53-command-injection.md` - Command injection attempt
- `54-xml-xxe.md` - XML External Entity attack attempt

### Resource Exhaustion
- `01-extremely-long-line.md` - Very long single line (10KB)
- `10-many-columns-table.md` - Table with many columns (50)
- `17-long-header.md` - Very long header (5KB)
- `19-long-url.md` - Very long URL (2KB)
- `20-long-alt-text.md` - Very long image alt text (5KB)
- `21-long-code-line.md` - Very long code line (10KB)
- `22-long-table-cell.md` - Very long table cell (5KB)
- `26-long-yaml-key.md` - Very long YAML key (1KB)
- `27-long-yaml-value.md` - Very long YAML value (10KB)
- `28-large-yaml-array.md` - YAML array with many items (1000)
- `37-many-headers.md` - Many headers (1000)
- `38-many-links.md` - Many links (1000)
- `46-many-hrules.md` - Many horizontal rules (1000)
- `50-many-spaces.md` - Many spaces (10KB)
- `51-many-tabs.md` - Many tabs (10KB)
- `55-long-word.md` - Very long word (5KB)
- `59-many-empty-lines.md` - Many empty lines (1000)

### Deep Nesting
- `06-deep-nesting.md` - Deeply nested blockquotes (20 levels)
- `23-deep-list.md` - Deeply nested lists (30 levels)
- `24-deep-blockquote.md` - Deeply nested blockquotes (30 levels)

### Invalid Characters
- `04-frontmatter-null-bytes.md` - Null bytes in frontmatter
- `05-frontmatter-control-chars.md` - Control characters in frontmatter
- `18-header-null-bytes.md` - Null bytes in header
- `29-invalid-utf8.md` - Invalid UTF-8 sequences
- `30-mixed-utf8.md` - Mixed valid/invalid UTF-8
- `31-control-chars.md` - Control characters in content
- `32-cr-only.md` - Carriage return only (no newline)
- `33-form-feed.md` - Form feed characters
- `34-vertical-tab.md` - Vertical tab characters
- `35-zero-width.md` - Zero-width characters
- `36-bidi-override.md` - Bidirectional override characters
- `60-binary-data.md` - Binary data mixed with text

### Edge Cases
- `39-circular-refs.md` - Circular reference links
- `40-self-ref-link.md` - Self-referencing link
- `43-many-backticks.md` - Many backticks (100)
- `47-mixed-line-endings.md` - Mixed CRLF and LF
- `48-only-lfs.md` - Only line feeds
- `49-only-crs.md` - Only carriage returns
- `56-empty.md` - Empty file
- `57-only-whitespace.md` - Only whitespace
- `58-only-special.md` - Only special characters

## Usage

These files should be tested to ensure:

1. **No crashes** - Parser should handle all files gracefully
2. **No security vulnerabilities** - Malicious content should be sanitized
3. **Reasonable resource usage** - Large files shouldn't cause memory issues
4. **Error handling** - Invalid input should produce appropriate errors
5. **Output safety** - Generated output should not contain executable code

## Testing

To test the converter with these files:

```bash
# Test a single corrupt file
go run cli/adc.go md2adoc testbed/corrupt/01-extremely-long-line.md

# Test all corrupt files (expect some to fail gracefully)
for file in testbed/corrupt/*.md; do
    echo "Testing: $file"
    go run cli/adc.go md2adoc "$file" 2>&1 | head -5
    echo "---"
done
```

## Notes

- Files are designed to test parser resilience, not to cause actual harm
- Resource limits are set to reasonable values (not extreme stress testing)
- Security tests use common attack patterns but are contained
- The parser should sanitize or reject malicious content appropriately
- Some files may legitimately fail to parse - that's expected and acceptable

## Regeneration

To regenerate all corrupt test files:

```bash
cd testbed/corrupt
./generate_corrupt.sh
```

