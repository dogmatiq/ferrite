package ferrite

import "github.com/dogmatiq/ferrite/variable"

// option is an implementation of various option interfaces.
//
// Functions that return options should return an anonymous interface type that
// embeds one or more of the option interfaces.
type option struct {
	ApplyToInitConfig func(*initConfig)

	ApplyToSetConfig           func(*variableSetConfig)
	ApplyToSpec                func(variable.SpecBuilder)
	ApplyToRequiredSetConfig   func(*variableSetConfig)
	ApplyToSpecInRequiredSet   func(variable.SpecBuilder)
	ApplyToOptionalSetConfig   func(*variableSetConfig)
	ApplyToSpecInOptionalSet   func(variable.SpecBuilder)
	ApplyToDeprecatedSetConfig func(*variableSetConfig)
	ApplyToSpecInDeprecatedSet func(variable.SpecBuilder)

	ApplyToRefersToRelationship   func(*variable.RefersTo)
	ApplyToSupersedesRelationship func(*variable.Supersedes)
}

func (o option) applyInitOption(opts *initConfig) {
	applyOption(opts, o.ApplyToInitConfig)
}

func (o option) applyRequiredOptionToConfig(opts *variableSetConfig) {
	applyOption(opts, o.ApplyToSetConfig, o.ApplyToRequiredSetConfig)
}

func (o option) applyRequiredOptionToSpec(b variable.SpecBuilder) {
	applyOption(b, o.ApplyToSpec, o.ApplyToSpecInRequiredSet)
}

func (o option) applyOptionalOptionToConfig(opts *variableSetConfig) {
	applyOption(opts, o.ApplyToSetConfig, o.ApplyToOptionalSetConfig)
}

func (o option) applyOptionalOptionToSpec(b variable.SpecBuilder) {
	applyOption(b, o.ApplyToSpec, o.ApplyToSpecInOptionalSet)
}

func (o option) applyDeprecatedOptionToConfig(opts *variableSetConfig) {
	applyOption(opts, o.ApplyToSetConfig, o.ApplyToDeprecatedSetConfig)
}

func (o option) applyDeprecatedOptionToSpec(b variable.SpecBuilder) {
	applyOption(b, o.ApplyToSpec, o.ApplyToSpecInDeprecatedSet)
}

func (o option) applyRefersToOption(r *variable.RefersTo) {
	applyOption(r, o.ApplyToRefersToRelationship)
}

func (o option) applySupersedesOption(r *variable.Supersedes) {
	applyOption(r, o.ApplyToSupersedesRelationship)
}

func applyOption[T any](cfg T, funcs ...func(T)) {
	for _, fn := range funcs {
		if fn != nil {
			fn(cfg)
		}
	}
}
