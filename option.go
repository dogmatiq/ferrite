package ferrite

import "github.com/dogmatiq/ferrite/internal/variable"

// option is an implementation of various option interfaces.
//
// Functions that return options should return an anonymous interface type that
// embeds one or more of the option interfaces.
type option struct {
	ApplyToInitConfig func(*initConfig)
	ApplyToRegistry   func(*variable.Registry)

	ApplyToSetConfig           func(*variableSetConfig)
	ApplyToSpec                func(variable.SpecBuilder)
	ApplyToRequiredSetConfig   func(*variableSetConfig)
	ApplyToSpecInRequiredSet   func(variable.SpecBuilder)
	ApplyToOptionalSetConfig   func(*variableSetConfig)
	ApplyToSpecInOptionalSet   func(variable.SpecBuilder)
	ApplyToDeprecatedSetConfig func(*variableSetConfig)
	ApplyToSpecInDeprecatedSet func(variable.SpecBuilder)
}

func (o option) applyInitOption(cfg *initConfig) {
	applyOption(cfg, o.ApplyToInitConfig)
}

func (o option) applyRegistryOption(reg *variable.Registry) {
	applyOption(reg, o.ApplyToRegistry)
}

func (o option) applyRequiredOptionToConfig(cfg *variableSetConfig) {
	applyOption(cfg, o.ApplyToSetConfig, o.ApplyToRequiredSetConfig)
}

func (o option) applyRequiredOptionToSpec(b variable.SpecBuilder) {
	applyOption(b, o.ApplyToSpec, o.ApplyToSpecInRequiredSet)
}

func (o option) applyOptionalOptionToConfig(cfg *variableSetConfig) {
	applyOption(cfg, o.ApplyToSetConfig, o.ApplyToOptionalSetConfig)
}

func (o option) applyOptionalOptionToSpec(b variable.SpecBuilder) {
	applyOption(b, o.ApplyToSpec, o.ApplyToSpecInOptionalSet)
}

func (o option) applyDeprecatedOptionToConfig(cfg *variableSetConfig) {
	applyOption(cfg, o.ApplyToSetConfig, o.ApplyToDeprecatedSetConfig)
}

func (o option) applyDeprecatedOptionToSpec(b variable.SpecBuilder) {
	applyOption(b, o.ApplyToSpec, o.ApplyToSpecInDeprecatedSet)
}

func applyOption[T any](cfg T, funcs ...func(T)) {
	for _, fn := range funcs {
		if fn != nil {
			fn(cfg)
		}
	}
}
