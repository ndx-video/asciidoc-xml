const asciidocFrame = document.getElementById('asciidoc-frame');
const xmlFrame = document.getElementById('xml-frame');
const xsltFrame = document.getElementById('xslt-frame');
const htmlFrame = document.getElementById('html-frame');
const statusEl = document.getElementById('status');

let currentAsciiDoc = '';
let currentXML = '';
let currentXSLT = '';
let currentHTML = '';
let startupAutoConvert = false; // Flag to track if we should auto-convert on startup

// Update iframe content
function updateFrameContent(frame, content, mimeType = 'text/html', useSourceView = false, sourceType = 'html') {
    // Set up auto-convert listener for asciidoc frame on startup
    if (frame === asciidocFrame && startupAutoConvert) {
        const handleLoad = () => {
            // Only auto-convert once at startup
            if (startupAutoConvert && currentAsciiDoc) {
                startupAutoConvert = false; // Reset flag
                setTimeout(() => {
                    convertAsciiDoc();
                }, 100);
            }
            frame.removeEventListener('load', handleLoad);
        };
        frame.addEventListener('load', handleLoad, { once: true });
    }
    
    if (useSourceView) {
        const html = createSourceView(content, sourceType);
        const blob = new Blob([html], { type: 'text/html' });
        frame.src = URL.createObjectURL(blob);
    } else {
        const blob = new Blob([content], { type: mimeType });
        frame.src = URL.createObjectURL(blob);
    }
}

// Show status message
function showStatus(message, type = '') {
    statusEl.textContent = message;
    statusEl.className = 'status ' + type;
    setTimeout(() => {
        statusEl.textContent = '';
        statusEl.className = 'status';
    }, 3000);
}

// escapeHtml is now in pretty.js

// Get content from AsciiDoc (read-only, so just return stored content)
function getAsciiDocContent() {
    return currentAsciiDoc;
}

// Initialize AsciiDoc display (read-only with syntax highlighting)
function initAsciiDocEditor(content = '') {
    currentAsciiDoc = content;
    // Display as read-only with syntax highlighting, same as other columns
    updateFrameContent(asciidocFrame, content, 'text/plain', true, 'asciidoc');
}

// Load XSLT template
async function loadXSLT(path = null) {
    try {
        let response;
        if (path) {
            response = await fetch(`/api/load-file?path=${encodeURIComponent(path)}`);
        } else {
            response = await fetch('/api/xslt');
            if (response.ok) {
                document.getElementById('xslt-path').value = '/xslt/asciidoc-to-html.xsl';
            }
        }
        if (!response.ok) throw new Error('Failed to load XSLT');
        currentXSLT = await response.text();
        updateFrameContent(xsltFrame, currentXSLT, 'application/xml', true, 'xslt');
        if (path) {
            document.getElementById('xslt-path').value = path;
        }
    } catch (error) {
        showStatus('Failed to load XSLT: ' + error.message, 'error');
    }
}

// Load AsciiDoc from server path
async function loadAsciiDocFromPath() {
    const path = document.getElementById('asciidoc-path').value.trim();
    if (!path) {
        showStatus('Please enter a path', 'error');
        return;
    }

    try {
        showStatus('Loading AsciiDoc...');
        const response = await fetch(`/api/load-file?path=${encodeURIComponent(path)}`);
        if (!response.ok) {
            throw new Error(`Failed to load file: ${response.statusText}`);
        }
        const content = await response.text();
        currentAsciiDoc = content;
        
        // Update AsciiDoc display with syntax highlighting
        initAsciiDocEditor(content);
        
        // Auto-convert (but not on startup - that's handled by the iframe load event)
        if (!startupAutoConvert) {
            await convertAsciiDoc();
        }
        showStatus('AsciiDoc loaded', 'success');
    } catch (error) {
        showStatus('Failed to load AsciiDoc: ' + error.message, 'error');
    }
}

