---
title: तालिकाएं - बुनियादी और उन्नत
categories:
  - markdown
  - tables
tags:
  - grid
  - data
  - formatting
---

# तालिकाएं

## बुनियादी तालिका

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |
| Cell 4   | Cell 5   | Cell 6   |

---

## Table with Header Only

| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|

---

## Table with Alignment

| Left Aligned | Center Aligned | Right Aligned |
|:-------------|:-------------:|--------------:|
| Left         | Center        | Right         |
| Left         | Center        | Right         |

---

## Table with Left Alignment

| Left | Left | Left |
|:-----|:-----|:-----|
| Data | Data | Data |
| Data | Data | Data |

---

## Table with Right Alignment

| Right | Right | Right |
|------:|------:|------:|
| Data  | Data  | Data  |
| Data  | Data  | Data  |

---

## Table with Center Alignment

| Center | Center | Center |
|:------:|:------:|:------:|
| Data   | Data   | Data   |
| Data   | Data   | Data   |

---

## Table with Mixed Alignment

| Left | Center | Right |
|:-----|:------:|------:|
| Data | Data   | Data  |
| Data | Data   | Data  |

---

## Table with Empty Cells

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Cell 1   |          | Cell 3   |
|          | Cell 5   |          |
| Cell 7   | Cell 8   | Cell 9   |

---

## Table with Long Content

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| This is a very long cell that contains multiple words and sentences to test how tables handle longer content that might wrap or extend beyond normal cell widths | Short | Another long cell with extensive content that demonstrates table formatting capabilities |
| Short | Medium length cell | Short |

---

## Table with Special Characters

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| "Quotes" | 'Single' | (Parentheses) |
| [Brackets] | {Braces} | <Tags> |
| $Dollar$ | %Percent% | &Ampersand& |

---

## Table with Numbers

| Number | Decimal | Percentage |
|--------|---------|------------|
| 123    | 123.456 | 50%        |
| 1,000  | 0.001   | 100%       |
| -42    | -3.14   | 0%         |

---

## Table with Links

| Column 1 | Column 2 |
|----------|----------|
| [Link](https://example.com) | [Another Link](https://example.com/2) |
| [Reference Link][ref] | <https://example.com> |

[ref]: https://example.com/reference

---

## Table with Images

| Column 1 | Column 2 |
|----------|----------|
| ![Image](image.png) | ![Another Image](image2.png) |
| ![Reference Image][imgref] | Text |

[imgref]: image.png

---

## Table with Emphasis

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

## Table with HTML

| Column 1 | Column 2 |
|----------|----------|
| <strong>HTML Bold</strong> | <em>HTML Italic</em> |
| <code>HTML Code</code> | <a href="https://example.com">HTML Link</a> |

---

## Table with Unicode

| Column 1 | Column 2 |
|----------|----------|
| café | naïve |
| 你好 | こんにちは |
| مرحبا | Привет |

---

## Table with Multiple Paragraphs (if supported)

| Column 1 | Column 2 |
|----------|----------|
| First paragraph. | Another paragraph. |
| Second paragraph. | More content. |

---

## Table with Line Breaks (if supported)

| Column 1 | Column 2 |
|----------|----------|
| Line 1<br>Line 2 | Line 1<br>Line 2 |

---

## Table Without Pipes

Column 1 | Column 2 | Column 3
----------|----------|----------
Cell 1   | Cell 2   | Cell 3
Cell 4   | Cell 5   | Cell 6

---

## Table with Leading/Trailing Pipes

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |

---

## Table with Extra Spaces

|   Column 1   |   Column 2   |   Column 3   |
|--------------|--------------|--------------|
|   Cell 1     |   Cell 2     |   Cell 3     |
|   Cell 4     |   Cell 5     |   Cell 6     |

---

## Table with Single Column

| Column 1 |
|----------|
| Cell 1   |
| Cell 2   |
| Cell 3   |

---

## Table with Many Columns

| Col 1 | Col 2 | Col 3 | Col 4 | Col 5 | Col 6 | Col 7 | Col 8 |
|-------|-------|-------|-------|-------|-------|-------|-------|
| Data  | Data  | Data  | Data  | Data  | Data  | Data  | Data  |

---

## Table with Many Rows

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

## Table in Lists

* List item before table

| Column 1 | Column 2 |
|----------|----------|
| Data     | Data     |

* List item after table

---

## Table in Blockquotes

> Blockquote before table
> 
> | Column 1 | Column 2 |
> |----------|----------|
> | Data     | Data     |
> 
> Blockquote after table

---

## Table with Headers Only

| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|

---

## Table with No Headers (if supported)

| Cell 1 | Cell 2 | Cell 3 |
|--------|--------|--------|
| Cell 4 | Cell 5 | Cell 6 |

---

## Table with Separator Row Variations

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

## Table with Complex Content

| Column 1 | Column 2 |
|----------|----------|
| **Bold** text with [link](https://example.com) | *Italic* text with `code` |
| ![Image](image.png) with caption | Mixed **bold** and *italic* |
| `Code` with **bold** and [link](https://example.com) | All together |

---

## Table Ending Document

| Column 1 | Column 2 |
|----------|----------|
| Final    | Data     |

---

## Table Starting Document

| Column 1 | Column 2 |
|----------|----------|
| First    | Data     |

---

## Adjacent Tables

| Table 1 Col 1 | Table 1 Col 2 |
|---------------|----------------|
| Data          | Data           |

| Table 2 Col 1 | Table 2 Col 2 |
|---------------|----------------|
| Data          | Data           |

---

## Table with Escaped Pipes

| Column 1 | Column 2 |
|----------|----------|
| Cell with \| pipe | Normal cell |

---

## Table with Empty Rows (if supported)

| Column 1 | Column 2 |
|----------|----------|
|          |          |
| Data     | Data     |

---

## Table with Very Wide Content

| Very Wide Column Header That Extends Beyond Normal Width | Another Wide Column |
|----------------------------------------------------------|---------------------|
| Very wide cell content that extends beyond normal cell width | Normal cell |

---

## Table with Narrow Content

| A | B | C |
|---|---|---|
| 1 | 2 | 3 |
| x | y | z |

---

## Table with Mixed Data Types

| String | Number | Boolean | Date |
|--------|--------|---------|------|
| Text   | 123    | true    | 2024-01-15 |
| More   | 456.78 | false   | 2024-12-31 |

---

## Table with Nested Content (if supported)

| Column 1 | Column 2 |
|----------|----------|
| - List item 1 | 1. Numbered item |
| - List item 2 | 2. Another item |

---

## Table with Code Blocks (if supported)

| Column 1 | Column 2 |
|----------|----------|
| `code` | More `code` |

---

## Table with Horizontal Rules (if supported)

| Column 1 | Column 2 |
|----------|----------|
| --- | --- |

---

## Table with Markdown in Cells

| Column 1 | Column 2 |
|----------|----------|
| **Bold** and *italic* | `code` and [link](https://example.com) |
| # Header | ## Subheader |

---

## Table Preserving Whitespace (if supported)

| Column 1 | Column 2 |
|----------|----------|
|   Indented   |   Content   |
| Normal       | Content     |

