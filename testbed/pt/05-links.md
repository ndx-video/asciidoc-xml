---
title: Links - Inline, Referência e Autolinks
categories:
  - markdown
  - links
tags:
  - hyperlinks
  - references
  - urls
---

# Links

## Links Inline

Este é um [link inline](https://example.com).

Este é um [link com título](https://example.com "Título do Link").

Este é um [link com aspas simples](https://example.com 'Título do Link').

Este é um [link com parênteses](https://example.com (Título do Link)).

---

## Links com Diferentes Protocolos

[Link HTTP](http://example.com)

[Link HTTPS](https://example.com)

[Link FTP](ftp://ftp.example.com)

[Link de email](mailto:user@example.com)

[Link de arquivo](file:///path/to/file)

---

## Links com Caminhos

[Caminho raiz](https://example.com/)

[Subdiretório](https://example.com/path/to/resource)

[Extensão de arquivo](https://example.com/file.html)

[String de consulta](https://example.com/page?param=value&other=123)

[Fragmento](https://example.com/page#section)

[URL complexa](https://example.com/path/to/resource?query=value&other=123#fragment)

---

## Links com Caracteres Especiais

[Link com espaços](https://example.com/path with spaces)

[Link com unicode](https://example.com/你好)

[Link com codificação](https://example.com/path%20with%20spaces)

---

## Links Estilo Referência

Este é um [link de referência][ref1].

Este é um [link de referência implícito][].

Este é um [link de referência curto].

[ref1]: https://example.com/reference1
[implicit reference link]: https://example.com/implicit
[shortcut reference link]: https://example.com/shortcut

---

## Links de Referência com Títulos

Este é um [link de referência com título][ref2].

[ref2]: https://example.com/reference2 "Título de Referência"

Este é um [link de referência com aspas simples][ref3].

[ref3]: https://example.com/reference3 'Título de Referência'

Este é um [link de referência com parênteses][ref4].

[ref4]: https://example.com/reference4 (Título de Referência)

---

## Links de Referência - Sensibilidade a Maiúsculas

[Link sensível a maiúsculas][case1]

[Link Sensível a Maiúsculas][Case1]

[case1]: https://example.com/lowercase
[Case1]: https://example.com/uppercase

---

## Links de Referência - Números

[Link 1][1]

[Link 2][2]

[1]: https://example.com/1
[2]: https://example.com/2

---

## Links de Referência - Caracteres Especiais

[Link com hífen][link-dash]

[Link com sublinhado][link_underscore]

[Link com ponto][link.dot]

[link-dash]: https://example.com/dash
[link_underscore]: https://example.com/underscore
[link.dot]: https://example.com/dot

---

## Autolinks

<https://example.com>

<http://example.com>

<mailto:user@example.com>

<user@example.com>

---

## Links com Ênfase

**[Link negrito](https://example.com)**

*[Link itálico](https://example.com)*

***[Link negrito e itálico](https://example.com)***

[**Texto negrito**](https://example.com) no link

[*Texto itálico*](https://example.com) no link

---

## Links em Listas

* Item com [link](https://example.com)
* Item com [link de referência][listref]
* Item com <https://example.com>

[listref]: https://example.com/list

---

## Links em Citações

> Este é um [link](https://example.com) em uma citação.

> Este é um [link de referência][blockquoteref] em uma citação.

[blockquoteref]: https://example.com/blockquote

---

## Links em Cabeçalhos

# Cabeçalho com [link](https://example.com)

## Cabeçalho com [link de referência][headerref]

[headerref]: https://example.com/header

---

## Links em Tabelas

| Column 1 | Column 2 |
|----------|----------|
| [Link](https://example.com) | [Referência][tableref] |
| <https://example.com> | [Link com título](https://example.com "Título") |

[tableref]: https://example.com/table

---

## Links com Imagens

[![Imagem](image.png)](https://example.com)

[![Imagem com alt](image.png "Título")](https://example.com "Título do Link")

---

## Links Relativos

[Link relativo](./relative.html)

[Diretório pai](../parent.html)

[Arquivo irmão](../sibling.html)

[Relativo à raiz](/root.html)

---

## Links de Âncora

[Link para seção](#section)

[Link para subseção](#subsection-name)

[Link com espaços](#section name)

[Link com caracteres especiais](#section-name_123)

---

## Links com Código

[Link com `código`](https://example.com)

`[Link de código](https://example.com)` (não deveria ser um link)

---

## Links com HTML

<a href="https://example.com">Link HTML</a>

[Link Markdown](https://example.com) e <a href="https://example.com">Link HTML</a>

---

## Links Quebrados

[Link quebrado](https://nonexistent.example.com)

[Referência quebrada][broken]

---

## Links com Caracteres Escapados

\[Link escapado\]\(https://example.com\)

[Link com \*asterisco\*](https://example.com)

[Link com \_sublinhado\_](https://example.com)

---

## Links com Unicode

[Link com café](https://example.com/café)

[Link com 你好](https://example.com/你好)

[Link com こんにちは](https://example.com/こんにちは)

---

## Múltiplos Links em Um Parágrafo

Este parágrafo tem [primeiro link](https://example.com/1), [segundo link](https://example.com/2) e [terceiro link](https://example.com/3).

Este parágrafo mistura [link inline](https://example.com/inline), [link de referência][multiref] e <https://example.com/autolink>.

[multiref]: https://example.com/reference

---

## Links com URLs Longas

[Link com URL muito longa](https://example.com/very/long/path/to/resource/that/spans/multiple/segments/and/includes/query/parameters?param1=value1&param2=value2&param3=value3#and-even-a-fragment)

---

## Links com Texto Vazio

[](https://example.com)

[ ](https://example.com)

---

## Links Apenas com Espaços

[   ](https://example.com)

---

## Links em Blocos de Código

```
[Link](https://example.com) não deveria funcionar em blocos de código
```

---

## Links com Esquemas de URL Especiais

[javascript:alert('XSS')](javascript:alert('XSS'))

[data:text/plain,Hello](data:text/plain,Hello)

[about:blank](about:blank)

---

## Links com Codificação Percentual

[Link com codificação](https://example.com/path%20with%20spaces)

[Link com codificação unicode](https://example.com/%E4%BD%A0%E5%A5%BD)

---

## Combinações Complexas de Links

**[Link negrito](https://example.com)** com texto *itálico* e [outro link](https://example.com/2).

[Link com **negrito** e *itálico*](https://example.com) no texto.

[Link](https://example.com) com `código` e **negrito** e *itálico*.

