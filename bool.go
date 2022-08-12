package ferrite

import (
	"fmt"

	"github.com/dogmatiq/ferrite/maybe"
	"github.com/dogmatiq/ferrite/variable"
)

// Bool configures an environment variable as a boolean.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Bool(name, desc string) BoolBuilder[bool] {
	return BoolAs[bool](name, desc)
}

// BoolAs configures an environment variable as a boolean using a user-defined
// type.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func BoolAs[T ~bool](name, desc string) BoolBuilder[T] {
	return BoolBuilder[T]{
		spec: variable.PendingSpecFor[T]{
			Name:        variable.Name(name),
			Description: desc,
		},
	}.WithLiterals(
		fmt.Sprint(T(true)),
		fmt.Sprint(T(false)),
	)
}

// BoolBuilder builds a specification for a boolean value.
type BoolBuilder[T ~bool] struct {
	spec variable.PendingSpecFor[T]
}

// WithLiterals overrides the default literals used to represent true and false.
//
// The default literals "true" and "false" are no longer valid values when using
// custom literals.
func (b BoolBuilder[T]) WithLiterals(t, f string) BoolBuilder[T] {
	s, err := variable.NewSet(
		func(v T) (variable.Literal, error) {
			if v {
				return variable.Literal(t), nil
			}
			return variable.Literal(f), nil
		},
		T(true),
		T(false),
	)
	if err != nil {
		b.spec.InvalidErr(err)
	}

	b.spec.Schema = s
	return b
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b BoolBuilder[T]) WithDefault(v T) BoolBuilder[T] {
	b.spec.Default = maybe.Some(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b BoolBuilder[T]) Required(options ...variable.RegisterOption) Required[T] {
	return req[T]{variable.Register(b.spec, options)}
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b BoolBuilder[T]) Optional(options ...variable.RegisterOption) Optional[T] {
	b.spec.IsOptional = true
	return opt[T]{variable.Register(b.spec, options)}
}
