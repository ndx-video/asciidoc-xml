---
title: Markdown 中的 HTML
categories:
  - markdown
  - html
tags:
  - mixed
  - html
---

# Markdown 中的 HTML

## 内联 HTML 标签

这个段落包含 <strong>HTML 粗体</strong> 文本。

这个段落包含 <em>HTML 斜体</em> 文本。

这个段落包含 <code>HTML 代码</code> 文本。

---

## HTML 链接

这个段落有 <a href="https://example.com">HTML 链接</a>。

这个段落有 <a href="https://example.com" title="Link Title">带标题的 HTML 链接</a>。

这个段落有 <a href="https://example.com" target="_blank">带目标的 HTML 链接</a>。

---

## HTML 图片

这个段落有 <img src="image.png" alt="HTML image">。

这个段落有 <img src="image.png" alt="HTML image" width="100" height="100">。

这个段落有 <img src="image.png" alt="HTML image" class="custom-class">。

---

## HTML 换行

这个段落有<br>一个换行。

这个段落有<br/>一个换行（自闭合）。

这个段落有<br />一个换行（带空格）。

---

## HTML 段落

<p>这是一个 HTML 段落。</p>

<p>这是另一个 HTML 段落。</p>

---

## HTML 标题

<h1>HTML 标题 1</h1>

<h2>HTML 标题 2</h2>

<h3>HTML 标题 3</h3>

---

## HTML 列表

<ul>
<li>项目 1</li>
<li>项目 2</li>
<li>项目 3</li>
</ul>

<ol>
<li>第一</li>
<li>第二</li>
<li>第三</li>
</ol>

---

## HTML 引用块

<blockquote>
这是一个 HTML 引用块。
</blockquote>

---

## HTML 代码块

<pre><code>function example() {
    return "code";
}</code></pre>

---

## HTML 表格

<table>
<tr>
<th>Header 1</th>
<th>Header 2</th>
</tr>
<tr>
<td>Data 1</td>
<td>Data 2</td>
</tr>
</table>

---

## 带属性的 HTML

<div class="container" id="main" data-value="123">
Content in div
</div>

<span style="color: red;">Styled span</span>

---

## HTML 注释

<!-- This is an HTML comment -->

注释后的内容。

---

## HTML 实体

&lt;div&gt;Content&lt;/div&gt;

&amp;copy; 2024

&quot;Quoted text&quot;

&apos;Single quoted&apos;

---

## HTML 与 Markdown 混合

这个段落有 **Markdown 粗体** 和 <strong>HTML 粗体</strong>。

这个段落有 *Markdown 斜体* 和 <em>HTML 斜体</em>。

这个段落有 `Markdown 代码` 和 <code>HTML 代码</code>。

---

## HTML 和 Markdown 链接

这个段落有 [Markdown 链接](https://example.com) 和 <a href="https://example.com">HTML 链接</a>。

---

## HTML 和 Markdown 图片

这个段落有 ![Markdown 图片](image.png) 和 <img src="image.png" alt="HTML image">。

---

## 列表中的 HTML

* 带 <strong>HTML 粗体</strong> 的项目
* 带 <em>HTML 斜体</em> 的项目
* 带 <code>HTML 代码</code> 的项目

---

## 引用块中的 HTML

> 带 <strong>HTML 粗体</strong> 的引用块
> 
> 带 <em>HTML 斜体</em> 的引用块

---

## 表格中的 HTML

| Column 1 | Column 2 |
|----------|----------|
| <strong>HTML bold</strong> | <em>HTML italic</em> |
| <code>HTML code</code> | <a href="https://example.com">HTML link</a> |

---

## HTML Script 标签

<script>
console.log("Script content");
</script>

---

## HTML Style 标签

<style>
body { font-family: Arial; }
</style>

---

## HTML 表单元素

<form>
<input type="text" name="field">
<button type="submit">Submit</button>
</form>

---

## HTML Details/Summary

<details>
<summary>Click to expand</summary>
Hidden content here.
</details>

---

## HTML 缩写

<abbr title="HyperText Markup Language">HTML</abbr>

---

## HTML 引用

<cite>Citation text</cite>

---

## HTML 定义

<dfn>Definition term</dfn>

---

## HTML 键盘输入

按 <kbd>Ctrl</kbd> + <kbd>C</kbd> 复制。

---

## HTML 标记文本

<mark>Highlighted text</mark>

---

## HTML 小文本

<small>Small text</small>

---

## HTML 下标和上标

H<sub>2</sub>O

E = mc<sup>2</sup>

---

## HTML 时间

<time datetime="2024-01-15">January 15, 2024</time>

---

## HTML 变量

<var>variable_name</var>

---

## HTML 示例输出

<samp>Sample output</samp>

---

## HTML 地址

<address>
123 Main St<br>
City, State 12345
</address>

---

## HTML 预格式化文本

<pre>
Preformatted
    text
        with
            indentation
</pre>

---

## 带嵌套标签的 HTML

<div>
<p>Paragraph in div</p>
<ul>
<li>List in div</li>
</ul>
</div>

---

## HTML 自闭合标签

<hr>

<hr/>

<hr />

<br>

<br/>

<br />

<img src="image.png" alt="Image">

<img src="image.png" alt="Image"/>

---

## 带数据属性的 HTML

<div data-id="123" data-name="example">
Content
</div>

---

## 带 ARIA 属性的 HTML

<button aria-label="Close">×</button>

<div role="button" tabindex="0">Clickable div</div>

---

## 带类和 ID 的 HTML

<div class="container main" id="content">
Content
</div>

---

## 带内联样式的 HTML

<div style="color: red; font-size: 18px;">
Styled content
</div>

---

## HTML 转义字符

&lt;div&gt; should display as &lt;div&gt;

&amp; should display as &

&quot;quotes&quot; should display as "quotes"

---

## HTML 混合内容

<div>
**Markdown bold** and <strong>HTML bold</strong>
*Markdown italic* and <em>HTML italic</em>
`Markdown code` and <code>HTML code</code>
</div>

---

## HTML 保留空白

<pre>
    Indented
        content
</pre>

---

## 带 Unicode 的 HTML

<div>
Content with café and 你好
</div>

---

## 内容中的 HTML 注释

注释前的内容 <!-- comment --> 注释后的内容。

---

## HTML 空标签

<div></div>

<span></span>

<p></p>

---

## 属性中包含引号的 HTML

<div title="Title with 'single quotes'">
Content
</div>

<div title='Title with "double quotes"'>
Content
</div>

---

## 属性中包含特殊字符的 HTML

<div data-value="value with spaces">
Content
</div>

<div data-value="value-with-dashes">
Content
</div>

<div data-value="value_with_underscores">
Content
</div>

---

## HTML 块级 vs 内联

<div>Block element</div>

<span>Inline element</span>

<p>Block paragraph</p>

<strong>Inline bold</strong>

---

## 正确嵌套的 HTML

<div>
<p>Paragraph</p>
<ul>
<li>Item</li>
</ul>
</div>

---

## HTML 与 Markdown 格式混合

**Bold** 与 <strong>HTML bold</strong> 和 *italic* 与 <em>HTML italic</em>。

---

## 代码块中的 HTML（应该是字面量）

```
<div>HTML in code block</div>
```

---

## Markdown 中的 HTML 转义

\&lt;div&gt; escaped HTML

---

## 带 Markdown 链接的 HTML

[Link](https://example.com) 和 <a href="https://example.com">HTML link</a> 一起。

---

## 带 Markdown 图片的 HTML

![Image](image.png) 和 <img src="image.png" alt="HTML image"> 一起。

