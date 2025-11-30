package lib

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectArchiveFormat(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"test.zip", "zip"},
		{"test.ZIP", "zip"},
		{"test.tar", "tar"},
		{"test.TAR", "tar"},
		{"test.tgz", "tar.gz"},
		{"test.TGZ", "tar.gz"},
		{"test.tar.gz", "tar.gz"},
		{"test.TAR.GZ", "tar.gz"},
		{"test.txt", ""},
		{"test", ""},
		{"archive.tar.gz.backup", ""}, // Only matches .tar.gz at end
	}

	for _, tt := range tests {
		result := DetectArchiveFormat(tt.filename)
		if result != tt.expected {
			t.Errorf("DetectArchiveFormat(%q) = %q, want %q", tt.filename, result, tt.expected)
		}
	}
}

func TestExtractArchive_ZIP(t *testing.T) {
	// Create a temporary ZIP file
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	zipPath := filepath.Join(tempDir, "test.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("Failed to create zip file: %v", err)
	}

	zw := zip.NewWriter(zipFile)
	w, err := zw.Create("test.txt")
	if err != nil {
		zipFile.Close()
		t.Fatalf("Failed to create zip entry: %v", err)
	}
	w.Write([]byte("test content"))
	zw.Close()
	zipFile.Close()

	// Reopen for reading
	zipFile, err = os.Open(zipPath)
	if err != nil {
		t.Fatalf("Failed to reopen zip file: %v", err)
	}
	defer zipFile.Close()

	// Extract
	extractDir := filepath.Join(tempDir, "extracted")
	err = ExtractArchive(zipFile, zipPath, extractDir)
	if err != nil {
		t.Fatalf("ExtractArchive failed: %v", err)
	}

	// Verify extracted file
	extractedFile := filepath.Join(extractDir, "test.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	if string(content) != "test content" {
		t.Errorf("Expected 'test content', got %q", string(content))
	}
}

func TestExtractArchive_TAR(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tarPath := filepath.Join(tempDir, "test.tar")
	tarFile, err := os.Create(tarPath)
	if err != nil {
		t.Fatalf("Failed to create tar file: %v", err)
	}

	tw := tar.NewWriter(tarFile)
	header := &tar.Header{
		Name: "test.txt",
		Size: 12,
		Mode: 0644,
	}
	tw.WriteHeader(header)
	tw.Write([]byte("test content"))
	tw.Close()
	tarFile.Close()

	// Reopen for reading
	tarFile, err = os.Open(tarPath)
	if err != nil {
		t.Fatalf("Failed to reopen tar file: %v", err)
	}
	defer tarFile.Close()

	// Extract
	extractDir := filepath.Join(tempDir, "extracted")
	err = ExtractArchive(tarFile, tarPath, extractDir)
	if err != nil {
		t.Fatalf("ExtractArchive failed: %v", err)
	}

	// Verify extracted file
	extractedFile := filepath.Join(extractDir, "test.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	if string(content) != "test content" {
		t.Errorf("Expected 'test content', got %q", string(content))
	}
}

func TestExtractArchive_TARGZ(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tarGzPath := filepath.Join(tempDir, "test.tar.gz")
	tarGzFile, err := os.Create(tarGzPath)
	if err != nil {
		t.Fatalf("Failed to create tar.gz file: %v", err)
	}

	gw := gzip.NewWriter(tarGzFile)
	tw := tar.NewWriter(gw)
	header := &tar.Header{
		Name: "test.txt",
		Size: 12,
		Mode: 0644,
	}
	tw.WriteHeader(header)
	tw.Write([]byte("test content"))
	tw.Close()
	gw.Close()
	tarGzFile.Close()

	// Reopen for reading
	tarGzFile, err = os.Open(tarGzPath)
	if err != nil {
		t.Fatalf("Failed to reopen tar.gz file: %v", err)
	}
	defer tarGzFile.Close()

	// Extract
	extractDir := filepath.Join(tempDir, "extracted")
	err = ExtractArchive(tarGzFile, tarGzPath, extractDir)
	if err != nil {
		t.Fatalf("ExtractArchive failed: %v", err)
	}

	// Verify extracted file
	extractedFile := filepath.Join(extractDir, "test.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	if string(content) != "test content" {
		t.Errorf("Expected 'test content', got %q", string(content))
	}
}

