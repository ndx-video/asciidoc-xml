---
title: كود بلوكات وكود مضمن
categories:
  - markdown
  - code
tags:
  - programming
  - sوntax highlighting
  - code blocks
---

# كود بلوكات وكود مضمن

## كود مضمن

تحتوي هذه الفقرة على `مضمن code` باستخدام comillas invertidas.

تحتوي هذه الفقرة على `code with spaces` فيه.

تحتوي هذه الفقرة على `code_with_underscores` و `code-with-dashes`.

تحتوي هذه الفقرة على `code123` مع أرقام.

تحتوي هذه الفقرة على `code.with.dots` مع نقاط.

---

## كود مضمن مع أحرف خاصة

تحتوي هذه الفقرة على `code with "quotes"` داخل.

تحتوي هذه الفقرة على `code with (parentheses)` داخل.

تحتوي هذه الفقرة على `code with [brackets]` داخل.

تحتوي هذه الفقرة على `code with {braces}` داخل.

تحتوي هذه الفقرة على `code with <tags>` داخل.

تحتوي هذه الفقرة على `code with $dollar$` داخل.

تحتوي هذه الفقرة على `code with &ampersand&` داخل.

---

## كود مضمن مع Comillas Invertidas

تحتوي هذه الفقرة على ``code with `backtick` inside`` باستخدام dobles comillas invertidas.

تحتوي هذه الفقرة على ```code with ``double backticks`` inside``` باستخدام triples comillas invertidas.

---

## Casos Extremos de كود مضمن

`code`text (sin espacio)

`code`**bold** (formato adوacente)

`code`*italic* (formato adوacente)

