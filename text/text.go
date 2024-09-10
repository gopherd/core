package text

import (
	"fmt"
	"regexp"
)

// ContainsWord returns true if the given word is found in the string s.
func ContainsWord(s, word string) bool {
	pattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(word))
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}
