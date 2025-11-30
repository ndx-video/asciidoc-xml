---
title: Imágenes - Estilo Inline y Referencia
categories:
  - markdown
  - images
tags:
  - graphics
  - media
  - pictures
---

# Imágenes

## Imágenes Inline

![Texto alternativo](image.png)

![Texto alternativo con espacios](image with spaces.png)

![Texto alternativo](image.png "Título de Imagen")

![Texto alternativo](image.png 'Título de Imagen')

![Texto alternativo](image.png (Título de Imagen))

---

## Imágenes con Diferentes Formatos

![Imagen PNG](image.png)

![Imagen JPEG](image.jpg)

![Imagen GIF](image.gif)

![Imagen SVG](image.svg)

![Imagen WebP](image.webp)

---

## Imágenes con Rutas

![Imagen raíz](/image.png)

![Imagen subdirectorio](./images/image.png)

![Directorio padre](../images/image.png)

![Ruta profunda](./images/subfolder/image.png)

---

## Imágenes con URLs

![Imagen remota](https://example.com/image.png)

![Imagen remota con ruta](https://example.com/images/photo.jpg)

![Imagen remota con consulta](https://example.com/image.png?size=large)

---

## Imágenes Estilo Referencia

![Imagen de referencia][img1]

![Imagen de referencia implícita][]

![Imagen de referencia corta]

[img1]: image1.png
[implicit reference image]: image2.png
[shortcut reference image]: image3.png

---

## Imágenes de Referencia con Títulos

![Imagen de referencia con título][img2]

[img2]: image.png "Título de Imagen"

![Imagen de referencia con comillas simples][img3]

[img3]: image.png 'Título de Imagen'

![Imagen de referencia con paréntesis][img4]

[img4]: image.png (Título de Imagen)

---

## Imágenes con Caracteres Especiales en Texto Alternativo

![Imagen con "comillas"](image.png)

![Imagen con (paréntesis)](image.png)

![Imagen con [corchetes]](image.png)

![Imagen con {llaves}](image.png)

![Imagen con <etiquetas>](image.png)

![Imagen con $dólar$](image.png)

---

## Imágenes con Unicode en Texto Alternativo

![Imagen con café](image.png)

![Imagen con 你好](image.png)

![Imagen con こんにちは](image.png)

![Imagen con مرحبا](image.png)

---

## Imágenes con Números en Texto Alternativo

![Imagen 123](image.png)

![Imagen v1.2.3](image.png)

![Imagen 50%](image.png)

---

## Imágenes con Énfasis en Texto Alternativo

![**Texto alternativo negrita**](image.png)

![*Texto alternativo cursiva*](image.png)

![***Texto alternativo negrita y cursiva***](image.png)

---

## Imágenes con Enlaces

[![Imagen enlazada](image.png)](https://example.com)

[![Imagen enlazada con alt](image.png "Título de Imagen")](https://example.com "Título del Enlace")

---

## Imágenes en Listas

* Elemento con ![imagen](image.png)
* Elemento con ![imagen de referencia][listimg]
* Elemento con ![imagen](image.png "Título")

[listimg]: image.png

---

## Imágenes en Citas

> Esta es una cita con ![imagen](image.png)

> Esta es una cita con ![imagen de referencia][blockquoteimg]

[blockquoteimg]: image.png

---

## Imágenes en Encabezados

# Encabezado con ![imagen](image.png)

## Encabezado con ![imagen de referencia][headerimg]

[headerimg]: image.png

---

## Imágenes en Tablas

| Column 1 | Column 2 |
|----------|----------|
| ![Imagen](image.png) | ![Referencia][tableimg] |
| Texto | ![Imagen con título](image.png "Título") |

[tableimg]: image.png

---

## Imágenes con Texto Alternativo Vacío

![](image.png)

![ ](image.png)

![   ](image.png)

---

## Imágenes Solo con Espacios en Texto Alternativo

![   ](image.png)

---

## Imágenes con Texto Alternativo Largo

![Este es un texto alternativo muy largo que describe la imagen en gran detalle y proporciona información completa sobre lo que contiene la imagen y su propósito en el documento](image.png)

---

## Múltiples Imágenes en Un Párrafo

Este párrafo tiene ![primera imagen](image1.png), ![segunda imagen](image2.png) y ![tercera imagen](image3.png).

Este párrafo mezcla ![imagen inline](image.png), ![imagen de referencia][multimg] y texto.

[multimg]: image.png

---

## Imágenes con HTML

<img src="image.png" alt="Imagen HTML">

![Imagen Markdown](image.png) y <img src="image.png" alt="Imagen HTML">

---

## Imágenes con Caracteres Escapados

\![Imagen escapada](image.png)

![Imagen con \*asterisco\*](image.png)

![Imagen con \_guion bajo\_](image.png)

---

## Imágenes en Bloques de Código

```
![Imagen](image.png) no debería funcionar en bloques de código
```

---

## Imágenes con Esquemas de URL Especiales

![Imagen Data URI](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==)

---

## Imágenes con Codificación Porcentual

![Imagen con codificación](https://example.com/image%20with%20spaces.png)

![Imagen con codificación unicode](https://example.com/%E4%BD%A0%E5%A5%BD.png)

---

## Combinaciones Complejas de Imágenes

**Texto negrita** con ![imagen](image.png) y texto *cursiva*.

![Imagen](image.png) con [enlace](https://example.com) y `código`.

![**Alt negrita**](image.png) con texto *cursiva* y [enlace](https://example.com).

---

## Imágenes con Diferentes Tamaños (atributos HTML)

<img src="image.png" alt="Imagen pequeña" width="100" height="100">

<img src="image.png" alt="Imagen grande" width="800" height="600">

---

## Imágenes con Clases CSS (HTML)

<img src="image.png" alt="Imagen con estilo" class="custom-class">

<img src="image.png" alt="Imagen centrada" style="display: block; margin: 0 auto;">

---

## Imágenes con Títulos y Subtítulos

![Imagen](image.png "Este es un título")

![Imagen](image.png "Título con 'comillas'")

![Imagen](image.png "Título con (paréntesis)")

---

## Imágenes con Rutas Relativas

![Mismo directorio](./image.png)

![Subdirectorio](./images/image.png)

![Directorio padre](../image.png)

![Relativo a raíz](/image.png)

---

## Imágenes con URLs Absolutas

![URL absoluta](https://example.com/image.png)

![URL absoluta con ruta](https://example.com/images/photo.jpg)

![URL absoluta con consulta](https://example.com/image.png?v=1&size=large)

![URL absoluta con fragmento](https://example.com/image.png#section)

---

## Imágenes en Estructuras Anidadas

### Imágenes en Listas Anidadas

* Primer nivel
  * Segundo nivel con ![imagen](image.png)
    * Tercer nivel con ![imagen](image.png)

### Imágenes en Citas Anidadas

> Cita externa
> > Cita interna con ![imagen](image.png)

---

## Imágenes con Referencias Rotas

![Imagen rota](nonexistent.png)

![Referencia rota][broken]

[broken]: nonexistent.png

---

## Imágenes con Nombres de Archivo Especiales

![Imagen con espacios](image with spaces.png)

![Imagen.con.puntos](image.with.dots.png)

![Imagen-con-guiones](image-with-dashes.png)

![Imagen_con_guiones_bajos](Image_with_underscores.png)

![Imagen123](Image123.png)

![Imagen@especial](Image@special.png)

