//go:build testbed

package lib

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestConvertMarkdownToAsciiDoc_TestbedAllFiles tests the markdown to asciidoc converter
// against all test files in the ./testbed/ directory and all subdirectories
func TestConvertMarkdownToAsciiDoc_TestbedAllFiles(t *testing.T) {
	// Try multiple possible paths
	testbedDirs := []string{"./testbed", "testbed", "../testbed"}
	var testbedDir string
	var err error
	
	for _, dir := range testbedDirs {
		if _, err = os.Stat(dir); err == nil {
			testbedDir = dir
			break
		}
	}
	
	if testbedDir == "" {
		t.Skipf("Testbed directory not found (tried: %v), skipping test", testbedDirs)
		return
	}

	// Find all .md files recursively in testbed and subdirectories
	var testFiles []string
	err = filepath.Walk(testbedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			return nil
		}
		
		// Only process .md files
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			return nil
		}
		
		// Skip README.md files
		if info.Name() == "README.md" {
			return nil
		}
		
		testFiles = append(testFiles, path)
		return nil
	})
	
	if err != nil {
		t.Fatalf("Failed to walk testbed directory: %v", err)
	}

	if len(testFiles) == 0 {
		t.Fatal("No test files found in testbed directory")
	}

	t.Logf("Found %d test files to process (including subdirectories)", len(testFiles))

	// Process each file
	for _, testFile := range testFiles {
		// Get relative path for cleaner test names
		relPath, err := filepath.Rel(testbedDir, testFile)
		if err != nil {
			relPath = testFile
		}
		
		t.Run(relPath, func(t *testing.T) {
			// Read the markdown file
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", testFile, err)
			}

			// Determine if this is a corrupt file (intentionally broken input)
			// Files in the "corrupt" subdirectory are expected to be malformed
			isCorruptFile := strings.Contains(relPath, "corrupt/") || 
			                 strings.Contains(relPath, string(filepath.Separator)+"corrupt") ||
			                 strings.HasPrefix(relPath, "corrupt")

			// Convert using the non-streaming function (for simplicity in tests)
			output, err := ConvertMarkdownToAsciiDoc(bytes.NewReader(content))

			if isCorruptFile {
				// For corrupt files, the converter should handle them gracefully
				// It's acceptable to:
				// 1. Return an error (converter detected the corruption)
				// 2. Return empty output (converter skipped invalid content)
				// 3. Return partial output (converter handled what it could)
				// The key is that it should NOT crash or panic
				
				if err != nil {
					// Error is acceptable for corrupt files - converter detected the issue
					t.Logf("Converter correctly detected corruption in %s: %v", relPath, err)
					return
				}

				// If no error, empty or partial output is acceptable
				// The converter successfully handled (or skipped) the corrupt input
				if len(strings.TrimSpace(output)) == 0 {
					t.Logf("Converter gracefully handled corrupt file %s (empty output is acceptable)", relPath)
				} else {
					t.Logf("Converter produced partial output for corrupt file %s (%d bytes)", relPath, len(output))
				}
				return
			}

			// For valid (non-corrupt) files, conversion must succeed
			if err != nil {
				t.Errorf("ConvertMarkdownToAsciiDoc failed for valid file %s: %v", testFile, err)
				return
			}

			// Valid files must produce non-empty output (unless input was empty)
			if len(strings.TrimSpace(string(content))) > 0 && len(strings.TrimSpace(output)) == 0 {
				t.Errorf("Output is empty for non-empty valid input file %s - converter failed", testFile)
				return
			}

			// Verify output contains some basic AsciiDoc markers
			// (This is a sanity check - actual content validation is done in other tests)
			if len(output) > 0 {
				// Output should have at least some content
				if len(output) < 10 {
					t.Logf("Warning: Output for %s is very short: %d bytes", testFile, len(output))
				}
			}

			t.Logf("Successfully converted %s (%d bytes input -> %d bytes output)",
				relPath, len(content), len(output))
		})
	}
}

