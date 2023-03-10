package ferrite

import (
	"github.com/dogmatiq/ferrite/variable"
)

// Deprecated is the application-facing interface for a value that is sourced
// from deprecated environment variables.
//
// It is obtained by calling Deprecated() on a variable builder.
type Deprecated[T any] interface {
}

// DeprecatedOption is an option that can be applied to a deprecated variable.
type DeprecatedOption interface {
	variable.RegisterOption
}

// deprecated is a convenience function that registers and returns a
// deprecated[T] that maps one-to-one with an environment variable of the same
// type.
func deprecated[T any, S variable.TypedSchema[T]](
	schema S,
	spec *variable.TypedSpecBuilder[T],
	options []DeprecatedOption,
) Deprecated[T] {
	spec.MarkDeprecated()

	variable.Register(
		spec.Done(schema),
		options...,
	)

	// interface is currently empty so we don't need an implementation
	return nil
}
