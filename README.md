# AsciiDoc XML Converter

**Version 0.5.0**

A pure Go package for converting AsciiDoc documents to a custom XML format, designed for easy transformation to HTML via XSLT. Also provides comprehensive Markdown to AsciiDoc conversion with full CommonMark and GitHub Flavored Markdown (GFM) support. Includes a web-based development harness for testing XML generation and XSLT template development.

## Overview

This package provides:

- **Custom XML Schema (XSD)**: A purpose-built XML schema specifically designed for AsciiDoc, avoiding the bloat of DocBook
- **Go XML Structures**: Type-safe Go structs matching the XSD schema
- **AsciiDoc Parser**: Comprehensive pure Go parser that converts AsciiDoc source to XML with support for inline formatting, macros, cross-references, attributes, and more
- **Markdown Converter**: Full CommonMark and GitHub Flavored Markdown (GFM) to AsciiDoc conversion with streaming support for large files
- **XSLT Template**: Comprehensive XSLT stylesheet for transforming XML to HTML
- **Web Development Harness**: Single Page Application (SPA) for interactive development and testing
- **Batch Processing**: Parallel file processing with archive support (ZIP, TAR, TAR.GZ)
- **Comprehensive Logging**: Structured logging with levels, file rotation, and request tracking

## Features

- ✅ Pure Go implementation (no external binaries required)
- ✅ Cross-platform support (no architecture-specific dependencies)
- ✅ Comprehensive AsciiDoc feature support (inline formatting, macros, cross-references, attributes, footnotes, etc.)
- ✅ **Markdown to AsciiDoc conversion** with full CommonMark and GFM support
- ✅ **Streaming conversion** for memory-efficient processing of large files
- ✅ **Batch processing** with parallel execution and progress tracking
- ✅ **Archive support** (ZIP, TAR, TAR.GZ) with automatic extraction
- ✅ **Comprehensive logging** with structured output, file rotation, and request tracking
- ✅ Well-formed, validatable XML output
- ✅ Semantic HTML output via XSLT
- ✅ Web-based development harness with live preview
- ✅ Command line tool for batch conversion
- ✅ JavaScript syntax highlighting and pretty-printing
- ✅ Comprehensive test suite (Go + JavaScript) with testbed coverage
- ✅ Extensible and customizable

## Package Structure

```
asciidoc-xml/
├── schema/
│   └── asciidoc.xsd          # XML Schema Definition
├── lib/
│   ├── adoc-parser.go        # Core AsciiDoc parser library
│   ├── adoc-parser_test.go   # Parser tests
│   ├── converter.go          # Converter functions (XML, HTML, XHTML)
│   ├── converter_test.go     # Converter tests
│   ├── converter_streaming.go # Markdown to AsciiDoc streaming converter
│   ├── converter_streaming_test.go # Streaming converter tests
│   ├── converter_testbed_test.go # Testbed test suite
│   ├── batch.go              # Batch processing with parallel execution
│   ├── batch_test.go         # Batch processing tests
│   ├── archive.go            # Archive extraction (ZIP, TAR, TAR.GZ)
│   ├── archive_test.go       # Archive extraction tests
│   ├── logger.go             # Comprehensive logging system
│   ├── logger_test.go       # Logging tests
│   ├── ast.go                # Abstract Syntax Tree definitions
│   ├── ast_test.go           # AST tests
│   └── version.go            # Library version
├── docs/
│   ├── asciidoc-xml.adoc     # User guide documentation
│   ├── api.adoc              # Web API documentation
│   ├── cms-upgrade-guide-v0.4.0.md # CMS upgrade guide
│   └── cms-agent-prompt.md   # CMS developer agent prompt
├── cli/
│   ├── adc.go                # Command line tool (AsciiDoc Converter)
│   ├── adc_test.go           # CLI tests
│   ├── adc.json              # CLI configuration file
│   └── version.go            # CLI version
├── web/
│   ├── main.go               # Web server
│   ├── main_test.go          # Server tests
│   ├── version.go            # Web server version
│   ├── static/
│   │   ├── app.js            # Main application logic
│   │   ├── app.css           # Styles
│   │   ├── pretty.js         # Syntax highlighting
│   │   └── js_test.go        # JavaScript tests (using Goja)
│   └── templates/
│       └── index.html        # SPA template
├── watcher/
│   ├── main.go               # File watcher daemon
│   └── version.go            # Watcher version
├── testbed/                  # Comprehensive test suite
│   ├── *.md                  # Markdown test files
│   └── corrupt/              # Corrupt file tests
├── xslt/
│   └── asciidoc-to-html.xsl  # XSLT transformation template
├── examples/
│   └── comprehensive.adoc    # Example file with all features
├── harness.sh                # Development server manager
├── Makefile                  # Build and distribution automation
├── xml.go                    # Go XML struct definitions
├── VERSION                   # Project version file
└── README.md                 # This file
```

