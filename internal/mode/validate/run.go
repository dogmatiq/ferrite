package validate

import (
	"io"

	"github.com/dogmatiq/ferrite/internal/mode"
	"github.com/dogmatiq/ferrite/internal/variable"
)

// Run validates the variables in the given registry.
//
// w is the target for human-readable usage and validation information intended
// for display in the console.
//
// It returns true if all variables are valid.
func Run(cfg mode.Config) {
	show := false
	valid := true

	t := table{}
	for _, v := range cfg.Registry.Variables() {
		t.AddRow(
			name(v),
			description(v),
			spec(v),
			value(v),
		)

		switch attentionNeeded(v) {
		case attentionWarning:
			show = true
		case attentionError:
			show = true
			valid = false
		}
	}

	if show {
		if _, err := io.WriteString(cfg.Err, "Environment Variables:\n\n"); err != nil {
			panic(err)
		}

		if _, err := t.WriteTo(cfg.Err); err != nil {
			panic(err)
		}

		if _, err := io.WriteString(cfg.Err, "\n"); err != nil {
			panic(err)
		}
	}

	if !valid {
		cfg.Exit(1)
	}
}

const (
	iconOK        = "✓"
	iconWarn      = "⚠"
	iconError     = "✗"
	iconNeutral   = "•"
	iconAttention = "❯"
)

type attentionLevel int

const (
	attentionNone attentionLevel = iota
	attentionWarning
	attentionError
)

// attentionNeeded returns true if v needs attention from the user.
func attentionNeeded(v variable.Any) attentionLevel {
	s := v.Spec()

	if err := v.Error(); err != nil {
		if v.Availability() != variable.AvailabilityIgnored {
			return attentionError
		}

		if _, ok := err.(variable.ValueError); ok {
			return attentionWarning
		}
	}

	if s.IsDeprecated() && v.Source() == variable.SourceEnvironment {
		return attentionWarning
	}

	return attentionNone
}
