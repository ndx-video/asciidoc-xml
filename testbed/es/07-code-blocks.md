---
title: Bloques de Código y Código Inline
categories:
  - markdown
  - code
tags:
  - programming
  - syntax highlighting
  - code blocks
---

# Bloques de Código y Código Inline

## Código Inline

Este párrafo contiene `inline code` usando comillas invertidas.

Este párrafo tiene `code with spaces` en él.

Este párrafo tiene `code_with_underscores` y `code-with-dashes`.

Este párrafo tiene `code123` con números.

Este párrafo tiene `code.with.dots` con puntos.

---

## Código Inline con Caracteres Especiales

Este párrafo tiene `code with "quotes"` dentro.

Este párrafo tiene `code with (parentheses)` dentro.

Este párrafo tiene `code with [brackets]` dentro.

Este párrafo tiene `code with {braces}` dentro.

Este párrafo tiene `code with <tags>` dentro.

Este párrafo tiene `code with $dollar$` dentro.

Este párrafo tiene `code with &ampersand&` dentro.

---

## Código Inline con Comillas Invertidas

Este párrafo tiene ``code with `backtick` inside`` usando dobles comillas invertidas.

Este párrafo tiene ```code with ``double backticks`` inside``` usando triples comillas invertidas.

---

## Casos Extremos de Código Inline

`code`text (sin espacio)

`code`**bold** (formato adyacente)

`code`*italic* (formato adyacente)

