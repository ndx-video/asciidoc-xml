---
title: 引用块 - 基础和嵌套
categories:
  - markdown
  - blockquotes
tags:
  - quotes
  - citations
  - nested
---

# 引用块

## 基础引用块

> This is a basic blockquote.
> It spans multiple lines.
> Each line starts with a > character.

---

## 单行引用块

> This is a single line blockquote.

---

## 带段落的引用块

> This is the first paragraph in a blockquote.
> 
> This is the second paragraph in the same blockquote.

---

## 带强调的引用块

> This blockquote contains **bold** text.

> This blockquote contains *italic* text.

> This blockquote contains ***bold and italic*** text.

> This blockquote contains `code` inline.

---

## 带链接的引用块

> This blockquote contains a [link](https://example.com).

> This blockquote contains a [reference link][blockquoteref].

[blockquoteref]: https://example.com/blockquote

---

## Blockquote with Images

> This blockquote contains an ![image](image.png).

> This blockquote contains a ![reference image][blockquoteimg].

[blockquoteimg]: image.png

---

## 带列表的引用块

> This blockquote contains a list:
> - Item 1
> - Item 2
> - Item 3

> This blockquote contains a numbered list:
> 1. First item
> 2. Second item
> 3. Third item

---

## 带代码块的引用块

> This blockquote contains a code block:
> 
> ```javascript
> console.log("Code in blockquote");
> ```

> This blockquote contains an indented code block:
> 
>     const code = "indented";
>     console.log(code);

---

## Blockquote with Headers

> # Header in Blockquote
> 
> ## Subheader in Blockquote
> 
> Content under headers.

---

## 带水平线的引用块

> Content before rule
> 
> ---
> 
> Content after rule

---

## Blockquote with Tables

> This blockquote contains a table:
> 
> | Column 1 | Column 2 |
> |----------|----------|
> | Data     | Data     |

---

## Nested Blockquotes

> This is the first level blockquote.
> 
> > This is the second level blockquote.
> 
> > > This is the third level blockquote.

---

## 带内容的嵌套引用块

> Outer blockquote content.
> 
> > Inner blockquote content.
> > 
> > > Deepest blockquote content.

---

## 带混合内容的引用块

> This blockquote has **bold**, *italic*, `code`, [links](https://example.com), and ![images](image.png).
> 
> It also has:
> - Lists
> - Multiple items
> 
> And code blocks:
> 
> ```javascript
> console.log("Example");
> ```

---

## 带特殊字符的引用块

> Blockquote with "quotes" and 'single quotes'.

> Blockquote with (parentheses) and [brackets].

> Blockquote with {braces} and <tags>.

> Blockquote with $dollar$ and %percent%.

---

## 带数字的引用块

> Blockquote with 123 numbers.

> Blockquote with v1.2.3 version.

> Blockquote with 50% percentage.

---

## 带 Unicode 的引用块

> Blockquote with café and naïve.

> Blockquote with 你好 and こんにちは.

> Blockquote with مرحبا and Привет.

---

## Blockquote with HTML

> Blockquote with <strong>HTML bold</strong>.

> Blockquote with <em>HTML italic</em>.

> Blockquote with <code>HTML code</code>.

---

## 带转义字符的引用块

> Blockquote with \*escaped asterisk\*.

> Blockquote with \_escaped underscore\_.

> Blockquote with \[escaped brackets\].

---

## 开始文档的引用块

> This blockquote starts the document.

---

## 结束文档的引用块

> This blockquote ends the document.

---

## 多个引用块

> First blockquote.

> Second blockquote.

> Third blockquote.

---

## 带归属的引用块（如果支持）

> This is a quote.
> 
> — Author Name

> This is another quote.
> 
> — Another Author

---

## 带引用的引用块（如果支持）

> This is a quote.
> 
> <cite>Source Name</cite>

---

## 带长内容的引用块

> This is a very long blockquote that contains multiple sentences and continues for a while to test how blockquotes handle longer content. It should maintain proper formatting and structure regardless of the length of the content within the blockquote. This is useful for testing how the converter handles blockquotes that might be reformatted or wrapped differently in the output format.

---

## 带换行的引用块

> This blockquote has a line break  
> using two spaces.

> This blockquote uses a backslash\
> for a line break.

---

## 带空行的引用块

> Content before empty line.
> 
> 
> Content after empty lines.

---

## 只有空格的引用块

>     
> 

---

## 只有换行的引用块

>


---

## 列表中的引用块

* List item with blockquote:
  
  > Blockquote in list item

* Another list item

---

## 带嵌套列表的引用块

> Blockquote with nested list:
> 
> 1. First item
>    * Nested item
>    * Another nested item
> 2. Second item

---

## 带复杂嵌套的引用块

> Outer blockquote.
> 
> > Inner blockquote with list:
> > - Item 1
> > - Item 2
> > 
> > > Deep blockquote with code:
> > > ```javascript
> > > console.log("Nested");
> > > ```

---

## 带多个段落和列表的引用块

> First paragraph in blockquote.
> 
> Second paragraph in blockquote.
> 
> List in blockquote:
> - Item 1
> - Item 2
> 
> Third paragraph after list.

---

## 带代码和列表的引用块

> Blockquote with code:
> 
> ```javascript
> const example = "code";
> ```
> 
> And a list:
> - Item 1
> - Item 2

---

## 带表格和列表的引用块

> Blockquote with table:
> 
> | Column 1 | Column 2 |
> |----------|----------|
> | Data     | Data     |
> 
> And a list:
> - Item 1
> - Item 2

---

## 带标题和内容的引用块

> # Main Header
> 
> Content under header.
> 
> ## Subheader
> 
> More content.

---

## 保留格式的引用块

> This blockquote preserves:
> 
> - List formatting
> - **Bold** and *italic*
> - `Code` blocks
> - [Links](https://example.com)

---

## 带 GitHub 风格 Markdown 功能的引用块

> Blockquote with ~~strikethrough~~.

> Blockquote with - [ ] task list:
> - [ ] Task 1
> - [x] Task 2

---

## 带自动链接的引用块

> Blockquote with <https://example.com>.

> Blockquote with <mailto:user@example.com>.

---

## 带引用链接的引用块

> Blockquote with [reference link][ref1] and [another][ref2].

[ref1]: https://example.com/1
[ref2]: https://example.com/2

---

## 带图片和链接的引用块

> Blockquote with ![image](image.png) and [link](https://example.com).

> Blockquote with [![linked image](image.png)](https://example.com).

---

## 带混合内联元素的引用块

> Blockquote with **bold**, *italic*, `code`, [link](https://example.com), and ![image](image.png) all together.

---

## 带 HTML 实体的引用块

> Blockquote with &lt;tags&gt; and &quot;quotes&quot;.

> Blockquote with &amp; ampersand.

---

## 带特殊标点的引用块

> Blockquote ending with period.

> Blockquote ending with exclamation!

> Blockquote ending with question?

> Blockquote with ellipsis...

> Blockquote with em-dash—and en-dash–.

---

## 带很长单词的引用块

> Blockquote with verylongwordthatextendsbeyondnormalwidth and anotherverylongword.

---

## 保留空白的引用块（如果支持）

>     Indented content
> Normal content

---

## 带混合语言的引用块

> Blockquote with English, français, español, and 中文.

---

## 带数学表达式的引用块（如果支持）

> Blockquote with $x = y + z$ inline math.

> Blockquote with block math:
> 
> $$
> x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}
> $$

