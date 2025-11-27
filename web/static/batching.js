// ./web/static/batching.js

const WATCHER_URL = 'http://localhost:8006'; // Default port

// Watcher Controls
document.getElementById('btn-start-watcher').addEventListener('click', async () => {
    const statusEl = document.getElementById('daemon-status');
    try {
        const response = await fetch(WATCHER_URL + '/start', { 
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
        });
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        const result = await response.json();
        if (statusEl) statusEl.textContent = "Starting...";
    } catch (e) {
        console.error("Failed to start watcher", e);
        if (statusEl) statusEl.textContent = "Error: " + e.message;
        alert("Failed to start watcher. Make sure the watcher daemon is running on port 8006.");
    }
});

document.getElementById('btn-stop-watcher').addEventListener('click', async () => {
    const statusEl = document.getElementById('daemon-status');
    try {
        const response = await fetch(WATCHER_URL + '/stop', { 
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
        });
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        const result = await response.json();
        if (statusEl) statusEl.textContent = "Stopping...";
    } catch (e) {
        console.error("Failed to stop watcher", e);
        if (statusEl) statusEl.textContent = "Error: " + e.message;
        alert("Failed to stop watcher. Make sure the watcher daemon is running on port 8006.");
    }
});

// SSE Connection
function connectSSE() {
    const statusEl = document.getElementById('status');
    const daemonStatusEl = document.getElementById('daemon-status');
    
    const evtSource = new EventSource(WATCHER_URL + '/events');
    
    evtSource.onmessage = function(event) {
        console.log("SSE Message:", event.data);
        
        const data = event.data;
        
        // Check if it's a file path (absolute path, likely contains / and ends with .adoc)
        // Also check for paths that don't contain status messages
        const isFilePath = data.includes('/') && 
                           (data.endsWith('.adoc') || data.endsWith('.asciidoc')) &&
                           !data.includes('Watcher') &&
                           !data.includes('Detected') &&
                           !data.includes('Error') &&
                           !data.includes('Running') &&
                           !data.includes('Processing');
        
        if (isFilePath) {
            // This is a file path from the watcher
            watcherQueue.add(data);
        } else {
            // Status message
            // Update navbar status
            if (statusEl) {
                statusEl.textContent = data;
                // Clear after 5 seconds if it's a success message or similar
                if (!data.includes("Error")) {
                     setTimeout(() => {
                         if (statusEl.textContent === data) {
                             statusEl.textContent = "";
                         }
                     }, 5000);
                }
            }

            // Check for status updates for the local indicator
            if (data.includes("Watcher started")) {
                if (daemonStatusEl) daemonStatusEl.textContent = "Running";
                watcherQueue.clear(); // Clear queue when watcher starts
            } else if (data.includes("Watcher stopped")) {
                if (daemonStatusEl) daemonStatusEl.textContent = "Stopped";
            } else if (data.includes("Watcher status:")) {
                if (daemonStatusEl) daemonStatusEl.textContent = data.replace("Watcher status: ", "");
            }
        }
    };

    evtSource.onerror = function(err) {
        console.error("EventSource failed:", err);
        if (daemonStatusEl) daemonStatusEl.textContent = "Disconnected";
        evtSource.close();
        // Try to reconnect in 5s
        setTimeout(connectSSE, 5000);
    };
}

// Initial connection
connectSSE();

// Set default temp path for Server Folder Batch
(async () => {
    try {
        const response = await fetch('/api/batch/default-temp-path');
        const result = await response.json();
        if (result.path) {
            document.getElementById('folderPath').value = result.path;
        }
    } catch (e) {
        console.error("Failed to fetch default temp path:", e);
    }
})();

// Original Event Listeners
document.getElementById('btn-upload-zip').addEventListener('click', async () => {
    const fileInput = document.getElementById('zipFile');
    const statusDiv = document.getElementById('zip-status');
    
    if (!fileInput.files[0]) {
        statusDiv.textContent = "Please select a file.";
        return;
    }

    const formData = new FormData();
    formData.append('zip', fileInput.files[0]);

    statusDiv.textContent = "Uploading...";
    try {
        const response = await fetch('/api/batch/upload-zip', {
            method: 'POST',
            body: formData
        });
        const result = await response.json();
        statusDiv.textContent = result.message || "Upload complete.";
        
        // Populate the Server Folder Batch path field with the extraction path
        if (result.path) {
            document.getElementById('folderPath').value = result.path;
        }
    } catch (e) {
        statusDiv.textContent = "Error: " + e.message;
    }
});

