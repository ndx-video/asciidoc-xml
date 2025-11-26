# AsciiDoc XML Converter

A pure Go package for converting AsciiDoc documents to a custom XML format, designed for easy transformation to HTML via XSLT. Includes a web-based development harness for testing XML generation and XSLT template development.

## Overview

This package provides:

- **Custom XML Schema (XSD)**: A purpose-built XML schema specifically designed for AsciiDoc, avoiding the bloat of DocBook
- **Go XML Structures**: Type-safe Go structs matching the XSD schema
- **AsciiDoc Parser**: Converts AsciiDoc source to XML using `libasciidoc`
- **XSLT Template**: Comprehensive XSLT stylesheet for transforming XML to HTML
- **Web Development Harness**: Single Page Application (SPA) for interactive development and testing

## Features

- ✅ Pure Go implementation (no external binaries required)
- ✅ Cross-platform support (no architecture-specific dependencies)
- ✅ Comprehensive AsciiDoc feature support
- ✅ Well-formed, validatable XML output
- ✅ Semantic HTML output via XSLT
- ✅ Web-based development harness with live preview
- ✅ JavaScript syntax highlighting and pretty-printing
- ✅ Comprehensive test suite (Go + JavaScript)
- ✅ Extensible and customizable

## Package Structure

```
asciidoc-xml/
├── schema/
│   └── asciidoc.xsd          # XML Schema Definition
├── converter/
│   ├── converter.go          # AsciiDoc to XML converter
│   └── converter_test.go     # Converter tests
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
├── xml.go                    # Go XML struct definitions
└── README.md                 # This file
```

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
- **4-column layout**: AsciiDoc source, Generated XML, XSLT template, HTML output
- **Resizable columns**: Drag column borders to adjust widths
- **Syntax highlighting**: Color-coded AsciiDoc, XML, XSLT, and HTML
- **Live conversion**: Automatic XML generation when AsciiDoc loads
- **File upload**: Upload AsciiDoc and XSLT files via the web interface
- **Path loading**: Load files from server paths

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

## Usage

### Basic Conversion

```go
package main

import (
    "bytes"
    "fmt"
    "asciidoc-xml/converter"
)

func main() {
    asciidoc := `= My Document
    
This is a paragraph.`
    
    xml, err := converter.ConvertToXML(bytes.NewReader([]byte(asciidoc)))
    if err != nil {
        panic(err)
    }
    
    fmt.Println(xml)
}
```

### Programmatic Access

```go
doc, err := converter.Convert(bytes.NewReader([]byte(asciidoc)))
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

- `GET /` - Main SPA
- `POST /api/convert` - Convert AsciiDoc to XML
- `POST /api/validate` - Validate AsciiDoc syntax
- `GET /api/xslt` - Get XSLT template
- `POST /api/upload` - Upload AsciiDoc or XSLT file
- `GET /api/load-file?path=...` - Load file from server path

## Example

See `examples/comprehensive.adoc` or `web/static/comprehensive.adoc` for a complete example demonstrating all AsciiDoc features.

## Dependencies

- `github.com/bytesparadise/libasciidoc` - AsciiDoc parser and validation
- `github.com/sirupsen/logrus` - Logging (for silencing libasciidoc logs)
- `github.com/dop251/goja` - JavaScript interpreter for testing JS files

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
3. Update converter in `converter/converter.go`
4. Update XSLT template
5. Add tests

## License

[Your License Here]

## Contributing

[Contributing Guidelines]
