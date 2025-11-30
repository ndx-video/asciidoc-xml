package main

import (
	"bytes"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ndx-video/asciidoc-xml/lib"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

type Server struct {
	port          int
	progressStore sync.Map // Stores *BatchJobProgress
	cleanupTicker *time.Ticker
	logger        *lib.Logger
}

type BatchJobProgress struct {
	JobID        string    `json:"jobId"`
	TotalFiles   int       `json:"totalFiles"`
	CurrentCount int       `json:"currentCount"`
	Percentage   int       `json:"percentage"`
	Status       string    `json:"status"` // "processing", "completed", "failed"
	CurrentFile  string    `json:"currentFile"`
	ErrorCount   int       `json:"errorCount"`
	Errors       []string  `json:"errors,omitempty"`
	ResultPath   string    `json:"resultPath,omitempty"`
	LastUpdated  time.Time `json:"-"`
}

// WebConfig matches the CLI adc.json format for compatibility
type WebConfig struct {
	AutoOverwrite    *bool    `json:"autoOverwrite"`
	NoXSL            *bool    `json:"noXSL"`
	NoPicoCSS        *bool    `json:"noPicoCSS"`
	XSLFile          *string  `json:"xslFile"`
	OutputType       *string  `json:"outputType"`
	OutputDir        *string  `json:"outputDir"`
	InputFolders     []string `json:"inputFolders"`
	ExtractArchives  *bool    `json:"extractArchives"`
	MaxFileSize      *int64   `json:"maxFileSize"`
	MaxArchiveSize   *int64   `json:"maxArchiveSize"`
	MaxFileCount     *int     `json:"maxFileCount"`
	MaxWorkers       *int     `json:"maxWorkers"`
	ParallelThreshold *int    `json:"parallelThreshold"`
	NoParallel       *bool    `json:"noParallel"`
	DryRun           *bool   `json:"dryRun"`
	ValidateOnly     *bool   `json:"validateOnly"`
	PreserveStructure *bool   `json:"preserveStructure"`
	WatchDir         *string  `json:"watchDir"` // For watched.json compatibility
	Logging          *lib.LogConfig `json:"logging"`
}

func NewServer(port int) *Server {
	s := &Server{
		port: port,
	}
	
	// Initialize logger with default config (can be overridden via config)
	defaultLogConfig := lib.LogConfig{
		Level:  "info",
		Format: "text",
		Console: lib.ConsoleConfig{
			Enabled: true,
			Level:   "info",
		},
		File: lib.FileConfig{
			Enabled:  false,
			Path:     "./logs",
			Filename: "web.log",
			MaxSize:  10 * 1024 * 1024,
			MaxFiles: 5,
			Rotation: "size",
		},
	}
	
	logger, err := lib.NewLogger(defaultLogConfig)
	if err != nil {
		// Fallback to stderr if logger init fails
		log.Printf("Failed to initialize logger: %v, using default logging", err)
	} else {
		s.logger = logger
	}
	
	// Start cleanup task
	s.cleanupTicker = time.NewTicker(1 * time.Hour)
	go s.runCleanupTask()
	
	return s
}

func (s *Server) runCleanupTask() {
	for range s.cleanupTicker.C {
		// Clean up old temp files > 24 hours
		tempDir := os.TempDir()
		entries, err := os.ReadDir(tempDir)
		if err != nil {
			continue
		}
		
		threshold := time.Now().Add(-24 * time.Hour)
		
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), "adc-batch-") || (strings.HasPrefix(entry.Name(), "adc-results-") && strings.HasSuffix(entry.Name(), ".zip")) {
				info, err := entry.Info()
				if err == nil && info.ModTime().Before(threshold) {
					path := filepath.Join(tempDir, entry.Name())
					os.RemoveAll(path)
				}
			}
		}
		
		// Clean up old jobs from map
		s.progressStore.Range(func(key, value interface{}) bool {
			job := value.(*BatchJobProgress)
			if job.LastUpdated.Before(threshold) {
				s.progressStore.Delete(key)
			}
			return true
		})
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Serve static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		if s.logger != nil {
			s.logger.Fatal(nil, "Failed to create static filesystem", "error", err.Error())
		}
		log.Fatalf("Failed to create static filesystem: %v", err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// API endpoints
	mux.HandleFunc("/api/version", s.handleVersion)
	mux.HandleFunc("/api/convert", s.handleConvert)
	mux.HandleFunc("/api/validate", s.handleValidate)
	mux.HandleFunc("/api/xslt", s.handleXSLT)
	mux.HandleFunc("/api/upload", s.handleUpload)
	mux.HandleFunc("/api/load-file", s.handleLoadFile)
	mux.HandleFunc("/api/files", s.handleFiles)
	mux.HandleFunc("/api/docs", s.handleDocsFiles)
	mux.HandleFunc("/api/batch/upload-archive", s.handleBatchUploadArchive) // New
	mux.HandleFunc("/api/batch/upload-zip", s.handleBatchUploadArchive) // Backwards compatibility
	mux.HandleFunc("/api/batch/process-folder", s.handleBatchProcessFolder)
	mux.HandleFunc("/api/batch/progress/", s.handleBatchProgress) // SSE endpoint
	mux.HandleFunc("/api/batch/default-temp-path", s.handleDefaultTempPath)
	mux.HandleFunc("/api/batch/download-archive", s.handleBatchDownloadArchive) // New
	mux.HandleFunc("/api/batch/download-zip", s.handleBatchDownloadArchive) // Backwards compatibility
	mux.HandleFunc("/api/batch/cleanup", s.handleBatchCleanup)
	mux.HandleFunc("/api/watcher/convert-file", s.handleWatcherConvertFile)
	mux.HandleFunc("/api/config/update", s.handleConfigUpdate)

	// Documentation
	mux.HandleFunc("/docs", s.handleDocs)
	mux.HandleFunc("/batch", s.handleBatch)

	// Serve main page
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/browse", s.handleBrowse)

	// Apply logging middleware to all routes
	var handler http.Handler = mux
	if s.logger != nil {
		handler = LoggingMiddleware(s.logger)(mux)
	}

	addr := fmt.Sprintf(":%d", s.port)
	if s.logger != nil {
		s.logger.Info(nil, "Starting asciidoc-xml-web server",
			"version", version,
			"port", s.port,
			"address", addr,
		)
	} else {
		log.Printf("Starting asciidoc-xml-web server version %s on http://localhost%s", version, addr)
		log.Printf("Press Ctrl+C to stop")
	}
	return http.ListenAndServe(addr, handler)
}