document.getElementById('btn-process-folder').addEventListener('click', async () => {
    const pathInput = document.getElementById('folderPath');
    const outputTypeSelect = document.getElementById('folder-output-type');
    const statusDiv = document.getElementById('folder-status');
    const actionsDiv = document.getElementById('folder-actions');
    const downloadBtn = document.getElementById('btn-download-zip');
    const cleanupBtn = document.getElementById('btn-cleanup');
    const outputFilesOnlyLabel = document.getElementById('output-files-only-label');
    const outputFilesOnlyCheckbox = document.getElementById('output-files-only');
    
    if (!pathInput.value) {
        statusDiv.textContent = "Please enter a path.";
        return;
    }

    // Hide action buttons
    actionsDiv.classList.remove('visible');
    downloadBtn.classList.add('hidden');
    cleanupBtn.classList.add('hidden');
    outputFilesOnlyLabel.classList.remove('visible');

    statusDiv.textContent = "Processing...";
    try {
        const response = await fetch('/api/batch/process-folder', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ 
                path: pathInput.value,
                outputType: outputTypeSelect.value
            })
        });
        const result = await response.json();
        if (result.error) {
            statusDiv.textContent = "Error: " + result.error;
        } else {
            statusDiv.textContent = result.message || "Processing completed.";
            
            // Show download and cleanup buttons if zip was created
            if (result.zipPath && result.successCount > 0 && result.errorCount === 0) {
                // Store data for download
                downloadBtn.dataset.zipPath = result.zipPath;
                downloadBtn.dataset.sourcePath = result.sourcePath || pathInput.value;
                downloadBtn.dataset.outputType = result.outputType || outputTypeSelect.value;
                
                // Show checkbox only if no subfolders
                if (!result.hasSubfolders) {
                    outputFilesOnlyLabel.classList.add('visible');
                } else {
                    outputFilesOnlyLabel.classList.remove('visible');
                    outputFilesOnlyCheckbox.checked = false;
                }
                
                downloadBtn.classList.remove('hidden');
                cleanupBtn.classList.remove('hidden');
                actionsDiv.classList.add('visible');
            }
        }
    } catch (e) {
        statusDiv.textContent = "Error: " + e.message;
    }
});

// Download zip button
document.getElementById('btn-download-zip').addEventListener('click', () => {
    const downloadBtn = document.getElementById('btn-download-zip');
    const zipPath = downloadBtn.dataset.zipPath;
    const sourcePath = downloadBtn.dataset.sourcePath;
    const outputType = downloadBtn.dataset.outputType;
    const outputFilesOnly = document.getElementById('output-files-only').checked;
    
    if (sourcePath && outputType) {
        // Create new zip with current options
        const url = `/api/batch/download-zip?sourcePath=${encodeURIComponent(sourcePath)}&outputType=${encodeURIComponent(outputType)}&outputFilesOnly=${outputFilesOnly}`;
        window.location.href = url;
    } else if (zipPath) {
        // Use existing zip
        window.location.href = `/api/batch/download-zip?path=${encodeURIComponent(zipPath)}`;
    }
});

// Cleanup button
document.getElementById('btn-cleanup').addEventListener('click', async () => {
    const statusDiv = document.getElementById('folder-status');
    const actionsDiv = document.getElementById('folder-actions');
    const pathInput = document.getElementById('folderPath');
    
    if (!confirm('Are you sure you want to clean up all batch processing files on the server?')) {
        return;
    }

    try {
        const response = await fetch('/api/batch/cleanup', {
            method: 'POST'
        });
        const result = await response.json();
        
        if (result.error) {
            statusDiv.textContent = "Error: " + result.error;
                } else {
                    statusDiv.textContent = result.message || "Cleanup completed.";
                    actionsDiv.classList.remove('visible');
                    // Reset form
                    pathInput.value = '';
                    document.getElementById('btn-download-zip').dataset.zipPath = '';
                }
    } catch (e) {
        statusDiv.textContent = "Error: " + e.message;
    }
});

document.getElementById('btn-set-watch').addEventListener('click', async () => {
    const pathInput = document.getElementById('watchPath');
    const statusDiv = document.getElementById('watch-status');
    
    if (!pathInput.value) {
        statusDiv.textContent = "Please enter a path.";
        return;
    }

    statusDiv.textContent = "Updating config...";
    try {
        const response = await fetch('/api/config/update', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ watchDir: pathInput.value })
        });
        const result = await response.json();
        statusDiv.textContent = result.message || "Configuration updated.";
    } catch (e) {
        statusDiv.textContent = "Error: " + e.message;
    }
});

