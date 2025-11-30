package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ndx-video/asciidoc-xml/lib"
)

var (
	overwriteAll      bool
	skipAll           bool
	autoOverwrite     bool
	noXSL             bool
	noPicoCSS         bool
	xslFile           string
	outputType        string
	outputDir         string
	filesListFile     string
	
	// Parallel processing & limits flags
	maxWorkers        int
	parallelThreshold int
	noParallel        bool
	maxFileSize       int64
	maxArchiveSize    int64
	maxFileCount      int
	extractArchives   bool
	inputFolders      string
	dryRun            bool
	validateOnly      bool
	preserveStructure bool
)

type Config struct {
	AutoOverwrite *bool   `json:"autoOverwrite"`
	NoXSL         *bool   `json:"noXSL"`
	NoPicoCSS     *bool   `json:"noPicoCSS"`
	XSLFile       *string `json:"xslFile"`
	OutputType    *string `json:"outputType"`
	OutputDir     *string `json:"outputDir"`
	
	// New batch processing fields
	InputFolders        []string        `json:"inputFolders"`
	ExtractArchives     *bool           `json:"extractArchives"`
	MaxFileSize         *int64           `json:"maxFileSize"`
	MaxArchiveSize      *int64           `json:"maxArchiveSize"`
	MaxFileCount        *int             `json:"maxFileCount"`
	MaxWorkers          *int             `json:"maxWorkers"`
	ParallelThreshold   *int            `json:"parallelThreshold"`
	NoParallel          *bool            `json:"noParallel"`
	DryRun              *bool            `json:"dryRun"`
	ValidateOnly        *bool            `json:"validateOnly"`
	PreserveStructure   *bool            `json:"preserveStructure"`
	Logging             *lib.LogConfig  `json:"logging"`
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

	// New flags
	flag.IntVar(&maxWorkers, "workers", runtime.GOMAXPROCS(0), "Maximum concurrent workers")
	flag.IntVar(&maxWorkers, "w", runtime.GOMAXPROCS(0), "Maximum concurrent workers (shorthand)")
	flag.IntVar(&parallelThreshold, "parallel-threshold", lib.DefaultThreshold, "Minimum files to enable parallelization")
	flag.BoolVar(&noParallel, "no-parallel", false, "Disable parallelization (force sequential processing)")
	flag.Int64Var(&maxFileSize, "max-file-size", lib.DefaultMaxFileSize, "Maximum file size in bytes")
	flag.Int64Var(&maxArchiveSize, "max-archive-size", lib.DefaultMaxArchiveSize, "Maximum archive size in bytes")
	flag.IntVar(&maxFileCount, "max-file-count", lib.DefaultMaxFileCount, "Maximum files per batch")
	flag.BoolVar(&extractArchives, "extract-archives", false, "Extract compressed files in sequence before processing")
	flag.StringVar(&inputFolders, "input-folders", "", "Comma-separated list of input folders")
	flag.BoolVar(&dryRun, "dry-run", false, "Check files without processing")
	flag.BoolVar(&validateOnly, "validate-only", false, "Validate files without generating output")
	flag.BoolVar(&preserveStructure, "preserve-structure", true, "Maintain directory structure in output")

	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("adc version %s\n", version)
		os.Exit(0)
	}

	// Load configuration from adc.json
	jsonConfig := loadConfig()

	// Initialize logger
	var logger *lib.Logger
	if jsonConfig.Logging != nil {
		var err error
		logger, err = lib.NewLogger(*jsonConfig.Logging)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
		defer logger.Close()
	} else {
		// Default logger config (console only, info level)
		defaultConfig := lib.LogConfig{
			Level:  "info",
			Format: "text",
			Console: lib.ConsoleConfig{
				Enabled: true,
				Level:   "info",
			},
			File: lib.FileConfig{
				Enabled: false,
			},
		}
		var err error
		logger, err = lib.NewLogger(defaultConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
		defer logger.Close()
	}

	logger.Infof("Starting asciidoc-xml CLI")

	if flag.NArg() == 0 && filesListFile == "" && inputFolders == "" && (jsonConfig.InputFolders == nil || len(jsonConfig.InputFolders) == 0) {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file.adoc|file.md|directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Supports: AsciiDoc (.adoc) and Markdown (.md, .markdown) files\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Merge config with flags (flags take precedence, then config, then defaults)
	batchConfig := lib.BatchConfig{
		MaxWorkers:        maxWorkers,
		ParallelThreshold: parallelThreshold,
		EnableParallel:    !noParallel,
		DryRun:            dryRun,
		ValidateOnly:      validateOnly,
	}
	
	limits := lib.ProcessingLimits{
		MaxFileSize:    maxFileSize,
		MaxArchiveSize: maxArchiveSize,
		MaxFileCount:   maxFileCount,
	}

	// Process input folders
	var searchPaths []string
	if flag.NArg() > 0 {
		for _, arg := range flag.Args() {
			searchPaths = append(searchPaths, arg)
		}
	}
	if inputFolders != "" {
		paths := strings.Split(inputFolders, ",")
		for _, p := range paths {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				searchPaths = append(searchPaths, trimmed)
			}
		}
	}
	if jsonConfig.InputFolders != nil {
		searchPaths = append(searchPaths, jsonConfig.InputFolders...)
	}

	var files []string

	// Process files from --files flag
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

	// Handle archive extraction if requested
	if extractArchives {
		for _, path := range searchPaths {
			info, err := os.Stat(path)
			if err == nil && !info.IsDir() {
				// Check if it's a supported archive
				format := lib.DetectArchiveFormat(path)
				if format != "" {
					logger.Info(nil, "Extracting archive",
						"archive", path,
						"format", format,
					)
					// Create extraction dir
					extDir := strings.TrimSuffix(path, filepath.Ext(path)) + "_extracted"
					if err := os.MkdirAll(extDir, 0755); err != nil {
						logger.Error(nil, "Failed to create extraction directory",
							"directory", extDir,
							"error", err.Error(),
						)
						continue
					}
					
					f, err := os.Open(path)
					if err != nil {
						logger.Error(nil, "Failed to open archive",
							"archive", path,
							"error", err.Error(),
						)
						continue
					}
					defer f.Close()
					
					if err := lib.ExtractArchive(f, path, extDir); err != nil {
						logger.Error(nil, "Failed to extract archive",
							"archive", path,
							"error", err.Error(),
						)
						continue
					}
					logger.Info(nil, "Archive extracted successfully",
						"archive", path,
						"extracted_to", extDir,
					)
					// Add extracted dir to search paths for subsequent processing
					// (We don't add to files directly, we let the walker find them)
					searchPaths = append(searchPaths, extDir)
				}
			}
		}
	}

	// Walk paths to find .adoc, .md, and .markdown files
	for _, path := range searchPaths {
		info, err := os.Stat(path)
		if err != nil {
			logger.Error(nil, "Error accessing path",
				"path", path,
				"error", err.Error(),
			)
			continue
		}

		if info.IsDir() {
			logger.Debug(nil, "Scanning directory for .adoc, .md, and .markdown files", "directory", path)
			err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
				if err != nil {
					logger.Warn(nil, "Error accessing file during directory walk",
						"file", p,
						"error", err.Error(),
					)
					return nil // Continue walking
				}
				if !d.IsDir() {
					lowerPath := strings.ToLower(p)
					if strings.HasSuffix(lowerPath, ".adoc") || 
					   strings.HasSuffix(lowerPath, ".md") || 
					   strings.HasSuffix(lowerPath, ".markdown") {
						files = append(files, p)
					}
				}
				return nil
			})
			if err != nil {
				logger.Error(nil, "Error traversing directory",
					"directory", path,
					"error", err.Error(),
				)
			}
		} else {
			lowerPath := strings.ToLower(path)
			if strings.HasSuffix(lowerPath, ".adoc") || 
			   strings.HasSuffix(lowerPath, ".md") || 
			   strings.HasSuffix(lowerPath, ".markdown") {
				files = append(files, path)
			}
		}
	}

	if len(files) == 0 {
		logger.Error(nil, "No supported files found to process (.adoc, .md, .markdown)")
		fmt.Fprintf(os.Stderr, "No supported files found to process (.adoc, .md, .markdown)\n")
		os.Exit(1)
	}

	logger.Info(nil, "Found files to process",
		"file_count", len(files),
	)

	// Check if we have any markdown files (they convert directly to AsciiDoc, skip XSLT)
	hasMarkdownFiles := false
	for _, file := range files {
		if isMarkdownFile(file) {
			hasMarkdownFiles = true
			break
		}
	}

	// Determine XSLT file (only needed for XML output with XSLT transformation)
	// Skip XSLT check if we have markdown files (they convert directly to AsciiDoc)
	var xsltPath string
	if !hasMarkdownFiles && !noXSL && outputType == "xml" && !validateOnly {
		if xslFile != "" {
			xsltPath = xslFile
		} else {
			xsltPath = "default.xsl"
		}
		if _, err := os.Stat(xsltPath); os.IsNotExist(err) {
			logger.Error(nil, "XSLT file not found",
				"xslt_file", xsltPath,
			)
			fmt.Fprintf(os.Stderr, "Error: XSLT file not found: %s\n", xsltPath)
			fmt.Fprintf(os.Stderr, "Use --no-xsl to generate XML only, or --xsl to specify a different XSLT file\n")
			os.Exit(1)
		}
		logger.Info(nil, "Using XSLT file",
			"xslt_file", xsltPath,
		)
	}

	logger.Info(nil, "Starting file processing",
		"file_count", len(files),
		"output_type", outputType,
		"dry_run", dryRun,
		"validate_only", validateOnly,
	)
	fmt.Printf("Processing %d files...\n", len(files))
	if dryRun {
		fmt.Println("DRY RUN: No files will be converted.")
	}

	// Process files using parallel processor
	results := lib.ProcessFilesParallel(files, func(file string) error {
		if validateOnly {
			// Validation only works for AsciiDoc files
			if isMarkdownFile(file) {
				// For markdown files, we can't validate them as AsciiDoc
				// Just check if the file is readable
				f, err := os.Open(file)
				if err != nil {
					return err
				}
				f.Close()
				return nil
			}
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			return lib.Validate(f)
		}
		return processFile(file, xsltPath, outputType, logger)
	}, batchConfig, limits, func(current, total int, file string, err error) {
		// Progress callback
		status := "OK"
		if err != nil {
			status = "ERR"
		}
		
		// Calculate percentage
		percent := 0
		if total > 0 {
			percent = int(float64(current) / float64(total) * 100)
		}
		
		// Clear line and print progress
		// Using carriage return \r to overwrite line
		// Limit filename length to keep bar readable
		shortName := filepath.Base(file)
		if len(shortName) > 30 {
			shortName = shortName[:27] + "..."
		}
		
		barWidth := 20
		filled := int(float64(barWidth) * float64(current) / float64(total))
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
		
		if current > 0 {
			// Only print updating line for progress
			fmt.Printf("\r[%s] %d%% (%d/%d) %s: %s    ", bar, percent, current, total, status, shortName)
		}
	}, logger)

	logger.Info(nil, "Processing completed",
		"total_files", results.TotalFiles,
		"success_count", results.SuccessCount,
		"error_count", results.ErrorCount,
		"duration_ms", results.Duration.Milliseconds(),
	)

	fmt.Printf("\n\nProcessing complete in %v\n", results.Duration)
	fmt.Printf("Success: %d\n", results.SuccessCount)
	fmt.Printf("Errors:  %d\n", results.ErrorCount)
	
	if results.ErrorCount > 0 {
		logger.Warn(nil, "Processing completed with errors",
			"error_count", results.ErrorCount,
		)
		fmt.Println("\nError details:")
		for _, err := range results.Errors {
			logger.Error(nil, "File processing error",
				"file", err.File,
				"error", err.Error.Error(),
			)
			fmt.Printf(" - %s: %v\n", err.File, err.Error)
		}
		os.Exit(1)
	}

	logger.Infof("All files processed successfully")
}

