package ferrite

import (
	"io"
	"os"
)

// Init initializes ferrite.
//
// By default it validates all of the environment variables that have been
// defined. If any environment variables are invalid it prints a report
// describing the problems and exits the process with a non-zero exit code.
//
// If the FERRITE_MODE environment variable is set to "usage/markdown" it prints
// information about the environment variables in Markdown format, then exists
// the process successfully.
func Init() {
	if result, ok := validate(); !ok {
		io.WriteString(output, result)
		exit(1)
	}
}

var (
	// output is the writer to which the validation result is written.
	output io.Writer = os.Stderr

	// exit is called to exit the process when validation fails.
	exit = os.Exit
)
