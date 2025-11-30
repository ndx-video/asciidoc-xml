---
title: 链接 - 内联、引用和自动链接
categories:
  - markdown
  - links
tags:
  - hyperlinks
  - references
  - urls
---

# 链接

## 内联链接

这是一个 [内联链接](https://example.com)。

这是一个 [带标题的链接](https://example.com "链接标题")。

这是一个 [带单引号的链接](https://example.com '链接标题')。

这是一个 [带括号的链接](https://example.com (链接标题))。

---

## 不同协议的链接

[HTTP 链接](http://example.com)

[HTTPS 链接](https://example.com)

[FTP 链接](ftp://ftp.example.com)

[电子邮件链接](mailto:user@example.com)

[文件链接](file:///path/to/file)

---

## 带路径的链接

[根路径](https://example.com/)

[子目录](https://example.com/path/to/resource)

[文件扩展名](https://example.com/file.html)

[查询字符串](https://example.com/page?param=value&other=123)

[片段](https://example.com/page#section)

[复杂 URL](https://example.com/path/to/resource?query=value&other=123#fragment)

---

## 带特殊字符的链接

[带空格的链接](https://example.com/path with spaces)

[带 unicode 的链接](https://example.com/你好)

[带编码的链接](https://example.com/path%20with%20spaces)

---

## 引用样式链接

这是一个 [引用链接][ref1]。

这是一个 [隐式引用链接][]。

这是一个 [快捷引用链接]。

[ref1]: https://example.com/reference1
[implicit reference link]: https://example.com/implicit
[shortcut reference link]: https://example.com/shortcut

---

## 带标题的引用链接

这是一个 [带标题的引用链接][ref2]。

[ref2]: https://example.com/reference2 "引用标题"

这是一个 [带单引号的引用链接][ref3]。

[ref3]: https://example.com/reference3 '引用标题'

这是一个 [带括号的引用链接][ref4]。

[ref4]: https://example.com/reference4 (引用标题)

---

## 引用链接 - 大小写敏感性

[区分大小写的链接][case1]

[区分大小写的链接][Case1]

[case1]: https://example.com/lowercase
[Case1]: https://example.com/uppercase

---

## 引用链接 - 数字

[链接 1][1]

[链接 2][2]

[1]: https://example.com/1
[2]: https://example.com/2

---

## 引用链接 - 特殊字符

[带破折号的链接][link-dash]

[带下划线的链接][link_underscore]

[带点的链接][link.dot]

[link-dash]: https://example.com/dash
[link_underscore]: https://example.com/underscore
[link.dot]: https://example.com/dot

---

## 自动链接

<https://example.com>

<http://example.com>

<mailto:user@example.com>

<user@example.com>

---

## 带强调的链接

**[粗体链接](https://example.com)**

*[斜体链接](https://example.com)*

***[粗体和斜体链接](https://example.com)***

[**粗体文本**](https://example.com) 在链接中

[*斜体文本*](https://example.com) 在链接中

---

## 列表中的链接

* 带 [链接](https://example.com) 的项目
* 带 [引用链接][listref] 的项目
* 带 <https://example.com> 的项目

[listref]: https://example.com/list

---

## 引用块中的链接

> 这是引用块中的 [链接](https://example.com)。

> 这是引用块中的 [引用链接][blockquoteref]。

[blockquoteref]: https://example.com/blockquote

---

## 标题中的链接

# 带 [链接](https://example.com) 的标题

## 带 [引用链接][headerref] 的标题

[headerref]: https://example.com/header

---

## 表格中的链接

| Column 1 | Column 2 |
|----------|----------|
| [链接](https://example.com) | [引用][tableref] |
| <https://example.com> | [带标题的链接](https://example.com "标题") |

[tableref]: https://example.com/table

---

## 带图片的链接

[![图片](image.png)](https://example.com)

[![带替代文本的图片](image.png "标题")](https://example.com "链接标题")

---

## 相对链接

[相对链接](./relative.html)

[父目录](../parent.html)

[同级文件](../sibling.html)

[根相对](/root.html)

---

## 锚点链接

[链接到章节](#section)

[链接到子章节](#subsection-name)

[带空格的链接](#section name)

[带特殊字符的链接](#section-name_123)

---

## 带代码的链接

[带 `代码` 的链接](https://example.com)

`[代码链接](https://example.com)`（不应该是链接）

---

## 带 HTML 的链接

<a href="https://example.com">HTML 链接</a>

[Markdown 链接](https://example.com) 和 <a href="https://example.com">HTML 链接</a>

---

## 断开的链接

[断开的链接](https://nonexistent.example.com)

[断开的引用][broken]

---

## 带转义字符的链接

\[转义的链接\]\(https://example.com\)

[带 \*星号\* 的链接](https://example.com)

[带 \_下划线\_ 的链接](https://example.com)

---

## 带 Unicode 的链接

[带 café 的链接](https://example.com/café)

[带 你好 的链接](https://example.com/你好)

[带 こんにちは 的链接](https://example.com/こんにちは)

---

## 一个段落中的多个链接

这个段落有 [第一个链接](https://example.com/1)、[第二个链接](https://example.com/2) 和 [第三个链接](https://example.com/3)。

这个段落混合了 [内联链接](https://example.com/inline)、[引用链接][multiref] 和 <https://example.com/autolink>。

[multiref]: https://example.com/reference

---

## 带长 URL 的链接

[带很长 URL 的链接](https://example.com/very/long/path/to/resource/that/spans/multiple/segments/and/includes/query/parameters?param1=value1&param2=value2&param3=value3#and-even-a-fragment)

---

## 带空文本的链接

[](https://example.com)

[ ](https://example.com)

---

## 只有空格的链接

[   ](https://example.com)

---

## 代码块中的链接

```
[链接](https://example.com) 在代码块中不应起作用
```

---

## 带特殊 URL 方案的链接

[javascript:alert('XSS')](javascript:alert('XSS'))

[data:text/plain,Hello](data:text/plain,Hello)

[about:blank](about:blank)

---

## 带百分比编码的链接

[带编码的链接](https://example.com/path%20with%20spaces)

[带 unicode 编码的链接](https://example.com/%E4%BD%A0%E5%A5%BD)

---

## 复杂的链接组合

**[粗体链接](https://example.com)** 与 *斜体* 文本和 [另一个链接](https://example.com/2)。

[带 **粗体** 和 *斜体* 的链接](https://example.com) 在文本中。

[链接](https://example.com) 与 `代码` 和 **粗体** 和 *斜体*。

