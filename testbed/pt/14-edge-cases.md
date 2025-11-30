---
title: Casos Extremos e Cenários Especiais
categories:
  - markdown
  - edge cases
tags:
  - testing
  - special cases
---

# Casos Extremos e Cenários Especiais

## Documento Vazio

---

## Documento com Apenas Espaços em Branco

    

---

## Documento com Apenas Quebras de Linha



---

## Documento que Começa com Cabeçalho

# First Header

Content.

---

## Documento que Termina com Cabeçalho

Content.

# Final Header

---

## Documento com Apenas Cabeçalho

# Only Header

---

## Documento com Apenas Parágrafo

This is the only paragraph.

---

## Cabeçalhos Adjacentes

# Header 1
## Header 2
### Header 3

---

## Cabeçalhos Sem Conteúdo

# Header

## Another Header

---

## Parágrafos Sem Separação

Paragraph one.
Paragraph two.
Paragraph three.

---

## Listas Sem Linhas em Branco

* Item 1
* Item 2
# Header
* Item 3

---

## Itens de Lista Vazios

* 
* Item with content
* 
* Another item

---

## Listas com Apenas Espaços

*     
*     

---

## Tabelas com Células Vazias

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
|          |          |          |
| Data     |          | Data     |

---

## Tabelas com Apenas Cabeçalhos

| Header 1 | Header 2 |
|----------|----------|

---

## Blocos de Código com Apenas Espaços em Branco

```
    
```

---

## Blocos de Código com Apenas Quebras de Linha

```

```

---

## Citações com Apenas Espaços

>     

---

## Citações com Apenas Quebras de Linha

>

---

## Links com Texto Vazio

[](https://example.com)

[ ](https://example.com)

---

## Imagens com Texto Alternativo Vazio

![](image.png)

![ ](image.png)

---

## Ênfase com Conteúdo Vazio

****

* *

** **

---

## Ênfase Adjacente ao Texto

**Bold**text

*Italic*text

`Code`text

---

## Links Adjacentes ao Texto

[Link](https://example.com)text

text[Link](https://example.com)

---

## Imagens Adjacentes ao Texto

![Image](image.png)text

text![Image](image.png)

---

## Cabeçalhos com Apenas Espaços

#     

##     

---

## Cabeçalhos com Apenas Caracteres Especiais

# ***

## ---

### ___

---

## Listas com Marcadores Mistos

* Item 1
- Item 2
+ Item 3

---

## Listas que Começam com Números Diferentes

5. Item 5
10. Item 10
1. Item 1

---

## Tabelas com Separadores Inconsistentes

| Column 1 | Column 2 |
|:---------|---------|
| Data     | Data    |

---

## Blocos de Código com Diferentes Comprimentos de Cerca

````
Code with four backticks
````

`````
Code with five backticks
`````

---

## Estruturas Aninhadas Profundidade Máxima

> Blockquote
> > Nested
> > > Deep nested
> > > > Very deep

---

## Casos Extremos de Formato Misto

**Bold*Italic**Bold*

*Italic**Bold*Italic*

---

## Links em Diferentes Posições

[Start](https://example.com) text.

Text [middle](https://example.com) text.

Text [end](https://example.com).

---

## Imagens em Diferentes Posições

![Start](image.png) text.

Text ![middle](image.png) text.

Text ![end](image.png).

---

## Código em Diferentes Posições

`Start` text.

Text `middle` text.

Text `end`.

---

## Ênfase em Diferentes Posições

**Start** text.

Text **middle** text.

Text **end**.

---

## Caracteres Especiais nos Limites

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

## Números em Diferentes Contextos

123

1,000

1.234

-42

+100

---

## URLs em Diferentes Formatos

https://example.com

http://example.com

ftp://ftp.example.com

file:///path/to/file

---

## Formatos de Email

user@example.com

user.name@example.com

user+tag@example.com

user_name@example.com

---

## Caminhos de Arquivo

/path/to/file

./relative/path

../parent/path

~/home/path

C:\Windows\Path

---

## Números de Versão

v1.0.0

1.2.3

1.2.3-beta.1

1.2.3-alpha.2+build.123

---

## Datas

2024-01-15

01/15/2024

January 15, 2024

2024-01-15T10:30:00Z

---

## Horários

10:30 AM

22:30

2h30m

5s

---

## Porcentagens

50%

100%

0%

123.45%

---

## Moeda

$19.99

€15.50

£10.00

¥1000

---

## Endereços IP

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

## Menções

@username

@user_name

@user-name

@user123

---

## Caracteres Especiais Mistos

$100 (50% off) @user #tag

Version 1.2.3 released on 2024-01-15

Email: user@example.com, Phone: +1-234-567-8900

---

## Frontmatter Vazio

---
---

---

## Frontmatter com Apenas Espaços

---
    
---

---

## Frontmatter com Apenas Quebras de Linha

---


---

---

## Múltiplos Blocos Frontmatter

---
title: First
---

---
title: Second
---

---

## Frontmatter Sem Separadores

title: No separators
author: Author Name

---

## Conteúdo que Começa Imediatamente

Content without frontmatter or header.

---

## Conteúdo que Termina Abruptamente

Content that ends without newline.

