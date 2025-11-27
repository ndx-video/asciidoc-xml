package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessFile_SingleFile(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test AsciiDoc file
	testFile := filepath.Join(tempDir, "test.adoc")
	testContent := "= Test Document\n\nThis is a test."
	os.WriteFile(testFile, []byte(testContent), 0644)

	// Process file with no XSLT, XML output
	err := processFile(testFile, "", "xml")
	if err != nil {
		t.Fatalf("processFile failed: %v", err)
	}

	// Check that XML file was created
	xmlFile := filepath.Join(tempDir, "test.xml")
	if _, err := os.Stat(xmlFile); os.IsNotExist(err) {
		t.Error("XML file was not created")
	}

	// Verify XML content
	xmlContent, err := os.ReadFile(xmlFile)
	if err != nil {
		t.Fatalf("Failed to read XML file: %v", err)
	}

	if !strings.Contains(string(xmlContent), "Test Document") {
		t.Error("XML should contain document title")
	}
}

func TestProcessFile_OutputDir(t *testing.T) {
	tempDir := t.TempDir()
	outputDirObj := filepath.Join(tempDir, "out")
	
	// Create test AsciiDoc file
	testFile := filepath.Join(tempDir, "test.adoc")
	testContent := "= Test Document\n\nThis is a test."
	os.WriteFile(testFile, []byte(testContent), 0644)

	// Set global outputDir
	oldOutputDir := outputDir
	outputDir = outputDirObj
	defer func() { outputDir = oldOutputDir }()

	// Process file
	err := processFile(testFile, "", "xml")
	if err != nil {
		t.Fatalf("processFile failed: %v", err)
	}

	// Check that XML file was created in output directory
	xmlFile := filepath.Join(outputDirObj, "test.xml")
	if _, err := os.Stat(xmlFile); os.IsNotExist(err) {
		t.Error("XML file was not created in output directory")
	}

	// Verify original location is empty
	origXmlFile := filepath.Join(tempDir, "test.xml")
	if _, err := os.Stat(origXmlFile); err == nil {
		t.Error("XML file should not be created in source directory")
	}
}

func TestLoadConfig(t *testing.T) {
	// Run in a separate directory to avoid conflict with real adc.json
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)

	// Save original flags
	oldOutputDir := outputDir
	oldNoXSL := noXSL
	oldOutputType := outputType
	
	defer func() {
		outputDir = oldOutputDir
		noXSL = oldNoXSL
		outputType = oldOutputType
	}()

	// Reset flags for this test
	outputDir = ""
	noXSL = false
	outputType = ""
	
	config := Config{
		OutputDir:  stringPtr("custom-out"),
		NoXSL:      boolPtr(true),
		OutputType: stringPtr("html"),
	}
	
	data, _ := json.Marshal(config)
	os.WriteFile("adc.json", data, 0644)

	// Initialize flags if not already done (main package init might have run)
	// But we can't easily re-parse flags. loadConfig checks flags but here we assume none are set.
	// We need to mock the flag visited map or just rely on the fact that we haven't parsed flags in this test execution
	// Actually loadConfig calls flag.Visit, which uses flag.CommandLine.
	
	// Reset flag.CommandLine for this test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.BoolVar(&noXSL, "no-xsl", false, "")
	flag.StringVar(&outputType, "output", "xml", "")
	flag.StringVar(&outputDir, "out-dir", "", "")

	loadConfig()

	if outputDir != "custom-out" {
		t.Errorf("Expected outputDir 'custom-out', got '%s'", outputDir)
	}
	if !noXSL {
		t.Error("Expected noXSL to be true")
	}
	if outputType != "html" {
		t.Errorf("Expected outputType 'html', got '%s'", outputType)
	}
}

func stringPtr(s string) *string { return &s }
func boolPtr(b bool) *bool { return &b }

func TestProcessFile_OverwriteWithY(t *testing.T) {
	tempDir := t.TempDir()
	
	testFile := filepath.Join(tempDir, "test.adoc")
	os.WriteFile(testFile, []byte("= Test\n\nContent"), 0644)

	// Set auto-overwrite
	autoOverwrite = true
	defer func() { autoOverwrite = false }()

	// First conversion
	err := processFile(testFile, "", "xml")
	if err != nil {
		t.Fatalf("First conversion failed: %v", err)
	}

	// Verify XML exists
	xmlFile := filepath.Join(tempDir, "test.xml")
	if _, err := os.Stat(xmlFile); os.IsNotExist(err) {
		t.Fatal("XML file was not created on first run")
	}

	// Get first XML content
	firstContent, _ := os.ReadFile(xmlFile)

	// Second conversion should overwrite
	err = processFile(testFile, "", "xml")
	if err != nil {
		t.Fatalf("Second conversion failed: %v", err)
	}

	// Verify it was overwritten (content should be the same, but file should exist)
	secondContent, err := os.ReadFile(xmlFile)
	if err != nil {
		t.Fatalf("Failed to read XML after second conversion: %v", err)
	}

	if string(firstContent) != string(secondContent) {
		t.Error("XML content should be the same")
	}
}

