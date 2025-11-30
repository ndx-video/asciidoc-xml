---
title: Variaciones de Frontmatter
categories:
  - frontmatter
  - yaml
tags:
  - metadata
  - frontmatter
  - yaml
author: Autor de Prueba
date: 2024-01-15
published: true
---

# Variaciones de Frontmatter

Este documento prueba varios formatos y estructuras de frontmatter.

---

## Frontmatter Simple

```yaml
---
title: Simple Document
author: Author Name
---
```

---

## Frontmatter con Arrays

```yaml
---
title: Document with Arrays
categories:
  - category1
  - category2
  - category3
tags:
  - tag1
  - tag2
---
```

---

## Frontmatter con Estructuras Anidadas

```yaml
---
title: Document with Nested Structure
metadata:
  author: Author Name
  date: 2024-01-15
  version: 1.0.0
---
```

---

## Frontmatter con Arrays Complejos

```yaml
---
title: Complex Arrays
categories:
  - microsoft
  - Windows
  - windows server
tags:
  - 2012 R2
  - autounattend.xml
  - packer
wordtwit_post_info:
  - 'O:8:"stdClass":14:{s:6:"manual";b:1;...}'
---
```

---

## Frontmatter con Valores Vacíos

```yaml
---
title: Empty Values
author:
description:
tags: []
---
```

---

## Frontmatter con Valores entre Comillas

```yaml
---
title: "Quoted Title"
author: 'Single Quoted Author'
description: "Description with \"quotes\" inside"
---
```

---

## Frontmatter con Caracteres Especiales

```yaml
---
title: Title with (parentheses) and [brackets]
author: Author with &ampersand&
description: Description with $dollar$ and %percent%
---
```

---

## Frontmatter con Números

```yaml
---
title: Title with Numbers
version: 1.2.3
year: 2024
percentage: 50%
---
```

---

## Frontmatter con Fechas

```yaml
---
title: Document with Dates
date: 2024-01-15
published_date: 2024-01-15T10:30:00Z
updated: 2024-12-31
---
```

---

## Frontmatter con Booleanos

```yaml
---
title: Document with Booleans
published: true
draft: false
featured: true
---
```

---

## Frontmatter con Valores Null

```yaml
---
title: Document with Null
author: null
description: null
tags: null
---
```

---

## Frontmatter con Tipos Mixtos

```yaml
---
title: Mixed Types Document
author: Author Name
version: 1.0.0
published: true
date: 2024-01-15
tags:
  - tag1
  - tag2
---
```

---

## Frontmatter con Valores Largos

```yaml
---
title: Document with Long Values
description: This is a very long description that spans multiple lines and contains extensive information about the document content and purpose. It tests how frontmatter handles longer string values that might wrap or extend beyond normal line lengths.
---
```

---

## Frontmatter con Unicode

```yaml
---
title: Document with Unicode
author: Author with café
description: Description with 你好 and こんにちは
tags:
  - unicode
  - 测试
---
```

---

## Frontmatter con URLs

```yaml
---
title: Document with URLs
url: https://example.com/document
image: https://example.com/image.png
canonical: https://example.com/canonical
---
```

---

## Frontmatter con Email

```yaml
---
title: Document with Email
author: Author Name
email: author@example.com
contact: contact@example.com
---
```

---

## Frontmatter con Rutas de Archivo

```yaml
---
title: Document with Paths
source: ./source/file.md
output: ./output/file.html
template: templates/default.html
---
```

---

## Frontmatter con Arrays de Objetos

```yaml
---
title: Arrays of Objects
authors:
  - name: Author One
    email: one@example.com
  - name: Author Two
    email: two@example.com
---
```

---

## Frontmatter con Anidamiento Profundo

```yaml
---
title: Deep Nesting
metadata:
  author:
    name: Author Name
    contact:
      email: author@example.com
      phone: +1-234-567-8900
  publishing:
    date: 2024-01-15
    status: published
---
```

---

## Frontmatter con Múltiples Arrays

```yaml
---
title: Multiple Arrays
categories:
  - tech
  - programming
tags:
  - go
  - markdown
keywords:
  - keyword1
  - keyword2
  - keyword3
---
```

---

## Frontmatter con Arrays Vacíos

```yaml
---
title: Empty Arrays
categories: []
tags: []
keywords: []
---
```

---

## Frontmatter con Arrays de Un Solo Elemento

```yaml
---
title: Single Item Arrays
categories:
  - single
tags:
  - only
---
```

---

## Frontmatter con Formatos de Array Mixtos

```yaml
---
title: Mixed Array Formats
inline_array: [item1, item2, item3]
multiline_array:
  - item1
  - item2
  - item3
---
```

---

## Frontmatter con Comentarios (si se admite)

```yaml
---
# This is a comment
title: Document with Comments
# Another comment
author: Author Name
---
```

---

## Frontmatter con Anclas y Alias (si se admite)

```yaml
---
title: Anchors and Aliases
author: &author
  name: Author Name
  email: author@example.com
coauthor: *author
---
```

