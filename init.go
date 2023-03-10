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
	opts := defaultInitOptions // copy
	for _, opt := range options {
		opt.applyInitOption(&opts)
	}

	switch m := os.Getenv("FERRITE_MODE"); m {
	case "usage/markdown":
		result := markdownmode.Run(
			filepath.Base(os.Args[0]),
			&variable.DefaultRegistry,
		)
		io.WriteString(opts.Err, result)
		opts.Exit(0)

	case "validate", "":
		if result, ok := validatemode.Run(&variable.DefaultRegistry); !ok {
			io.WriteString(opts.Err, result)
			opts.Exit(1)
		}

	default:
		fmt.Fprintf(opts.Err, "unrecognized FERRITE_MODE (%s)\n", m)
		opts.Exit(1)
	}
}

// An InitOption changes the behavior of Init().
type InitOption interface {
	applyInitOption(*initOptions)
}

type initOptions struct {
	Out  io.Writer
	Err  io.Writer
	Exit func(int)
}

var defaultInitOptions = initOptions{
	Out:  os.Stdout,
	Err:  os.Stderr,
	Exit: os.Exit,
}
