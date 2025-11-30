# CMS Developer Agent: Upgrade Guide for asciidoc-xml v0.4.0

## Overview

The `github.com/ndx-video/asciidoc-xml` library has been upgraded to **v0.4.0**, introducing comprehensive Markdown conversion capabilities, enhanced logging, batch processing, and archive support. This guide will help you upgrade your Go dependency and take advantage of the new features.

## Upgrade Instructions

### 1. Update Go Module Dependency

```bash
go get github.com/ndx-video/asciidoc-xml@v0.4.0
go mod tidy
```

### 2. Verify Version

```go
import "github.com/ndx-video/asciidoc-xml/lib"

fmt.Println("Library version:", lib.Version) // Should output "0.4.0"
```

## New Features Summary

### 1. Comprehensive Markdown to AsciiDoc Conversion

**Previous Version**: Basic Markdown conversion with limited features  
**New Version**: Full CommonMark specification + GitHub Flavored Markdown (GFM) support

#### New Conversion Function

```go
import (
    "github.com/ndx-video/asciidoc-xml/lib"
    "bytes"
    "io"
)

// Streaming conversion (memory-efficient for large files)
func ConvertMarkdownToAsciiDocStreaming(reader io.Reader, writer io.Writer) error

// Non-streaming conversion (backward compatible)
func ConvertMarkdownToAsciiDoc(reader io.Reader) (string, error)
```

#### Supported Markdown Features

**CommonMark Features:**
- ✅ Setext-style headers (`===` and `---`)
- ✅ Reference-style links and images (`[text][ref]`)
- ✅ Link reference definitions (`[ref]: url "title"`)
- ✅ Autolinks (`<url>` and `<email@example.com>`)
- ✅ Hard line breaks (two spaces + newline)
- ✅ Indented code blocks (4+ spaces)
- ✅ Escaped characters (backslash escaping)

**GitHub Flavored Markdown (GFM):**
- ✅ Task lists (`- [ ]` and `- [x]`)
- ✅ Strikethrough (`~~text~~`)
- ✅ Table alignment (left, center, right)
- ✅ HTML blocks and spans (with passthrough, GFM-specific)

#### Usage Example

```go
package main

import (
    "bytes"
    "fmt"
    "github.com/ndx-video/asciidoc-xml/lib"
    "os"
)

func main() {
    // Read Markdown file
    mdContent, err := os.ReadFile("document.md")
    if err != nil {
        panic(err)
    }
    
    // Convert to AsciiDoc
    asciidoc, err := lib.ConvertMarkdownToAsciiDoc(bytes.NewReader(mdContent))
    if err != nil {
        panic(err)
    }
    
    fmt.Println(asciidoc)
    
    // Or use streaming for large files
    var output bytes.Buffer
    err = lib.ConvertMarkdownToAsciiDocStreaming(
        bytes.NewReader(mdContent),
        &output,
    )
    if err != nil {
        panic(err)
    }
    
    fmt.Println(output.String())
}
```

### 2. Comprehensive Logging System

**New Package**: `lib.Logger` with structured logging support

#### Features

- Log levels: DEBUG, INFO, WARN, ERROR, FATAL
- Structured logging (text and JSON formats)
- File logging with size-based rotation
- Context-based request IDs
- Thread-safe concurrent logging

#### Usage Example

```go
import "github.com/ndx-video/asciidoc-xml/lib"

// Create logger configuration
config := lib.LogConfig{
    Level:  "info",
    Format: "text", // or "json"
    File: lib.FileLogConfig{
        Enabled:  true,
        Path:     "./logs",
        Filename: "app.log",
        MaxSize:  10 * 1024 * 1024, // 10MB
        MaxFiles: 5,
        Rotation: "size",
    },
    Console: lib.ConsoleLogConfig{
        Enabled: true,
        Level:   "info",
    },
}

// Initialize logger
logger, err := lib.NewLogger(config)
if err != nil {
    panic(err)
}
defer logger.Close()

// Use logger
ctx := context.Background()
logger.Info(ctx, "Processing file", "file", "document.md", "size", 1024)
logger.Error(ctx, "Conversion failed", "error", err.Error())
```

### 3. Batch Processing with Parallelization

**New Package**: `lib.ProcessFilesParallel` for efficient batch operations

#### Features

- Parallel file processing with worker pool
- Configurable worker count and thresholds
- File size and count limits
- Dry-run and validation-only modes
- Directory structure preservation

#### Usage Example

```go
import "github.com/ndx-video/asciidoc-xml/lib"

// Configure batch processing
config := lib.BatchConfig{
    MaxWorkers:        4,
    ParallelThreshold: 2,
    EnableParallel:    true,
    DryRun:           false,
    ValidateOnly:     false,
    PreserveStructure: true,
}

limits := lib.ProcessingLimits{
    MaxFileSize:    10 * 1024 * 1024, // 10MB
    MaxArchiveSize: 10 * 1024 * 1024, // 10MB
    MaxFileCount:   10000,
}

// Process files
files := []string{"file1.md", "file2.md", "file3.md"}

results, err := lib.ProcessFilesParallel(
    files,
    config,
    limits,
    func(filePath string) error {
        // Your conversion logic here
        content, _ := os.ReadFile(filePath)
        asciidoc, _ := lib.ConvertMarkdownToAsciiDoc(bytes.NewReader(content))
        // Save asciidoc...
        return nil
    },
    nil, // optional logger
)
```

