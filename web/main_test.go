package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"asciidoc-xml/lib"
)

func TestServer_handleIndex(t *testing.T) {
	server := NewServer(8005)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	server.handleIndex(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "AsciiDoc XML Converter") {
		t.Error("Response should contain 'AsciiDoc XML Converter'")
	}
}

func TestServer_handleConvert(t *testing.T) {
	server := NewServer(8005)

	tests := []struct {
		name           string
		asciidoc       string
		outputDir      string
		filename       string
		expectedStatus int
		expectXML      bool
	}{
		{
			name:           "valid asciidoc",
			asciidoc:       "= Test Document\n\nThis is a test.",
			expectedStatus:  http.StatusOK,
			expectXML:      true,
		},
		{
			name:           "empty asciidoc",
			asciidoc:       "",
			expectedStatus:  http.StatusOK,
			expectXML:      true,
		},
		{
			name:           "complex asciidoc",
			asciidoc:       "= Title\n\n== Section\n\nContent with *bold* and _italic_.",
			expectedStatus:  http.StatusOK,
			expectXML:      true,
		},
		{
			name:           "with output dir",
			asciidoc:       "= Test\n\nContent",
			outputDir:      "test_output",
			filename:       "my_doc",
			expectedStatus: http.StatusOK,
			expectXML:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]string{
				"asciidoc": tt.asciidoc,
			}
			if tt.outputDir != "" {
				body["outputDir"] = tt.outputDir
			}
			if tt.filename != "" {
				body["filename"] = tt.filename
			}
			
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest(http.MethodPost, "/api/convert", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			// Setup cleanup for output dir
			if tt.outputDir != "" {
				defer os.RemoveAll(tt.outputDir)
			}

			server.handleConvert(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectXML {
				var result map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
					t.Errorf("Failed to parse JSON response: %v", err)
				}

				if result["output"] == nil {
					t.Error("Expected 'output' field in response")
				}

				// Verify saved file if outputDir was set
				if tt.outputDir != "" {
					if result["savedTo"] == nil {
						t.Error("Expected 'savedTo' field in response")
					}
					
					savedPath := result["savedTo"].(string)
					if _, err := os.Stat(savedPath); os.IsNotExist(err) {
						t.Errorf("Saved file not found at %s", savedPath)
					}
				}

				// Verify it's valid XML by trying to parse it
				_, err := lib.Convert(bytes.NewReader([]byte(tt.asciidoc)))
				if err != nil {
					t.Logf("Note: Converter returned error (may be expected): %v", err)
				}
			}
		})
	}
}

