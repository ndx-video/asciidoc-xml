---
title: Imagens - Estilo Inline e Referência
categories:
  - markdown
  - images
tags:
  - graphics
  - media
  - pictures
---

# Imagens

## Imagens Inline

![Texto alternativo](image.png)

![Texto alternativo com espaços](image with spaces.png)

![Texto alternativo](image.png "Título da Imagem")

![Texto alternativo](image.png 'Título da Imagem')

![Texto alternativo](image.png (Título da Imagem))

---

## Imagens com Diferentes Formatos

![Imagem PNG](image.png)

![Imagem JPEG](image.jpg)

![Imagem GIF](image.gif)

![Imagem SVG](image.svg)

![Imagem WebP](image.webp)

---

## Imagens com Caminhos

![Imagem raiz](/image.png)

![Imagem subdiretório](./images/image.png)

![Diretório pai](../images/image.png)

![Caminho profundo](./images/subfolder/image.png)

---

## Imagens com URLs

![Imagem remota](https://example.com/image.png)

![Imagem remota com caminho](https://example.com/images/photo.jpg)

![Imagem remota com consulta](https://example.com/image.png?size=large)

---

## Imagens Estilo Referência

![Imagem de referência][img1]

![Imagem de referência implícita][]

![Imagem de referência curta]

[img1]: image1.png
[implicit reference image]: image2.png
[shortcut reference image]: image3.png

---

## Imagens de Referência com Títulos

![Imagem de referência com título][img2]

[img2]: image.png "Título da Imagem"

![Imagem de referência com aspas simples][img3]

[img3]: image.png 'Título da Imagem'

![Imagem de referência com parênteses][img4]

[img4]: image.png (Título da Imagem)

---

## Imagens com Caracteres Especiais em Texto Alternativo

![Imagem com "aspas"](image.png)

![Imagem com (parênteses)](image.png)

![Imagem com [colchetes]](image.png)

![Imagem com {chaves}](image.png)

![Imagem com <tags>](image.png)

![Imagem com $dólar$](image.png)

---

## Imagens com Unicode em Texto Alternativo

![Imagem com café](image.png)

![Imagem com 你好](image.png)

![Imagem com こんにちは](image.png)

![Imagem com مرحبا](image.png)

---

## Imagens com Números em Texto Alternativo

![Imagem 123](image.png)

![Imagem v1.2.3](image.png)

![Imagem 50%](image.png)

---

## Imagens com Ênfase em Texto Alternativo

![**Texto alternativo negrito**](image.png)

![*Texto alternativo itálico*](image.png)

![***Texto alternativo negrito e itálico***](image.png)

---

## Imagens com Links

[![Imagem vinculada](image.png)](https://example.com)

[![Imagem vinculada com alt](image.png "Título da Imagem")](https://example.com "Título do Link")

---

## Imagens em Listas

* Item com ![imagem](image.png)
* Item com ![imagem de referência][listimg]
* Item com ![imagem](image.png "Título")

[listimg]: image.png

---

## Imagens em Citações

> Esta é uma citação com ![imagem](image.png)

> Esta é uma citação com ![imagem de referência][blockquoteimg]

[blockquoteimg]: image.png

---

## Imagens em Cabeçalhos

# Cabeçalho com ![imagem](image.png)

## Cabeçalho com ![imagem de referência][headerimg]

[headerimg]: image.png

---

## Imagens em Tabelas

| Column 1 | Column 2 |
|----------|----------|
| ![Imagem](image.png) | ![Referência][tableimg] |
| Texto | ![Imagem com título](image.png "Título") |

[tableimg]: image.png

---

## Imagens com Texto Alternativo Vazio

![](image.png)

![ ](image.png)

![   ](image.png)

---

## Imagens Apenas com Espaços em Texto Alternativo

![   ](image.png)

---

## Imagens com Texto Alternativo Longo

![Este é um texto alternativo muito longo que descreve a imagem em grande detalhe e fornece informações abrangentes sobre o que a imagem contém e seu propósito no documento](image.png)

---

## Múltiplas Imagens em Um Parágrafo

Este parágrafo tem ![primeira imagem](image1.png), ![segunda imagem](image2.png) e ![terceira imagem](image3.png).

Este parágrafo mistura ![imagem inline](image.png), ![imagem de referência][multimg] e texto.

[multimg]: image.png

---

## Imagens com HTML

<img src="image.png" alt="Imagem HTML">

![Imagem Markdown](image.png) e <img src="image.png" alt="Imagem HTML">

---

## Imagens com Caracteres Escapados

\![Imagem escapada](image.png)

![Imagem com \*asterisco\*](image.png)

![Imagem com \_sublinhado\_](image.png)

---

## Imagens em Blocos de Código

```
![Imagem](image.png) não deveria funcionar em blocos de código
```

---

## Imagens com Esquemas de URL Especiais

![Imagem Data URI](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==)

---

## Imagens com Codificação Percentual

![Imagem com codificação](https://example.com/image%20with%20spaces.png)

![Imagem com codificação unicode](https://example.com/%E4%BD%A0%E5%A5%BD.png)

---

## Combinações Complexas de Imagens

**Texto negrito** com ![imagem](image.png) e texto *itálico*.

![Imagem](image.png) com [link](https://example.com) e `código`.

![**Alt negrito**](image.png) com texto *itálico* e [link](https://example.com).

---

## Imagens com Diferentes Tamanhos (atributos HTML)

<img src="image.png" alt="Imagem pequena" width="100" height="100">

<img src="image.png" alt="Imagem grande" width="800" height="600">

---

## Imagens com Classes CSS (HTML)

<img src="image.png" alt="Imagem estilizada" class="custom-class">

<img src="image.png" alt="Imagem centralizada" style="display: block; margin: 0 auto;">

---

## Imagens com Títulos e Legendas

![Imagem](image.png "Este é um título")

![Imagem](image.png "Título com 'aspas'")

![Imagem](image.png "Título com (parênteses)")

---

## Imagens com Caminhos Relativos

![Mesmo diretório](./image.png)

![Subdiretório](./images/image.png)

![Diretório pai](../image.png)

![Relativo à raiz](/image.png)

---

## Imagens com URLs Absolutas

![URL absoluta](https://example.com/image.png)

![URL absoluta com caminho](https://example.com/images/photo.jpg)

![URL absoluta com consulta](https://example.com/image.png?v=1&size=large)

![URL absoluta com fragmento](https://example.com/image.png#section)

---

## Imagens em Estruturas Aninhadas

### Imagens em Listas Aninhadas

* Primeiro nível
  * Segundo nível com ![imagem](image.png)
    * Terceiro nível com ![imagem](image.png)

### Imagens em Citações Aninhadas

> Citação externa
> > Citação interna com ![imagem](image.png)

---

## Imagens com Referências Quebradas

![Imagem quebrada](nonexistent.png)

![Referência quebrada][broken]

[broken]: nonexistent.png

---

## Imagens com Nomes de Arquivo Especiais

![Imagem com espaços](image with spaces.png)

![Imagem.com.pontos](image.with.dots.png)

![Imagem-com-hifens](image-with-dashes.png)

![Imagem_com_sublinhados](Image_with_underscores.png)

![Imagem123](Image123.png)

![Imagem@especial](Image@special.png)

