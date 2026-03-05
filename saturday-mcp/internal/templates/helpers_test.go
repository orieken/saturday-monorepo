package templates

import (
	"testing"
)

func TestPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "HelloWorld"},
		{"hello_world", "HelloWorld"},
		{"hello-world", "HelloWorld"},
		{"helloWorld", "HelloWorld"},
		{"HelloWorld", "HelloWorld"},
		{"", ""},
		{"a", "A"},
		{"my_var_name", "MyVarName"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := PascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("PascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "helloWorld"},
		{"hello_world", "helloWorld"},
		{"hello-world", "helloWorld"},
		{"HelloWorld", "helloWorld"},
		{"", ""},
		{"A", "a"},
		{"my_var_name", "myVarName"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := CamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("CamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello_world"},
		{"helloWorld", "hello_world"},
		{"hello world", "hello_world"},
		{"hello-world", "hello_world"},
		{"", ""},
		{"A", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := SnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("SnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello-world"},
		{"helloWorld", "hello-world"},
		{"hello world", "hello-world"},
		{"hello_world", "hello-world"},
		{"", ""},
		{"A", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := KebabCase(tt.input)
			if result != tt.expected {
				t.Errorf("KebabCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDefaultFuncMap(t *testing.T) {
	funcMap := DefaultFuncMap()

	expectedFuncs := []string{
		"pascalCase", "camelCase", "snakeCase", "kebabCase",
		"upper", "lower", "title", "trim",
		"join", "split", "contains", "hasPrefix", "hasSuffix",
		"replace", "repeat",
	}

	for _, name := range expectedFuncs {
		if _, exists := funcMap[name]; !exists {
			t.Errorf("Expected function %q to exist in FuncMap", name)
		}
	}
}
