package trie

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gopherd/core/container/tree"
)

func TestNew(t *testing.T) {
	tr := New()
	if tr == nil {
		t.Fatal("New() returned nil")
	}
}

func TestAdd(t *testing.T) {
	tr := New()
	words := []string{"apple", "app", "application", "banana", ""}
	for _, word := range words {
		tr.Add(word)
	}
	for _, word := range words {
		if !tr.Has(word) {
			t.Errorf("Add() failed for word: %s", word)
		}
	}
}

func TestRemove(t *testing.T) {
	tr := New()
	words := []string{"apple", "app", "application", "banana", ""}
	for _, word := range words {
		tr.Add(word)
	}

	// Test removing existing words
	for _, word := range words {
		if !tr.Remove(word) {
			t.Errorf("Remove() failed for existing word: %s", word)
		}
		if tr.Has(word) {
			t.Errorf("Word still exists after removal: %s", word)
		}
	}

	// Test removing non-existing word
	if tr.Remove("nonexistent") {
		t.Error("Remove() returned true for non-existing word")
	}

	// Test removing empty string
	tr.Add("")
	if !tr.Remove("") {
		t.Error("Remove() failed for empty string")
	}
}

func TestHas(t *testing.T) {
	tr := New()
	words := []string{"apple", "app", "application", "banana", ""}
	for _, word := range words {
		tr.Add(word)
	}

	// Test existing words
	for _, word := range words {
		if !tr.Has(word) {
			t.Errorf("Has() returned false for existing word: %s", word)
		}
	}

	// Test non-existing words
	nonExisting := []string{"ap", "appl", "banan", "cherry"}
	for _, word := range nonExisting {
		if tr.Has(word) {
			t.Errorf("Has() returned true for non-existing word: %s", word)
		}
	}
}

func TestHasPrefix(t *testing.T) {
	tr := New()
	words := []string{"apple", "app", "application", "banana", ""}
	for _, word := range words {
		tr.Add(word)
	}

	// Test valid prefixes
	validPrefixes := []string{"", "a", "ap", "app", "appl", "ban", "banana"}
	for _, prefix := range validPrefixes {
		if !tr.HasPrefix(prefix) {
			t.Errorf("HasPrefix() returned false for valid prefix: %s", prefix)
		}
	}

	// Test invalid prefixes
	invalidPrefixes := []string{"bx", "c", "cherry"}
	for _, prefix := range invalidPrefixes {
		if tr.HasPrefix(prefix) {
			t.Errorf("HasPrefix() returned true for invalid prefix: %s", prefix)
		}
	}
}

func TestSearch(t *testing.T) {
	tr := New()
	words := []string{"apple", "app", "application", "banana", "ban", "bandana"}
	for _, word := range words {
		tr.Add(word)
	}

	testCases := []struct {
		prefix string
		limit  int
		expect []string
	}{
		{"app", 0, []string{"app", "apple", "application"}},
		{"app", 2, []string{"app", "apple"}},
		{"ban", 0, []string{"ban", "banana", "bandana"}},
		{"ban", 1, []string{"ban"}},
		{"c", 0, []string{}},
		{"", 0, []string{"app", "apple", "application", "ban", "banana", "bandana"}},
		{"", 3, []string{"app", "apple", "application"}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("prefix=%s,limit=%d", tc.prefix, tc.limit), func(t *testing.T) {
			result := tr.Search(tc.prefix, tc.limit)
			if !equal(result, tc.expect) {
				t.Errorf("Search(%q, %d) = %v; want %v", tc.prefix, tc.limit, result, tc.expect)
			}
		})
	}
}

func TestSearchAppend(t *testing.T) {
	tr := New()
	words := []string{"apple", "app", "application", "banana", "ban", "bandana"}
	for _, word := range words {
		tr.Add(word)
	}

	testCases := []struct {
		prefix   string
		limit    int
		initial  []string
		expected []string
	}{
		{"app", 0, []string{"existing"}, []string{"existing", "app", "apple", "application"}},
		{"ban", 2, []string{"prefix"}, []string{"prefix", "ban", "banana"}},
		{"c", 0, []string{"no match"}, []string{"no match"}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("prefix=%s,limit=%d", tc.prefix, tc.limit), func(t *testing.T) {
			result := tr.SearchAppend(tc.initial, tc.prefix, tc.limit)
			if !equal(result, tc.expected) {
				t.Errorf("SearchAppend(%v, %q, %d) = %v; want %v", tc.initial, tc.prefix, tc.limit, result, tc.expected)
			}
		})
	}
}

func TestString(t *testing.T) {
	tr := New()
	words := []string{"app", "apple", "banana"}
	for _, word := range words {
		tr.Add(word)
	}

	result := tr.String()
	if result == "" {
		t.Error("String() returned empty string")
	}
}

func TestStringify(t *testing.T) {
	tr := New()
	words := []string{"app", "apple", "banana"}
	for _, word := range words {
		tr.Add(word)
	}

	options := &tree.Options{
		Prefix: "  ",
	}
	result := tr.Stringify(options)
	if result == "" {
		t.Error("Stringify() returned empty string")
	}

	// Test with nil options
	nilResult := tr.Stringify(nil)
	if nilResult == "" {
		t.Error("Stringify(nil) returned empty string")
	}
}

// Helper function to compare two string slices
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestRemoveEmptyString(t *testing.T) {
	tr := New()

	if tr.Remove("") {
		t.Error("Remove() returned true for non-existing empty string")
	}

	tr.Add("")
	if !tr.Remove("") {
		t.Error("Remove() returned false for existing empty string")
	}

	if tr.Remove("") {
		t.Error("Remove() returned true for already removed empty string")
	}
}

func TestSearchAppendEdgeCases(t *testing.T) {
	tr := New()
	words := []string{"apple", "app", "application"}
	for _, word := range words {
		tr.Add(word)
	}

	result := tr.SearchAppend([]string{"existing"}, "app", 0)
	expected := []string{"existing", "app", "apple", "application"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SearchAppend with limit 0 failed. Got %v, want %v", result, expected)
	}

	result = tr.SearchAppend([]string{"existing", "word"}, "app", 1)
	expected = []string{"existing", "word"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SearchAppend with limit <= len(dst) failed. Got %v, want %v", result, expected)
	}

	result = tr.SearchAppend([]string{"existing"}, "banana", 5)
	expected = []string{"existing"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SearchAppend with non-matching prefix failed. Got %v, want %v", result, expected)
	}
}

func TestNodeAddExistingRune(t *testing.T) {
	tr := New()
	tr.Add("test")

	root := tr.root
	child := root.add('t')

	sameChild := root.add('t')

	if child != sameChild {
		t.Error("Adding existing rune should return the same node")
	}
}

func TestRemoveNonExistentRune(t *testing.T) {
	tr := New()
	tr.Add("test")

	root := tr.root
	tNode := root.search('t')

	if tNode == nil {
		t.Fatal("Failed to get 't' node")
	}

	tNode.remove('x')
	if tNode.NumChild() != 1 || tNode.search('e') == nil {
		t.Error("Removing non-existent rune should not affect existing children")
	}

	tNode.remove('s')
	if tNode.NumChild() != 1 || tNode.search('e') == nil {
		t.Error("Removing non-child rune should not affect existing children")
	}
}