## Installation

### Pre-built Binaries

Pre-built binaries are available for multiple platforms. Download the appropriate package for your system:

**CLI-only packages** (includes only the `adc` command-line tool):
- Linux (amd64): `asciidoc-xml-cli-v1.0.0-linux-amd64.tar.gz`
- Linux (arm64): `asciidoc-xml-cli-v1.0.0-linux-arm64.tar.gz`
- macOS (amd64): `asciidoc-xml-cli-v1.0.0-darwin-amd64.tar.gz`
- macOS (arm64): `asciidoc-xml-cli-v1.0.0-darwin-arm64.tar.gz`
- Windows (amd64): `asciidoc-xml-cli-v1.0.0-windows-amd64.zip`

**CLI + Watcher packages** (includes `adc` CLI tool and `adc-watcher` daemon):
- Linux (amd64): `asciidoc-xml-cli-watcher-v1.0.0-linux-amd64.tar.gz`
- Linux (arm64): `asciidoc-xml-cli-watcher-v1.0.0-linux-arm64.tar.gz`
- macOS (amd64): `asciidoc-xml-cli-watcher-v1.0.0-darwin-amd64.tar.gz`
- macOS (arm64): `asciidoc-xml-cli-watcher-v1.0.0-darwin-arm64.tar.gz`
- Windows (amd64): `asciidoc-xml-cli-watcher-v1.0.0-windows-amd64.zip`

**Full packages** (includes CLI tool, web server, and watcher):
- Linux (amd64): `asciidoc-xml-full-v1.0.0-linux-amd64.tar.gz`
- Linux (arm64): `asciidoc-xml-full-v1.0.0-linux-arm64.tar.gz`
- macOS (amd64): `asciidoc-xml-full-v1.0.0-darwin-amd64.tar.gz`
- macOS (arm64): `asciidoc-xml-full-v1.0.0-darwin-arm64.tar.gz`
- Windows (amd64): `asciidoc-xml-full-v1.0.0-windows-amd64.zip`

#### Installing CLI-only Package

**Linux/macOS:**
```bash
# Extract the archive
tar -xzf asciidoc-xml-cli-v1.0.0-linux-amd64.tar.gz

# Move binary to PATH (optional)
sudo cp asciidoc-xml-cli-v1.0.0-linux-amd64/bin/adc /usr/local/bin/
```

**Windows:**
```powershell
# Extract the ZIP file
Expand-Archive asciidoc-xml-cli-v1.0.0-windows-amd64.zip

# Add to PATH or use from extracted directory
```

#### Installing CLI + Watcher Package

**Linux/macOS:**
```bash
# Extract the archive
tar -xzf asciidoc-xml-cli-watcher-v1.0.0-linux-amd64.tar.gz

# Move binaries to PATH (optional)
sudo cp asciidoc-xml-cli-watcher-v1.0.0-linux-amd64/bin/* /usr/local/bin/

# Start the watcher daemon
adc-watcher --watch /path/to/watch --port 8006
```

**Windows:**
```powershell
# Extract the ZIP file
Expand-Archive asciidoc-xml-cli-watcher-v1.0.0-windows-amd64.zip

# Add to PATH or use from extracted directory
```

#### Installing Full Package

**Linux/macOS:**
```bash
# Extract the archive
tar -xzf asciidoc-xml-full-v1.0.0-linux-amd64.tar.gz

# Move binaries to PATH (optional)
sudo cp asciidoc-xml-full-v1.0.0-linux-amd64/bin/* /usr/local/bin/

# Copy XSLT template (optional)
sudo cp -r asciidoc-xml-full-v1.0.0-linux-amd64/xslt /usr/local/share/asciidoc-xml/
```

**Windows:**
```powershell
# Extract the ZIP file
Expand-Archive asciidoc-xml-full-v1.0.0-windows-amd64.zip

# Add to PATH or use from extracted directory
```

#### Installing Docker Compose Stack

**Prerequisites:**
- Docker and Docker Compose installed
- Docker Hub account (or access to your registry)

**Steps:**
```bash
# Option 1: Use published images
# First, publish images to your registry
./scripts/publish-docker.sh docker.io/username v1.0.0

# Then use the production compose file
docker-compose -f docker-compose.prod.yml up -d

# Option 2: Build and use local images
docker-compose build
docker-compose up -d

# Access services
# Web: http://localhost:8005
# Watcher API: http://localhost:8006
```

## Building from Source

### Prerequisites

- Go 1.21 or later
- Git (for cloning the repository)

### Build Instructions

#### Using Make (Recommended)

The project includes a `Makefile` for easy cross-compilation:

