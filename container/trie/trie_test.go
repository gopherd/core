package trie_test

import (
	"testing"

	"github.com/gopherd/core/container/trie"
)

func TestAdd(t *testing.T) {
	words := []string{"hello", "world", "hi", "hit", "holo", "中文"}
	tr := trie.New()
	for _, word := range words {
		tr.Add(word)
	}
	for _, word := range words {
		if !tr.Has(word) {
			t.Fatalf("word %q not found", word)
		}
		if !tr.HasPrefix(word) {
			t.Fatalf("prefix %q not found", word)
		}
	}
	type testCase struct {
		word   string
		prefix string
		has    bool
	}
	for i, tc := range []testCase{
		{"hi", "h", true},
		{"hit", "hi", true},
		{"hello", "ho", true},
		{"holo", "hol", true},
		{"holo", "holo", true},
		{"world", "w", true},
		{"world", "wo", true},
		{"world", "wor", true},
		{"world", "worl", true},
		{"world", "world", true},
		{"w", "i", false},
		{"wo", "ij", false},
		{"wor", "ello", false},
		{"worl", "it", false},
		{"中", "文字", false},
		{"中国", "文字", false},
		{"中文", "中", true},
		{"中文", "中文", true},
		{"文字", "文", false},
	} {
		if tr.Has(tc.word) != tc.has {
			t.Fatalf("%dth: Has: want %v, got %v", i, tc.has, !tc.has)
		}
		if tr.HasPrefix(tc.prefix) != tc.has {
			t.Fatalf("%dth: HasPrefix: want %v, got %v", i, tc.has, !tc.has)
		}
	}
	t.Logf("trie:\n%v", tr.String())
}

func TestRemove(t *testing.T) {
	words := []string{"hello", "world", "hi", "hit", "holo"}
	tr := trie.New()
	for _, word := range words {
		tr.Add(word)
	}
	type removeTestCase struct {
		word    string
		removed bool
	}
	for i, tc := range []removeTestCase{
		{"h", false},
		{"hi", true},
		{"he", false},
		{"hello", true},
	} {
		if tr.Remove(tc.word) != tc.removed {
			t.Fatalf("%dth: Remove %q: want %v, got %v", i, tc.word, tc.removed, !tc.removed)
		}
	}

	type testCase struct {
		word   string
		prefix string
		has    bool
	}
	for i, tc := range []testCase{
		{"hit", "hi", true},
		{"holo", "hol", true},
		{"holo", "holo", true},
		{"world", "w", true},
		{"world", "wo", true},
		{"world", "wor", true},
		{"world", "worl", true},
		{"world", "world", true},
		{"hi", "he", false},
		{"hello", "hell", false},
		{"w", "i", false},
		{"wo", "ij", false},
		{"wor", "ello", false},
		{"worl", "it", false},
	} {
		if tr.Has(tc.word) != tc.has {
			t.Fatalf("%dth: Has: want %v, got %v", i, tc.has, !tc.has)
		}
		if tr.HasPrefix(tc.prefix) != tc.has {
			t.Fatalf("%dth: HasPrefix: want %v, got %v", i, tc.has, !tc.has)
		}
	}
	t.Logf("trie:\n%v", tr.String())
}

func TestSearch(t *testing.T) {
	tr := trie.New()
	for _, word := range []string{"hello", "world", "hi", "hit"} {
		tr.Add(word)
	}
	stringseq := func(s1, s2 []string) bool {
		if len(s1) != len(s2) {
			return false
		}
		for i := range s1 {
			if s2[i] != s1[i] {
				return false
			}
		}
		return true
	}
	type testCase struct {
		word   string
		limit  int
		result []string
	}
	for i, tc := range []testCase{
		{"w", -1, []string{"world"}},
		{"h", -1, []string{"hello", "hi", "hit"}},
		{"h", 1, []string{"hello"}},
		{"h", 2, []string{"hello", "hi"}},
		{"h", 3, []string{"hello", "hi", "hit"}},
		{"he", -1, []string{"hello"}},
		{"he", 1, []string{"hello"}},
		{"hi", -1, []string{"hi", "hit"}},
		{"hi", 1, []string{"hi"}},
		{"hi", 2, []string{"hi", "hit"}},
		{"hit", -1, []string{"hit"}},
	} {
		if got := tr.Search(tc.word, tc.limit); !stringseq(got, tc.result) {
			t.Fatalf("%dth: want %q, but got %q", i, tc.result, got)
		}
	}
}