`code`[link](https://example.com) (enlace adyacente)

---

## Bloques de Código con Cerca - Sin Lenguaje

```
This is a code block
with multiple lines
of code
```

---

## Bloques de Código con Cerca - JavaScript

```javascript
function greet(name) {
    console.log("Hello, " + name + "!");
}

greet("World");
```

---

## Bloques de Código con Cerca - Python

```python
def greet(name):
    print(f"Hello, {name}!")

greet("World")
```

---

## Bloques de Código con Cerca - Go

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

---

## Bloques de Código con Cerca - HTML

```html
<!DOCTYPE html>
<html>
<head>
    <title>Example</title>
</head>
<body>
    <h1>Hello, World!</h1>
</body>
</html>
```

---

## Bloques de Código con Cerca - CSS

```css
body {
    font-family: Arial, sans-serif;
    margin: 0;
    padding: 20px;
}

h1 {
    color: #333;
}
```

---

## Bloques de Código con Cerca - JSON

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

## Bloques de Código con Cerca - YAML

```yaml
name: Example
version: 1.0.0
dependencies:
  package: ^1.2.3
```

---

## Bloques de Código con Cerca - Shell

```bash
#!/bin/bash
echo "Hello, World!"
ls -la
```

---

## Bloques de Código con Cerca - SQL

```sql
SELECT * FROM users
WHERE age > 18
ORDER BY name;
```

---

## Bloques de Código con Cerca - XML

```xml
<?xml version="1.0" encoding="UTF-8"?>
<root>
    <element attribute="value">Content</element>
</root>
```

---

## Bloques de Código con Cerca - Markdown

```markdown
# Header

**Bold** and *italic* text.

- List item
```

---

## Bloques de Código con Cerca - Múltiples Lenguajes

```javascript
// JavaScript code
```

```python
# Python code
```

```go
// Go code
```

---

## Bloques de Código con Sangría

    This is an indented code block
    with multiple lines
    using four spaces

---

## Bloques de Código con Líneas Vacías

```javascript
function example() {
    // First line
    
    // Third line after empty line
}
```

---

## Bloques de Código con Caracteres Especiales

```javascript
const special = {
    "quotes": 'single quotes',
    "parentheses": (value),
    "brackets": [array],
    "braces": {object},
    "tags": <html>,
    "dollar": $value$,
    "ampersand": &value&
};
```

---

## Bloques de Código con Unicode

```javascript
const unicode = {
    "café": "café",
    "你好": "你好",
    "こんにちは": "こんにちは"
};
```

---

## Bloques de Código con Comillas Invertidas Dentro

````markdown
```
Code block with backticks
```
````

---

## Bloques de Código con Tres Comillas Invertidas Dentro

`````markdown
```
Code block content
```
`````

---

## Bloques de Código en Listas

* Elemento de lista con bloque de código:

  ```javascript
  console.log("Code in list");
  ```

* Otro elemento

---

## Bloques de Código en Citas

> Cita con bloque de código:
> 
> ```javascript
> console.log("Code in blockquote");
> ```

---

## Bloques de Código con Números de Línea (si se admite)

```javascript:1-5
function example() {
    return "Hello";
}
```

---

## Bloques de Código con Nombres de Archivo (si se admite)

```javascript:example.js
function example() {
    return "Hello";
}
```

---

## Bloques de Código con Metadatos

```javascript
// File: example.js
// Author: Test Author
// Date: 2024-01-15
function example() {
    return "Hello";
}
```

---

## Bloques de Código con Comentarios

```javascript
// Single line comment
function example() {
    /* Multi-line
       comment */
    return "Hello";
}
```

---

## Bloques de Código con Cadenas

```javascript
const single = 'Single quoted string';
const double = "Double quoted string";
const template = `Template string with ${variable}`;
```

---

## Bloques de Código con Expresiones Regulares

```javascript
const regex1 = /pattern/g;
const regex2 = /pattern with \(escaped\)/g;
const regex3 = new RegExp("pattern");
```

---

## Bloques de Código con Entidades HTML

```html
&lt;div&gt;Content&lt;/div&gt;
&amp;copy; 2024
&quot;Quoted text&quot;
```

---

## Bloques de Código con Caracteres Escapados

```javascript
const escaped = "String with \\n newline";
const tab = "String with \\t tab";
const quote = "String with \\\" quote";
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

## Bloques de Código Adyacentes

```javascript
// First block
```

```python
# Second block
```

---

## Bloques de Código con Solo Espacios en Blanco

```
    
```

---

## Bloques de Código con Solo Saltos de Línea

```

```

---

## Bloques de Código con Contenido Mixto

```javascript
function mixed() {
    // Comments
    const code = "code";
    /* Multi-line
       comment */
    return code;
}
```

---

## Bloques de Código con Líneas Muy Largas

```javascript
const veryLongLine = "This is a very long line of code that extends beyond the normal width and tests how the converter handles long lines in code blocks without breaking the formatting or structure";
```

---

## Bloques de Código con Sangría

```python
def nested():
    if True:
        if True:
            return "Deeply nested"
```

---

## Bloques de Código con Tabulaciones

```javascript
function withTabs() {
	return "Tab indented";
}
```

---

## Bloques de Código con Espacios Finales

```javascript
function example() {
    return "Line with trailing spaces    ";
}
```

---

## Bloques de Código con Caracteres Markdown Especiales

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
    const obj = {
        key1: "value1",
        key2: "value2",
        key3: {
            nested: "value"
        }
    };
    return obj;
}
```

---

## Código Inline en Diferentes Contextos

Este párrafo tiene `code` inline.

**Bold** con `code` inline.

*Italic* con `code` inline.

[Link](https://example.com) con `code` inline.

![Image](image.png) con `code` inline.

---

## Escape de Código Inline

\`Escaped backtick\` no debería ser código.

`code with \`escaped backtick\`` dentro.

---

## Bloques de Código con Alias de Lenguaje

```js
// JavaScript using 'js' alias
```

```py
# Python using 'py' alias
```

```sh
# Shell using 'sh' alias
```

---

## Bloques de Código con Lenguaje Inválido

```invalidlanguage
This is code with an invalid language identifier
```

---

## Bloques de Código con Lenguaje Vacío

```
This is a code block with empty language identifier
```

