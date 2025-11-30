---
title: Images - Inline and Reference Style
categories:
  - markdown
  - images
tags:
  - graphics
  - media
  - pictures
---

# Images

## Inline Images

![Alt text](image.png)

![Alt text with spaces](image with spaces.png)

![Alt text](image.png "Image Title")

![Alt text](image.png 'Image Title')

![Alt text](image.png (Image Title))

---

## Images with Different Formats

![PNG image](image.png)

![JPEG image](image.jpg)

![GIF image](image.gif)

![SVG image](image.svg)

![WebP image](image.webp)

---

## Images with Paths

![Root image](/image.png)

![Subdirectory image](./images/image.png)

![Parent directory](../images/image.png)

![Deep path](./images/subfolder/image.png)

---

## Images with URLs

![Remote image](https://example.com/image.png)

![Remote image with path](https://example.com/images/photo.jpg)

![Remote image with query](https://example.com/image.png?size=large)

---

## Reference-Style Images

![Reference image][img1]

![Implicit reference image][]

![Shortcut reference image]

[img1]: image1.png
[implicit reference image]: image2.png
[shortcut reference image]: image3.png

---

## Reference Images with Titles

![Reference image with title][img2]

[img2]: image.png "Image Title"

![Reference image with single quotes][img3]

[img3]: image.png 'Image Title'

![Reference image with parentheses][img4]

[img4]: image.png (Image Title)

---

## Images with Special Characters in Alt Text

![Image with "quotes"](image.png)

![Image with (parentheses)](image.png)

![Image with [brackets]](image.png)

![Image with {braces}](image.png)

![Image with <tags>](image.png)

![Image with $dollar$](image.png)

---

## Images with Unicode in Alt Text

![Image with café](image.png)

![Image with 你好](image.png)

![Image with こんにちは](image.png)

![Image with مرحبا](image.png)

---

## Images with Numbers in Alt Text

![Image 123](image.png)

![Image v1.2.3](image.png)

![Image 50%](image.png)

---

## Images with Emphasis in Alt Text

![**Bold alt text**](image.png)

![*Italic alt text*](image.png)

![***Bold and italic alt text***](image.png)

---

## Images with Links

[![Linked image](image.png)](https://example.com)

[![Linked image with alt](image.png "Image Title")](https://example.com "Link Title")

---

## Images in Lists

* Item with ![image](image.png)
* Item with ![reference image][listimg]
* Item with ![image](image.png "Title")

[listimg]: image.png

---

## Images in Blockquotes

> This is a blockquote with ![image](image.png)

> This is a blockquote with ![reference image][blockquoteimg]

[blockquoteimg]: image.png

---

## Images in Headers

# Header with ![image](image.png)

## Header with ![reference image][headerimg]

[headerimg]: image.png

---

## Images in Tables

| Column 1 | Column 2 |
|----------|----------|
| ![Image](image.png) | ![Reference][tableimg] |
| Text | ![Image with title](image.png "Title") |

[tableimg]: image.png

---

## Images with Empty Alt Text

![](image.png)

![ ](image.png)

![   ](image.png)

---

## Images with Only Spaces in Alt Text

![   ](image.png)

---

## Images with Long Alt Text

![This is a very long alt text that describes the image in great detail and provides comprehensive information about what the image contains and its purpose in the document](image.png)

---

## Multiple Images in One Paragraph

This paragraph has ![first image](image1.png), ![second image](image2.png), and ![third image](image3.png).

This paragraph mixes ![inline image](image.png), ![reference image][multimg], and text.

[multimg]: image.png

---

## Images with HTML

<img src="image.png" alt="HTML image">

![Markdown image](image.png) and <img src="image.png" alt="HTML image">

---

## Images with Escaped Characters

\![Escaped image](image.png)

![Image with \*asterisk\*](image.png)

![Image with \_underscore\_](image.png)

---

## Images in Code Blocks

```
![Image](image.png) should not work in code blocks
```

---

## Images with Special URL Schemes

![Data URI image](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==)

---

## Images with Percent Encoding

![Image with encoding](https://example.com/image%20with%20spaces.png)

![Image with unicode encoding](https://example.com/%E4%BD%A0%E5%A5%BD.png)

---

## Complex Image Combinations

**Bold text** with ![image](image.png) and *italic* text.

![Image](image.png) with [link](https://example.com) and `code`.

![**Bold alt**](image.png) with *italic* text and [link](https://example.com).

---

## Images with Different Sizes (HTML attributes)

<img src="image.png" alt="Small image" width="100" height="100">

<img src="image.png" alt="Large image" width="800" height="600">

---

## Images with CSS Classes (HTML)

<img src="image.png" alt="Styled image" class="custom-class">

<img src="image.png" alt="Centered image" style="display: block; margin: 0 auto;">

---

## Images with Titles and Captions

![Image](image.png "This is a title")

![Image](image.png "Title with 'quotes'")

![Image](image.png "Title with (parentheses)")

---

## Images with Relative Paths

![Same directory](./image.png)

![Subdirectory](./images/image.png)

![Parent directory](../image.png)

![Root relative](/image.png)

---

## Images with Absolute URLs

![Absolute URL](https://example.com/image.png)

![Absolute URL with path](https://example.com/images/photo.jpg)

![Absolute URL with query](https://example.com/image.png?v=1&size=large)

![Absolute URL with fragment](https://example.com/image.png#section)

---

## Images in Nested Structures

### Images in Nested Lists

* First level
  * Second level with ![image](image.png)
    * Third level with ![image](image.png)

### Images in Nested Blockquotes

> Outer blockquote
> > Inner blockquote with ![image](image.png)

---

## Images with Broken References

![Broken image](nonexistent.png)

![Broken reference][broken]

[broken]: nonexistent.png

---

## Images with Special File Names

![Image with spaces](image with spaces.png)

![Image.with.dots](image.with.dots.png)

![Image-with-dashes](image-with-dashes.png)

![Image_with_underscores](Image_with_underscores.png)

![Image123](Image123.png)

![Image@special](Image@special.png)