// isMarkdownFile checks if a file is a markdown file based on its extension
func isMarkdownFile(filename string) bool {
	lower := strings.ToLower(filename)
	return strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".markdown")
}

func processFile(adocFile, xsltPath, outputType string, logger *lib.Logger) error {
	if logger != nil {
		logger.Debug(nil, "Processing file",
			"file", adocFile,
			"output_type", outputType,
		)
	}

	// Handle markdown files - convert to AsciiDoc first
	if isMarkdownFile(adocFile) {
		return processMarkdownFile(adocFile, logger)
	}

	// Read AsciiDoc file
	adocContent, err := os.ReadFile(adocFile)
	if err != nil {
		if logger != nil {
			logger.Error(nil, "Failed to read file",
				"file", adocFile,
				"error", err.Error(),
			)
		}
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Convert based on output type
	var output string
	var outputFile string
	var extension string

	// Use CDN link for PicoCSS in CLI
	usePicoCSS := !noPicoCSS
	picoCDNPath := ""
	if usePicoCSS && (outputType == "html" || outputType == "xhtml") {
		picoCDNPath = "https://cdn.jsdelivr.net/npm/@picocss/pico@2.1.1/css/pico.min.css"
	}

	switch outputType {
	case "xml":
		output, err = lib.ConvertToXML(strings.NewReader(string(adocContent)))
		if err != nil {
			if logger != nil {
				logger.Error(nil, "XML conversion failed",
					"file", adocFile,
					"error", err.Error(),
				)
			}
			return fmt.Errorf("conversion failed: %w", err)
		}
		extension = ".xml"
	case "html":
		output, err = lib.ConvertToHTML(strings.NewReader(string(adocContent)), false, usePicoCSS, picoCDNPath, "")
		if err != nil {
			if logger != nil {
				logger.Error(nil, "HTML conversion failed",
					"file", adocFile,
					"error", err.Error(),
				)
			}
			return fmt.Errorf("conversion failed: %w", err)
		}
		extension = ".html"
	case "xhtml":
		output, err = lib.ConvertToHTML(strings.NewReader(string(adocContent)), true, usePicoCSS, picoCDNPath, "")
		if err != nil {
			if logger != nil {
				logger.Error(nil, "XHTML conversion failed",
					"file", adocFile,
					"error", err.Error(),
				)
			}
			return fmt.Errorf("conversion failed: %w", err)
		}
		extension = ".xhtml"
	default:
		if logger != nil {
			logger.Error(nil, "Unsupported output type",
				"file", adocFile,
				"output_type", outputType,
			)
		}
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
		
		if preserveStructure {
			// If inputFolders was used, we need to handle relative paths, but here we simplified
			// CLI logic usually implies explicit output dir overrides structure unless complex mapping
			// For now, simple flat output if structure not handled by caller logic
			// TODO: Enhance structure preservation for CLI recursive walks if needed
			outputFile = filepath.Join(outputDir, baseName+extension)
		} else {
			outputFile = filepath.Join(outputDir, baseName+extension)
		}
	} else {
		outputFile = strings.TrimSuffix(adocFile, filepath.Ext(adocFile)) + extension
	}

	// Check if output file exists
	if _, err := os.Stat(outputFile); err == nil {
		if !autoOverwrite && !overwriteAll && !skipAll {
			// For parallel processing, we can't easily prompt interactive
			// So we assume skip if not forced
			return fmt.Errorf("file %s exists (use -y to overwrite)", outputFile)
		} else if skipAll {
			return nil // Skip
		}
	}

	// Write output file
	if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
		if logger != nil {
			logger.Error(nil, "Failed to write output file",
				"file", adocFile,
				"output_file", outputFile,
				"error", err.Error(),
			)
		}
		return fmt.Errorf("failed to write output file: %w", err)
	}

	if logger != nil {
		logger.Debug(nil, "File converted successfully",
			"file", adocFile,
			"output_file", outputFile,
		)
	}

	// Apply XSLT transformation if requested (only for XML output)
	if !noXSL && xsltPath != "" && outputType == "xml" {
		var htmlFile string
		if outputDir != "" {
			htmlFile = filepath.Join(outputDir, baseName+".html")
		} else {
			htmlFile = strings.TrimSuffix(adocFile, filepath.Ext(adocFile)) + ".html"
		}
		
		if logger != nil {
			logger.Debug(nil, "Applying XSLT transformation",
				"xml_file", outputFile,
				"html_file", htmlFile,
				"xslt_file", xsltPath,
			)
		}
		
		if err := applyXSLT(outputFile, xsltPath, htmlFile); err != nil {
			if logger != nil {
				logger.Error(nil, "XSLT transformation failed",
					"xml_file", outputFile,
					"html_file", htmlFile,
					"xslt_file", xsltPath,
					"error", err.Error(),
				)
			}
			return fmt.Errorf("XSLT transformation failed: %w", err)
		}
		
		if logger != nil {
			logger.Debug(nil, "XSLT transformation completed",
				"html_file", htmlFile,
			)
		}
	}

	return nil
}

