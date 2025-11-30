---
title: Code Blocks and Inline Code
categories:
  - markdown
  - code
tags:
  - programming
  - syntax highlighting
  - code blocks
---

# Code Blocks and Inline Code

## Inline Code

This paragraph contains `inline code` using backticks.

This paragraph has `code with spaces` in it.

This paragraph has `code_with_underscores` and `code-with-dashes`.

This paragraph has `code123` with numbers.

This paragraph has `code.with.dots` with dots.

---

## Inline Code with Special Characters

This paragraph has `code with "quotes"` inside.

This paragraph has `code with (parentheses)` inside.

This paragraph has `code with [brackets]` inside.

This paragraph has `code with {braces}` inside.

This paragraph has `code with <tags>` inside.

This paragraph has `code with $dollar$` inside.

This paragraph has `code with &ampersand&` inside.

---

## Inline Code with Backticks

This paragraph has ``code with `backtick` inside`` using double backticks.

This paragraph has ```code with ``double backticks`` inside``` using triple backticks.

---

## Inline Code Edge Cases

`code`text (no space)

`code`**bold** (adjacent formatting)

`code`*italic* (adjacent formatting)

`code`[link](https://example.com) (adjacent link)

---

## Fenced Code Blocks - No Language

```
This is a code block
with multiple lines
of code
```

---

## Fenced Code Blocks - JavaScript

```javascript
function greet(name) {
    console.log("Hello, " + name + "!");
}

greet("World");
```

---

## Fenced Code Blocks - Python

```python
def greet(name):
    print(f"Hello, {name}!")

greet("World")
```

---

## Fenced Code Blocks - Go

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

---

## Fenced Code Blocks - HTML

```html
<!DOCTYPE html>
<html>
<head>
    <title>Example</title>
</head>
<body>
    <h1>Hello, World!</h1>
</body>
</html>
```

---

## Fenced Code Blocks - CSS

```css
body {
    font-family: Arial, sans-serif;
    margin: 0;
    padding: 20px;
}

h1 {
    color: #333;
}
```

---

## Fenced Code Blocks - JSON

```json
{
    "name": "Example",
    "version": "1.0.0",
    "dependencies": {
        "package": "^1.2.3"
    }
}
```

---

## Fenced Code Blocks - YAML

```yaml
name: Example
version: 1.0.0
dependencies:
  package: ^1.2.3
```

---

## Fenced Code Blocks - Shell

```bash
#!/bin/bash
echo "Hello, World!"
ls -la
```

---

## Fenced Code Blocks - SQL

```sql
SELECT * FROM users
WHERE age > 18
ORDER BY name;
```

---

## Fenced Code Blocks - XML

```xml
<?xml version="1.0" encoding="UTF-8"?>
<root>
    <element attribute="value">Content</element>
</root>
```

---

## Fenced Code Blocks - Markdown

```markdown
# Header

**Bold** and *italic* text.

- List item
```

---

## Fenced Code Blocks - Multiple Languages

```javascript
// JavaScript code
```

```python
# Python code
```

```go
// Go code
```

---

## Indented Code Blocks

    This is an indented code block
    with multiple lines
    using four spaces

---

## Code Blocks with Empty Lines

```javascript
function example() {
    // First line
    
    // Third line after empty line
}
```

---

## Code Blocks with Special Characters

```javascript
const special = {
    "quotes": 'single quotes',
    "parentheses": (value),
    "brackets": [array],
    "braces": {object},
    "tags": <html>,
    "dollar": $value$,
    "ampersand": &value&
};
```

---

## Code Blocks with Unicode

```javascript
const unicode = {
    "café": "café",
    "你好": "你好",
    "こんにちは": "こんにちは"
};
```

---

## Code Blocks with Backticks Inside

````markdown
```
Code block with backticks
```
````

---

## Code Blocks with Triple Backticks Inside

`````markdown
```
Code block content
```
`````

---

## Code Blocks in Lists

* List item with code block:

  ```javascript
  console.log("Code in list");
  ```

* Another item

---

## Code Blocks in Blockquotes

> Blockquote with code block:
> 
> ```javascript
> console.log("Code in blockquote");
> ```

---

## Code Blocks with Line Numbers (if supported)

```javascript:1-5
function example() {
    return "Hello";
}
```

---

## Code Blocks with File Names (if supported)

```javascript:example.js
function example() {
    return "Hello";
}
```

---

## Code Blocks with Metadata

```javascript
// File: example.js
// Author: Test Author
// Date: 2024-01-15
function example() {
    return "Hello";
}
```

---

## Code Blocks with Comments

```javascript
// Single line comment
function example() {
    /* Multi-line
       comment */
    return "Hello";
}
```

---

## Code Blocks with Strings

```javascript
const single = 'Single quoted string';
const double = "Double quoted string";
const template = `Template string with ${variable}`;
```

---

## Code Blocks with Regular Expressions

```javascript
const regex1 = /pattern/g;
const regex2 = /pattern with \(escaped\)/g;
const regex3 = new RegExp("pattern");
```

---

## Code Blocks with HTML Entities

```html
&lt;div&gt;Content&lt;/div&gt;
&amp;copy; 2024
&quot;Quoted text&quot;
```

---

## Code Blocks with Escaped Characters

```javascript
const escaped = "String with \\n newline";
const tab = "String with \\t tab";
const quote = "String with \\\" quote";
```

---

## Code Blocks Ending Document

```javascript
function final() {
    return "End of document";
}
```

---

## Code Blocks Starting Document

```javascript
function first() {
    return "Start of document";
}
```

---

## Adjacent Code Blocks

```javascript
// First block
```

```python
# Second block
```

---

## Code Blocks with Only Whitespace

```
    
```

---

## Code Blocks with Only Newlines

```

```

---

## Code Blocks with Mixed Content

```javascript
function mixed() {
    // Comments
    const code = "code";
    /* Multi-line
       comment */
    return code;
}
```

---

## Code Blocks with Very Long Lines

```javascript
const veryLongLine = "This is a very long line of code that extends beyond the normal width and tests how the converter handles long lines in code blocks without breaking the formatting or structure";
```

---

## Code Blocks with Indentation

```python
def nested():
    if True:
        if True:
            return "Deeply nested"
```

---

## Code Blocks with Tabs

```javascript
function withTabs() {
	return "Tab indented";
}
```

---

## Code Blocks with Trailing Spaces

```javascript
function example() {
    return "Line with trailing spaces    ";
}
```

---

## Code Blocks with Special Markdown Characters

```markdown
# This should not be a header
**This should not be bold**
*This should not be italic*
[This should not be a link](https://example.com)
```

---

## Code Blocks Preserving Formatting

```javascript
function formatted() {
    const obj = {
        key1: "value1",
        key2: "value2",
        key3: {
            nested: "value"
        }
    };
    return obj;
}
```

---

## Inline Code in Different Contexts

This paragraph has `code` inline.

**Bold** with `code` inline.

*Italic* with `code` inline.

[Link](https://example.com) with `code` inline.

![Image](image.png) with `code` inline.

---

## Inline Code Escaping

\`Escaped backtick\` should not be code.

`code with \`escaped backtick\`` inside.

---

## Code Blocks with Language Aliases

```js
// JavaScript using 'js' alias
```

```py
# Python using 'py' alias
```

```sh
# Shell using 'sh' alias
```

---

## Code Blocks with Invalid Language

```invalidlanguage
This is code with an invalid language identifier
```

---

## Code Blocks with Empty Language

```
This is a code block with empty language identifier
```