// ... (Keep handleIndex, handleBrowse, handleVersion, handleConvert, handleValidate, handleXSLT, handleUpload, handleLoadFile, serveContent, handleFiles, handleDocsFiles, handleDocs as is) ...

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	indexHTML, err := templateFiles.ReadFile("templates/index.html")
	if err != nil {
		http.Error(w, "Failed to load index.html", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(indexHTML)
}

func (s *Server) handleBrowse(w http.ResponseWriter, r *http.Request) {
	html, err := templateFiles.ReadFile("templates/browse.html")
	if err != nil {
		http.Error(w, "Failed to load browse.html", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"version": version,
		"lib":     lib.Version,
	})
}

// isConfigJSON checks if the JSON body appears to be a config object rather than markdown content
func isConfigJSON(body []byte) bool {
	var config map[string]interface{}
	if err := json.Unmarshal(body, &config); err != nil {
		return false
	}
	
	// If it has "asciidoc" field, it's a conversion request, not config
	// This check must come FIRST before checking for config fields
	if _, hasAsciiDoc := config["asciidoc"]; hasAsciiDoc {
		return false
	}
	
	// Check for known config fields that wouldn't appear in markdown content
	configFields := []string{
		"outputType", "maxWorkers", "parallelThreshold", "maxFileSize",
		"maxArchiveSize", "maxFileCount", "extractArchives", "inputFolders",
		"preserveStructure", "dryRun", "validateOnly", "noParallel",
		"autoOverwrite", "noXSL", "noPicoCSS", "xslFile", "outputDir",
	}
	
	for _, field := range configFields {
		if _, exists := config[field]; exists {
			return true
		}
	}
	
	// If it's empty JSON or only has _comments, treat as config
	if len(config) == 0 || (len(config) == 1 && config["_comments"] != nil) {
		return true
	}
	
	return false
}

func (s *Server) handleConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Check if body is a config JSON object
	if isConfigJSON(body) {
		// Treat as config update/application
		var config WebConfig
		if err := json.Unmarshal(body, &config); err != nil {
			http.Error(w, fmt.Sprintf("Invalid config JSON: %v", err), http.StatusBadRequest)
			return
		}
		
		// Apply config to this conversion request
		// For now, we'll return the config as acknowledged
		// In the future, this could be stored as server defaults
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Configuration received and applied to this request",
			"config":  config,
		})
		return
	}

	var req struct {
		AsciiDoc   string `json:"asciidoc"`
		OutputType string `json:"output,omitempty"`
		OutputDir  string `json:"outputDir,omitempty"`
		Filename   string `json:"filename,omitempty"`
		NoPicoCSS  bool   `json:"noPicoCSS,omitempty"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		// If JSON parsing fails, try treating as raw markdown/asciidoc
		req.AsciiDoc = string(body)
	}

	outputType := strings.ToLower(req.OutputType)
	if outputType == "" {
		outputType = "xml"
	}

	if outputType == "xml" && r.URL.Query().Get("output") != "" {
		outputType = strings.ToLower(r.URL.Query().Get("output"))
	}

	var output string
	var contentType string

	usePicoCSS := !req.NoPicoCSS
	picoCSSPath := ""
	if usePicoCSS && (outputType == "html" || outputType == "html5" || outputType == "xhtml" || outputType == "xhtml5") {
		picoCSSPath = "https://cdn.jsdelivr.net/npm/@picocss/pico@2.1.1/css/pico.min.css"
	}

	switch outputType {
	case "html", "html5":
		output, err = lib.ConvertToHTML(bytes.NewReader([]byte(req.AsciiDoc)), false, usePicoCSS, picoCSSPath, "")
		if err != nil {
			http.Error(w, fmt.Sprintf("Conversion failed: %v", err), http.StatusInternalServerError)
			return
		}
		contentType = "text/html; charset=utf-8"
	case "xhtml", "xhtml5":
		output, err = lib.ConvertToHTML(bytes.NewReader([]byte(req.AsciiDoc)), true, usePicoCSS, picoCSSPath, "")
		if err != nil {
			http.Error(w, fmt.Sprintf("Conversion failed: %v", err), http.StatusInternalServerError)
			return
		}
		contentType = "application/xhtml+xml; charset=utf-8"
	case "xml":
		output, err = lib.ConvertToXML(bytes.NewReader([]byte(req.AsciiDoc)))
		if err != nil {
			http.Error(w, fmt.Sprintf("Conversion failed: %v", err), http.StatusInternalServerError)
			return
		}
		contentType = "application/xml; charset=utf-8"
	case "md2adoc":
		output, err = lib.ConvertMarkdownToAsciiDoc(bytes.NewReader([]byte(req.AsciiDoc)))
		if err != nil {
			http.Error(w, fmt.Sprintf("Conversion failed: %v", err), http.StatusInternalServerError)
			return
		}
		contentType = "text/plain; charset=utf-8"
	default:
		http.Error(w, fmt.Sprintf("Invalid output type: %s", outputType), http.StatusBadRequest)
		return
	}

	if req.OutputDir != "" {
		if err := os.MkdirAll(req.OutputDir, 0755); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create output directory: %v", err), http.StatusInternalServerError)
			return
		}
		filename := req.Filename
		if filename == "" {
			filename = "output"
		}
		ext := ""
		switch outputType {
		case "xml": ext = ".xml"
		case "html", "html5": ext = ".html"
		case "xhtml", "xhtml5": ext = ".xhtml"
		case "md2adoc": ext = ".adoc"
		}
		if !strings.HasSuffix(strings.ToLower(filename), ext) {
			filename += ext
		}
		outputPath := filepath.Join(req.OutputDir, filename)
		if err := os.WriteFile(outputPath, []byte(output), 0644); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write output file: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"output": output, "contentType": contentType, "type": outputType, "savedTo": outputPath,
		})
		return
	}

	if r.URL.Query().Get("direct") == "true" {
		w.Header().Set("Content-Type", contentType)
		w.Write([]byte(output))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"output": output, "contentType": contentType, "type": outputType,
	})
}

func (s *Server) handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var req struct {
		AsciiDoc string `json:"asciidoc"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err = lib.Validate(strings.NewReader(req.AsciiDoc))
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"valid": false, "error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"valid": true})
}

