---
title: Escaped Characters and Special Cases
categories:
  - markdown
  - escaping
tags:
  - special characters
  - edge cases
---

# Escaped Characters and Special Cases

## Escaping Asterisks

\*This should not be italic\*

\*\*This should not be bold\*\*

\*\*\*This should not be bold and italic\*\*\*

---

## Escaping Underscores

\_This should not be italic\_

\_\_This should not be bold\_\_

\_\_\_This should not be bold and italic\_\_\_

---

## Escaping Backticks

\`This should not be code\`

\`\`This should not be code\`\`

---

## Escaping Brackets

\[This should not be a link\](https://example.com)

\!\[This should not be an image\](image.png)

---

## Escaping Hashtags

\# This should not be a header

\## This should not be a header

---

## Escaping Plus Signs

\+ This should not be a list item

---

## Escaping Minus Signs

\- This should not be a list item

---

## Escaping Periods

\. This should not be a numbered list

---

## Escaping Greater Than

\> This should not be a blockquote

---

## Escaping Pipes

\| This should not be a table cell

---

## Escaping Tildes

\~\~This should not be strikethrough\~\~

---

## Escaping Parentheses

\(Parentheses\)

\[Brackets\]

\{Braces\}

---

## Escaping Special Characters

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

## Escaping in Different Contexts

### In Paragraphs

\*Escaped asterisk\* in paragraph.

### In Lists

* Item with \*escaped asterisk\*
* Item with \_escaped underscore\_

### In Blockquotes

> Blockquote with \*escaped asterisk\*

### In Code Blocks

```
\*Escaped asterisk\* in code block
```

---

## Escaping URLs

\https://example.com should not be a link

---

## Escaping Email

\user@example.com should not be a link

---

## Escaping HTML

\&lt;div&gt; should display as &lt;div&gt;

\&amp; should display as &

---

## Escaping Multiple Characters

\*\*Bold\*\* and \*italic\* escaped.

---

## Escaping at Start of Line

\# Not a header

\* Not a list

\- Not a list

---

## Escaping at End of Line

Text with \*escaped asterisk\*

---

## Escaping with Spaces

\* Escaped with space

\*Escaped without space\*

---

## Escaping Special Sequences

\*\*\* Triple asterisk escaped

\_\_\_ Triple underscore escaped

\`\`\` Triple backtick escaped

---

## Escaping in Links

[Link with \*asterisk\*](https://example.com)

[Link with \_underscore\_](https://example.com)

---

## Escaping in Images

![Image with \*asterisk\*](image.png)

![Image with \_underscore\_](image.png)

---

## Escaping in Headers

# Header with \*asterisk\*

## Header with \_underscore\_

---

## Escaping in Tables

| Column 1 | Column 2 |
|----------|----------|
| \*Escaped\* | \_Escaped\_ |

---

## Escaping Backslashes

\\ Backslash escaped

\\\\ Double backslash

---

## Escaping Combinations

\*\*Bold\*\* and \*italic\* and \`code\` all escaped.

---

## Escaping Unicode

\你好 should not be processed

\café should not be processed

---

## Escaping Numbers

\123 should not be processed

\1.2.3 should not be processed

---

## Escaping Quotes

\"Double quotes\"

\'Single quotes\'

---

## Escaping Dollar Signs

\$Dollar sign\$

---

## Escaping Percent Signs

\%Percent sign\%

---

## Escaping Ampersands

\&Ampersand\&

---

## Escaping Angle Brackets

\<Less than\>

\>Greater than\>

---

## Escaping Square Brackets

\[Left bracket\]

\]Right bracket\]

---

## Escaping Curly Braces

\{Left brace\}

\}Right brace\}

---

## Escaping Parentheses

\(Left parenthesis\)

\)Right parenthesis\)

---

## Escaping Exclamation Marks

\!Exclamation mark\!

---

## Escaping Question Marks

\?Question mark\?

---

## Escaping Colons

\:Colon\:

---

## Escaping Semicolons

\;Semicolon\;

---

## Escaping Commas

\,Comma\,

---

## Escaping Periods

\.Period\.

---

## Escaping Multiple Special Characters

\*\*Bold\*\* \*italic\* \`code\` \[link\] all escaped.

---

## Escaping in Different Positions

Start: \*escaped\*

Middle: text \*escaped\* text

End: text \*escaped\*

---

## Escaping with Adjacent Text

\*escaped\*text

text\*escaped\*

\*escaped\*text\*escaped\*

---

## Escaping Preserving Literal Meaning

\* should display as *

\_ should display as _

\` should display as `

\# should display as #

---

## Escaping Edge Cases

\*\* should display as **

\*\*\* should display as ***

\_\_ should display as __

\_\_\_ should display as ___

