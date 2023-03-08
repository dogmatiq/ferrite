package ferrite

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dogmatiq/ferrite/internal/markdownmode"
	"github.com/dogmatiq/ferrite/internal/validatemode"
	"github.com/dogmatiq/ferrite/variable"
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
	switch m := os.Getenv("FERRITE_MODE"); m {
	case "usage/markdown":
		result := markdownmode.Run(
			filepath.Base(os.Args[0]),
			&variable.DefaultRegistry,
		)
		io.WriteString(output, result)
		exit(0)

	case "validate", "":
		if result, ok := validatemode.Run(&variable.DefaultRegistry); !ok {
			io.WriteString(output, result)
			exit(1)
		}

	default:
		fmt.Fprintf(output, "unrecognized FERRITE_MODE (%s)\n", m)
		exit(1)
	}
}

var (
	// output is the writer to which the validation result is written.
	output io.Writer = os.Stderr

	// exit is called to exit the process when validation fails.
	exit = os.Exit
)

// An InitOption changes the behavior of Init().
//
// WARNING: The signature of this function is not considered part of Ferrite's
// public API. It may change at any time and without warning.
type InitOption struct{}