---

## Frontmatter con Cadenas Multilínea

```yaml
---
title: Multiline Strings
description: |
  This is a multiline description
  that spans multiple lines
  and preserves line breaks
---
```

---

## Frontmatter con Cadenas Plegadas

```yaml
---
title: Folded Strings
description: >
  This is a folded string
  that converts line breaks
  to spaces
---
```

---

## Frontmatter con Escalares de Bloque Literal

```yaml
---
title: Literal Block Scalars
code: |
  function example() {
      return "code";
  }
---
```

---

## Frontmatter con Campos Personalizados

```yaml
---
title: Custom Fields
custom_field_1: Custom Value 1
custom_field_2: Custom Value 2
amazon-product-content-location: location
amazon-product-excerpt-hook-override: override
amazon-product-content-hook-override: override
amazon-product-newwindow: true
wordtwit_posted_tweets:
  - tweet1
  - tweet2
---
```

---

## Frontmatter con Números de Versión

```yaml
---
title: Version Numbers
version: 1.0.0
api_version: 2.1.3
semver: 1.2.3-beta.1
---
```

---

## Frontmatter con Moneda

```yaml
---
title: Currency Values
price: $19.99
currency: USD
amount: 100.50
---
```

---

## Frontmatter con Porcentajes

```yaml
---
title: Percentages
completion: 75%
discount: 20%
tax_rate: 8.5%
---
```

---

## Frontmatter con Valores de Tiempo

```yaml
---
title: Time Values
duration: 2h30m
timeout: 5s
interval: 1m
---
```

---

## Frontmatter con Direcciones IP

```yaml
---
title: IP Addresses
ipv4: 192.168.1.1
ipv6: 2001:0db8:85a3:0000:0000:8a2e:0370:7334
---
```

---

## Frontmatter con Control de Versiones

```yaml
---
title: Version Control
git_hash: abc123def456
branch: main
commit_date: 2024-01-15
---
```

---

## Frontmatter con Redes Sociales

```yaml
---
title: Social Media
twitter: @username
github: github.com/username
linkedin: linkedin.com/in/username
---
```

---

## Frontmatter con Datos Geográficos

```yaml
---
title: Geographic Data
latitude: 40.7128
longitude: -74.0060
location: New York, NY
---
```

---

## Frontmatter con Calificaciones

```yaml
---
title: Ratings
rating: 4.5
stars: 5
score: 8.5/10
---
```

---

## Frontmatter con Arrays Anidados Complejos

```yaml
---
title: Complex Nested Arrays
data:
  - name: Item 1
    values:
      - value1
      - value2
  - name: Item 2
    values:
      - value3
      - value4
---
```

---

## Frontmatter con Caracteres Escapados

```yaml
---
title: Escaped Characters
description: "Description with \"quotes\""
path: "C:\\Users\\Name\\File"
---
```

---

## Frontmatter con Estructuras Tipo JSON

```yaml
---
title: JSON-like Structures
config: '{"key": "value", "number": 123}'
settings: '{"enabled": true, "count": 5}'
---
```

---

## Frontmatter con Base64 (si es necesario)

```yaml
---
title: Base64 Data
encoded: SGVsbG8gV29ybGQ=
---
```

---

## Frontmatter con HTML

```yaml
---
title: HTML in Frontmatter
description: "<strong>Bold</strong> and <em>italic</em>"
---
```

---

## Frontmatter con Markdown

```yaml
---
title: Markdown in Frontmatter
description: "**Bold** and *italic* text"
---
```

---

## Frontmatter con Nombres de Campo Muy Largos

```yaml
---
title: Very Long Field Names
very_long_field_name_that_extends_beyond_normal_width: value
another_extremely_long_field_name_with_many_words: another value
---
```

---

## Frontmatter con Características Especiales de YAML

```yaml
---
title: Special YAML Features
null_value: null
boolean_true: true
boolean_false: false
integer: 123
float: 123.456
string: "string"
---
```

---

## Frontmatter que Termina con Tres Puntos

```yaml
---
title: Three Dots Ending
author: Author Name
...
```

---

## Frontmatter Estilo TOML (si se admite)

```toml
+++
title = "TOML Frontmatter"
author = "Author Name"
+++
```

---

## Frontmatter con JSON (si se admite)

```json
{
  "title": "JSON Frontmatter",
  "author": "Author Name"
}
```

---

## Frontmatter con Múltiples Separadores de Documento

```yaml
---
title: Multiple Separators
author: Author Name
---
---
second_section: value
---
```

---

## Frontmatter con Arrays Inline

```yaml
---
title: Inline Arrays
tags: [tag1, tag2, tag3]
categories: [cat1, cat2]
---
```

---

## Frontmatter con Comillas Mixtas

```yaml
---
title: 'Single quoted title'
description: "Double quoted description"
author: Unquoted author
---
```

---

## Frontmatter Preservando Orden

```yaml
---
z_field: Last
a_field: First
m_field: Middle
---
```

