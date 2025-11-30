---
title: कोड ब्लॉक और इनलाइन कोड
categories:
  - markdown
  - code
tags:
  - programming
  - syntax highlighting
  - code blocks
---

# कोड ब्लॉक और इनलाइन कोड

## इनलाइन कोड

इस अनुच्छेद में `inline code` बैकटिक का उपयोग करता है।

इस अनुच्छेद में `code with spaces` है।

इस अनुच्छेद में `code_with_underscores` और `code-with-dashes` है।

इस अनुच्छेद में संख्याओं के साथ `code123` है।

इस अनुच्छेद में बिंदुओं के साथ `code.with.dots` है।

---

## विशेष वर्णों के साथ इनलाइन कोड

इस अनुच्छेद के अंदर `code with "quotes"` है।

इस अनुच्छेद के अंदर `code with (parentheses)` है।

इस अनुच्छेद के अंदर `code with [brackets]` है।

इस अनुच्छेद के अंदर `code with {braces}` है।

इस अनुच्छेद के अंदर `code with <tags>` है।

इस अनुच्छेद के अंदर `code with $dollar$` है।

इस अनुच्छेद के अंदर `code with &ampersand&` है।

---

## बैकटिक के साथ इनलाइन कोड

इस अनुच्छेद में ``code with `backtick` inside`` दोहरे बैकटिक का उपयोग करता है।

इस अनुच्छेद में ```code with ``double backticks`` inside``` तिगुने बैकटिक का उपयोग करता है।

---

## इनलाइन कोड के किनारे के मामले

`code`text (रिक्त स्थान नहीं)

`code`**bold** (आसन्न प्रारूप)

`code`*italic* (आसन्न प्रारूप)

