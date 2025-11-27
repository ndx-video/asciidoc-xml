// ./web/static/docs.js

(function() {
    const urlParams = new URLSearchParams(window.location.search);
    const currentPage = urlParams.get('page') || null;
    
    // Flatten nodes for horizontal strip
    function flattenNodes(nodes) {
        const flat = [];
        nodes.forEach(node => {
            if (node.type === 'file') {
                flat.push(node);
            } else if (node.type === 'dir' && node.children) {
                flat.push(...flattenNodes(node.children));
            }
        });
        return flat;
    }
    
    async function loadDocsTOC() {
        try {
            const response = await fetch('/api/docs');
            if (!response.ok) {
                throw new Error('Failed to load docs list');
            }
            const nodes = await response.json();
            
            // Create sticky ToC strip
            const tocStrip = document.getElementById('docs-toc-strip');
            if (tocStrip) {
                const flatNodes = flattenNodes(nodes);
                const ul = document.createElement('ul');
                
                flatNodes.forEach(node => {
                    const li = document.createElement('li');
                    const a = document.createElement('a');
                    
                    const pathParts = node.path.split('/');
                    const filename = pathParts[pathParts.length - 1];
                    a.href = '/docs?page=' + encodeURIComponent(filename);
                    a.textContent = node.name;
                    
                    if (currentPage && filename === currentPage) {
                        a.classList.add('active');
                    }
                    
                    li.appendChild(a);
                    ul.appendChild(li);
                });
                
                tocStrip.innerHTML = '';
                tocStrip.appendChild(ul);
            }
            
            // Also populate sidebar ToC (for mobile/fallback)
            const toc = document.getElementById('docs-toc');
            if (toc) {
                toc.innerHTML = '';
                
                function renderNodes(nodes, parent) {
                    nodes.forEach(node => {
                        if (node.type === 'file') {
                            const li = document.createElement('li');
                            const a = document.createElement('a');
                            
                            const pathParts = node.path.split('/');
                            const filename = pathParts[pathParts.length - 1];
                            a.href = '/docs?page=' + encodeURIComponent(filename);
                            a.textContent = node.name;
                            
                            if (currentPage && filename === currentPage) {
                                a.classList.add('active');
                            }
                            li.appendChild(a);
                            parent.appendChild(li);
                        } else if (node.type === 'dir' && node.children) {
                            const li = document.createElement('li');
                            const span = document.createElement('span');
                            span.textContent = node.name;
                            span.style.fontWeight = 'bold';
                            span.style.display = 'block';
                            span.style.marginTop = '0.75rem';
                            span.style.marginBottom = '0.4rem';
                            span.style.fontSize = '0.85rem';
                            li.appendChild(span);
                            const ul = document.createElement('ul');
                            ul.style.listStyle = 'none';
                            ul.style.paddingLeft = '0.75rem';
                            renderNodes(node.children, ul);
                            li.appendChild(ul);
                            parent.appendChild(li);
                        }
                    });
                }
                
                renderNodes(nodes, toc);
            }
        } catch (error) {
            console.error('Error loading docs TOC:', error);
            const toc = document.getElementById('docs-toc');
            if (toc) {
                toc.innerHTML = '<li>Error loading table of contents</li>';
            }
            const tocStrip = document.getElementById('docs-toc-strip');
            if (tocStrip) {
                tocStrip.innerHTML = '<ul><li>Error loading table of contents</li></ul>';
            }
        }
    }
    
    // Always load TOC for navigation
    loadDocsTOC();
})();

