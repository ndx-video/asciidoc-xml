---
title: 表格 - 基础和高级
categories:
  - markdown
  - tables
tags:
  - grid
  - data
  - formatting
---

# 表格

## 基础表格

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |
| Cell 4   | Cell 5   | Cell 6   |

---

## 只有标题的表格

| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|

---

## 带对齐的表格

| Left Aligned | Center Aligned | Right Aligned |
|:-------------|:-------------:|--------------:|
| Left         | Center        | Right         |
| Left         | Center        | Right         |

---

## 左对齐的表格

| Left | Left | Left |
|:-----|:-----|:-----|
| Data | Data | Data |
| Data | Data | Data |

---

## 右对齐的表格

| Right | Right | Right |
|------:|------:|------:|
| Data  | Data  | Data  |
| Data  | Data  | Data  |

---

## 居中对齐的表格

| Center | Center | Center |
|:------:|:------:|:------:|
| Data   | Data   | Data   |
| Data   | Data   | Data   |

---

## 混合对齐的表格

| Left | Center | Right |
|:-----|:------:|------:|
| Data | Data   | Data  |
| Data | Data   | Data  |

---

## 带空单元格的表格

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Cell 1   |          | Cell 3   |
|          | Cell 5   |          |
| Cell 7   | Cell 8   | Cell 9   |

---

## 带长内容的表格

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| This is a very long cell that contains multiple words and sentences to test how tables handle longer content that might wrap or extend beyond normal cell widths | Short | Another long cell with extensive content that demonstrates table formatting capabilities |
| Short | Medium length cell | Short |

---

## 带特殊字符的表格

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| "Quotes" | 'Single' | (Parentheses) |
| [Brackets] | {Braces} | <Tags> |
| $Dollar$ | %Percent% | &Ampersand& |

---

## 带数字的表格

| Number | Decimal | Percentage |
|--------|---------|------------|
| 123    | 123.456 | 50%        |
| 1,000  | 0.001   | 100%       |
| -42    | -3.14   | 0%         |

---

## 带链接的表格