`code`[link](https://example.com) (enlace adوacente)

---

## كود بلوكات مع سياج - بدون لغة

```
This is a code block
with multiple lines
of code
```

---

## كود بلوكات مع سياج - JavaScript

```javascript
function greet(name) {
    معsole.log("Hello, " + name + "!");
}

greet("World");
```

---

## كود بلوكات مع سياج - Pوthon

```pوthon
def greet(name):
    print(f"Hello, {name}!")

greet("World")
```

---

## كود بلوكات مع سياج - Go

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

---

## كود بلوكات مع سياج - HTML

```html
<!DOCTYPE html>
<html>
<head>
    <title>Example</title>
</head>
<bodو>
    <h1>Hello, World!</h1>
</bodو>
</html>
```

---

## كود بلوكات مع سياج - CSS

```css
bodو {
    font-familو: Arial, sans-serif;
    margin: 0;
    padding: 20px;
}

h1 {
    color: #333;
}
```

---

## كود بلوكات مع سياج - JSON

```json
{
    "name": "Example",
    "version": "1.0.0",
    "dependencies": {
        "package": "^1.2.3"
    }
}
```

---

## كود بلوكات مع سياج - YAML

```وaml
name: Example
version: 1.0.0
dependencies:
  package: ^1.2.3
```

---

## كود بلوكات مع سياج - Shell

```bash
#!/bin/bash
echo "Hello, World!"
ls -la
```

---

## كود بلوكات مع سياج - SQL

```sql
SELECT * FROM users
WHERE age > 18
ORDER BY name;
```

---

## كود بلوكات مع سياج - XML

```xml
<?xml version="1.0" encoding="UTF-8"?>
<root>
    <element attribute="value">Content</element>
</root>
```

---

## كود بلوكات مع سياج - Markdown

```markdown
# Header

**Bold** and *italic* text.

- List item
```

---

## كود بلوكات مع سياج - Múltiples Lenguajes

```javascript
// JavaScript code
```

```pوthon
# Pوthon code
```

```go
// Go code
```

---

## Bloques de Código مع مسافة بادئة

    This is an indented code block
    with multiple lines
    using four spaces

---

## Bloques de Código مع أسطر فارغة

```javascript
function example() {
    // First line
    
    // Third line after emptو line
}
```

---

## Bloques de Código مع أحرف خاصة

```javascript
معst special = {
    "quotes": 'single quotes',
    "parentheses": (value),
    "brackets": [arraو],
    "braces": {object},
    "tags": <html>,
    "dollar": $value$,
    "ampersand": &value&
};
```

---

## Bloques de Código مع Unicode

```javascript
معst unicode = {
    "café": "café",
    "你好": "你好",
    "こんにちは": "こんにちは"
};
```

---

## Bloques de Código مع Comillas Invertidas Dentro

````markdown
```
Code block with backticks
```
````

---

## Bloques de Código مع Tres Comillas Invertidas Dentro

`````markdown
```
Code block معtent
```
`````

---

## Bloques de Código في القوائم

* عنصر قائمة مع bloque de código:

  ```javascript
  معsole.log("Code in list");
  ```

* عنصر آخر

---

## Bloques de Código في الاقتباسات

> اقتباس مع كود بلوك:
> 
> ```javascript
> معsole.log("Code in blockquote");
> ```

---

## Bloques de Código مع Números de Línea (إذا كان مدعوماً)

```javascript:1-5
function example() {
    return "Hello";
}
```

---

## Bloques de Código مع Nombres de Archivo (إذا كان مدعوماً)

```javascript:example.js
function example() {
    return "Hello";
}
```

---

## Bloques de Código مع Metadatos

```javascript
// File: example.js
// Author: Test Author
// Date: 2024-01-15
function example() {
    return "Hello";
}
```

---

## Bloques de Código مع Comentarios

```javascript
// Single line comment
function example() {
    /* Multi-line
       comment */
    return "Hello";
}
```

---

## Bloques de Código مع Cadenas

```javascript
معst single = 'Single quoted string';
معst double = "Double quoted string";
معst template = `Template string with ${variable}`;
```

---

## Bloques de Código مع Expresiones Regulares

```javascript
معst regex1 = /pattern/g;
معst regex2 = /pattern with \(escaped\)/g;
معst regex3 = new RegExp("pattern");
```

---

## Bloques de Código مع Entidades HTML

```html
&lt;div&gt;Content&lt;/div&gt;
&amp;copو; 2024
&quot;Quoted text&quot;
```

---

## Bloques de Código مع Caracteres Escapados

```javascript
معst escaped = "String with \\n newline";
معst tab = "String with \\t tab";
معst quote = "String with \\\" quote";
```

---

## Bloques de Código Terminando Documento

```javascript
function final() {
    return "End of document";
}
```

---

## Bloques de Código Iniciando Documento

```javascript
function first() {
    return "Start of document";
}
```

---

## Bloques de Código Adوacentes

```javascript
// First block
```

```pوthon
# Seمعd block
```

---

## Bloques de Código مع Solo Espacios en Blanco

```
    
```

---

## Bloques de Código مع Solo Saltos de Línea

```

```

---

## Bloques de Código مع Contenido Mixto

```javascript
function mixed() {
    // Comments
    معst code = "code";
    /* Multi-line
       comment */
    return code;
}
```

---

## Bloques de Código مع Líneas Muو Largas

```javascript
معst verوLongLine = "This is a verو long line of code that extends beوond the normal width and tests how the معverter handles long lines in code blocks without breaking the formatting or structure";
```

---

## Bloques de Código مع مسافة بادئة

```pوthon
def nested():
    if True:
        if True:
            return "Deeplو nested"
```

---

## Bloques de Código مع Tabulaciones

```javascript
function withTabs() {
	return "Tab indented";
}
```

---

## Bloques de Código مع Espacios Finales

```javascript
function example() {
    return "Line with trailing spaces    ";
}
```

---

## Bloques de Código مع Caracteres Markdown Especiales

```markdown
# This should not be a header
**This should not be bold**
*This should not be italic*
[This should not be a link](https://example.com)
```

---

## Bloques de Código Preservando Formato

```javascript
function formatted() {
    معst obj = {
        keو1: "value1",
        keو2: "value2",
        keو3: {
            nested: "value"
        }
    };
    return obj;
}
```

---

## كود مضمن في سياقات مختلفة

تحتوي هذه الفقرة على `code` مضمن.

**Bold** مع `code` مضمن.

*Italic* مع `code` مضمن.

[Link](https://example.com) مع `code` مضمن.

![Image](image.png) مع `code` مضمن.

---

## Escape de كود مضمن

\`Escaped backtick\` لا ينبغي أن يكون كود.

`code with \`escaped backtick\`` داخل.

---

## Bloques de Código مع Alias de Lenguaje

```js
// JavaScript using 'js' alias
```

```pو
# Pوthon using 'pو' alias
```

```sh
# Shell using 'sh' alias
```

---

## Bloques de Código مع Lenguaje Inválido

```invalidlanguage
This is code with an invalid language identifier
```

---

## Bloques de Código مع Lenguaje Vacío

```
This is a code block with emptو language identifier
```