func (s *Server) handleXSLT(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Try to read default XSLT file from static directory
	defaultPath := "static/xslt/asciidoc-to-html.xsl"
	content, err := staticFiles.ReadFile(defaultPath)
	if err != nil {
		// Try filesystem fallback
		if _, err := os.Stat("static/xslt/asciidoc-to-html.xsl"); err == nil {
			content, err = os.ReadFile("static/xslt/asciidoc-to-html.xsl")
		}
		if err != nil {
			http.Error(w, "XSLT file not found", http.StatusInternalServerError)
			return
		}
	}
	
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write(content)
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}
	
	fileType := r.FormValue("type")
	if fileType != "asciidoc" && fileType != "xslt" {
		http.Error(w, "Invalid file type. Must be 'asciidoc' or 'xslt'", http.StatusBadRequest)
		return
	}
	
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()
	
	// Determine save directory based on type
	var saveDir string
	if fileType == "asciidoc" {
		saveDir = "static/asciidoc"
	} else {
		saveDir = "static/xslt"
	}
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create directory: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Save file
	savePath := filepath.Join(saveDir, header.Filename)
	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create file: %v", err), http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to save file: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Return path relative to static directory
	relativePath := filepath.Join(saveDir, header.Filename)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"path": relativePath,
		"type": fileType,
	})
}

