package ferrite

import "github.com/dogmatiq/ferrite/variable"

// An Option changes the behavior of an environment variable specification.
//
// Options are used to configure aspects of the variable that are non specific
// to its type.
type Option = variable.RegisterOption //func(*optionsX)

type optionsX struct {
	RegisterOptions []variable.RegisterOption
}

func resolveOptions(options []Option) optionsX {
	var opts optionsX
	opts.RegisterOptions = options
	// for _, opt := range options {
	// 	opt(&opts)
	// }
	return opts
}
