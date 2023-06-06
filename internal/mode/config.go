package mode

import (
	"io"
	"os"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// Config is the configuration used when running a mode.
type Config struct {
	Registries variable.RegistrySet
	Args       []string
	Out        io.Writer
	Err        io.Writer
	Exit       func(int)
}

// DefaultConfig is the default configuration for running a mode.
var DefaultConfig Config

// ResetDefaultConfig resets DefaultConfig to its initial value. This is
// largely intended for tearing down tests.
func ResetDefaultConfig() {
	DefaultConfig = Config{
		Args: os.Args,
		Out:  os.Stdout,
		Err:  os.Stderr,
		Exit: os.Exit,
	}
}

func init() {
	ResetDefaultConfig()
}