func TestExtractArchive_PathTraversal(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	zipPath := filepath.Join(tempDir, "malicious.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("Failed to create zip file: %v", err)
	}

	zw := zip.NewWriter(zipFile)
	// Try to create a file outside the extraction directory
	w, err := zw.Create("../outside.txt")
	if err != nil {
		zipFile.Close()
		t.Fatalf("Failed to create zip entry: %v", err)
	}
	w.Write([]byte("malicious"))
	zw.Close()
	zipFile.Close()

	// Reopen for reading
	zipFile, err = os.Open(zipPath)
	if err != nil {
		t.Fatalf("Failed to reopen zip file: %v", err)
	}
	defer zipFile.Close()

	// Extract - should fail with path traversal error
	extractDir := filepath.Join(tempDir, "extracted")
	err = ExtractArchive(zipFile, zipPath, extractDir)
	if err == nil {
		t.Fatal("Expected error for path traversal, got nil")
	}

	if !strings.Contains(err.Error(), "illegal file path") {
		t.Errorf("Expected 'illegal file path' error, got: %v", err)
	}

	// Verify file was not created outside
	outsideFile := filepath.Join(tempDir, "outside.txt")
	if _, err := os.Stat(outsideFile); err == nil {
		t.Error("Malicious file was created outside extraction directory")
	}
}

func TestExtractArchive_NestedDirectories(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	zipPath := filepath.Join(tempDir, "nested.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("Failed to create zip file: %v", err)
	}

	zw := zip.NewWriter(zipFile)
	// Create nested structure
	zw.Create("dir1/")
	w, _ := zw.Create("dir1/file1.txt")
	w.Write([]byte("file1"))
	zw.Create("dir1/dir2/")
	w, _ = zw.Create("dir1/dir2/file2.txt")
	w.Write([]byte("file2"))
	zw.Close()
	zipFile.Close()

	// Reopen for reading
	zipFile, err = os.Open(zipPath)
	if err != nil {
		t.Fatalf("Failed to reopen zip file: %v", err)
	}
	defer zipFile.Close()

	// Extract
	extractDir := filepath.Join(tempDir, "extracted")
	err = ExtractArchive(zipFile, zipPath, extractDir)
	if err != nil {
		t.Fatalf("ExtractArchive failed: %v", err)
	}

	// Verify nested structure
	file1 := filepath.Join(extractDir, "dir1", "file1.txt")
	file2 := filepath.Join(extractDir, "dir1", "dir2", "file2.txt")

	if _, err := os.Stat(file1); err != nil {
		t.Errorf("File1 not found: %v", err)
	}
	if _, err := os.Stat(file2); err != nil {
		t.Errorf("File2 not found: %v", err)
	}
}

