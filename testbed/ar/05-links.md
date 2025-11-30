---
title: الروابط - مضمنة ومرجعية وتلقائية
categories:
  - markdown
  - links
tags:
  - hyperlinks
  - references
  - urls
---

# الروابط

## روابط مضمنة

هذا [رابط مضمن](https://example.com).

هذا [رابط بعنوان](https://example.com "عنوان الرابط").

هذا [رابط بعلامات اقتباس مفردة](https://example.com 'عنوان الرابط').

هذا [رابط بأقواس](https://example.com (عنوان الرابط)).

---

## روابط مع بروتوكولات مختلفة

[رابط HTTP](http://example.com)

[رابط HTTPS](https://example.com)

[رابط FTP](ftp://ftp.example.com)

[رابط بريد إلكتروني](mailto:user@example.com)

[رابط ملف](file:///path/to/file)

---

## روابط مع مسارات

[مسار الجذر](https://example.com/)

[مجلد فرعي](https://example.com/path/to/resource)

[امتداد الملف](https://example.com/file.html)

[سلسلة استعلام](https://example.com/page?param=value&other=123)

[جزء](https://example.com/page#section)

[عنوان URL معقد](https://example.com/path/to/resource?query=value&other=123#fragment)

---

## روابط مع أحرف خاصة

[رابط مع مسافات](https://example.com/path with spaces)

[رابط مع unicode](https://example.com/你好)

[رابط مع ترميز](https://example.com/path%20with%20spaces)

---

## روابط نمط المرجع

هذا [رابط مرجعي][ref1].

هذا [رابط مرجعي ضمني][].

هذا [رابط مرجعي مختصر].

[ref1]: https://example.com/reference1
[implicit reference link]: https://example.com/implicit
[shortcut reference link]: https://example.com/shortcut

---

## روابط مرجعية مع عناوين

هذا [رابط مرجعي بعنوان][ref2].

[ref2]: https://example.com/reference2 "عنوان المرجع"

هذا [رابط مرجعي بعلامات اقتباس مفردة][ref3].

[ref3]: https://example.com/reference3 'عنوان المرجع'

هذا [رابط مرجعي بأقواس][ref4].

[ref4]: https://example.com/reference4 (عنوان المرجع)

---

## روابط مرجعية - حساسية الأحرف

[رابط حساس للحالة][case1]

[رابط حساس للحالة][Case1]

[case1]: https://example.com/lowercase
[Case1]: https://example.com/uppercase

---

## روابط مرجعية - أرقام

[رابط 1][1]

[رابط 2][2]

[1]: https://example.com/1
[2]: https://example.com/2

---

## روابط مرجعية - أحرف خاصة

[رابط مع شرطة][link-dash]

[رابط مع شرطة سفلية][link_underscore]

[رابط مع نقطة][link.dot]

[link-dash]: https://example.com/dash
[link_underscore]: https://example.com/underscore
[link.dot]: https://example.com/dot

---

## روابط تلقائية

<https://example.com>

<http://example.com>

<mailto:user@example.com>

<user@example.com>

---

## روابط مع التأكيد

**[رابط عريض](https://example.com)**

*[رابط مائل](https://example.com)*

***[رابط عريض ومائل](https://example.com)***

[**نص عريض**](https://example.com) في الرابط

[*نص مائل*](https://example.com) في الرابط

---

## روابط في القوائم

* عنصر مع [رابط](https://example.com)
* عنصر مع [رابط مرجعي][listref]
* عنصر مع <https://example.com>

[listref]: https://example.com/list

---

## روابط في الاقتباسات

> هذا [رابط](https://example.com) في اقتباس.

> هذا [رابط مرجعي][blockquoteref] في اقتباس.

[blockquoteref]: https://example.com/blockquote

---

## روابط في العناوين

# عنوان مع [رابط](https://example.com)

## عنوان مع [رابط مرجعي][headerref]

[headerref]: https://example.com/header

---

## روابط في الجداول

| Column 1 | Column 2 |
|----------|----------|
| [رابط](https://example.com) | [مرجع][tableref] |
| <https://example.com> | [رابط بعنوان](https://example.com "عنوان") |

[tableref]: https://example.com/table

---

## روابط مع صور

[![صورة](image.png)](https://example.com)

[![صورة مع نص بديل](image.png "عنوان")](https://example.com "عنوان الرابط")

---

## روابط نسبية

[رابط نسبي](./relative.html)

[مجلد أب](../parent.html)

[ملف شقيق](../sibling.html)

[نسبي للجذر](/root.html)

---

## روابط المرساة

[رابط إلى قسم](#section)

[رابط إلى قسم فرعي](#subsection-name)

[رابط مع مسافات](#section name)

[رابط مع أحرف خاصة](#section-name_123)

---

## روابط مع كود

[رابط مع `كود`](https://example.com)

`[رابط كود](https://example.com)` (لا ينبغي أن يكون رابطًا)

---

## روابط مع HTML

<a href="https://example.com">رابط HTML</a>

[رابط Markdown](https://example.com) و <a href="https://example.com">رابط HTML</a>

---

## روابط مكسورة

[رابط مكسور](https://nonexistent.example.com)

[مرجع مكسور][broken]

---

## روابط مع أحرف مُهربة

\[رابط مُهرب\]\(https://example.com\)

[رابط مع \*نجمة\*](https://example.com)

[رابط مع \_شرطة سفلية\_](https://example.com)

---

## روابط مع Unicode

[رابط مع café](https://example.com/café)

[رابط مع 你好](https://example.com/你好)

[رابط مع こんにちは](https://example.com/こんにちは)

---

## روابط متعددة في فقرة واحدة

هذه الفقرة لها [الرابط الأول](https://example.com/1) و [الرابط الثاني](https://example.com/2) و [الرابط الثالث](https://example.com/3).

تجمع هذه الفقرة [رابط مضمن](https://example.com/inline) و [رابط مرجعي][multiref] و <https://example.com/autolink>.

[multiref]: https://example.com/reference

---

## روابط مع عناوين URL طويلة

[رابط مع عنوان URL طويل جدًا](https://example.com/very/long/path/to/resource/that/spans/multiple/segments/and/includes/query/parameters?param1=value1&param2=value2&param3=value3#and-even-a-fragment)

---

## روابط مع نص فارغ

[](https://example.com)

[ ](https://example.com)

---

## روابط فقط مع مسافات

[   ](https://example.com)

---

## روابط في كتل الكود

```
[رابط](https://example.com) لا ينبغي أن يعمل في كتل الكود
```

---

## روابط مع مخططات URL خاصة

[javascript:alert('XSS')](javascript:alert('XSS'))

[data:text/plain,Hello](data:text/plain,Hello)

[about:blank](about:blank)

---

## روابط مع ترميز النسبة المئوية

[رابط مع ترميز](https://example.com/path%20with%20spaces)

[رابط مع ترميز unicode](https://example.com/%E4%BD%A0%E5%A5%BD)

---

## مجموعات روابط معقدة

**[رابط عريض](https://example.com)** مع نص *مائل* و [رابط آخر](https://example.com/2).

[رابط مع **عريض** و *مائل*](https://example.com) في النص.

[رابط](https://example.com) مع `كود` و **عريض** و *مائل*.

