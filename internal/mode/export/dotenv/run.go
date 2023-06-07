package dotenv

import (
	"github.com/dogmatiq/ferrite/internal/mode"
	"github.com/dogmatiq/ferrite/internal/mode/internal/render"
	"github.com/dogmatiq/ferrite/internal/variable"
	"github.com/dogmatiq/iago/must"
)

// Run generates and env file describing the environment variables and their
// current values.
func Run(cfg mode.Config) {
	for i, v := range cfg.Registry.Variables() {
		s := v.Spec()

		if i > 0 {
			must.Fprintf(cfg.Out, "\n")
		}

		must.Fprintf(cfg.Out, "# %s (", s.Description())

		if def, ok := s.Default(); ok {
			must.WriteString(cfg.Out, "default: ")
			must.WriteString(cfg.Out, render.Value(s, def))
		} else if s.IsDeprecated() {
			must.Fprintf(cfg.Out, "deprecated")
		} else if s.IsRequired() {
			must.Fprintf(cfg.Out, "required")
		} else {
			must.Fprintf(cfg.Out, "optional")
		}

		if s.IsSensitive() {
			must.Fprintf(cfg.Out, ", sensitive")
		}

		must.Fprintf(cfg.Out, ")\n")
		must.Fprintf(cfg.Out, "export %s=", s.Name())

		if v.Source() == variable.SourceEnvironment {
			err := v.Error()
			if err, ok := err.(variable.ValueError); ok {
				must.Fprintf(
					cfg.Out,
					" # %s is invalid: %s",
					err.Literal().Quote(),
					err.Unwrap(),
				)
			} else {
				value := v.Value()

				must.Fprintf(
					cfg.Out,
					"%s",
					value.Verbatim().Quote(),
				)

				if value.Verbatim() != value.Canonical() {
					must.Fprintf(
						cfg.Out,
						" # equivalent to %s",
						value.Canonical().Quote(),
					)
				}
			}
		}

		must.Fprintf(cfg.Out, "\n")
	}

	cfg.Exit(0)
}