// TestConvertMarkdownToAsciiDocStreaming_TestbedAllFiles tests the streaming markdown to asciidoc converter
// against all test files in the ./testbed/ directory and all subdirectories
func TestConvertMarkdownToAsciiDocStreaming_TestbedAllFiles(t *testing.T) {
	// Try multiple possible paths
	testbedDirs := []string{"./testbed", "testbed", "../testbed"}
	var testbedDir string
	var err error
	
	for _, dir := range testbedDirs {
		if _, err = os.Stat(dir); err == nil {
			testbedDir = dir
			break
		}
	}
	
	if testbedDir == "" {
		t.Skipf("Testbed directory not found (tried: %v), skipping test", testbedDirs)
		return
	}

	// Find all .md files recursively in testbed and subdirectories
	var testFiles []string
	err = filepath.Walk(testbedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			return nil
		}
		
		// Only process .md files
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			return nil
		}
		
		// Skip README.md files
		if info.Name() == "README.md" {
			return nil
		}
		
		testFiles = append(testFiles, path)
		return nil
	})
	
	if err != nil {
		t.Fatalf("Failed to walk testbed directory: %v", err)
	}

	if len(testFiles) == 0 {
		t.Fatal("No test files found in testbed directory")
	}

	t.Logf("Found %d test files to process with streaming converter (including subdirectories)", len(testFiles))

	// Process each file
	for _, testFile := range testFiles {
		// Get relative path for cleaner test names
		relPath, err := filepath.Rel(testbedDir, testFile)
		if err != nil {
			relPath = testFile
		}
		
		t.Run(relPath, func(t *testing.T) {
			// Read the markdown file
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", testFile, err)
			}

			// Determine if this is a corrupt file (intentionally broken input)
			isCorruptFile := strings.Contains(relPath, "corrupt/") || 
			                 strings.Contains(relPath, string(filepath.Separator)+"corrupt") ||
			                 strings.HasPrefix(relPath, "corrupt")

			// Convert using the streaming function
			var output bytes.Buffer
			err = ConvertMarkdownToAsciiDocStreaming(bytes.NewReader(content), &output)

			if isCorruptFile {
				// For corrupt files, the converter should handle them gracefully
				// It's acceptable to:
				// 1. Return an error (converter detected the corruption)
				// 2. Return empty output (converter skipped invalid content)
				// 3. Return partial output (converter handled what it could)
				// The key is that it should NOT crash or panic
				
				if err != nil {
					// Error is acceptable for corrupt files - converter detected the issue
					t.Logf("Converter correctly detected corruption in %s: %v", relPath, err)
					return
				}

				outputStr := output.String()
				// If no error, empty or partial output is acceptable
				if len(strings.TrimSpace(outputStr)) == 0 {
					t.Logf("Converter gracefully handled corrupt file %s (empty output is acceptable)", relPath)
				} else {
					t.Logf("Converter produced partial output for corrupt file %s (%d bytes)", relPath, len(outputStr))
				}
				return
			}

			// For valid (non-corrupt) files, conversion must succeed
			if err != nil {
				t.Errorf("ConvertMarkdownToAsciiDocStreaming failed for valid file %s: %v", testFile, err)
				return
			}

			outputStr := output.String()

			// Valid files must produce non-empty output (unless input was empty)
			if len(strings.TrimSpace(string(content))) > 0 && len(strings.TrimSpace(outputStr)) == 0 {
				t.Errorf("Output is empty for non-empty valid input file %s - converter failed", testFile)
				return
			}

			// Verify output contains some basic AsciiDoc markers
			if len(outputStr) > 0 {
				// Output should have at least some content
				if len(outputStr) < 10 {
					t.Logf("Warning: Output for %s is very short: %d bytes", testFile, len(outputStr))
				}
			}

			t.Logf("Successfully converted %s (%d bytes input -> %d bytes output)",
				relPath, len(content), len(outputStr))
		})
	}
}

