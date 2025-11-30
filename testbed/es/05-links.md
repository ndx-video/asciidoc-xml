---
title: Enlaces - Inline, Referencia y Autolinks
categories:
  - markdown
  - links
tags:
  - hyperlinks
  - references
  - urls
---

# Enlaces

## Enlaces Inline

Este es un [enlace inline](https://example.com).

Este es un [enlace con título](https://example.com "Título del Enlace").

Este es un [enlace con comillas simples](https://example.com 'Título del Enlace').

Este es un [enlace con paréntesis](https://example.com (Título del Enlace)).

---

## Enlaces con Diferentes Protocolos

[Enlace HTTP](http://example.com)

[Enlace HTTPS](https://example.com)

[Enlace FTP](ftp://ftp.example.com)

[Enlace de correo](mailto:user@example.com)

[Enlace de archivo](file:///path/to/file)

---

## Enlaces con Rutas

[Ruta raíz](https://example.com/)

[Subdirectorio](https://example.com/path/to/resource)

[Extensión de archivo](https://example.com/file.html)

[Cadena de consulta](https://example.com/page?param=value&other=123)

[Fragmento](https://example.com/page#section)

[URL compleja](https://example.com/path/to/resource?query=value&other=123#fragment)

---

## Enlaces con Caracteres Especiales

[Enlace con espacios](https://example.com/path with spaces)

[Enlace con unicode](https://example.com/你好)

[Enlace con codificación](https://example.com/path%20with%20spaces)

---

## Enlaces Estilo Referencia

Este es un [enlace de referencia][ref1].

Este es un [enlace de referencia implícito][].

Este es un [enlace de referencia corto].

[ref1]: https://example.com/reference1
[implicit reference link]: https://example.com/implicit
[shortcut reference link]: https://example.com/shortcut

---

## Enlaces de Referencia con Títulos

Este es un [enlace de referencia con título][ref2].

[ref2]: https://example.com/reference2 "Título de Referencia"

Este es un [enlace de referencia con comillas simples][ref3].

[ref3]: https://example.com/reference3 'Título de Referencia'

Este es un [enlace de referencia con paréntesis][ref4].

[ref4]: https://example.com/reference4 (Título de Referencia)

---

## Enlaces de Referencia - Sensibilidad a Mayúsculas

[Enlace sensible a mayúsculas][case1]

[Enlace Sensible a Mayúsculas][Case1]

[case1]: https://example.com/lowercase
[Case1]: https://example.com/uppercase

---

## Enlaces de Referencia - Números

[Enlace 1][1]

[Enlace 2][2]

[1]: https://example.com/1
[2]: https://example.com/2

---

## Enlaces de Referencia - Caracteres Especiales

[Enlace con guion][link-dash]

[Enlace con guion bajo][link_underscore]

[Enlace con punto][link.dot]

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

## Enlaces con Énfasis

**[Enlace negrita](https://example.com)**

*[Enlace cursiva](https://example.com)*

***[Enlace negrita y cursiva](https://example.com)***

[**Texto negrita**](https://example.com) en enlace

[*Texto cursiva*](https://example.com) en enlace

---

## Enlaces en Listas

* Elemento con [enlace](https://example.com)
* Elemento con [enlace de referencia][listref]
* Elemento con <https://example.com>

[listref]: https://example.com/list

---

## Enlaces en Citas

> Este es un [enlace](https://example.com) en una cita.

> Este es un [enlace de referencia][blockquoteref] en una cita.

[blockquoteref]: https://example.com/blockquote

---

## Enlaces en Encabezados

# Encabezado con [enlace](https://example.com)

## Encabezado con [enlace de referencia][headerref]

[headerref]: https://example.com/header

---

## Enlaces en Tablas

| Column 1 | Column 2 |
|----------|----------|
| [Enlace](https://example.com) | [Referencia][tableref] |
| <https://example.com> | [Enlace con título](https://example.com "Título") |

[tableref]: https://example.com/table

---

## Enlaces con Imágenes

[![Imagen](image.png)](https://example.com)

[![Imagen con alt](image.png "Título")](https://example.com "Título del Enlace")

---

## Enlaces Relativos

[Enlace relativo](./relative.html)

[Directorio padre](../parent.html)

[Archivo hermano](../sibling.html)

[Relativo a raíz](/root.html)

---

## Enlaces de Anclaje

[Enlace a sección](#section)

[Enlace a subsección](#subsection-name)

[Enlace con espacios](#section name)

[Enlace con caracteres especiales](#section-name_123)

---

## Enlaces con Código

[Enlace con `código`](https://example.com)

`[Enlace de código](https://example.com)` (no debería ser un enlace)

---

## Enlaces con HTML

<a href="https://example.com">Enlace HTML</a>

[Enlace Markdown](https://example.com) y <a href="https://example.com">Enlace HTML</a>

---

## Enlaces Rotos

[Enlace roto](https://nonexistent.example.com)

[Referencia rota][broken]

---

## Enlaces con Caracteres Escapados

\[Enlace escapado\]\(https://example.com\)

[Enlace con \*asterisco\*](https://example.com)

[Enlace con \_guion bajo\_](https://example.com)

---

## Enlaces con Unicode

[Enlace con café](https://example.com/café)

[Enlace con 你好](https://example.com/你好)

[Enlace con こんにちは](https://example.com/こんにちは)

---

## Múltiples Enlaces en Un Párrafo

Este párrafo tiene [primer enlace](https://example.com/1), [segundo enlace](https://example.com/2) y [tercer enlace](https://example.com/3).

Este párrafo mezcla [enlace inline](https://example.com/inline), [enlace de referencia][multiref] y <https://example.com/autolink>.

[multiref]: https://example.com/reference

---

## Enlaces con URLs Largas

[Enlace con URL muy larga](https://example.com/very/long/path/to/resource/that/spans/multiple/segments/and/includes/query/parameters?param1=value1&param2=value2&param3=value3#and-even-a-fragment)

---

## Enlaces con Texto Vacío

[](https://example.com)

[ ](https://example.com)

---

## Enlaces Solo con Espacios

[   ](https://example.com)

---

## Enlaces en Bloques de Código

```
[Enlace](https://example.com) no debería funcionar en bloques de código
```

---

## Enlaces con Esquemas de URL Especiales

[javascript:alert('XSS')](javascript:alert('XSS'))

[data:text/plain,Hello](data:text/plain,Hello)

[about:blank](about:blank)

---

## Enlaces con Codificación Porcentual

[Enlace con codificación](https://example.com/path%20with%20spaces)

[Enlace con codificación unicode](https://example.com/%E4%BD%A0%E5%A5%BD)

---

## Combinaciones Complejas de Enlaces

**[Enlace negrita](https://example.com)** con texto *cursiva* y [otro enlace](https://example.com/2).

[Enlace con **negrita** y *cursiva*](https://example.com) en el texto.

[Enlace](https://example.com) con `código` y **negrita** y *cursiva*.

