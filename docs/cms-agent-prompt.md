# Prompt for CMS Developer Agent: asciidoc-xml v0.4.0 Upgrade

## Task

Upgrade the Go dependency `github.com/ndx-video/asciidoc-xml` from the previous version to **v0.4.0** and integrate the new features into the CMS codebase.

## Current Situation

The asciidoc-xml library has been significantly enhanced with new features. Your CMS currently uses an older version and should be upgraded to take advantage of:

1. **Comprehensive Markdown Conversion**: Full CommonMark specification and GitHub Flavored Markdown (GFM) support
2. **Logging System**: Structured logging with file rotation and context tracking
3. **Batch Processing**: Parallel file processing with configurable limits
4. **Archive Support**: Extract and process ZIP, TAR, and TAR.GZ archives

## Upgrade Steps

### Step 1: Update Dependency

```bash
go get github.com/ndx-video/asciidoc-xml@v0.4.0
go mod tidy
```

### Step 2: Verify No Breaking Changes

The upgrade is **backward compatible**. All existing code should continue to work. Test your current Markdown conversion functionality to ensure it still works correctly.

### Step 3: Review New Features

#### A. Enhanced Markdown Conversion

The library now supports:
- Setext-style headers (`===` and `---`)
- Reference-style links (`[text][ref]`)
- Autolinks (`<url>` and `<email@example.com>`)
- Hard line breaks (two spaces + newline)
- Indented code blocks (4+ spaces)
- Task lists (`- [ ]` and `- [x]`)
- Strikethrough (`~~text~~`)
- Table alignment
- HTML blocks/spans (GFM-specific)

**Action**: No code changes needed - these features work automatically. However, you may want to update your CMS documentation to inform users about the expanded Markdown support.

#### B. Logging System

New structured logging with:
- Log levels (DEBUG, INFO, WARN, ERROR, FATAL)
- Text and JSON formats
- File rotation
- Context-based request IDs

**Action**: Consider integrating the logger into your CMS for better observability:

```go
import "github.com/ndx-video/asciidoc-xml/lib"

// Initialize logger
config := lib.LogConfig{
    Level:  "info",
    Format: "json",
    File: lib.FileLogConfig{
        Enabled:  true,
        Path:     "./logs",
        Filename: "cms.log",
        MaxSize:  10 * 1024 * 1024,
        MaxFiles: 5,
    },
}
logger, err := lib.NewLogger(config)
// Use logger in your conversion functions
```

#### C. Batch Processing

Parallel file processing with worker pools:

**Action**: If your CMS processes multiple files, consider using batch processing:

```go
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

#### D. Archive Extraction

Support for ZIP, TAR, TAR.GZ archives:

**Action**: If your CMS accepts archive uploads, use the new extraction function:

```go
err := lib.ExtractArchive("upload.zip", "/tmp/extracted", 10*1024*1024)
```

### Step 4: Use Streaming for Large Files (Optional)

If you process large Markdown files, switch to streaming conversion:

```go
// Before
content, _ := os.ReadFile("large.md")
asciidoc, _ := lib.ConvertMarkdownToAsciiDoc(bytes.NewReader(content))

// After (memory-efficient)
file, _ := os.Open("large.md")
defer file.Close()
var output bytes.Buffer
lib.ConvertMarkdownToAsciiDocStreaming(file, &output)
```

## Testing Checklist

After upgrading, verify:

- [ ] Existing Markdown conversion still works
- [ ] New Markdown features (task lists, strikethrough, etc.) are supported
- [ ] No compilation errors
- [ ] All tests pass
- [ ] Logging works (if integrated)
- [ ] Batch processing works (if integrated)
- [ ] Archive extraction works (if integrated)

## Documentation

- See `docs/cms-upgrade-guide-v0.4.0.md` for detailed API reference
- See `docs/markdown-asciidoc-limitations.md` for feature limitations
- See `testbed/` directory for Markdown test examples

## Questions to Consider

1. **Do you process multiple files?** → Consider batch processing
2. **Do you handle large files?** → Use streaming conversion
3. **Do you need better observability?** → Integrate logging
4. **Do you accept archive uploads?** → Use archive extraction
5. **Do users need GFM features?** → Update documentation

## Expected Outcome

After this upgrade:
- ✅ CMS uses v0.4.0 of asciidoc-xml
- ✅ All existing functionality preserved
- ✅ New Markdown features automatically available
- ✅ Optional: Logging integrated
- ✅ Optional: Batch processing enabled
- ✅ Optional: Archive support added

## Notes

- **No breaking changes**: All existing code continues to work
- **No new dependencies**: All features use Go standard library
- **Backward compatible**: Safe to upgrade without code changes
- **Optional features**: Logging, batch processing, and archives are optional enhancements

Proceed with the upgrade and integrate features as appropriate for your CMS architecture.

