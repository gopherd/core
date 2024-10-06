package stringutil

import (
	"strings"
	"unicode"
)

// Capitalize returns the string with the first letter capitalized.
func Capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	return string(unicode.ToUpper(r[0])) + string(r[1:])
}

// Uncapitalize returns the string with the first letter uncapitalized.
func Uncapitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	return string(unicode.ToLower(r[0])) + string(r[1:])
}

// Rename renames the string with the given convert function and separator.
func Rename(s string, convert func(int, string) string, sep string) string {
	if s == "" {
		return ""
	}

	var result strings.Builder
	var word strings.Builder
	var count int

	runes := []rune(s)
	n := len(runes)
	i := 0

	for i < n {
		// Skip non-alphanumeric characters, treat them as word boundaries
		for i < n && !unicode.IsLetter(runes[i]) && !unicode.IsDigit(runes[i]) {
			i++
		}

		if i >= n {
			break
		}

		word.Reset()

		// Collect characters for the current word
		for i < n && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i])) {
			r := runes[i]
			word.WriteRune(r)
			i++

			if i < n && isWordBoundary(r, runes[i], i, runes) {
				break
			}
		}

		if word.Len() > 0 {
			if result.Len() > 0 && sep != "" {
				result.WriteString(sep)
			}

			convertedWord := convert(count, word.String())
			result.WriteString(convertedWord)
			count++
		}
	}

	return result.String()
}

func isWordBoundary(prev rune, curr rune, index int, runes []rune) bool {
	// If previous character is lowercase and current is uppercase, it's a boundary
	if unicode.IsLower(prev) && unicode.IsUpper(curr) {
		return true
	}

	// Handle acronyms (e.g., "HTTPServer")
	if unicode.IsUpper(prev) && unicode.IsUpper(curr) {
		// If next character exists and is lowercase, split before current character
		if index+1 < len(runes) && unicode.IsLower(runes[index+1]) {
			return true
		}
		return false
	}

	// If previous is digit and current is letter, it's a boundary
	if unicode.IsDigit(prev) && unicode.IsLetter(curr) {
		return true
	}

	// Do not split when transitioning from letter to digit
	if unicode.IsLetter(prev) && unicode.IsDigit(curr) {
		return false
	}

	// If both are letters (regardless of case), do not split
	if unicode.IsLetter(prev) && unicode.IsLetter(curr) {
		return false
	}

	// If both are digits, do not split
	if unicode.IsDigit(prev) && unicode.IsDigit(curr) {
		return false
	}

	// If one is letter/digit and the other is not, it's a boundary
	if (unicode.IsLetter(prev) || unicode.IsDigit(prev)) != (unicode.IsLetter(curr) || unicode.IsDigit(curr)) {
		return true
	}

	// Default to no boundary
	return false
}

func lowerAll(i int, s string) string {
	return strings.ToLower(s)
}

func capitalizeAll(i int, s string) string {
	return Capitalize(s)
}

func capitalizeExceptFirst(i int, s string) string {
	if i == 0 {
		return strings.ToLower(s)
	}
	return Capitalize(s)
}

// SnakeCase converts the string to snake_case.
func SnakeCase(s string) string {
	return Rename(s, lowerAll, "_")
}

// KebabCase converts the string to kebab-case.
func KebabCase(s string) string {
	return Rename(s, lowerAll, "-")
}

// CamelCase converts the string to camelCase.
func CamelCase(s string) string {
	return Rename(s, capitalizeExceptFirst, "")
}

// PascalCase converts the string to PascalCase.
func PascalCase(s string) string {
	return Rename(s, capitalizeAll, "")
}
