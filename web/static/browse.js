// ./web/static/browse.js

const outputFrame = document.getElementById('output-frame');
const statusEl = document.getElementById('status');
const outputTypeSelect = document.getElementById('output-type');
const fileTreeEl = document.getElementById('file-tree');
const currentPathEl = document.getElementById('current-path');
const htmlTabs = document.getElementById('html-tabs');

let currentAsciiDoc = '';
let currentOutput = '';
let currentPath = '';
let currentType = 'html';

// Status helper
function showStatus(message, type = '') {
    statusEl.textContent = message;
    statusEl.className = 'status ' + type;
    setTimeout(() => {
        statusEl.textContent = '';
        statusEl.className = 'status';
    }, 3000);
}

// Fetch file tree
async function loadFileTree() {
    try {
        const response = await fetch('/api/files');
        if (!response.ok) throw new Error('Failed to load file list');
        const files = await response.json();
        renderFileTree(files, fileTreeEl);
    } catch (error) {
        fileTreeEl.innerHTML = `<div style="color: red; padding: 10px;">Error: ${error.message}</div>`;
    }
}

// Render file tree
function renderFileTree(nodes, container) {
    container.innerHTML = '';
    
    if (!nodes || nodes.length === 0) {
        container.innerHTML = '<div style="padding: 10px;">No files found</div>';
        return;
    }

    const ul = document.createElement('ul');
    ul.style.listStyle = 'none';
    ul.style.paddingLeft = '0';
    ul.style.margin = '0';

    nodes.forEach(node => {
        const li = document.createElement('li');
        
        const div = document.createElement('div');
        div.className = 'tree-node';
        div.dataset.path = node.path;
        
        // Icon
        const icon = document.createElement('span');
        icon.className = 'tree-icon';
        icon.textContent = node.type === 'dir' ? '📁' : '📄';
        
        // Toggle for dirs
        const toggle = document.createElement('span');
        toggle.className = 'tree-toggle';
        toggle.textContent = node.type === 'dir' ? '▶' : ''; // Right arrow for collapsed
        
        div.appendChild(toggle);
        div.appendChild(icon);
        div.appendChild(document.createTextNode(node.name));
        
        li.appendChild(div);

        if (node.type === 'dir') {
            const childrenContainer = document.createElement('div');
            childrenContainer.className = 'tree-children';
            li.appendChild(childrenContainer);

            if (node.children) {
                renderFileTree(node.children, childrenContainer);
            }

            div.addEventListener('click', (e) => {
                e.stopPropagation();
                const isExpanded = childrenContainer.classList.contains('expanded');
                if (isExpanded) {
                    childrenContainer.classList.remove('expanded');
                    toggle.textContent = '▶';
                } else {
                    childrenContainer.classList.add('expanded');
                    toggle.textContent = '▼';
                }
            });
        } else {
            div.addEventListener('click', (e) => {
                e.stopPropagation();
                // Highlight selection
                document.querySelectorAll('.tree-node').forEach(n => n.classList.remove('active'));
                div.classList.add('active');
                loadFile(node.path);
            });
        }

        ul.appendChild(li);
    });

    container.appendChild(ul);
}

// Load selected file
async function loadFile(path) {
    currentPath = path;
    currentPathEl.textContent = path;
    
    // Prefix with examples/ for the API call if it's not already there
    // Actually the API expects path relative to static or examples root
    // But handleLoadFile logic is custom.
    // If the file listing returns "examples/foo/bar.adoc", then we pass "examples/foo/bar.adoc"
    // Our Go handler supports "examples/" prefix.
    
    try {
        showStatus('Loading...');
        const response = await fetch(`/api/load-file?path=${encodeURIComponent(path)}`);
        if (!response.ok) throw new Error('Failed to load file');
        currentAsciiDoc = await response.text();
        await convertAsciiDoc();
    } catch (error) {
        showStatus('Error: ' + error.message, 'error');
        outputFrame.src = 'about:blank';
    }
}

// Convert loaded file
async function convertAsciiDoc() {
    if (!currentAsciiDoc) return;

    const outputType = outputTypeSelect.value;
    currentType = outputType;

    try {
        showStatus('Converting...');
        const response = await fetch('/api/convert', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                asciidoc: currentAsciiDoc,
                output: outputType
            })
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error);
        }

        const result = await response.json();
        currentOutput = result.output;
        
        updateOutputView();
        showStatus('Loaded', 'success');

        // Show tabs only for HTML output
        if (outputType.startsWith('html') || outputType.startsWith('xhtml')) {
            htmlTabs.classList.add('visible');
        } else if (outputType === 'md2adoc') {
            // For md2adoc, show the converted AsciiDoc in the frame
            htmlTabs.classList.remove('visible');
        } else {
            htmlTabs.classList.remove('visible');
        }

    } catch (error) {
        showStatus('Conversion error: ' + error.message, 'error');
        outputFrame.src = 'data:text/plain;charset=utf-8,' + encodeURIComponent(error.message);
    }
}

// Update iframe content
function updateOutputView() {
    const activeTab = document.querySelector('.html-tabs button.active');
    const viewMode = activeTab ? activeTab.dataset.view : 'rendered';
    let contentType = 'text/html';
    if (currentType === 'xml') {
        contentType = 'application/xml';
    } else if (currentType === 'md2adoc') {
        contentType = 'text/plain';
    }

    if (viewMode === 'source' || currentType === 'xml') {
        // Source view
        const html = createSourceView(currentOutput, currentType);
        const blob = new Blob([html], { type: 'text/html' });
        outputFrame.src = URL.createObjectURL(blob);
    } else {
        // Rendered view
        const blob = new Blob([currentOutput], { type: contentType });
        outputFrame.src = URL.createObjectURL(blob);
    }
}

// Tab handling
document.querySelectorAll('.html-tabs button').forEach(btn => {
    btn.addEventListener('click', () => {
        document.querySelectorAll('.html-tabs button').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        updateOutputView();
    });
});

// Output type change
outputTypeSelect.addEventListener('change', () => {
    if (currentAsciiDoc) {
        convertAsciiDoc();
    }
});

// Resize handling (simplified from app.js)
function initResizable() {
    const panel = document.getElementById('browser-panel');
    const resizer = document.createElement('div');
    resizer.className = 'resizer';
    panel.appendChild(resizer);
    
    let isResizing = false;
    
    resizer.addEventListener('mousedown', (e) => {
        isResizing = true;
        document.body.style.cursor = 'col-resize';
        outputFrame.style.pointerEvents = 'none';
        e.preventDefault();
    });

    document.addEventListener('mousemove', (e) => {
        if (!isResizing) return;
        const newWidth = e.clientX - panel.getBoundingClientRect().left;
        if (newWidth > 150 && newWidth < 600) {
            panel.style.flex = `0 0 ${newWidth}px`;
        }
    });

    document.addEventListener('mouseup', () => {
        if (isResizing) {
            isResizing = false;
            document.body.style.cursor = '';
            outputFrame.style.pointerEvents = 'auto';
        }
    });
}

// Init
document.addEventListener('DOMContentLoaded', () => {
    loadFileTree();
    initResizable();
});

