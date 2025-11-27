package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	// Documentation
	mux.HandleFunc("/docs", s.handleDocs)

	// Serve main page
	mux.HandleFunc("/", s.handleIndex)

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
		AsciiDoc  string `json:"asciidoc"`
		OutputType string `json:"output,omitempty"` // "xml", "html", "xhtml"
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

	// Remove leading slash and ensure it's in static directory
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimPrefix(path, "static/")

	// Try embedded static files first
	var content []byte
	var err error
	
	// Try embedded files
	embeddedPath := filepath.Join("static", path)
	content, err = staticFiles.ReadFile(embeddedPath)
	if err != nil {
		// Try filesystem locations
		locations := []string{
			filepath.Join("..", "web", "static", path),
			filepath.Join("static", path),
			filepath.Join("..", "static", path),
			path,
		}

		for _, loc := range locations {
			content, err = os.ReadFile(loc)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read file: %v", err), http.StatusNotFound)
		return
	}

	// Determine content type
	contentType := "text/plain"
	if strings.HasSuffix(strings.ToLower(path), ".adoc") || strings.HasSuffix(strings.ToLower(path), ".asciidoc") {
		contentType = "text/plain"
	} else if strings.HasSuffix(strings.ToLower(path), ".xsl") || strings.HasSuffix(strings.ToLower(path), ".xslt") {
		contentType = "application/xml"
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(content)
}

func (s *Server) handleDocs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the AsciiDoc documentation file
	docPath := filepath.Join("..", "docs", "asciidoc-xml.adoc")
	docContent, err := os.ReadFile(docPath)
	if err != nil {
		// Try current directory
		docPath = filepath.Join("docs", "asciidoc-xml.adoc")
		docContent, err = os.ReadFile(docPath)
		if err != nil {
			// Try parent directory
			docPath = filepath.Join("..", "..", "docs", "asciidoc-xml.adoc")
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