```bash
# Build CLI tool for current platform
make cli

# Build web server for current platform
make web

# Build watcher for current platform
make watcher

# Build CLI for all target platforms
make build-cli

# Build web server for all target platforms
make build-web

# Build watcher for all target platforms
make build-watcher

# Build CLI, web, and watcher for all platforms
make build-all

# Create CLI-only distribution packages
make dist-cli

# Create CLI + watcher distribution packages
make dist-cli-watcher

# Create full distribution packages (CLI + web + watcher)
make dist-full VERSION=0.5.0

# Clean build artifacts
make clean

# Run tests
make test

# Install CLI to local system (current platform only)
make install-cli
```

Binaries are output to `bin/[GOOS-GOARCH]/` directory.

#### Manual Build

```bash
# Clone the repository
git clone https://github.com/yourusername/asciidoc-xml.git
cd asciidoc-xml

# Build CLI for current platform
go build -o adc ./cli

# Build CLI for specific platform (cross-compilation)
GOOS=linux GOARCH=amd64 go build -o adc-linux-amd64 ./cli

# Build web server
go build -o asciidoc-xml-web ./web

# Install CLI globally
go install ./cli

# Install web server globally
go install ./web
```

### Cross-Compilation

Go supports cross-compilation out of the box. Set `GOOS` and `GOARCH` environment variables:

```bash
# Linux amd64
GOOS=linux GOARCH=amd64 go build -o adc-linux-amd64 ./cli

# macOS arm64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o adc-darwin-arm64 ./cli

# Windows amd64
GOOS=windows GOARCH=amd64 go build -o adc-windows-amd64.exe ./cli
```

## Distribution

The project provides four distribution package types:

### 1. CLI-only Package

Contains only the `adc` command-line tool. Ideal for users who only need batch conversion functionality.

**Contents:**
- `bin/adc` - Command-line converter
- `LICENSE` - License file
- `README.md` - Documentation
- `examples/` - Example AsciiDoc files

**Build:** `make dist-cli`

### 2. CLI + Watcher Bundle

Contains the `adc` CLI tool and the `adc-watcher` daemon. Ideal for users who want automatic file watching and conversion.

**Contents:**
- `bin/adc` - Command-line converter
- `bin/adc-watcher` - File watcher daemon
- `LICENSE` - License file
- `README.md` - Documentation
- `examples/` - Example AsciiDoc files

**Build:** `make dist-cli-watcher`

**Usage:**
```bash
# Start the watcher daemon
./adc-watcher --watch /path/to/watch --port 8006

# The watcher will automatically run adc on changed .adoc files
```

### 3. Full Package

Contains the CLI tool, web server, watcher daemon, and XSLT templates. Ideal for users who want the complete feature set including the web interface and file watching.

**Contents:**
- `bin/adc` - Command-line converter
- `bin/asciidoc-xml-web` - Web server
- `bin/adc-watcher` - File watcher daemon
- `xslt/` - XSLT transformation templates
- `LICENSE` - License file
- `README.md` - Documentation
- `examples/` - Example AsciiDoc files

**Build:** `make dist-full`

### 4. Docker Compose Stack

Pre-built Docker images for the web server and watcher daemon, ready to deploy with Docker Compose.

**Contents:**
- Web server container (port 8005)
- Watcher daemon container (port 8006)
- Docker Compose configuration
- Health checks and automatic restarts

**Build and Publish:**
```bash
# Build and publish to Docker Hub (or your registry)
./scripts/publish-docker.sh docker.io/username v1.0.0

# Or use default registry
./scripts/publish-docker.sh
```

**Usage:**
```bash
# Using local images (development)
docker-compose up -d

# Using published images (production)
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

The Docker Compose stack includes:
- **Web service**: Accessible at `http://localhost:8005`
- **Watcher service**: API accessible at `http://localhost:8006`
- **Volume mounts**: 
  - `./examples` → `/app/examples` (read-only)
  - `./docs` → `/app/docs` (read-only)
  - `./xslt` → `/app/xslt` (read-only)
  - `./watch` → `/watch` (for watcher to monitor)

## Quick Start

### Development Harness

The easiest way to get started is using the web-based development harness:

```bash
# Start the development server
./harness.sh start

# The server will:
# 1. Run all tests (Go + JavaScript)
# 2. Start the web server on http://localhost:8005
# 3. Auto-load comprehensive.adoc example
```

The web interface provides:
- **Dynamic column layout**: AsciiDoc source (always visible), plus optional XML, XSLT, and HTML output columns
- **Output type selector**: Choose between XML, HTML, HTML5, XHTML, or XHTML5 output
- **Smart column visibility**: 
  - XML column: Only shown when XML output is selected
  - XSLT column: Only shown when XML, XHTML, or XHTML5 is selected (XSLT can process XHTML)
  - HTML column: Always visible, shows rendered or source view
