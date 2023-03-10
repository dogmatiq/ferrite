package validate

import (
	"io"

	"github.com/dogmatiq/ferrite/internal/mode"
	"github.com/dogmatiq/ferrite/variable"
)

// Run validates the variables in the given registry.
//
// w is the target for human-readable usage and validation information intended
// for display in the console.
//
// It returns true if all variables are valid.
func Run(opts mode.Options) {
	show := false
	valid := true

	t := table{}
	for _, v := range opts.Registry.Variables() {
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
		if _, err := io.WriteString(opts.Err, "Environment Variables:\n\n"); err != nil {
			panic(err)
		}

		if _, err := t.WriteTo(opts.Err); err != nil {
			panic(err)
		}

		if _, err := io.WriteString(opts.Err, "\n"); err != nil {
			panic(err)
		}
	}

	if !valid {
		opts.Exit(1)
	}
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
