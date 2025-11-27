package main

import (
	"archive/zip"
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
	"sort"
	"strconv"
	"strings"

	"asciidoc-xml/lib"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

type Server struct {
	port int
}

func NewServer(port int) *Server {
	return &Server{port: port}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Serve static files
	// Create a subdirectory FS for the static directory
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
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
	mux.HandleFunc("/api/batch/upload-zip", s.handleBatchUploadZip)
	mux.HandleFunc("/api/batch/process-folder", s.handleBatchProcessFolder)
	mux.HandleFunc("/api/batch/default-temp-path", s.handleDefaultTempPath)
	mux.HandleFunc("/api/batch/download-zip", s.handleBatchDownloadZip)
	mux.HandleFunc("/api/batch/cleanup", s.handleBatchCleanup)
	mux.HandleFunc("/api/watcher/convert-file", s.handleWatcherConvertFile)
	mux.HandleFunc("/api/config/update", s.handleConfigUpdate)

	// Documentation
	mux.HandleFunc("/docs", s.handleDocs)
	mux.HandleFunc("/batch", s.handleBatch)

	// Serve main page
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/browse", s.handleBrowse)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting asciidoc-xml-web server version %s on http://localhost%s", version, addr)
	log.Printf("Press Ctrl+C to stop")
	return http.ListenAndServe(addr, mux)
}

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

	var req struct {
		AsciiDoc   string `json:"asciidoc"`
		OutputType string `json:"output,omitempty"` // "xml", "html", "xhtml"
		OutputDir  string `json:"outputDir,omitempty"`
		Filename   string `json:"filename,omitempty"`
		NoPicoCSS  bool   `json:"noPicoCSS,omitempty"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Default to XML if not specified
	outputType := strings.ToLower(req.OutputType)
	if outputType == "" {
		outputType = "xml"
	}

	// Support query parameter as alternative
	if outputType == "xml" && r.URL.Query().Get("output") != "" {
		outputType = strings.ToLower(r.URL.Query().Get("output"))
	}

	var output string
	var contentType string

	// Determine PicoCSS usage (enabled by default unless noPicoCSS is true)
	// Use CDN for iframe content since blob URLs can't resolve relative paths
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
		// Convert Markdown to AsciiDoc
		output, err = lib.ConvertMarkdownToAsciiDoc(bytes.NewReader([]byte(req.AsciiDoc)))
		if err != nil {
			http.Error(w, fmt.Sprintf("Conversion failed: %v", err), http.StatusInternalServerError)
			return
		}
		contentType = "text/plain; charset=utf-8"

	default:
		http.Error(w, fmt.Sprintf("Invalid output type: %s. Supported types: xml, html, xhtml, md2adoc", outputType), http.StatusBadRequest)
		return
	}

	// Save to file if output directory is specified
	if req.OutputDir != "" {
		// Ensure output directory exists
		if err := os.MkdirAll(req.OutputDir, 0755); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create output directory: %v", err), http.StatusInternalServerError)
			return
		}

		// Determine filename
		filename := req.Filename
		if filename == "" {
			filename = "output"
		}

		// Ensure extension
		ext := ""
		switch outputType {
		case "xml":
			ext = ".xml"
		case "html", "html5":
			ext = ".html"
		case "xhtml", "xhtml5":
			ext = ".xhtml"
		case "md2adoc":
			ext = ".adoc"
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
			"output":      output,
			"contentType": contentType,
			"type":        outputType,
			"savedTo":     outputPath,
		})
		return
	}

	// Return JSON response for API compatibility, or direct output if requested
	if r.URL.Query().Get("direct") == "true" {
		w.Header().Set("Content-Type", contentType)
		w.Write([]byte(output))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"output":      output,
		"contentType": contentType,
		"type":        outputType,
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

	// Validate by attempting to parse
	err = lib.Validate(strings.NewReader(req.AsciiDoc))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	// If we get here, it's valid
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid": true,
	})
}

func (s *Server) handleXSLT(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Try to read from filesystem first (for development)
	xsltPath := filepath.Join("..", "xslt", "asciidoc-to-html.xsl")
	xsltContent, err := os.ReadFile(xsltPath)
	if err != nil {
		// Try current directory
		xsltPath = filepath.Join("xslt", "asciidoc-to-html.xsl")
		xsltContent, err = os.ReadFile(xsltPath)
		if err != nil {
			// Try parent directory
			xsltPath = filepath.Join("..", "..", "xslt", "asciidoc-to-html.xsl")
			xsltContent, err = os.ReadFile(xsltPath)
		}
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read XSLT file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(xsltContent)
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (10MB max)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	fileType := r.FormValue("type") // "asciidoc" or "xslt"
	if fileType != "asciidoc" && fileType != "xslt" {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Determine filename and extension
	var filename string
	if fileType == "asciidoc" {
		filename = header.Filename
		if !strings.HasSuffix(strings.ToLower(filename), ".adoc") && !strings.HasSuffix(strings.ToLower(filename), ".asciidoc") {
			filename += ".adoc"
		}
	} else {
		filename = header.Filename
		if !strings.HasSuffix(strings.ToLower(filename), ".xsl") && !strings.HasSuffix(strings.ToLower(filename), ".xslt") {
			filename += ".xsl"
		}
	}

	// Save to static directory
	staticDir := filepath.Join("..", "web", "static")
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		staticDir = filepath.Join("static")
	}

	targetPath := filepath.Join(staticDir, filename)
	dst, err := os.Create(targetPath)
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

	// Return the path relative to /static/
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"path":     "/static/" + filename,
		"filename": filename,
		"type":     fileType,
	})
}

func (s *Server) handleLoadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path parameter required", http.StatusBadRequest)
		return
	}

	// Security check
	if strings.Contains(path, "..") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Handle prefix
	path = strings.TrimPrefix(path, "/")

	var fsPath string

	if strings.HasPrefix(path, "static/") {
		path = strings.TrimPrefix(path, "static/")
		// Try embedded static files first
		embeddedPath := filepath.Join("static", path)
		if content, err := staticFiles.ReadFile(embeddedPath); err == nil {
			serveContent(w, path, content)
			return
		}

		// Try filesystem locations for static
		locations := []string{
			filepath.Join("..", "web", "static", path),
			filepath.Join("static", path),
			filepath.Join("..", "static", path),
			path,
		}

		for _, loc := range locations {
			if content, err := os.ReadFile(loc); err == nil {
				serveContent(w, path, content)
				return
			}
		}
	} else if strings.HasPrefix(path, "examples/") {
		// For examples, we look in the examples directory
		path = strings.TrimPrefix(path, "examples/")
		fsPath = filepath.Join("..", "examples", path)
		if _, err := os.Stat(fsPath); os.IsNotExist(err) {
			// Try current directory
			fsPath = filepath.Join("examples", path)
		}
	} else {
		// Fallback for simple names, assume static
		fsPath = filepath.Join("..", "web", "static", path)
		if _, err := os.Stat(fsPath); os.IsNotExist(err) {
			fsPath = filepath.Join("static", path)
		}
	}

	content, err := os.ReadFile(fsPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read file: %v", err), http.StatusNotFound)
		return
	}

	serveContent(w, path, content)
}

func serveContent(w http.ResponseWriter, path string, content []byte) {
	// Determine content type
	contentType := "text/plain"
	if strings.HasSuffix(strings.ToLower(path), ".adoc") || strings.HasSuffix(strings.ToLower(path), ".asciidoc") {
		contentType = "text/plain"
	} else if strings.HasSuffix(strings.ToLower(path), ".xsl") || strings.HasSuffix(strings.ToLower(path), ".xslt") {
		contentType = "application/xml"
	} else if strings.HasSuffix(strings.ToLower(path), ".html") {
		contentType = "text/html"
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(content)
}

type FileNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	Type     string     `json:"type"` // "file" or "dir"
	Children []FileNode `json:"children,omitempty"`
}

func (s *Server) handleFiles(w http.ResponseWriter, r *http.Request) {
	root := "examples"
	// Adjust root path for dev vs deployment
	if _, err := os.Stat(root); os.IsNotExist(err) {
		root = filepath.Join("..", "examples")
	}

	nodes, err := walkDir(root, "examples")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to walk directory: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

// extractTitleFromAdoc reads the first level 1 title from an AsciiDoc file
func extractTitleFromAdoc(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "" // Return empty if file can't be read
	}

	// Read line by line to find the first level 1 title (= Title)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Check for level 1 title: starts with "= " followed by non-whitespace
		if strings.HasPrefix(trimmed, "= ") && len(trimmed) > 2 {
			title := strings.TrimSpace(trimmed[2:])
			if title != "" {
				return title
			}
		}
	}

	return "" // No title found
}

func (s *Server) handleDocsFiles(w http.ResponseWriter, r *http.Request) {
	root := "docs"
	// Adjust root path for dev vs deployment
	if _, err := os.Stat(root); os.IsNotExist(err) {
		root = filepath.Join("..", "docs")
		if _, err := os.Stat(root); os.IsNotExist(err) {
			root = filepath.Join("..", "..", "docs")
		}
	}

	nodes, err := walkDirWithTitles(root, "docs")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to walk directory: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

func walkDir(fsPath, relPath string) ([]FileNode, error) {
	entries, err := os.ReadDir(fsPath)
	if err != nil {
		return nil, err
	}

	var nodes []FileNode
	for _, entry := range entries {
		if entry.Name() == "." || entry.Name() == ".." || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Filter: show dirs and .adoc files
		if !entry.IsDir() && !strings.HasSuffix(entry.Name(), ".adoc") && !strings.HasSuffix(entry.Name(), ".asciidoc") {
			continue
		}

		node := FileNode{
			Name: entry.Name(),
			Path: filepath.Join(relPath, entry.Name()),
			Type: "file",
		}

		if entry.IsDir() {
			node.Type = "dir"
			children, err := walkDir(filepath.Join(fsPath, entry.Name()), filepath.Join(relPath, entry.Name()))
			if err != nil {
				return nil, err
			}
			// Only add dirs if they have content
			if len(children) > 0 {
				node.Children = children
				nodes = append(nodes, node)
			}
		} else {
			nodes = append(nodes, node)
		}
	}

	// Sort: Dirs first, then files
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Type != nodes[j].Type {
			return nodes[i].Type == "dir"
		}
		return nodes[i].Name < nodes[j].Name
	})

	return nodes, nil
}

// walkDirWithTitles is like walkDir but extracts level 1 titles from .adoc files for display names
func walkDirWithTitles(fsPath, relPath string) ([]FileNode, error) {
	entries, err := os.ReadDir(fsPath)
	if err != nil {
		return nil, err
	}

	var nodes []FileNode
	for _, entry := range entries {
		if entry.Name() == "." || entry.Name() == ".." || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Filter: show dirs and .adoc files
		if !entry.IsDir() && !strings.HasSuffix(entry.Name(), ".adoc") && !strings.HasSuffix(entry.Name(), ".asciidoc") {
			continue
		}

		node := FileNode{
			Path: filepath.Join(relPath, entry.Name()),
			Type: "file",
		}

		if entry.IsDir() {
			node.Type = "dir"
			node.Name = entry.Name()
			children, err := walkDirWithTitles(filepath.Join(fsPath, entry.Name()), filepath.Join(relPath, entry.Name()))
			if err != nil {
				return nil, err
			}
			// Only add dirs if they have content
			if len(children) > 0 {
				node.Children = children
				nodes = append(nodes, node)
			}
		} else {
			// For files, try to extract title from the .adoc file
			filePath := filepath.Join(fsPath, entry.Name())
			title := extractTitleFromAdoc(filePath)
			if title != "" {
				node.Name = title
			} else {
				// Fallback to filename if no title found
				node.Name = entry.Name()
			}
			nodes = append(nodes, node)
		}
	}

	// Sort: Dirs first, then files
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Type != nodes[j].Type {
			return nodes[i].Type == "dir"
		}
		return nodes[i].Name < nodes[j].Name
	})

	return nodes, nil
}

func (s *Server) handleDocs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if requesting specific page
	page := r.URL.Query().Get("page")
	
	// Load the template first
	templateHTML, err := templateFiles.ReadFile("templates/docs.html")
	if err != nil {
		http.Error(w, "Failed to load docs.html", http.StatusInternalServerError)
		return
	}

	// If no page specified, serve empty template (ToC will load via JS)
	if page == "" {
		templateStr := string(templateHTML)
		templateStr = strings.Replace(templateStr, `<div id="docs-content-placeholder"></div>`, "", 1)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(templateStr))
		return
	}

	// Sanitize page path to prevent directory traversal
	page = filepath.Base(page)
	if !strings.HasSuffix(page, ".adoc") && !strings.HasSuffix(page, ".asciidoc") {
		// Invalid page format, serve empty template
		templateStr := string(templateHTML)
		templateStr = strings.Replace(templateStr, `<div id="docs-content-placeholder"></div>`, "", 1)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(templateStr))
		return
	}

	// Read the AsciiDoc documentation file
	var docPath string
	var docContent []byte

	// Try multiple locations for the docs directory
	docPath = filepath.Join("..", "docs", page)
	docContent, err = os.ReadFile(docPath)
	if err != nil {
		// Try current directory
		docPath = filepath.Join("docs", page)
		docContent, err = os.ReadFile(docPath)
		if err != nil {
			// Try parent directory
			docPath = filepath.Join("..", "..", "docs", page)
			docContent, err = os.ReadFile(docPath)
		}
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read documentation: %v", err), http.StatusNotFound)
		return
	}

	// Convert AsciiDoc directly to HTML5 (with PicoCSS enabled by default)
	htmlContent, err := lib.ConvertToHTML(bytes.NewReader(docContent), false, true, "/static/pico.min.css", "")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to convert documentation: %v", err), http.StatusInternalServerError)
		return
	}

	// Extract body content from converted HTML
	htmlLines := strings.Split(htmlContent, "\n")
	inBody := false
	bodyContent := strings.Builder{}

	for _, line := range htmlLines {
		if strings.Contains(line, "<body>") {
			inBody = true
			continue
		}
		if strings.Contains(line, "</body>") {
			break
		}
		if inBody {
			bodyContent.WriteString(line)
			bodyContent.WriteString("\n")
		}
	}

	// Replace placeholder in template with actual content
	// We'll add a placeholder div in the template
	templateStr := string(templateHTML)
	templateStr = strings.Replace(templateStr, `<div id="docs-content-placeholder"></div>`, bodyContent.String(), 1)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(templateStr))
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

func (s *Server) handleBatchUploadZip(w http.ResponseWriter, r *http.Request) {
	// Max 50MB
	r.ParseMultipartForm(50 << 20)
	file, header, err := r.FormFile("zip")
	if err != nil {
		http.Error(w, "Failed to get zip file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save to a temp dir
	tempDir, err := os.MkdirTemp("", "adc-batch-")
	if err != nil {
		http.Error(w, "Failed to create temp dir", http.StatusInternalServerError)
		return
	}

	zipPath := filepath.Join(tempDir, header.Filename)
	dst, err := os.Create(zipPath)
	if err != nil {
		http.Error(w, "Failed to save zip", http.StatusInternalServerError)
		return
	}
	io.Copy(dst, file)
	dst.Close()

	// Unzip
	r_zip, err := zip.OpenReader(zipPath)
	if err != nil {
		http.Error(w, "Failed to open zip", http.StatusInternalServerError)
		return
	}
	defer r_zip.Close()

	for _, f := range r_zip.File {
		fpath := filepath.Join(tempDir, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(tempDir)+string(os.PathSeparator)) {
			continue // illegal path
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			continue
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			continue
		}
		io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
	}

	// Respond with success and path
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Zip uploaded and extracted to " + tempDir,
		"path":    tempDir,
	})
}

func (s *Server) handleBatchProcessFolder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var req struct {
		Path       string `json:"path"`
		OutputType string `json:"outputType"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid JSON",
		})
		return
	}

	// Validate path exists
	if req.Path == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Path is required",
		})
		return
	}

	info, err := os.Stat(req.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": fmt.Sprintf("Path does not exist: %v", err),
		})
		return
	}

	if !info.IsDir() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Path must be a directory",
		})
		return
	}

	// Default to xml if not specified
	outputType := req.OutputType
	if outputType == "" {
		outputType = "xml"
	}

	// Find files based on output type
	var files []string
	var errorMsg string

	if outputType == "md2adoc" {
		// For md2adoc, find .md files
		errorMsg = "No .md files found in directory"
		err = filepath.WalkDir(req.Path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".md") {
				files = append(files, path)
			}
			return nil
		})
	} else {
		// For other types, find .adoc files
		errorMsg = "No .adoc files found in directory"
		err = filepath.WalkDir(req.Path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".adoc") {
				files = append(files, path)
			}
			return nil
		})
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": fmt.Sprintf("Failed to traverse directory: %v", err),
		})
		return
	}

	if len(files) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": errorMsg,
		})
		return
	}

	// Process each file
	successCount := 0
	errorCount := 0
	var errors []string

	for _, file := range files {
		if outputType == "md2adoc" {
			if err := s.processMarkdownFile(file, outputType); err != nil {
				errorCount++
				errors = append(errors, fmt.Sprintf("%s: %v", file, err))
			} else {
				successCount++
			}
		} else {
			if err := s.processAdocFile(file, outputType); err != nil {
				errorCount++
				errors = append(errors, fmt.Sprintf("%s: %v", file, err))
			} else {
				successCount++
			}
		}
	}

	// Check if directory has subfolders
	hasSubfolders := s.hasSubfolders(req.Path)

	// Format output type for display
	outputTypeDisplay := strings.ToUpper(outputType)
	if outputType == "html5" {
		outputTypeDisplay = "HTML5"
	} else if outputType == "xhtml5" {
		outputTypeDisplay = "XHTML5"
	} else if outputType == "md2adoc" {
		outputTypeDisplay = "MD2ADoc"
	}

	response := map[string]interface{}{
		"message":        fmt.Sprintf("Batch processing completed for %s (%s output)", req.Path, outputTypeDisplay),
		"successCount":   successCount,
		"errorCount":     errorCount,
		"totalFiles":     len(files),
		"hasSubfolders":  hasSubfolders,
		"outputType":     outputType,
	}

	// Create zip file if processing was successful
	var zipPath string
	if successCount > 0 && errorCount == 0 {
		// Default to including all files (outputFilesOnly=false)
		outputFilesOnly := false
		zipPath, err = s.createResultsZip(req.Path, outputType, outputFilesOnly)
		if err != nil {
			response["zipError"] = fmt.Sprintf("Failed to create zip file: %v", err)
		} else {
			response["zipPath"] = zipPath
			response["sourcePath"] = req.Path
		}
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	json.NewEncoder(w).Encode(response)
}

