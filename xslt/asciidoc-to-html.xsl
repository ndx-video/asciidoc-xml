<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="1.0"
    xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns:ad="https://github.com/ndx-video/asciidoc-xml"
    exclude-result-prefixes="ad">

    <xsl:output method="html" encoding="UTF-8" indent="yes" omit-xml-declaration="yes"/>

    <!-- Root template - creates HTML wrapper -->
    <xsl:template match="/">
        <html xmlns="http://www.w3.org/1999/xhtml">
            <head>
                <meta charset="UTF-8"/>
                <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
                <title>
                    <xsl:choose>
                        <xsl:when test="//ad:document/@title">
                            <xsl:value-of select="//ad:document/@title"/>
                        </xsl:when>
                        <xsl:otherwise>AsciiDoc Document</xsl:otherwise>
                    </xsl:choose>
                </title>
                <style>
                    body { 
                        font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
                        line-height: 1.6;
                        max-width: 1200px;
                        margin: 0 auto;
                        padding: 2rem;
                        color: #333;
                    }
                    .document-header { 
                        margin-bottom: 2rem;
                        border-bottom: 2px solid #e0e0e0;
                        padding-bottom: 1rem;
                    }
                    .document-title { 
                        font-size: 2.5rem;
                        margin: 0 0 0.5rem 0;
                        color: #1a1a1a;
                    }
                    .document-author { 
                        color: #666;
                        font-size: 1.1rem;
                    }
                    .document-author a { color: #0066cc; }
                    .document-revision { 
                        color: #888;
                        font-size: 0.9rem;
                        margin-top: 0.5rem;
                    }
                    .preamble { 
                        margin: 2rem 0;
                        font-size: 1.1rem;
                        color: #555;
                    }
                    section { margin: 2rem 0; }
                    h2, h3, h4, h5, h6 { 
                        margin-top: 1.5rem;
                        margin-bottom: 0.5rem;
                        color: #1a1a1a;
                    }
                    h2 { font-size: 2rem; border-bottom: 1px solid #e0e0e0; padding-bottom: 0.3rem; }
                    h3 { font-size: 1.6rem; }
                    h4 { font-size: 1.3rem; }
                    h5 { font-size: 1.1rem; }
                    h6 { font-size: 1rem; }
                    p { margin: 1rem 0; }
                    code { 
                        background: #f4f4f4;
                        padding: 0.2rem 0.4rem;
                        border-radius: 3px;
                        font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
                        font-size: 0.9em;
                    }
                    pre { 
                        background: #f8f8f8;
                        border: 1px solid #ddd;
                        border-left: 3px solid #0066cc;
                        padding: 1rem;
                        overflow-x: auto;
                        border-radius: 4px;
                        margin: 1rem 0;
                    }
                    pre code {
                        background: transparent;
                        padding: 0;
                    }
                    .codeblock-title {
                        background: #e8e8e8;
                        padding: 0.3rem 0.5rem;
                        font-weight: bold;
                        border: 1px solid #ddd;
                        border-bottom: none;
                        border-radius: 4px 4px 0 0;
                        margin-bottom: -1rem;
                        margin-top: 1rem;
                    }
                    .example { 
                        background: #f0f8ff;
                        border: 1px solid #b8d4e8;
                        border-left: 4px solid #0066cc;
                        padding: 1rem;
                        margin: 1rem 0;
                        border-radius: 4px;
                    }
                    .example-title {
                        font-weight: bold;
                        color: #0066cc;
                        margin-bottom: 0.5rem;
                    }
                    .sidebar { 
                        background: #fffef0;
                        border: 1px solid #e8e0c0;
                        border-left: 4px solid #ccaa00;
                        padding: 1rem;
                        margin: 1rem 0;
                        border-radius: 4px;
                    }
                    .sidebar-title {
                        font-weight: bold;
                        color: #997700;
                        margin-bottom: 0.5rem;
                    }
                    blockquote { 
                        border-left: 4px solid #ddd;
                        padding-left: 1rem;
                        margin: 1rem 0;
                        color: #666;
                        font-style: italic;
                    }
                    .quote-attribution {
                        display: block;
                        margin-top: 0.5rem;
                        font-size: 0.9rem;
                        text-align: right;
                        font-weight: bold;
                    }
                    .verseblock {
                        white-space: pre-wrap;
                        font-family: serif;
                        margin: 1rem 0;
                        padding: 1rem;
                        background: #fafafa;
                        border-left: 3px solid #888;
                    }
                    .admonition {
                        padding: 1rem;
                        margin: 1rem 0;
                        border-left: 4px solid;
                        border-radius: 4px;
                    }
                    .admonition-note { 
                        background: #e7f3ff;
                        border-left-color: #0066cc;
                    }
                    .admonition-tip { 
                        background: #e7ffe7;
                        border-left-color: #00cc00;
                    }
                    .admonition-important { 
                        background: #fff4e6;
                        border-left-color: #ff9900;
                    }
                    .admonition-warning { 
                        background: #fff0f0;
                        border-left-color: #ff3333;
                    }
                    .admonition-caution { 
                        background: #ffe6e6;
                        border-left-color: #cc0000;
                    }
                    .admonition-title {
                        font-weight: bold;
                        margin-bottom: 0.5rem;
                    }
                    table { 
                        border-collapse: collapse;
                        width: 100%;
                        margin: 1rem 0;
                    }
                    th, td { 
                        border: 1px solid #ddd;
                        padding: 0.5rem;
                        text-align: left;
                    }
                    th { 
                        background: #f4f4f4;
                        font-weight: bold;
                    }
                    tr:nth-child(even) { background: #fafafa; }
                    .table-title {
                        font-weight: bold;
                        margin-bottom: 0.5rem;
                        color: #333;
                    }
                    ul, ol { 
                        margin: 0.5rem 0;
                        padding-left: 2rem;
                    }
                    li { margin: 0.3rem 0; }
                    dl { margin: 1rem 0; }
                    dt { 
                        font-weight: bold;
                        margin-top: 0.5rem;
                    }
                    dd { 
                        margin-left: 2rem;
                        margin-bottom: 0.5rem;
                    }
                    hr { 
                        border: none;
                        border-top: 2px solid #ddd;
                        margin: 2rem 0;
                    }
                    .pagebreak {
                        page-break-after: always;
                        border-top: 3px dashed #ccc;
                        margin: 2rem 0;
                        padding-top: 2rem;
                    }
                    mark { 
                        background: #ffeb3b;
                        padding: 0.1rem 0.2rem;
                    }
                    a { 
                        color: #0066cc;
                        text-decoration: none;
                    }
                    a:hover { text-decoration: underline; }
                    kbd {
                        background: #f4f4f4;
                        border: 1px solid #ccc;
                        border-radius: 3px;
                        padding: 0.1rem 0.4rem;
                        font-family: monospace;
                        font-size: 0.85em;
                    }
                    .menu-item {
                        font-weight: bold;
                    }
                    .menu-separator {
                        margin: 0 0.3rem;
                        color: #666;
                    }
                    .footnote {
                        font-size: 0.85em;
                        vertical-align: super;
                        color: #0066cc;
                    }
                    figure {
                        margin: 1rem 0;
                        text-align: center;
                    }
                    figcaption {
                        margin-top: 0.5rem;
                        font-style: italic;
                        color: #666;
                    }
                    img {
                        max-width: 100%;
                        height: auto;
                    }
                </style>
            </head>
            <body>
                <xsl:apply-templates/>
            </body>
        </html>
    </xsl:template>

    <!-- Document root -->
    <xsl:template match="ad:document">
        <article class="document">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            
            <!-- Document header -->
            <xsl:if test="@title or @author or @revnumber">
                <div class="document-header">
                    <xsl:if test="@title">
                        <h1 class="document-title"><xsl:value-of select="@title"/></h1>
                    </xsl:if>
                    <xsl:if test="@author">
                        <div class="document-author">
                            <xsl:value-of select="@author"/>
                            <xsl:if test="@email">
                                <xsl:text> </xsl:text>
                                <a href="mailto:{@email}"><xsl:value-of select="@email"/></a>
                            </xsl:if>
                        </div>
                    </xsl:if>
                    <xsl:if test="@revnumber or @revdate or @revremark">
                        <div class="document-revision">
                            <xsl:if test="@revnumber">
                                <span class="revision-number">Version <xsl:value-of select="@revnumber"/></span>
                            </xsl:if>
                            <xsl:if test="@revdate">
                                <xsl:if test="@revnumber"><xsl:text>, </xsl:text></xsl:if>
                                <span class="revision-date"><xsl:value-of select="@revdate"/></span>
                            </xsl:if>
                            <xsl:if test="@revremark">
                                <xsl:text> — </xsl:text>
                                <span class="revision-remark"><xsl:value-of select="@revremark"/></span>
                            </xsl:if>
                        </div>
                    </xsl:if>
                </div>
            </xsl:if>
            
            <xsl:apply-templates/>
        </article>
    </xsl:template>

    <!-- Preamble -->
    <xsl:template match="ad:preamble">
        <div class="preamble">
            <xsl:apply-templates/>
        </div>
    </xsl:template>

    <!-- Sections -->
    <xsl:template match="ad:section">
        <section>
            <xsl:attribute name="class">
                <xsl:text>section level-</xsl:text>
                <xsl:value-of select="@level"/>
                <xsl:if test="@role">
                    <xsl:text> </xsl:text>
                    <xsl:value-of select="@role"/>
                </xsl:if>
            </xsl:attribute>
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            
            <!-- Section title -->
            <xsl:if test="@title">
                <xsl:variable name="level" select="number(@level) + 1"/>
                <xsl:choose>
                    <xsl:when test="$level = 1">
                        <h1><xsl:value-of select="@title"/></h1>
                    </xsl:when>
                    <xsl:when test="$level = 2">
                        <h2><xsl:value-of select="@title"/></h2>
                    </xsl:when>
                    <xsl:when test="$level = 3">
                        <h3><xsl:value-of select="@title"/></h3>
                    </xsl:when>
                    <xsl:when test="$level = 4">
                        <h4><xsl:value-of select="@title"/></h4>
                    </xsl:when>
                    <xsl:when test="$level = 5">
                        <h5><xsl:value-of select="@title"/></h5>
                    </xsl:when>
                    <xsl:otherwise>
                        <h6><xsl:value-of select="@title"/></h6>
                    </xsl:otherwise>
                </xsl:choose>
            </xsl:if>
            
            <xsl:apply-templates/>
        </section>
    </xsl:template>

    <!-- Paragraphs -->
    <xsl:template match="ad:paragraph">
        <p>
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@role">
                <xsl:attribute name="class"><xsl:value-of select="@role"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates/>
        </p>
    </xsl:template>

    <!-- Code blocks -->
    <xsl:template match="ad:codeblock">
        <div class="codeblock">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@title">
                <div class="codeblock-title"><xsl:value-of select="@title"/></div>
            </xsl:if>
            <pre>
                <xsl:if test="@language">
                    <xsl:attribute name="class">language-<xsl:value-of select="@language"/></xsl:attribute>
                </xsl:if>
                <code><xsl:value-of select="."/></code>
            </pre>
        </div>
    </xsl:template>

    <!-- Literal blocks -->
    <xsl:template match="ad:literalblock">
        <pre class="literal">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:value-of select="."/>
        </pre>
    </xsl:template>

    <!-- Example blocks -->
    <xsl:template match="ad:example">
        <div class="example">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@title">
                <div class="example-title"><xsl:value-of select="@title"/></div>
            </xsl:if>
            <div class="example-content">
                <xsl:apply-templates/>
            </div>
        </div>
    </xsl:template>

    <!-- Sidebar blocks -->
    <xsl:template match="ad:sidebar">
        <aside class="sidebar">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@title">
                <div class="sidebar-title"><xsl:value-of select="@title"/></div>
            </xsl:if>
            <div class="sidebar-content">
                <xsl:apply-templates/>
            </div>
        </aside>
    </xsl:template>

    <!-- Quote blocks -->
    <xsl:template match="ad:quote">
        <blockquote>
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates/>
            <xsl:if test="@attribution">
                <footer class="quote-attribution">
                    <cite>— <xsl:value-of select="@attribution"/></cite>
                    <xsl:if test="@citation">
                        <xsl:text>, </xsl:text>
                        <cite><xsl:value-of select="@citation"/></cite>
                    </xsl:if>
                </footer>
            </xsl:if>
        </blockquote>
    </xsl:template>

    <!-- Verse blocks -->
    <xsl:template match="ad:verseblock">
        <div class="verseblock">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates/>
            <xsl:if test="@attribution">
                <footer class="quote-attribution">
                    <cite>— <xsl:value-of select="@attribution"/></cite>
                </footer>
            </xsl:if>
        </div>
    </xsl:template>

    <!-- Open blocks -->
    <xsl:template match="ad:openblock">
        <div class="openblock">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@role">
                <xsl:attribute name="class">openblock <xsl:value-of select="@role"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates/>
        </div>
    </xsl:template>

    <!-- Tables -->
    <xsl:template match="ad:table">
        <div class="table-wrapper">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@title">
                <div class="table-title"><xsl:value-of select="@title"/></div>
            </xsl:if>
            <table>
                <xsl:if test="@role">
                    <xsl:attribute name="class"><xsl:value-of select="@role"/></xsl:attribute>
                </xsl:if>
                <xsl:choose>
                    <xsl:when test="ad:row[1]/@role = 'header'">
                        <thead>
                            <xsl:apply-templates select="ad:row[1]" mode="header"/>
                        </thead>
                        <tbody>
                            <xsl:apply-templates select="ad:row[position() > 1]"/>
                        </tbody>
                    </xsl:when>
                    <xsl:otherwise>
                        <tbody>
                            <xsl:apply-templates select="ad:row"/>
                        </tbody>
                    </xsl:otherwise>
                </xsl:choose>
            </table>
        </div>
    </xsl:template>

    <!-- Table rows -->
    <xsl:template match="ad:row">
        <tr>
            <xsl:apply-templates select="ad:cell"/>
        </tr>
    </xsl:template>

    <xsl:template match="ad:row" mode="header">
        <tr>
            <xsl:apply-templates select="ad:cell" mode="header"/>
        </tr>
    </xsl:template>

    <!-- Table cells -->
    <xsl:template match="ad:cell">
        <td>
            <xsl:if test="@align">
                <xsl:attribute name="style">text-align: <xsl:value-of select="@align"/>;</xsl:attribute>
            </xsl:if>
            <xsl:if test="@colspan">
                <xsl:attribute name="colspan"><xsl:value-of select="@colspan"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@rowspan">
                <xsl:attribute name="rowspan"><xsl:value-of select="@rowspan"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates/>
        </td>
    </xsl:template>

    <xsl:template match="ad:cell" mode="header">
        <th>
            <xsl:if test="@align">
                <xsl:attribute name="style">text-align: <xsl:value-of select="@align"/>;</xsl:attribute>
            </xsl:if>
            <xsl:if test="@colspan">
                <xsl:attribute name="colspan"><xsl:value-of select="@colspan"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@rowspan">
                <xsl:attribute name="rowspan"><xsl:value-of select="@rowspan"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates/>
        </th>
    </xsl:template>

    <!-- Lists -->
    <xsl:template match="ad:list[@style='unordered']">
        <ul>
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@role">
                <xsl:attribute name="class"><xsl:value-of select="@role"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates select="ad:listitem"/>
        </ul>
    </xsl:template>

    <xsl:template match="ad:list[@style='ordered']">
        <ol>
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@role">
                <xsl:attribute name="class"><xsl:value-of select="@role"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates select="ad:listitem"/>
        </ol>
    </xsl:template>

    <xsl:template match="ad:list[@style='labeled']">
        <dl>
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@role">
                <xsl:attribute name="class"><xsl:value-of select="@role"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates select="ad:item" mode="labeled"/>
        </dl>
    </xsl:template>

    <xsl:template match="ad:list[@style='callout']">
        <ol class="callout-list">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates select="ad:listitem"/>
        </ol>
    </xsl:template>

    <!-- List items -->
    <xsl:template match="ad:listitem">
        <li>
            <xsl:if test="@callout">
                <xsl:attribute name="data-callout"><xsl:value-of select="@callout"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates/>
        </li>
    </xsl:template>

    <xsl:template match="ad:listitem" mode="labeled">
        <xsl:if test="@term">
            <dt><xsl:value-of select="@term"/></dt>
        </xsl:if>
        <dd>
            <xsl:apply-templates/>
        </dd>
    </xsl:template>

    <!-- Admonitions -->
    <xsl:template match="ad:admonition">
        <div>
            <xsl:attribute name="class">
                <xsl:text>admonition admonition-</xsl:text>
                <xsl:value-of select="@type"/>
                <xsl:if test="@role">
                    <xsl:text> </xsl:text>
                    <xsl:value-of select="@role"/>
                </xsl:if>
            </xsl:attribute>
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@title">
                <div class="admonition-title">
                    <xsl:value-of select="@title"/>
                </div>
            </xsl:if>
            <xsl:if test="not(@title)">
                <div class="admonition-title">
                    <xsl:choose>
                        <xsl:when test="@type='note'">Note</xsl:when>
                        <xsl:when test="@type='tip'">Tip</xsl:when>
                        <xsl:when test="@type='important'">Important</xsl:when>
                        <xsl:when test="@type='warning'">Warning</xsl:when>
                        <xsl:when test="@type='caution'">Caution</xsl:when>
                        <xsl:otherwise><xsl:value-of select="@type"/></xsl:otherwise>
                    </xsl:choose>
                </div>
            </xsl:if>
            <div class="admonition-content">
                <xsl:apply-templates/>
            </div>
        </div>
    </xsl:template>

    <!-- Images -->
    <xsl:template match="ad:image">
        <figure class="image">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <img src="{@src}">
                <xsl:if test="@alt">
                    <xsl:attribute name="alt"><xsl:value-of select="@alt"/></xsl:attribute>
                </xsl:if>
                <xsl:if test="@width">
                    <xsl:attribute name="width"><xsl:value-of select="@width"/></xsl:attribute>
                </xsl:if>
                <xsl:if test="@height">
                    <xsl:attribute name="height"><xsl:value-of select="@height"/></xsl:attribute>
                </xsl:if>
            </img>
            <xsl:if test="@title">
                <figcaption><xsl:value-of select="@title"/></figcaption>
            </xsl:if>
        </figure>
    </xsl:template>

    <!-- Thematic break -->
    <xsl:template match="ad:thematicbreak">
        <hr/>
    </xsl:template>

    <!-- Page break -->
    <xsl:template match="ad:pagebreak">
        <div class="pagebreak"></div>
    </xsl:template>

    <!-- Macros -->
    <xsl:template match="ad:macro[@type='block' and @name='image']">
        <figure class="image">
            <img src="{@src}">
                <xsl:if test="@alt">
                    <xsl:attribute name="alt"><xsl:value-of select="@alt"/></xsl:attribute>
                </xsl:if>
            </img>
            <xsl:if test="@title">
                <figcaption><xsl:value-of select="@title"/></figcaption>
            </xsl:if>
        </figure>
    </xsl:template>

    <!-- Generic block macros -->
    <xsl:template match="ad:macro[@type='block']">
        <div class="macro macro-{@name}">
            <xsl:apply-templates/>
        </div>
    </xsl:template>

    <!-- Generic inline macros -->
    <xsl:template match="ad:macro[@type='inline']">
        <span class="macro macro-{@name}">
            <xsl:apply-templates/>
        </span>
    </xsl:template>

    <!-- Inline elements -->
    <xsl:template match="ad:strong">
        <strong><xsl:apply-templates/></strong>
    </xsl:template>

    <xsl:template match="ad:emphasis">
        <em><xsl:apply-templates/></em>
    </xsl:template>

    <xsl:template match="ad:monospace">
        <code><xsl:apply-templates/></code>
    </xsl:template>

    <xsl:template match="ad:superscript">
        <sup><xsl:apply-templates/></sup>
    </xsl:template>

    <xsl:template match="ad:subscript">
        <sub><xsl:apply-templates/></sub>
    </xsl:template>

    <xsl:template match="ad:highlight">
        <mark><xsl:apply-templates/></mark>
    </xsl:template>

    <!-- Links -->
    <xsl:template match="ad:link">
        <a href="{@href}">
            <xsl:if test="@title">
                <xsl:attribute name="title"><xsl:value-of select="@title"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@window = '_blank'">
                <xsl:attribute name="target">_blank</xsl:attribute>
                <xsl:attribute name="rel">noopener noreferrer</xsl:attribute>
            </xsl:if>
            <xsl:choose>
                <xsl:when test="node()">
                    <xsl:apply-templates/>
                </xsl:when>
                <xsl:otherwise>
                    <xsl:value-of select="@href"/>
                </xsl:otherwise>
            </xsl:choose>
        </a>
    </xsl:template>

    <!-- Anchors -->
    <xsl:template match="ad:anchor">
        <a id="{@id}"></a>
    </xsl:template>

    <!-- Footnotes -->
    <xsl:template match="ad:footnote">
        <sup class="footnote">
            <xsl:text>[</xsl:text>
            <xsl:choose>
                <xsl:when test="@ref">
                    <xsl:value-of select="@ref"/>
                </xsl:when>
                <xsl:otherwise>*</xsl:otherwise>
            </xsl:choose>
            <xsl:text>]</xsl:text>
        </sup>
    </xsl:template>

    <!-- Passthrough -->
    <xsl:template match="ad:passthrough">
        <xsl:value-of select="." disable-output-escaping="yes"/>
    </xsl:template>

    <!-- Text nodes -->
    <xsl:template match="text()">
        <xsl:value-of select="."/>
    </xsl:template>

</xsl:stylesheet>
