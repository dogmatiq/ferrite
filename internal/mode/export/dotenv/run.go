package dotenv

import (
	"fmt"
	"io"

	"github.com/dogmatiq/ferrite/internal/mode"
	"github.com/dogmatiq/ferrite/internal/mode/internal/render"
	"github.com/dogmatiq/ferrite/internal/variable"
)

// Run generates and env file describing the environment variables and their
// current values.
func Run(cfg mode.Config) {
	for i, v := range cfg.Registries.Variables() {
		s := v.Spec()

		if i > 0 {
			fprintf(cfg.Out, "\n")
		}

		fprintf(cfg.Out, "# %s (", s.Description())

		if def, ok := s.Default(); ok {
			fprintf(cfg.Out, "default: %s", render.Value(s, def))
		} else if s.IsDeprecated() {
			fprintf(cfg.Out, "deprecated")
		} else if s.IsRequired() {
			fprintf(cfg.Out, "required")
		} else {
			fprintf(cfg.Out, "optional")
		}

		if s.IsSensitive() {
			fprintf(cfg.Out, ", sensitive")
		}

		fprintf(cfg.Out, ")\n")
		fprintf(cfg.Out, "export %s=", s.Name())

		if v.Source() == variable.SourceEnvironment {
			err := v.Error()
			if err, ok := err.(variable.ValueError); ok {
				fprintf(
					cfg.Out,
					" # %s is invalid: %s",
					err.Literal().Quote(),
					err.Unwrap(),
				)
			} else {
				value := v.Value()

				fprintf(
					cfg.Out,
					"%s",
					value.Verbatim().Quote(),
				)

				if value.Verbatim() != value.Canonical() {
					fprintf(
						cfg.Out,
						" # equivalent to %s",
						value.Canonical().Quote(),
					)
				}
			}
		}

		fprintf(cfg.Out, "\n")
	}

	cfg.Exit(0)
}

func fprintf(w io.Writer, format string, args ...any) {
	if _, err := fmt.Fprintf(w, format, args...); err != nil {
		panic(err)
	}
}
