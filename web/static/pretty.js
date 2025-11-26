// Pretty print functionality for AsciiDoc, XML, XSLT, and HTML

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// AsciiDoc syntax highlighter
function highlightAsciiDoc(content) {
    const lines = content.split('\n');
    const highlighted = lines.map((line) => {
        let html = escapeHtml(line);
        
        // Document title (level 0)
        if (/^= [^=]/.test(line)) {
            html = html.replace(/^(= )(.+)$/, '<span class="title title-level-1">$1$2</span>');
        }
        // Section titles
        else if (/^==+ /.test(line)) {
            const level = (line.match(/^=+/)[0].length);
            html = html.replace(/^(=+ )(.+)$/, `<span class="title title-level-${level}">$1$2</span>`);
        }
        // Attributes
        else if (/^:/.test(line)) {
            html = html.replace(/^(:[^:]+:)(.*)$/, '<span class="attribute">$1</span>$2');
        }
        // Block delimiters
        else if (/^(----|\.\.\.\.|====|\*\*\*\*|____|\|\|===|```)/.test(line.trim())) {
            html = `<span class="block-delimiter">${html}</span>`;
        }
        // List markers
        else if (/^[\*\-\+\.] /.test(line) || /^[\*\-\+\.]{2,} /.test(line)) {
            html = html.replace(/^([\*\-\+\.]+ )/, '<span class="list-marker">$1</span>');
        }
        // Table delimiters
        else if (/^\|/.test(line)) {
            html = html.replace(/\|/g, '<span class="table-delimiter">|</span>');
        }
        // Inline formatting
        html = html.replace(/\*([^*]+)\*/g, '<span class="bold">*$1*</span>');
        html = html.replace(/_([^_]+)_/g, '<span class="italic">_$1_</span>');
        html = html.replace(/`([^`]+)`/g, '<span class="monospace">`$1`</span>');
        html = html.replace(/(https?:\/\/[^\s\[\]]+)(\[([^\]]+)\])?/g, '<span class="link">$1</span>');
        
        return `<span class="line">${html}</span>`;
    }).join('\n');
    
    return highlighted;
}

// XML syntax highlighter
function highlightXML(content) {
    return escapeHtml(content)
        .replace(/(&lt;\/?)([\w:]+)([^&]*?)(\/?&gt;)/g, (match, open, tag, attrs, close) => {
            let result = `<span class="tag">${open}</span><span class="tag-name">${tag}</span>`;
            // Highlight attributes
            attrs = attrs.replace(/(\w+)="([^"]*)"/g, '<span class="attribute">$1</span>="<span class="attribute-value">$2</span>"');
            result += attrs;
            result += `<span class="tag">${close}</span>`;
            return result;
        })
        .replace(/&lt;!--[\s\S]*?--&gt;/g, '<span class="comment">$&</span>')
        .replace(/&lt;!\[CDATA\[[\s\S]*?\]\]&gt;/g, '<span class="cdata">$&</span>');
}

// HTML syntax highlighter
function highlightHTML(content) {
    return highlightXML(content); // Same as XML for now
}

// Create pretty-printed source view
function createSourceView(content, type) {
    let highlighted;
    let className;
    
    if (type === 'asciidoc') {
        highlighted = highlightAsciiDoc(content);
        className = 'asciidoc-source';
    } else if (type === 'xml' || type === 'xslt') {
        highlighted = highlightXML(content);
        className = 'xml-source';
    } else if (type === 'html') {
        highlighted = highlightHTML(content);
        className = 'html-source';
    } else {
        highlighted = escapeHtml(content);
        className = 'xml-source';
    }
    
    return `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body { margin: 0; padding: 0; overflow: hidden; }
        .${className} {
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            font-size: 14px;
            line-height: 1.6;
            padding: 1rem;
            background: #fff;
            color: #333;
            white-space: pre;
            overflow: auto;
            height: 100vh;
            margin: 0;
        }
        ${type === 'asciidoc' ? `
        .asciidoc-source .title { color: #0066cc; font-weight: bold; }
        .asciidoc-source .title-level-1 { color: #0066cc; font-size: 1.2em; }
        .asciidoc-source .title-level-2 { color: #0088cc; }
        .asciidoc-source .title-level-3 { color: #00aacc; }
        .asciidoc-source .attribute { color: #990099; }
        .asciidoc-source .bold { color: #cc0000; font-weight: bold; }
        .asciidoc-source .italic { color: #006600; font-style: italic; }
        .asciidoc-source .monospace { color: #cc6600; }
        .asciidoc-source .link { color: #0066cc; text-decoration: underline; }
        .asciidoc-source .block-delimiter { color: #999999; }
        .asciidoc-source .list-marker { color: #666666; }
        .asciidoc-source .table-delimiter { color: #999999; }
        ` : ''}
        ${type === 'xml' || type === 'xslt' ? `
        .xml-source .tag { color: #0066cc; }
        .xml-source .tag-name { color: #0066cc; font-weight: bold; }
        .xml-source .attribute { color: #990099; }
        .xml-source .attribute-value { color: #cc0000; }
        .xml-source .comment { color: #999999; font-style: italic; }
        .xml-source .cdata { color: #006600; }
        ` : ''}
        ${type === 'html' ? `
        .html-source .tag { color: #0066cc; }
        .html-source .tag-name { color: #0066cc; font-weight: bold; }
        .html-source .attribute { color: #990099; }
        .html-source .attribute-value { color: #cc0000; }
        .html-source .comment { color: #999999; font-style: italic; }
        ` : ''}
    </style>
</head>
<body>
    <div class="${className}">${highlighted}</div>
</body>
</html>`;
}