- **Resizable columns**: Drag column borders to adjust widths
- **Syntax highlighting**: Color-coded AsciiDoc, XML, XSLT, and HTML
- **Live conversion**: Automatic conversion when AsciiDoc loads or output type changes
- **File upload**: Upload AsciiDoc and XSLT files via the web interface
- **Path loading**: Load files from server paths
- **Direct HTML conversion**: HTML/HTML5/XHTML output is generated directly without XML/XSLT pipeline

### Server Management

```bash
# Start server (runs tests first)
./harness.sh start

# Start server without running tests
./harness.sh start --no-test

# Stop server
./harness.sh stop

# Restart server (no port check)
./harness.sh restart

# Reload server (alias for restart)
./harness.sh reload

# Show server status
./harness.sh status

# Run all tests without starting server
./harness.sh test-all

# Use custom port
PORT=8080 ./harness.sh start
```

### Command Line Tool (`adc`)

The `adc` (AsciiDoc Converter) command line tool provides batch conversion of AsciiDoc files to XML and optionally to HTML via XSLT transformation. It also supports Markdown to AsciiDoc conversion, parallel processing, archive extraction, and comprehensive logging.

#### Building the Tool

```bash
# Build the command line tool
go build -o adc ./cli

# Or install globally
go install ./cli
```

#### Basic Usage

```bash
# Convert a single AsciiDoc file to XML
./adc document.adoc

# Convert all .adoc files in a directory (recursively)
./adc /path/to/documents/

# Convert Markdown file to AsciiDoc
./adc --output md2adoc document.md

# Convert Markdown file to AsciiDoc, then to XML
./adc --output md2adoc document.md
./adc document.adoc

# Generate XML only (skip XSLT transformation)
./adc --no-xsl document.adoc

# Specify a custom XSLT file
./adc --xsl xslt/asciidoc-to-html.xsl document.adoc

# Auto-overwrite existing files without prompting
./adc -y document.adoc

# Batch process with parallel execution
./adc --input-folders ./docs,./examples --workers 4

# Extract and process archives
./adc --extract-archives --input-folders ./archives

# Process with custom limits and logging
./adc --max-file-size 5242880 --max-file-count 1000 ./docs/
```

#### Command Options

**Basic Options:**
- `-y`: Automatically overwrite existing XML/HTML files without prompting
- `--no-xsl`: Generate XML only, skip XSLT transformation to HTML
- `--xsl <path>`: Path to XSLT file (default: `./default.xsl`)
- `--out-dir <path>` or `-d <path>`: Specify output directory (files are created here instead of source directory)
- `--files <path>`: Path to a file containing a list of files to process (one per line)
- `--output <type>` or `-o <type>`: Output type: `xml`, `html`, `xhtml`, or `md2adoc` (default: `xml`)

**Batch Processing Options:**
- `--input-folders <paths>`: Comma-separated list of input folders to process
- `--workers <n>` or `-w <n>`: Maximum concurrent workers (default: CPU count)
- `--parallel-threshold <n>`: Minimum files to enable parallelization (default: 2)
- `--no-parallel`: Force sequential processing
- `--extract-archives`: Extract compressed files (ZIP, TAR, TAR.GZ) before processing
- `--preserve-structure`: Maintain directory structure in output (default: true)

**Safety Limits:**
- `--max-file-size <bytes>`: Maximum file size in bytes (default: 10MB)
- `--max-archive-size <bytes>`: Maximum archive size in bytes (default: 10MB)
- `--max-file-count <n>`: Maximum files per batch (default: 10000)

**Validation Options:**
- `--dry-run`: Preview mode - scan files and validate limits without processing
- `--validate-only`: Validation mode - parse files to check for errors without generating output

**Logging Options:**
- Configure via `adc.json` file (see Configuration File section)

#### Configuration File

The tool supports a comprehensive configuration file `adc.json` in the current directory. Example:

```json
{
  "autoOverwrite": true,
  "noXSL": false,
  "xslFile": "custom.xsl",
  "outputType": "xml",
  "outputDir": "dist",
  "inputFolders": ["./docs", "./examples"],
  "extractArchives": false,
  "maxFileSize": 10485760,
  "maxArchiveSize": 10485760,
  "maxFileCount": 10000,
  "maxWorkers": 0,
  "parallelThreshold": 2,
  "noParallel": false,
  "dryRun": false,
  "validateOnly": false,
  "preserveStructure": true,
  "logging": {
    "level": "info",
    "format": "text",
    "file": {
      "enabled": true,
      "path": "./logs",
      "filename": "adc.log",
      "maxSize": 10485760,
      "maxFiles": 5,
      "rotation": "size"
    },
    "console": {
      "enabled": true,
      "level": "info"
    }
  }
}
```

Command line arguments override configuration file settings. See `cli/adc.json` for detailed comments on each parameter.

#### File Processing

