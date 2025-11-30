---
title: Casos Extremos y Escenarios Especiales
categories:
  - markdown
  - edge cases
tags:
  - testing
  - special cases
---

# Casos Extremos y Escenarios Especiales

## Documento Vacío

---

## Documento con Solo Espacios en Blanco

    

---

## Documento con Solo Saltos de Línea



---

## Documento que Comienza con Encabezado

# First Header

Content.

---

## Documento que Termina con Encabezado

Content.

# Final Header

---

## Documento con Solo Encabezado

# Only Header

---

## Documento con Solo Párrafo

This is the only paragraph.

---

## Encabezados Adyacentes

# Header 1
## Header 2
### Header 3

---

## Encabezados Sin Contenido

# Header

## Another Header

---

## Párrafos Sin Separación

Paragraph one.
Paragraph two.
Paragraph three.

---

## Listas Sin Líneas en Blanco

* Item 1
* Item 2
# Header
* Item 3

---

## Elementos de Lista Vacíos

* 
* Item with content
* 
* Another item

---

## Listas con Solo Espacios

*     
*     

---

## Tablas con Celdas Vacías

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
|          |          |          |
| Data     |          | Data     |

---

## Tablas con Solo Encabezados

| Header 1 | Header 2 |
|----------|----------|

---

## Bloques de Código con Solo Espacios en Blanco

```
    
```

---

## Bloques de Código con Solo Saltos de Línea

```

```

---

## Citas con Solo Espacios

>     

---

## Citas con Solo Saltos de Línea

>

---

## Enlaces con Texto Vacío

[](https://example.com)

[ ](https://example.com)

---

## Imágenes con Texto Alternativo Vacío

![](image.png)

![ ](image.png)

---

## Énfasis con Contenido Vacío

****

* *

** **

---

## Énfasis Adyacente al Texto

**Bold**text

*Italic*text

`Code`text

---

## Enlaces Adyacentes al Texto

[Link](https://example.com)text

text[Link](https://example.com)

---

## Imágenes Adyacentes al Texto

![Image](image.png)text

text![Image](image.png)

---

## Encabezados con Solo Espacios

#     

##     

---

## Encabezados con Solo Caracteres Especiales

# ***

## ---

### ___

---

## Listas con Marcadores Mixtos

* Item 1
- Item 2
+ Item 3

---

## Listas que Comienzan con Diferentes Números

5. Item 5
10. Item 10
1. Item 1

---

## Tablas con Separadores Inconsistentes

| Column 1 | Column 2 |
|:---------|---------|
| Data     | Data    |

---

## Bloques de Código con Diferentes Longitudes de Cerca

````
Code with four backticks
````

`````
Code with five backticks
`````

---

## Estructuras Anidadas Profundidad Máxima

> Blockquote
> > Nested
> > > Deep nested
> > > > Very deep

---

## Casos Extremos de Formato Mixto

**Bold*Italic**Bold*

*Italic**Bold*Italic*

---

## Enlaces en Diferentes Posiciones

[Start](https://example.com) text.

Text [middle](https://example.com) text.

Text [end](https://example.com).

---

## Imágenes en Diferentes Posiciones

![Start](image.png) text.

Text ![middle](image.png) text.

Text ![end](image.png).

---

## Código en Diferentes Posiciones

`Start` text.

Text `middle` text.

Text `end`.

---

## Énfasis en Diferentes Posiciones

**Start** text.

Text **middle** text.

Text **end**.

---

## Caracteres Especiales en los Límites

*Item with "quotes"*

*Item with (parentheses)*

*Item with [brackets]*

---

## Casos Extremos Unicode

café naïve résumé

你好世界

こんにちは世界

مرحبا بالعالم

Привет мир

---

## Números en Diferentes Contextos

123

1,000

1.234

-42

+100

---

## URLs en Diferentes Formatos

https://example.com

http://example.com

ftp://ftp.example.com

file:///path/to/file

---

## Formatos de Correo Electrónico

user@example.com

user.name@example.com

user+tag@example.com

user_name@example.com

---

## Rutas de Archivo

/path/to/file

./relative/path

../parent/path

~/home/path

C:\Windows\Path

---

## Números de Versión

v1.0.0

1.2.3

1.2.3-beta.1

1.2.3-alpha.2+build.123

---

## Fechas

2024-01-15

01/15/2024

January 15, 2024

2024-01-15T10:30:00Z

---

## Tiempos

10:30 AM

22:30

2h30m

5s

---

## Porcentajes

50%

100%

0%

123.45%

---

## Moneda

$19.99

€15.50

£10.00

¥1000

---

## Direcciones IP

192.168.1.1

127.0.0.1

2001:0db8:85a3:0000:0000:8a2e:0370:7334

---

## Etiquetas

#hashtag

#hashtag-with-dashes

#hashtag_with_underscores

#123hashtag

---

## Menciones

@username

@user_name

@user-name

@user123

---

## Caracteres Especiales Mixtos

$100 (50% off) @user #tag

Version 1.2.3 released on 2024-01-15

Email: user@example.com, Phone: +1-234-567-8900

---

## Frontmatter Vacío

---
---

---

## Frontmatter con Solo Espacios

---
    
---

---

## Frontmatter con Solo Saltos de Línea

---


---

---

## Múltiples Bloques Frontmatter

---
title: First
---

---
title: Second
---

---

## Frontmatter Sin Separadores

title: No separators
author: Author Name

---

## Contenido que Comienza Inmediatamente

Content without frontmatter or header.

---

## Contenido que Termina Abruptamente

Content that ends without newline.