### 4. Archive Extraction Support

**New Package**: `lib.ExtractArchive` for handling compressed files

#### Supported Formats

- ZIP (`.zip`)
- TAR (`.tar`)
- TAR.GZ (`.tar.gz`, `.tgz`)

#### Features

- Path traversal security checks
- Size limits
- Sequential extraction
- Standard library only (no external dependencies)

#### Usage Example

```go
import "github.com/ndx-video/asciidoc-xml/lib"

// Extract archive
err := lib.ExtractArchive("archive.zip", "/tmp/extracted", 10*1024*1024)
if err != nil {
    panic(err)
}

// Detect archive format
format := lib.DetectArchiveFormat("archive.tar.gz")
fmt.Println("Format:", format) // "tar.gz"
```

## Migration Guide

### No Breaking Changes

**Good News**: v0.4.0 maintains backward compatibility. All existing code will continue to work without changes.

### Optional Upgrades

#### 1. Add Logging to Existing Code

If you want to add logging to your CMS:

```go
// Before
func processFile(file string) error {
    // ... conversion logic
    return nil
}

// After (with logging)
func processFile(file string, logger *lib.Logger) error {
    ctx := context.Background()
    logger.Info(ctx, "Processing file", "file", file)
    
    // ... conversion logic
    
    logger.Info(ctx, "File processed", "file", file, "size", size)
    return nil
}
```

#### 2. Use Streaming for Large Files

If you're processing large Markdown files:

```go
// Before (loads entire file into memory)
content, _ := os.ReadFile("large.md")
asciidoc, _ := lib.ConvertMarkdownToAsciiDoc(bytes.NewReader(content))

// After (streaming, memory-efficient)
file, _ := os.Open("large.md")
defer file.Close()

var output bytes.Buffer
lib.ConvertMarkdownToAsciiDocStreaming(file, &output)
```

#### 3. Enable Batch Processing

If you're processing multiple files:

```go
// Before (sequential)
for _, file := range files {
    processFile(file)
}

// After (parallel with limits)
config := lib.BatchConfig{
    MaxWorkers: 4,
    EnableParallel: true,
}
limits := lib.ProcessingLimits{
    MaxFileSize: 10 * 1024 * 1024,
    MaxFileCount: 10000,
}

lib.ProcessFilesParallel(files, config, limits, processFile, logger)
```

## New Dependencies

**None!** All new features use only Go's standard library. No additional dependencies required.

## Testing

The library now includes comprehensive test coverage:

- 150+ Markdown test files in `testbed/`
- Unit tests for all new components
- Integration tests for batch processing
- Security tests for archive extraction

## API Reference

### New Exported Types

```go
// Logging
type Logger struct { ... }
type LogConfig struct { ... }
type LogLevel int

// Batch Processing
type BatchConfig struct { ... }
type ProcessingLimits struct { ... }

// Archive
type ArchiveFormat string
```

### New Exported Functions

```go
// Markdown Conversion
func ConvertMarkdownToAsciiDocStreaming(reader io.Reader, writer io.Writer) error
func ConvertMarkdownToAsciiDoc(reader io.Reader) (string, error)

// Logging
func NewLogger(config LogConfig) (*Logger, error)

// Batch Processing
func ProcessFilesParallel(
    files []string,
    config BatchConfig,
    limits ProcessingLimits,
    processor func(string) error,
    logger *Logger,
) ([]ProcessingResult, error)

// Archive
func ExtractArchive(archivePath, destDir string, maxSize int64) error
func DetectArchiveFormat(filename string) ArchiveFormat
func CreateArchive(files []string, archivePath string, format ArchiveFormat) error
```

## Recommended Upgrade Path

1. **Update dependency**: `go get github.com/ndx-video/asciidoc-xml@v0.4.0`
2. **Test existing code**: Verify all existing functionality still works
3. **Add logging** (optional): Integrate logger for better observability
4. **Enable batch processing** (optional): Use parallel processing for multiple files
5. **Use streaming** (optional): Switch to streaming conversion for large files

## Support

For questions or issues:
- Check the documentation in `docs/` directory
- Review test files in `testbed/` for usage examples
- See `docs/markdown-asciidoc-limitations.md` for feature details

## Summary

**Version**: 0.4.0  
**Breaking Changes**: None  
**New Features**: 
- Full CommonMark + GFM Markdown support
- Comprehensive logging system
- Batch processing with parallelization
- Archive extraction support
- Enhanced test coverage

**Action Required**: Update dependency, optionally integrate new features.