// Load XSLT from server path
async function loadXSLTFromPath() {
    const path = document.getElementById('xslt-path').value.trim();
    if (!path) {
        showStatus('Please enter a path', 'error');
        return;
    }

    try {
        showStatus('Loading XSLT...');
        await loadXSLT(path);
        if (currentXML) {
            transformXMLToHTML();
        }
        showStatus('XSLT loaded', 'success');
    } catch (error) {
        showStatus('Failed to load XSLT: ' + error.message, 'error');
    }
}

// Validate AsciiDoc
async function validateAsciiDoc() {
    const asciidoc = getAsciiDocContent();
    if (!asciidoc.trim()) {
        showStatus('No AsciiDoc content to validate', 'error');
        return;
    }

    try {
        showStatus('Validating...');
        const response = await fetch('/api/validate', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ asciidoc })
        });

        const result = await response.json();
        if (result.valid) {
            showStatus('✓ Valid AsciiDoc', 'success');
        } else {
            showStatus('✗ Invalid: ' + result.error, 'error');
        }
    } catch (error) {
        showStatus('Validation error: ' + error.message, 'error');
    }
}

// Convert AsciiDoc to XML
async function convertAsciiDoc() {
    const asciidoc = getAsciiDocContent();
    if (!asciidoc.trim()) {
        showStatus('No AsciiDoc content to convert', 'error');
        return;
    }

    try {
        showStatus('Converting...');
        const response = await fetch('/api/convert', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ asciidoc })
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error);
        }

        const result = await response.json();
        currentXML = result.xml;
        currentAsciiDoc = asciidoc;

        // Update XML frame with syntax highlighting
        updateFrameContent(xmlFrame, currentXML, 'application/xml', true, 'xml');

        // Transform XML to HTML
        transformXMLToHTML();
    } catch (error) {
        showStatus('Conversion error: ' + error.message, 'error');
    }
}

// Transform XML to HTML using browser XSLT
function transformXMLToHTML() {
    if (!currentXML || !currentXSLT) return;

    try {
        const parser = new DOMParser();
        const xmlDoc = parser.parseFromString(currentXML, 'application/xml');
        
        const parserError = xmlDoc.querySelector('parsererror');
        if (parserError) {
            throw new Error('XML parsing error: ' + parserError.textContent);
        }

        const xsltDoc = parser.parseFromString(currentXSLT, 'application/xml');
        const processor = new XSLTProcessor();
        processor.importStylesheet(xsltDoc);
        
        const resultDoc = processor.transformToDocument(xmlDoc);
        const serializer = new XMLSerializer();
        const html = serializer.serializeToString(resultDoc);
        currentHTML = html;

        updateHTMLOutput();
    } catch (error) {
        showStatus('XSLT transformation error: ' + error.message, 'error');
        updateFrameContent(htmlFrame, '<pre>' + escapeHtml(error.message) + '</pre>');
    }
}

// Update HTML output frame
function updateHTMLOutput() {
    const view = document.querySelector('.html-tabs button.active').dataset.view;
    
    if (view === 'rendered') {
        updateFrameContent(htmlFrame, currentHTML, 'text/html', false);
    } else {
        updateFrameContent(htmlFrame, currentHTML, 'text/html', true, 'html');
    }
}

// Load example file
async function loadExample() {
    const path = '/static/comprehensive.adoc';
    document.getElementById('asciidoc-path').value = path;
    await loadAsciiDocFromPath();
}

// Upload file
async function uploadFile(type) {
    const fileInput = type === 'asciidoc' ? document.getElementById('asciidocFile') : document.getElementById('xsltFile');
    const file = fileInput.files[0];
    if (!file) {
        showStatus('Please select a file', 'error');
        return;
    }

    const formData = new FormData();
    formData.append('file', file);
    formData.append('type', type);

    try {
        showStatus(`Uploading ${type}...`);
        const response = await fetch('/api/upload', {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error);
        }

        const result = await response.json();
        showStatus(`${type} uploaded successfully`, 'success');

        if (type === 'asciidoc') {
            document.getElementById('asciidoc-path').value = result.path;
            await loadAsciiDocFromPath();
        } else {
            document.getElementById('xslt-path').value = result.path;
            await loadXSLTFromPath();
        }

        fileInput.value = '';
        document.getElementById(`btn-upload-${type}`).disabled = true;
    } catch (error) {
        showStatus(`Upload failed: ${error.message}`, 'error');
    }
}

