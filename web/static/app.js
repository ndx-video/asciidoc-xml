const asciidocFrame = document.getElementById('asciidoc-frame');
const xmlFrame = document.getElementById('xml-frame');
const xsltFrame = document.getElementById('xslt-frame');
const htmlFrame = document.getElementById('html-frame');
const statusEl = document.getElementById('status');
const outputTypeSelect = document.getElementById('output-type');

let currentAsciiDoc = '';
let currentXML = '';
let currentXSLT = '';
let currentHTML = '';
let startupAutoConvert = false;

// Get current output type
function getOutputType() {
    return outputTypeSelect.value;
}

// Check if XSLT should be available for current output type
function shouldShowXSLT(outputType) {
    return outputType === 'xml' || outputType === 'xhtml' || outputType === 'xhtml5';
}

// Update column visibility based on output type
function updateColumnVisibility() {
    const outputType = getOutputType();
    const xmlPanel = document.getElementById('xml-panel');
    const xsltPanel = document.getElementById('xslt-panel');
    const htmlPanel = document.getElementById('html-panel');
    const xsltUploadSection = document.getElementById('xslt-upload-section');

    // Show/hide XML panel
    if (outputType === 'xml') {
        xmlPanel.classList.remove('hidden');
    } else {
        xmlPanel.classList.add('hidden');
    }

    // Show/hide XSLT panel and upload section
    const showXSLT = shouldShowXSLT(outputType);
    if (showXSLT) {
        xsltPanel.classList.remove('hidden');
        xsltUploadSection.style.display = 'block';
    } else {
        xsltPanel.classList.add('hidden');
        xsltUploadSection.style.display = 'none';
    }

    // HTML panel is always visible
    htmlPanel.classList.remove('hidden');

    // Update resizable columns after visibility changes
    setTimeout(() => {
        initResizableColumns();
    }, 0);
}

// Update panel header based on output type
function updatePanelHeaders() {
    const outputType = getOutputType();
    const htmlPanel = document.getElementById('html-panel');
    const htmlPanelHeader = htmlPanel.querySelector('.panel-header');
    
    let headerText = 'HTML Output';
    if (outputType === 'xhtml' || outputType === 'xhtml5') {
        headerText = 'XHTML Output';
    } else if (outputType === 'html5') {
        headerText = 'HTML5 Output';
    }
    
    htmlPanelHeader.textContent = headerText;
}

