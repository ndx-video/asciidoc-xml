//go:build testbed

package lib

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestConvertMarkdownToAsciiDoc_TestbedFiles tests the markdown to asciidoc converter
// against all test files in the ./testbed/ directory
func TestConvertMarkdownToAsciiDoc_TestbedFiles(t *testing.T) {
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

	// Find all .md files in testbed (excluding subdirectories and README.md)
	entries, err := os.ReadDir(testbedDir)
	if err != nil {
		t.Fatalf("Failed to read testbed directory: %v", err)
	}

	var testFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip subdirectories like corrupt/, etc.
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			continue
		}
		if entry.Name() == "README.md" {
			continue // Skip README
		}
		testFiles = append(testFiles, filepath.Join(testbedDir, entry.Name()))
	}

	if len(testFiles) == 0 {
		t.Fatal("No test files found in testbed directory")
	}

	t.Logf("Found %d test files to process", len(testFiles))

	// Process each file
	for _, testFile := range testFiles {
		t.Run(filepath.Base(testFile), func(t *testing.T) {
			// Read the markdown file
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", testFile, err)
			}

			// Root directory testbed files are all valid (non-corrupt) files
			// They must convert successfully
			output, err := ConvertMarkdownToAsciiDoc(bytes.NewReader(content))
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
				filepath.Base(testFile), len(content), len(output))
		})
	}
}

// TestConvertMarkdownToAsciiDocStreaming_TestbedFiles tests the streaming markdown to asciidoc converter
// against all test files in the ./testbed/ directory
func TestConvertMarkdownToAsciiDocStreaming_TestbedFiles(t *testing.T) {
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

	// Find all .md files in testbed (excluding subdirectories and README.md)
	entries, err := os.ReadDir(testbedDir)
	if err != nil {
		t.Fatalf("Failed to read testbed directory: %v", err)
	}

	var testFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip subdirectories
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			continue
		}
		if entry.Name() == "README.md" {
			continue // Skip README
		}
		testFiles = append(testFiles, filepath.Join(testbedDir, entry.Name()))
	}

	if len(testFiles) == 0 {
		t.Fatal("No test files found in testbed directory")
	}

	t.Logf("Found %d test files to process with streaming converter", len(testFiles))

	// Process each file
	for _, testFile := range testFiles {
		t.Run(filepath.Base(testFile), func(t *testing.T) {
			// Read the markdown file
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", testFile, err)
			}

			// Root directory testbed files are all valid (non-corrupt) files
			// They must convert successfully
			var output bytes.Buffer
			err = ConvertMarkdownToAsciiDocStreaming(bytes.NewReader(content), &output)
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
				filepath.Base(testFile), len(content), len(outputStr))
		})
	}
}

// TestConvertMarkdownToAsciiDoc_TestbedFiles_Consistency verifies that both
// streaming and non-streaming converters produce equivalent output.
// Note: This test may reveal differences due to implementation details (e.g., placeholder handling
// in inline markdown conversion). The main testbed tests verify that all files convert successfully,
// which is the primary goal. This consistency test helps identify potential bugs or inconsistencies.
func TestConvertMarkdownToAsciiDoc_TestbedFiles_Consistency(t *testing.T) {
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

	// Find all .md files in testbed (excluding subdirectories and README.md)
	entries, err := os.ReadDir(testbedDir)
	if err != nil {
		t.Fatalf("Failed to read testbed directory: %v", err)
	}

	var testFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			continue
		}
		if entry.Name() == "README.md" {
			continue
		}
		testFiles = append(testFiles, filepath.Join(testbedDir, entry.Name()))
	}

	if len(testFiles) == 0 {
		t.Fatal("No test files found in testbed directory")
	}

	// Process each file with both converters and compare
	for _, testFile := range testFiles {
		t.Run(filepath.Base(testFile), func(t *testing.T) {
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", testFile, err)
			}

			// Non-streaming conversion
			output1, err := ConvertMarkdownToAsciiDoc(bytes.NewReader(content))
			if err != nil {
				t.Fatalf("ConvertMarkdownToAsciiDoc failed: %v", err)
			}

			// Streaming conversion
			var output2 bytes.Buffer
			err = ConvertMarkdownToAsciiDocStreaming(bytes.NewReader(content), &output2)
			if err != nil {
				t.Fatalf("ConvertMarkdownToAsciiDocStreaming failed: %v", err)
			}

			output2Str := output2.String()

			// Compare outputs
			// Note: Due to Go map iteration order, frontmatter attributes may appear in different orders
			// between streaming and non-streaming versions. We'll do a more lenient comparison.
			if output1 != output2Str {
				// Check if the difference is only in attribute ordering (common with frontmatter)
				// by comparing line-by-line and allowing attribute lines to be in different orders
				lines1 := strings.Split(output1, "\n")
				lines2 := strings.Split(output2Str, "\n")
				
				// Extract attribute lines (lines starting with :)
				var attrs1, attrs2 []string
				var otherLines1, otherLines2 []string
				
				for _, line := range lines1 {
					if strings.HasPrefix(line, ":") {
						attrs1 = append(attrs1, line)
					} else {
						otherLines1 = append(otherLines1, line)
					}
				}
				
				for _, line := range lines2 {
					if strings.HasPrefix(line, ":") {
						attrs2 = append(attrs2, line)
					} else {
						otherLines2 = append(otherLines2, line)
					}
				}
				
				// Compare non-attribute lines (must match exactly)
				if len(otherLines1) != len(otherLines2) {
					t.Errorf("Outputs differ in number of non-attribute lines: %d vs %d",
						len(otherLines1), len(otherLines2))
					return
				}
				
				for i := range otherLines1 {
					if otherLines1[i] != otherLines2[i] {
						t.Errorf("Outputs differ at non-attribute line %d:\n  non-streaming: %q\n  streaming:     %q",
							i, otherLines1[i], otherLines2[i])
						return
					}
				}
				
				// For attributes, check that we have the same set (order may differ)
				if len(attrs1) != len(attrs2) {
					t.Errorf("Outputs differ in number of attributes: %d vs %d",
						len(attrs1), len(attrs2))
					return
				}
				
				// Create maps to compare attribute sets
				attrMap1 := make(map[string]bool)
				attrMap2 := make(map[string]bool)
				for _, attr := range attrs1 {
					attrMap1[attr] = true
				}
				for _, attr := range attrs2 {
					attrMap2[attr] = true
				}
				
				// Check if attribute sets match
				for attr := range attrMap1 {
					if !attrMap2[attr] {
						t.Errorf("Attribute %q found in non-streaming but not in streaming", attr)
						return
					}
				}
				for attr := range attrMap2 {
					if !attrMap1[attr] {
						t.Errorf("Attribute %q found in streaming but not in non-streaming", attr)
						return
					}
				}
				
				// If we get here, the only difference is attribute ordering, which is acceptable
				t.Logf("Outputs differ only in attribute ordering (acceptable due to Go map iteration order)")
			}
		})
	}
}

