package markdown

import (
	"strings"
	"unicode"
)

// headingSlug converts a markdown heading text to a GitHub-compatible anchor
// slug, using the same algorithm as github.com/Flet/github-slugger.
//
// It lowercases the input, strips characters that match the github-slugger
// removal regex (keeping letters, digits, hyphens, underscores, and spaces),
// and converts spaces to hyphens.
func headingSlug(text string) string {
	var b strings.Builder

	for _, r := range strings.ToLower(text) {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' || r == '-':
			b.WriteRune(r)
		case r == ' ':
			b.WriteRune('-')
		}
	}

	return b.String()
}