- **Single File**: Processes one `.adoc` or `.md` file, outputs `.xml` (and optionally `.html`) or `.adoc` (for Markdown input) in the same directory
- **Directory**: Recursively traverses directory, processes all `.adoc` files found
- **Markdown Files**: Use `--output md2adoc` to convert Markdown to AsciiDoc
- **Archives**: Supports ZIP, TAR, and TAR.GZ archives with `--extract-archives` flag
- **Parallel Processing**: Automatically parallelizes when processing multiple files (configurable via `--workers` and `--parallel-threshold`)
- **Output Files**: 
  - XML: `filename.xml` (same name as input, `.xml` extension)
  - HTML: `filename.html` (generated when XSLT transformation is applied)
  - AsciiDoc: `filename.adoc` (when converting from Markdown with `--output md2adoc`)

#### Overwrite Behavior

By default, `adc` prompts before overwriting existing files:
- `y` or `yes`: Overwrite this file
- `n` or `no`: Skip this file
- `a` or `all`: Overwrite all remaining files
- `q` or `quit`: Cancel operation

Use `-y` flag to skip prompts and automatically overwrite.

#### XSLT Transformation

When XSLT transformation is enabled (default), `adc` uses `xsltproc` to transform XML to HTML. Ensure `xsltproc` is installed:

```bash
# On Ubuntu/Debian
sudo apt-get install xsltproc

# On macOS
brew install libxslt

# On Fedora/RHEL
sudo dnf install libxslt
```

#### Examples

```bash
# Convert single AsciiDoc file with default XSLT
./adc document.adoc
# Output: document.xml, document.html

# Convert Markdown to AsciiDoc
./adc --output md2adoc document.md
# Output: document.adoc

# Convert directory, XML only
./adc --no-xsl ./docs/
# Output: All .xml files, no HTML

# Batch convert with custom XSLT, auto-overwrite
./adc -y --xsl custom.xsl ./docs/
# Output: All .xml and .html files

# Process multiple folders in parallel
./adc --input-folders ./docs,./examples --workers 4
# Output: All files processed with 4 concurrent workers

# Extract and process archives
./adc --extract-archives --input-folders ./archives
# Output: Archives extracted and processed

# Dry run to preview what would be processed
./adc --dry-run --input-folders ./docs
# Output: File list and validation without processing

# Process with custom limits
./adc --max-file-size 5242880 --max-file-count 500 ./docs/
# Output: Files processed with 5MB and 500 file limits

# Process summary
Processed: 5 successful, 0 errors
```

## Using as a Library

The `lib` package can be used as a dependency in your own Go projects. Only the library code will be compiled into your project—CLI tool, web server, and test dependencies are excluded.

### Installation

```bash
go get github.com/yourusername/asciidoc-xml/lib
```

### Quick Example

```go
import (
    "github.com/ndx-video/asciidoc-xml/lib"
    "bytes"
    "os"
)

// Convert AsciiDoc to HTML5 (with PicoCSS enabled by default)
html, err := lib.ConvertToHTML(strings.NewReader(asciidoc), false, true, "", "")

// Convert to XML
xml, err := lib.ConvertToXML(strings.NewReader(asciidoc))

// Convert Markdown to AsciiDoc
mdContent, _ := os.ReadFile("document.md")
asciidoc, err := lib.ConvertMarkdownToAsciiDoc(bytes.NewReader(mdContent))

// Convert Markdown to AsciiDoc (streaming for large files)
var output bytes.Buffer
err = lib.ConvertMarkdownToAsciiDocStreaming(
    bytes.NewReader(mdContent),
    &output,
)
```

### What Gets Included?

✅ **Included:**
- Only the `lib` package (parser, converter, DOM)
- Standard library dependencies only

❌ **Excluded:**
- CLI tool (`cli/` package)
- Web server (`web/` package)  
- Test files and test-only dependencies

The `lib` package has **zero external runtime dependencies**—only Go standard library.

See [Library Usage Guide](docs/library-usage.md) for detailed examples and API reference.

## Usage

### Basic Conversion

```go
package main

import (
    "bytes"
    "fmt"
    "os"
    "github.com/ndx-video/asciidoc-xml/lib"
)

func main() {
    asciidoc := `= My Document
    
This is a paragraph.`
    
    // Convert to XML
    xml, err := lib.ConvertToXML(bytes.NewReader([]byte(asciidoc)))
    if err != nil {
        panic(err)
    }
    
    // Convert to HTML5 (with PicoCSS enabled by default)
    html, err := lib.ConvertToHTML(bytes.NewReader([]byte(asciidoc)), false, true, "", "")
    if err != nil {
        panic(err)
    }
    
    // Convert to XHTML5 (with PicoCSS enabled by default)
    xhtml, err := lib.ConvertToHTML(bytes.NewReader([]byte(asciidoc)), true, true, "", "")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(xml)
    fmt.Println(html)
    fmt.Println(xhtml)
}
```

### Markdown to AsciiDoc Conversion

