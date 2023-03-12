package ferrite

import "github.com/dogmatiq/ferrite/variable"

// WithRegistry is an option that sets the variable registry to use.
func WithRegistry(reg *variable.Registry) interface {
	InitOption
	DeprecatedOption
	OptionalOption
	RequiredOption
} {
	if reg == nil {
		panic("registry must not be nil")
	}

	return option{
		Init: func(opts *initConfig) {
			opts.ModeConfig.Registry = reg
		},
		Input: func(opts *inputConfig) {
			opts.Registry = reg
		},
	}
}

// option is an implementation of the InitOption, DeprecatedOption,
// RequiredOption and OptionalOption interfaces.
//
// It is used to implement the various options. Functions that return options
// should return an anonymous interface type that embeds one or more of the
// option interfaces.
type option struct {
	Init       func(*initConfig)
	Input      func(*inputConfig)
	Deprecated func(*deprecatedConfig)
	Optional   func(*optionalConfig)
	Required   func(*requiredConfig)
}

func (o option) applyInitOption(opts *initConfig) {
	if o.Init != nil {
		o.Init(opts)
	}
}

func (o option) applyDeprecatedOption(opts *deprecatedConfig) {
	if o.Input != nil {
		o.Input(&opts.inputConfig)
	}

	if o.Deprecated != nil {
		o.Deprecated(opts)
	}
}

func (o option) applyOptionalOption(opts *optionalConfig) {
	if o.Input != nil {
		o.Input(&opts.inputConfig)
	}

	if o.Optional != nil {
		o.Optional(opts)
	}
}

func (o option) applyRequiredOption(opts *requiredConfig) {
	if o.Input != nil {
		o.Input(&opts.inputConfig)
	}

	if o.Required != nil {
		o.Required(opts)
	}
}
