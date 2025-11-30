---
title: छवियां - इनलाइन और संदर्भ शैली
categories:
  - markdown
  - images
tags:
  - graphics
  - media
  - pictures
---

# छवियां

## इनलाइन छवियां

![वैकल्पिक पाठ](image.png)

![रिक्त स्थान के साथ वैकल्पिक पाठ](image with spaces.png)

![वैकल्पिक पाठ](image.png "छवि शीर्षक")

![वैकल्पिक पाठ](image.png 'छवि शीर्षक')

![वैकल्पिक पाठ](image.png (छवि शीर्षक))

---

## विभिन्न प्रारूपों के साथ छवियां

![PNG छवि](image.png)

![JPEG छवि](image.jpg)

![GIF छवि](image.gif)

![SVG छवि](image.svg)

![WebP छवि](image.webp)

---

## पथों के साथ छवियां

![रूट छवि](/image.png)

![उपनिर्देशिका छवि](./images/image.png)

![पैरेंट निर्देशिका](../images/image.png)

![गहरा पथ](./images/subfolder/image.png)

---

## URL के साथ छवियां

![दूरस्थ छवि](https://example.com/image.png)

![पथ के साथ दूरस्थ छवि](https://example.com/images/photo.jpg)

![क्वेरी के साथ दूरस्थ छवि](https://example.com/image.png?size=large)

---

## संदर्भ-शैली छवियां

![संदर्भ छवि][img1]

![अंतर्निहित संदर्भ छवि][]

![शॉर्टकट संदर्भ छवि]

[img1]: image1.png
[implicit reference image]: image2.png
[shortcut reference image]: image3.png

---

## शीर्षकों के साथ संदर्भ छवियां

![शीर्षक के साथ संदर्भ छवि][img2]

[img2]: image.png "छवि शीर्षक"

![एकल उद्धरण के साथ संदर्भ छवि][img3]

[img3]: image.png 'छवि शीर्षक'

![कोष्ठक के साथ संदर्भ छवि][img4]

[img4]: image.png (छवि शीर्षक)

---

## वैकल्पिक पाठ में विशेष वर्णों के साथ छवियां

!["उद्धरण" के साथ छवि](image.png)

![(कोष्ठक) के साथ छवि](image.png)

!["[वर्ग कोष्ठक]" के साथ छवि](image.png)

![{ब्रेस} के साथ छवि](image.png)

![<टैग> के साथ छवि](image.png)

![$डॉलर$ के साथ छवि](image.png)

---

## वैकल्पिक पाठ में Unicode के साथ छवियां

![café के साथ छवि](image.png)

![你好 के साथ छवि](image.png)

![こんにちは के साथ छवि](image.png)

![مرحبا के साथ छवि](image.png)

---

## वैकल्पिक पाठ में संख्याओं के साथ छवियां

![छवि 123](image.png)

![छवि v1.2.3](image.png)

![छवि 50%](image.png)

---

## वैकल्पिक पाठ में जोर के साथ छवियां

![**मोटा वैकल्पिक पाठ**](image.png)

![*तिरछा वैकल्पिक पाठ*](image.png)

![***मोटा और तिरछा वैकल्पिक पाठ***](image.png)

---

## लिंक के साथ छवियां

[![लिंक की गई छवि](image.png)](https://example.com)

[![alt के साथ लिंक की गई छवि](image.png "छवि शीर्षक")](https://example.com "लिंक शीर्षक")

---

## सूचियों में छवियां

* ![छवि](image.png) के साथ आइटम
* ![संदर्भ छवि][listimg] के साथ आइटम
* ![छवि](image.png "शीर्षक") के साथ आइटम

[listimg]: image.png

---

## ब्लॉककोट्स में छवियां

> यह ![छवि](image.png) के साथ एक ब्लॉककोट है

> यह ![संदर्भ छवि][blockquoteimg] के साथ एक ब्लॉककोट है

[blockquoteimg]: image.png

---

## शीर्षकों में छवियां

# ![छवि](image.png) के साथ शीर्षक

## ![संदर्भ छवि][headerimg] के साथ शीर्षक

[headerimg]: image.png

---

## तालिकाओं में छवियां

| Column 1 | Column 2 |
|----------|----------|
| ![छवि](image.png) | ![संदर्भ][tableimg] |
| पाठ | ![शीर्षक के साथ छवि](image.png "शीर्षक") |

[tableimg]: image.png

---

## खाली वैकल्पिक पाठ के साथ छवियां

![](image.png)

![ ](image.png)

![   ](image.png)

---

## केवल रिक्त स्थान के साथ वैकल्पिक पाठ वाली छवियां

![   ](image.png)

---

## लंबे वैकल्पिक पाठ के साथ छवियां

![यह एक बहुत लंबा वैकल्पिक पाठ है जो छवि का विस्तार से वर्णन करता है और छवि में क्या शामिल है और दस्तावेज़ में इसके उद्देश्य के बारे में व्यापक जानकारी प्रदान करता है](image.png)

---

## एक अनुच्छेद में कई छवियां

इस अनुच्छेद में ![पहली छवि](image1.png), ![दूसरी छवि](image2.png), और ![तीसरी छवि](image3.png) हैं।

यह अनुच्छेद ![इनलाइन छवि](image.png), ![संदर्भ छवि][multimg], और पाठ को मिलाता है।

[multimg]: image.png

---

## HTML के साथ छवियां

<img src="image.png" alt="HTML छवि">

![Markdown छवि](image.png) और <img src="image.png" alt="HTML छवि">

---

## एस्केप किए गए वर्णों के साथ छवियां

\![एस्केप की गई छवि](image.png)

![\*तारांकन\* के साथ छवि](image.png)

![\_अंडरस्कोर\_ के साथ छवि](image.png)

---

## कोड ब्लॉक में छवियां

```
![छवि](image.png) कोड ब्लॉक में काम नहीं करना चाहिए
```

---

## विशेष URL योजनाओं के साथ छवियां

![Data URI छवि](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==)

---

## प्रतिशत एन्कोडिंग के साथ छवियां

![एन्कोडिंग के साथ छवि](https://example.com/image%20with%20spaces.png)

![unicode एन्कोडिंग के साथ छवि](https://example.com/%E4%BD%A0%E5%A5%BD.png)

---

## जटिल छवि संयोजन

**मोटा पाठ** ![छवि](image.png) और *तिरछा* पाठ के साथ।

![छवि](image.png) [लिंक](https://example.com) और `कोड` के साथ।

![**मोटा alt**](image.png) *तिरछा* पाठ और [लिंक](https://example.com) के साथ।

---

## विभिन्न आकारों के साथ छवियां (HTML विशेषताएं)

<img src="image.png" alt="छोटी छवि" width="100" height="100">

<img src="image.png" alt="बड़ी छवि" width="800" height="600">

---

## CSS कक्षाओं के साथ छवियां (HTML)

<img src="image.png" alt="स्टाइल की गई छवि" class="custom-class">

<img src="image.png" alt="केंद्रित छवि" style="display: block; margin: 0 auto;">

---

## शीर्षकों और कैप्शन के साथ छवियां

![छवि](image.png "यह एक शीर्षक है")

![छवि](image.png "'उद्धरण' के साथ शीर्षक")

![छवि](image.png "(कोष्ठक) के साथ शीर्षक")

---

## सापेक्ष पथों के साथ छवियां

![समान निर्देशिका](./image.png)

![उपनिर्देशिका](./images/image.png)

![पैरेंट निर्देशिका](../image.png)

![रूट सापेक्ष](/image.png)

---

## निरपेक्ष URL के साथ छवियां

![निरपेक्ष URL](https://example.com/image.png)

![पथ के साथ निरपेक्ष URL](https://example.com/images/photo.jpg)

![क्वेरी के साथ निरपेक्ष URL](https://example.com/image.png?v=1&size=large)

![फ़्रैगमेंट के साथ निरपेक्ष URL](https://example.com/image.png#section)

---

## नेस्टेड संरचनाओं में छवियां

### नेस्टेड सूचियों में छवियां

* प्रथम स्तर
  * ![छवि](image.png) के साथ दूसरा स्तर
    * ![छवि](image.png) के साथ तीसरा स्तर

### नेस्टेड ब्लॉककोट्स में छवियां

> बाहरी ब्लॉककोट
> > ![छवि](image.png) के साथ आंतरिक ब्लॉककोट

---

## टूटे हुए संदर्भों के साथ छवियां

![टूटी हुई छवि](nonexistent.png)

![टूटा हुआ संदर्भ][broken]

[broken]: nonexistent.png

---

## विशेष फ़ाइल नामों के साथ छवियां

![रिक्त स्थान के साथ छवि](image with spaces.png)

![बिंदुओं के साथ छवि](image.with.dots.png)

![डैश के साथ छवि](image-with-dashes.png)

![अंडरस्कोर के साथ छवि](Image_with_underscores.png)

![छवि123](Image123.png)

![छवि@विशेष](Image@special.png)