```go
package main

import (
    "bytes"
    "fmt"
    "os"
    "github.com/ndx-video/asciidoc-xml/lib"
)

func main() {
    // Read Markdown file
    mdContent, err := os.ReadFile("document.md")
    if err != nil {
        panic(err)
    }
    
    // Convert to AsciiDoc (non-streaming)
    asciidoc, err := lib.ConvertMarkdownToAsciiDoc(bytes.NewReader(mdContent))
    if err != nil {
        panic(err)
    }
    
    fmt.Println(asciidoc)
    
    // Or use streaming for large files (memory-efficient)
    var output bytes.Buffer
    err = lib.ConvertMarkdownToAsciiDocStreaming(
        bytes.NewReader(mdContent),
        &output,
    )
    if err != nil {
        panic(err)
    }
    
    // Write output to file
    os.WriteFile("document.adoc", output.Bytes(), 0644)
}
```

### Batch Processing

```go
package main

import (
    "fmt"
    "github.com/ndx-video/asciidoc-xml/lib"
)

func main() {
    files := []string{"file1.adoc", "file2.adoc", "file3.adoc"}
    
    config := lib.BatchConfig{
        MaxWorkers:        4,
        ParallelThreshold: 2,
        EnableParallel:    true,
        DryRun:            false,
        ValidateOnly:      false,
    }
    
    limits := lib.ProcessingLimits{
        MaxFileSize:    10 * 1024 * 1024, // 10MB
        MaxArchiveSize: 10 * 1024 * 1024, // 10MB
        MaxFileCount:   10000,
    }
    
    results := lib.ProcessFilesParallel(
        files,
        func(file string) error {
            // Process each file
            return processFile(file)
        },
        config,
        limits,
        func(current, total int, file string, err error) {
            // Progress callback
            fmt.Printf("Processing %d/%d: %s\n", current, total, file)
        },
    )
    
    fmt.Printf("Processed: %d successful, %d errors\n", 
        results.SuccessCount, results.ErrorCount)
}
```

### Programmatic Access

```go
doc, err := lib.Convert(bytes.NewReader([]byte(asciidoc)))
if err != nil {
    panic(err)
}

// Access document structure
fmt.Println("DocType:", doc.DocType)
fmt.Println("Title:", doc.Header.Title)
```

### XSLT Transformation

The generated XML can be transformed to HTML using the provided XSLT template:

```bash
xsltproc xslt/asciidoc-to-html.xsl document.xml > document.html
```

Or programmatically using an XSLT processor library. The web harness uses the browser's built-in XSLT processor for live preview.

## Testing

The project includes comprehensive test coverage:

### Go Tests

```bash
# Run all Go tests
go test ./...

# Run tests for specific package
go test ./lib/...
go test ./web/...
go test ./cli/...

# Run testbed test suite (all Markdown files in testbed/)
go test ./lib -run TestConvertMarkdownToAsciiDoc_Testbed

# Run comprehensive testbed (all subfolders)
go test ./lib -run TestConvertMarkdownToAsciiDoc_TestbedAll
```

### Testbed Suite

The project includes a comprehensive testbed (`./testbed/`) with:
- Markdown test files covering all CommonMark and GFM features
- Corrupt file tests for parser resilience
- Multi-language test files
- Edge case coverage

Run the testbed suite:
```bash
# Run testbed.sh script
./testbed.sh

# Run with all subfolders
./testbed.sh -all

# Run with main logging system
./testbed.sh -log
```

### JavaScript Tests

JavaScript files are tested using Goja (ECMAScript 5.1+ implementation in Go):

```bash
# JavaScript tests run automatically with Go tests
go test ./web/static/...

# Or run all tests via harness
./harness.sh test-all
```

The JavaScript test suite validates:
- Syntax correctness of `app.js` and `pretty.js`
- Function existence and callability
- Syntax highlighting functionality
- Source view generation

## XML Schema

The custom XSD schema (`schema/asciidoc.xsd`) defines a comprehensive, purpose-built structure for AsciiDoc documents. It includes:

**Document Structure:**
- Document root with attributes (id, role, doctype)
- Preamble support
- Header with title, authors, revision, attributes

**Block Elements:**
- Sections (levels 0-5) with id, role, appendix, discrete attributes
- Paragraphs with id and role
- Code blocks and literal blocks with language, id, role
- Example blocks, sidebars, quotes with id and role
- Verse blocks and open blocks (generic containers)
- Admonitions (note, tip, warning, caution, important)

**Lists:**
- Ordered and unordered lists with id and role
- List items with callout, role, id, term attributes
- List continuations

**Tables:**
- Tables with id and role
- Table rows with role (header/footer)
- Table cells with align, colspan, rowspan, role, id attributes

