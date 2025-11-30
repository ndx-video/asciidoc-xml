---
title: الأحرف المُهربة والحالات الخاصة
categories:
  - markdown
  - escaping
tags:
  - special characters
  - edge cases
---

# الأحرف المُهربة والحالات الخاصة

## تهريب النجمات

\*This should not be italic\*

\*\*This should not be bold\*\*

\*\*\*This should not be bold and italic\*\*\*

---

## تهريب الشرطات السفلية

\_This should not be italic\_

\_\_This should not be bold\_\_

\_\_\_This should not be bold and italic\_\_\_

---

## تهريب علامات الاقتباس المقلوبة

\`This should not be code\`

\`\`This should not be code\`\`

---

## تهريب الأقواس المربعة

\[This should not be a link\](https://example.com)

\!\[This should not be an image\](image.png)

---

## تهريب علامات الهاشتاغ

\# This should not be a header

\## This should not be a header

---

## تهريب علامات الجمع

\+ This should not be a list item

---

## تهريب علامات الطرح

\- This should not be a list item

---

## تهريب النقاط

\. This should not be a numbered list

---

## تهريب علامة أكبر من

\> This should not be a blockquote

---

## تهريب الأنابيب

\| This should not be a table cell

---

## تهريب التيلدا

\~\~This should not be strikethrough\~\~

---

## تهريب الأقواس

\(Parentheses\)

\[Brackets\]

\{Braces\}

---

## تهريب الأحرف الخاصة

\* Asterisk

\_ Underscore

\` Backtick

\# Hashtag

\+ Plus

\- Minus

\. Period

\> Greater than

\| Pipe

\~ Tilde

---

## التهريب في سياقات مختلفة

### في الفقرات

\*Escaped asterisk\* in paragraph.

### في القوائم

* Item with \*escaped asterisk\*
* Item with \_escaped underscore\_

### في الاقتباسات

> Blockquote with \*escaped asterisk\*

### في كتل الكود

```
\*Escaped asterisk\* in code block
```

---

## تهريب عناوين URL

\https://example.com should not be a link

---

## تهريب البريد الإلكتروني

\user@example.com should not be a link

---

## تهريب HTML

\&lt;div&gt; should display as &lt;div&gt;

\&amp; should display as &

---

## تهريب أحرف متعددة

\*\*Bold\*\* and \*italic\* escaped.

---

## التهريب في بداية السطر

\# Not a header

\* Not a list

\- Not a list

---

## التهريب في نهاية السطر

Text with \*escaped asterisk\*

---

## التهريب مع المسافات

\* Escaped with space

\*Escaped without space\*

---

## تهريب التسلسلات الخاصة

\*\*\* Triple asterisk escaped

\_\_\_ Triple underscore escaped

\`\`\` Triple backtick escaped

---

## التهريب في الروابط

[Link with \*asterisk\*](https://example.com)

[Link with \_underscore\_](https://example.com)

---

## التهريب في الصور

![Image with \*asterisk\*](image.png)

![Image with \_underscore\_](image.png)

---

## التهريب في العناوين

# Header with \*asterisk\*

## Header with \_underscore\_

---

## التهريب في الجداول

| Column 1 | Column 2 |
|----------|----------|
| \*Escaped\* | \_Escaped\_ |

---

## تهريب الشرطة المائلة للخلف

\\ Backslash escaped

\\\\ Double backslash

---

## تهريب المجموعات

\*\*Bold\*\* and \*italic\* and \`code\` all escaped.

---

## تهريب Unicode

\你好 should not be processed

\café should not be processed

---

## تهريب الأرقام

\123 should not be processed

\1.2.3 should not be processed

---

## تهريب علامات الاقتباس

\"Double quotes\"

\'Single quotes\'

---

## تهريب علامات الدولار

\$Dollar sign\$

---

## تهريب علامات النسبة المئوية

\%Percent sign\%

---

## تهريب علامات العطف

\&Ampersand\&

---

## تهريب الأقواس الزاوية

\<Less than\>

\>Greater than\>

---

## تهريب الأقواس المربعة

\[Left bracket\]

\]Right bracket\]

---

## تهريب الأقواس المعقوفة

\{Left brace\}

\}Right brace\}

---

## تهريب الأقواس

\(Left parenthesis\)

\)Right parenthesis\)

---

## تهريب علامات التعجب

\!Exclamation mark\!

---

## تهريب علامات الاستفهام

\?Question mark\?

---

## تهريب النقطتين

\:Colon\:

---

## تهريب الفاصلة المنقوطة

\;Semicolon\;

---

## تهريب الفواصل

\,Comma\,

---

## تهريب النقاط

\.Period\.

---

## تهريب أحرف خاصة متعددة

\*\*Bold\*\* \*italic\* \`code\` \[link\] all escaped.

---

## التهريب في مواضع مختلفة

Start: \*escaped\*

Middle: text \*escaped\* text

End: text \*escaped\*

---

## التهريب مع النص المجاور

\*escaped\*text

text\*escaped\*

\*escaped\*text\*escaped\*

---

## التهريب مع الحفاظ على المعنى الحرفي

\* should display as *

\_ should display as _

\` should display as `

\# should display as #

---

## حالات حد التهريب

\*\* should display as **

\*\*\* should display as ***

\_\_ should display as __

\_\_\_ should display as ___

