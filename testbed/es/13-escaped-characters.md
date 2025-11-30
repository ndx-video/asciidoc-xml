---
title: Caracteres Escapados y Casos Especiales
categories:
  - markdown
  - escaping
tags:
  - special characters
  - edge cases
---

# Caracteres Escapados y Casos Especiales

## Escapar Asteriscos

\*This should not be italic\*

\*\*This should not be bold\*\*

\*\*\*This should not be bold and italic\*\*\*

---

## Escapar Guiones Bajos

\_This should not be italic\_

\_\_This should not be bold\_\_

\_\_\_This should not be bold and italic\_\_\_

---

## Escapar Comillas Invertidas

\`This should not be code\`

\`\`This should not be code\`\`

---

## Escapar Corchetes

\[This should not be a link\](https://example.com)

\!\[This should not be an image\](image.png)

---

## Escapar Etiquetas

\# This should not be a header

\## This should not be a header

---

## Escapar Signos Más

\+ This should not be a list item

---

## Escapar Signos Menos

\- This should not be a list item

---

## Escapar Puntos

\. This should not be a numbered list

---

## Escapar Mayor Que

\> This should not be a blockquote

---

## Escapar Barras Verticales

\| This should not be a table cell

---

## Escapar Tildes

\~\~This should not be strikethrough\~\~

---

## Escapar Paréntesis

\(Parentheses\)

\[Brackets\]

\{Braces\}

---

## Escapar Caracteres Especiales

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

## Escapar en Diferentes Contextos

### En Párrafos

\*Escaped asterisk\* in paragraph.

### En Listas

* Item with \*escaped asterisk\*
* Item with \_escaped underscore\_

### En Citas

> Blockquote with \*escaped asterisk\*

### En Bloques de Código

```
\*Escaped asterisk\* in code block
```

---

## Escapar URLs

\https://example.com should not be a link

---

## Escapar Correo Electrónico

\user@example.com should not be a link

---

## Escapar HTML

\&lt;div&gt; should display as &lt;div&gt;

\&amp; should display as &

---

## Escapar Múltiples Caracteres

\*\*Bold\*\* and \*italic\* escaped.

---

## Escapar al Inicio de Línea

\# Not a header

\* Not a list

\- Not a list

---

## Escapar al Final de Línea

Text with \*escaped asterisk\*

---

## Escapar con Espacios

\* Escaped with space

\*Escaped without space\*

---

## Escapar Secuencias Especiales

\*\*\* Triple asterisk escaped

\_\_\_ Triple underscore escaped

\`\`\` Triple backtick escaped

---

## Escapar en Enlaces

[Link with \*asterisk\*](https://example.com)

[Link with \_underscore\_](https://example.com)

---

## Escapar en Imágenes

![Image with \*asterisk\*](image.png)

![Image with \_underscore\_](image.png)

---

## Escapar en Encabezados

# Header with \*asterisk\*

## Header with \_underscore\_

---

## Escapar en Tablas

| Column 1 | Column 2 |
|----------|----------|
| \*Escaped\* | \_Escaped\_ |

---

## Escapar Barras Invertidas

\\ Backslash escaped

\\\\ Double backslash

---

## Escapar Combinaciones

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

## Escapar Comillas

\"Double quotes\"

\'Single quotes\'

---

## Escapar Signos de Dólar

\$Dollar sign\$

---

## Escapar Signos de Porcentaje

\%Percent sign\%

---

## Escapar Y Comerciales

\&Ampersand\&

---

## Escapar Corchetes Angulares

\<Less than\>

\>Greater than\>

---

## Escapar Corchetes Cuadrados

\[Left bracket\]

\]Right bracket\]

---

## Escapar Llaves

\{Left brace\}

\}Right brace\}

---

## Escapar Paréntesis

\(Left parenthesis\)

\)Right parenthesis\)

---

## Escapar Signos de Exclamación

\!Exclamation mark\!

---

## Escapar Signos de Interrogación

\?Question mark\?

---

## Escapar Dos Puntos

\:Colon\:

---

## Escapar Puntos y Comas

\;Semicolon\;

---

## Escapar Comas

\,Comma\,

---

## Escapar Puntos

\.Period\.

---

## Escapar Múltiples Caracteres Especiales

\*\*Bold\*\* \*italic\* \`code\` \[link\] all escaped.

---

## Escapar en Diferentes Posiciones

Start: \*escaped\*

Middle: text \*escaped\* text

End: text \*escaped\*

---

## Escapar con Texto Adyacente

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

## Casos Límite de Escape

\*\* should display as **

\*\*\* should display as ***

\_\_ should display as __

\_\_\_ should display as ___

