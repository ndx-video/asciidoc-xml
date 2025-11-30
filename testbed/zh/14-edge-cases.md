---
title: 边缘情况和特殊场景
categories:
  - markdown
  - edge cases
tags:
  - testing
  - special cases
---

# 边缘情况和特殊场景

## 空文档

---

## 只有空白的文档

    

---

## 只有换行的文档



---

## 以标题开始的文档

# First Header

Content.

---

## 以标题结束的文档

Content.

# Final Header

---

## 只有标题的文档

# Only Header

---

## 只有段落的文档

This is the only paragraph.

---

## 相邻的标题

# Header 1
## Header 2
### Header 3

---

## 没有内容的标题

# Header

## Another Header

---

## 没有分隔的段落

Paragraph one.
Paragraph two.
Paragraph three.

---

## 没有空行的列表

* Item 1
* Item 2
# Header
* Item 3

---

## 空列表项

* 
* Item with content
* 
* Another item

---

## 只有空格的列表

*     
*     

---

## 带空单元格的表格

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
|          |          |          |
| Data     |          | Data     |

---

## 只有标题的表格

| Header 1 | Header 2 |
|----------|----------|

---

## 只有空白的代码块

```
    
```

---

## 只有换行的代码块

```

```

---

## 只有空白的引用块

>     

---

## 只有换行的引用块

>

---

## 空文本的链接

[](https://example.com)

[ ](https://example.com)

---

## 空替代文本的图片

![](image.png)

![ ](image.png)

---

## 空内容的强调

****

* *

** **

---

## 与文本相邻的强调

**Bold**text

*Italic*text

`Code`text

---

## 与文本相邻的链接

[Link](https://example.com)text

text[Link](https://example.com)

---

## 与文本相邻的图片

![Image](image.png)text

text![Image](image.png)

---

## 只有空格的标题

#     

##     

---

## 只有特殊字符的标题

# ***

## ---

### ___

---

## 混合标记的列表

* Item 1
- Item 2
+ Item 3

---

## 以不同数字开始的列表

5. Item 5
10. Item 10
1. Item 1

---

## 不一致分隔符的表格

| Column 1 | Column 2 |
|:---------|---------|
| Data     | Data    |

---

## 不同围栏长度的代码块

````
Code with four backticks
````

`````
Code with five backticks
`````

---

## 嵌套结构最大深度

> Blockquote
> > Nested
> > > Deep nested
> > > > Very deep

---

## 混合格式边缘情况

**Bold*Italic**Bold*

*Italic**Bold*Italic*

---

## 不同位置的链接

[Start](https://example.com) text.

Text [middle](https://example.com) text.

Text [end](https://example.com).

---

## 不同位置的图片

![Start](image.png) text.

Text ![middle](image.png) text.

Text ![end](image.png).

---

## 不同位置的代码

`Start` text.

Text `middle` text.

Text `end`.

---

## 不同位置的强调

**Start** text.

Text **middle** text.

Text **end**.

---

## 边界处的特殊字符

*Item with "quotes"*

*Item with (parentheses)*

*Item with [brackets]*

---

## Unicode 边缘情况

café naïve résumé

你好世界

こんにちは世界

مرحبا بالعالم

Привет мир

---

## 不同上下文中的数字

123

1,000

1.234

-42

+100

---

## 不同格式的 URL

https://example.com

http://example.com

ftp://ftp.example.com

file:///path/to/file

---

## 电子邮件格式

user@example.com

user.name@example.com

user+tag@example.com

user_name@example.com

---

## 文件路径

/path/to/file

./relative/path

../parent/path

~/home/path

C:\Windows\Path

---

## 版本号

v1.0.0

1.2.3

1.2.3-beta.1

1.2.3-alpha.2+build.123

---

## 日期

2024-01-15

01/15/2024

January 15, 2024

2024-01-15T10:30:00Z

---

## 时间

10:30 AM

22:30

2h30m

5s

---

## 百分比

50%

100%

0%

123.45%

---

## 货币

$19.99

€15.50

£10.00

¥1000

---

## IP 地址

192.168.1.1

127.0.0.1

2001:0db8:85a3:0000:0000:8a2e:0370:7334

---

## 话题标签

#hashtag

#hashtag-with-dashes

#hashtag_with_underscores

#123hashtag

---

## 提及

@username

@user_name

@user-name

@user123

---

## 混合特殊字符

$100 (50% off) @user #tag

Version 1.2.3 released on 2024-01-15

Email: user@example.com, Phone: +1-234-567-8900

---

## 空 Frontmatter

---
---

---

## 只有空格的 Frontmatter

---
    
---

---

## 只有换行的 Frontmatter

---


---

---

## 多个 Frontmatter 块

---
title: First
---

---
title: Second
---

---

## 没有分隔符的 Frontmatter

title: No separators
author: Author Name

---

## 立即开始的内容

Content without frontmatter or header.

---

## 突然结束的内容

Content that ends without newline.

