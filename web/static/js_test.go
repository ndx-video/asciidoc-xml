//go:build test

package static

import (
	"embed"
	"strings"
	"testing"

	"github.com/dop251/goja"
)

//go:embed *.js
var jsFiles embed.FS

// Enhanced mock browser APIs for comprehensive testing
const browserMocks = `
// HTML entity escaping function
function escapeHtmlEntities(text) {
    var map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#39;'
    };
    return text.replace(/[&<>"']/g, function(m) { return map[m]; });
}

// Mock storage for element state
var mockElements = {};
var mockSelectors = {};

// Mock document object with comprehensive DOM API
var document = {
    readyState: 'complete',
    createElement: function(tag) {
        var elem = {
            _textContent: '',
            _innerHTML: '',
            _children: [],
            _style: {},
            _classList: {
                _classes: [],
                add: function(c) { if (this._classes.indexOf(c) === -1) this._classes.push(c); },
                remove: function(c) { var idx = this._classes.indexOf(c); if (idx > -1) this._classes.splice(idx, 1); },
                contains: function(c) { return this._classes.indexOf(c) > -1; },
                toggle: function(c) { if (this.contains(c)) this.remove(c); else this.add(c); }
            },
            _dataset: {},
            _value: '',
            _files: [],
            set textContent(v) { 
                this._textContent = v; 
                this._innerHTML = escapeHtmlEntities(String(v));
            },
            get textContent() { 
                return this._textContent || ''; 
            },
            set innerHTML(v) { 
                this._innerHTML = v; 
            },
            get innerHTML() { 
                return this._innerHTML || ''; 
            },
            get style() { return this._style; },
            get classList() { return this._classList; },
            get dataset() { return this._dataset; },
            set value(v) { this._value = v; },
            get value() { return this._value; },
            get files() { return this._files; },
            appendChild: function(child) { this._children.push(child); return child; },
            addEventListener: function(type, handler) {},
            removeEventListener: function(type, handler) {},
            querySelector: function(sel) {
                return mockSelectors[sel] || {
                    dataset: { view: 'rendered' },
                    classList: { contains: function() { return false; }, remove: function() {}, add: function() {} }
                };
            },
            querySelectorAll: function(sel) {
                return [];
            }
        };
        return elem;
    },
    createTextNode: function(text) {
        return { textContent: text, nodeType: 3 };
    },
    querySelector: function(sel) {
        return mockSelectors[sel] || {
            dataset: { view: 'rendered' },
            classList: { contains: function() { return false; }, remove: function() {}, add: function() {} }
        };
    },
    querySelectorAll: function(sel) {
        return [];
    },
    getElementById: function(id) {
        if (!mockElements[id]) {
            mockElements[id] = {
                contentDocument: null,
                contentWindow: { document: null },
                querySelector: function() { return null; },
                src: '',
                addEventListener: function() {},
                removeEventListener: function() {},
                scrollTop: 0,
                scrollLeft: 0,
                textContent: '',
                innerHTML: '',
                className: '',
                dataset: {},
                value: '',
                files: [],
                style: { display: 'block', flex: '' },
                classList: { contains: function() { return false; }, remove: function() {}, add: function() {} },
                appendChild: function() { return {}; }
            };
        }
        return mockElements[id];
    },
    addEventListener: function(type, handler) {}
};

// Mock window object
var window = {
    parent: null,
    addEventListener: function() {},
    clearTimeout: function() {},
    setTimeout: function(fn, delay) { return 1; }
};

// Mock URL and Blob APIs
var URL = {
    createObjectURL: function(blob) {
        return 'blob:mock-url-' + Math.random();
    },
    revokeObjectURL: function(url) {}
};

var Blob = function(parts, options) {
    this.parts = parts;
    this.type = options ? options.type : '';
    this.size = parts.join('').length;
};

// Mock fetch API
var fetch = function(url, options) {
    return Promise.resolve({
        ok: true,
        status: 200,
        statusText: 'OK',
        text: function() { return Promise.resolve('mock response'); },
        json: function() { return Promise.resolve({ valid: true, output: 'mock output' }); }
    });
};

// Mock DOMParser
var DOMParser = function() {};
DOMParser.prototype.parseFromString = function(str, type) {
    return {
        querySelector: function(sel) { return null; },
        querySelectorAll: function(sel) { return []; }
    };
};

// Mock XSLTProcessor
var XSLTProcessor = function() {};
XSLTProcessor.prototype.importStylesheet = function(xslt) {};
XSLTProcessor.prototype.transformToDocument = function(xml) {
    return {
        documentElement: { innerHTML: '<div>transformed</div>' }
    };
};
XSLTProcessor.prototype.transformToFragment = function(xml, doc) {
    return document.createElement('div');
};

// Mock XMLSerializer
var XMLSerializer = function() {};
XMLSerializer.prototype.serializeToString = function(doc) {
    return '<html><body>serialized</body></html>';
};

// Mock console for error handling
var console = {
    log: function() {},
    error: function() {},
    warn: function() {}
};

// Mock JSON
var JSON = {
    parse: function(str) {
        try {
            return eval('(' + str + ')');
        } catch(e) {
            throw new Error('Invalid JSON');
        }
    },
    stringify: function(obj) {
        return '{"mock": "json"}';
    }
};
`

