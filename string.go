package ferrite

import (
	"github.com/dogmatiq/ferrite/maybe"
	"github.com/dogmatiq/ferrite/variable"
)

// String configures an environment variable as a string.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func String(name, desc string) StringBuilder[string] {
	return StringAs[string](name, desc)
}

// StringAs configures an environment variable as a string using a user-defined
// type.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func StringAs[T ~string](name, desc string) StringBuilder[T] {
	return StringBuilder[T]{
		name: name,
		desc: desc,
	}
}

// StringBuilder builds a specification for a string variable.
type StringBuilder[T ~string] struct {
	name, desc string
	def        maybe.Value[T]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b StringBuilder[T]) WithDefault(v T) StringBuilder[T] {
	b.def = maybe.Some(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b StringBuilder[T]) Required(options ...variable.RegisterOption) Required[T] {
	return registerRequired(b.spec(true), options)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b StringBuilder[T]) Optional(options ...variable.RegisterOption) Optional[T] {
	return registerOptional(b.spec(false), options)
}

func (b StringBuilder[T]) spec(req bool) variable.TypedSpec[T] {
	s, err := variable.NewSpec(
		b.name,
		b.desc,
		b.def,
		req,
		variable.TypedString[T]{},
	)
	if err != nil {
		panic(err.Error())
	}

	return s
}
