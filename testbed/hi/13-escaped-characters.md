---
title: एस्केप किए गए वर्ण और विशेष मामले
categories:
  - markdown
  - escaping
tags:
  - special characters
  - edge cases
---

# एस्केप किए गए वर्ण और विशेष मामले

## तारांकन एस्केप करना

\*This should not be italic\*

\*\*This should not be bold\*\*

\*\*\*This should not be bold and italic\*\*\*

---

## अंडरस्कोर एस्केप करना

\_This should not be italic\_

\_\_This should not be bold\_\_

\_\_\_This should not be bold and italic\_\_\_

---

## बैकटिक एस्केप करना

\`This should not be code\`

\`\`This should not be code\`\`

---

## ब्रैकेट एस्केप करना

\[This should not be a link\](https://example.com)

\!\[This should not be an image\](image.png)

---

## हैशटैग एस्केप करना

\# This should not be a header

\## This should not be a header

---

## प्लस साइन एस्केप करना

\+ This should not be a list item

---

## माइनस साइन एस्केप करना

\- This should not be a list item

---

## पीरियड एस्केप करना

\. This should not be a numbered list

---

## ग्रेटर दैन एस्केप करना

\> This should not be a blockquote

---

## पाइप एस्केप करना

\| This should not be a table cell

---

## टिल्ड एस्केप करना

\~\~This should not be strikethrough\~\~

---

## कोष्ठक एस्केप करना

\(Parentheses\)

\[Brackets\]

\{Braces\}

---

## विशेष वर्ण एस्केप करना

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

## विभिन्न संदर्भों में एस्केप करना

### अनुच्छेदों में

\*Escaped asterisk\* in paragraph.

### सूचियों में

* Item with \*escaped asterisk\*
* Item with \_escaped underscore\_

### ब्लॉककोट्स में

> Blockquote with \*escaped asterisk\*

### कोड ब्लॉक में

```
\*Escaped asterisk\* in code block
```

---

## URL एस्केप करना

\https://example.com should not be a link

---

## ईमेल एस्केप करना

\user@example.com should not be a link

---

## HTML एस्केप करना

\&lt;div&gt; should display as &lt;div&gt;

\&amp; should display as &

---

## कई वर्ण एस्केप करना

\*\*Bold\*\* and \*italic\* escaped.

---

## लाइन की शुरुआत में एस्केप करना

\# Not a header

\* Not a list

\- Not a list

---

## लाइन के अंत में एस्केप करना

Text with \*escaped asterisk\*

---

## रिक्त स्थान के साथ एस्केप करना

\* Escaped with space

\*Escaped without space\*

---

## विशेष अनुक्रम एस्केप करना

\*\*\* Triple asterisk escaped

\_\_\_ Triple underscore escaped

\`\`\` Triple backtick escaped

---

## लिंक में एस्केप करना

[Link with \*asterisk\*](https://example.com)

[Link with \_underscore\_](https://example.com)

---

## छवियों में एस्केप करना

![Image with \*asterisk\*](image.png)

![Image with \_underscore\_](image.png)

---

## शीर्षकों में एस्केप करना

# Header with \*asterisk\*

## Header with \_underscore\_

---

## तालिकाओं में एस्केप करना

| Column 1 | Column 2 |
|----------|----------|
| \*Escaped\* | \_Escaped\_ |

---

## बैकस्लैश एस्केप करना

\\ Backslash escaped

\\\\ Double backslash

---

## संयोजन एस्केप करना

\*\*Bold\*\* and \*italic\* and \`code\` all escaped.

---

## Unicode एस्केप करना

\你好 should not be processed

\café should not be processed

---

## संख्याएं एस्केप करना

\123 should not be processed

\1.2.3 should not be processed

---

## उद्धरण एस्केप करना

\"Double quotes\"

\'Single quotes\'

---

## डॉलर साइन एस्केप करना

\$Dollar sign\$

---

## प्रतिशत चिह्न एस्केप करना

\%Percent sign\%

---

## एम्परसेंड एस्केप करना

\&Ampersand\&

---

## कोणीय ब्रैकेट एस्केप करना

\<Less than\>

\>Greater than\>

---

## वर्ग ब्रैकेट एस्केप करना

\[Left bracket\]

\]Right bracket\]

---

## ब्रेस एस्केप करना

\{Left brace\}

\}Right brace\}

---

## कोष्ठक एस्केप करना

\(Left parenthesis\)

\)Right parenthesis\)

---

## विस्मयादिबोधक चिह्न एस्केप करना

\!Exclamation mark\!

---

## प्रश्न चिह्न एस्केप करना

\?Question mark\?

---

## कोलन एस्केप करना

\:Colon\:

---

## सेमीकोलन एस्केप करना

\;Semicolon\;

---

## कॉमा एस्केप करना

\,Comma\,

---

## पीरियड एस्केप करना

\.Period\.

---

## कई विशेष वर्ण एस्केप करना

\*\*Bold\*\* \*italic\* \`code\` \[link\] all escaped.

---

## विभिन्न स्थितियों में एस्केप करना

Start: \*escaped\*

Middle: text \*escaped\* text

End: text \*escaped\*

---

## आसन्न पाठ के साथ एस्केप करना

\*escaped\*text

text\*escaped\*

\*escaped\*text\*escaped\*

---

## शाब्दिक अर्थ संरक्षित करते हुए एस्केप करना

\* should display as *

\_ should display as _

\` should display as `

\# should display as #

---

## एस्केप किनारे के मामले

\*\* should display as **

\*\*\* should display as ***

\_\_ should display as __

\_\_\_ should display as ___

