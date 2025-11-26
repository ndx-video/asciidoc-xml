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

	"github.com/bytesparadise/libasciidoc"
	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"asciidoc-xml/converter"
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
		AsciiDoc string `json:"asciidoc"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	xml, err := converter.ConvertToXML(bytes.NewReader([]byte(req.AsciiDoc)))
	if err != nil {
		http.Error(w, fmt.Sprintf("Conversion failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"xml": xml,
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

	// Use libasciidoc to validate by attempting to convert
	// We'll use a null writer to just validate without generating output
	var buf bytes.Buffer
	config := configuration.NewConfiguration(configuration.WithFilename(""))
	_, err = libasciidoc.Convert(strings.NewReader(req.AsciiDoc), &buf, config)
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
	if strings.HasPrefix(path, "static/") {
		path = strings.TrimPrefix(path, "static/")
	}

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

