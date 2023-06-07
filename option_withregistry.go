package ferrite

import "github.com/dogmatiq/ferrite/internal/variable"

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
		ApplyToInitConfig: func(cfg *initConfig) {
			cfg.ModeConfig.Registry = reg
		},
		ApplyToSetConfig: func(cfg *variableSetConfig) {
			cfg.Registry = reg
		},
	}
}