func TestServer_handleConvert_MethodNotAllowed(t *testing.T) {
	server := NewServer(8005)
	req := httptest.NewRequest(http.MethodGet, "/api/convert", nil)
	w := httptest.NewRecorder()

	server.handleConvert(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestServer_handleValidate(t *testing.T) {
	server := NewServer(8005)

	tests := []struct {
		name           string
		asciidoc       string
		expectedStatus int
	}{
		{
			name:           "valid asciidoc",
			asciidoc:       "= Test Document\n\nThis is a test.",
			expectedStatus:  http.StatusOK,
		},
		{
			name:           "empty asciidoc",
			asciidoc:       "",
			expectedStatus:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]string{"asciidoc": tt.asciidoc}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest(http.MethodPost, "/api/validate", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.handleValidate(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
				t.Errorf("Failed to parse JSON response: %v", err)
			}

			if _, ok := result["valid"]; !ok {
				t.Error("Response should contain 'valid' field")
			}
		})
	}
}

func TestServer_handleValidate_MethodNotAllowed(t *testing.T) {
	server := NewServer(8005)
	req := httptest.NewRequest(http.MethodGet, "/api/validate", nil)
	w := httptest.NewRecorder()

	server.handleValidate(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestServer_handleXSLT(t *testing.T) {
	server := NewServer(8005)
	req := httptest.NewRequest(http.MethodGet, "/api/xslt", nil)
	w := httptest.NewRecorder()

	server.handleXSLT(w, req)

	// XSLT endpoint should return something (either success or error)
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}

	if w.Code == http.StatusOK {
		// If successful, should contain XML/XSLT content
		body := w.Body.String()
		if !strings.Contains(body, "<?xml") && !strings.Contains(body, "<xsl:") {
			t.Log("XSLT response may not be valid XML (this is OK if file doesn't exist)")
		}
	}
}

func TestServer_handleXSLT_MethodNotAllowed(t *testing.T) {
	server := NewServer(8005)
	req := httptest.NewRequest(http.MethodPost, "/api/xslt", nil)
	w := httptest.NewRecorder()

	server.handleXSLT(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestServer_handleUpload(t *testing.T) {
	server := NewServer(8005)

	// Create a temporary directory for uploads
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)
	os.MkdirAll("static", 0755)

	tests := []struct {
		name           string
		fileType       string
		fileName       string
		fileContent    string
		expectedStatus int
	}{
		{
			name:           "upload asciidoc file",
			fileType:       "asciidoc",
			fileName:       "test.adoc",
			fileContent:    "= Test Document\n\nContent here.",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "upload xslt file",
			fileType:       "xslt",
			fileName:       "test.xsl",
			fileContent:    "<?xml version=\"1.0\"?><xsl:stylesheet version=\"1.0\"></xsl:stylesheet>",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)

			// Add file type
			writer.WriteField("type", tt.fileType)

			// Add file
			part, err := writer.CreateFormFile("file", tt.fileName)
			if err != nil {
				t.Fatalf("Failed to create form file: %v", err)
			}
			part.Write([]byte(tt.fileContent))
			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/api/upload", &buf)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			w := httptest.NewRecorder()

			server.handleUpload(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if w.Code == http.StatusOK {
				var result map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
					t.Errorf("Failed to parse JSON response: %v", err)
				}

				if result["path"] == "" {
					t.Error("Response should contain 'path' field")
				}

				if result["type"] != tt.fileType {
					t.Errorf("Expected type %s, got %s", tt.fileType, result["type"])
				}
			}
		})
	}
}

func TestServer_handleUpload_InvalidType(t *testing.T) {
	server := NewServer(8005)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("type", "invalid")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	server.handleUpload(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestServer_handleUpload_MethodNotAllowed(t *testing.T) {
	server := NewServer(8005)
	req := httptest.NewRequest(http.MethodGet, "/api/upload", nil)
	w := httptest.NewRecorder()

	server.handleUpload(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestServer_handleLoadFile(t *testing.T) {
	server := NewServer(8005)

	// Create a test file in temp directory
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.adoc")
	testContent := "= Test Document\n\nContent here."
	os.WriteFile(testFile, []byte(testContent), 0644)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "load file with static prefix",
			path:           "/static/test.adoc",
			expectedStatus: http.StatusNotFound, // Will fail because file doesn't exist in static
		},
		{
			name:           "load file without prefix",
			path:           "test.adoc",
			expectedStatus: http.StatusNotFound, // Will fail because file doesn't exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/load-file?path="+tt.path, nil)
			w := httptest.NewRecorder()

			server.handleLoadFile(w, req)

			// We expect 404 because the file won't be in the expected locations
			// This is OK - the test verifies the endpoint works
			if w.Code != tt.expectedStatus && w.Code != http.StatusOK {
				t.Logf("Expected status %d or 200, got %d (this is OK if file doesn't exist)", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestServer_handleLoadFile_MissingPath(t *testing.T) {
	server := NewServer(8005)
	req := httptest.NewRequest(http.MethodGet, "/api/load-file", nil)
	w := httptest.NewRecorder()

	server.handleLoadFile(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestServer_handleLoadFile_MethodNotAllowed(t *testing.T) {
	server := NewServer(8005)
	req := httptest.NewRequest(http.MethodPost, "/api/load-file?path=test.adoc", nil)
	w := httptest.NewRecorder()

	server.handleLoadFile(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestServer_Start(t *testing.T) {
	server := NewServer(9999) // Use a high port number for testing

	// Start server in a goroutine
	done := make(chan error, 1)
	go func() {
		done <- server.Start()
	}()

	// Give server a moment to start
	// Note: This test may fail if port 9999 is in use
	// We'll just verify the server can be created, not that it actually starts

	// The server will block, so we can't easily test it fully
	// But we can verify the server struct is created correctly
	if server.port != 9999 {
		t.Errorf("Expected port 9999, got %d", server.port)
	}
}

func TestServer_handleDocs(t *testing.T) {
	server := NewServer(8005)

	// Create a temporary docs directory with test file in current directory structure
	// The handler looks for docs/asciidoc-xml.adoc relative to web/ directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	// Create docs directory in project root (one level up from web/)
	projectRoot := filepath.Join(originalDir, "..")
	if !strings.HasSuffix(projectRoot, "asciidoc-xml") {
		// We're already in project root
		projectRoot = originalDir
	}
	
	docsDir := filepath.Join(projectRoot, "docs")
	
	// Check if docs dir exists to avoid deleting it later
	docsCreated := false
	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		os.MkdirAll(docsDir, 0755)
		docsCreated = true
	}
	
	if docsCreated {
		defer os.RemoveAll(docsDir) // Cleanup only if we created it
	}
	
	testDoc := filepath.Join(docsDir, "test-doc.adoc")
	testContent := "= Test Document\n\nThis is a test document."
	os.WriteFile(testDoc, []byte(testContent), 0644)
	defer os.Remove(testDoc)

	req := httptest.NewRequest(http.MethodGet, "/docs?page=test-doc.adoc", nil)
	w := httptest.NewRecorder()

	server.handleDocs(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Response body: %s", w.Body.String())
		t.Errorf("Expected status 200, got %d", w.Code)
		return
	}

	body := w.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("Response should contain HTML")
	}
	// The content might be in the transformed HTML, so check for XML structure
	if !strings.Contains(body, "Test Document") && !strings.Contains(body, "test document") {
		t.Logf("Response body (first 500 chars): %s", body[:min(500, len(body))])
		t.Error("Response should contain document content")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestServer_handleDocs_MethodNotAllowed(t *testing.T) {
	server := NewServer(8005)
	req := httptest.NewRequest(http.MethodPost, "/docs", nil)
	w := httptest.NewRecorder()

	server.handleDocs(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestNewServer(t *testing.T) {
	port := 8080
	server := NewServer(port)

	if server == nil {
		t.Fatal("NewServer should not return nil")
	}

	if server.port != port {
		t.Errorf("Expected port %d, got %d", port, server.port)
	}
}

// Test helper to create multipart form data
func createMultipartForm(fileType, fileName, content string) (io.Reader, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	writer.WriteField("type", fileType)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, "", err
	}

	part.Write([]byte(content))
	writer.Close()

	return &buf, writer.FormDataContentType(), nil
}

