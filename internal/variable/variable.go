package variable

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/dogmatiq/ferrite/internal/environment"
)

// Availability is an enumeration describing why a variable is or is not
// available.
type Availability int

const (
	// AvailabilityNone indicates that the variable's value is not available
	// because it has not been explicitly defined in the environment and no
	// default value has been specified.
	AvailabilityNone Availability = iota

	// AvailabilityInvalid indicates that the variable is defined with an
	// invalid value.
	AvailabilityInvalid

	// AvailabilityIgnored indicates that the variable's value valid but it
	// should not be made available to the user.
	AvailabilityIgnored

	// AvailabilityOK indicates that the variable's value is valid and available
	// to the user.
	AvailabilityOK
)

// Source is an enumeration of the possible sources of an environment variable's
// value.
type Source int

const (
	// SourceNone indicates that there is no value, and hence no source.
	SourceNone Source = iota

	// SourceDefault indicates that the value is the default value in the
	// variable's specification.
	SourceDefault

	// SourceEnvironment indicates that the value was obtained from the
	// environment.
	SourceEnvironment
)

// Any is an interface for an environment variable of any type.
type Any interface {
	Spec() Spec
	Availability() Availability
	Source() Source
	Value() Value
	Error() Error
}

// OfType is an environment variable depicted by type T.
type OfType[T any] struct {
	TypedSpec *TypedSpec[T]

	m          sync.Mutex
	resolution atomic.Pointer[resolution[T]]
}

// resolution holds the cached result of resolving an environment variable.
type resolution[T any] struct {
	lit    string
	source Source
	value  valueOf[T]
	err    Error
}

// Spec returns the variable's specification.
func (v *OfType[T]) Spec() Spec {
	return v.TypedSpec
}

// Availability returns the variable's availability.
func (v *OfType[T]) Availability() Availability {
	for _, fn := range v.TypedSpec.preconditions {
		if !fn() {
			return AvailabilityIgnored
		}
	}

	r := v.resolve()

	if _, ok := r.err.(valueError); ok {
		return AvailabilityInvalid
	}

	if r.source == SourceNone {
		return AvailabilityNone
	}

	return AvailabilityOK
}

// Source returns the source of the variable's value.
func (v *OfType[T]) Source() Source {
	return v.resolve().source
}

// Value returns the variable's value.
//
// If no value is available it returns a zero-value. It is the caller's
// responsibility to check the variable's availability before using the value.
func (v *OfType[T]) Value() Value {
	return v.resolve().value
}

// NativeValue returns the variable's value.
//
// If no value is available it returns a zero-value. It is the caller's
// responsibility to check the variable's availability before using the value.
func (v *OfType[T]) NativeValue() T {
	return v.resolve().value.native
}

// Error returns an error describing the variable's state.
//
// If the variable's availability is AvailabilityInvalid the error is guaranteed
// to be a ValueError.
//
// The error is nil if the variable is in a valid state, which occurs when it
// has an availability of AvailabilityOK, or if it has an availability of
// AvailabilityNone and v.Spec().IsRequired() is false.
func (v *OfType[T]) Error() Error {
	return v.resolve().err
}

func (v *OfType[T]) resolve() *resolution[T] {
	lit := environment.Get(v.TypedSpec.name)

	if r := v.resolution.Load(); r != nil {
		if r.lit == lit {
			return r
		}
	}

	v.m.Lock()
	defer v.m.Unlock()

	if r := v.resolution.Load(); r != nil {
		if r.lit == lit {
			return r
		}
	}

	r := &resolution[T]{
		lit: lit,
	}

	if lit == "" {
		if def, ok := v.TypedSpec.def.Get(); ok {
			r.source = SourceDefault
			r.value = def
		} else if v.TypedSpec.required {
			r.err = undefinedError{v.TypedSpec.Name()}
		}
	} else {
		r.source = SourceEnvironment

		n, c, err := v.TypedSpec.Unmarshal(ConstraintContextFinal, Literal{String: lit})
		if err != nil {
			r.err = valueError{
				name:    v.TypedSpec.name,
				literal: Literal{String: lit},
				cause:   err,
			}
		} else {
			r.value = valueOf[T]{
				verbatim:  Literal{String: lit},
				native:    n,
				canonical: c,
			}
		}
	}

	v.resolution.Store(r)
	return r
}

// undefinedError is an Error that indicates that a variable is undefined and
// does not have a default value.
type undefinedError struct {
	name string
}

func (e undefinedError) Name() string {
	return e.name
}

func (e undefinedError) Error() string {
	return fmt.Sprintf(
		"%s is undefined and does not have a default value",
		e.name,
	)
}
