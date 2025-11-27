package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"asciidoc-xml/lib"
)

var (
	overwriteAll    bool
	skipAll         bool
	autoOverwrite   bool
	noXSL           bool
	xslFile         string
	outputType      string
)

func main() {
	flag.BoolVar(&autoOverwrite, "y", false, "Automatically overwrite existing files without prompting")
	flag.BoolVar(&noXSL, "no-xsl", false, "Generate XML only, skip XSLT transformation")
	flag.StringVar(&xslFile, "xsl", "", "Path to XSLT file (default: ./default.xsl)")
	flag.StringVar(&outputType, "output", "xml", "Output type: xml, html, or xhtml (default: xml)")
	flag.StringVar(&outputType, "o", "xml", "Output type: xml, html, or xhtml (shorthand for --output)")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file.adoc|directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	path := flag.Arg(0)
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Validate output type
	outputType = strings.ToLower(outputType)
	if outputType != "xml" && outputType != "html" && outputType != "xhtml" {
		fmt.Fprintf(os.Stderr, "Error: Invalid output type: %s. Supported types: xml, html, xhtml\n", outputType)
		os.Exit(1)
	}

	// Determine XSLT file (only needed for XML output with XSLT transformation)
	var xsltPath string
	if !noXSL && outputType == "xml" {
		if xslFile != "" {
			xsltPath = xslFile
		} else {
			// Look for default.xsl in current directory
			xsltPath = "default.xsl"
		}

		// Verify XSLT file exists
		if _, err := os.Stat(xsltPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: XSLT file not found: %s\n", xsltPath)
			fmt.Fprintf(os.Stderr, "Use --no-xsl to generate XML only, or --xsl to specify a different XSLT file\n")
			os.Exit(1)
		}
	}

	var files []string
	if info.IsDir() {
		// Traverse directory for .adoc files
		err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(strings.ToLower(p), ".adoc") {
				files = append(files, p)
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error traversing directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Single file
		if !strings.HasSuffix(strings.ToLower(path), ".adoc") {
			fmt.Fprintf(os.Stderr, "Error: file must have .adoc extension\n")
			os.Exit(1)
		}
		files = []string{path}
	}

	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "No .adoc files found\n")
		os.Exit(1)
	}

	// Process each file
	successCount := 0
	errorCount := 0

	for _, adocFile := range files {
		if err := processFile(adocFile, xsltPath, outputType); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", adocFile, err)
			errorCount++
		} else {
			successCount++
		}
	}

	// Summary
	fmt.Printf("\nProcessed: %d successful, %d errors\n", successCount, errorCount)
	if errorCount > 0 {
		os.Exit(1)
	}
}

func processFile(adocFile, xsltPath, outputType string) error {
	// Read AsciiDoc file
	adocContent, err := os.ReadFile(adocFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Convert based on output type
	var output string
	var outputFile string
	var extension string

	switch outputType {
	case "xml":
		output, err = lib.ConvertToXML(strings.NewReader(string(adocContent)))
		if err != nil {
			return fmt.Errorf("conversion failed: %w", err)
		}
		extension = ".xml"
	case "html":
		output, err = lib.ConvertToHTML(strings.NewReader(string(adocContent)), false)
		if err != nil {
			return fmt.Errorf("conversion failed: %w", err)
		}
		extension = ".html"
	case "xhtml":
		output, err = lib.ConvertToHTML(strings.NewReader(string(adocContent)), true)
		if err != nil {
			return fmt.Errorf("conversion failed: %w", err)
		}
		extension = ".xhtml"
	default:
		return fmt.Errorf("unsupported output type: %s", outputType)
	}

	// Determine output file path
	outputFile = strings.TrimSuffix(adocFile, filepath.Ext(adocFile)) + extension

	// Check if output file exists
	if _, err := os.Stat(outputFile); err == nil {
		// File exists, check if we should overwrite
		if !autoOverwrite && !overwriteAll && !skipAll {
			response := promptOverwrite(outputFile)
			switch response {
			case "a": // all
				overwriteAll = true
			case "n": // no/all
				skipAll = true
				fmt.Printf("Skipping %s\n", outputFile)
				return nil
			case "y": // yes
				// Continue to overwrite
			case "q": // quit
				fmt.Println("Cancelled by user")
				os.Exit(0)
			}
		} else if skipAll {
			fmt.Printf("Skipping %s\n", outputFile)
			return nil
		}
	}

	// Write output file
	if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Converted: %s -> %s\n", adocFile, outputFile)

	// Apply XSLT transformation if requested (only for XML output)
	if !noXSL && xsltPath != "" && outputType == "xml" {
		htmlFile := strings.TrimSuffix(adocFile, filepath.Ext(adocFile)) + ".html"
		
		// Check if HTML file exists
		if _, err := os.Stat(htmlFile); err == nil {
			if !autoOverwrite && !overwriteAll && !skipAll {
				response := promptOverwrite(htmlFile)
				switch response {
				case "a":
					overwriteAll = true
				case "n":
					skipAll = true
					fmt.Printf("Skipping %s\n", htmlFile)
					return nil
				case "y":
					// Continue
				case "q":
					fmt.Println("Cancelled by user")
					os.Exit(0)
				}
			} else if skipAll {
				fmt.Printf("Skipping %s\n", htmlFile)
				return nil
			}
		}

		// Apply XSLT transformation using xsltproc (external tool)
		if err := applyXSLT(outputFile, xsltPath, htmlFile); err != nil {
			return fmt.Errorf("XSLT transformation failed: %w", err)
		}
		fmt.Printf("Transformed: %s -> %s\n", outputFile, htmlFile)
	}

	return nil
}

func applyXSLT(xmlFile, xsltFile, htmlFile string) error {
	// Use xsltproc for XSLT transformation
	cmd := exec.Command("xsltproc", "-o", htmlFile, xsltFile, xmlFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("xsltproc failed: %v, output: %s", err, string(output))
	}
	return nil
}

func promptOverwrite(filename string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("File %s already exists. Overwrite? [y]es/[n]o/[a]ll/[q]uit: ", filename)
		response, err := reader.ReadString('\n')
		if err != nil {
			return "q"
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response == "y" || response == "yes" {
			return "y"
		}
		if response == "n" || response == "no" {
			return "n"
		}
		if response == "a" || response == "all" {
			return "a"
		}
		if response == "q" || response == "quit" {
			return "q"
		}
		fmt.Println("Invalid response. Please enter y, n, a, or q")
	}
}