func (s *Server) processMarkdownFile(mdFile string, outputType string) error {
	// Read the Markdown file
	file, err := os.Open(mdFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Convert Markdown to AsciiDoc
	adocOutput, err := lib.ConvertMarkdownToAsciiDoc(file)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	// Determine output file path
	baseName := strings.TrimSuffix(mdFile, filepath.Ext(mdFile))
	outputFile := baseName + ".adoc"

	// Write the AsciiDoc output
	if err := os.WriteFile(outputFile, []byte(adocOutput), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func (s *Server) processAdocFile(adocFile string, outputType string) error {
	// Read the AsciiDoc file
	file, err := os.Open(adocFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Determine output file path and extension
	baseName := strings.TrimSuffix(adocFile, filepath.Ext(adocFile))
	var outputFile string
	var output string

	// Use CDN link for batch processing (version pinned to match local static file)
	picoCDNPath := "https://cdn.jsdelivr.net/npm/@picocss/pico@2.1.1/css/pico.min.css"

	switch outputType {
	case "xml":
		outputFile = baseName + ".xml"
		output, err = lib.ConvertToXML(file)
	case "html", "html5":
		outputFile = baseName + ".html"
		output, err = lib.ConvertToHTML(file, false, true, picoCDNPath, "")
	case "xhtml", "xhtml5":
		outputFile = baseName + ".xhtml"
		output, err = lib.ConvertToHTML(file, true, true, picoCDNPath, "")
	default:
		return fmt.Errorf("unsupported output type: %s", outputType)
	}

	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	// Write output file
	if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func (s *Server) hasSubfolders(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() {
			return true
		}
	}
	return false
}

func (s *Server) createResultsZip(sourceDir string, outputType string, outputFilesOnly bool) (string, error) {
	// Create zip file in temp directory
	zipFile, err := os.CreateTemp("", "adc-results-*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}
	zipPath := zipFile.Name()

	zipWriter := zip.NewWriter(zipFile)

	// Determine file extension for filtering
	var ext string
	if outputFilesOnly {
		switch outputType {
		case "xml":
			ext = ".xml"
		case "html", "html5":
			ext = ".html"
		case "xhtml", "xhtml5":
			ext = ".xhtml"
		case "md2adoc":
			ext = ".adoc"
		}
	}

	// Walk the directory and add files to zip
	err = filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Filter files if outputFilesOnly is true
		if outputFilesOnly {
			if !strings.HasSuffix(strings.ToLower(path), ext) {
				return nil // Skip non-output files
			}
		}

		// Get relative path from source directory
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Get file info
		info, err := file.Stat()
		if err != nil {
			return err
		}

		// Create zip file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath
		header.Method = zip.Deflate

		// Write file to zip
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		zipWriter.Close()
		zipFile.Close()
		os.Remove(zipPath)
		return "", fmt.Errorf("failed to add files to zip: %w", err)
	}

	// Close zip writer to flush data
	if err := zipWriter.Close(); err != nil {
		zipFile.Close()
		os.Remove(zipPath)
		return "", fmt.Errorf("failed to close zip writer: %w", err)
	}

	// Close file
	if err := zipFile.Close(); err != nil {
		os.Remove(zipPath)
		return "", fmt.Errorf("failed to close zip file: %w", err)
	}

	return zipPath, nil
}

func (s *Server) handleBatchDownloadZip(w http.ResponseWriter, r *http.Request) {
	sourcePath := r.URL.Query().Get("sourcePath")
	outputType := r.URL.Query().Get("outputType")
	outputFilesOnly := r.URL.Query().Get("outputFilesOnly") == "true"
	zipPath := r.URL.Query().Get("path")

	// If we need to create a new zip with different options
	if sourcePath != "" && outputType != "" {
		// Security check: ensure source path is in temp directory
		if !strings.HasPrefix(sourcePath, os.TempDir()) {
			http.Error(w, "Invalid source path", http.StatusBadRequest)
			return
		}

		var err error
		zipPath, err = s.createResultsZip(sourcePath, outputType, outputFilesOnly)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create zip: %v", err), http.StatusInternalServerError)
			return
		}
	}

	if zipPath == "" {
		http.Error(w, "Path parameter required", http.StatusBadRequest)
		return
	}

	// Security check: ensure path is in temp directory and matches pattern
	if !strings.HasPrefix(zipPath, os.TempDir()) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	if !strings.Contains(filepath.Base(zipPath), "adc-results-") {
		http.Error(w, "Invalid zip file", http.StatusBadRequest)
		return
	}

	// Check if file exists
	info, err := os.Stat(zipPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Open the zip file
	file, err := os.Open(zipPath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set headers for download
	filename := filepath.Base(zipPath)
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

	// Stream the file
	io.Copy(w, file)
}

func (s *Server) handleBatchCleanup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tempDir := os.TempDir()
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": fmt.Sprintf("Failed to read temp directory: %v", err),
		})
		return
	}

	cleanedCount := 0
	var errors []string

	for _, entry := range entries {
		// Remove adc-batch-* directories (equivalent to rm -rf /tmp/adc-batch-*)
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "adc-batch-") {
			dirPath := filepath.Join(tempDir, entry.Name())
			if err := os.RemoveAll(dirPath); err != nil {
				errors = append(errors, fmt.Sprintf("Failed to remove %s: %v", entry.Name(), err))
			} else {
				cleanedCount++
			}
		}
		// Remove adc-results-*.zip files (equivalent to rm /tmp/adc-results-*.zip)
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "adc-results-") && strings.HasSuffix(entry.Name(), ".zip") {
			filePath := filepath.Join(tempDir, entry.Name())
			if err := os.Remove(filePath); err != nil {
				errors = append(errors, fmt.Sprintf("Failed to remove %s: %v", entry.Name(), err))
			} else {
				cleanedCount++
			}
		}
	}

	response := map[string]interface{}{
		"message":      fmt.Sprintf("Cleaned up %d items", cleanedCount),
		"cleanedCount": cleanedCount,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleDefaultTempPath(w http.ResponseWriter, r *http.Request) {
	// Generate a UUID v4
	uuid := make([]byte, 16)
	if _, err := rand.Read(uuid); err != nil {
		http.Error(w, "Failed to generate UUID", http.StatusInternalServerError)
		return
	}
	// Set version (4) and variant bits
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant 10

	// Format as UUID string
	uuidStr := fmt.Sprintf("%s-%s-%s-%s-%s",
		hex.EncodeToString(uuid[0:4]),
		hex.EncodeToString(uuid[4:6]),
		hex.EncodeToString(uuid[6:8]),
		hex.EncodeToString(uuid[8:10]),
		hex.EncodeToString(uuid[10:16]))

	// Get temp directory (works cross-platform: /tmp on Unix, %TEMP% on Windows)
	tempDir := os.TempDir()
	defaultPath := filepath.Join(tempDir, uuidStr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"path": defaultPath,
	})
}

