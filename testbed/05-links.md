---
title: Links - Inline, Reference, and Autolinks
categories:
  - markdown
  - links
tags:
  - hyperlinks
  - references
  - urls
---

# Links

## Inline Links

This is an [inline link](https://example.com).

This is a [link with title](https://example.com "Link Title").

This is a [link with single quotes](https://example.com 'Link Title').

This is a [link with parentheses](https://example.com (Link Title)).

---

## Links with Different Protocols

[HTTP link](http://example.com)

[HTTPS link](https://example.com)

[FTP link](ftp://ftp.example.com)

[Email link](mailto:user@example.com)

[File link](file:///path/to/file)

---

## Links with Paths

[Root path](https://example.com/)

[Subdirectory](https://example.com/path/to/resource)

[File extension](https://example.com/file.html)

[Query string](https://example.com/page?param=value&other=123)

[Fragment](https://example.com/page#section)

[Complex URL](https://example.com/path/to/resource?query=value&other=123#fragment)

---

## Links with Special Characters

[Link with spaces](https://example.com/path with spaces)

[Link with unicode](https://example.com/你好)

[Link with encoded](https://example.com/path%20with%20spaces)

---

## Reference-Style Links

This is a [reference link][ref1].

This is an [implicit reference link][].

This is a [shortcut reference link].

[ref1]: https://example.com/reference1
[implicit reference link]: https://example.com/implicit
[shortcut reference link]: https://example.com/shortcut

---

## Reference Links with Titles

This is a [reference link with title][ref2].

[ref2]: https://example.com/reference2 "Reference Title"

This is a [reference link with single quotes][ref3].

[ref3]: https://example.com/reference3 'Reference Title'

This is a [reference link with parentheses][ref4].

[ref4]: https://example.com/reference4 (Reference Title)

---

## Reference Links - Case Sensitivity

[Case sensitive link][case1]

[Case Sensitive Link][Case1]

[case1]: https://example.com/lowercase
[Case1]: https://example.com/uppercase

---

## Reference Links - Numbers

[Link 1][1]

[Link 2][2]

[1]: https://example.com/1
[2]: https://example.com/2

---

## Reference Links - Special Characters

[Link with dash][link-dash]

[Link with underscore][link_underscore]

[Link with dot][link.dot]

[link-dash]: https://example.com/dash
[link_underscore]: https://example.com/underscore
[link.dot]: https://example.com/dot

---

## Autolinks

<https://example.com>

<http://example.com>

<mailto:user@example.com>

<user@example.com>

---

## Links with Emphasis

**[Bold link](https://example.com)**

*[Italic link](https://example.com)*

***[Bold and italic link](https://example.com)***

[**Bold text**](https://example.com) in link

[*Italic text*](https://example.com) in link

---

## Links in Lists

* Item with [link](https://example.com)
* Item with [reference link][listref]
* Item with <https://example.com>

[listref]: https://example.com/list

---

## Links in Blockquotes

> This is a [link](https://example.com) in a blockquote.

> This is a [reference link][blockquoteref] in a blockquote.

[blockquoteref]: https://example.com/blockquote

---

## Links in Headers

# Header with [link](https://example.com)

## Header with [reference link][headerref]

[headerref]: https://example.com/header

---

## Links in Tables

| Column 1 | Column 2 |
|----------|----------|
| [Link](https://example.com) | [Reference][tableref] |
| <https://example.com> | [Link with title](https://example.com "Title") |

[tableref]: https://example.com/table

---

## Links with Images

[![Image](image.png)](https://example.com)

[![Image with alt](image.png "Title")](https://example.com "Link Title")

---

## Relative Links

[Relative link](./relative.html)

[Parent directory](../parent.html)

[Sibling file](../sibling.html)

[Root relative](/root.html)

---

## Anchor Links

[Link to section](#section)

[Link to subsection](#subsection-name)

[Link with spaces](#section name)

[Link with special chars](#section-name_123)

---

## Links with Code

[Link with `code`](https://example.com)

`[Code link](https://example.com)` (should not be a link)

---

## Links with HTML

<a href="https://example.com">HTML link</a>

[Markdown link](https://example.com) and <a href="https://example.com">HTML link</a>

---

## Broken Links

[Broken link](https://nonexistent.example.com)

[Broken reference][broken]

---

## Links with Escaped Characters

\[Escaped link\]\(https://example.com\)

[Link with \*asterisk\*](https://example.com)

[Link with \_underscore\_](https://example.com)

---

## Links with Unicode

[Link with café](https://example.com/café)

[Link with 你好](https://example.com/你好)

[Link with こんにちは](https://example.com/こんにちは)

---

## Multiple Links in One Paragraph

This paragraph has [first link](https://example.com/1), [second link](https://example.com/2), and [third link](https://example.com/3).

This paragraph mixes [inline link](https://example.com/inline), [reference link][multiref], and <https://example.com/autolink>.

[multiref]: https://example.com/reference

---

## Links with Long URLs

[Link with very long URL](https://example.com/very/long/path/to/resource/that/spans/multiple/segments/and/includes/query/parameters?param1=value1&param2=value2&param3=value3#and-even-a-fragment)

---

## Links with Empty Text

[](https://example.com)

[ ](https://example.com)

---

## Links with Only Spaces

[   ](https://example.com)

---

## Links in Code Blocks

```
[Link](https://example.com) should not work in code blocks
```

---

## Links with Special URL Schemes

[javascript:alert('XSS')](javascript:alert('XSS'))

[data:text/plain,Hello](data:text/plain,Hello)

[about:blank](about:blank)

---

## Links with Percent Encoding

[Link with encoding](https://example.com/path%20with%20spaces)

[Link with unicode encoding](https://example.com/%E4%BD%A0%E5%A5%BD)

---

## Complex Link Combinations

**[Bold link](https://example.com)** with *italic* text and [another link](https://example.com/2).

[Link with **bold** and *italic*](https://example.com) in the text.

[Link](https://example.com) with `code` and **bold** and *italic*.

