package lib

import (
	"regexp"
	"strings"
)

// SubstituteAttributes replaces attribute references in text with their values
// Attribute references are in the format {attr-name} or {attr-name}
func SubstituteAttributes(text string, attrs map[string]string) string {
	if attrs == nil || len(attrs) == 0 {
		return text
	}

	// Match attribute references: {attr-name} or {attr-name}
	// This matches { followed by attribute name (word chars, hyphens, underscores) and }
	attrRegex := regexp.MustCompile(`\{([\w\-_]+)\}`)

	result := attrRegex.ReplaceAllStringFunc(text, func(match string) string {
		// Extract attribute name (remove { and })
		attrName := match[1 : len(match)-1]
		
		// Look up attribute value
		if value, ok := attrs[attrName]; ok {
			return value
		}
		
		// If not found, return empty string (or could return the original match)
		return ""
	})

	return result
}

// GetBuiltInAttribute returns a built-in attribute value based on document attributes
func GetBuiltInAttribute(name string, docAttrs map[string]string) string {
	switch name {
	case "title":
		if title, ok := docAttrs["title"]; ok {
			return title
		}
		return ""
	case "author":
		if author, ok := docAttrs["author"]; ok {
			return author
		}
		return ""
	case "revnumber":
		if revnumber, ok := docAttrs["revnumber"]; ok {
			return revnumber
		}
		return ""
	case "revdate":
		if revdate, ok := docAttrs["revdate"]; ok {
			return revdate
		}
		return ""
	case "revremark":
		if revremark, ok := docAttrs["revremark"]; ok {
			return revremark
		}
		return ""
	default:
		// Check if it's a custom attribute (with : prefix)
		if strings.HasPrefix(name, ":") {
			key := name[1:]
			if value, ok := docAttrs[":"+key]; ok {
				return value
			}
		}
		// Check custom attributes without : prefix
		if value, ok := docAttrs[name]; ok {
			return value
		}
		return ""
	}
}

// MergeAttributes merges two attribute maps, with the second map taking precedence
func MergeAttributes(base, override map[string]string) map[string]string {
	result := make(map[string]string)
	
	// Copy base attributes
	for k, v := range base {
		result[k] = v
	}
	
	// Override with override attributes
	for k, v := range override {
		result[k] = v
	}
	
	return result
}

