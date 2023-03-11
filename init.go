package ferrite

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite/internal/mode"
	"github.com/dogmatiq/ferrite/internal/mode/export/dotenv"
	"github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	"github.com/dogmatiq/ferrite/internal/mode/validate"
)

// Init initializes ferrite.
//
// By default it validates all of the environment variables that have been
// defined. If any environment variables are invalid it prints a report
// describing the problems and exits the process with a non-zero exit code.
//
// If the FERRITE_MODE environment variable is set to "usage/markdown" it prints
// information about the environment variables in Markdown format, then exits
// the process successfully.
func Init(options ...InitOption) {
	opts := mode.DefaultOptions
	for _, opt := range options {
		opt.applyInitOption(&opts)
	}

	switch m := os.Getenv("FERRITE_MODE"); m {
	case "validate", "":
		validate.Run(opts)
	case "usage/markdown":
		markdown.Run(opts)
	case "export/dotenv":
		dotenv.Run(opts)
	default:
		fmt.Fprintf(opts.Err, "unrecognized FERRITE_MODE (%s)\n", m)
		opts.Exit(1)
	}
}

// An InitOption changes the behavior of Init().
type InitOption interface {
	applyInitOption(*mode.Options)
}
