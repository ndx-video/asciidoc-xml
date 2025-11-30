---
title: 转义字符和特殊情况
categories:
  - markdown
  - escaping
tags:
  - special characters
  - edge cases
---

# 转义字符和特殊情况

## 转义星号

\*This should not be italic\*

\*\*This should not be bold\*\*

\*\*\*This should not be bold and italic\*\*\*

---

## 转义下划线

\_This should not be italic\_

\_\_This should not be bold\_\_

\_\_\_This should not be bold and italic\_\_\_

---

## 转义反引号

\`This should not be code\`

\`\`This should not be code\`\`

---

## 转义方括号

\[This should not be a link\](https://example.com)

\!\[This should not be an image\](image.png)

---

## 转义井号

\# This should not be a header

\## This should not be a header

---

## 转义加号

\+ This should not be a list item

---

## 转义减号

\- This should not be a list item

---

## 转义句号

\. This should not be a numbered list

---

## 转义大于号

\> This should not be a blockquote

---

## 转义管道符

\| This should not be a table cell

---

## 转义波浪号

\~\~This should not be strikethrough\~\~

---

## 转义括号

\(Parentheses\)

\[Brackets\]

\{Braces\}

---

## 转义特殊字符

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

## 不同上下文中的转义

### 在段落中

\*Escaped asterisk\* in paragraph.

### 在列表中

* Item with \*escaped asterisk\*
* Item with \_escaped underscore\_

### 在引用块中

> Blockquote with \*escaped asterisk\*

### 在代码块中

```
\*Escaped asterisk\* in code block
```

---

## 转义 URL

\https://example.com should not be a link

---

## 转义电子邮件

\user@example.com should not be a link

---

## 转义 HTML

\&lt;div&gt; should display as &lt;div&gt;

\&amp; should display as &

---

## 转义多个字符

\*\*Bold\*\* and \*italic\* escaped.

---

## 在行首转义

\# Not a header

\* Not a list

\- Not a list

---

## 在行尾转义

Text with \*escaped asterisk\*

---

## 带空格的转义

\* Escaped with space

\*Escaped without space\*

---

## 转义特殊序列

\*\*\* Triple asterisk escaped

\_\_\_ Triple underscore escaped

\`\`\` Triple backtick escaped

---

## 在链接中转义

[Link with \*asterisk\*](https://example.com)

[Link with \_underscore\_](https://example.com)

---

## 在图片中转义

![Image with \*asterisk\*](image.png)

![Image with \_underscore\_](image.png)

---

## 在标题中转义

# Header with \*asterisk\*

## Header with \_underscore\_

---

## 在表格中转义

| Column 1 | Column 2 |
|----------|----------|
| \*Escaped\* | \_Escaped\_ |

---

## 转义反斜杠

\\ Backslash escaped

\\\\ Double backslash

---

## 转义组合

\*\*Bold\*\* and \*italic\* and \`code\` all escaped.

---

## 转义 Unicode

\你好 should not be processed

\café should not be processed

---

## 转义数字

\123 should not be processed

\1.2.3 should not be processed

---

## 转义引号

\"Double quotes\"

\'Single quotes\'

---

## 转义美元符号

\$Dollar sign\$

---

## 转义百分号

\%Percent sign\%

---

## 转义和号

\&Ampersand\&

---

## 转义尖括号

\<Less than\>

\>Greater than\>

---

## 转义方括号

\[Left bracket\]

\]Right bracket\]

---

## 转义大括号

\{Left brace\}

\}Right brace\}

---

## 转义圆括号

\(Left parenthesis\)

\)Right parenthesis\)

---

## 转义感叹号

\!Exclamation mark\!

---

## 转义问号

\?Question mark\?

---

## 转义冒号

\:Colon\:

---

## 转义分号

\;Semicolon\;

---

## 转义逗号

\,Comma\,

---

## 转义句号

\.Period\.

---

## 转义多个特殊字符

\*\*Bold\*\* \*italic\* \`code\` \[link\] all escaped.

---

## 在不同位置转义

Start: \*escaped\*

Middle: text \*escaped\* text

End: text \*escaped\*

---

## 与相邻文本转义

\*escaped\*text

text\*escaped\*

\*escaped\*text\*escaped\*

---

## 转义保留字面意思

\* should display as *

\_ should display as _

\` should display as `

\# should display as #

---

## 转义边缘情况

\*\* should display as **

\*\*\* should display as ***

\_\_ should display as __

\_\_\_ should display as ___

