package lib

import (
	"testing"
)

func TestSubstituteAttributes(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		attrs    map[string]string
		expected string
	}{
		{
			name:     "empty text",
			text:     "",
			attrs:    map[string]string{"attr1": "value1"},
			expected: "",
		},
		{
			name:     "no attributes",
			text:     "This is plain text",
			attrs:    nil,
			expected: "This is plain text",
		},
		{
			name:     "empty attributes map",
			text:     "This is plain text",
			attrs:    map[string]string{},
			expected: "This is plain text",
		},
		{
			name:     "single attribute replacement",
			text:     "Hello {name}",
			attrs:    map[string]string{"name": "World"},
			expected: "Hello World",
		},
		{
			name:     "multiple attribute replacements",
			text:     "{greeting} {name}, welcome to {place}",
			attrs:    map[string]string{"greeting": "Hello", "name": "Alice", "place": "Wonderland"},
			expected: "Hello Alice, welcome to Wonderland",
		},
		{
			name:     "attribute with hyphens",
			text:     "Version {version-number}",
			attrs:    map[string]string{"version-number": "1.0.0"},
			expected: "Version 1.0.0",
		},
		{
			name:     "attribute with underscores",
			text:     "User {user_name}",
			attrs:    map[string]string{"user_name": "john_doe"},
			expected: "User john_doe",
		},
		{
			name:     "attribute not found",
			text:     "Hello {missing}",
			attrs:    map[string]string{"other": "value"},
			expected: "Hello ",
		},
		{
			name:     "multiple same attribute",
			text:     "{attr} and {attr} again",
			attrs:    map[string]string{"attr": "value"},
			expected: "value and value again",
		},
		{
			name:     "attribute in middle of text",
			text:     "Start {middle} end",
			attrs:    map[string]string{"middle": "MIDDLE"},
			expected: "Start MIDDLE end",
		},
		{
			name:     "empty attribute value",
			text:     "Text {empty} here",
			attrs:    map[string]string{"empty": ""},
			expected: "Text  here",
		},
		{
			name:     "mixed text with and without attributes",
			text:     "Plain text {attr} more plain text",
			attrs:    map[string]string{"attr": "replaced"},
			expected: "Plain text replaced more plain text",
		},
		{
			name:     "attribute with special characters in value",
			text:     "Value: {special}",
			attrs:    map[string]string{"special": "a{b}c"},
			expected: "Value: a{b}c",
		},
		{
			name:     "invalid attribute syntax (not replaced)",
			text:     "{incomplete and {valid}",
			attrs:    map[string]string{"valid": "OK"},
			expected: "{incomplete and OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SubstituteAttributes(tt.text, tt.attrs)
			if result != tt.expected {
				t.Errorf("SubstituteAttributes(%q, %v) = %q, want %q", tt.text, tt.attrs, result, tt.expected)
			}
		})
	}
}

func TestGetBuiltInAttribute(t *testing.T) {
	tests := []struct {
		name     string
		attrName string
		docAttrs map[string]string
		expected string
	}{
		{
			name:     "title attribute found",
			attrName: "title",
			docAttrs: map[string]string{"title": "My Document"},
			expected: "My Document",
		},
		{
			name:     "title attribute not found",
			attrName: "title",
			docAttrs: map[string]string{"other": "value"},
			expected: "",
		},
		{
			name:     "author attribute found",
			attrName: "author",
			docAttrs: map[string]string{"author": "John Doe"},
			expected: "John Doe",
		},
		{
			name:     "revnumber attribute found",
			attrName: "revnumber",
			docAttrs: map[string]string{"revnumber": "1.0"},
			expected: "1.0",
		},
		{
			name:     "revdate attribute found",
			attrName: "revdate",
			docAttrs: map[string]string{"revdate": "2024-01-01"},
			expected: "2024-01-01",
		},
		{
			name:     "revremark attribute found",
			attrName: "revremark",
			docAttrs: map[string]string{"revremark": "Initial release"},
			expected: "Initial release",
		},
		{
			name:     "custom attribute with colon prefix",
			attrName: ":custom",
			docAttrs: map[string]string{":custom": "custom-value"},
			expected: "custom-value",
		},
		{
			name:     "custom attribute without colon prefix",
			attrName: "custom-attr",
			docAttrs: map[string]string{"custom-attr": "value"},
			expected: "value",
		},
		{
			name:     "custom attribute not found",
			attrName: "missing",
			docAttrs: map[string]string{"other": "value"},
			expected: "",
		},
		{
			name:     "empty docAttrs",
			attrName: "title",
			docAttrs: map[string]string{},
			expected: "",
		},
		{
			name:     "nil docAttrs",
			attrName: "title",
			docAttrs: nil,
			expected: "",
		},
		{
			name:     "colon prefix but key stored without colon",
			attrName: ":test",
			docAttrs: map[string]string{":test": "value"},
			expected: "value",
		},
		{
			name:     "colon prefix with key stored with colon",
			attrName: ":test",
			docAttrs: map[string]string{":test": "value"},
			expected: "value",
		},
		{
			name:     "unknown built-in attribute",
			attrName: "unknown",
			docAttrs: map[string]string{"unknown": "value"},
			expected: "value",
		},
		{
			name:     "unknown built-in attribute not found",
			attrName: "unknown",
			docAttrs: map[string]string{"other": "value"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetBuiltInAttribute(tt.attrName, tt.docAttrs)
			if result != tt.expected {
				t.Errorf("GetBuiltInAttribute(%q, %v) = %q, want %q", tt.attrName, tt.docAttrs, result, tt.expected)
			}
		})
	}
}