**Inline Elements:**
- Bold, italic, monospace
- Superscript, subscript, highlight
- Links with href, title, class, window (target) attributes
- Inline macros (kbd, btn, menu, generic)
- Anchors and cross-references
- Footnotes and footnote references
- Passthrough (CDATA)

**Special Elements:**
- Thematic breaks and page breaks
- Block macros (include, TOC, video, audio, etc.)
- Images with alt, width, height, link attributes

All AsciiDoc features are represented as XML elements with appropriate attributes for configuration options. Content is stored in text nodes, while metadata and features are stored as attributes to keep XML lean and avoid unnecessary nesting.

## XSLT Template

The XSLT template (`xslt/asciidoc-to-html.xsl`) transforms the XML to semantic HTML with CSS classes:

- `asciidoc-document` - Document container
- `asciidoc-section` - Sections with level classes
- `asciidoc-paragraph` - Paragraphs
- `asciidoc-codeblock` - Code blocks
- `asciidoc-table` - Tables
- `asciidoc-list` - Lists
- And more...

All elements include semantic classes that can be styled with CSS (e.g., Tailwind CSS).

## Web Development Harness

The web harness (`web/`) is a Single Page Application (SPA) that provides:

### Features

- **4-Column Layout**: 
  - AsciiDoc source (read-only, syntax highlighted)
  - Generated XML (syntax highlighted)
  - XSLT template (syntax highlighted)
  - HTML output (rendered + source view)

- **Interactive Features**:
  - Resizable columns (drag borders)
  - File upload (AsciiDoc and XSLT)
  - Path-based file loading
  - Auto-conversion on startup
  - Live XSLT transformation

- **Syntax Highlighting**:
  - Custom AsciiDoc highlighter with color coding
  - XML/XSLT highlighting
  - HTML highlighting

### API Endpoints

See the complete [API Documentation](docs/api.adoc) for detailed endpoint documentation, request/response formats, and examples.

- `GET /` - Main SPA
- `POST /api/convert` - Convert AsciiDoc to XML, HTML, or XHTML (supports `output` parameter: "xml", "html", "xhtml", "md2adoc")
- `POST /api/validate` - Validate AsciiDoc syntax
- `GET /api/xslt` - Get XSLT template
- `POST /api/upload` - Upload AsciiDoc, Markdown, or XSLT file
- `GET /api/load-file?path=...` - Load file from server path
- `GET /docs` - User guide documentation (generated from AsciiDoc)
- `POST /api/batch` - Batch process multiple files (supports configuration JSON body)

## Markdown to AsciiDoc Conversion

The library provides comprehensive Markdown to AsciiDoc conversion with support for:

### CommonMark Features
- Setext-style headers (`===` and `---`)
- Reference-style links and images (`[text][ref]`)
- Link reference definitions (`[ref]: url "title"`)
- Autolinks (`<url>` and `<email@example.com>`)
- Hard line breaks (two spaces + newline)
- Indented code blocks (4+ spaces)
- Escaped characters (backslash escaping)
- YAML frontmatter (arrays, dictionaries, nested structures)

### GitHub Flavored Markdown (GFM)
- Task lists (`- [ ]` and `- [x]`)
- Strikethrough (`~~text~~`)
- Table alignment (left, center, right)
- HTML blocks and spans (with passthrough, GFM-specific)

### Usage

**CLI:**
```bash
# Convert Markdown to AsciiDoc
./adc --output md2adoc document.md

# Convert Markdown to AsciiDoc, then to XML
./adc --output md2adoc document.md
./adc document.adoc
```

**Library:**
```go
// Non-streaming conversion
asciidoc, err := lib.ConvertMarkdownToAsciiDoc(markdownReader)

// Streaming conversion (memory-efficient for large files)
err := lib.ConvertMarkdownToAsciiDocStreaming(markdownReader, asciidocWriter)
```

**Web API:**
```bash
curl -X POST http://localhost:8005/api/convert \
  -H "Content-Type: application/json" \
  -d '{"asciidoc": "# Markdown content", "output": "md2adoc"}'
```

## Enhanced AsciiDoc Parser Features

The AsciiDoc parser has been significantly enhanced with support for:

### Inline Formatting
- **Superscript**: `^text^`
- **Subscript**: `~text~`
- **Highlight**: `#text#`
- **Inline passthrough**: `pass:[content]`

### Inline Macros
- **Keyboard**: `kbd:[Ctrl+C]`
- **Button**: `btn:[Save]`
- **Menu**: `menu:File[New]`
- **Generic**: `macro:name[target,attributes]`

### Cross-References and Anchors
- **Block anchors**: `[[anchor-id]]` or `[#anchor-id]`
- **Section ID generation**: Automatic ID generation from titles
- **Cross-references**: `<<anchor-id>>` or `xref:anchor-id[]`
- **Anchor registry**: Tracks all anchors for resolution

