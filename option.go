package ferrite

import "github.com/dogmatiq/ferrite/variable"

// An Option changes the behavior of a variable. Options are applicable to all
// variable types.
type Option = variable.RegisterOption

// (*options)

// type options struct {
// 	SpecOptions     []variable.TypedSpecOption
// 	RegisterOptions []variable.RegisterOption
// }

// // WithRegistry is an option that sets the registry that an environment variable
// // specification is placed into.
// func WithRegistry(r *variable.Registry) Option {
// 	return func(o *options) {
// 		o.RegisterOptions = append(
// 			o.RegisterOptions,
// 			variable.WithRegistry(r),
// 		)
// 	}
// }
