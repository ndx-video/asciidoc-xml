---
title: 代码块和内联代码
categories:
  - markdown
  - code
tags:
  - programming
  - syntax highlighting
  - code blocks
---

# 代码块和内联代码

## 内联代码

这个段落包含使用反引号的 `内联代码`。

这个段落中有 `带空格的代码`。

这个段落有 `带下划线的代码` 和 `带破折号的代码`。

这个段落有带数字的 `代码123`。

这个段落有带点的 `代码.带.点`。

---

## 带特殊字符的内联代码

这个段落内部有 `带 "引号" 的代码`。

这个段落内部有 `带（括号）的代码`。

这个段落内部有 `带 [方括号] 的代码`。

这个段落内部有 `带 {大括号} 的代码`。

这个段落内部有 `带 <标签> 的代码`。

这个段落内部有 `带 $美元$ 的代码`。

这个段落内部有 `带 &和号& 的代码`。

---

## 带反引号的内联代码

这个段落使用双反引号有 ``内部带 `反引号` 的代码``。

这个段落使用三个反引号有 ```内部带 ``双反引号`` 的代码```。

---

## 内联代码边缘情况

`代码`文本（无空格）

`代码`**粗体**（相邻格式）

`代码`*斜体*（相邻格式）

`代码`[链接](https://example.com)（相邻链接）

---

## 围栏代码块 - 无语言

```
这是一个代码块
有多行
代码
```

---

## 围栏代码块 - JavaScript

```javascript
function greet(name) {
    console.log("Hello, " + name + "!");
}

greet("World");
```

---

## 围栏代码块 - Python

```python
def greet(name):
    print(f"Hello, {name}!")

greet("World")
```

---

## 围栏代码块 - Go

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

---

## 围栏代码块 - HTML

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

## 围栏代码块 - CSS

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

## 围栏代码块 - JSON

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

## 围栏代码块 - YAML

```yaml
name: Example
version: 1.0.0
dependencies:
  package: ^1.2.3
```

---

## 围栏代码块 - Shell

```bash
#!/bin/bash
echo "Hello, World!"
ls -la
```

---

## 围栏代码块 - SQL

```sql
SELECT * FROM users
WHERE age > 18
ORDER BY name;
```

---

## 围栏代码块 - XML

```xml
<?xml version="1.0" encoding="UTF-8"?>
<root>
    <element attribute="value">Content</element>
</root>
```

---

## 围栏代码块 - Markdown

```markdown
# Header

**Bold** and *italic* text.

- List item
```

---

## 围栏代码块 - 多种语言

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

## 缩进代码块

    This is an indented code block
    with multiple lines
    using four spaces

---

## 带空行的代码块

```javascript
function example() {
    // First line
    
    // Third line after empty line
}
```

---

## 带特殊字符的代码块

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

## 带 Unicode 的代码块

```javascript
const unicode = {
    "café": "café",
    "你好": "你好",
    "こんにちは": "こんにちは"
};
```

---

## 内部带反引号的代码块

````markdown
```
Code block with backticks
```
````

---

## 内部带三个反引号的代码块

`````markdown
```
Code block content
```
`````

---

## 列表中的代码块

* 带代码块的列表项：

  ```javascript
  console.log("Code in list");
  ```

* 另一个项目

---

## 引用块中的代码块

> 带代码块的引用块：
> 
> ```javascript
> console.log("Code in blockquote");
> ```

---

## 带行号的代码块（如果支持）

```javascript:1-5
function example() {
    return "Hello";
}
```

---

## 带文件名的代码块（如果支持）

```javascript:example.js
function example() {
    return "Hello";
}
```

---

## 带元数据的代码块

```javascript
// File: example.js
// Author: Test Author
// Date: 2024-01-15
function example() {
    return "Hello";
}
```

---

## 带注释的代码块

```javascript
// Single line comment
function example() {
    /* Multi-line
       comment */
    return "Hello";
}
```

---

## 带字符串的代码块

```javascript
const single = 'Single quoted string';
const double = "Double quoted string";
const template = `Template string with ${variable}`;
```

---

## 带正则表达式的代码块

```javascript
const regex1 = /pattern/g;
const regex2 = /pattern with \(escaped\)/g;
const regex3 = new RegExp("pattern");
```

---

## 带 HTML 实体的代码块

```html
&lt;div&gt;Content&lt;/div&gt;
&amp;copy; 2024
&quot;Quoted text&quot;
```

---

## 带转义字符的代码块

```javascript
const escaped = "String with \\n newline";
const tab = "String with \\t tab";
const quote = "String with \\\" quote";
```

---

## 结束文档的代码块

```javascript
function final() {
    return "End of document";
}
```

---

## 开始文档的代码块

```javascript
function first() {
    return "Start of document";
}
```

---

## 相邻的代码块

```javascript
// First block
```

```python
# Second block
```

---

## 只有空白的代码块

```
    
```

---

## 只有换行的代码块

```

```

---

## 带混合内容的代码块

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

## 带很长行的代码块

```javascript
const veryLongLine = "This is a very long line of code that extends beyond the normal width and tests how the converter handles long lines in code blocks without breaking the formatting or structure";
```

---

## 带缩进的代码块

```python
def nested():
    if True:
        if True:
            return "Deeply nested"
```

---

## 带制表符的代码块

```javascript
function withTabs() {
	return "Tab indented";
}
```

---

## 带尾随空格的代码块

```javascript
function example() {
    return "Line with trailing spaces    ";
}
```

---

## 带特殊 Markdown 字符的代码块

```markdown
# This should not be a header
**This should not be bold**
*This should not be italic*
[This should not be a link](https://example.com)
```

---

## 保留格式的代码块

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

## 不同上下文中的内联代码

这个段落有 `代码` 内联。

**粗体** 与 `代码` 内联。

*斜体* 与 `代码` 内联。

[链接](https://example.com) 与 `代码` 内联。

![图片](image.png) 与 `代码` 内联。

---

## 内联代码转义

\`转义的反引号\` 不应该是代码。

内部有 `带 \`转义反引号\` 的代码`。

---

## 带语言别名的代码块

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

## 带无效语言的代码块

```invalidlanguage
This is code with an invalid language identifier
```

---

## 带空语言的代码块

```
This is a code block with empty language identifier
```

