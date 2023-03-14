package variable

import (
	"fmt"
	"sync"
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
	spec *TypedSpec[T]
	env  Environment

	once         sync.Once
	availability Availability
	source       Source
	value        valueOf[T]
	err          Error
}

// Spec returns the variable's specification.
func (v *OfType[T]) Spec() Spec {
	return v.spec
}

// Availability returns the variable's availability.
func (v *OfType[T]) Availability() Availability {
	v.resolve()
	return v.availability
}

// Source returns the source of the variable's value.
func (v *OfType[T]) Source() Source {
	v.resolve()
	return v.source
}

// Value returns the variable's value.
//
// If no value is available it returns a zero-value. It is the caller's
// responsibility to check the variable's availability before using the value.
func (v *OfType[T]) Value() Value {
	v.resolve()
	return v.value
}

// NativeValue returns the variable's value.
//
// If no value is available it returns a zero-value. It is the caller's
// responsibility to check the variable's availability before using the value.
func (v *OfType[T]) NativeValue() T {
	v.resolve()
	return v.value.native
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
	v.resolve()
	return v.err
}

func (v *OfType[T]) resolve() {
	v.once.Do(func() {
		// Override the availability to AvailabilityIgnored if any of the
		// preconditions fail.
		defer func() {
			for _, fn := range v.spec.preconditions {
				if !fn() {
					v.availability = AvailabilityIgnored
					break
				}
			}
		}()

		lit := v.env.Get(v.spec.name)

		if lit.String == "" {
			if def, ok := v.spec.def.Get(); ok {
				v.availability = AvailabilityOK
				v.source = SourceDefault
				v.value = def
			} else if v.spec.required {
				v.availability = AvailabilityNone
				v.err = undefinedError{v.spec.Name()}
			}
			return
		}

		v.source = SourceEnvironment

		n, c, err := v.spec.Unmarshal(lit)
		if err != nil {
			v.availability = AvailabilityInvalid
			v.err = valueError{
				name:    v.spec.name,
				literal: lit,
				cause:   err,
			}
			return
		}

		v.availability = AvailabilityOK
		v.value = valueOf[T]{
			verbatim:  lit,
			native:    n,
			canonical: c,
		}
	})
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
