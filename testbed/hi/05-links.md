---
title: लिंक - इनलाइन, संदर्भ और ऑटोलिंक
categories:
  - markdown
  - links
tags:
  - hyperlinks
  - references
  - urls
---

# लिंक

## इनलाइन लिंक

यह एक [इनलाइन लिंक](https://example.com) है।

यह एक [शीर्षक के साथ लिंक](https://example.com "लिंक शीर्षक") है।

यह एक [एकल उद्धरण के साथ लिंक](https://example.com 'लिंक शीर्षक') है।

यह एक [कोष्ठक के साथ लिंक](https://example.com (लिंक शीर्षक)) है।

---

## विभिन्न प्रोटोकॉल के साथ लिंक

[HTTP लिंक](http://example.com)

[HTTPS लिंक](https://example.com)

[FTP लिंक](ftp://ftp.example.com)

[ईमेल लिंक](mailto:user@example.com)

[फ़ाइल लिंक](file:///path/to/file)

---

## पथों के साथ लिंक

[रूट पथ](https://example.com/)

[उपनिर्देशिका](https://example.com/path/to/resource)

[फ़ाइल एक्सटेंशन](https://example.com/file.html)

[क्वेरी स्ट्रिंग](https://example.com/page?param=value&other=123)

[फ़्रैगमेंट](https://example.com/page#section)

[जटिल URL](https://example.com/path/to/resource?query=value&other=123#fragment)

---

## विशेष वर्णों के साथ लिंक

[रिक्त स्थान के साथ लिंक](https://example.com/path with spaces)

[unicode के साथ लिंक](https://example.com/你好)

[एन्कोडेड के साथ लिंक](https://example.com/path%20with%20spaces)

---

## संदर्भ-शैली लिंक

यह एक [संदर्भ लिंक][ref1] है।

यह एक [अंतर्निहित संदर्भ लिंक][] है।

यह एक [शॉर्टकट संदर्भ लिंक] है।

[ref1]: https://example.com/reference1
[implicit reference link]: https://example.com/implicit
[shortcut reference link]: https://example.com/shortcut

---

## शीर्षकों के साथ संदर्भ लिंक

यह एक [शीर्षक के साथ संदर्भ लिंक][ref2] है।

[ref2]: https://example.com/reference2 "संदर्भ शीर्षक"

यह एक [एकल उद्धरण के साथ संदर्भ लिंक][ref3] है।

[ref3]: https://example.com/reference3 'संदर्भ शीर्षक'

यह एक [कोष्ठक के साथ संदर्भ लिंक][ref4] है।

[ref4]: https://example.com/reference4 (संदर्भ शीर्षक)

---

## संदर्भ लिंक - केस संवेदनशीलता

[केस संवेदनशील लिंक][case1]

[केस संवेदनशील लिंक][Case1]

[case1]: https://example.com/lowercase
[Case1]: https://example.com/uppercase

---

## संदर्भ लिंक - संख्याएं

[लिंक 1][1]

[लिंक 2][2]

[1]: https://example.com/1
[2]: https://example.com/2

---

## संदर्भ लिंक - विशेष वर्ण

[डैश के साथ लिंक][link-dash]

[अंडरस्कोर के साथ लिंक][link_underscore]

[डॉट के साथ लिंक][link.dot]

[link-dash]: https://example.com/dash
[link_underscore]: https://example.com/underscore
[link.dot]: https://example.com/dot

---

## ऑटोलिंक

<https://example.com>

<http://example.com>

<mailto:user@example.com>

<user@example.com>

---

## जोर के साथ लिंक

**[मोटा लिंक](https://example.com)**

*[तिरछा लिंक](https://example.com)*

***[मोटा और तिरछा लिंक](https://example.com)***

[**मोटा पाठ**](https://example.com) लिंक में

[*तिरछा पाठ*](https://example.com) लिंक में

---

## सूचियों में लिंक

* [लिंक](https://example.com) के साथ आइटम
* [संदर्भ लिंक][listref] के साथ आइटम
* <https://example.com> के साथ आइटम

[listref]: https://example.com/list

---

## ब्लॉककोट्स में लिंक

> यह ब्लॉककोट में एक [लिंक](https://example.com) है।

> यह ब्लॉककोट में एक [संदर्भ लिंक][blockquoteref] है।

[blockquoteref]: https://example.com/blockquote

---

## शीर्षकों में लिंक

# [लिंक](https://example.com) के साथ शीर्षक

## [संदर्भ लिंक][headerref] के साथ शीर्षक

[headerref]: https://example.com/header

---

## तालिकाओं में लिंक

| Column 1 | Column 2 |
|----------|----------|
| [लिंक](https://example.com) | [संदर्भ][tableref] |
| <https://example.com> | [शीर्षक के साथ लिंक](https://example.com "शीर्षक") |

[tableref]: https://example.com/table

---

## छवियों के साथ लिंक

[![छवि](image.png)](https://example.com)

[![alt के साथ छवि](image.png "शीर्षक")](https://example.com "लिंक शीर्षक")

---

## सापेक्ष लिंक

[सापेक्ष लिंक](./relative.html)

[पैरेंट निर्देशिका](../parent.html)

[सहोदर फ़ाइल](../sibling.html)

[रूट सापेक्ष](/root.html)

---

## एंकर लिंक

[अनुभाग के लिए लिंक](#section)

[उपअनुभाग के लिए लिंक](#subsection-name)

[रिक्त स्थान के साथ लिंक](#section name)

[विशेष वर्णों के साथ लिंक](#section-name_123)

---

## कोड के साथ लिंक

[`कोड` के साथ लिंक](https://example.com)

`[कोड लिंक](https://example.com)` (लिंक नहीं होना चाहिए)

---

## HTML के साथ लिंक

<a href="https://example.com">HTML लिंक</a>

[Markdown लिंक](https://example.com) और <a href="https://example.com">HTML लिंक</a>

---

## टूटे हुए लिंक

[टूटा हुआ लिंक](https://nonexistent.example.com)

[टूटा हुआ संदर्भ][broken]

---

## एस्केप किए गए वर्णों के साथ लिंक

\[एस्केप किया गया लिंक\]\(https://example.com\)

[\*तारांकन\* के साथ लिंक](https://example.com)

[\_अंडरस्कोर\_ के साथ लिंक](https://example.com)

---

## Unicode के साथ लिंक

[café के साथ लिंक](https://example.com/café)

[你好 के साथ लिंक](https://example.com/你好)

[こんにちは के साथ लिंक](https://example.com/こんにちは)

---

## एक अनुच्छेद में कई लिंक

इस अनुच्छेद में [पहला लिंक](https://example.com/1), [दूसरा लिंक](https://example.com/2), और [तीसरा लिंक](https://example.com/3) हैं।

यह अनुच्छेद [इनलाइन लिंक](https://example.com/inline), [संदर्भ लिंक][multiref], और <https://example.com/autolink> को मिलाता है।

[multiref]: https://example.com/reference

---

## लंबे URL के साथ लिंक

[बहुत लंबे URL के साथ लिंक](https://example.com/very/long/path/to/resource/that/spans/multiple/segments/and/includes/query/parameters?param1=value1&param2=value2&param3=value3#and-even-a-fragment)

---

## खाली पाठ के साथ लिंक

[](https://example.com)

[ ](https://example.com)

---

## केवल रिक्त स्थान के साथ लिंक

[   ](https://example.com)

---

## कोड ब्लॉक में लिंक

```
[लिंक](https://example.com) कोड ब्लॉक में काम नहीं करना चाहिए
```

---

## विशेष URL योजनाओं के साथ लिंक

[javascript:alert('XSS')](javascript:alert('XSS'))

[data:text/plain,Hello](data:text/plain,Hello)

[about:blank](about:blank)

---

## प्रतिशत एन्कोडिंग के साथ लिंक

[एन्कोडिंग के साथ लिंक](https://example.com/path%20with%20spaces)

[unicode एन्कोडिंग के साथ लिंक](https://example.com/%E4%BD%A0%E5%A5%BD)

---

## जटिल लिंक संयोजन

**[मोटा लिंक](https://example.com)** *तिरछा* पाठ और [एक और लिंक](https://example.com/2) के साथ।

[**मोटा** और *तिरछा* के साथ लिंक](https://example.com) पाठ में।

[लिंक](https://example.com) `कोड` और **मोटा** और *तिरछा* के साथ।

