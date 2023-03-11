package mode

import (
	"io"
	"os"

	"github.com/dogmatiq/ferrite/variable"
)

// Config is the configuration used when running a mode.
type Config struct {
	Registry *variable.Registry
	Args     []string
	Out      io.Writer
	Err      io.Writer
	Exit     func(int)
}

// DefaultConfig is the default configuration for running a mode.
var DefaultConfig Config

// ResetDefaultConfig resets DefaultConfig to its initial value. This is
// largely intended for tearing down tests.
func ResetDefaultConfig() {
	DefaultConfig = Config{
		&variable.DefaultRegistry,
		os.Args,
		os.Stdout,
		os.Stderr,
		os.Exit,
	}
}

func init() {
	ResetDefaultConfig()
}
