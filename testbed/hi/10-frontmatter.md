---
title: Frontmatter विविधताएं
categories:
  - frontmatter
  - yaml
tags:
  - metadata
  - frontmatter
  - yaml
author: परीक्षण लेखक
date: 2024-01-15
published: true
---

# Frontmatter विविधताएं

यह दस्तावेज़ विभिन्न frontmatter प्रारूपों और संरचनाओं का परीक्षण करता है।

---

## सरल Frontmatter

```yaml
---
title: Simple Document
author: Author Name
---
```

---

## सरणियों के साथ Frontmatter

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

## नेस्टेड संरचनाओं के साथ Frontmatter

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

## जटिल सरणियों के साथ Frontmatter

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

## खाली मानों के साथ Frontmatter

```yaml
---
title: Empty Values
author:
description:
tags: []
---
```

---

## उद्धृत मानों के साथ Frontmatter

```yaml
---
title: "Quoted Title"
author: 'Single Quoted Author'
description: "Description with \"quotes\" inside"
---
```

---

## विशेष वर्णों के साथ Frontmatter

```yaml
---
title: Title with (parentheses) and [brackets]
author: Author with &ampersand&
description: Description with $dollar$ and %percent%
---
```

---

## संख्याओं के साथ Frontmatter

```yaml
---
title: Title with Numbers
version: 1.2.3
year: 2024
percentage: 50%
---
```

---

## तारीखों के साथ Frontmatter

```yaml
---
title: Document with Dates
date: 2024-01-15
published_date: 2024-01-15T10:30:00Z
updated: 2024-12-31
---
```

---

## बूलियन के साथ Frontmatter

```yaml
---
title: Document with Booleans
published: true
draft: false
featured: true
---
```

---

## Null मानों के साथ Frontmatter

```yaml
---
title: Document with Null
author: null
description: null
tags: null
---
```

---

## मिश्रित प्रकारों के साथ Frontmatter

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

## लंबे मानों के साथ Frontmatter

```yaml
---
title: Document with Long Values
description: This is a very long description that spans multiple lines and contains extensive information about the document content and purpose. It tests how frontmatter handles longer string values that might wrap or extend beyond normal line lengths.
---
```

---

## Unicode के साथ Frontmatter

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

## URL के साथ Frontmatter

```yaml
---
title: Document with URLs
url: https://example.com/document
image: https://example.com/image.png
canonical: https://example.com/canonical
---
```

---

## ईमेल के साथ Frontmatter

```yaml
---
title: Document with Email
author: Author Name
email: author@example.com
contact: contact@example.com
---
```

---

## फ़ाइल पथों के साथ Frontmatter

```yaml
---
title: Document with Paths
source: ./source/file.md
output: ./output/file.html
template: templates/default.html
---
```

---

## वस्तुओं की सरणियों के साथ Frontmatter

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

## गहरे नेस्टिंग के साथ Frontmatter

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

## कई सरणियों के साथ Frontmatter

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

## खाली सरणियों के साथ Frontmatter

```yaml
---
title: Empty Arrays
categories: []
tags: []
keywords: []
---
```

---

## एकल आइटम सरणियों के साथ Frontmatter

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

## मिश्रित सरणी प्रारूपों के साथ Frontmatter

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

## टिप्पणियों के साथ Frontmatter (यदि समर्थित)

```yaml
---
# This is a comment
title: Document with Comments
# Another comment
author: Author Name
---
```

---

## एंकर और उपनामों के साथ Frontmatter (यदि समर्थित)

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

## बहु-पंक्ति स्ट्रिंग्स के साथ Frontmatter

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

## मुड़ी हुई स्ट्रिंग्स के साथ Frontmatter

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

## शाब्दिक ब्लॉक स्केलर के साथ Frontmatter

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

## कस्टम फ़ील्ड के साथ Frontmatter

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

## संस्करण संख्याओं के साथ Frontmatter

