package stringutil

import (
	"testing"
)

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"a", "A"},
		{"A", "A"},
		{"hello", "Hello"},
		{"Hello", "Hello"},
		{"1hello", "1hello"},
		{" hello", " hello"},
		{"hELLO", "HELLO"},
		{"h", "H"},
	}

	for _, tt := range tests {
		got := Capitalize(tt.input)
		if got != tt.want {
			t.Errorf("Capitalize(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}

func TestUncapitalize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"A", "a"},
		{"a", "a"},
		{"Hello", "hello"},
		{"hello", "hello"},
		{"1Hello", "1Hello"},
		{" Hello", " Hello"},
		{"HELLO", "hELLO"},
		{"H", "h"},
	}

	for _, tt := range tests {
		got := Uncapitalize(tt.input)
		if got != tt.want {
			t.Errorf("Uncapitalize(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"simple", "simple"},
		{"SimpleTestCase", "simple_test_case"},
		{"HTTPServer", "http_server"},
		{"xmlHTTPRequest", "xml_http_request"},
		{"MyID", "my_id"},
		{"My123ID", "my123_id"},
		{"HelloWorld", "hello_world"},
		{"helloWorld", "hello_world"},
		{"hello_world", "hello_world"},
		{"Hello_World", "hello_world"},
		{"hello-world", "hello_world"},
	}

	for _, tt := range tests {
		got := SnakeCase(tt.input)
		if got != tt.want {
			t.Errorf("SnakeCase(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}

func TestKebabCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"simple", "simple"},
		{"SimpleTestCase", "simple-test-case"},
		{"HTTPServer", "http-server"},
		{"xmlHTTPRequest", "xml-http-request"},
		{"MyID", "my-id"},
		{"My123ID", "my123-id"},
		{"HelloWorld", "hello-world"},
		{"helloWorld", "hello-world"},
		{"hello_world", "hello-world"},
		{"Hello_World", "hello-world"},
	}

	for _, tt := range tests {
		got := KebabCase(tt.input)
		if got != tt.want {
			t.Errorf("KebabCase(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}

func TestCamelCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"simple", "simple"},
		{"SimpleTestCase", "simpleTestCase"},
		{"HTTPServer", "httpServer"},
		{"xmlHTTPRequest", "xmlHTTPRequest"},
		{"MyID", "myID"},
		{"My123ID", "my123ID"},
		{"HelloWorld", "helloWorld"},
		{"helloWorld", "helloWorld"},
		{"hello_world", "helloWorld"},
		{"Hello_World", "helloWorld"},
		{"hello world", "helloWorld"},
		{"hello World", "helloWorld"},
		{"Hello World", "helloWorld"},
		{"hello-world", "helloWorld"},
		{"hello-World", "helloWorld"},
		{"Hello-World", "helloWorld"},
	}

	for _, tt := range tests {
		got := CamelCase(tt.input)
		if got != tt.want {
			t.Errorf("CamelCase(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}

func TestPascalCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"simple", "Simple"},
		{"SimpleTestCase", "SimpleTestCase"},
		{"HTTPServer", "HTTPServer"},
		{"xmlHTTPRequest", "XmlHTTPRequest"},
		{"MyID", "MyID"},
		{"My123ID", "My123ID"},
		{"HelloWorld", "HelloWorld"},
		{"helloWorld", "HelloWorld"},
		{"hello_world", "HelloWorld"},
		{"Hello_World", "HelloWorld"},
	}

	for _, tt := range tests {
		got := PascalCase(tt.input)
		if got != tt.want {
			t.Errorf("PascalCase(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}
