package validatemode

import (
	"io"

	"github.com/dogmatiq/ferrite/variable"
)

// Run validates the variables in the given registry.
//
// w is the target for human-readable usage and validation information intended
// for display in the console.
//
// It returns true if all variables are valid.
func Run(reg *variable.Registry, w io.Writer) bool {
	show := false
	valid := true

	t := table{}
	for _, v := range reg.Variables() {
		t.AddRow(
			name(v),
			description(v),
			spec(v),
			value(v),
		)

		if needsAttention(v) {
			show = true
		}

		if !v.IsValid() {
			valid = false
		}
	}

	if show {
		if _, err := io.WriteString(w, "Environment Variables:\n\n"); err != nil {
			panic(err)
		}

		if _, err := t.WriteTo(w); err != nil {
			panic(err)
		}

		if _, err := io.WriteString(w, "\n"); err != nil {
			panic(err)
		}
	}

	return valid
}

const (
	iconOK        = "✓"
	iconWarn      = "⚠"
	iconError     = "✗"
	iconNeutral   = "•"
	iconAttention = "❯"
)

// needsAttention returns true if v needs attention from the user.
func needsAttention(v variable.Any) bool {
	s := v.Spec()

	if !v.IsValid() {
		return true
	}

	if s.IsDeprecated() && v.IsExplicit() {
		return true
	}

	return false
}
