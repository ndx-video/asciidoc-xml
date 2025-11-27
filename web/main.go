package main

import (
	"archive/zip"
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
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
	mux.HandleFunc("/api/convert", s.handleConvert)
	mux.HandleFunc("/api/validate", s.handleValidate)
	mux.HandleFunc("/api/xslt", s.handleXSLT)
	mux.HandleFunc("/api/upload", s.handleUpload)
	mux.HandleFunc("/api/load-file", s.handleLoadFile)
	mux.HandleFunc("/api/files", s.handleFiles)
	mux.HandleFunc("/api/batch/upload-zip", s.handleBatchUploadZip)
	mux.HandleFunc("/api/batch/process-folder", s.handleBatchProcessFolder)
	mux.HandleFunc("/api/config/update", s.handleConfigUpdate)

	// Documentation
	mux.HandleFunc("/docs", s.handleDocs)
	mux.HandleFunc("/batch", s.handleBatch)

	// Serve main page
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/browse", s.handleBrowse)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting server on http://localhost%s", addr)
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

	switch outputType {
	case "html", "html5":
		output, err = lib.ConvertToHTML(bytes.NewReader([]byte(req.AsciiDoc)), false)
		if err != nil {
			http.Error(w, fmt.Sprintf("Conversion failed: %v", err), http.StatusInternalServerError)
			return
		}
		contentType = "text/html; charset=utf-8"

	case "xhtml", "xhtml5":
		output, err = lib.ConvertToHTML(bytes.NewReader([]byte(req.AsciiDoc)), true)
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

	default:
		http.Error(w, fmt.Sprintf("Invalid output type: %s. Supported types: xml, html, xhtml", outputType), http.StatusBadRequest)
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

func (s *Server) handleDocs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the AsciiDoc documentation file
	var docPath string
	var docContent []byte
	var err error

	// Check if requesting specific page or default
	page := r.URL.Query().Get("page")
	if page == "" || page == "main" {
		page = "index.adoc"
	}

	// Sanitize page path to prevent directory traversal
	page = filepath.Base(page)
	if !strings.HasSuffix(page, ".adoc") {
		http.Error(w, "Invalid page format", http.StatusBadRequest)
		return
	}

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
		http.Error(w, fmt.Sprintf("Failed to read documentation: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert AsciiDoc directly to HTML5
	htmlContent, err := lib.ConvertToHTML(bytes.NewReader(docContent), false)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to convert documentation: %v", err), http.StatusInternalServerError)
		return
	}

	// Inject navigation and styling into the HTML
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>AsciiDoc XML Converter - Documentation</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 2rem;
            line-height: 1.6;
        }
        .navbar {
            background: #2c3e50;
            color: white;
            padding: 0.75rem 1rem;
            margin: -2rem -2rem 2rem -2rem;
            display: flex;
            align-items: center;
            gap: 1rem;
        }
        .navbar a {
            color: white;
            text-decoration: none;
            padding: 0.5rem 1rem;
            border-radius: 4px;
            transition: background 0.2s;
        }
        .navbar a:hover {
            background: rgba(255, 255, 255, 0.1);
        }
        h1, h2, h3, h4, h5, h6 {
            color: #2c3e50;
            margin-top: 2rem;
            margin-bottom: 1rem;
        }
        code {
            background: #f4f4f4;
            padding: 0.2em 0.4em;
            border-radius: 3px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            font-size: 0.9em;
        }
        pre {
            background: #f4f4f4;
            padding: 1rem;
            border-radius: 4px;
            overflow-x: auto;
        }
        pre code {
            background: none;
            padding: 0;
        }
        table {
            border-collapse: collapse;
            width: 100%;
            margin: 1rem 0;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 0.5rem;
            text-align: left;
        }
        th {
            background: #f4f4f4;
            font-weight: 600;
        }
        .admonition {
            border-left: 4px solid #3498db;
            padding: 1rem;
            margin: 1rem 0;
            background: #f8f9fa;
        }
        .admonition.admonition-warning {
            border-left-color: #f39c12;
        }
        .admonition.admonition-important {
            border-left-color: #e74c3c;
        }
        .sidebar {
            background: #f8f9fa;
            border: 1px solid #dee2e6;
            padding: 1rem;
            margin: 1rem 0;
            border-radius: 4px;
        }
        .example {
            background: #fff;
            border: 1px solid #dee2e6;
            padding: 1rem;
            margin: 1rem 0;
            border-radius: 4px;
        }
    </style>
</head>
<body>
    <nav class="navbar">
        <h1 style="margin: 0; font-size: 1.25rem;">AsciiDoc XML Converter</h1>
        <a href="/">Home</a>
        <a href="/docs">Docs</a>
        <a href="/batch">Batching</a>
        <a href="/browse">Browse</a>
    </nav>`

	// Insert the generated HTML content (skip the DOCTYPE and html/head tags)
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

	html += bodyContent.String()
	html += `</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
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
	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Assuming adc is in path. If not, this might fail.
	// We try to find 'adc' executable.
	cmd := exec.Command("adc", req.Path)
	if err := cmd.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to start batch process: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Batch processing started for " + req.Path,
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