func (s *Server) handleWatcherConvertFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		FilePath   string `json:"filePath"`
		OutputType string `json:"outputType,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid JSON",
		})
		return
	}

	if req.FilePath == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "filePath is required",
		})
		return
	}

	// Default to xml if not specified
	outputType := req.OutputType
	if outputType == "" {
		outputType = "xml"
	}

	// Process the file - check if it's a markdown file for md2adoc conversion
	var err error
	if outputType == "md2adoc" && (strings.HasSuffix(strings.ToLower(req.FilePath), ".md")) {
		err = s.processMarkdownFile(req.FilePath, outputType)
	} else {
		err = s.processAdocFile(req.FilePath, outputType)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   fmt.Sprintf("Failed to process file: %v", err),
			"filePath": req.FilePath,
		})
		return
	}

	// Determine output file path
	baseName := strings.TrimSuffix(req.FilePath, filepath.Ext(req.FilePath))
	var outputFile string
	switch outputType {
	case "md2adoc":
		outputFile = baseName + ".adoc"
	case "xml":
		outputFile = baseName + ".xml"
	case "html", "html5":
		outputFile = baseName + ".html"
	case "xhtml", "xhtml5":
		outputFile = baseName + ".xhtml"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"filePath":   req.FilePath,
		"outputFile": outputFile,
		"outputType": outputType,
	})
}

func (s *Server) handleConfigUpdate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WatchDir string `json:"watchDir"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Read existing config
	config := make(map[string]interface{})
	if data, err := os.ReadFile("adc.json"); err == nil {
		json.Unmarshal(data, &config)
	}

	config["watchDir"] = req.WatchDir

	data, _ := json.MarshalIndent(config, "", "  ")
	if err := os.WriteFile("adc.json", data, 0644); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Configuration updated",
	})
}

func main() {
	port := 8005

	// Check environment variable first, then command line argument
	if portStr := os.Getenv("PORT"); portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("Invalid port number in PORT environment variable: %s", portStr)
		}
	} else if len(os.Args) > 1 {
		var err error
		port, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("Invalid port number: %s", os.Args[1])
		}
	}

	server := NewServer(port)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