func TestMergeAttributes(t *testing.T) {
	tests := []struct {
		name     string
		base     map[string]string
		override map[string]string
		expected map[string]string
	}{
		{
			name:     "both nil",
			base:     nil,
			override: nil,
			expected: map[string]string{},
		},
		{
			name:     "base nil, override has values",
			base:     nil,
			override: map[string]string{"a": "1", "b": "2"},
			expected: map[string]string{"a": "1", "b": "2"},
		},
		{
			name:     "base has values, override nil",
			base:     map[string]string{"a": "1", "b": "2"},
			override: nil,
			expected: map[string]string{"a": "1", "b": "2"},
		},
		{
			name:     "base empty, override has values",
			base:     map[string]string{},
			override: map[string]string{"a": "1", "b": "2"},
			expected: map[string]string{"a": "1", "b": "2"},
		},
		{
			name:     "base has values, override empty",
			base:     map[string]string{"a": "1", "b": "2"},
			override: map[string]string{},
			expected: map[string]string{"a": "1", "b": "2"},
		},
		{
			name:     "no overlap",
			base:     map[string]string{"a": "1", "b": "2"},
			override: map[string]string{"c": "3", "d": "4"},
			expected: map[string]string{"a": "1", "b": "2", "c": "3", "d": "4"},
		},
		{
			name:     "override replaces base values",
			base:     map[string]string{"a": "1", "b": "2"},
			override: map[string]string{"b": "20", "c": "3"},
			expected: map[string]string{"a": "1", "b": "20", "c": "3"},
		},
		{
			name:     "override completely replaces",
			base:     map[string]string{"a": "1", "b": "2"},
			override: map[string]string{"a": "10", "b": "20"},
			expected: map[string]string{"a": "10", "b": "20"},
		},
		{
			name:     "empty string override",
			base:     map[string]string{"a": "1", "b": "2"},
			override: map[string]string{"b": ""},
			expected: map[string]string{"a": "1", "b": ""},
		},
		{
			name:     "single key in both",
			base:     map[string]string{"key": "base"},
			override: map[string]string{"key": "override"},
			expected: map[string]string{"key": "override"},
		},
		{
			name:     "complex merge",
			base:     map[string]string{"a": "1", "b": "2", "c": "3", "d": "4"},
			override: map[string]string{"b": "20", "e": "5", "f": "6"},
			expected: map[string]string{"a": "1", "b": "20", "c": "3", "d": "4", "e": "5", "f": "6"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeAttributes(tt.base, tt.override)
			
			// Check length
			if len(result) != len(tt.expected) {
				t.Errorf("MergeAttributes(%v, %v) length = %d, want %d", tt.base, tt.override, len(result), len(tt.expected))
			}
			
			// Check all expected keys and values
			for key, expectedValue := range tt.expected {
				if actualValue, ok := result[key]; !ok {
					t.Errorf("MergeAttributes(%v, %v) missing key %q", tt.base, tt.override, key)
				} else if actualValue != expectedValue {
					t.Errorf("MergeAttributes(%v, %v)[%q] = %q, want %q", tt.base, tt.override, key, actualValue, expectedValue)
				}
			}
			
			// Check no extra keys
			for key := range result {
				if _, ok := tt.expected[key]; !ok {
					t.Errorf("MergeAttributes(%v, %v) has unexpected key %q", tt.base, tt.override, key)
				}
			}
		})
	}
	
	// Test that base and override are not modified
	t.Run("base and override not modified", func(t *testing.T) {
		base := map[string]string{"a": "1", "b": "2"}
		override := map[string]string{"b": "20", "c": "3"}
		originalBase := make(map[string]string)
		originalOverride := make(map[string]string)
		for k, v := range base {
			originalBase[k] = v
		}
		for k, v := range override {
			originalOverride[k] = v
		}
		
		_ = MergeAttributes(base, override)
		
		// Check base unchanged
		for k, v := range base {
			if originalBase[k] != v {
				t.Errorf("Base map was modified: key %q changed from %q to %q", k, originalBase[k], v)
			}
		}
		
		// Check override unchanged
		for k, v := range override {
			if originalOverride[k] != v {
				t.Errorf("Override map was modified: key %q changed from %q to %q", k, originalOverride[k], v)
			}
		}
	})
}

