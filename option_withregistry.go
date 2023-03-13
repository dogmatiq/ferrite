package ferrite

import "github.com/dogmatiq/ferrite/variable"

// WithRegistry is an option that sets the variable registry to use.
func WithRegistry(reg *variable.Registry) interface {
	InitOption
	RequiredOption
	OptionalOption
	DeprecatedOption
} {
	if reg == nil {
		panic("registry must not be nil")
	}

	return option{
		ApplyToInitConfig: func(opts *initConfig) {
			opts.ModeConfig.Registry = reg
		},
		ApplyToSetConfig: func(opts *variableSetConfig) {
			opts.Registry = reg
		},
	}
}
