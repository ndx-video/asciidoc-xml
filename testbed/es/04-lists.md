---
title: Listas - Ordenadas, No Ordenadas y Anidadas
categories:
  - markdown
  - lists
tags:
  - ordered
  - unordered
  - nested
  - task lists
---

# Listas

## Listas No Ordenadas - Asterisco

* Primer elemento
* Segundo elemento
* Tercer elemento

---

## Listas No Ordenadas - Guion

- Primer elemento
- Segundo elemento
- Tercer elemento

---

## Listas No Ordenadas - Más

+ Primer elemento
+ Segundo elemento
+ Tercer elemento

---

## Listas Ordenadas - Números

1. Primer elemento
2. Segundo elemento
3. Tercer elemento

---

## Listas Ordenadas - Números Secuenciales

1. Primer elemento
2. Segundo elemento
3. Tercer elemento
4. Cuarto elemento
5. Quinto elemento

---

## Listas Ordenadas - Números No Secuenciales

1. Primer elemento
5. Segundo elemento (numerado 5)
3. Tercer elemento (numerado 3)
10. Cuarto elemento (numerado 10)

---

## Listas Ordenadas que Comienzan desde Diferentes Números

5. Quinto elemento
6. Sexto elemento
7. Séptimo elemento

---

## Listas No Ordenadas Anidadas

* Elemento de primer nivel
  * Elemento de segundo nivel
    * Elemento de tercer nivel
  * Otro elemento de segundo nivel
* Otro elemento de primer nivel

---

## Listas Ordenadas Anidadas

1. Elemento de primer nivel
   1. Elemento de segundo nivel
      1. Elemento de tercer nivel
   2. Otro elemento de segundo nivel
2. Otro elemento de primer nivel

---

## Listas Anidadas Mixtas

* Elemento no ordenado
  1. Elemento ordenado anidado
  2. Otro elemento ordenado anidado
* Otro elemento no ordenado

1. Elemento ordenado
   * Elemento no ordenado anidado
   * Otro elemento no ordenado anidado
2. Otro elemento ordenado

---

## Listas con Múltiples Párrafos

* Primer elemento

  Este es un segundo párrafo en el primer elemento.

* Segundo elemento

  Este es un segundo párrafo en el segundo elemento.

  Este es un tercer párrafo en el segundo elemento.

---

## Listas con Bloques de Código

* Primer elemento

  ```javascript
  console.log("Bloque de código en elemento de lista");
  ```

* Segundo elemento

---

## Listas con Citas

* Primer elemento

  > Esta es una cita en un elemento de lista

* Segundo elemento

---

## Listas con Énfasis

* Elemento con texto **negrita**
* Elemento con texto *cursiva*
* Elemento con texto ***negrita y cursiva***
* Elemento con `código` inline

---

## Listas con Enlaces

* Elemento con [enlace](https://example.com)
* Elemento con [enlace de referencia][ref]
* Elemento con [enlace](https://example.com "título")

[ref]: https://example.com/reference

---

## Listas con Imágenes

* Elemento con ![imagen](image.png)
* Elemento con ![imagen con alt](image.png "título")
* Elemento con ![imagen de referencia][imgref]

[imgref]: image.png "Imagen de referencia"

---

## Listas de Tareas (GitHub Flavored Markdown)

- [ ] Tarea no marcada
- [x] Tarea marcada
- [X] Tarea marcada (mayúscula)
- [ ] Otra tarea no marcada
- [x] Otra tarea marcada

---

## Listas de Tareas Anidadas

- [ ] Tarea principal
  - [ ] Subtarea 1
  - [x] Subtarea 2 (completada)
  - [ ] Subtarea 3
- [x] Otra tarea principal (completada)
  - [x] Subtarea completada
  - [ ] Subtarea incompleta

---

## Listas de Tareas con Contenido

- [ ] Tarea no marcada con texto **negrita**
- [x] Tarea marcada con texto *cursiva*
- [ ] Tarea con [enlace](https://example.com)
- [x] Tarea con `código` inline

---

## Listas con HTML

* Elemento con <strong>HTML negrita</strong>
* Elemento con <em>HTML cursiva</em>
* Elemento con <code>HTML código</code>

---

## Listas con Caracteres Especiales

* Elemento con "comillas"
* Elemento con (paréntesis)
* Elemento con [corchetes]
* Elemento con {llaves}
* Elemento con <etiquetas>
* Elemento con $dólar$
* Elemento con &ampersand&

---

## Listas con Números

* Elemento con números 123
* Elemento con versión v1.2.3
* Elemento con porcentaje 50%
* Elemento con número grande 1,000,000

---

## Listas con Unicode

* Elemento con café
* Elemento con naïve
* Elemento con 你好
* Elemento con こんにちは
* Elemento con مرحبا

---

## Listas Apretadas (Sin Líneas en Blanco)

* Primer elemento
* Segundo elemento
* Tercer elemento

1. Primer elemento
2. Segundo elemento
3. Tercer elemento

---

## Listas Sueltas (Con Líneas en Blanco)

* Primer elemento

* Segundo elemento

* Tercer elemento

1. Primer elemento

2. Segundo elemento

3. Tercer elemento

---

## Listas que Comienzan a Mitad del Documento

Algún texto de párrafo antes de la lista.

* Elemento de lista 1
* Elemento de lista 2

Algún texto de párrafo después de la lista.

---

## Listas que Terminan el Documento

* Elemento final 1
* Elemento final 2
* Elemento final 3

---

## Elementos de Lista Vacíos

* 
* Elemento con contenido
* 
* Otro elemento

---

## Listas con Contenido Largo

* Este es un elemento de lista muy largo que contiene múltiples oraciones y continúa por un tiempo para probar cómo el convertidor maneja contenido más largo dentro de los elementos de lista. Debe mantener el formato y la estructura adecuados.

* Otro elemento largo que abarca múltiples líneas en la fuente pero debe renderizarse como un solo párrafo dentro de la estructura del elemento de lista.

---

## Listas Anidadas Complejas

1. Primer elemento ordenado
   * Elemento no ordenado anidado
     - Elemento anidado profundo
       + Elemento anidado más profundo
   * Otro elemento no ordenado anidado
2. Segundo elemento ordenado
   * Elemento no ordenado anidado
     1. Elemento ordenado anidado
     2. Otro elemento ordenado anidado
3. Tercer elemento ordenado

---

## Listas con Reglas Horizontales

* Elemento antes de la regla

---

* Elemento después de la regla

---

## Listas con Encabezados

* Elemento antes del encabezado

## Encabezado en Contexto de Lista

* Elemento después del encabezado

---

## Listas de Definición (si se admite)

Término 1
: Definición 1

Término 2
: Definición 2
: Definición alternativa 2

Término 3
: Definición 3 con **negrita** y *cursiva*