// Update iframe content
function updateFrameContent(frame, content, mimeType = 'text/html', useSourceView = false, sourceType = 'html') {
    if (frame === asciidocFrame && startupAutoConvert) {
        const handleLoad = () => {
            if (startupAutoConvert && currentAsciiDoc) {
                startupAutoConvert = false;
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

// Get content from AsciiDoc
function getAsciiDocContent() {
    return currentAsciiDoc;
}

// Initialize AsciiDoc display
function initAsciiDocEditor(content = '') {
    currentAsciiDoc = content;
    updateFrameContent(asciidocFrame, content, 'text/plain', true, 'asciidoc');
}

// Load XSLT template
async function loadXSLT(path = null) {
    if (!shouldShowXSLT(getOutputType())) {
        return;
    }

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
        
        initAsciiDocEditor(content);
        
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
    if (!shouldShowXSLT(getOutputType())) {
        showStatus('XSLT is not available for the selected output type', 'error');
        return;
    }

    const path = document.getElementById('xslt-path').value.trim();
    if (!path) {
        showStatus('Please enter a path', 'error');
        return;
    }

    try {
        showStatus('Loading XSLT...');
        await loadXSLT(path);
        if (currentXML && shouldShowXSLT(getOutputType())) {
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

// Convert AsciiDoc based on selected output type
async function convertAsciiDoc() {
    const asciidoc = getAsciiDocContent();
    if (!asciidoc.trim()) {
        showStatus('No AsciiDoc content to convert', 'error');
        return;
    }

    const outputType = getOutputType();
    
    try {
        showStatus('Converting...');
        const response = await fetch('/api/convert', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                asciidoc: asciidoc,
                output: outputType
            })
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error);
        }

        const result = await response.json();
        currentAsciiDoc = asciidoc;

        // Handle XML output
        if (outputType === 'xml') {
            currentXML = result.output;
            updateFrameContent(xmlFrame, currentXML, 'application/xml', true, 'xml');
            
            // If XSLT is loaded, transform to HTML
            if (currentXSLT && shouldShowXSLT(outputType)) {
                transformXMLToHTML();
            } else {
                // Clear HTML output if no XSLT
                currentHTML = '';
                updateHTMLOutput();
            }
        } else if (outputType === 'md2adoc') {
            // Handle MD2ADoc output (Markdown to AsciiDoc conversion)
            // Show the converted AsciiDoc in the AsciiDoc frame
            currentAsciiDoc = result.output;
            updateFrameContent(asciidocFrame, currentAsciiDoc, 'text/plain', false, 'asciidoc');
            
            // Clear other outputs
            currentXML = '';
            currentHTML = '';
            updateHTMLOutput();
        } else {
            // Handle HTML/XHTML output (direct conversion)
            currentHTML = result.output;
            
            // Clear XML if not showing XML panel
            if (outputType !== 'xml') {
                currentXML = '';
            }
            
            updateHTMLOutput();
        }

        showStatus('Conversion complete', 'success');
    } catch (error) {
        showStatus('Conversion error: ' + error.message, 'error');
    }
}

// Transform XML to HTML using browser XSLT
function transformXMLToHTML() {
    if (!currentXML || !currentXSLT || !shouldShowXSLT(getOutputType())) return;

    try {
        const parser = new DOMParser();
        const xmlDoc = parser.parseFromString(currentXML, 'application/xml');
        
        const parserError = xmlDoc.querySelector('parsererror');
        if (parserError) {
            const errorMessage = 'XML parsing error: ' + parserError.textContent;
            throw new Error(errorMessage);
        }

        const xsltDoc = parser.parseFromString(currentXSLT, 'application/xml');
        const xsltParserError = xsltDoc.querySelector('parsererror');
        if (xsltParserError) {
            const errorMessage = 'XSLT parsing error: ' + xsltParserError.textContent;
            throw new Error(errorMessage);
        }

        const processor = new XSLTProcessor();
        processor.importStylesheet(xsltDoc);
        
        const resultDoc = processor.transformToDocument(xmlDoc);
        const serializer = new XMLSerializer();
        const html = serializer.serializeToString(resultDoc);
        currentHTML = html;

        updateHTMLOutput();
    } catch (error) {
        console.error('XSLT transformation error:', error);
        showStatus('XSLT transformation error: ' + error.message, 'error');
        updateFrameContent(htmlFrame, '<pre>' + escapeHtml(error.message) + '</pre>', 'text/html');
        
        // Report to server for logging
        if (window.reportErrorToServer) {
            window.reportErrorToServer({
                message: 'XSLT/XML transformation error: ' + error.message,
                stack: error.stack || '',
                url: window.location.href,
                lineNumber: 0,
                colNumber: 0,
                userAgent: navigator.userAgent,
                timestamp: new Date().toISOString(),
                location: window.location.href
            });
        }
    }
}

// Update HTML output frame
function updateHTMLOutput() {
    const view = document.querySelector('.html-tabs button.active').dataset.view;
    const outputType = getOutputType();
    
    if (view === 'rendered' && currentHTML) {
        updateFrameContent(htmlFrame, currentHTML, 'text/html', false);
    } else if (currentHTML) {
        updateFrameContent(htmlFrame, currentHTML, 'text/html', true, 'html');
    }
}

// Load example file
async function loadExample() {
    const path = 'examples/comprehensive.adoc';
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

    if (type === 'xslt' && !shouldShowXSLT(getOutputType())) {
        showStatus('XSLT is not available for the selected output type', 'error');
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
    // Remove existing resizers
    document.querySelectorAll('.resizer').forEach(resizer => resizer.remove());

    const panels = Array.from(document.querySelectorAll('.panel:not(.hidden)'));
    
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
                    
                    document.querySelectorAll('iframe').forEach(iframe => {
                        iframe.style.pointerEvents = 'auto';
                    });
                }
            });
        }
    });
}

