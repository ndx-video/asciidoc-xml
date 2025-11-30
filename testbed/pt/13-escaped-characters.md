---
title: Caracteres Escapados e Casos Especiais
categories:
  - markdown
  - escaping
tags:
  - special characters
  - edge cases
---

# Caracteres Escapados e Casos Especiais

## Escapar Asteriscos

\*This should not be italic\*

\*\*This should not be bold\*\*

\*\*\*This should not be bold and italic\*\*\*

---

## Escapar Sublinhados

\_This should not be italic\_

\_\_This should not be bold\_\_

\_\_\_This should not be bold and italic\_\_\_

---

## Escapar Crases

\`This should not be code\`

\`\`This should not be code\`\`

---

## Escapar Colchetes

\[This should not be a link\](https://example.com)

\!\[This should not be an image\](image.png)

---

## Escapar Hashtags

\# This should not be a header

\## This should not be a header

---

## Escapar Sinais de Mais

\+ This should not be a list item

---

## Escapar Sinais de Menos

\- This should not be a list item

---

## Escapar Pontos

\. This should not be a numbered list

---

## Escapar Maior Que

\> This should not be a blockquote

---

## Escapar Barras Verticais

\| This should not be a table cell

---

## Escapar Tildes

\~\~This should not be strikethrough\~\~

---

## Escapar Parênteses

\(Parentheses\)

\[Brackets\]

\{Braces\}

---

## Escapar Caracteres Especiais

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

## Escapar em Diferentes Contextos

### Em Parágrafos

\*Escaped asterisk\* in paragraph.

### Em Listas

* Item with \*escaped asterisk\*
* Item with \_escaped underscore\_

### Em Citações

> Blockquote with \*escaped asterisk\*

### Em Blocos de Código

```
\*Escaped asterisk\* in code block
```

---

## Escapar URLs

\https://example.com should not be a link

---

## Escapar Email

\user@example.com should not be a link

---

## Escapar HTML

\&lt;div&gt; should display as &lt;div&gt;

\&amp; should display as &

---

## Escapar Múltiplos Caracteres

\*\*Bold\*\* and \*italic\* escaped.

---

## Escapar no Início da Linha

\# Not a header

\* Not a list

\- Not a list

---

## Escapar no Final da Linha

Text with \*escaped asterisk\*

---

## Escapar com Espaços

\* Escaped with space

\*Escaped without space\*

---

## Escapar Sequências Especiais

\*\*\* Triple asterisk escaped

\_\_\_ Triple underscore escaped

\`\`\` Triple backtick escaped

---

## Escapar em Links

[Link with \*asterisk\*](https://example.com)

[Link with \_underscore\_](https://example.com)

---

## Escapar em Imagens

![Image with \*asterisk\*](image.png)

![Image with \_underscore\_](image.png)

---

## Escapar em Cabeçalhos

# Header with \*asterisk\*

## Header with \_underscore\_

---

## Escapar em Tabelas

| Column 1 | Column 2 |
|----------|----------|
| \*Escaped\* | \_Escaped\_ |

---

## Escapar Barras Invertidas

\\ Backslash escaped

\\\\ Double backslash

---

## Escapar Combinações

\*\*Bold\*\* and \*italic\* and \`code\` all escaped.

---

## Escapar Unicode

\你好 should not be processed

\café should not be processed

---

## Escapar Números

\123 should not be processed

\1.2.3 should not be processed

---

## Escapar Aspas

\"Double quotes\"

\'Single quotes\'

---

## Escapar Sinais de Dólar

\$Dollar sign\$

---

## Escapar Sinais de Porcentagem

\%Percent sign\%

---

## Escapar E Comerciais

\&Ampersand\&

---

## Escapar Colchetes Angulares

\<Less than\>

\>Greater than\>

---

## Escapar Colchetes Quadrados

\[Left bracket\]

\]Right bracket\]

---

## Escapar Chaves

\{Left brace\}

\}Right brace\}

---

## Escapar Parênteses

\(Left parenthesis\)

\)Right parenthesis\)

---

## Escapar Sinais de Exclamação

\!Exclamation mark\!

---

## Escapar Sinais de Interrogação

\?Question mark\?

---

## Escapar Dois Pontos

\:Colon\:

---

## Escapar Ponto e Vírgula

\;Semicolon\;

---

## Escapar Vírgulas

\,Comma\,

---

## Escapar Pontos

\.Period\.

---

## Escapar Múltiplos Caracteres Especiais

\*\*Bold\*\* \*italic\* \`code\` \[link\] all escaped.

---

## Escapar em Diferentes Posições

Start: \*escaped\*

Middle: text \*escaped\* text

End: text \*escaped\*

---

## Escapar com Texto Adjacente

\*escaped\*text

text\*escaped\*

\*escaped\*text\*escaped\*

---

## Escapar Preservando Significado Literal

\* should display as *

\_ should display as _

\` should display as `

\# should display as #

---

## Casos Extremos de Escape

\*\* should display as **

\*\*\* should display as ***

\_\_ should display as __

\_\_\_ should display as ___

