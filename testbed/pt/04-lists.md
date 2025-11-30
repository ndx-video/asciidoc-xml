---
title: Listas - Ordenadas, Não Ordenadas e Aninhadas
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

## Listas Não Ordenadas - Asterisco

* Primeiro item
* Segundo item
* Terceiro item

---

## Listas Não Ordenadas - Hífen

- Primeiro item
- Segundo item
- Terceiro item

---

## Listas Não Ordenadas - Mais

+ Primeiro item
+ Segundo item
+ Terceiro item

---

## Listas Ordenadas - Números

1. Primeiro item
2. Segundo item
3. Terceiro item

---

## Listas Ordenadas - Números Sequenciais

1. Primeiro item
2. Segundo item
3. Terceiro item
4. Quarto item
5. Quinto item

---

## Listas Ordenadas - Números Não Sequenciais

1. Primeiro item
5. Segundo item (numerado 5)
3. Terceiro item (numerado 3)
10. Quarto item (numerado 10)

---

## Listas Ordenadas que Começam de Números Diferentes

5. Quinto item
6. Sexto item
7. Sétimo item

---

## Listas Não Ordenadas Aninhadas

* Item de primeiro nível
  * Item de segundo nível
    * Item de terceiro nível
  * Outro item de segundo nível
* Outro item de primeiro nível

---

## Listas Ordenadas Aninhadas

1. Item de primeiro nível
   1. Item de segundo nível
      1. Item de terceiro nível
   2. Outro item de segundo nível
2. Outro item de primeiro nível

---

## Listas Aninhadas Mistas

* Item não ordenado
  1. Item ordenado aninhado
  2. Outro item ordenado aninhado
* Outro item não ordenado

1. Item ordenado
   * Item não ordenado aninhado
   * Outro item não ordenado aninhado
2. Outro item ordenado

---

## Listas com Múltiplos Parágrafos

* Primeiro item

  Este é um segundo parágrafo no primeiro item.

* Segundo item

  Este é um segundo parágrafo no segundo item.

  Este é um terceiro parágrafo no segundo item.

---

## Listas com Blocos de Código

* Primeiro item

  ```javascript
  console.log("Bloco de código em item de lista");
  ```

* Segundo item

---

## Listas com Citações

* Primeiro item

  > Esta é uma citação em um item de lista

* Segundo item

---

## Listas com Ênfase

* Item com texto **negrito**
* Item com texto *itálico*
* Item com texto ***negrito e itálico***
* Item com `código` inline

---

## Listas com Links

* Item com [link](https://example.com)
* Item com [link de referência][ref]
* Item com [link](https://example.com "título")

[ref]: https://example.com/reference

---

## Listas com Imagens

* Item com ![imagem](image.png)
* Item com ![imagem com alt](image.png "título")
* Item com ![imagem de referência][imgref]

[imgref]: image.png "Imagem de referência"

---

## Listas de Tarefas (GitHub Flavored Markdown)

- [ ] Tarefa não marcada
- [x] Tarefa marcada
- [X] Tarefa marcada (maiúscula)
- [ ] Outra tarefa não marcada
- [x] Outra tarefa marcada

---

## Listas de Tarefas Aninhadas

- [ ] Tarefa principal
  - [ ] Subtarefa 1
  - [x] Subtarefa 2 (concluída)
  - [ ] Subtarefa 3
- [x] Outra tarefa principal (concluída)
  - [x] Subtarefa concluída
  - [ ] Subtarefa incompleta

---

## Listas de Tarefas com Conteúdo

- [ ] Tarefa não marcada com texto **negrito**
- [x] Tarefa marcada com texto *itálico*
- [ ] Tarefa com [link](https://example.com)
- [x] Tarefa com `código` inline

---

## Listas com HTML

* Item com <strong>HTML negrito</strong>
* Item com <em>HTML itálico</em>
* Item com <code>HTML código</code>

---

## Listas com Caracteres Especiais

* Item com "aspas"
* Item com (parênteses)
* Item com [colchetes]
* Item com {chaves}
* Item com <tags>
* Item com $dólar$
* Item com &ampersand&

---

## Listas com Números

* Item com números 123
* Item com versão v1.2.3
* Item com porcentagem 50%
* Item com número grande 1,000,000

---

## Listas com Unicode

* Item com café
* Item com naïve
* Item com 你好
* Item com こんにちは
* Item com مرحبا

---

## Listas Apertadas (Sem Linhas em Branco)

* Primeiro item
* Segundo item
* Terceiro item

1. Primeiro item
2. Segundo item
3. Terceiro item

---

## Listas Soltas (Com Linhas em Branco)

* Primeiro item

* Segundo item

* Terceiro item

1. Primeiro item

2. Segundo item

3. Terceiro item

---

## Listas que Começam no Meio do Documento

Algum texto de parágrafo antes da lista.

* Item de lista 1
* Item de lista 2

Algum texto de parágrafo depois da lista.

---

## Listas que Terminam o Documento

* Item final 1
* Item final 2
* Item final 3

---

## Itens de Lista Vazios

* 
* Item com conteúdo
* 
* Outro item

---

## Listas com Conteúdo Longo

* Este é um item de lista muito longo que contém múltiplas sentenças e continua por um tempo para testar como o conversor lida com conteúdo mais longo dentro dos itens de lista. Deve manter a formatação e estrutura adequadas.

* Outro item longo que abrange múltiplas linhas na fonte mas deve ser renderizado como um único parágrafo dentro da estrutura do item de lista.

---

## Listas Aninhadas Complexas

1. Primeiro item ordenado
   * Item não ordenado aninhado
     - Item aninhado profundo
       + Item aninhado mais profundo
   * Outro item não ordenado aninhado
2. Segundo item ordenado
   * Item não ordenado aninhado
     1. Item ordenado aninhado
     2. Outro item ordenado aninhado
3. Terceiro item ordenado

---

## Listas com Regras Horizontais

* Item antes da regra

---

* Item depois da regra

---

## Listas com Cabeçalhos

* Item antes do cabeçalho

## Cabeçalho em Contexto de Lista

* Item depois do cabeçalho

---

## Listas de Definição (se suportado)

Termo 1
: Definição 1

Termo 2
: Definição 2
: Definição alternativa 2

Termo 3
: Definição 3 com **negrito** e *itálico*