// Event listeners - wrapped in try/catch for error reporting
try {
    document.getElementById('btn-validate').addEventListener('click', function() {
        try {
            validateAsciiDoc();
        } catch (error) {
            console.error('Error in validate button handler:', error);
            showStatus('Validation error: ' + error.message, 'error');
        }
    });
    
    document.getElementById('btn-convert').addEventListener('click', function() {
        try {
            convertAsciiDoc();
        } catch (error) {
            console.error('Error in convert button handler:', error);
            showStatus('Conversion error: ' + error.message, 'error');
        }
    });
    
    document.getElementById('btn-load-example').addEventListener('click', function() {
        try {
            loadExample();
        } catch (error) {
            console.error('Error in load example handler:', error);
            showStatus('Load error: ' + error.message, 'error');
        }
    });
    
    document.getElementById('btn-upload').addEventListener('click', function() {
        try {
            document.getElementById('uploadModal').style.display = 'block';
        } catch (error) {
            console.error('Error showing upload modal:', error);
        }
    });
    
    document.getElementById('closeModal').addEventListener('click', function() {
        try {
            document.getElementById('uploadModal').style.display = 'none';
        } catch (error) {
            console.error('Error closing modal:', error);
        }
    });
    
    document.getElementById('btn-upload-asciidoc').addEventListener('click', function() {
        try {
            uploadFile('asciidoc');
        } catch (error) {
            console.error('Error in upload asciidoc handler:', error);
            showStatus('Upload error: ' + error.message, 'error');
        }
    });
    
    document.getElementById('btn-upload-xslt').addEventListener('click', function() {
        try {
            uploadFile('xslt');
        } catch (error) {
            console.error('Error in upload xslt handler:', error);
            showStatus('Upload error: ' + error.message, 'error');
        }
    });

    document.getElementById('asciidocFile').addEventListener('change', function(e) {
        try {
            document.getElementById('btn-upload-asciidoc').disabled = !e.target.files[0];
        } catch (error) {
            console.error('Error in asciidoc file change handler:', error);
        }
    });
    
    document.getElementById('xsltFile').addEventListener('change', function(e) {
        try {
            document.getElementById('btn-upload-xslt').disabled = !e.target.files[0];
        } catch (error) {
            console.error('Error in xslt file change handler:', error);
        }
    });

    window.addEventListener('click', function(e) {
        try {
            const modal = document.getElementById('uploadModal');
            if (e.target === modal) {
                modal.style.display = 'none';
            }
        } catch (error) {
            console.error('Error in window click handler:', error);
        }
    });

    document.querySelectorAll('.html-tabs button').forEach(function(btn) {
        btn.addEventListener('click', function() {
            try {
                document.querySelectorAll('.html-tabs button').forEach(function(b) {
                    b.classList.remove('active');
                });
                btn.classList.add('active');
                if (currentHTML) {
                    updateHTMLOutput();
                }
            } catch (error) {
                console.error('Error in HTML tabs handler:', error);
            }
        });
    });

    // Output type change handler
    outputTypeSelect.addEventListener('change', function() {
        try {
            updateColumnVisibility();
            updatePanelHeaders();
            
            // If we have AsciiDoc content, reconvert with new output type
            if (currentAsciiDoc) {
                convertAsciiDoc();
            }
            
            // Load XSLT if needed for new output type
            if (shouldShowXSLT(getOutputType()) && !currentXSLT) {
                loadXSLT();
            }
        } catch (error) {
            console.error('Error in output type change handler:', error);
            showStatus('Error changing output type: ' + error.message, 'error');
        }
    });
} catch (error) {
    console.error('Error setting up event listeners:', error);
    if (window.reportErrorToServer) {
        window.reportErrorToServer({
            message: 'Error setting up event listeners: ' + error.message,
            stack: error.stack || '',
            url: window.location.href,
            lineNumber: 0,
            colNumber: 0,
            userAgent: navigator.userAgent,
            timestamp: new Date().toISOString(),
            location: window.location.href
        });
    }
}

// Initialize when DOM is ready
try {
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', function() {
            try {
                initialize();
            } catch (error) {
                console.error('Error during DOMContentLoaded:', error);
                if (window.reportErrorToServer) {
                    window.reportErrorToServer({
                        message: 'Error during DOMContentLoaded: ' + error.message,
                        stack: error.stack || '',
                        url: window.location.href,
                        lineNumber: 0,
                        colNumber: 0,
                        userAgent: navigator.userAgent,
                        timestamp: new Date().toISOString(),
                        location: window.location.href
                    });
                }
            }
        });
    } else {
        initialize();
    }
} catch (error) {
    console.error('Error setting up initialization:', error);
}

function initialize() {
    try {
        // Set initial column visibility
        updateColumnVisibility();
        updatePanelHeaders();
        initResizableColumns();

        // Load XSLT if needed for initial output type
        (async () => {
            try {
                if (shouldShowXSLT(getOutputType())) {
                    await loadXSLT();
                }
                
                // Enable auto-convert on startup
                startupAutoConvert = true;
                
                initAsciiDocEditor();
                // Wait a bit for display to initialize, then load example which will auto-convert
                setTimeout(async () => {
                    try {
                        await loadExample();
                    } catch (error) {
                        console.error('Error loading example:', error);
                        showStatus('Failed to load example: ' + error.message, 'error');
                    }
                }, 300);
            } catch (error) {
                console.error('Error in async initialization:', error);
                showStatus('Initialization error: ' + error.message, 'error');
            }
        })();
    } catch (error) {
        console.error('Error in initialize:', error);
        if (window.reportErrorToServer) {
            window.reportErrorToServer({
                message: 'Error in initialize: ' + error.message,
                stack: error.stack || '',
                url: window.location.href,
                lineNumber: 0,
                colNumber: 0,
                userAgent: navigator.userAgent,
                timestamp: new Date().toISOString(),
                location: window.location.href
            });
        }
    }
}