// Resizable columns
function initResizableColumns() {
    const panels = document.querySelectorAll('.panel');
    panels.forEach((panel, index) => {
        if (index < panels.length - 1) {
            const resizer = document.createElement('div');
            resizer.className = 'resizer';
            panel.appendChild(resizer);
            
            let isResizing = false;
            let startX = 0;
            let startWidth = 0;
            let nextStartWidth = 0;
            
            resizer.addEventListener('mousedown', (e) => {
                isResizing = true;
                startX = e.clientX;
                startWidth = panel.offsetWidth;
                const nextPanel = panels[index + 1];
                nextStartWidth = nextPanel.offsetWidth;
                resizer.classList.add('active');
                document.body.style.cursor = 'col-resize';
                document.body.style.userSelect = 'none';
                
                // Disable pointer events on iframes during resize
                document.querySelectorAll('iframe').forEach(iframe => {
                    iframe.style.pointerEvents = 'none';
                });
                
                e.preventDefault();
                e.stopPropagation();
            });
            
            document.addEventListener('mousemove', (e) => {
                if (!isResizing) return;
                
                const diff = e.clientX - startX;
                const newWidth = startWidth + diff;
                const nextPanel = panels[index + 1];
                const newNextWidth = nextStartWidth - diff;
                
                if (newWidth >= 150 && newNextWidth >= 150) {
                    panel.style.flex = `0 0 ${newWidth}px`;
                    nextPanel.style.flex = `0 0 ${newNextWidth}px`;
                }
            });
            
            document.addEventListener('mouseup', () => {
                if (isResizing) {
                    isResizing = false;
                    resizer.classList.remove('active');
                    document.body.style.cursor = '';
                    document.body.style.userSelect = '';
                    
                    // Re-enable pointer events on iframes after resize
                    document.querySelectorAll('iframe').forEach(iframe => {
                        iframe.style.pointerEvents = 'auto';
                    });
                }
            });
        }
    });
}

// Event listeners
document.getElementById('btn-validate').addEventListener('click', validateAsciiDoc);
document.getElementById('btn-convert').addEventListener('click', convertAsciiDoc);
document.getElementById('btn-load-example').addEventListener('click', loadExample);
document.getElementById('btn-upload').addEventListener('click', () => {
    document.getElementById('uploadModal').style.display = 'block';
});
document.getElementById('closeModal').addEventListener('click', () => {
    document.getElementById('uploadModal').style.display = 'none';
});
document.getElementById('btn-upload-asciidoc').addEventListener('click', () => uploadFile('asciidoc'));
document.getElementById('btn-upload-xslt').addEventListener('click', () => uploadFile('xslt'));

document.getElementById('asciidocFile').addEventListener('change', (e) => {
    document.getElementById('btn-upload-asciidoc').disabled = !e.target.files[0];
});
document.getElementById('xsltFile').addEventListener('change', (e) => {
    document.getElementById('btn-upload-xslt').disabled = !e.target.files[0];
});

window.addEventListener('click', (e) => {
    const modal = document.getElementById('uploadModal');
    if (e.target === modal) {
        modal.style.display = 'none';
    }
});

document.querySelectorAll('.html-tabs button').forEach(btn => {
    btn.addEventListener('click', () => {
        document.querySelectorAll('.html-tabs button').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        if (currentHTML) {
            updateHTMLOutput();
        }
    });
});

// Initialize when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initialize);
} else {
    initialize();
}

function initialize() {
    initResizableColumns();

    // Load XSLT first, then initialize AsciiDoc display, then load example
    (async () => {
        await loadXSLT();
        
        // Enable auto-convert on startup
        startupAutoConvert = true;
        
        initAsciiDocEditor();
        // Wait a bit for display to initialize, then load example which will auto-convert
        setTimeout(async () => {
            await loadExample();
        }, 300);
    })();
}