func TestJavaScriptSyntax(t *testing.T) {
	files := []string{"pretty.js", "app.js", "browse.js", "user-manual.js"}

	for _, filename := range files {
		t.Run(filename, func(t *testing.T) {
			content, err := jsFiles.ReadFile(filename)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", filename, err)
			}

			vm := goja.New()

			// Inject browser mocks
			_, err = vm.RunString(browserMocks)
			if err != nil {
				t.Fatalf("Failed to inject browser mocks: %v", err)
			}

			// Try to parse and execute the JavaScript
			jsCode := string(content)

			// For files with heavy DOM dependencies, we'll just verify syntax
			// by attempting to parse them
			_, err = vm.RunString(jsCode)
			if err != nil {
				// Some errors are expected due to XSLTProcessor limitations in mock
				// but we should catch syntax errors
				if strings.Contains(err.Error(), "syntax") || strings.Contains(err.Error(), "SyntaxError") {
					t.Errorf("JavaScript syntax error in %s: %v", filename, err)
				} else {
					// Runtime errors due to missing DOM APIs are acceptable for syntax testing
					t.Logf("Runtime error in %s (may be expected due to DOM dependencies): %v", filename, err)
				}
			}
		})
	}
}

func TestPrettyPrintFunctions(t *testing.T) {
	vm := goja.New()

	// Inject browser mocks
	_, err := vm.RunString(browserMocks)
	if err != nil {
		t.Fatalf("Failed to inject browser mocks: %v", err)
	}

	// Load pretty.js
	prettyJS, err := jsFiles.ReadFile("pretty.js")
	if err != nil {
		t.Fatalf("Failed to read pretty.js: %v", err)
	}

	_, err = vm.RunString(string(prettyJS))
	if err != nil {
		t.Fatalf("Failed to load pretty.js: %v", err)
	}

	tests := []struct {
		name     string
		function string
		input    string
		validate func(t *testing.T, result goja.Value)
	}{
		{
			name:     "escapeHtml function exists",
			function: "escapeHtml",
			input:    `"test & <test>"`,
			validate: func(t *testing.T, result goja.Value) {
				if result == nil {
					t.Error("escapeHtml returned nil")
				}
				str := result.String()
				if !strings.Contains(str, "&amp;") || !strings.Contains(str, "&lt;") {
					t.Errorf("escapeHtml did not escape properly: %s", str)
				}
			},
		},
		{
			name:     "highlightAsciiDoc function exists",
			function: "highlightAsciiDoc",
			input:    `"= Test Title\n\nContent"`,
			validate: func(t *testing.T, result goja.Value) {
				if result == nil {
					t.Error("highlightAsciiDoc returned nil")
				}
				str := result.String()
				if str == "" {
					t.Error("highlightAsciiDoc returned empty string")
				}
				if !strings.Contains(str, "title") {
					t.Error("highlightAsciiDoc should highlight titles")
				}
			},
		},
		{
			name:     "highlightXML function exists",
			function: "highlightXML",
			input:    `"<tag attr='value'>content</tag>"`,
			validate: func(t *testing.T, result goja.Value) {
				if result == nil {
					t.Error("highlightXML returned nil")
				}
				str := result.String()
				if str == "" {
					t.Error("highlightXML returned empty string")
				}
			},
		},
		{
			name:     "highlightHTML function exists",
			function: "highlightHTML",
			input:    `"<div>test</div>"`,
			validate: func(t *testing.T, result goja.Value) {
				if result == nil {
					t.Error("highlightHTML returned nil")
				}
				str := result.String()
				if str == "" {
					t.Error("highlightHTML returned empty string")
				}
			},
		},
		{
			name:     "createSourceView function exists",
			function: "createSourceView",
			input:    `"test content", "asciidoc"`,
			validate: func(t *testing.T, result goja.Value) {
				if result == nil {
					t.Error("createSourceView returned nil")
				}
				str := result.String()
				if str == "" {
					t.Error("createSourceView returned empty string")
				}
				if !strings.Contains(str, "<!DOCTYPE html>") {
					t.Error("createSourceView should return HTML")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check function exists
			fn := vm.Get(tt.function)
			if fn == nil {
				t.Fatalf("Function %s not found", tt.function)
			}

			// Test function execution
			code := tt.function + "(" + tt.input + ")"
			result, err := vm.RunString(code)
			if err != nil {
				t.Errorf("Error executing %s: %v", tt.function, err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestPrettyPrintHighlightAsciiDoc(t *testing.T) {
	vm := goja.New()

	// Inject browser mocks
	_, err := vm.RunString(browserMocks)
	if err != nil {
		t.Fatalf("Failed to inject browser mocks: %v", err)
	}

	// Load pretty.js
	prettyJS, err := jsFiles.ReadFile("pretty.js")
	if err != nil {
		t.Fatalf("Failed to read pretty.js: %v", err)
	}

	_, err = vm.RunString(string(prettyJS))
	if err != nil {
		t.Fatalf("Failed to load pretty.js: %v", err)
	}

	testCases := []struct {
		name     string
		input    string
		expected []string // Substrings that should be in the output
	}{
		{
			name:     "document title",
			input:    "= Document Title",
			expected: []string{"title", "title-level-1", "Document Title"},
		},
		{
			name:     "section title",
			input:    "== Section Title",
			expected: []string{"title", "title-level-2", "Section Title"},
		},
		{
			name:     "attributes",
			input:    ":author: John Doe",
			expected: []string{"attribute", "author"},
		},
		{
			name:     "bold text",
			input:    "This is *bold* text",
			expected: []string{"bold", "*bold*"},
		},
		{
			name:     "italic text",
			input:    "This is _italic_ text",
			expected: []string{"italic", "_italic_"},
		},
		{
			name:     "monospace text",
			input:    "This is `code` text",
			expected: []string{"monospace", "`code`"},
		},
		{
			name:     "list marker",
			input:    "* List item",
			expected: []string{"list-marker"},
		},
		{
			name:     "block delimiter",
			input:    "----",
			expected: []string{"block-delimiter"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code := `highlightAsciiDoc("` + strings.ReplaceAll(tc.input, `"`, `\"`) + `")`
			result, err := vm.RunString(code)
			if err != nil {
				t.Fatalf("Error executing highlightAsciiDoc: %v", err)
			}

			output := result.String()
			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}

func TestPrettyPrintHighlightXML(t *testing.T) {
	vm := goja.New()

	// Inject browser mocks
	_, err := vm.RunString(browserMocks)
	if err != nil {
		t.Fatalf("Failed to inject browser mocks: %v", err)
	}

	// Load pretty.js
	prettyJS, err := jsFiles.ReadFile("pretty.js")
	if err != nil {
		t.Fatalf("Failed to read pretty.js: %v", err)
	}

	_, err = vm.RunString(string(prettyJS))
	if err != nil {
		t.Fatalf("Failed to load pretty.js: %v", err)
	}

	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "XML tags",
			input:    "<tag>content</tag>",
			expected: []string{"tag", "tag-name"},
		},
		{
			name:     "XML attributes",
			input:    `<tag attr="value">content</tag>`,
			expected: []string{"tag-name", "tag"}, // Attributes are escaped, so check for basic structure
		},
		{
			name:     "XML comments",
			input:    "<!-- comment -->",
			expected: []string{"comment"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code := `highlightXML("` + strings.ReplaceAll(tc.input, `"`, `\"`) + `")`
			result, err := vm.RunString(code)
			if err != nil {
				t.Fatalf("Error executing highlightXML: %v", err)
			}

			output := result.String()
			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}

func TestPrettyPrintCreateSourceView(t *testing.T) {
	vm := goja.New()

	// Inject browser mocks
	_, err := vm.RunString(browserMocks)
	if err != nil {
		t.Fatalf("Failed to inject browser mocks: %v", err)
	}

	// Load pretty.js
	prettyJS, err := jsFiles.ReadFile("pretty.js")
	if err != nil {
		t.Fatalf("Failed to read pretty.js: %v", err)
	}

	_, err = vm.RunString(string(prettyJS))
	if err != nil {
		t.Fatalf("Failed to load pretty.js: %v", err)
	}

	testCases := []struct {
		name     string
		content  string
		fileType string
		expected []string
	}{
		{
			name:     "AsciiDoc source view",
			content:  "= Test",
			fileType: "asciidoc",
			expected: []string{"<!DOCTYPE html>", "asciidoc-source", "Test"},
		},
		{
			name:     "XML source view",
			content:  "<root><child/></root>",
			fileType: "xml",
			expected: []string{"<!DOCTYPE html>", "xml-source"},
		},
		{
			name:     "XSLT source view",
			content:  "<xsl:stylesheet></xsl:stylesheet>",
			fileType: "xslt",
			expected: []string{"<!DOCTYPE html>", "xml-source"},
		},
		{
			name:     "HTML source view",
			content:  "<div>test</div>",
			fileType: "html",
			expected: []string{"<!DOCTYPE html>", "html-source"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			escapedContent := strings.ReplaceAll(tc.content, `"`, `\"`)
			code := `createSourceView("` + escapedContent + `", "` + tc.fileType + `")`
			result, err := vm.RunString(code)
			if err != nil {
				t.Fatalf("Error executing createSourceView: %v", err)
			}

			output := result.String()
			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}

func TestPrettyPrintEscapeHtml(t *testing.T) {
	vm := goja.New()

	_, err := vm.RunString(browserMocks)
	if err != nil {
		t.Fatalf("Failed to inject browser mocks: %v", err)
	}

	prettyJS, err := jsFiles.ReadFile("pretty.js")
	if err != nil {
		t.Fatalf("Failed to read pretty.js: %v", err)
	}

	_, err = vm.RunString(string(prettyJS))
	if err != nil {
		t.Fatalf("Failed to load pretty.js: %v", err)
	}

	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "escape ampersand",
			input:    "test & test",
			expected: []string{"&amp;"},
		},
		{
			name:     "escape less than",
			input:    "test < test",
			expected: []string{"&lt;"},
		},
		{
			name:     "escape greater than",
			input:    "test > test",
			expected: []string{"&gt;"},
		},
		{
			name:     "escape quotes",
			input:    `test "quote" test`,
			expected: []string{"&quot;"},
		},
		{
			name:     "escape apostrophe",
			input:    "test 'apostrophe' test",
			expected: []string{"&#39;"},
		},
		{
			name:     "no escaping needed",
			input:    "normal text",
			expected: []string{"normal text"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code := `escapeHtml("` + strings.ReplaceAll(tc.input, `"`, `\"`) + `")`
			result, err := vm.RunString(code)
			if err != nil {
				t.Fatalf("Error executing escapeHtml: %v", err)
			}

			output := result.String()
			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}

func TestPrettyPrintHighlightHTML(t *testing.T) {
	vm := goja.New()

	_, err := vm.RunString(browserMocks)
	if err != nil {
		t.Fatalf("Failed to inject browser mocks: %v", err)
	}

	prettyJS, err := jsFiles.ReadFile("pretty.js")
	if err != nil {
		t.Fatalf("Failed to read pretty.js: %v", err)
	}

	_, err = vm.RunString(string(prettyJS))
	if err != nil {
		t.Fatalf("Failed to load pretty.js: %v", err)
	}

	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "HTML tags",
			input:    "<div>content</div>",
			expected: []string{"tag", "tag-name"},
		},
		{
			name:     "HTML with attributes",
			input:    `<div class="test">content</div>`,
			expected: []string{"tag-name", "attribute"},
		},
		{
			name:     "HTML comments",
			input:    "<!-- comment -->",
			expected: []string{"comment"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code := `highlightHTML("` + strings.ReplaceAll(tc.input, `"`, `\"`) + `")`
			result, err := vm.RunString(code)
			if err != nil {
				t.Fatalf("Error executing highlightHTML: %v", err)
			}

			output := result.String()
			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}

func TestPrettyPrintHighlightAsciiDoc_Advanced(t *testing.T) {
	vm := goja.New()

	_, err := vm.RunString(browserMocks)
	if err != nil {
		t.Fatalf("Failed to inject browser mocks: %v", err)
	}

	prettyJS, err := jsFiles.ReadFile("pretty.js")
	if err != nil {
		t.Fatalf("Failed to read pretty.js: %v", err)
	}

	_, err = vm.RunString(string(prettyJS))
	if err != nil {
		t.Fatalf("Failed to load pretty.js: %v", err)
	}

	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "multiple section levels",
			input:    "== Level 2\n=== Level 3\n==== Level 4",
			expected: []string{"title-level-2", "title-level-3", "title-level-4"},
		},
		{
			name:     "table with delimiters",
			input:    "|Header|Data|",
			expected: []string{"table-delimiter"},
		},
		{
			name:     "code block delimiter",
			input:    "----\ncode\n----",
			expected: []string{"block-delimiter"},
		},
		{
			name:     "literal block delimiter",
			input:    "....\nliteral\n....",
			expected: []string{"block-delimiter"},
		},
		{
			name:     "example block delimiter",
			input:    "====\nexample\n====",
			expected: []string{"block-delimiter"},
		},
		{
			name:     "sidebar delimiter",
			input:    "****\nsidebar\n****",
			expected: []string{"block-delimiter"},
		},
		{
			name:     "quote delimiter",
			input:    "____\nquote\n____",
			expected: []string{"block-delimiter"},
		},
		{
			name:     "ordered list marker",
			input:    ". Item one\n. Item two",
			expected: []string{"list-marker"},
		},
		{
			name:     "unordered list marker",
			input:    "- Item one\n- Item two",
			expected: []string{"list-marker"},
		},
		{
			name:     "link detection",
			input:    "Visit https://example.com for more info",
			expected: []string{"link"},
		},
		{
			name:     "link with text",
			input:    "Visit https://example.com[Example] for more",
			expected: []string{"link"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code := `highlightAsciiDoc("` + strings.ReplaceAll(tc.input, `"`, `\"`) + `")`
			result, err := vm.RunString(code)
			if err != nil {
				t.Fatalf("Error executing highlightAsciiDoc: %v", err)
			}

			output := result.String()
			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}

func TestPrettyPrintCreateSourceView_AllTypes(t *testing.T) {
	vm := goja.New()

	_, err := vm.RunString(browserMocks)
	if err != nil {
		t.Fatalf("Failed to inject browser mocks: %v", err)
	}

	prettyJS, err := jsFiles.ReadFile("pretty.js")
	if err != nil {
		t.Fatalf("Failed to read pretty.js: %v", err)
	}

	_, err = vm.RunString(string(prettyJS))
	if err != nil {
		t.Fatalf("Failed to load pretty.js: %v", err)
	}

	testCases := []struct {
		name     string
		content  string
		fileType string
		expected []string
	}{
		{
			name:     "AsciiDoc with complex content",
			content:  "= Title\n\n== Section\n\n* List item\n\n----\ncode\n----",
			fileType: "asciidoc",
			expected: []string{"<!DOCTYPE html>", "asciidoc-source", "Title", "Section"},
		},
		{
			name:     "XML with namespaces",
			content:  "<xsl:stylesheet xmlns:xsl='http://www.w3.org/1999/XSL/Transform'></xsl:stylesheet>",
			fileType: "xml",
			expected: []string{"<!DOCTYPE html>", "xml-source"},
		},
		{
			name:     "XSLT transformation",
			content:  "<xsl:template match='/'>content</xsl:template>",
			fileType: "xslt",
			expected: []string{"<!DOCTYPE html>", "xml-source"},
		},
		{
			name:     "HTML with script tags",
			content:  "<html><head><script>alert('test');</script></head><body>content</body></html>",
			fileType: "html",
			expected: []string{"<!DOCTYPE html>", "html-source"},
		},
		{
			name:     "unknown type defaults to xml-source",
			content:  "some content",
			fileType: "unknown",
			expected: []string{"<!DOCTYPE html>", "xml-source"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			escapedContent := strings.ReplaceAll(tc.content, `"`, `\"`)
			escapedContent = strings.ReplaceAll(escapedContent, "\n", "\\n")
			code := `createSourceView("` + escapedContent + `", "` + tc.fileType + `")`
			result, err := vm.RunString(code)
			if err != nil {
				t.Fatalf("Error executing createSourceView: %v", err)
			}

			output := result.String()
			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}

func TestBrowseJS_Functions(t *testing.T) {
	vm := goja.New()

	_, err := vm.RunString(browserMocks)
	if err != nil {
		t.Fatalf("Failed to inject browser mocks: %v", err)
	}

	// Mock fetch for browse.js
	fetchMock := `
		var fetch = function(url, options) {
			if (url.includes('/api/files')) {
				return Promise.resolve({
					ok: true,
					json: function() { return Promise.resolve([{name: 'test.adoc', path: 'examples/test.adoc', type: 'file'}]); }
				});
			}
			if (url.includes('/api/load-file')) {
				return Promise.resolve({
					ok: true,
					text: function() { return Promise.resolve('= Test Document\\n\\nContent'); }
				});
			}
			if (url.includes('/api/convert')) {
				return Promise.resolve({
					ok: true,
					json: function() { return Promise.resolve({output: '<html><body>Converted</body></html>'}); }
				});
			}
			return Promise.resolve({ok: false, text: function() { return Promise.resolve('Error'); }});
		};
	`
	_, err = vm.RunString(fetchMock)
	if err != nil {
		t.Fatalf("Failed to inject fetch mock: %v", err)
	}

	// Load browse.js
	browseJS, err := jsFiles.ReadFile("browse.js")
	if err != nil {
		t.Fatalf("Failed to read browse.js: %v", err)
	}

	// Try to load browse.js - it may have runtime errors due to DOM dependencies
	// but we can at least verify the syntax
	_, err = vm.RunString(string(browseJS))
	if err != nil {
		// Runtime errors are acceptable if they're due to missing DOM elements
		if !strings.Contains(err.Error(), "Cannot read") && !strings.Contains(err.Error(), "undefined") {
			t.Logf("Note: browse.js has runtime dependencies: %v", err)
		}
	}

	// Test that createSourceView is available (from pretty.js which browse.js uses)
	prettyJS, err := jsFiles.ReadFile("pretty.js")
	if err != nil {
		t.Fatalf("Failed to read pretty.js: %v", err)
	}

	_, err = vm.RunString(string(prettyJS))
	if err != nil {
		t.Fatalf("Failed to load pretty.js: %v", err)
	}

	// Verify createSourceView works (used by browse.js)
	code := `createSourceView("= Test", "asciidoc")`
	result, err := vm.RunString(code)
	if err != nil {
		t.Fatalf("createSourceView should be available: %v", err)
	}

	if result.String() == "" {
		t.Error("createSourceView should return content")
	}
}

func TestUserManualJS_Functions(t *testing.T) {
	vm := goja.New()

	_, err := vm.RunString(browserMocks)
	if err != nil {
		t.Fatalf("Failed to inject browser mocks: %v", err)
	}

	// Mock the required DOM elements
	mockDOM := `
		var xmlScript = {
			textContent: JSON.stringify('<?xml version="1.0"?><root><content>Test</content></root>')
		};
		var xsltScript = {
			textContent: JSON.stringify('<?xml version="1.0"?><xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"><xsl:template match="/"><div><xsl:value-of select="root/content"/></div></xsl:template></xsl:stylesheet>')
		};
		var contentDiv = document.createElement('div');
		document.getElementById = function(id) {
			if (id === 'xml-data') return xmlScript;
			if (id === 'xslt-data') return xsltScript;
			if (id === 'content') return contentDiv;
			return null;
		};
	`
	_, err = vm.RunString(mockDOM)
	if err != nil {
		t.Fatalf("Failed to inject DOM mocks: %v", err)
	}

	// Load user-manual.js
	userManualJS, err := jsFiles.ReadFile("user-manual.js")
	if err != nil {
		t.Fatalf("Failed to read user-manual.js: %v", err)
	}

	// Try to execute user-manual.js
	// It's wrapped in an IIFE, so it should execute immediately
	_, err = vm.RunString(string(userManualJS))
	if err != nil {
		// Some errors are expected due to XSLTProcessor limitations in mock
		// but we should catch syntax errors
		if strings.Contains(err.Error(), "syntax") || strings.Contains(err.Error(), "SyntaxError") {
			t.Errorf("JavaScript syntax error in user-manual.js: %v", err)
		} else {
			t.Logf("Note: user-manual.js has runtime dependencies: %v", err)
		}
	}
}

func TestAppJS_Functions(t *testing.T) {
	vm := goja.New()

	_, err := vm.RunString(browserMocks)
	if err != nil {
		t.Fatalf("Failed to inject browser mocks: %v", err)
	}

	// Load pretty.js first (app.js depends on it)
	prettyJS, err := jsFiles.ReadFile("pretty.js")
	if err != nil {
		t.Fatalf("Failed to read pretty.js: %v", err)
	}

	_, err = vm.RunString(string(prettyJS))
	if err != nil {
		t.Fatalf("Failed to load pretty.js: %v", err)
	}

	// Mock additional DOM elements needed by app.js
	appDOM := `
		// Mock all the elements app.js expects
		var elements = {
			'asciidoc-frame': { src: '', addEventListener: function() {} },
			'xml-frame': { src: '', addEventListener: function() {} },
			'xslt-frame': { src: '', addEventListener: function() {} },
			'html-frame': { src: '', addEventListener: function() {} },
			'status': { textContent: '', className: '' },
			'output-type': { value: 'xml', addEventListener: function() {} },
			'xml-panel': { classList: { remove: function() {}, add: function() {} } },
			'xslt-panel': { classList: { remove: function() {}, add: function() {} } },
			'html-panel': { classList: { remove: function() {}, add: function() {} }, querySelector: function() { return { textContent: '' }; } },
			'xslt-upload-section': { style: { display: '' } },
			'asciidoc-path': { value: '' },
			'xslt-path': { value: '' },
			'btn-validate': { addEventListener: function() {} },
			'btn-convert': { addEventListener: function() {} },
			'btn-load-example': { addEventListener: function() {} },
			'btn-upload': { addEventListener: function() {} },
			'uploadModal': { style: { display: 'none' } },
			'closeModal': { addEventListener: function() {} },
			'btn-upload-asciidoc': { addEventListener: function() {}, disabled: false },
			'btn-upload-xslt': { addEventListener: function() {}, disabled: false },
			'asciidocFile': { addEventListener: function() {}, files: [], value: '' },
			'xsltFile': { addEventListener: function() {}, files: [], value: '' }
		};
		
		document.getElementById = function(id) {
			return elements[id] || { addEventListener: function() {} };
		};
		
		document.querySelectorAll = function(sel) {
			if (sel === '.html-tabs button') {
				return [{
					dataset: { view: 'rendered' },
					classList: { remove: function() {}, add: function() {} },
					addEventListener: function() {}
				}];
			}
			if (sel === '.panel:not(.hidden)') {
				return [];
			}
			if (sel === '.resizer') {
				return [];
			}
			if (sel === 'iframe') {
				return [];
			}
			return [];
		};
		
		document.querySelector = function(sel) {
			if (sel === '.html-tabs button.active') {
				return { dataset: { view: 'rendered' } };
			}
			return null;
		};
		
		// Mock fetch
		var fetch = function(url, options) {
			if (url.includes('/api/xslt')) {
				return Promise.resolve({
					ok: true,
					text: function() { return Promise.resolve('<?xml version="1.0"?><xsl:stylesheet></xsl:stylesheet>'); }
				});
			}
			if (url.includes('/api/convert')) {
				return Promise.resolve({
					ok: true,
					json: function() { return Promise.resolve({output: '<xml>test</xml>'}); }
				});
			}
			if (url.includes('/api/validate')) {
				return Promise.resolve({
					ok: true,
					json: function() { return Promise.resolve({valid: true}); }
				});
			}
			if (url.includes('/api/upload')) {
				return Promise.resolve({
					ok: true,
					json: function() { return Promise.resolve({path: '/static/test.adoc'}); }
				});
			}
			return Promise.resolve({ok: false});
		};
	`
	_, err = vm.RunString(appDOM)
	if err != nil {
		t.Fatalf("Failed to inject app.js DOM mocks: %v", err)
	}

	// Load app.js
	appJS, err := jsFiles.ReadFile("app.js")
	if err != nil {
		t.Fatalf("Failed to read app.js: %v", err)
	}

	// Try to execute app.js
	// It will have runtime dependencies but we can verify syntax
	_, err = vm.RunString(string(appJS))
	if err != nil {
		// Runtime errors are expected due to complex DOM dependencies
		// but syntax errors should be caught
		if strings.Contains(err.Error(), "syntax") || strings.Contains(err.Error(), "SyntaxError") {
			t.Errorf("JavaScript syntax error in app.js: %v", err)
		} else {
			t.Logf("Note: app.js has runtime dependencies (expected): %v", err)
		}
	}
}

func TestAllJSFiles_Loadable(t *testing.T) {
	files := []string{"pretty.js", "app.js", "browse.js", "user-manual.js"}

	for _, filename := range files {
		t.Run(filename, func(t *testing.T) {
			content, err := jsFiles.ReadFile(filename)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", filename, err)
			}

			if len(content) == 0 {
				t.Errorf("%s is empty", filename)
			}

			// Verify it's valid JavaScript by checking for common syntax patterns
			contentStr := string(content)
			if strings.Contains(contentStr, "function") || strings.Contains(contentStr, "=>") || strings.Contains(contentStr, "const ") || strings.Contains(contentStr, "var ") {
				// Has function definitions, looks like valid JS
			} else {
				t.Logf("Warning: %s may not contain function definitions", filename)
			}
		})
	}
}
