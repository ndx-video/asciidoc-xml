---
title: الصور - مضمنة ونمط المرجع
categories:
  - markdown
  - images
tags:
  - graphics
  - media
  - pictures
---

# الصور

## صور مضمنة

![نص بديل](image.png)

![نص بديل مع مسافات](image with spaces.png)

![نص بديل](image.png "عنوان الصورة")

![نص بديل](image.png 'عنوان الصورة')

![نص بديل](image.png (عنوان الصورة))

---

## صور بصيغ مختلفة

![صورة PNG](image.png)

![صورة JPEG](image.jpg)

![صورة GIF](image.gif)

![صورة SVG](image.svg)

![صورة WebP](image.webp)

---

## صور مع مسارات

![صورة الجذر](/image.png)

![صورة مجلد فرعي](./images/image.png)

![مجلد أب](../images/image.png)

![مسار عميق](./images/subfolder/image.png)

---

## صور مع عناوين URL

![صورة بعيدة](https://example.com/image.png)

![صورة بعيدة مع مسار](https://example.com/images/photo.jpg)

![صورة بعيدة مع استعلام](https://example.com/image.png?size=large)

---

## صور نمط المرجع

![صورة مرجعية][img1]

![صورة مرجعية ضمنية][]

![صورة مرجعية مختصرة]

[img1]: image1.png
[implicit reference image]: image2.png
[shortcut reference image]: image3.png

---

## صور مرجعية مع عناوين

![صورة مرجعية مع عنوان][img2]

[img2]: image.png "عنوان الصورة"

![صورة مرجعية مع علامات اقتباس مفردة][img3]

[img3]: image.png 'عنوان الصورة'

![صورة مرجعية مع أقواس][img4]

[img4]: image.png (عنوان الصورة)

---

## صور مع أحرف خاصة في النص البديل

![صورة مع "اقتباسات"](image.png)

![صورة مع (أقواس)](image.png)

![صورة مع [أقواس مربعة]](image.png)

![صورة مع {أقواس معقوفة}](image.png)

![صورة مع <علامات>](image.png)

![صورة مع $دولار$](image.png)

---

## صور مع Unicode في النص البديل

![صورة مع café](image.png)

![صورة مع 你好](image.png)

![صورة مع こんにちは](image.png)

![صورة مع مرحبا](image.png)

---

## صور مع أرقام في النص البديل

![صورة 123](image.png)

![صورة v1.2.3](image.png)

![صورة 50%](image.png)

---

## صور مع التأكيد في النص البديل

![**نص بديل عريض**](image.png)

![*نص بديل مائل*](image.png)

![***نص بديل عريض ومائل***](image.png)

---

## صور مع روابط

[![صورة مرتبطة](image.png)](https://example.com)

[![صورة مرتبطة مع alt](image.png "عنوان الصورة")](https://example.com "عنوان الرابط")

---

## صور في القوائم

* عنصر مع ![صورة](image.png)
* عنصر مع ![صورة مرجعية][listimg]
* عنصر مع ![صورة](image.png "عنوان")

[listimg]: image.png

---

## صور في الاقتباسات

> هذا اقتباس مع ![صورة](image.png)

> هذا اقتباس مع ![صورة مرجعية][blockquoteimg]

[blockquoteimg]: image.png

---

## صور في العناوين

# عنوان مع ![صورة](image.png)

## عنوان مع ![صورة مرجعية][headerimg]

[headerimg]: image.png

---

## صور في الجداول

| Column 1 | Column 2 |
|----------|----------|
| ![صورة](image.png) | ![مرجع][tableimg] |
| نص | ![صورة مع عنوان](image.png "عنوان") |

[tableimg]: image.png

---

## صور مع نص بديل فارغ

![](image.png)

![ ](image.png)

![   ](image.png)

---

## صور فقط مع مسافات في النص البديل

![   ](image.png)

---

## صور مع نص بديل طويل

![هذا نص بديل طويل جدًا يصف الصورة بالتفصيل ويوفر معلومات شاملة حول ما تحتويه الصورة والغرض منها في المستند](image.png)

---

## صور متعددة في فقرة واحدة

هذه الفقرة لها ![الصورة الأولى](image1.png) و ![الصورة الثانية](image2.png) و ![الصورة الثالثة](image3.png).

تجمع هذه الفقرة ![صورة مضمنة](image.png) و ![صورة مرجعية][multimg] ونص.

[multimg]: image.png

---

## صور مع HTML

<img src="image.png" alt="صورة HTML">

![صورة Markdown](image.png) و <img src="image.png" alt="صورة HTML">

---

## صور مع أحرف مُهربة

\![صورة مُهربة](image.png)

![صورة مع \*نجمة\*](image.png)

![صورة مع \_شرطة سفلية\_](image.png)

---

## صور في كتل الكود

```
![صورة](image.png) لا ينبغي أن تعمل في كتل الكود
```

---

## صور مع مخططات URL خاصة

![صورة Data URI](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==)

---

## صور مع ترميز النسبة المئوية

![صورة مع ترميز](https://example.com/image%20with%20spaces.png)

![صورة مع ترميز unicode](https://example.com/%E4%BD%A0%E5%A5%BD.png)

---

## مجموعات صور معقدة

**نص عريض** مع ![صورة](image.png) ونص *مائل*.

![صورة](image.png) مع [رابط](https://example.com) و `كود`.

![**Alt عريض**](image.png) مع نص *مائل* و [رابط](https://example.com).

---

## صور بأحجام مختلفة (سمات HTML)

<img src="image.png" alt="صورة صغيرة" width="100" height="100">

<img src="image.png" alt="صورة كبيرة" width="800" height="600">

---

## صور مع فئات CSS (HTML)

<img src="image.png" alt="صورة منمقة" class="custom-class">

<img src="image.png" alt="صورة م centrada" style="display: block; margin: 0 auto;">

---

## صور مع عناوين وتسميات توضيحية

![صورة](image.png "هذا عنوان")

![صورة](image.png "عنوان مع 'اقتباسات'")

![صورة](image.png "عنوان مع (أقواس)")

---

## صور مع مسارات نسبية

![نفس المجلد](./image.png)

![مجلد فرعي](./images/image.png)

![مجلد أب](../image.png)

![نسبي للجذر](/image.png)

---

## صور مع عناوين URL مطلقة

![URL مطلقة](https://example.com/image.png)

![URL مطلقة مع مسار](https://example.com/images/photo.jpg)

![URL مطلقة مع استعلام](https://example.com/image.png?v=1&size=large)

![URL مطلقة مع جزء](https://example.com/image.png#section)

---

## صور في هياكل متداخلة

### صور في قوائم متداخلة

* المستوى الأول
  * المستوى الثاني مع ![صورة](image.png)
    * المستوى الثالث مع ![صورة](image.png)

### صور في اقتباسات متداخلة

> اقتباس خارجي
> > اقتباس داخلي مع ![صورة](image.png)

---

## صور مع مراجع مكسورة

![صورة مكسورة](nonexistent.png)

![مرجع مكسور][broken]

[broken]: nonexistent.png

---

## صور مع أسماء ملفات خاصة

![صورة مع مسافات](image with spaces.png)

![صورة.مع.نقاط](image.with.dots.png)

![صورة-مع-شرطات](image-with-dashes.png)

![صورة_مع_شرطات_سفلية](Image_with_underscores.png)

![صورة123](Image123.png)

![صورة@خاصة](Image@special.png)

