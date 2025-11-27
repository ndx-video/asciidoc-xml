package main

import (
	"bufio"
	"encoding/json"
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
	noPicoCSS       bool
	xslFile         string
	outputType      string
	outputDir       string
	filesListFile   string
)

type Config struct {
	AutoOverwrite *bool   `json:"autoOverwrite"`
	NoXSL         *bool   `json:"noXSL"`
	NoPicoCSS     *bool   `json:"noPicoCSS"`
	XSLFile       *string `json:"xslFile"`
	OutputType    *string `json:"outputType"`
	OutputDir     *string `json:"outputDir"`
}

func main() {
	showVersion := flag.Bool("version", false, "Show version information and exit")
	flag.BoolVar(&autoOverwrite, "y", false, "Automatically overwrite existing files without prompting")
	flag.BoolVar(&noXSL, "no-xsl", false, "Generate XML only, skip XSLT transformation")
	flag.BoolVar(&noPicoCSS, "no-picocss", false, "Disable PicoCSS styling in HTML output (PicoCSS is enabled by default)")
	flag.StringVar(&xslFile, "xsl", "", "Path to XSLT file (default: ./default.xsl)")
	flag.StringVar(&outputType, "output", "xml", "Output type: xml, html, or xhtml (default: xml)")
	flag.StringVar(&outputType, "o", "xml", "Output type: xml, html, or xhtml (shorthand for --output)")
	flag.StringVar(&outputDir, "out-dir", "", "Output directory (default: same as input file)")
	flag.StringVar(&outputDir, "d", "", "Output directory (shorthand for --out-dir)")
	flag.StringVar(&filesListFile, "files", "", "Path to file containing list of files to process")
	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("adc version %s\n", version)
		os.Exit(0)
	}

	// Load configuration from adc.json
	loadConfig()

	if flag.NArg() == 0 && filesListFile == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file.adoc|directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var files []string

	// Process files from --files flag if provided
	if filesListFile != "" {
		content, err := os.ReadFile(filesListFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading files list: %v\n", err)
			os.Exit(1)
		}
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				files = append(files, line)
			}
		}
	}

	// Process command line arguments
	if flag.NArg() > 0 {
		path := flag.Arg(0)
		info, err := os.Stat(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

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
			files = append(files, path)
		}
	}

	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "No .adoc files found\n")
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

	// Use CDN link for PicoCSS in CLI (version pinned to match local static file)
	usePicoCSS := !noPicoCSS
	picoCDNPath := ""
	if usePicoCSS && (outputType == "html" || outputType == "xhtml") {
		// Version pinned to 2.1.1 to match what's in web/static/pico.min.css
		picoCDNPath = "https://cdn.jsdelivr.net/npm/@picocss/pico@2.1.1/css/pico.min.css"
	}

	switch outputType {
	case "xml":
		output, err = lib.ConvertToXML(strings.NewReader(string(adocContent)))
		if err != nil {
			return fmt.Errorf("conversion failed: %w", err)
		}
		extension = ".xml"
	case "html":
		output, err = lib.ConvertToHTML(strings.NewReader(string(adocContent)), false, usePicoCSS, picoCDNPath, "")
		if err != nil {
			return fmt.Errorf("conversion failed: %w", err)
		}
		extension = ".html"
	case "xhtml":
		output, err = lib.ConvertToHTML(strings.NewReader(string(adocContent)), true, usePicoCSS, picoCDNPath, "")
		if err != nil {
			return fmt.Errorf("conversion failed: %w", err)
		}
		extension = ".xhtml"
	default:
		return fmt.Errorf("unsupported output type: %s", outputType)
	}

	// Determine output file path
	fileName := filepath.Base(adocFile)
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	if outputDir != "" {
		// Ensure output directory exists
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
		outputFile = filepath.Join(outputDir, baseName+extension)
	} else {
		outputFile = strings.TrimSuffix(adocFile, filepath.Ext(adocFile)) + extension
	}

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
		fileName := filepath.Base(adocFile)
		baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		var htmlFile string

		if outputDir != "" {
			htmlFile = filepath.Join(outputDir, baseName+".html")
		} else {
			htmlFile = strings.TrimSuffix(adocFile, filepath.Ext(adocFile)) + ".html"
		}
		
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

func loadConfig() {
	file, err := os.ReadFile("adc.json")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Fprintf(os.Stderr, "Error reading adc.json: %v\n", err)
		return
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing adc.json: %v\n", err)
		return
	}

	visited := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		visited[f.Name] = true
	})

	// Helper to check if flag was set
	isSet := func(names ...string) bool {
		for _, n := range names {
			if visited[n] {
				return true
			}
		}
		return false
	}

	if config.AutoOverwrite != nil && !isSet("y") {
		autoOverwrite = *config.AutoOverwrite
	}
	if config.NoXSL != nil && !isSet("no-xsl") {
		noXSL = *config.NoXSL
	}
	if config.NoPicoCSS != nil && !isSet("no-picocss") {
		noPicoCSS = *config.NoPicoCSS
	}
	if config.XSLFile != nil && !isSet("xsl") {
		xslFile = *config.XSLFile
	}
	if config.OutputType != nil && !isSet("output", "o") {
		outputType = *config.OutputType
	}
	if config.OutputDir != nil && !isSet("out-dir", "d") {
		outputDir = *config.OutputDir
	}
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

