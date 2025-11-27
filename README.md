# AsciiDoc XML Converter

A pure Go package for converting AsciiDoc documents to a custom XML format, designed for easy transformation to HTML via XSLT. Includes a web-based development harness for testing XML generation and XSLT template development.

## Overview

This package provides:

- **Custom XML Schema (XSD)**: A purpose-built XML schema specifically designed for AsciiDoc, avoiding the bloat of DocBook
- **Go XML Structures**: Type-safe Go structs matching the XSD schema
- **AsciiDoc Parser**: Pure Go parser that converts AsciiDoc source to XML
- **XSLT Template**: Comprehensive XSLT stylesheet for transforming XML to HTML
- **Web Development Harness**: Single Page Application (SPA) for interactive development and testing

## Features

- ✅ Pure Go implementation (no external binaries required)
- ✅ Cross-platform support (no architecture-specific dependencies)
- ✅ Comprehensive AsciiDoc feature support
- ✅ Well-formed, validatable XML output
- ✅ Semantic HTML output via XSLT
- ✅ Web-based development harness with live preview
- ✅ Command line tool for batch conversion
- ✅ JavaScript syntax highlighting and pretty-printing
- ✅ Comprehensive test suite (Go + JavaScript)
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
│   └── converter_test.go     # Converter tests
├── docs/
│   ├── asciidoc-xml.adoc     # User guide documentation
│   └── api.adoc              # Web API documentation
├── cli/
│   ├── adc.go                # Command line tool (AsciiDoc Converter)
│   └── adc_test.go           # CLI tests
├── web/
│   ├── main.go               # Web server
│   ├── main_test.go          # Server tests
│   ├── static/
│   │   ├── app.js            # Main application logic
│   │   ├── app.css           # Styles
│   │   ├── pretty.js         # Syntax highlighting
│   │   ├── js_test.go        # JavaScript tests (using Goja)
│   │   └── comprehensive.adoc # Example file
│   └── templates/
│       └── index.html        # SPA template
├── xslt/
│   └── asciidoc-to-html.xsl  # XSLT transformation template
├── examples/
│   └── comprehensive.adoc    # Example file with all features
├── harness.sh                # Development server manager
├── Makefile                  # Build and distribution automation
├── xml.go                    # Go XML struct definitions
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

**Full packages** (includes both CLI tool and web server):
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

# Build CLI for all target platforms
make build-cli

# Build web server for all target platforms
make build-web

# Build both CLI and web for all platforms
make build-all

# Create CLI-only distribution packages
make dist-cli

# Create full distribution packages (CLI + web)
make dist-full VERSION=1.0.0

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

The project provides two distribution package types:

### CLI-only Package

Contains only the `adc` command-line tool. Ideal for users who only need batch conversion functionality.

**Contents:**
- `bin/adc` - Command-line converter
- `LICENSE` - License file
- `README.md` - Documentation
- `examples/` - Example AsciiDoc files

### Full Package

Contains both the CLI tool and web server, plus XSLT templates. Ideal for users who want the complete feature set including the web interface.

**Contents:**
- `bin/adc` - Command-line converter
- `bin/asciidoc-xml-web` - Web server
- `xslt/` - XSLT transformation templates
- `LICENSE` - License file
- `README.md` - Documentation
- `examples/` - Example AsciiDoc files

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

The `adc` (AsciiDoc Converter) command line tool provides batch conversion of AsciiDoc files to XML and optionally to HTML via XSLT transformation.

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

# Generate XML only (skip XSLT transformation)
./adc --no-xsl document.adoc

# Specify a custom XSLT file
./adc --xsl xslt/asciidoc-to-html.xsl document.adoc

# Auto-overwrite existing files without prompting
./adc -y document.adoc
```

#### Command Options

- `-y`: Automatically overwrite existing XML/HTML files without prompting
- `--no-xsl`: Generate XML only, skip XSLT transformation to HTML
- `--xsl <path>`: Path to XSLT file (default: `./default.xsl`)

#### File Processing

- **Single File**: Processes one `.adoc` file, outputs `.xml` (and optionally `.html`) in the same directory
- **Directory**: Recursively traverses directory, processes all `.adoc` files found
- **Output Files**: 
  - XML: `filename.xml` (same name as input, `.xml` extension)
  - HTML: `filename.html` (generated when XSLT transformation is applied)

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
# Convert single file with default XSLT
./adc document.adoc
# Output: document.xml, document.html

# Convert directory, XML only
./adc --no-xsl ./docs/
# Output: All .xml files, no HTML

# Batch convert with custom XSLT, auto-overwrite
./adc -y --xsl custom.xsl ./docs/
# Output: All .xml and .html files

# Process summary
Processed: 5 successful, 0 errors
```

## Usage

### Basic Conversion

```go
package main

import (
    "bytes"
    "fmt"
    "asciidoc-xml/lib"
)

func main() {
    asciidoc := `= My Document
    
This is a paragraph.`
    
    // Convert to XML
    xml, err := lib.ConvertToXML(bytes.NewReader([]byte(asciidoc)))
    if err != nil {
        panic(err)
    }
    
    // Convert to HTML5
    html, err := lib.ConvertToHTML(bytes.NewReader([]byte(asciidoc)), false)
    if err != nil {
        panic(err)
    }
    
    // Convert to XHTML5
    xhtml, err := lib.ConvertToHTML(bytes.NewReader([]byte(asciidoc)), true)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(xml)
    fmt.Println(html)
    fmt.Println(xhtml)
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
go test ./converter/...
go test ./web/...
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

The custom XSD schema (`schema/asciidoc.xsd`) defines a minimal, purpose-built structure for AsciiDoc documents. It includes:

- Document structure (header, content)
- Sections (levels 0-5)
- Text blocks (paragraphs, code blocks, literal blocks)
- Lists (ordered, unordered, labeled, callout)
- Tables
- Media (images, video, audio)
- Admonitions
- Inline formatting
- And more...

All AsciiDoc features are represented as XML elements with appropriate attributes for configuration options.

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
- `POST /api/convert` - Convert AsciiDoc to XML, HTML, or XHTML (supports `output` parameter: "xml", "html", "xhtml")
- `POST /api/validate` - Validate AsciiDoc syntax
- `GET /api/xslt` - Get XSLT template
- `POST /api/upload` - Upload AsciiDoc or XSLT file
- `GET /api/load-file?path=...` - Load file from server path
- `GET /docs` - User guide documentation (generated from AsciiDoc)

## Example

See `examples/comprehensive.adoc` or `web/static/comprehensive.adoc` for a complete example demonstrating all AsciiDoc features.

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
