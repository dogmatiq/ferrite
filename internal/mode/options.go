package mode

import (
	"io"
	"os"

	"github.com/dogmatiq/ferrite/variable"
)

// Options is a set of options for running a mode.
type Options struct {
	Registry *variable.Registry
	Args     []string
	Out      io.Writer
	Err      io.Writer
	Exit     func(int)
}

// DefaultOptions is the default set of options for running a mode.
var DefaultOptions Options

// ResetDefaultOptions resets DefaultOptions to its default value. This is
// largely intended for tearing down tests.
func ResetDefaultOptions() {
	DefaultOptions = Options{
		&variable.DefaultRegistry,
		os.Args,
		os.Stdout,
		os.Stderr,
		os.Exit,
	}
}

func init() {
	ResetDefaultOptions()
}