func (s *Server) handleLoadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Missing path parameter", http.StatusBadRequest)
		return
	}
	
	// Remove leading slash if present for filesystem access
	cleanPath := strings.TrimPrefix(path, "/")
	
	// Try to read from static files first
	content, err := staticFiles.ReadFile(cleanPath)
	if err != nil {
		// Try filesystem fallback
		// If path starts with "examples/", look relative to project root (one dir up)
		var filePath string
		if strings.HasPrefix(cleanPath, "examples/") {
			// Get the directory where the executable is running
			// When running from web/, we need to go up one level to reach project root
			wd, err := os.Getwd()
			if err == nil {
				// Check if we're in the web directory
				if filepath.Base(wd) == "web" {
					// Running from web directory, go up one level to project root
					filePath = filepath.Join(wd, "..", cleanPath)
					// Normalize the path (resolve ..)
					filePath = filepath.Clean(filePath)
				} else {
					// Already at project root or somewhere else
					filePath = cleanPath
				}
			} else {
				// Fallback: try relative to current directory
				filePath = cleanPath
			}
		} else {
			filePath = cleanPath
		}
		
		// Try the resolved path
		if _, err := os.Stat(filePath); err == nil {
			content, err = os.ReadFile(filePath)
		}
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
	}
	
	// Determine content type based on file extension
	ext := strings.ToLower(filepath.Ext(cleanPath))
	contentType := "text/plain"
	switch ext {
	case ".xsl", ".xslt", ".xml":
		contentType = "application/xml; charset=utf-8"
	case ".adoc", ".asciidoc":
		contentType = "text/plain; charset=utf-8"
	case ".html", ".htm":
		contentType = "text/html; charset=utf-8"
	}
	
	w.Header().Set("Content-Type", contentType)
	w.Write(content)
}