`code`[link](https://example.com) (आसन्न लिंक)

---

## बाड़ के साथ कोड ब्लॉक - बिना भाषा

```
This is a code block
with multiple lines
of code
```

---

## बाड़ के साथ कोड ब्लॉक - JavaScript

```javascript
function greet(name) {
    console.log("Hello, " + name + "!");
}

greet("World");
```

---

## बाड़ के साथ कोड ब्लॉक - Python

```python
def greet(name):
    print(f"Hello, {name}!")

greet("World")
```

---

## बाड़ के साथ कोड ब्लॉक - Go

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

---

## बाड़ के साथ कोड ब्लॉक - HTML

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

## बाड़ के साथ कोड ब्लॉक - CSS

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

## बाड़ के साथ कोड ब्लॉक - JSON

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

## बाड़ के साथ कोड ब्लॉक - YAML

```yaml
name: Example
version: 1.0.0
dependencies:
  package: ^1.2.3
```

---

## बाड़ के साथ कोड ब्लॉक - Shell

```bash
#!/bin/bash
echo "Hello, World!"
ls -la
```

---

## बाड़ के साथ कोड ब्लॉक - SQL

```sql
SELECT * FROM users
WHERE age > 18
ORDER BY name;
```

---

## बाड़ के साथ कोड ब्लॉक - XML

```xml
<?xml version="1.0" encoding="UTF-8"?>
<root>
    <element attribute="value">Content</element>
</root>
```

---

## बाड़ के साथ कोड ब्लॉक - Markdown

```markdown
# Header

**Bold** and *italic* text.

- List item
```

---

## बाड़ के साथ कोड ब्लॉक - कई भाषाएं

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

## सिंडेंटेशन के साथ कोड ब्लॉक

    This is an indented code block
    with multiple lines
    using four spaces

---

## रिक्त पंक्तियों के साथ कोड ब्लॉक

```javascript
function example() {
    // First line
    
    // Third line after empty line
}
```

---

## विशेष वर्णों के साथ कोड ब्लॉक

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

## Unicode के साथ कोड ब्लॉक

```javascript
const unicode = {
    "café": "café",
    "你好": "你好",
    "こんにちは": "こんにちは"
};
```

---

## अंदर बैकटिक के साथ कोड ब्लॉक

````markdown
```
Code block with backticks
```
````

---

## अंदर तीन बैकटिक के साथ कोड ब्लॉक

`````markdown
```
Code block content
```
`````

---

## सूचियों में कोड ब्लॉक

* कोड ब्लॉक के साथ सूची आइटम:

  ```javascript
  console.log("Code in list");
  ```

* दूसरा आइटम

---

## ब्लॉककोट्स में कोड ब्लॉक

> कोड ब्लॉक के साथ ब्लॉककोट:
> 
> ```javascript
> console.log("Code in blockquote");
> ```

---

## लाइन नंबर के साथ कोड ब्लॉक (यदि समर्थित)

```javascript:1-5
function example() {
    return "Hello";
}
```

---

## फ़ाइल नामों के साथ कोड ब्लॉक (यदि समर्थित)

```javascript:example.js
function example() {
    return "Hello";
}
```

---

## मेटाडेटा के साथ कोड ब्लॉक

```javascript
// File: example.js
// Author: Test Author
// Date: 2024-01-15
function example() {
    return "Hello";
}
```

---

## टिप्पणियों के साथ कोड ब्लॉक

```javascript
// Single line comment
function example() {
    /* Multi-line
       comment */
    return "Hello";
}
```

---

## स्ट्रिंग्स के साथ कोड ब्लॉक

```javascript
const single = 'Single quoted string';
const double = "Double quoted string";
const template = `Template string with ${variable}`;
```

---

## नियमित अभिव्यक्तियों के साथ कोड ब्लॉक

```javascript
const regex1 = /pattern/g;
const regex2 = /pattern with \(escaped\)/g;
const regex3 = new RegExp("pattern");
```

---

## HTML इकाइयों के साथ कोड ब्लॉक

```html
&lt;div&gt;Content&lt;/div&gt;
&amp;copy; 2024
&quot;Quoted text&quot;
```

---

## एस्केप किए गए वर्णों के साथ कोड ब्लॉक

```javascript
const escaped = "String with \\n newline";
const tab = "String with \\t tab";
const quote = "String with \\\" quote";
```

---

## दस्तावेज़ समाप्त करने वाले कोड ब्लॉक

```javascript
function final() {
    return "End of document";
}
```

---

## दस्तावेज़ शुरू करने वाले कोड ब्लॉक

```javascript
function first() {
    return "Start of document";
}
```

---

## आसन्न कोड ब्लॉक

```javascript
// First block
```

```python
# Second block
```

---

## केवल रिक्त स्थान वाले कोड ब्लॉक

```
    
```

---

## केवल नई पंक्तियों वाले कोड ब्लॉक

```

```

---

## मिश्रित सामग्री वाले कोड ब्लॉक

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

## बहुत लंबी पंक्तियों वाले कोड ब्लॉक

```javascript
const veryLongLine = "This is a very long line of code that extends beyond the normal width and tests how the converter handles long lines in code blocks without breaking the formatting or structure";
```

---

## सिंडेंटेशन के साथ कोड ब्लॉक

```python
def nested():
    if True:
        if True:
            return "Deeply nested"
```

---

## टैब के साथ कोड ब्लॉक

```javascript
function withTabs() {
	return "Tab indented";
}
```

---

## अंतिम रिक्त स्थान वाले कोड ब्लॉक

```javascript
function example() {
    return "Line with trailing spaces    ";
}
```

---

## विशेष Markdown वर्णों के साथ कोड ब्लॉक

```markdown
# This should not be a header
**This should not be bold**
*This should not be italic*
[This should not be a link](https://example.com)
```

---

## प्रारूप संरक्षित करने वाले कोड ब्लॉक

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

## विभिन्न संदर्भों में इनलाइन कोड

इस अनुच्छेद में `code` इनलाइन है।

**Bold** के साथ `code` इनलाइन।

*Italic* के साथ `code` इनलाइन।

[Link](https://example.com) के साथ `code` इनलाइन।

![Image](image.png) के साथ `code` इनलाइन।

---

## इनलाइन कोड एस्केप

\`Escaped backtick\` कोड नहीं होना चाहिए।

`code with \`escaped backtick\`` अंदर।

---

## भाषा उपनामों के साथ कोड ब्लॉक

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

## अमान्य भाषा के साथ कोड ब्लॉक

```invalidlanguage
This is code with an invalid language identifier
```

---

## खाली भाषा के साथ कोड ब्लॉक

```
This is a code block with empty language identifier
```

