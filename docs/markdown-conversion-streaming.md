# Markdown to AsciiDoc Conversion - Streaming and Parallelization

## Current Implementation Analysis

### Memory Usage

The `ConvertMarkdownToAsciiDoc` function uses **line-by-line streaming** for input but accumulates the entire output in memory:

1. **Input**: Uses `bufio.Scanner` - streams line-by-line ✅
2. **Output**: Accumulates in `bytes.Buffer`, returns as `string` ❌
3. **Temporary buffers**: 
   - `frontmatterLines []string` - accumulates frontmatter
   - `codeBlockLines []string` - accumulates code block content

### Memory Footprint

- **Small files (< 1MB)**: Minimal impact
- **Large files (> 100MB)**: Entire output held in memory
- **Very large files (> 1GB)**: May cause memory issues

## Streaming Implementation

### New Function: `ConvertMarkdownToAsciiDocStreaming`

A new streaming version has been added that writes directly to an `io.Writer`:

```go
func ConvertMarkdownToAsciiDocStreaming(reader io.Reader, writer io.Writer) error
```

**Benefits:**
- ✅ Writes output as it processes (streaming output)
- ✅ Lower memory footprint for large files
- ✅ Can write directly to files, network connections, etc.
- ✅ Same functionality as original

**Usage Example:**
```go
// Stream directly to a file
file, _ := os.Create("output.adoc")
defer file.Close()
err := lib.ConvertMarkdownToAsciiDocStreaming(inputReader, file)

// Stream to HTTP response
err := lib.ConvertMarkdownToAsciiDocStreaming(inputReader, httpResponseWriter)

// Still get string result if needed
var buf bytes.Buffer
err := lib.ConvertMarkdownToAsciiDocStreaming(inputReader, &buf)
result := buf.String()
```

### Backward Compatibility

The original `ConvertMarkdownToAsciiDoc` function now uses the streaming version internally:

```go
func ConvertMarkdownToAsciiDoc(reader io.Reader) (string, error) {
    var result bytes.Buffer
    if err := ConvertMarkdownToAsciiDocStreaming(reader, &result); err != nil {
        return "", err
    }
    return result.String(), nil
}
```

This maintains backward compatibility while using the improved implementation.

## Parallelization Considerations

### Why Chunking is Difficult

Markdown parsing is **inherently sequential** due to multi-line constructs:

1. **Frontmatter** - Must be processed first (at file start)
2. **Code blocks** - Span multiple lines (` ```...``` `)
3. **Tables** - Span multiple lines (rows)
4. **Blockquotes** - Can span multiple lines (`> ...`)
5. **Lists** - Can span multiple lines with nested content
6. **Context-dependent** - Line meaning depends on previous lines

### Parallelization Strategies

#### 1. **Process Multiple Files in Parallel** ✅ Recommended

This is the most practical approach:

```go
// Process multiple files concurrently
files := []string{"file1.md", "file2.md", "file3.md"}
var wg sync.WaitGroup
results := make(chan result, len(files))

for _, file := range files {
    wg.Add(1)
    go func(f string) {
        defer wg.Done()
        // Process file
        input, _ := os.Open(f)
        output, _ := os.Create(f + ".adoc")
        err := ConvertMarkdownToAsciiDocStreaming(input, output)
        results <- result{file: f, err: err}
    }(file)
}

wg.Wait()
close(results)
```

**Benefits:**
- ✅ No changes to parser logic needed
- ✅ Scales with number of files
- ✅ Each file processes independently
- ✅ Easy to implement

#### 2. **Chunking with Context Preservation** ⚠️ Complex

Theoretically possible but requires:

- **Context tracking** - Maintain state across chunks
- **Boundary detection** - Identify safe split points
- **State synchronization** - Share context between workers
- **Reassembly** - Merge results maintaining order

**Challenges:**
- Code blocks can be thousands of lines
- Tables can span many rows
- Nested structures complicate boundaries
- Overhead may exceed benefits

**Example Approach:**
```go
// Hypothetical chunked processing
func ProcessChunked(reader io.Reader, writer io.Writer, chunkSize int) error {
    // 1. Read first chunk (includes frontmatter)
    // 2. Process frontmatter
    // 3. Identify safe split points (empty lines, end of blocks)
    // 4. Process chunks in parallel with context
    // 5. Merge results maintaining order
}
```

**Not Recommended** - Complexity outweighs benefits for most use cases.

#### 3. **Pipeline Processing** ✅ Good for Large Files

Process different stages in parallel:

```go
// Stage 1: Read and parse frontmatter
// Stage 2: Process content blocks
// Stage 3: Convert inline formatting
// Stage 4: Write output

// Each stage can process different parts concurrently
```

**Benefits:**
- ✅ Can parallelize independent stages
- ✅ Maintains sequential order
- ✅ Good for very large files

**Limitations:**
- Requires significant refactoring
- May not provide much benefit for typical file sizes

## Recommendations

### For Most Use Cases

**Use `ConvertMarkdownToAsciiDocStreaming`** for:
- Large files (> 10MB)
- Network streaming
- File-to-file conversion
- Memory-constrained environments

**Use `ConvertMarkdownToAsciiDoc`** for:
- Small files (< 1MB)
- When string result is needed
- Simpler API requirements

### For Batch Processing

**Process multiple files in parallel:**

```go
func ProcessBatch(files []string, maxWorkers int) error {
    sem := make(chan struct{}, maxWorkers)
    var wg sync.WaitGroup
    
    for _, file := range files {
        wg.Add(1)
        go func(f string) {
            defer wg.Done()
            sem <- struct{}{} // Acquire
            defer func() { <-sem }() // Release
            
            // Process file
            input, _ := os.Open(f)
            defer input.Close()
            output, _ := os.Create(f + ".adoc")
            defer output.Close()
            ConvertMarkdownToAsciiDocStreaming(input, output)
        }(file)
    }
    
    wg.Wait()
    return nil
}
```

### For Very Large Single Files

1. **Use streaming** - `ConvertMarkdownToAsciiDocStreaming`
2. **Increase scanner buffer** - For very long lines:
   ```go
   scanner := bufio.NewScanner(reader)
   buf := make([]byte, 0, 64*1024) // 64KB buffer
   scanner.Buffer(buf, 1024*1024)  // Max 1MB line
   ```
3. **Consider chunked reading** - Only if file is extremely large (> 1GB)

## Performance Characteristics

### Current Implementation

- **Input**: O(n) time, O(1) memory per line
- **Output**: O(n) time, O(n) memory (entire output)
- **Total**: O(n) time, O(n) memory

### Streaming Implementation

- **Input**: O(n) time, O(1) memory per line
- **Output**: O(n) time, O(1) memory (streaming)
- **Total**: O(n) time, O(1) memory (excluding buffers)

### Parallel File Processing

- **Time**: O(n/m) where m = number of workers
- **Memory**: O(1) per file (with streaming)
- **Scalability**: Linear with number of files

## Conclusion

1. **Streaming**: ✅ Implemented - Use `ConvertMarkdownToAsciiDocStreaming` for large files
2. **Chunking single files**: ❌ Not practical - Markdown is inherently sequential
3. **Parallel file processing**: ✅ Recommended - Process multiple files concurrently
4. **Memory efficiency**: ✅ Improved - Streaming reduces memory footprint

The current implementation is well-suited for most use cases. For batch processing of many files, parallelize at the file level rather than chunking individual files.