// processMarkdownFile converts a markdown file to AsciiDoc format
func processMarkdownFile(mdFile string, logger *lib.Logger) error {
	if logger != nil {
		logger.Debug(nil, "Processing markdown file",
			"file", mdFile,
		)
	}

	// Open input markdown file
	inFile, err := os.Open(mdFile)
	if err != nil {
		if logger != nil {
			logger.Error(nil, "Failed to open markdown file",
				"file", mdFile,
				"error", err.Error(),
			)
		}
		return fmt.Errorf("failed to open markdown file: %w", err)
	}
	defer inFile.Close()

	// Determine output file path
	fileName := filepath.Base(mdFile)
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	var outputFile string

	if outputDir != "" {
		// Ensure output directory exists
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
		
		if preserveStructure {
			// For now, simple flat output if structure not handled by caller logic
			// TODO: Enhance structure preservation for CLI recursive walks if needed
			outputFile = filepath.Join(outputDir, baseName+".adoc")
		} else {
			outputFile = filepath.Join(outputDir, baseName+".adoc")
		}
	} else {
		// Output in same directory as input file
		outputFile = strings.TrimSuffix(mdFile, filepath.Ext(mdFile)) + ".adoc"
	}

	// Check if output file exists
	if _, err := os.Stat(outputFile); err == nil {
		if !autoOverwrite && !overwriteAll && !skipAll {
			// For parallel processing, we can't easily prompt interactive
			// So we assume skip if not forced
			return fmt.Errorf("file %s exists (use -y to overwrite)", outputFile)
		} else if skipAll {
			return nil // Skip
		}
	}

	// Create output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		if logger != nil {
			logger.Error(nil, "Failed to create output file",
				"file", mdFile,
				"output_file", outputFile,
				"error", err.Error(),
			)
		}
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Use streaming converter for markdown to AsciiDoc conversion
	if err := lib.ConvertMarkdownToAsciiDocStreaming(inFile, outFile); err != nil {
		if logger != nil {
			logger.Error(nil, "Markdown to AsciiDoc conversion failed",
				"file", mdFile,
				"output_file", outputFile,
				"error", err.Error(),
			)
		}
		return fmt.Errorf("markdown to AsciiDoc conversion failed: %w", err)
	}

	if logger != nil {
		logger.Debug(nil, "Markdown file converted successfully",
			"file", mdFile,
			"output_file", outputFile,
		)
	}

	return nil
}

