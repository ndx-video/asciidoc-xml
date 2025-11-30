#!/bin/bash
# Generate corrupt/malicious test files for parser security testing

CORRUPT_DIR="$(dirname "$0")"

# Function to create a corrupt file
create_file() {
    local filename="$1"
    local content="$2"
    echo -e "$content" > "$CORRUPT_DIR/$filename"
    echo "Created: $filename"
}

# 1. Extremely long line (but not too long - 10KB)
create_file "01-extremely-long-line.md" "$(printf 'A%.0s' {1..10000})"

# 2. Malformed frontmatter - unclosed
create_file "02-unclosed-frontmatter.md" "---
title: Test
author: Test
# Missing closing ---"

# 3. Malformed frontmatter - nested incorrectly
create_file "03-nested-frontmatter.md" "---
---
title: Test
---
---"

# 4. Frontmatter with null bytes
create_file "04-frontmatter-null-bytes.md" "$(printf '---\ntitle: Test\x00author: Test\n---')"

# 5. Frontmatter with control characters
create_file "05-frontmatter-control-chars.md" "$(printf '---\ntitle: Test\x01\x02\x03author: Test\n---')"

# 6. Extremely deep nesting (but reasonable - 20 levels)
create_file "06-deep-nesting.md" "$(
    for i in {1..20}; do
        printf '%*s> ' $i ''
    done
    echo 'Deeply nested blockquote'
)"

# 7. Unclosed code block
create_file "07-unclosed-code-block.md" "\`\`\`
Code block that never closes"

# 8. Unclosed HTML tags
create_file "08-unclosed-html.md" "<div>
<p>Unclosed paragraph
<strong>Unclosed bold"

# 9. Malformed table - inconsistent columns
create_file "09-malformed-table.md" "| Col 1 | Col 2 |
|-------|
| Data  | Data  | Data  |"

# 10. Table with extreme column count (but reasonable - 50)
create_file "10-many-columns-table.md" "$(
    echo -n "|"
    for i in {1..50}; do
        echo -n " Col$i |"
    done
    echo
    echo -n "|"
    for i in {1..50}; do
        echo -n "-----|"
    done
    echo
)"

# 11. Injection attempt - XSS in links
create_file "11-xss-link.md" "[Click me](javascript:alert('XSS'))"

# 12. Injection attempt - XSS in images
create_file "12-xss-image.md" "![Alt](javascript:alert('XSS'))"

# 13. Injection attempt - HTML script tag
create_file "13-html-script.md" "<script>alert('XSS')</script>"

# 14. Injection attempt - HTML with event handlers
create_file "14-html-events.md" "<img src=x onerror=alert('XSS')>"

# 15. Path traversal attempt in links
create_file "15-path-traversal.md" "[Link](../../../../etc/passwd)"

# 16. Path traversal attempt in images
create_file "16-path-traversal-image.md" "![Image](../../../../etc/passwd)"

# 17. Extremely long header (but reasonable - 5KB)
create_file "17-long-header.md" "$(printf '# %s\n' "$(printf 'A%.0s' {1..5000})")"

# 18. Header with null bytes
create_file "18-header-null-bytes.md" "$(printf '# Header\x00with\x00null\x00bytes')"

# 19. Link with extremely long URL (but reasonable - 2KB)
create_file "19-long-url.md" "[Link]($(printf 'https://example.com/%s' "$(printf 'a%.0s' {1..2000})"))"

# 20. Image with extremely long alt text (but reasonable - 5KB)
create_file "20-long-alt-text.md" "![$(printf 'A%.0s' {1..5000})](image.png)"

