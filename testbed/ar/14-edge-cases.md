---
title: حالات الحد والسيناريوهات الخاصة
categories:
  - markdown
  - edge cases
tags:
  - testing
  - special cases
---

# حالات الحد والسيناريوهات الخاصة

## مستند فارغ

---

## مستند يحتوي فقط على مسافات بيضاء

    

---

## مستند يحتوي فقط على أسطر جديدة



---

## مستند يبدأ بعنوان

# First Header

Content.

---

## مستند ينتهي بعنوان

Content.

# Final Header

---

## مستند يحتوي فقط على عنوان

# Only Header

---

## مستند يحتوي فقط على فقرة

This is the only paragraph.

---

## عناوين متجاورة

# Header 1
## Header 2
### Header 3

---

## عناوين بدون محتوى

# Header

## Another Header

---

## فقرات بدون فصل

Paragraph one.
Paragraph two.
Paragraph three.

---

## قوائم بدون أسطر فارغة

* Item 1
* Item 2
# Header
* Item 3

---

## عناصر قائمة فارغة

* 
* Item with content
* 
* Another item

---

## قوائم تحتوي فقط على مسافات

*     
*     

---

## جداول مع خلايا فارغة

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
|          |          |          |
| Data     |          | Data     |

---

## جداول تحتوي فقط على عناوين

| Header 1 | Header 2 |
|----------|----------|

---

## كتل كود تحتوي فقط على مسافات بيضاء

```
    
```

---

## كتل كود تحتوي فقط على أسطر جديدة

```

```

---

## اقتباسات تحتوي فقط على مسافات

>     

---

## اقتباسات تحتوي فقط على أسطر جديدة

>

---

## روابط مع نص فارغ

[](https://example.com)

[ ](https://example.com)

---

## صور مع نص بديل فارغ

![](image.png)

![ ](image.png)

---

## تأكيد مع محتوى فارغ

****

* *

** **

---

## تأكيد مجاور للنص

**Bold**text

*Italic*text

`Code`text

---

## روابط مجاورة للنص

[Link](https://example.com)text

text[Link](https://example.com)

---

## صور مجاورة للنص

![Image](image.png)text

text![Image](image.png)

---

## عناوين تحتوي فقط على مسافات

#     

##     

---

## عناوين تحتوي فقط على أحرف خاصة

# ***

## ---

### ___

---

## قوائم مع علامات مختلطة

* Item 1
- Item 2
+ Item 3

---

## قوائم تبدأ بأرقام مختلفة

5. Item 5
10. Item 10
1. Item 1

---

## جداول مع فواصل غير متسقة

| Column 1 | Column 2 |
|:---------|---------|
| Data     | Data    |

---

## كتل كود بأطوال سياج مختلفة

````
Code with four backticks
````

`````
Code with five backticks
`````

---

## هياكل متداخلة بعمق أقصى

> Blockquote
> > Nested
> > > Deep nested
> > > > Very deep

---

## حالات حد تنسيق مختلط

**Bold*Italic**Bold*

*Italic**Bold*Italic*

---

## روابط في مواضع مختلفة

[Start](https://example.com) text.

Text [middle](https://example.com) text.

Text [end](https://example.com).

---

## صور في مواضع مختلفة

![Start](image.png) text.

Text ![middle](image.png) text.

Text ![end](image.png).

---

## كود في مواضع مختلفة

`Start` text.

Text `middle` text.

Text `end`.

---

## تأكيد في مواضع مختلفة

**Start** text.

Text **middle** text.

Text **end**.

---

## أحرف خاصة عند الحدود

*Item with "quotes"*

*Item with (parentheses)*

*Item with [brackets]*

---

## حالات حد Unicode

café naïve résumé

你好世界

こんにちは世界

مرحبا بالعالم

Привет мир

---

## أرقام في سياقات مختلفة

123

1,000

1.234

-42

+100

---

## عناوين URL بصيغ مختلفة

https://example.com

http://example.com

ftp://ftp.example.com

file:///path/to/file

---

## تنسيقات البريد الإلكتروني

user@example.com

user.name@example.com

user+tag@example.com

user_name@example.com

---

## مسارات الملفات

/path/to/file

./relative/path

../parent/path

~/home/path

C:\Windows\Path

---

## أرقام الإصدار

v1.0.0

1.2.3

1.2.3-beta.1

1.2.3-alpha.2+build.123

---

## التواريخ

2024-01-15

01/15/2024

January 15, 2024

2024-01-15T10:30:00Z

---

## الأوقات

10:30 AM

22:30

2h30m

5s

---

## النسب المئوية

50%

100%

0%

123.45%

---

## العملة

$19.99

€15.50

£10.00

¥1000

---

## عناوين IP

192.168.1.1

127.0.0.1

2001:0db8:85a3:0000:0000:8a2e:0370:7334

---

## علامات الهاشتاغ

#hashtag

#hashtag-with-dashes

#hashtag_with_underscores

#123hashtag

---

## الإشارات

@username

@user_name

@user-name

@user123

---

## أحرف خاصة مختلطة

$100 (50% off) @user #tag

Version 1.2.3 released on 2024-01-15

Email: user@example.com, Phone: +1-234-567-8900

---

## Frontmatter فارغ

---
---

---

## Frontmatter يحتوي فقط على مسافات

---
    
---

---

## Frontmatter يحتوي فقط على أسطر جديدة

---


---

---

## كتل Frontmatter متعددة

---
title: First
---

---
title: Second
---

---

## Frontmatter بدون فواصل

title: No separators
author: Author Name

---

## محتوى يبدأ فورًا

Content without frontmatter or header.

---

## محتوى ينتهي فجأة

Content that ends without newline.

