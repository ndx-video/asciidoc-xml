---
title: किनारे के मामले और विशेष परिदृश्य
categories:
  - markdown
  - edge cases
tags:
  - testing
  - special cases
---

# किनारे के मामले और विशेष परिदृश्य

## खाली दस्तावेज़

---

## केवल रिक्त स्थान वाला दस्तावेज़

    

---

## केवल नई पंक्तियों वाला दस्तावेज़



---

## शीर्षक से शुरू होने वाला दस्तावेज़

# First Header

Content.

---

## शीर्षक के साथ समाप्त होने वाला दस्तावेज़

Content.

# Final Header

---

## केवल शीर्षक वाला दस्तावेज़

# Only Header

---

## केवल अनुच्छेद वाला दस्तावेज़

This is the only paragraph.

---

## आसन्न शीर्षक

# Header 1
## Header 2
### Header 3

---

## बिना सामग्री के शीर्षक

# Header

## Another Header

---

## बिना अलगाव के अनुच्छेद

Paragraph one.
Paragraph two.
Paragraph three.

---

## बिना रिक्त पंक्तियों की सूचियां

* Item 1
* Item 2
# Header
* Item 3

---

## खाली सूची आइटम

* 
* Item with content
* 
* Another item

---

## केवल रिक्त स्थान वाली सूचियां

*     
*     

---

## खाली कोशिकाओं वाली तालिकाएं

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
|          |          |          |
| Data     |          | Data     |

---

## केवल शीर्षक वाली तालिकाएं

| Header 1 | Header 2 |
|----------|----------|

---

## केवल रिक्त स्थान वाले कोड ब्लॉक

```
    
```

---

## केवल नई पंक्तियों वाले कोड ब्लॉक

```

```

---

## केवल रिक्त स्थान वाले ब्लॉककोट्स

>     

---

## केवल नई पंक्तियों वाले ब्लॉककोट्स

>

---

## खाली पाठ वाले लिंक

[](https://example.com)

[ ](https://example.com)

---

## खाली वैकल्पिक पाठ वाली छवियां

![](image.png)

![ ](image.png)

---

## खाली सामग्री वाला जोर

****

* *

** **

---

## पाठ के आसन्न जोर

**Bold**text

*Italic*text

`Code`text

---

## पाठ के आसन्न लिंक

[Link](https://example.com)text

text[Link](https://example.com)

---

## पाठ के आसन्न छवियां

![Image](image.png)text

text![Image](image.png)

---

## केवल रिक्त स्थान वाले शीर्षक

#     

##     

---

## केवल विशेष वर्ण वाले शीर्षक

# ***

## ---

### ___

---

## मिश्रित मार्कर वाली सूचियां

* Item 1
- Item 2
+ Item 3

---

## विभिन्न संख्याओं से शुरू होने वाली सूचियां

5. Item 5
10. Item 10
1. Item 1

---

## असंगत विभाजक वाली तालिकाएं

| Column 1 | Column 2 |
|:---------|---------|
| Data     | Data    |

---

## विभिन्न बाड़ लंबाई वाले कोड ब्लॉक

````
Code with four backticks
````

`````
Code with five backticks
`````

---

## नेस्टेड संरचनाएं अधिकतम गहराई

> Blockquote
> > Nested
> > > Deep nested
> > > > Very deep

---

## मिश्रित प्रारूप किनारे के मामले

**Bold*Italic**Bold*

*Italic**Bold*Italic*

---

## विभिन्न स्थितियों में लिंक

[Start](https://example.com) text.

Text [middle](https://example.com) text.

Text [end](https://example.com).

---

## विभिन्न स्थितियों में छवियां

![Start](image.png) text.

Text ![middle](image.png) text.

Text ![end](image.png).

---

## विभिन्न स्थितियों में कोड

`Start` text.

Text `middle` text.

Text `end`.

---

## विभिन्न स्थितियों में जोर

**Start** text.

Text **middle** text.

Text **end**.

---

## सीमाओं पर विशेष वर्ण

*Item with "quotes"*

*Item with (parentheses)*

*Item with [brackets]*

---

## Unicode किनारे के मामले

café naïve résumé

你好世界

こんにちは世界

مرحبا بالعالم

Привет мир

---

## विभिन्न संदर्भों में संख्याएं

123

1,000

1.234

-42

+100

---

## विभिन्न प्रारूपों में URL

https://example.com

http://example.com

ftp://ftp.example.com

file:///path/to/file

---

## ईमेल प्रारूप

user@example.com

user.name@example.com

user+tag@example.com

user_name@example.com

---

## फ़ाइल पथ

/path/to/file

./relative/path

../parent/path

~/home/path

C:\Windows\Path

---

## संस्करण संख्याएं

v1.0.0

1.2.3

1.2.3-beta.1

1.2.3-alpha.2+build.123

---

## तारीखें

2024-01-15

01/15/2024

January 15, 2024

2024-01-15T10:30:00Z

---

## समय

10:30 AM

22:30

2h30m

5s

---

## प्रतिशत

50%

100%

0%

123.45%

---

## मुद्रा

$19.99

€15.50

£10.00

¥1000

---

## IP पते

192.168.1.1

127.0.0.1

2001:0db8:85a3:0000:0000:8a2e:0370:7334

---

## हैशटैग

#hashtag

#hashtag-with-dashes

#hashtag_with_underscores

#123hashtag

---

## उल्लेख

@username

@user_name

@user-name

@user123

---

## मिश्रित विशेष वर्ण

$100 (50% off) @user #tag

Version 1.2.3 released on 2024-01-15

Email: user@example.com, Phone: +1-234-567-8900

---

## खाली Frontmatter

---
---

---

## केवल रिक्त स्थान वाला Frontmatter

---
    
---

---

## केवल नई पंक्तियों वाला Frontmatter

---


---

---

## कई Frontmatter ब्लॉक

---
title: First
---

---
title: Second
---

---

## बिना विभाजक के Frontmatter

title: No separators
author: Author Name

---

## तुरंत शुरू होने वाली सामग्री

Content without frontmatter or header.

---

## अचानक समाप्त होने वाली सामग्री

Content that ends without newline.

