package spec

import (
	"regexp"
	"strings"
)

var needsQuotes = regexp.MustCompile(`[^\w@%+=:,./-]`)

// Escape escapes a string for use in an environment variable.
func Escape(s string) string {
	if needsQuotes.MatchString(s) {
		return `'` + strings.ReplaceAll(s, `'`, `'"'"'`) + `'`
	}

	return s
}
