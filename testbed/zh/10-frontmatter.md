---
title: Frontmatter 变体
categories:
  - frontmatter
  - yaml
tags:
  - metadata
  - frontmatter
  - yaml
author: 测试作者
date: 2024-01-15
published: true
---

# Frontmatter 变体

本文档测试各种 frontmatter 格式和结构。

---

## 简单 Frontmatter

```yaml
---
title: Simple Document
author: Author Name
---
```

---

## 带数组的 Frontmatter

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

## 带嵌套结构的 Frontmatter

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

## 带复杂数组的 Frontmatter

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

## 带空值的 Frontmatter

```yaml
---
title: Empty Values
author:
description:
tags: []
---
```

---

## 带引号值的 Frontmatter

```yaml
---
title: "Quoted Title"
author: 'Single Quoted Author'
description: "Description with \"quotes\" inside"
---
```

---

## 带特殊字符的 Frontmatter

```yaml
---
title: Title with (parentheses) and [brackets]
author: Author with &ampersand&
description: Description with $dollar$ and %percent%
---
```

---

## 带数字的 Frontmatter

```yaml
---
title: Title with Numbers
version: 1.2.3
year: 2024
percentage: 50%
---
```

---

## 带日期的 Frontmatter

```yaml
---
title: Document with Dates
date: 2024-01-15
published_date: 2024-01-15T10:30:00Z
updated: 2024-12-31
---
```

---

## 带布尔值的 Frontmatter

```yaml
---
title: Document with Booleans
published: true
draft: false
featured: true
---
```

---

## 带 Null 值的 Frontmatter

```yaml
---
title: Document with Null
author: null
description: null
tags: null
---
```

---

## 带混合类型的 Frontmatter

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

## 带长值的 Frontmatter

```yaml
---
title: Document with Long Values
description: This is a very long description that spans multiple lines and contains extensive information about the document content and purpose. It tests how frontmatter handles longer string values that might wrap or extend beyond normal line lengths.
---
```

---

## 带 Unicode 的 Frontmatter

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

## 带 URL 的 Frontmatter

```yaml
---
title: Document with URLs
url: https://example.com/document
image: https://example.com/image.png
canonical: https://example.com/canonical
---
```

---

## 带电子邮件的 Frontmatter

```yaml
---
title: Document with Email
author: Author Name
email: author@example.com
contact: contact@example.com
---
```

---

## 带文件路径的 Frontmatter

```yaml
---
title: Document with Paths
source: ./source/file.md
output: ./output/file.html
template: templates/default.html
---
```

---

## 带对象数组的 Frontmatter

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

## 带深度嵌套的 Frontmatter

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

## 带多个数组的 Frontmatter

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

## 带空数组的 Frontmatter

```yaml
---
title: Empty Arrays
categories: []
tags: []
keywords: []
---
```

---

## 带单项数组的 Frontmatter

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

## 带混合数组格式的 Frontmatter

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

## 带注释的 Frontmatter（如果支持）

```yaml
---
# This is a comment
title: Document with Comments
# Another comment
author: Author Name
---
```

---

## 带锚点和别名的 Frontmatter（如果支持）

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

## 带多行字符串的 Frontmatter

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

## 带折叠字符串的 Frontmatter

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

## 带字面块标量的 Frontmatter

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

## 带自定义字段的 Frontmatter

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

## 带版本号的 Frontmatter

```yaml
---
title: Version Numbers
version: 1.0.0
api_version: 2.1.3
semver: 1.2.3-beta.1
---
```

---

## 带货币的 Frontmatter

```yaml
---
title: Currency Values
price: $19.99
currency: USD
amount: 100.50
---
```

---

## 带百分比的 Frontmatter

```yaml
---
title: Percentages
completion: 75%
discount: 20%
tax_rate: 8.5%
---
```

---

## 带时间值的 Frontmatter

```yaml
---
title: Time Values
duration: 2h30m
timeout: 5s
interval: 1m
---
```

---

## 带 IP 地址的 Frontmatter

```yaml
---
title: IP Addresses
ipv4: 192.168.1.1
ipv6: 2001:0db8:85a3:0000:0000:8a2e:0370:7334
---
```

---

## 带版本控制的 Frontmatter

```yaml
---
title: Version Control
git_hash: abc123def456
branch: main
commit_date: 2024-01-15
---
```

---

## 带社交媒体的 Frontmatter

```yaml
---
title: Social Media
twitter: @username
github: github.com/username
linkedin: linkedin.com/in/username
---
```

---

## 带地理数据的 Frontmatter

```yaml
---
title: Geographic Data
latitude: 40.7128
longitude: -74.0060
location: New York, NY
---
```

---

## 带评分的 Frontmatter

```yaml
---
title: Ratings
rating: 4.5
stars: 5
score: 8.5/10
---
```

---

## 带复杂嵌套数组的 Frontmatter

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

## 带转义字符的 Frontmatter

```yaml
---
title: Escaped Characters
description: "Description with \"quotes\""
path: "C:\\Users\\Name\\File"
---
```

---

## 带类似 JSON 结构的 Frontmatter

```yaml
---
title: JSON-like Structures
config: '{"key": "value", "number": 123}'
settings: '{"enabled": true, "count": 5}'
---
```

---

## 带 Base64 的 Frontmatter（如果需要）

```yaml
---
title: Base64 Data
encoded: SGVsbG8gV29ybGQ=
---
```

---

## 带 HTML 的 Frontmatter

```yaml
---
title: HTML in Frontmatter
description: "<strong>Bold</strong> and <em>italic</em>"
---
```

---

## 带 Markdown 的 Frontmatter

```yaml
---
title: Markdown in Frontmatter
description: "**Bold** and *italic* text"
---
```

---

## 带很长字段名的 Frontmatter

```yaml
---
title: Very Long Field Names
very_long_field_name_that_extends_beyond_normal_width: value
another_extremely_long_field_name_with_many_words: another value
---
```

---

## 带特殊 YAML 功能的 Frontmatter

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

## 以三个点结尾的 Frontmatter

```yaml
---
title: Three Dots Ending
author: Author Name
...
```

---

## 带 TOML 样式的 Frontmatter（如果支持）

```toml
+++
title = "TOML Frontmatter"
author = "Author Name"
+++
```

---

## 带 JSON 的 Frontmatter（如果支持）

```json
{
  "title": "JSON Frontmatter",
  "author": "Author Name"
}
```

---

## 带多个文档分隔符的 Frontmatter

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

## 带内联数组的 Frontmatter

```yaml
---
title: Inline Arrays
tags: [tag1, tag2, tag3]
categories: [cat1, cat2]
---
```

---

## 带混合引号的 Frontmatter

```yaml
---
title: 'Single quoted title'
description: "Double quoted description"
author: Unquoted author
---
```

---

## 保留顺序的 Frontmatter

```yaml
---
z_field: Last
a_field: First
m_field: Middle
---
```

