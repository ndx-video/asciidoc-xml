<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="1.0"
    xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns:ad="http://asciidoc.org/ns"
    exclude-result-prefixes="ad">

    <xsl:output method="html" encoding="UTF-8" indent="yes" omit-xml-declaration="yes"/>

    <!-- Root template -->
    <xsl:template match="/">
        <xsl:apply-templates/>
    </xsl:template>

    <!-- Document root -->
    <xsl:template match="ad:asciidoc">
        <div class="asciidoc-document" data-doctype="{@doctype}">
            <xsl:apply-templates/>
        </div>
    </xsl:template>

    <!-- Document header -->
    <xsl:template match="ad:header">
        <header class="asciidoc-header">
            <xsl:if test="ad:title">
                <h1 class="asciidoc-title">
                    <xsl:apply-templates select="ad:title"/>
                </h1>
            </xsl:if>
            <xsl:if test="ad:author">
                <div class="asciidoc-authors">
                    <xsl:for-each select="ad:author">
                        <div class="asciidoc-author">
                            <span class="author-name"><xsl:value-of select="ad:name"/></span>
                            <xsl:if test="ad:email">
                                <span class="author-email"> &lt;<xsl:value-of select="ad:email"/>&gt;</span>
                            </xsl:if>
                        </div>
                    </xsl:for-each>
                </div>
            </xsl:if>
            <xsl:if test="ad:revision">
                <div class="asciidoc-revision">
                    <xsl:if test="ad:revision/ad:number">
                        <span class="revision-number">Version <xsl:value-of select="ad:revision/ad:number"/></span>
                    </xsl:if>
                    <xsl:if test="ad:revision/ad:date">
                        <span class="revision-date"><xsl:value-of select="ad:revision/ad:date"/></span>
                    </xsl:if>
                    <xsl:if test="ad:revision/ad:remark">
                        <span class="revision-remark"><xsl:value-of select="ad:revision/ad:remark"/></span>
                    </xsl:if>
                </div>
            </xsl:if>
        </header>
    </xsl:template>

    <!-- Content container -->
    <xsl:template match="ad:content">
        <div class="asciidoc-content">
            <xsl:apply-templates/>
        </div>
    </xsl:template>

    <!-- Sections -->
    <xsl:template match="ad:section">
        <section class="asciidoc-section asciidoc-section-level-{@level}" 
                 data-level="{@level}">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@role">
                <xsl:attribute name="class">
                    <xsl:value-of select="concat('asciidoc-section asciidoc-section-level-', @level, ' ', @role)"/>
                </xsl:attribute>
            </xsl:if>
            <xsl:if test="ad:title">
                <xsl:element name="h{@level + 1}">
                    <xsl:attribute name="class">asciidoc-section-title</xsl:attribute>
                    <xsl:apply-templates select="ad:title"/>
                </xsl:element>
            </xsl:if>
            <xsl:if test="ad:content">
                <div class="asciidoc-section-content">
                    <xsl:apply-templates select="ad:content"/>
                </div>
            </xsl:if>
        </section>
    </xsl:template>

    <!-- Paragraphs -->
    <xsl:template match="ad:paragraph">
        <p class="asciidoc-paragraph">
            <xsl:if test="@role">
                <xsl:attribute name="class">asciidoc-paragraph <xsl:value-of select="@role"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates/>
        </p>
    </xsl:template>

    <!-- Code blocks -->
    <xsl:template match="ad:codeblock">
        <div class="asciidoc-codeblock">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@title">
                <div class="asciidoc-codeblock-title"><xsl:value-of select="@title"/></div>
            </xsl:if>
            <pre class="asciidoc-code">
                <xsl:if test="@language">
                    <xsl:attribute name="data-language"><xsl:value-of select="@language"/></xsl:attribute>
                    <xsl:attribute name="class">asciidoc-code language-<xsl:value-of select="@language"/></xsl:attribute>
                </xsl:if>
                <code>
                    <xsl:if test="@language">
                        <xsl:attribute name="class">language-<xsl:value-of select="@language"/></xsl:attribute>
                    </xsl:if>
                    <xsl:value-of select="text()"/>
                </code>
            </pre>
            <xsl:if test="@source">
                <div class="asciidoc-codeblock-source">
                    <small>Source: <xsl:value-of select="@source"/></small>
                </div>
            </xsl:if>
        </div>
    </xsl:template>

    <!-- Literal blocks -->
    <xsl:template match="ad:literalblock">
        <div class="asciidoc-literalblock">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <pre class="asciidoc-literal">
                <xsl:if test="@style">
                    <xsl:attribute name="class">asciidoc-literal asciidoc-literal-<xsl:value-of select="@style"/></xsl:attribute>
                </xsl:if>
                <xsl:value-of select="text()"/>
            </pre>
        </div>
    </xsl:template>

    <!-- Listing blocks -->
    <xsl:template match="ad:listingblock">
        <div class="asciidoc-listingblock">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@title">
                <div class="asciidoc-listingblock-title"><xsl:value-of select="@title"/></div>
            </xsl:if>
            <pre class="asciidoc-listing">
                <xsl:if test="@language">
                    <xsl:attribute name="data-language"><xsl:value-of select="@language"/></xsl:attribute>
                </xsl:if>
                <code>
                    <xsl:value-of select="text()"/>
                </code>
            </pre>
        </div>
    </xsl:template>

    <!-- Example blocks -->
    <xsl:template match="ad:example">
        <div class="asciidoc-example">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@role">
                <xsl:attribute name="class">asciidoc-example <xsl:value-of select="@role"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="ad:title">
                <div class="asciidoc-example-title"><xsl:value-of select="ad:title"/></div>
            </xsl:if>
            <div class="asciidoc-example-content">
                <xsl:apply-templates select="ad:content"/>
            </div>
        </div>
    </xsl:template>

    <!-- Sidebar blocks -->
    <xsl:template match="ad:sidebar">
        <aside class="asciidoc-sidebar">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@role">
                <xsl:attribute name="class">asciidoc-sidebar <xsl:value-of select="@role"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="ad:title">
                <div class="asciidoc-sidebar-title"><xsl:value-of select="ad:title"/></div>
            </xsl:if>
            <div class="asciidoc-sidebar-content">
                <xsl:apply-templates select="ad:content"/>
            </div>
        </aside>
    </xsl:template>

    <!-- Quote blocks -->
    <xsl:template match="ad:quote">
        <blockquote class="asciidoc-quote">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@role">
                <xsl:attribute name="class">asciidoc-quote <xsl:value-of select="@role"/></xsl:attribute>
            </xsl:if>
            <div class="asciidoc-quote-content">
                <xsl:apply-templates select="ad:content"/>
            </div>
            <xsl:if test="ad:attribution or ad:citation">
                <footer class="asciidoc-quote-footer">
                    <xsl:if test="ad:attribution">
                        <cite class="asciidoc-quote-attribution">
                            <xsl:apply-templates select="ad:attribution"/>
                        </cite>
                    </xsl:if>
                    <xsl:if test="ad:citation">
                        <span class="asciidoc-quote-citation">— <xsl:value-of select="ad:citation"/></span>
                    </xsl:if>
                </footer>
            </xsl:if>
        </blockquote>
    </xsl:template>

    <!-- Verse blocks -->
    <xsl:template match="ad:verse">
        <div class="asciidoc-verse">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <pre class="asciidoc-verse-content">
                <xsl:apply-templates select="ad:content"/>
            </pre>
            <xsl:if test="ad:attribution or ad:citation">
                <div class="asciidoc-verse-footer">
                    <xsl:if test="ad:attribution">
                        <cite class="asciidoc-verse-attribution">
                            <xsl:apply-templates select="ad:attribution"/>
                        </cite>
                    </xsl:if>
                    <xsl:if test="ad:citation">
                        <span class="asciidoc-verse-citation">— <xsl:value-of select="ad:citation"/></span>
                    </xsl:if>
                </div>
            </xsl:if>
        </div>
    </xsl:template>

    <!-- Tables -->
    <xsl:template match="ad:table">
        <div class="asciidoc-table-wrapper">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="ad:title">
                <div class="asciidoc-table-title"><xsl:value-of select="ad:title"/></div>
            </xsl:if>
            <table class="asciidoc-table">
                <xsl:if test="@frame">
                    <xsl:attribute name="class">asciidoc-table asciidoc-table-frame-<xsl:value-of select="@frame"/></xsl:attribute>
                </xsl:if>
                <xsl:if test="@grid">
                    <xsl:attribute name="class">asciidoc-table asciidoc-table-grid-<xsl:value-of select="@grid"/></xsl:attribute>
                </xsl:if>
                <xsl:if test="@stripes">
                    <xsl:attribute name="class">asciidoc-table asciidoc-table-stripes-<xsl:value-of select="@stripes"/></xsl:attribute>
                </xsl:if>
                <xsl:if test="@role">
                    <xsl:attribute name="class">asciidoc-table <xsl:value-of select="@role"/></xsl:attribute>
                </xsl:if>
                <xsl:if test="ad:header">
                    <thead class="asciidoc-table-header">
                        <xsl:apply-templates select="ad:header"/>
                    </thead>
                </xsl:if>
                <tbody class="asciidoc-table-body">
                    <xsl:apply-templates select="ad:row"/>
                </tbody>
            </table>
        </div>
    </xsl:template>

    <xsl:template match="ad:table/ad:header/ad:row | ad:row">
        <tr class="asciidoc-table-row">
            <xsl:apply-templates select="ad:cell"/>
        </tr>
    </xsl:template>

    <xsl:template match="ad:cell">
        <td class="asciidoc-table-cell">
            <xsl:if test="@colspan">
                <xsl:attribute name="colspan"><xsl:value-of select="@colspan"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@rowspan">
                <xsl:attribute name="rowspan"><xsl:value-of select="@rowspan"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@align">
                <xsl:attribute name="style">text-align: <xsl:value-of select="@align"/>;</xsl:attribute>
            </xsl:if>
            <xsl:if test="@valign">
                <xsl:attribute name="style">vertical-align: <xsl:value-of select="@valign"/>;</xsl:attribute>
            </xsl:if>
            <xsl:if test="@style">
                <xsl:attribute name="class">asciidoc-table-cell <xsl:value-of select="@style"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@role">
                <xsl:attribute name="class">asciidoc-table-cell <xsl:value-of select="@role"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates/>
        </td>
    </xsl:template>

    <!-- Lists -->
    <xsl:template match="ad:list">
        <xsl:choose>
            <xsl:when test="@style = 'unordered'">
                <ul class="asciidoc-list asciidoc-list-unordered">
                    <xsl:if test="@id">
                        <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
                    </xsl:if>
                    <xsl:if test="@role">
                        <xsl:attribute name="class">asciidoc-list asciidoc-list-unordered <xsl:value-of select="@role"/></xsl:attribute>
                    </xsl:if>
                    <xsl:apply-templates select="ad:item"/>
                </ul>
            </xsl:when>
            <xsl:when test="@style = 'ordered'">
                <ol class="asciidoc-list asciidoc-list-ordered">
                    <xsl:if test="@id">
                        <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
                    </xsl:if>
                    <xsl:if test="@start">
                        <xsl:attribute name="start"><xsl:value-of select="@start"/></xsl:attribute>
                    </xsl:if>
                    <xsl:if test="@role">
                        <xsl:attribute name="class">asciidoc-list asciidoc-list-ordered <xsl:value-of select="@role"/></xsl:attribute>
                    </xsl:if>
                    <xsl:apply-templates select="ad:item"/>
                </ol>
            </xsl:when>
            <xsl:when test="@style = 'labeled'">
                <dl class="asciidoc-list asciidoc-list-labeled">
                    <xsl:if test="@id">
                        <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
                    </xsl:if>
                    <xsl:if test="@role">
                        <xsl:attribute name="class">asciidoc-list asciidoc-list-labeled <xsl:value-of select="@role"/></xsl:attribute>
                    </xsl:if>
                    <xsl:apply-templates select="ad:item"/>
                </dl>
            </xsl:when>
            <xsl:when test="@style = 'callout'">
                <ol class="asciidoc-list asciidoc-list-callout">
                    <xsl:if test="@id">
                        <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
                    </xsl:if>
                    <xsl:apply-templates select="ad:item"/>
                </ol>
            </xsl:when>
        </xsl:choose>
    </xsl:template>

    <xsl:template match="ad:list[@style='labeled']/ad:item">
        <xsl:if test="ad:term">
            <dt class="asciidoc-list-term">
                <xsl:apply-templates select="ad:term"/>
            </dt>
        </xsl:if>
        <dd class="asciidoc-list-description">
            <xsl:apply-templates select="*[not(self::ad:term)]"/>
        </dd>
    </xsl:template>

    <xsl:template match="ad:list[@style!='labeled']/ad:item">
        <li class="asciidoc-list-item">
            <xsl:if test="@marker">
                <xsl:attribute name="data-marker"><xsl:value-of select="@marker"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates/>
        </li>
    </xsl:template>

    <!-- Images -->
    <xsl:template match="ad:image">
        <figure class="asciidoc-image">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:choose>
                <xsl:when test="@link">
                    <a href="{@link}" class="asciidoc-image-link">
                        <img src="{@src}" 
                             alt="{@alt}" 
                             title="{@title}">
                            <xsl:if test="@width">
                                <xsl:attribute name="width"><xsl:value-of select="@width"/></xsl:attribute>
                            </xsl:if>
                            <xsl:if test="@height">
                                <xsl:attribute name="height"><xsl:value-of select="@height"/></xsl:attribute>
                            </xsl:if>
                        </img>
                    </a>
                </xsl:when>
                <xsl:otherwise>
                    <img src="{@src}" 
                         alt="{@alt}" 
                         title="{@title}">
                        <xsl:if test="@width">
                            <xsl:attribute name="width"><xsl:value-of select="@width"/></xsl:attribute>
                        </xsl:if>
                        <xsl:if test="@height">
                            <xsl:attribute name="height"><xsl:value-of select="@height"/></xsl:attribute>
                        </xsl:if>
                    </img>
                </xsl:otherwise>
            </xsl:choose>
            <xsl:if test="@title">
                <figcaption class="asciidoc-image-caption"><xsl:value-of select="@title"/></figcaption>
            </xsl:if>
        </figure>
    </xsl:template>

    <!-- Video -->
    <xsl:template match="ad:video">
        <div class="asciidoc-video">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <video src="{@src}">
                <xsl:if test="@poster">
                    <xsl:attribute name="poster"><xsl:value-of select="@poster"/></xsl:attribute>
                </xsl:if>
                <xsl:if test="@width">
                    <xsl:attribute name="width"><xsl:value-of select="@width"/></xsl:attribute>
                </xsl:if>
                <xsl:if test="@height">
                    <xsl:attribute name="height"><xsl:value-of select="@height"/></xsl:attribute>
                </xsl:if>
                <xsl:if test="@autoplay = 'true'">
                    <xsl:attribute name="autoplay">autoplay</xsl:attribute>
                </xsl:if>
                <xsl:if test="@loop = 'true'">
                    <xsl:attribute name="loop">loop</xsl:attribute>
                </xsl:if>
                <xsl:if test="@controls = 'false'">
                    <xsl:attribute name="controls">false</xsl:attribute>
                </xsl:if>
                <xsl:if test="not(@controls) or @controls = 'true'">
                    <xsl:attribute name="controls">controls</xsl:attribute>
                </xsl:if>
            </video>
        </div>
    </xsl:template>

    <!-- Audio -->
    <xsl:template match="ad:audio">
        <div class="asciidoc-audio">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <audio src="{@src}">
                <xsl:if test="@autoplay = 'true'">
                    <xsl:attribute name="autoplay">autoplay</xsl:attribute>
                </xsl:if>
                <xsl:if test="@loop = 'true'">
                    <xsl:attribute name="loop">loop</xsl:attribute>
                </xsl:if>
                <xsl:if test="@controls = 'false'">
                    <xsl:attribute name="controls">false</xsl:attribute>
                </xsl:if>
                <xsl:if test="not(@controls) or @controls = 'true'">
                    <xsl:attribute name="controls">controls</xsl:attribute>
                </xsl:if>
            </audio>
        </div>
    </xsl:template>

    <!-- Page break -->
    <xsl:template match="ad:pagebreak">
        <div class="asciidoc-pagebreak"></div>
    </xsl:template>

    <!-- Thematic break -->
    <xsl:template match="ad:thematicbreak">
        <hr class="asciidoc-thematic-break"/>
    </xsl:template>

    <!-- Admonitions -->
    <xsl:template match="ad:admonition">
        <div class="asciidoc-admonition asciidoc-admonition-{@type}">
            <xsl:if test="@id">
                <xsl:attribute name="id"><xsl:value-of select="@id"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@role">
                <xsl:attribute name="class">asciidoc-admonition asciidoc-admonition-{@type} <xsl:value-of select="@role"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="ad:title">
                <div class="asciidoc-admonition-title"><xsl:value-of select="ad:title"/></div>
            </xsl:if>
            <div class="asciidoc-admonition-content">
                <xsl:apply-templates select="ad:content"/>
            </div>
        </div>
    </xsl:template>

    <!-- Passthrough -->
    <xsl:template match="ad:passthrough">
        <div class="asciidoc-passthrough">
            <xsl:value-of select="text()" disable-output-escaping="yes"/>
        </div>
    </xsl:template>

    <!-- Inline content -->
    <xsl:template match="ad:title | ad:term | ad:attribution">
        <xsl:apply-templates/>
    </xsl:template>

    <!-- Inline elements -->
    <xsl:template match="ad:text">
        <xsl:value-of select="text()"/>
    </xsl:template>

    <xsl:template match="ad:strong">
        <strong class="asciidoc-strong">
            <xsl:apply-templates/>
        </strong>
    </xsl:template>

    <xsl:template match="ad:emphasis">
        <em class="asciidoc-emphasis">
            <xsl:apply-templates/>
        </em>
    </xsl:template>

    <xsl:template match="ad:monospace">
        <code class="asciidoc-monospace">
            <xsl:apply-templates/>
        </code>
    </xsl:template>

    <xsl:template match="ad:superscript">
        <sup class="asciidoc-superscript">
            <xsl:apply-templates/>
        </sup>
    </xsl:template>

    <xsl:template match="ad:subscript">
        <sub class="asciidoc-subscript">
            <xsl:apply-templates/>
        </sub>
    </xsl:template>

    <xsl:template match="ad:mark">
        <mark class="asciidoc-mark">
            <xsl:apply-templates/>
        </mark>
    </xsl:template>

    <xsl:template match="ad:link">
        <a href="{@href}" class="asciidoc-link">
            <xsl:if test="@title">
                <xsl:attribute name="title"><xsl:value-of select="@title"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@role">
                <xsl:attribute name="class">asciidoc-link <xsl:value-of select="@role"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@window">
                <xsl:attribute name="target"><xsl:value-of select="@window"/></xsl:attribute>
            </xsl:if>
            <xsl:apply-templates/>
        </a>
    </xsl:template>

    <xsl:template match="ad:xref">
        <a href="#{@refid}" class="asciidoc-xref">
            <xsl:if test="@path">
                <xsl:attribute name="href"><xsl:value-of select="@path"/>#<xsl:value-of select="@refid"/></xsl:attribute>
            </xsl:if>
            <xsl:choose>
                <xsl:when test="ad:*">
                    <xsl:apply-templates/>
                </xsl:when>
                <xsl:otherwise>
                    <xsl:value-of select="@refid"/>
                </xsl:otherwise>
            </xsl:choose>
        </a>
    </xsl:template>

    <xsl:template match="ad:image[parent::ad:*[local-name()='title' or local-name()='term' or local-name()='attribution' or local-name()='link' or local-name()='xref']]">
        <img src="{@src}" 
             alt="{@alt}" 
             title="{@title}"
             class="asciidoc-inline-image">
            <xsl:if test="@width">
                <xsl:attribute name="width"><xsl:value-of select="@width"/></xsl:attribute>
            </xsl:if>
            <xsl:if test="@height">
                <xsl:attribute name="height"><xsl:value-of select="@height"/></xsl:attribute>
            </xsl:if>
        </img>
    </xsl:template>

    <xsl:template match="ad:kbd">
        <kbd class="asciidoc-kbd">
            <xsl:apply-templates/>
        </kbd>
    </xsl:template>

    <xsl:template match="ad:button">
        <button class="asciidoc-button">
            <xsl:apply-templates/>
        </button>
    </xsl:template>

    <xsl:template match="ad:menu">
        <span class="asciidoc-menu">
            <xsl:for-each select="ad:menuitem">
                <span class="asciidoc-menuitem">
                    <xsl:value-of select="."/>
                </span>
                <xsl:if test="position() != last()">
                    <span class="asciidoc-menu-separator">→</span>
                </xsl:if>
            </xsl:for-each>
        </span>
    </xsl:template>

    <xsl:template match="ad:attribute">
        <span class="asciidoc-attribute-ref" data-attribute="{@name}">
            <xsl:text>{</xsl:text><xsl:value-of select="@name"/><xsl:text>}</xsl:text>
        </span>
    </xsl:template>

    <!-- Default template for unknown elements -->
    <xsl:template match="*">
        <xsl:apply-templates/>
    </xsl:template>

    <!-- Text nodes -->
    <xsl:template match="text()">
        <xsl:value-of select="."/>
    </xsl:template>

</xsl:stylesheet>

