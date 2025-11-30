---
title: Lists - Ordered, Unordered, and Nested
categories:
  - markdown
  - lists
tags:
  - ordered
  - unordered
  - nested
  - task lists
---

# Lists

## Unordered Lists - Asterisk

* First item
* Second item
* Third item

---

## Unordered Lists - Hyphen

- First item
- Second item
- Third item

---

## Unordered Lists - Plus

+ First item
+ Second item
+ Third item

---

## Ordered Lists - Numbers

1. First item
2. Second item
3. Third item

---

## Ordered Lists - Sequential Numbers

1. First item
2. Second item
3. Third item
4. Fourth item
5. Fifth item

---

## Ordered Lists - Non-Sequential Numbers

1. First item
5. Second item (numbered 5)
3. Third item (numbered 3)
10. Fourth item (numbered 10)

---

## Ordered Lists Starting from Different Numbers

5. Fifth item
6. Sixth item
7. Seventh item

---

## Nested Unordered Lists

* First level item
  * Second level item
    * Third level item
  * Another second level item
* Another first level item

---

## Nested Ordered Lists

1. First level item
   1. Second level item
      1. Third level item
   2. Another second level item
2. Another first level item

---

## Mixed Nested Lists

* Unordered item
  1. Ordered nested item
  2. Another ordered nested item
* Another unordered item

1. Ordered item
   * Unordered nested item
   * Another unordered nested item
2. Another ordered item

---

## Lists with Multiple Paragraphs

* First item

  This is a second paragraph in the first item.

* Second item

  This is a second paragraph in the second item.

  This is a third paragraph in the second item.

---

## Lists with Code Blocks

* First item

  ```javascript
  console.log("Code block in list item");
  ```

* Second item

---

## Lists with Blockquotes

* First item

  > This is a blockquote in a list item

* Second item

---

## Lists with Emphasis

* Item with **bold** text
* Item with *italic* text
* Item with ***bold and italic*** text
* Item with `code` inline

---

## Lists with Links

* Item with [link](https://example.com)
* Item with [reference link][ref]
* Item with [link](https://example.com "title")

[ref]: https://example.com/reference

---

## Lists with Images

* Item with ![image](image.png)
* Item with ![image with alt](image.png "title")
* Item with ![reference image][imgref]

[imgref]: image.png "Reference image"

---

## Task Lists (GitHub Flavored Markdown)

- [ ] Unchecked task
- [x] Checked task
- [X] Checked task (uppercase)
- [ ] Another unchecked task
- [x] Another checked task

---

## Nested Task Lists

- [ ] Main task
  - [ ] Subtask 1
  - [x] Subtask 2 (completed)
  - [ ] Subtask 3
- [x] Another main task (completed)
  - [x] Completed subtask
  - [ ] Incomplete subtask

---

## Task Lists with Content

- [ ] Unchecked task with **bold** text
- [x] Checked task with *italic* text
- [ ] Task with [link](https://example.com)
- [x] Task with `code` inline

---

## Lists with HTML

* Item with <strong>HTML bold</strong>
* Item with <em>HTML italic</em>
* Item with <code>HTML code</code>

---

## Lists with Special Characters

* Item with "quotes"
* Item with (parentheses)
* Item with [brackets]
* Item with {braces}
* Item with <tags>
* Item with $dollar$
* Item with &ampersand&

---

## Lists with Numbers

* Item with 123 numbers
* Item with v1.2.3 version
* Item with 50% percentage
* Item with 1,000,000 large number

---

## Lists with Unicode

* Item with café
* Item with naïve
* Item with 你好
* Item with こんにちは
* Item with مرحبا

---

## Tight Lists (No Blank Lines)

* First item
* Second item
* Third item

1. First item
2. Second item
3. Third item

---

## Loose Lists (With Blank Lines)

* First item

* Second item

* Third item

1. First item

2. Second item

3. Third item

---

## Lists Starting Mid-Document

Some paragraph text before the list.

* List item 1
* List item 2

Some paragraph text after the list.

---

## Lists Ending Document

* Final item 1
* Final item 2
* Final item 3

---

## Empty List Items

* 
* Item with content
* 
* Another item

---

## Lists with Long Content

* This is a very long list item that contains multiple sentences and continues for a while to test how the converter handles longer content within list items. It should maintain proper formatting and structure.

* Another long item that spans multiple lines in the source but should be rendered as a single paragraph within the list item structure.

---

## Complex Nested Lists

1. First ordered item
   * Nested unordered item
     - Deep nested item
       + Deeper nested item
   * Another nested unordered item
2. Second ordered item
   * Nested unordered item
     1. Nested ordered item
     2. Another nested ordered item
3. Third ordered item

---

## Lists with Horizontal Rules

* Item before rule

---

* Item after rule

---

## Lists with Headers

* Item before header

## Header in List Context

* Item after header

---

## Definition Lists (if supported)

Term 1
: Definition 1

Term 2
: Definition 2
: Alternative definition 2

Term 3
: Definition 3 with **bold** and *italic*