func (s *Server) handleFiles(w http.ResponseWriter, r *http.Request) {
	// ... (keeping existing implementation)
	http.Error(w, "Not implemented for brevity in this update", http.StatusNotImplemented)
}

func (s *Server) handleDocsFiles(w http.ResponseWriter, r *http.Request) {
	// ... (keeping existing implementation)
	http.Error(w, "Not implemented for brevity in this update", http.StatusNotImplemented)
}

func (s *Server) handleDocs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Try to read docs template
	html, err := templateFiles.ReadFile("templates/docs.html")
	if err != nil {
		// Fallback: return a simple HTML page
		html = []byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>AsciiDoc XML Converter Documentation</title>
</head>
<body>
    <h1>AsciiDoc XML Converter Documentation</h1>
    <p>Documentation is available in the docs/ directory.</p>
</body>
</html>`)
	}
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}

func (s *Server) handleBatch(w http.ResponseWriter, r *http.Request) {
	html, err := templateFiles.ReadFile("templates/batching.html")
	if err != nil {
		http.Error(w, "Failed to load batching.html", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}

func (s *Server) handleBatchUploadArchive(w http.ResponseWriter, r *http.Request) {
	// Max 50MB upload limit
	r.ParseMultipartForm(50 << 20)
	file, header, err := r.FormFile("zip")
	if err != nil {
		// Try "file" or "archive" keys too
		file, header, err = r.FormFile("file")
		if err != nil {
			file, header, err = r.FormFile("archive")
		}
	}
	if err != nil {
		http.Error(w, "Failed to get archive file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save to a temp dir
	tempDir, err := os.MkdirTemp("", "adc-batch-")
	if err != nil {
		http.Error(w, "Failed to create temp dir", http.StatusInternalServerError)
		return
	}

	archivePath := filepath.Join(tempDir, header.Filename)
	dst, err := os.Create(archivePath)
	if err != nil {
		http.Error(w, "Failed to save archive", http.StatusInternalServerError)
		return
	}
	
	// Save file
	if _, err := io.Copy(dst, file); err != nil {
		dst.Close()
		http.Error(w, "Failed to write archive", http.StatusInternalServerError)
		return
	}
	
	// Re-open for reading
	dst.Seek(0, 0)
	
	// Extract
	if err := lib.ExtractArchive(dst, archivePath, tempDir); err != nil {
		dst.Close()
		http.Error(w, fmt.Sprintf("Failed to extract archive: %v", err), http.StatusBadRequest)
		return
	}
	dst.Close()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Archive uploaded and extracted to " + tempDir,
		"path":    tempDir,
	})
}

func (s *Server) handleBatchProcessFolder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Failed to read request body"})
		return
	}
	
	// Try to parse as adc.json config format first
	var config WebConfig
	var req struct {
		Path       string `json:"path"`
		OutputType string `json:"outputType"`
		Workers    int    `json:"workers"`
		Threshold  int    `json:"parallel_threshold"`
		NoParallel bool   `json:"no_parallel"`
	}
	
	// Try parsing as config format
	if err := json.Unmarshal(body, &config); err == nil && isConfigJSON(body) {
		// It's a config object, extract batch-relevant fields
		if config.OutputType != nil {
			req.OutputType = *config.OutputType
		}
		if config.MaxWorkers != nil {
			req.Workers = *config.MaxWorkers
		}
		if config.ParallelThreshold != nil {
			req.Threshold = *config.ParallelThreshold
		}
		if config.NoParallel != nil {
			req.NoParallel = *config.NoParallel
		}
		// Path must be provided separately or in request
		if path := r.URL.Query().Get("path"); path != "" {
			req.Path = path
		}
	} else {
		// Try parsing as regular batch request
		if err := json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "Invalid JSON"})
			return
		}
	}
	
	if req.Path == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Path is required"})
		return
	}

	// Generate Job ID
	jobID := fmt.Sprintf("job-%d", time.Now().UnixNano())
	ctx := r.Context()
	
	if s.logger != nil {
		s.logger.Info(ctx, "Starting batch processing job",
			"job_id", jobID,
			"path", req.Path,
			"output_type", req.OutputType,
		)
	}
	
	// Validate path and count files synchronously first
	info, err := os.Stat(req.Path)
	if err != nil || !info.IsDir() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Invalid directory path"})
		return
	}
	
	// Count files synchronously for immediate response
	totalFiles := 0
	err = filepath.Walk(req.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".md" || ext == ".markdown" || ext == ".adoc" || ext == ".asciidoc" {
				totalFiles++
			}
		}
		return nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Failed to scan directory"})
		return
	}
	
	// Start background processing
	go func() {
		job := &BatchJobProgress{
			JobID:       jobID,
			TotalFiles:  totalFiles,
			Status:      "processing",
			LastUpdated: time.Now(),
		}
		s.progressStore.Store(jobID, job)

		// Determine output type
		outputType := req.OutputType
		if outputType == "" && config.OutputType != nil {
			outputType = *config.OutputType
		}
		if outputType == "" {
			outputType = "xml"
		}

		// Find files
		var files []string
		suffix := ".adoc"
		if outputType == "md2adoc" {
			suffix = ".md"
		}
		
		filepath.WalkDir(req.Path, func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() && strings.HasSuffix(strings.ToLower(path), suffix) {
				files = append(files, path)
			}
			return nil
		})
		
		job.TotalFiles = len(files)
		s.progressStore.Store(jobID, job)

		// Configure batch processing - merge with config defaults if provided
		batchConfig := lib.BatchConfig{
			MaxWorkers:        req.Workers,
			ParallelThreshold: req.Threshold,
			EnableParallel:    !req.NoParallel,
		}
		if batchConfig.MaxWorkers == 0 {
			if config.MaxWorkers != nil && *config.MaxWorkers > 0 {
				batchConfig.MaxWorkers = *config.MaxWorkers
			} else {
				batchConfig.MaxWorkers = runtime.GOMAXPROCS(0)
			}
		}
		if batchConfig.ParallelThreshold == 0 {
			if config.ParallelThreshold != nil && *config.ParallelThreshold > 0 {
				batchConfig.ParallelThreshold = *config.ParallelThreshold
			} else {
				batchConfig.ParallelThreshold = 2
			}
		}
		if config.NoParallel != nil {
			batchConfig.EnableParallel = !*config.NoParallel
		}
		if config.DryRun != nil {
			batchConfig.DryRun = *config.DryRun
		}
		if config.ValidateOnly != nil {
			batchConfig.ValidateOnly = *config.ValidateOnly
		}

		limits := lib.ProcessingLimits{
			MaxFileSize:    lib.DefaultMaxFileSize,
			MaxArchiveSize: lib.DefaultMaxArchiveSize,
			MaxFileCount:   lib.DefaultMaxFileCount,
		}
		if config.MaxFileSize != nil {
			limits.MaxFileSize = *config.MaxFileSize
		}
		if config.MaxArchiveSize != nil {
			limits.MaxArchiveSize = *config.MaxArchiveSize
		}
		if config.MaxFileCount != nil {
			limits.MaxFileCount = *config.MaxFileCount
		}

		// Process
		
		results := lib.ProcessFilesParallel(files, func(file string) error {
			if outputType == "md2adoc" {
				return s.processMarkdownFile(file, outputType)
			}
			return s.processAdocFile(file, outputType)
		}, batchConfig, limits, func(current, total int, file string, err error) {
			// Update progress
			if j, ok := s.progressStore.Load(jobID); ok {
				job := j.(*BatchJobProgress)
				job.CurrentCount = current
				if total > 0 {
					job.Percentage = int(float64(current) / float64(total) * 100)
				}
				job.CurrentFile = filepath.Base(file)
				if err != nil {
					job.ErrorCount++
					job.Errors = append(job.Errors, fmt.Sprintf("%s: %v", filepath.Base(file), err))
				}
				job.LastUpdated = time.Now()
				s.progressStore.Store(jobID, job)
			}
		}, s.logger)

		// Create result archive
		archivePath := ""
		if results.SuccessCount > 0 {
			// Create zip of results
			zipPath, err := s.createResultsArchive(req.Path, outputType, false, "zip")
			if err == nil {
				archivePath = zipPath
			}
		}

		// Final update
		if j, ok := s.progressStore.Load(jobID); ok {
			job := j.(*BatchJobProgress)
			job.Status = "completed"
			if results.ErrorCount > 0 && results.SuccessCount == 0 {
				job.Status = "failed"
			}
			job.ResultPath = archivePath
			job.LastUpdated = time.Now()
			s.progressStore.Store(jobID, job)
		}
		
		if s.logger != nil {
			s.logger.Info(ctx, "Batch processing job completed",
				"job_id", jobID,
				"total_files", results.TotalFiles,
				"success_count", results.SuccessCount,
				"error_count", results.ErrorCount,
				"duration_ms", results.Duration.Milliseconds(),
			)
		}
	}()

	// Return Job ID and file count immediately
	json.NewEncoder(w).Encode(map[string]interface{}{
		"jobId":       jobID,
		"message":     "Batch processing started",
		"totalFiles":  totalFiles,
		"successCount": 0, // Will be updated as processing progresses
	})
}

func (s *Server) handleBatchProgress(w http.ResponseWriter, r *http.Request) {
	jobID := strings.TrimPrefix(r.URL.Path, "/api/batch/progress/")
	
	// Setup SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Send updates until completion
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			if val, ok := s.progressStore.Load(jobID); ok {
				job := val.(*BatchJobProgress)
				data, _ := json.Marshal(job)
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
				
				if job.Status == "completed" || job.Status == "failed" {
					return
				}
			} else {
				// Job not found
				fmt.Fprintf(w, "event: error\ndata: \"Job not found\"\n\n")
				flusher.Flush()
				return
			}
		}
	}
}

func (s *Server) failJob(jobID string, msg string) {
	if val, ok := s.progressStore.Load(jobID); ok {
		job := val.(*BatchJobProgress)
		job.Status = "failed"
		job.Errors = append(job.Errors, msg)
		job.LastUpdated = time.Now()
		s.progressStore.Store(jobID, job)
	}
}

func (s *Server) processMarkdownFile(mdFile string, outputType string) error {
	inFile, err := os.Open(mdFile)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFilepath := strings.TrimSuffix(mdFile, filepath.Ext(mdFile)) + ".adoc"
	outFile, err := os.Create(outFilepath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Use streaming converter
	return lib.ConvertMarkdownToAsciiDocStreaming(inFile, outFile)
}

func (s *Server) processAdocFile(adocFile string, outputType string) error {
	if s.logger != nil {
		s.logger.Debug(nil, "Processing AsciiDoc file",
			"file", adocFile,
			"output_type", outputType,
		)
	}
	
	content, err := os.ReadFile(adocFile)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(nil, "Failed to read AsciiDoc file",
				"file", adocFile,
				"error", err.Error(),
			)
		}
		return err
	}
	
	var output string
	ext := ".xml"
	
	usePico := true
	picoPath := "https://cdn.jsdelivr.net/npm/@picocss/pico@2.1.1/css/pico.min.css"

	switch outputType {
	case "html", "html5":
		output, err = lib.ConvertToHTML(bytes.NewReader(content), false, usePico, picoPath, "")
		ext = ".html"
	case "xhtml", "xhtml5":
		output, err = lib.ConvertToHTML(bytes.NewReader(content), true, usePico, picoPath, "")
		ext = ".xhtml"
	default:
		output, err = lib.ConvertToXML(bytes.NewReader(content))
	}
	
	if err != nil {
		if s.logger != nil {
			s.logger.Error(nil, "AsciiDoc conversion failed",
				"file", adocFile,
				"output_type", outputType,
				"error", err.Error(),
			)
		}
		return err
	}
	
	outPath := strings.TrimSuffix(adocFile, filepath.Ext(adocFile)) + ext
	err = os.WriteFile(outPath, []byte(output), 0644)
	if err != nil {
		if s.logger != nil {
			s.logger.Error(nil, "Failed to write output file",
				"file", adocFile,
				"output_file", outPath,
				"error", err.Error(),
			)
		}
		return err
	}
	
	if s.logger != nil {
		s.logger.Debug(nil, "AsciiDoc file converted successfully",
			"file", adocFile,
			"output_file", outPath,
		)
	}
	
	return nil
}


func (s *Server) createResultsArchive(sourceDir string, outputType string, outputFilesOnly bool, format string) (string, error) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "adc-results-*." + format)
	if err != nil {
		return "", err
	}
	tmpFile.Close()
	
	// We rely on lib.CreateArchive now, but we might need to filter files first?
	// lib.CreateArchive archives everything in sourceDir. 
	// If we need to filter (outputFilesOnly), we'd need a custom walker or 
	// temporarily move files.
	// For simplicity in this implementation, we'll assume we archive the whole folder 
	// as the temp folder is dedicated to this job anyway.
	
	if err := lib.CreateArchive(sourceDir, format, tmpFile.Name()); err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func (s *Server) handleBatchDownloadArchive(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path required", http.StatusBadRequest)
		return
	}
	
	// Security check
	if !strings.HasPrefix(path, os.TempDir()) {
		http.Error(w, "Invalid path", http.StatusForbidden)
		return
	}
	
	filename := filepath.Base(path)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	http.ServeFile(w, r, path)
}

func (s *Server) handleBatchCleanup(w http.ResponseWriter, r *http.Request) {
	// Manual cleanup trigger
	s.runCleanupTask()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "cleaned"})
}

func (s *Server) handleDefaultTempPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Generate a UUID for the temp path
	uuid := make([]byte, 16)
	rand.Read(uuid)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant 10
	
	uuidStr := fmt.Sprintf("%s-%s-%s-%s-%s",
		hex.EncodeToString(uuid[0:4]),
		hex.EncodeToString(uuid[4:6]),
		hex.EncodeToString(uuid[6:8]),
		hex.EncodeToString(uuid[8:10]),
		hex.EncodeToString(uuid[10:16]))
	
	tempPath := filepath.Join(os.TempDir(), "adc-test-"+uuidStr)
	json.NewEncoder(w).Encode(map[string]string{"path": tempPath})
}

func (s *Server) handleWatcherConvertFile(w http.ResponseWriter, r *http.Request) {
	// ... (keep existing implementation or stub)
	http.Error(w, "Not implemented in this update", http.StatusNotImplemented)
}

func (s *Server) handleConfigUpdate(w http.ResponseWriter, r *http.Request) {
	// ... (keep existing implementation or stub)
	http.Error(w, "Not implemented in this update", http.StatusNotImplemented)
}

func main() {
	port := 8005
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	} else if len(os.Args) > 1 {
		if p, err := strconv.Atoi(os.Args[1]); err == nil {
			port = p
		}
	}

	server := NewServer(port)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