### Attribute Substitution
- **Document attributes**: `:attr-name: value`
- **Header attributes**: Custom attributes in document header
- **Built-in attributes**: `{author}`, `{revnumber}`, etc.
- **Attribute substitution**: `{attr-name}` in content

### List Enhancements
- **List continuations**: `+` for continuing list items
- **Callout lists**: `<1>`, `<2>`, etc. for code block callouts
- **List item attributes**: `[.class]#item#` syntax

### Block Enhancements
- **Verse blocks**: `[verse]____...____` with attribution
- **Open blocks**: `--` generic containers
- **Enhanced attributes**: `[#id.role]` syntax for all blocks

### Tables
- **Cell alignment**: Left, center, right
- **Cell spanning**: Colspan and rowspan
- **Row roles**: Header and footer rows
- **Cell attributes**: ID, role, alignment

### Footnotes
- **Inline footnotes**: `footnote:[text]`
- **Reference footnotes**: `footnote:id[]` and `footnoteref:id[]`
- **Automatic numbering**: When no ID provided

### Standard Block Macros
- **Include**: `include::file.adoc[]`
- **Table of Contents**: `toc::[]`
- **Video**: `video::url[]`
- **Audio**: `audio::url[]`

## Examples

See `examples/comprehensive.adoc` for a complete example demonstrating all AsciiDoc features.

See `testbed/*.md` for comprehensive Markdown examples covering all supported features.

## Logging

The library includes a comprehensive logging system using only Go's standard library:

### Features
- **Log Levels**: DEBUG, INFO, WARN, ERROR, FATAL
- **Structured Logging**: Text and JSON formats
- **File Logging**: With automatic rotation (size-based, daily, hourly)
- **Console Logging**: Configurable console output with separate log levels
- **Request Tracking**: Context-based request IDs for web API
- **HTTP Middleware**: Automatic request/response logging for web server

### Configuration

Configure logging via `adc.json`:

```json
{
  "logging": {
    "level": "info",
    "format": "text",
    "file": {
      "enabled": true,
      "path": "./logs",
      "filename": "adc.log",
      "maxSize": 10485760,
      "maxFiles": 5,
      "rotation": "size"
    },
    "console": {
      "enabled": true,
      "level": "info"
    }
  }
}
```

### Usage

```go
import "github.com/ndx-video/asciidoc-xml/lib"

// Create logger
config := lib.LogConfig{
    Level:  "info",
    Format: "text",
    Console: lib.ConsoleConfig{
        Enabled: true,
        Level:   "info",
    },
    File: lib.FileConfig{
        Enabled: true,
        Path:     "./logs",
        Filename: "app.log",
        MaxSize:  10 * 1024 * 1024,
        MaxFiles: 5,
        Rotation: "size",
    },
}

logger, err := lib.NewLogger(config)
if err != nil {
    panic(err)
}
defer logger.Close()

// Use logger
logger.Debug(nil, "Debug message", "key", "value")
logger.Info(nil, "Info message", "key", "value")
logger.Warn(nil, "Warning message", "key", "value")
logger.Error(nil, "Error message", "key", "value")

// With context (for request tracking)
ctx := context.WithValue(context.Background(), "request_id", "abc123")
logger.Info(ctx, "Request processed", "method", "POST", "path", "/api/convert")
```

## Dependencies

### Runtime Dependencies

None! This is a pure Go implementation with zero external runtime dependencies for production use.

### Test Dependencies

- `github.com/dop251/goja` - JavaScript runtime for testing JavaScript files (test-only, not included in production binaries)

The `goja` dependency is only used in test files and is automatically excluded from production builds using Go build constraints (`//go:build test`).

## Development Workflow

1. **Start the development server**:
   ```bash
   ./harness.sh start
   ```

2. **Open browser** to `http://localhost:8005`

3. **Load an AsciiDoc file**:
   - Use the path input field below "AsciiDoc Source"
   - Or upload a file via the Upload button
   - Or the example loads automatically

4. **View results**:
   - XML is generated automatically
   - HTML is transformed using the XSLT template
   - All columns are syntax highlighted

5. **Edit XSLT**:
   - Load XSLT from path or upload
   - Changes are reflected immediately in HTML output

6. **Test changes**:
   ```bash
   ./harness.sh test-all
   ```

## Customization

### Custom XSLT Templates

You can create custom XSLT templates for different output formats:

1. Copy `xslt/asciidoc-to-html.xsl`
2. Modify the templates to match your HTML structure
3. Use your custom template for transformation
4. Test in the web harness by loading your custom XSLT

### Extending the Schema

To add new features:

1. Update `schema/asciidoc.xsd`
2. Update Go structs in `xml.go`
3. Update parser in `lib/adoc-parser.go`
4. Update converter in `lib/converter.go`
5. Update XSLT template (if using XML pipeline)
6. Add tests

## License

[Your License Here]

## Contributing

[Contributing Guidelines]