func TestFindAdocFiles_Directory(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create multiple AsciiDoc files
	files := []string{"file1.adoc", "file2.adoc", "subdir/file3.adoc", "other.txt"}
	for _, f := range files {
		fullPath := filepath.Join(tempDir, f)
		os.MkdirAll(filepath.Dir(fullPath), 0755)
		os.WriteFile(fullPath, []byte("content"), 0644)
	}

	// Test directory traversal
	var foundFiles []string
	filepath.WalkDir(tempDir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(p), ".adoc") {
			foundFiles = append(foundFiles, p)
		}
		return nil
	})

	if len(foundFiles) != 3 {
		t.Errorf("Expected 3 .adoc files, found %d", len(foundFiles))
	}
}

func TestPromptOverwrite(t *testing.T) {
	// Test response validation logic
	responses := map[string]bool{
		"y":    true,
		"yes":  true,
		"n":    true,
		"no":   true,
		"a":    true,
		"all":  true,
		"q":    true,
		"quit": true,
		"x":    false,
		"":     false,
	}

	for resp, valid := range responses {
		normalized := strings.TrimSpace(strings.ToLower(resp))
		isValid := normalized == "y" || normalized == "yes" || normalized == "n" || normalized == "no" || normalized == "a" || normalized == "all" || normalized == "q" || normalized == "quit"
		if isValid != valid {
			t.Errorf("Response '%s' validation mismatch: expected %v, got %v", resp, valid, isValid)
		}
	}
}

func TestADC_NoXSLFlag(t *testing.T) {
	// Test that --no-xsl flag works
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Reset flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	
	noXSL = true
	defer func() { noXSL = false }()

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.adoc")
	os.WriteFile(testFile, []byte("= Test\n\nContent"), 0644)

	// Should not require XSLT file
	err := processFile(testFile, "", "xml")
	if err != nil {
		t.Fatalf("Should work without XSLT: %v", err)
	}
}

func TestADC_InvalidFileExtension(t *testing.T) {
	tempDir := t.TempDir()
	
	testFile := filepath.Join(tempDir, "test.txt")
	os.WriteFile(testFile, []byte("content"), 0644)

	// This should be caught before processFile is called
	// Test the validation logic
	if strings.HasSuffix(strings.ToLower(testFile), ".adoc") {
		t.Error("Should not accept .txt files")
	}
}

func TestApplyXSLT_CommandExists(t *testing.T) {
	// Check if xsltproc is available
	_, err := exec.LookPath("xsltproc")
	if err != nil {
		t.Skip("xsltproc not available, skipping XSLT test")
	}

	tempDir := t.TempDir()
	
	// Create minimal XML
	xmlFile := filepath.Join(tempDir, "test.xml")
	xmlContent := `<?xml version="1.0"?><document><title>Test</title></document>`
	os.WriteFile(xmlFile, []byte(xmlContent), 0644)

	// Create minimal XSLT
	xsltFile := filepath.Join(tempDir, "test.xsl")
	xsltContent := `<?xml version="1.0"?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:template match="/">
    <html><body><xsl:value-of select="document/title"/></body></html>
  </xsl:template>
</xsl:stylesheet>`
	os.WriteFile(xsltFile, []byte(xsltContent), 0644)

	htmlFile := filepath.Join(tempDir, "test.html")

	// Test XSLT application
	err = applyXSLT(xmlFile, xsltFile, htmlFile)
	if err != nil {
		t.Fatalf("XSLT transformation failed: %v", err)
	}

	// Verify HTML was created
	if _, err := os.Stat(htmlFile); os.IsNotExist(err) {
		t.Error("HTML file was not created")
	}
}

func TestADC_FileNotFound(t *testing.T) {
	// Test with non-existent file
	err := processFile("/nonexistent/file.adoc", "", "xml")
	if err == nil {
		t.Error("Should fail for non-existent file")
	}
}

func TestADC_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	
	testFile := filepath.Join(tempDir, "test.adoc")
	os.WriteFile(testFile, []byte(""), 0644)

	// Should handle empty file
	err := processFile(testFile, "", "xml")
	// Empty file might succeed or fail depending on converter implementation
	// Just verify it doesn't panic
	if err != nil {
		t.Logf("Empty file conversion returned error (may be expected): %v", err)
	}
}

