// Watcher Queue Management
// This module manages the queue of files to be processed by the watcher

const watcherQueue = {
    items: [],
    processing: false,
    outputType: 'xml', // Default output type
    
    add: function(filePath) {
        // Check if already in queue
        if (this.items.some(item => item.path === filePath)) {
            return;
        }
        
        this.items.push({
            path: filePath,
            status: 'pending', // pending, processing, completed, error
            error: null
        });
        this.updateDisplay();
        this.processNext();
    },
    
    updateDisplay: function() {
        const container = document.getElementById('watcher-queue-container');
        const queueDiv = document.getElementById('watcher-queue');
        
        if (!container || !queueDiv) {
            return; // Elements not available (e.g., in test environment)
        }
        
        if (this.items.length === 0) {
            container.style.display = 'none';
            queueDiv.innerHTML = '<p style="margin: 0; color: #666;">No files in queue</p>';
            return;
        }
        
        container.style.display = 'block';
        let html = '<div style="display: flex; flex-direction: column; gap: 5px;">';
        
        this.items.forEach((item, index) => {
            const statusColor = {
                'pending': '#ffc107',
                'processing': '#007bff',
                'completed': '#28a745',
                'error': '#dc3545'
            }[item.status] || '#666';
            
            const fileName = item.path.split('/').pop() || item.path;
            html += `
                <div style="display: flex; align-items: center; gap: 10px; padding: 5px; background: white; border-radius: 4px; border-left: 3px solid ${statusColor};">
                    <span style="flex: 1; font-family: monospace; font-size: 0.9em; word-break: break-all;">${fileName}</span>
                    <span style="padding: 2px 8px; background: ${statusColor}; color: white; border-radius: 3px; font-size: 0.8em; text-transform: capitalize;">${item.status}</span>
                    ${item.error ? `<span style="color: #dc3545; font-size: 0.8em;">${item.error}</span>` : ''}
                </div>
            `;
        });
        
        html += '</div>';
        queueDiv.innerHTML = html;
    },
    
    processNext: async function() {
        if (this.processing) return;
        
        const pendingItem = this.items.find(item => item.status === 'pending');
        if (!pendingItem) return;
        
        this.processing = true;
        pendingItem.status = 'processing';
        this.updateDisplay();
        
        try {
            const response = await fetch('/api/watcher/convert-file', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    filePath: pendingItem.path,
                    outputType: this.outputType
                })
            });
            
            const result = await response.json();
            
            if (result.error) {
                pendingItem.status = 'error';
                pendingItem.error = result.error;
            } else {
                pendingItem.status = 'completed';
            }
        } catch (error) {
            pendingItem.status = 'error';
            pendingItem.error = error.message;
        }
        
        this.updateDisplay();
        this.processing = false;
        
        // Process next item
        setTimeout(() => this.processNext(), 100);
    },
    
    clear: function() {
        this.items = [];
        this.updateDisplay();
    }
};

