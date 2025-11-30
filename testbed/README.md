# Markdown Test Suite

This directory contains comprehensive test files for the Markdown to AsciiDoc converter. Each file focuses on specific Markdown features and includes extensive examples and edge cases.

## Test Files

### 01-headers-and-structure.md
- ATX style headers (# through ######)
- Setext style headers (=== and ---)
- Headers with emphasis, links, images
- Headers with special characters, numbers, emojis
- Multiple headers in sequence

### 02-paragraphs-and-text.md
- Basic paragraphs
- Paragraphs with line breaks
- Paragraphs with special characters
- Paragraphs with URLs, code references
- Paragraphs with emphasis combinations
- Paragraphs with links and images
- Paragraphs with HTML
- Paragraphs with escaped characters
- Paragraphs with Unicode
- Very long paragraphs
- Mixed content paragraphs

### 03-emphasis-and-formatting.md
- Bold text (** and __)
- Italic text (* and _)
- Bold and italic combined (*** and ___)
- Emphasis edge cases
- Emphasis with punctuation and special characters
- Emphasis with numbers, links, images
- Emphasis with code and HTML
- Strikethrough (~~)
- Emphasis escaping
- Emphasis with Unicode
- Emphasis combinations
- Emphasis in headers, lists, blockquotes, tables

### 04-lists.md
- Unordered lists (*, -, +)
- Ordered lists (1., 2., 3.)
- Nested lists (mixed ordered/unordered)
- Lists with multiple paragraphs
- Lists with code blocks and blockquotes
- Lists with emphasis, links, images
- Task lists (- [ ] and - [x])
- Nested task lists
- Lists with HTML
- Lists with special characters and Unicode
- Tight vs loose lists
- Empty list items
- Complex nested lists

### 05-links.md
- Inline links [text](url)
- Links with titles
- Links with different protocols (http, https, ftp, mailto, file)
- Links with paths, query strings, fragments
- Reference-style links [text][ref]
- Implicit and shortcut reference links
- Reference links with titles
- Autolinks <https://example.com>
- Links with emphasis
- Links in lists, blockquotes, headers, tables
- Links with images
- Relative links
- Anchor links (#section)
- Broken links
- Links with escaped characters
- Links with Unicode
- Multiple links in one paragraph
- Links with empty text

### 06-images.md
- Inline images ![alt](src)
- Images with titles
- Images with different formats (png, jpg, gif, svg, webp)
- Images with paths (relative, absolute)
- Images with URLs
- Reference-style images ![alt][ref]
- Images with special characters in alt text
- Images with Unicode in alt text
- Images with emphasis in alt text
- Images with links
- Images in lists, blockquotes, headers, tables
- Images with empty alt text
- Multiple images in one paragraph
- Images with HTML
- Images with escaped characters
- Images with special URL schemes
- Complex image combinations

### 07-code-blocks.md
- Inline code `code`
- Inline code with special characters
- Inline code with backticks
- Fenced code blocks ```
- Code blocks with language identifiers (javascript, python, go, html, css, json, yaml, bash, sql, xml, markdown)
- Indented code blocks (4 spaces)
- Code blocks with empty lines
- Code blocks with special characters and Unicode
- Code blocks with backticks inside
- Code blocks in lists and blockquotes
- Code blocks with line numbers and file names
- Code blocks with metadata and comments
- Code blocks with strings and regular expressions
- Code blocks with HTML entities
- Code blocks with escaped characters
- Code blocks preserving formatting
- Code blocks with language aliases
- Code blocks with invalid/empty languages

### 08-tables.md
- Basic tables
- Tables with header only
- Tables with alignment (left, center, right)
- Tables with mixed alignment
- Tables with empty cells
- Tables with long content
- Tables with special characters
- Tables with numbers
- Tables with links and images
- Tables with emphasis and code
- Tables with HTML and Unicode
- Tables without pipes
- Tables with extra spaces
- Single column tables
- Tables with many columns/rows
- Tables in lists and blockquotes
- Tables with complex content
- Tables with escaped pipes
- Tables with very wide/narrow content
- Tables with mixed data types

### 09-blockquotes.md
- Basic blockquotes
- Blockquotes with single line
- Blockquotes with paragraphs
- Blockquotes with emphasis, links, images
- Blockquotes with lists and code blocks
- Blockquotes with headers and horizontal rules
- Blockquotes with tables
- Nested blockquotes (multiple levels)
- Blockquotes with mixed content
- Blockquotes with special characters and Unicode
- Blockquotes with HTML and escaped characters
- Blockquotes starting/ending document
- Multiple blockquotes
- Blockquotes with attribution and citation
- Blockquotes with long content
- Blockquotes with line breaks
- Blockquotes with empty lines
- Blockquotes in lists
- Blockquotes with complex nesting
- Blockquotes preserving formatting

### 10-frontmatter.md
- Simple frontmatter
- Frontmatter with arrays
- Frontmatter with nested structures
- Frontmatter with complex arrays
- Frontmatter with empty values
- Frontmatter with quoted values
- Frontmatter with special characters
- Frontmatter with numbers, dates, booleans
- Frontmatter with null values
- Frontmatter with mixed types
- Frontmatter with long values and Unicode
- Frontmatter with URLs and email
- Frontmatter with file paths
- Frontmatter with arrays of objects
- Frontmatter with deep nesting
- Frontmatter with multiple arrays
- Frontmatter with empty arrays
- Frontmatter with single item arrays
- Frontmatter with mixed array formats
- Frontmatter with multiline strings
- Frontmatter with custom fields
- Frontmatter with version numbers, currency, percentages
- Frontmatter with time values, IP addresses
- Frontmatter with version control, social media
- Frontmatter with geographic data, ratings
- Frontmatter with complex nested arrays
- Frontmatter with escaped characters
- Frontmatter with JSON-like structures
- Frontmatter with HTML and Markdown
- Frontmatter ending with three dots
- Frontmatter with inline arrays

### 11-horizontal-rules.md
- Three hyphens (---)
- Three asterisks (***)
- Three underscores (___)
- Four, five, and more characters
- Horizontal rules with spaces
- Horizontal rules between paragraphs and sections
- Horizontal rules in lists and blockquotes
- Horizontal rules ending/starting document
- Multiple horizontal rules
- Horizontal rules with different characters mixed
- Horizontal rules with text on same line (should not work)
- Horizontal rules in code blocks (should not work)
- Horizontal rules with emphasis/links (should not work)
- Horizontal rules with trailing spaces
- Horizontal rules minimum valid (three characters)
- Horizontal rules maximum reasonable
- Horizontal rules in different contexts
- Horizontal rules with code and escaped characters
- Horizontal rules with HTML

### 12-html-in-markdown.md
- Inline HTML tags (strong, em, code)
- HTML links and images
- HTML line breaks
- HTML paragraphs and headers
- HTML lists and blockquotes
- HTML code blocks and tables
- HTML with attributes
- HTML comments
- HTML entities
- HTML mixed with Markdown
- HTML in lists, blockquotes, tables
- HTML script and style tags
- HTML form elements
- HTML details/summary
- HTML abbreviations, citations, definitions
- HTML keyboard input, marked text
- HTML subscript and superscript
- HTML time, variables, sample output
- HTML address and preformatted text
- HTML with nested tags
- HTML self-closing tags
- HTML with data attributes
- HTML with ARIA attributes
- HTML with class and ID
- HTML with inline styles
- HTML escaped characters
- HTML mixed content
- HTML preserving whitespace
- HTML with Unicode
- HTML comments in content
- HTML empty tags
- HTML with attributes containing quotes
- HTML with special characters in attributes
- HTML block vs inline
- HTML nested properly
- HTML mixed with Markdown formatting
- HTML in code blocks (should be literal)
- HTML escaping in Markdown

### 13-escaped-characters.md
- Escaping asterisks, underscores, backticks
- Escaping brackets, hashtags
- Escaping plus, minus, periods
- Escaping greater than, pipes, tildes
- Escaping parentheses, brackets, braces
- Escaping in different contexts (paragraphs, lists, blockquotes, code blocks)
- Escaping URLs and email
- Escaping HTML
- Escaping multiple characters
- Escaping at start/end of line
- Escaping with spaces
- Escaping special sequences
- Escaping in links and images
- Escaping in headers and tables
- Escaping backslashes
- Escaping combinations
- Escaping Unicode and numbers
- Escaping quotes, dollar signs, percent signs
- Escaping ampersands, angle brackets
- Escaping square brackets, curly braces
- Escaping parentheses
- Escaping punctuation marks
- Escaping multiple special characters
- Escaping in different positions
- Escaping with adjacent text
- Escaping preserving literal meaning
- Escaping edge cases

### 14-edge-cases.md
- Empty document
- Document with only whitespace/newlines
- Document starting/ending with header
- Document with only header/paragraph
- Adjacent headers
- Headers without content
- Paragraphs without separation
- Lists without blank lines
- Empty list items
- Lists with only spaces
- Tables with empty cells
- Tables with only headers
- Code blocks with only whitespace/newlines
- Blockquotes with only spaces/newlines
- Links/images with empty text
- Emphasis with empty content
- Emphasis adjacent to text
- Links/images adjacent to text
- Headers with only spaces/special characters
- Lists with mixed markers
- Lists starting with different numbers
- Tables with inconsistent separators
- Code blocks with different fence lengths
- Nested structures maximum depth
- Mixed formatting edge cases
- Links/images/code in different positions
- Emphasis in different positions
- Special characters at boundaries
- Unicode edge cases
- Numbers in different contexts
- URLs in different formats
- Email formats
- File paths
- Version numbers
- Dates and times
- Percentages and currency
- IP addresses
- Hashtags and mentions
- Mixed special characters
- Empty frontmatter
- Frontmatter with only spaces/newlines
- Multiple frontmatter blocks
- Frontmatter without separators
- Content starting/ending immediately

### 15-complex-documents.md
- Technical documentation example
- Blog post example
- README example
- Meeting notes example
- Recipe example
- Product documentation example
- Academic paper example
- Mixed content example
- Comprehensive document with all features

## Usage

These test files can be used to:

1. **Test conversion accuracy** - Verify that Markdown features convert correctly to AsciiDoc
2. **Identify edge cases** - Find scenarios where conversion might fail
3. **Validate fixes** - Ensure fixes work across all test cases
4. **Regression testing** - Prevent previously fixed issues from reoccurring
5. **Documentation** - Serve as examples of supported Markdown features

## Running Tests

To test the converter with these files:

```bash
# Convert a single file
go run cli/adc.go md2adoc testbed/01-headers-and-structure.md

# Convert all test files
for file in testbed/*.md; do
    go run cli/adc.go md2adoc "$file"
done
```

## Notes

- Files are organized by feature category
- Each file contains comprehensive examples and edge cases
- Files avoid stress testing (no extremely long content)
- Real-world examples are included in the complex documents file
- All standard Markdown features are covered
- GitHub Flavored Markdown extensions are included where applicable