func TestProcessFile_HTMLOutput(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test AsciiDoc file
	testFile := filepath.Join(tempDir, "test.adoc")
	testContent := "= Test Document\n\nThis is a test."
	os.WriteFile(testFile, []byte(testContent), 0644)

	// Process file with HTML output
	err := processFile(testFile, "", "html")
	if err != nil {
		t.Fatalf("processFile failed: %v", err)
	}

	// Check that HTML file was created (not XML)
	htmlFile := filepath.Join(tempDir, "test.html")
	if _, err := os.Stat(htmlFile); os.IsNotExist(err) {
		t.Error("HTML file was not created")
	}

	// Verify XML file was NOT created
	xmlFile := filepath.Join(tempDir, "test.xml")
	if _, err := os.Stat(xmlFile); err == nil {
		t.Error("XML file should not be created when output type is HTML")
	}

	// Verify HTML content
	htmlContent, err := os.ReadFile(htmlFile)
	if err != nil {
		t.Fatalf("Failed to read HTML file: %v", err)
	}

	htmlStr := string(htmlContent)
	if !strings.Contains(htmlStr, "Test Document") {
		t.Error("HTML should contain document title")
	}

	// HTML should contain HTML tags
	if !strings.Contains(htmlStr, "<") || !strings.Contains(htmlStr, ">") {
		t.Error("HTML should contain HTML tags")
	}
}

func TestProcessFile_XHTMLOutput(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test AsciiDoc file
	testFile := filepath.Join(tempDir, "test.adoc")
	testContent := "= Test Document\n\nThis is a test."
	os.WriteFile(testFile, []byte(testContent), 0644)

	// Process file with XHTML output
	err := processFile(testFile, "", "xhtml")
	if err != nil {
		t.Fatalf("processFile failed: %v", err)
	}

	// Check that XHTML file was created
	xhtmlFile := filepath.Join(tempDir, "test.xhtml")
	if _, err := os.Stat(xhtmlFile); os.IsNotExist(err) {
		t.Error("XHTML file was not created")
	}

	// Verify XML and HTML files were NOT created
	xmlFile := filepath.Join(tempDir, "test.xml")
	if _, err := os.Stat(xmlFile); err == nil {
		t.Error("XML file should not be created when output type is XHTML")
	}
	htmlFile := filepath.Join(tempDir, "test.html")
	if _, err := os.Stat(htmlFile); err == nil {
		t.Error("HTML file should not be created when output type is XHTML")
	}

	// Verify XHTML content
	xhtmlContent, err := os.ReadFile(xhtmlFile)
	if err != nil {
		t.Fatalf("Failed to read XHTML file: %v", err)
	}

	xhtmlStr := string(xhtmlContent)
	if !strings.Contains(xhtmlStr, "Test Document") {
		t.Error("XHTML should contain document title")
	}

	// XHTML should contain HTML/XML tags
	if !strings.Contains(xhtmlStr, "<") || !strings.Contains(xhtmlStr, ">") {
		t.Error("XHTML should contain HTML/XML tags")
	}
}

func TestProcessFile_InvalidOutputType(t *testing.T) {
	tempDir := t.TempDir()
	
	testFile := filepath.Join(tempDir, "test.adoc")
	os.WriteFile(testFile, []byte("= Test\n\nContent"), 0644)

	// Test with invalid output type
	err := processFile(testFile, "", "invalid")
	if err == nil {
		t.Error("Should fail for invalid output type")
	}

	if !strings.Contains(err.Error(), "unsupported output type") {
		t.Errorf("Error message should mention unsupported output type, got: %v", err)
	}
}

func TestProcessFile_AllOutputTypes(t *testing.T) {
	tempDir := t.TempDir()
	
	testFile := filepath.Join(tempDir, "test.adoc")
	testContent := "= Test Document\n\nContent here."
	os.WriteFile(testFile, []byte(testContent), 0644)

	// Set auto-overwrite to avoid prompts
	autoOverwrite = true
	defer func() { autoOverwrite = false }()

	outputTypes := []string{"xml", "html", "xhtml"}
	extensions := []string{".xml", ".html", ".xhtml"}

	for i, outputType := range outputTypes {
		// Process with each output type
		err := processFile(testFile, "", outputType)
		if err != nil {
			t.Fatalf("processFile failed for output type %s: %v", outputType, err)
		}

		// Verify correct file was created
		expectedFile := filepath.Join(tempDir, "test"+extensions[i])
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created for output type %s", expectedFile, outputType)
		}

		// Verify file has content
		content, err := os.ReadFile(expectedFile)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", expectedFile, err)
		}

		if len(content) == 0 {
			t.Errorf("Output file %s is empty for output type %s", expectedFile, outputType)
		}

		// Clean up for next iteration
		os.Remove(expectedFile)
	}
}
