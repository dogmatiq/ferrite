package ferrite

import (
	"fmt"

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
		spec: variable.NewSpec[T](name, desc),
	}.WithLiterals(
		variable.Literal(fmt.Sprint(T(true))),
		variable.Literal(fmt.Sprint(T(false))),
	)
}

// BoolBuilder builds a specification for a boolean value.
type BoolBuilder[T ~bool] struct {
	spec variable.SpecFor[T]
}

// WithLiterals overrides the default literals used to represent true and false.
//
// The default literals "true" and "false" are no longer valid values when using
// custom literals.
func (b BoolBuilder[T]) WithLiterals(t, f variable.Literal) BoolBuilder[T] {
	set, err := variable.NewSet(
		func(v T) variable.Literal {
			if v {
				return t
			}
			return f
		},
		true,
		false,
	)
	if err != nil {
		b.spec.InvalidErr(err)
	}

	b.spec.SetClass(set)
	return b
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b BoolBuilder[T]) WithDefault(v T) BoolBuilder[T] {
	b.spec.SetDefault(v)
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
	b.spec.MarkOptional()
	return opt[T]{variable.Register(b.spec, options)}
}