```yaml
---
title: Version Numbers
version: 1.0.0
api_version: 2.1.3
semver: 1.2.3-beta.1
---
```

---

## मुद्रा के साथ Frontmatter

```yaml
---
title: Currency Values
price: $19.99
currency: USD
amount: 100.50
---
```

---

## प्रतिशत के साथ Frontmatter

```yaml
---
title: Percentages
completion: 75%
discount: 20%
tax_rate: 8.5%
---
```

---

## समय मानों के साथ Frontmatter

```yaml
---
title: Time Values
duration: 2h30m
timeout: 5s
interval: 1m
---
```

---

## IP पतों के साथ Frontmatter

```yaml
---
title: IP Addresses
ipv4: 192.168.1.1
ipv6: 2001:0db8:85a3:0000:0000:8a2e:0370:7334
---
```

---

## संस्करण नियंत्रण के साथ Frontmatter

```yaml
---
title: Version Control
git_hash: abc123def456
branch: main
commit_date: 2024-01-15
---
```

---

## सोशल मीडिया के साथ Frontmatter

```yaml
---
title: Social Media
twitter: @username
github: github.com/username
linkedin: linkedin.com/in/username
---
```

---

## भौगोलिक डेटा के साथ Frontmatter

```yaml
---
title: Geographic Data
latitude: 40.7128
longitude: -74.0060
location: New York, NY
---
```

---

## रेटिंग के साथ Frontmatter

```yaml
---
title: Ratings
rating: 4.5
stars: 5
score: 8.5/10
---
```

---

## जटिल नेस्टेड सरणियों के साथ Frontmatter

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

## एस्केप किए गए वर्णों के साथ Frontmatter

```yaml
---
title: Escaped Characters
description: "Description with \"quotes\""
path: "C:\\Users\\Name\\File"
---
```

---

## JSON-जैसी संरचनाओं के साथ Frontmatter

```yaml
---
title: JSON-like Structures
config: '{"key": "value", "number": 123}'
settings: '{"enabled": true, "count": 5}'
---
```

---

## Base64 के साथ Frontmatter (यदि आवश्यक)

```yaml
---
title: Base64 Data
encoded: SGVsbG8gV29ybGQ=
---
```

---

## HTML के साथ Frontmatter

```yaml
---
title: HTML in Frontmatter
description: "<strong>Bold</strong> and <em>italic</em>"
---
```

---

## Markdown के साथ Frontmatter

```yaml
---
title: Markdown in Frontmatter
description: "**Bold** and *italic* text"
---
```

---

## बहुत लंबे फ़ील्ड नामों के साथ Frontmatter

```yaml
---
title: Very Long Field Names
very_long_field_name_that_extends_beyond_normal_width: value
another_extremely_long_field_name_with_many_words: another value
---
```

---

## विशेष YAML सुविधाओं के साथ Frontmatter

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

## तीन बिंदुओं के साथ समाप्त होने वाला Frontmatter

```yaml
---
title: Three Dots Ending
author: Author Name
...
```

---

## TOML-शैली के साथ Frontmatter (यदि समर्थित)

```toml
+++
title = "TOML Frontmatter"
author = "Author Name"
+++
```

---

## JSON के साथ Frontmatter (यदि समर्थित)

```json
{
  "title": "JSON Frontmatter",
  "author": "Author Name"
}
```

---

## कई दस्तावेज़ विभाजकों के साथ Frontmatter

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

## इनलाइन सरणियों के साथ Frontmatter

```yaml
---
title: Inline Arrays
tags: [tag1, tag2, tag3]
categories: [cat1, cat2]
---
```

---

## मिश्रित उद्धरणों के साथ Frontmatter

```yaml
---
title: 'Single quoted title'
description: "Double quoted description"
author: Unquoted author
---
```

---

## क्रम संरक्षित करने वाला Frontmatter

```yaml
---
z_field: Last
a_field: First
m_field: Middle
---
```