func loadConfig() Config {
	var config Config
	file, err := os.ReadFile("adc.json")
	if err != nil {
		if os.IsNotExist(err) {
			return config
		}
		fmt.Fprintf(os.Stderr, "Error reading adc.json: %v\n", err)
		return config
	}

	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing adc.json: %v\n", err)
		return config
	}

	visited := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		visited[f.Name] = true
	})

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
	
	// New config fields
	if config.MaxWorkers != nil && !isSet("workers", "w") {
		maxWorkers = *config.MaxWorkers
	}
	if config.ParallelThreshold != nil && !isSet("parallel-threshold") {
		parallelThreshold = *config.ParallelThreshold
	}
	if config.NoParallel != nil && !isSet("no-parallel") {
		noParallel = *config.NoParallel
	}
	if config.MaxFileSize != nil && !isSet("max-file-size") {
		maxFileSize = *config.MaxFileSize
	}
	if config.MaxArchiveSize != nil && !isSet("max-archive-size") {
		maxArchiveSize = *config.MaxArchiveSize
	}
	if config.MaxFileCount != nil && !isSet("max-file-count") {
		maxFileCount = *config.MaxFileCount
	}
	if config.ExtractArchives != nil && !isSet("extract-archives") {
		extractArchives = *config.ExtractArchives
	}
	if config.DryRun != nil && !isSet("dry-run") {
		dryRun = *config.DryRun
	}
	if config.ValidateOnly != nil && !isSet("validate-only") {
		validateOnly = *config.ValidateOnly
	}
	if config.PreserveStructure != nil && !isSet("preserve-structure") {
		preserveStructure = *config.PreserveStructure
	}
	
	return config
}

func applyXSLT(xmlFile, xsltFile, htmlFile string) error {
	cmd := exec.Command("xsltproc", "-o", htmlFile, xsltFile, xmlFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("xsltproc failed: %v, output: %s", err, string(output))
	}
	return nil
}
