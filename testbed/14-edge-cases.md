---
title: Edge Cases and Special Scenarios
categories:
  - markdown
  - edge cases
tags:
  - testing
  - special cases
---

# Edge Cases and Special Scenarios

## Empty Document

---

## Document with Only Whitespace

    

---

## Document with Only Newlines



---

## Document Starting with Header

# First Header

Content.

---

## Document Ending with Header

Content.

# Final Header

---

## Document with Only Header

# Only Header

---

## Document with Only Paragraph

This is the only paragraph.

---

## Adjacent Headers

# Header 1
## Header 2
### Header 3

---

## Headers Without Content

# Header

## Another Header

---

## Paragraphs Without Separation

Paragraph one.
Paragraph two.
Paragraph three.

---

## Lists Without Blank Lines

* Item 1
* Item 2
# Header
* Item 3

---

## Empty List Items

* 
* Item with content
* 
* Another item

---

## Lists with Only Spaces

*     
*     

---

## Tables with Empty Cells

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
|          |          |          |
| Data     |          | Data     |

---

## Tables with Only Headers

| Header 1 | Header 2 |
|----------|----------|

---

## Code Blocks with Only Whitespace

```
    
```

---

## Code Blocks with Only Newlines

```

```

---

## Blockquotes with Only Spaces

>     

---

## Blockquotes with Only Newlines

>

---

## Links with Empty Text

[](https://example.com)

[ ](https://example.com)

---

## Images with Empty Alt Text

![](image.png)

![ ](image.png)

---

## Emphasis with Empty Content

****

* *

** **

---

## Emphasis Adjacent to Text

**Bold**text

*Italic*text

`Code`text

---

## Links Adjacent to Text

[Link](https://example.com)text

text[Link](https://example.com)

---

## Images Adjacent to Text

![Image](image.png)text

text![Image](image.png)

---

## Headers with Only Spaces

#     

##     

---

## Headers with Only Special Characters

# ***

## ---

### ___

---

## Lists with Mixed Markers

* Item 1
- Item 2
+ Item 3

---

## Lists Starting with Different Numbers

5. Item 5
10. Item 10
1. Item 1

---

## Tables with Inconsistent Separators

| Column 1 | Column 2 |
|:---------|---------|
| Data     | Data    |

---

## Code Blocks with Different Fence Lengths

````
Code with four backticks
````

`````
Code with five backticks
`````

---

## Nested Structures Maximum Depth

> Blockquote
> > Nested
> > > Deep nested
> > > > Very deep

---

## Mixed Formatting Edge Cases

**Bold*Italic**Bold*

*Italic**Bold*Italic*

---

## Links in Different Positions

[Start](https://example.com) text.

Text [middle](https://example.com) text.

Text [end](https://example.com).

---

## Images in Different Positions

![Start](image.png) text.

Text ![middle](image.png) text.

Text ![end](image.png).

---

## Code in Different Positions

`Start` text.

Text `middle` text.

Text `end`.

---

## Emphasis in Different Positions

**Start** text.

Text **middle** text.

Text **end**.

---

## Special Characters at Boundaries

*Item with "quotes"*

*Item with (parentheses)*

*Item with [brackets]*

---

## Unicode Edge Cases

café naïve résumé

你好世界

こんにちは世界

مرحبا بالعالم

Привет мир

---

## Numbers in Different Contexts

123

1,000

1.234

-42

+100

---

## URLs in Different Formats

https://example.com

http://example.com

ftp://ftp.example.com

file:///path/to/file

---

## Email Formats

user@example.com

user.name@example.com

user+tag@example.com

user_name@example.com

---

## File Paths

/path/to/file

./relative/path

../parent/path

~/home/path

C:\Windows\Path

---

## Version Numbers

v1.0.0

1.2.3

1.2.3-beta.1

1.2.3-alpha.2+build.123

---

## Dates

2024-01-15

01/15/2024

January 15, 2024

2024-01-15T10:30:00Z

---

## Times

10:30 AM

22:30

2h30m

5s

---

## Percentages

50%

100%

0%

123.45%

---

## Currency

$19.99

€15.50

£10.00

¥1000

---

## IP Addresses

192.168.1.1

127.0.0.1

2001:0db8:85a3:0000:0000:8a2e:0370:7334

---

## Hashtags

#hashtag

#hashtag-with-dashes

#hashtag_with_underscores

#123hashtag

---

## Mentions

@username

@user_name

@user-name

@user123

---

## Mixed Special Characters

$100 (50% off) @user #tag

Version 1.2.3 released on 2024-01-15

Email: user@example.com, Phone: +1-234-567-8900

---

## Empty Frontmatter

---
---

---

## Frontmatter with Only Spaces

---
    
---

---

## Frontmatter with Only Newlines

---


---

---

## Multiple Frontmatter Blocks

---
title: First
---

---
title: Second
---

---

## Frontmatter Without Separators

title: No separators
author: Author Name

---

## Content Starting Immediately

Content without frontmatter or header.

---

## Content Ending Abruptly

Content that ends without newline.

