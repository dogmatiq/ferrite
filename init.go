package ferrite

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite/internal/mode"
	"github.com/dogmatiq/ferrite/internal/mode/export/dotenv"
	"github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	"github.com/dogmatiq/ferrite/internal/mode/validate"
)

// Init initializes Ferrite.
//
// Different modes can be selected by setting the `FERRITE_MODE` environment
// variable.
//
// "validate" mode: This is the default mode. If one or more environment
// variables are invalid, this mode renders a description of all declared
// environment variables and their associated values and validation failures to
// `STDERR`, then exits the process with a non-zero exit code.
//
// It also shows warnings if deprecated environment variables are used.
//
// "usage/markdown" mode: This mode renders Markdown documentation about the
// environment variables to `STDOUT`. The output is designed to be included in
// the application's `README.md` file or a similar file.
//
// "export/dotenv" mode: This mode renders environment variables to `STDOUT` in
// a format suitable for use as a `.env` file.
func Init(options ...InitOption) {
	cfg := initConfig{
		mode.DefaultConfig,
	}
	for _, opt := range options {
		opt.applyInitOption(&cfg)
	}

	switch m := os.Getenv("FERRITE_MODE"); m {
	case "validate", "":
		validate.Run(cfg.ModeConfig)
	case "usage/markdown":
		markdown.Run(cfg.ModeConfig)
	case "export/dotenv":
		dotenv.Run(cfg.ModeConfig)
	default:
		fmt.Fprintf(cfg.ModeConfig.Err, "unrecognized FERRITE_MODE (%s)\n", m)
		cfg.ModeConfig.Exit(1)
	}
}

// An InitOption changes the behavior of the Init() function.
type InitOption interface {
	applyInitOption(*initConfig)
}

// initConfig is the configuration for the Init() function, built from
// InitOption values.
type initConfig struct {
	ModeConfig mode.Config
}