| Column 1 | Column 2 |
|----------|----------|
| [Link](https://example.com) | [Another Link](https://example.com/2) |
| [Reference Link][ref] | <https://example.com> |

[ref]: https://example.com/reference

---

## 带图片的表格

| Column 1 | Column 2 |
|----------|----------|
| ![Image](image.png) | ![Another Image](image2.png) |
| ![Reference Image][imgref] | Text |

[imgref]: image.png

---

## 带强调的表格

| Column 1 | Column 2 |
|----------|----------|
| **Bold** | *Italic* |
| ***Bold and Italic*** | `Code` |
| **Bold** and *Italic* | Mixed formatting |

---

## Table with Code

| Column 1 | Column 2 |
|----------|----------|
| `inline code` | `more code` |
| ``code with `backtick`` | `code123` |

---

## 带 HTML 的表格

| Column 1 | Column 2 |
|----------|----------|
| <strong>HTML Bold</strong> | <em>HTML Italic</em> |
| <code>HTML Code</code> | <a href="https://example.com">HTML Link</a> |

---

## 带 Unicode 的表格

| Column 1 | Column 2 |
|----------|----------|
| café | naïve |
| 你好 | こんにちは |
| مرحبا | Привет |

---

## 带多个段落的表格（如果支持）

| Column 1 | Column 2 |
|----------|----------|
| First paragraph. | Another paragraph. |
| Second paragraph. | More content. |

---

## 带换行的表格（如果支持）

| Column 1 | Column 2 |
|----------|----------|
| Line 1<br>Line 2 | Line 1<br>Line 2 |

---

## 没有管道符的表格

Column 1 | Column 2 | Column 3
----------|----------|----------
Cell 1   | Cell 2   | Cell 3
Cell 4   | Cell 5   | Cell 6

---

## 带前导/尾随管道符的表格

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |

---

## 带额外空格的表格

|   Column 1   |   Column 2   |   Column 3   |
|--------------|--------------|--------------|
|   Cell 1     |   Cell 2     |   Cell 3     |
|   Cell 4     |   Cell 5     |   Cell 6     |

---

## 单列表格

| Column 1 |
|----------|
| Cell 1   |
| Cell 2   |
| Cell 3   |

---

## 多列表格

| Col 1 | Col 2 | Col 3 | Col 4 | Col 5 | Col 6 | Col 7 | Col 8 |
|-------|-------|-------|-------|-------|-------|-------|-------|
| Data  | Data  | Data  | Data  | Data  | Data  | Data  | Data  |

---

## 多行表格

| Column 1 | Column 2 |
|----------|----------|
| Row 1    | Data     |
| Row 2    | Data     |
| Row 3    | Data     |
| Row 4    | Data     |
| Row 5    | Data     |
| Row 6    | Data     |
| Row 7    | Data     |
| Row 8    | Data     |
| Row 9    | Data     |
| Row 10   | Data     |

---

## 列表中的表格

* List item before table

| Column 1 | Column 2 |
|----------|----------|
| Data     | Data     |

* List item after table

---

## 引用块中的表格

> Blockquote before table
> 
> | Column 1 | Column 2 |
> |----------|----------|
> | Data     | Data     |
> 
> Blockquote after table

---

## 只有标题的表格

| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|

---

## 没有标题的表格（如果支持）

| Cell 1 | Cell 2 | Cell 3 |
|--------|--------|--------|
| Cell 4 | Cell 5 | Cell 6 |

---

## 带分隔符行变化的表格

| Column 1 | Column 2 |
|:---------|---------:|
| Left     | Right    |

| Column 1 | Column 2 |
|:---------|:---------|
| Left     | Left     |

| Column 1 | Column 2 |
|---------:|---------:|
| Right    | Right    |

---

## 带复杂内容的表格

| Column 1 | Column 2 |
|----------|----------|
| **Bold** text with [link](https://example.com) | *Italic* text with `code` |
| ![Image](image.png) with caption | Mixed **bold** and *italic* |
| `Code` with **bold** and [link](https://example.com) | All together |

---

## 结束文档的表格

| Column 1 | Column 2 |
|----------|----------|
| Final    | Data     |

---

## 开始文档的表格

| Column 1 | Column 2 |
|----------|----------|
| First    | Data     |

---

## 相邻的表格

| Table 1 Col 1 | Table 1 Col 2 |
|---------------|----------------|
| Data          | Data           |

| Table 2 Col 1 | Table 2 Col 2 |
|---------------|----------------|
| Data          | Data           |

---

## 带转义管道符的表格

| Column 1 | Column 2 |
|----------|----------|
| Cell with \| pipe | Normal cell |

---

## 带空行的表格（如果支持）

| Column 1 | Column 2 |
|----------|----------|
|          |          |
| Data     | Data     |

---

## 带很宽内容的表格

| Very Wide Column Header That Extends Beyond Normal Width | Another Wide Column |
|----------------------------------------------------------|---------------------|
| Very wide cell content that extends beyond normal cell width | Normal cell |

---

## 带窄内容的表格

| A | B | C |
|---|---|---|
| 1 | 2 | 3 |
| x | y | z |

---

## 带混合数据类型的表格

| String | Number | Boolean | Date |
|--------|--------|---------|------|
| Text   | 123    | true    | 2024-01-15 |
| More   | 456.78 | false   | 2024-12-31 |

---

## 带嵌套内容的表格（如果支持）

| Column 1 | Column 2 |
|----------|----------|
| - List item 1 | 1. Numbered item |
| - List item 2 | 2. Another item |

---

## 带代码块的表格（如果支持）

| Column 1 | Column 2 |
|----------|----------|
| `code` | More `code` |

---

## 带水平线的表格（如果支持）

| Column 1 | Column 2 |
|----------|----------|
| --- | --- |

---

## 单元格中带 Markdown 的表格

| Column 1 | Column 2 |
|----------|----------|
| **Bold** and *italic* | `code` and [link](https://example.com) |
| # Header | ## Subheader |

---

## 保留空白的表格（如果支持）

| Column 1 | Column 2 |
|----------|----------|
|   Indented   |   Content   |
| Normal       | Content     |

