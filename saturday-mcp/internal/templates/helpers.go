package templates

import (
	"strings"
	"text/template"
	"unicode"
)

// DefaultFuncMap returns the default template function map
func DefaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"pascalCase":  PascalCase,
		"camelCase":   CamelCase,
		"snakeCase":   SnakeCase,
		"kebabCase":   KebabCase,
		"upper":       strings.ToUpper,
		"lower":       strings.ToLower,
		"title":       strings.Title,
		"trim":        strings.TrimSpace,
		"join":        strings.Join,
		"split":       strings.Split,
		"contains":    strings.Contains,
		"hasPrefix":   strings.HasPrefix,
		"hasSuffix":   strings.HasSuffix,
		"replace":     strings.ReplaceAll,
		"repeat":      strings.Repeat,
	}
}

// PascalCase converts a string to PascalCase
func PascalCase(s string) string {
	if s == "" {
		return ""
	}

	words := splitWords(s)
	var result strings.Builder

	for _, word := range words {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(string(word[0])))
			if len(word) > 1 {
				result.WriteString(strings.ToLower(word[1:]))
			}
		}
	}

	return result.String()
}

// CamelCase converts a string to camelCase
func CamelCase(s string) string {
	if s == "" {
		return ""
	}

	pascal := PascalCase(s)
	if len(pascal) == 0 {
		return ""
	}

	return strings.ToLower(string(pascal[0])) + pascal[1:]
}

// SnakeCase converts a string to snake_case
func SnakeCase(s string) string {
	if s == "" {
		return ""
	}

	words := splitWords(s)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}

	return strings.Join(words, "_")
}

// KebabCase converts a string to kebab-case
func KebabCase(s string) string {
	if s == "" {
		return ""
	}

	words := splitWords(s)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}

	return strings.Join(words, "-")
}

// splitWords splits a string into words based on various delimiters and case changes
func splitWords(s string) []string {
	var words []string
	var current strings.Builder

	for i, r := range s {
		if r == '_' || r == '-' || r == ' ' || r == '.' {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
			continue
		}

		if i > 0 && unicode.IsUpper(r) && unicode.IsLower(rune(s[i-1])) {
			// Transition from lowercase to uppercase (e.g., "myVar" -> "my", "Var")
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
		}

		current.WriteRune(r)
	}

	if current.Len() > 0 {
		words = append(words, current.String())
	}

	return words
}