# 21. Code block with extremely long line (but reasonable - 10KB)
create_file "21-long-code-line.md" "\`\`\`
$(printf 'A%.0s' {1..10000})
\`\`\`"

# 22. Table with extremely long cell (but reasonable - 5KB)
create_file "22-long-table-cell.md" "| Column 1 | Column 2 |
|----------|----------|
| $(printf 'A%.0s' {1..5000}) | Data |"

# 23. List with extremely deep nesting (but reasonable - 30 levels)
create_file "23-deep-list.md" "$(
    for i in {1..30}; do
        printf '%*s* ' $i ''
        echo "Level $i"
    done
)"

# 24. Blockquote with extremely deep nesting (but reasonable - 30 levels)
create_file "24-deep-blockquote.md" "$(
    for i in {1..30}; do
        printf '%*s> ' $i ''
        echo "Level $i"
    done
)"

# 25. Malformed YAML in frontmatter
create_file "25-malformed-yaml.md" "---
title: Test
invalid: [unclosed array
another: {unclosed object
---"

# 26. Frontmatter with extremely long key (but reasonable - 1KB)
create_file "26-long-yaml-key.md" "---
$(printf 'A%.0s' {1..1000}): value
---"

# 27. Frontmatter with extremely long value (but reasonable - 10KB)
create_file "27-long-yaml-value.md" "---
title: $(printf 'A%.0s' {1..10000})
---"

# 28. YAML array with extremely many items (but reasonable - 1000)
create_file "28-large-yaml-array.md" "$(
    echo "---"
    echo "tags:"
    for i in {1..1000}; do
        echo "  - tag$i"
    done
    echo "---"
)"

# 29. Invalid UTF-8 sequences
create_file "29-invalid-utf8.md" "$(printf '\xFF\xFE\xFD')"

# 30. Mixed valid and invalid UTF-8
create_file "30-mixed-utf8.md" "Valid text $(printf '\xFF\xFE') more valid text"

# 31. Control characters in content
create_file "31-control-chars.md" "$(printf 'Text\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F')"

# 32. Carriage return only (no newline)
create_file "32-cr-only.md" "$(printf 'Line 1\rLine 2\r')"

# 33. Form feed characters
create_file "33-form-feed.md" "$(printf 'Text\x0CText')"

# 34. Vertical tab characters
create_file "34-vertical-tab.md" "$(printf 'Text\x0BText')"

# 35. Zero-width characters
create_file "35-zero-width.md" "$(printf 'Text\u200B\u200C\u200DText')"

# 36. Bidirectional override characters
create_file "36-bidi-override.md" "$(printf 'Text\u202Ereversed\u202CText')"

# 37. Extremely many headers (but reasonable - 1000)
create_file "37-many-headers.md" "$(
    for i in {1..1000}; do
        echo "# Header $i"
    done
)"

# 38. Extremely many links (but reasonable - 1000)
create_file "38-many-links.md" "$(
    for i in {1..1000}; do
        echo "[Link $i](https://example.com/$i)"
    done
)"

# 39. Circular reference links
create_file "39-circular-refs.md" "[Link 1][ref1]

[ref1]: https://example.com
[ref2]: See [Link 1][ref1]"

# 40. Self-referencing link
create_file "40-self-ref-link.md" "[Link][link]

[link]: #link"

# 41. Malformed emphasis - unclosed
create_file "41-unclosed-bold.md" "**Bold text that never closes"

# 42. Malformed emphasis - nested incorrectly
create_file "42-nested-emphasis.md" "***Bold*italic**"

# 43. Extremely many backticks (but reasonable - 100)
create_file "43-many-backticks.md" "$(printf '\`%.0s' {1..100})"

# 44. Code block with mismatched fences
create_file "44-mismatched-fences.md" "\`\`\`
Code
\`\`\`\`"

# 45. Table with mismatched pipes
create_file "45-mismatched-pipes.md" "| Col 1 | Col 2
|-------|-------|
| Data  | Data  |"

# 46. Extremely many horizontal rules (but reasonable - 1000)
create_file "46-many-hrules.md" "$(
    for i in {1..1000}; do
        echo "---"
    done
)"

# 47. Mixed line endings (CRLF and LF)
create_file "47-mixed-line-endings.md" "$(printf 'Line 1\r\nLine 2\nLine 3\r\n')"

# 48. Only line feeds (no content)
create_file "48-only-lfs.md" "$(printf '\n%.0s' {1..100})"

# 49. Only carriage returns
create_file "49-only-crs.md" "$(printf '\r%.0s' {1..100})"

# 50. Extremely many spaces (but reasonable - 10KB)
create_file "50-many-spaces.md" "$(printf ' %.0s' {1..10000})"

# 51. Extremely many tabs (but reasonable - 10KB)
create_file "51-many-tabs.md" "$(printf '\t%.0s' {1..10000})"

# 52. SQL injection attempt in frontmatter
create_file "52-sql-injection.md" "---
title: '; DROP TABLE users; --
---"

# 53. Command injection attempt
create_file "53-command-injection.md" "\`\`\`bash
\$(rm -rf /)
\`\`\`"

# 54. XML/XXE attempt
create_file "54-xml-xxe.md" "<?xml version=\"1.0\"?>
<!DOCTYPE foo [<!ENTITY xxe SYSTEM \"file:///etc/passwd\">]>
<foo>&xxe;</foo>"

# 55. Extremely long word (but reasonable - 5KB)
create_file "55-long-word.md" "$(printf 'A%.0s' {1..5000})"

# 56. Empty file
create_file "56-empty.md" ""

# 57. Only whitespace
create_file "57-only-whitespace.md" "     "

# 58. Only special characters
create_file "58-only-special.md" "!@#$%^&*()_+-=[]{}|;':\",./<>?"

# 59. Extremely many empty lines (but reasonable - 1000)
create_file "59-many-empty-lines.md" "$(printf '\n%.0s' {1..1000})"

# 60. Binary data mixed with text
create_file "60-binary-data.md" "$(printf 'Text\x00\x01\x02\x03\x04\x05Text')"

echo "Done generating corrupt test files!"