func TestExtractArchive_UnsupportedFormat(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	file, err := os.CreateTemp(tempDir, "test.*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	file.Close()

	file, err = os.Open(file.Name())
	if err != nil {
		t.Fatalf("Failed to reopen file: %v", err)
	}
	defer file.Close()

	err = ExtractArchive(file, "test.txt", tempDir)
	if err == nil {
		t.Fatal("Expected error for unsupported format, got nil")
	}

	if !strings.Contains(err.Error(), "unsupported archive format") {
		t.Errorf("Expected 'unsupported archive format' error, got: %v", err)
	}
}

func TestCreateArchive_ZIP(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create source directory with files
	sourceDir := filepath.Join(tempDir, "source")
	os.MkdirAll(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("content1"), 0644)
	os.MkdirAll(filepath.Join(sourceDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(sourceDir, "subdir", "file2.txt"), []byte("content2"), 0644)
	
	// Ensure sourceDir exists and is a directory
	if info, err := os.Stat(sourceDir); err != nil || !info.IsDir() {
		t.Fatalf("Source directory not properly created: %v", err)
	}

	// Create archive
	archivePath := filepath.Join(tempDir, "output.zip")
	err = CreateArchive(sourceDir, "zip", archivePath)
	if err != nil {
		t.Fatalf("CreateArchive failed: %v", err)
	}

	// Verify archive exists
	if _, err := os.Stat(archivePath); err != nil {
		t.Fatalf("Archive file not created: %v", err)
	}

	// Extract and verify
	extractDir := filepath.Join(tempDir, "extracted")
	os.MkdirAll(extractDir, 0755)
	archiveFile, err := os.Open(archivePath)
	if err != nil {
		t.Fatalf("Failed to open archive: %v", err)
	}
	defer archiveFile.Close()

	err = ExtractArchive(archiveFile, archivePath, extractDir)
	if err != nil {
		t.Fatalf("Failed to extract created archive: %v", err)
	}

	// Verify files
	content1, _ := os.ReadFile(filepath.Join(extractDir, "file1.txt"))
	content2, _ := os.ReadFile(filepath.Join(extractDir, "subdir", "file2.txt"))

	if string(content1) != "content1" {
		t.Errorf("Expected 'content1', got %q", string(content1))
	}
	if string(content2) != "content2" {
		t.Errorf("Expected 'content2', got %q", string(content2))
	}
}

func TestCreateArchive_TAR(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "source")
	os.MkdirAll(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "file.txt"), []byte("content"), 0644)

	archivePath := filepath.Join(tempDir, "output.tar")
	err = CreateArchive(sourceDir, "tar", archivePath)
	if err != nil {
		t.Fatalf("CreateArchive failed: %v", err)
	}

	if _, err := os.Stat(archivePath); err != nil {
		t.Fatalf("Archive file not created: %v", err)
	}
}

func TestCreateArchive_TARGZ(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "source")
	os.MkdirAll(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "file.txt"), []byte("content"), 0644)

	archivePath := filepath.Join(tempDir, "output.tar.gz")
	err = CreateArchive(sourceDir, "tar.gz", archivePath)
	if err != nil {
		t.Fatalf("CreateArchive failed: %v", err)
	}

	if _, err := os.Stat(archivePath); err != nil {
		t.Fatalf("Archive file not created: %v", err)
	}
}

func TestCreateArchive_TGZ(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "source")
	os.MkdirAll(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "file.txt"), []byte("content"), 0644)

	archivePath := filepath.Join(tempDir, "output.tgz")
	err = CreateArchive(sourceDir, "tgz", archivePath)
	if err != nil {
		t.Fatalf("CreateArchive failed: %v", err)
	}

	if _, err := os.Stat(archivePath); err != nil {
		t.Fatalf("Archive file not created: %v", err)
	}
}

func TestCreateArchive_UnsupportedFormat(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "source")
	os.MkdirAll(sourceDir, 0755)

	archivePath := filepath.Join(tempDir, "output.rar")
	err = CreateArchive(sourceDir, "rar", archivePath)
	if err == nil {
		t.Fatal("Expected error for unsupported format, got nil")
	}

	if !strings.Contains(err.Error(), "unsupported archive format") {
		t.Errorf("Expected 'unsupported archive format' error, got: %v", err)
	}
}

func TestExtractArchive_EmptyArchive(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create empty ZIP
	zipPath := filepath.Join(tempDir, "empty.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("Failed to create zip file: %v", err)
	}
	zw := zip.NewWriter(zipFile)
	zw.Close()
	zipFile.Close()

	// Reopen for reading
	zipFile, err = os.Open(zipPath)
	if err != nil {
		t.Fatalf("Failed to reopen zip file: %v", err)
	}
	defer zipFile.Close()

	// Extract empty archive - should succeed
	extractDir := filepath.Join(tempDir, "extracted")
	err = ExtractArchive(zipFile, zipPath, extractDir)
	if err != nil {
		t.Fatalf("ExtractArchive should handle empty archive: %v", err)
	}
}

