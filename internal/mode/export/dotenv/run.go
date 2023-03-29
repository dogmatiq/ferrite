package dotenv

import (
	"strings"

	"github.com/dogmatiq/ferrite/internal/mode"
	"github.com/dogmatiq/ferrite/variable"
	"github.com/dogmatiq/iago/must"
)

// Run generates and env file describing the environment variables and their
// current values.
func Run(opts mode.Config) {
	for i, v := range opts.Registry.Variables() {
		s := v.Spec()

		if i > 0 {
			must.Fprintf(opts.Out, "\n")
		}

		must.Fprintf(opts.Out, "# %s (", s.Description())

		if def, ok := s.Default(); ok {
			x := def.Quote()
			if s.IsSensitive() {
				x = strings.Repeat("*", len(def.String))
			}
			must.Fprintf(opts.Out, "default: %s", x)
		} else if s.IsDeprecated() {
			must.Fprintf(opts.Out, "deprecated")
		} else if s.IsRequired() {
			must.Fprintf(opts.Out, "required")
		} else {
			must.Fprintf(opts.Out, "optional")
		}

		if s.IsSensitive() {
			must.Fprintf(opts.Out, ", sensitive")
		}

		must.Fprintf(opts.Out, ")\n")
		must.Fprintf(opts.Out, "export %s=", s.Name())

		if v.Source() == variable.SourceEnvironment && !s.IsSensitive() {
			err := v.Error()
			if err, ok := err.(variable.ValueError); ok {
				must.Fprintf(
					opts.Out,
					" # %s is invalid: %s",
					err.Literal().Quote(),
					err.Unwrap(),
				)
			} else {
				value := v.Value()

				must.Fprintf(
					opts.Out,
					"%s",
					value.Verbatim().Quote(),
				)

				if value.Verbatim() != value.Canonical() {
					must.Fprintf(
						opts.Out,
						" # equivalent to %s",
						value.Canonical().Quote(),
					)
				}
			}
		}

		must.Fprintf(opts.Out, "\n")
	}

	opts.Exit(0)
}
