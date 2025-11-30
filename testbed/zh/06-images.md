---
title: 图片 - 内联和引用样式
categories:
  - markdown
  - images
tags:
  - graphics
  - media
  - pictures
---

# 图片

## 内联图片

![替代文本](image.png)

![带空格的替代文本](image with spaces.png)

![替代文本](image.png "图片标题")

![替代文本](image.png '图片标题')

![替代文本](image.png (图片标题))

---

## 不同格式的图片

![PNG 图片](image.png)

![JPEG 图片](image.jpg)

![GIF 图片](image.gif)

![SVG 图片](image.svg)

![WebP 图片](image.webp)

---

## 带路径的图片

![根图片](/image.png)

![子目录图片](./images/image.png)

![父目录](../images/image.png)

![深层路径](./images/subfolder/image.png)

---

## 带 URL 的图片

![远程图片](https://example.com/image.png)

![带路径的远程图片](https://example.com/images/photo.jpg)

![带查询的远程图片](https://example.com/image.png?size=large)

---

## 引用样式图片

![引用图片][img1]

![隐式引用图片][]

![快捷引用图片]

[img1]: image1.png
[implicit reference image]: image2.png
[shortcut reference image]: image3.png

---

## 带标题的引用图片

![带标题的引用图片][img2]

[img2]: image.png "图片标题"

![带单引号的引用图片][img3]

[img3]: image.png '图片标题'

![带括号的引用图片][img4]

[img4]: image.png (图片标题)

---

## 替代文本中带特殊字符的图片

![带 "引号" 的图片](image.png)

![带（括号）的图片](image.png)

![带 [方括号] 的图片](image.png)

![带 {大括号} 的图片](image.png)

![带 <标签> 的图片](image.png)

![带 $美元$ 的图片](image.png)

---

## 替代文本中带 Unicode 的图片

![带 café 的图片](image.png)

![带 你好 的图片](image.png)

![带 こんにちは 的图片](image.png)

![带 مرحبا 的图片](image.png)

---

## 替代文本中带数字的图片

![图片 123](image.png)

![图片 v1.2.3](image.png)

![图片 50%](image.png)

---

## 替代文本中带强调的图片

![**粗体替代文本**](image.png)

![*斜体替代文本*](image.png)

![***粗体和斜体替代文本***](image.png)

---

## 带链接的图片

[![链接图片](image.png)](https://example.com)

[![带替代文本的链接图片](image.png "图片标题")](https://example.com "链接标题")

---

## 列表中的图片

* 带 ![图片](image.png) 的项目
* 带 ![引用图片][listimg] 的项目
* 带 ![图片](image.png "标题") 的项目

[listimg]: image.png

---

## 引用块中的图片

> 这是带 ![图片](image.png) 的引用块

> 这是带 ![引用图片][blockquoteimg] 的引用块

[blockquoteimg]: image.png

---

## 标题中的图片

# 带 ![图片](image.png) 的标题

## 带 ![引用图片][headerimg] 的标题

[headerimg]: image.png

---

## 表格中的图片

| Column 1 | Column 2 |
|----------|----------|
| ![图片](image.png) | ![引用][tableimg] |
| 文本 | ![带标题的图片](image.png "标题") |

[tableimg]: image.png

---

## 空替代文本的图片

![](image.png)

![ ](image.png)

![   ](image.png)

---

## 替代文本只有空格的图片

![   ](image.png)

---

## 长替代文本的图片

![这是一个非常长的替代文本，详细描述了图片，并提供了关于图片内容和其在文档中目的的全面信息](image.png)

---

## 一个段落中的多个图片

这个段落有 ![第一张图片](image1.png)、![第二张图片](image2.png) 和 ![第三张图片](image3.png)。

这个段落混合了 ![内联图片](image.png)、![引用图片][multimg] 和文本。

[multimg]: image.png

---

## 带 HTML 的图片

<img src="image.png" alt="HTML 图片">

![Markdown 图片](image.png) 和 <img src="image.png" alt="HTML 图片">

---

## 带转义字符的图片

\![转义的图片](image.png)

![带 \*星号\* 的图片](image.png)

![带 \_下划线\_ 的图片](image.png)

---

## 代码块中的图片

```
![图片](image.png) 在代码块中不应起作用
```

---

## 带特殊 URL 方案的图片

![Data URI 图片](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==)

---

## 带百分比编码的图片

![带编码的图片](https://example.com/image%20with%20spaces.png)

![带 unicode 编码的图片](https://example.com/%E4%BD%A0%E5%A5%BD.png)

---

## 复杂的图片组合

**粗体文本** 与 ![图片](image.png) 和 *斜体* 文本。

![图片](image.png) 与 [链接](https://example.com) 和 `代码`。

![**粗体替代文本**](image.png) 与 *斜体* 文本和 [链接](https://example.com)。

---

## 不同大小的图片（HTML 属性）

<img src="image.png" alt="小图片" width="100" height="100">

<img src="image.png" alt="大图片" width="800" height="600">

---

## 带 CSS 类的图片（HTML）

<img src="image.png" alt="样式化图片" class="custom-class">

<img src="image.png" alt="居中图片" style="display: block; margin: 0 auto;">

---

## 带标题和说明的图片

![图片](image.png "这是一个标题")

![图片](image.png "带 '引号' 的标题")

![图片](image.png "带（括号）的标题")

---

## 带相对路径的图片

![同一目录](./image.png)

![子目录](./images/image.png)

![父目录](../image.png)

![根相对](/image.png)

---

## 带绝对 URL 的图片

![绝对 URL](https://example.com/image.png)

![带路径的绝对 URL](https://example.com/images/photo.jpg)

![带查询的绝对 URL](https://example.com/image.png?v=1&size=large)

![带片段的绝对 URL](https://example.com/image.png#section)

---

## 嵌套结构中的图片

### 嵌套列表中的图片

* 第一级
  * 带 ![图片](image.png) 的第二级
    * 带 ![图片](image.png) 的第三级

### 嵌套引用块中的图片

> 外部引用块
> > 带 ![图片](image.png) 的内部引用块

---

## 断开的引用图片

![断开的图片](nonexistent.png)

![断开的引用][broken]

[broken]: nonexistent.png

---

## 带特殊文件名的图片

![带空格的图片](image with spaces.png)

![带点的图片](image.with.dots.png)

![带破折号的图片](image-with-dashes.png)

![带下划线的图片](Image_with_underscores.png)

![图片123](Image123.png)

![图片@特殊](Image@special.png)