func TestExtractArchive_InvalidZIP(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create invalid ZIP (just random bytes)
	zipPath := filepath.Join(tempDir, "invalid.zip")
	os.WriteFile(zipPath, []byte("not a zip file"), 0644)

	zipFile, err := os.Open(zipPath)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer zipFile.Close()

	extractDir := filepath.Join(tempDir, "extracted")
	err = ExtractArchive(zipFile, zipPath, extractDir)
	if err == nil {
		t.Fatal("Expected error for invalid ZIP, got nil")
	}
}

func TestExtractArchive_InvalidTAR(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create invalid TAR
	tarPath := filepath.Join(tempDir, "invalid.tar")
	os.WriteFile(tarPath, []byte("not a tar file"), 0644)

	tarFile, err := os.Open(tarPath)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer tarFile.Close()

	extractDir := filepath.Join(tempDir, "extracted")
	err = ExtractArchive(tarFile, tarPath, extractDir)
	if err == nil {
		t.Fatal("Expected error for invalid TAR, got nil")
	}
}

func TestExtractArchive_InvalidTARGZ(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create invalid TAR.GZ
	tarGzPath := filepath.Join(tempDir, "invalid.tar.gz")
	os.WriteFile(tarGzPath, []byte("not a tar.gz file"), 0644)

	tarGzFile, err := os.Open(tarGzPath)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer tarGzFile.Close()

	extractDir := filepath.Join(tempDir, "extracted")
	err = ExtractArchive(tarGzFile, tarGzPath, extractDir)
	if err == nil {
		t.Fatal("Expected error for invalid TAR.GZ, got nil")
	}
}

func TestExtractArchive_RelativePathTraversal(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	zipPath := filepath.Join(tempDir, "traversal.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("Failed to create zip file: %v", err)
	}

	zw := zip.NewWriter(zipFile)
	// Try relative path traversal
	w, err := zw.Create("../../outside.txt")
	if err != nil {
		zipFile.Close()
		t.Fatalf("Failed to create zip entry: %v", err)
	}
	w.Write([]byte("malicious"))
	zw.Close()
	zipFile.Close()

	zipFile, err = os.Open(zipPath)
	if err != nil {
		t.Fatalf("Failed to reopen zip file: %v", err)
	}
	defer zipFile.Close()

	extractDir := filepath.Join(tempDir, "extracted")
	err = ExtractArchive(zipFile, zipPath, extractDir)
	if err == nil {
		t.Fatal("Expected error for path traversal, got nil")
	}
}

func TestCreateArchive_EmptyDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "source")
	os.MkdirAll(sourceDir, 0755)

	archivePath := filepath.Join(tempDir, "empty.zip")
	err = CreateArchive(sourceDir, "zip", archivePath)
	if err != nil {
		t.Fatalf("CreateArchive should handle empty directory: %v", err)
	}

	if _, err := os.Stat(archivePath); err != nil {
		t.Fatalf("Archive file should be created even for empty directory: %v", err)
	}
}

func TestExtractArchive_FilePermissions(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "archive-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	zipPath := filepath.Join(tempDir, "perms.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("Failed to create zip file: %v", err)
	}

	zw := zip.NewWriter(zipFile)
	w, err := zw.Create("executable.sh")
	if err != nil {
		zipFile.Close()
		t.Fatalf("Failed to create zip entry: %v", err)
	}
	w.Write([]byte("#!/bin/sh"))
	zw.Close()
	zipFile.Close()

	zipFile, err = os.Open(zipPath)
	if err != nil {
		t.Fatalf("Failed to reopen zip file: %v", err)
	}
	defer zipFile.Close()

	extractDir := filepath.Join(tempDir, "extracted")
	err = ExtractArchive(zipFile, zipPath, extractDir)
	if err != nil {
		t.Fatalf("ExtractArchive failed: %v", err)
	}

	// Verify file was extracted
	extractedFile := filepath.Join(extractDir, "executable.sh")
	if _, err := os.Stat(extractedFile); err != nil {
		t.Fatalf("Extracted file not found: %v", err)
	}
}

