package variable

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/dogmatiq/ferrite/internal/limits"
	"github.com/dogmatiq/ferrite/maybe"
	"golang.org/x/exp/constraints"
)

// Numeric is a schema that allows numeric values.
type Numeric interface {
	Schema

	// Min returns the minimum permitted value as a literal.
	Min() maybe.Value[Literal]

	// Max returns the maximum permitted value as a literal.
	Max() maybe.Value[Literal]

	// Limits returns the range of permitted values.
	//
	// explicit is true if both the minimum and the maximum limits are specified
	// by the application, or false if either of the limits is simply the limit
	// of the underlying type.
	Limits() (min, max Literal, explicit bool)

	// Bits is the number of bits used to store the number.
	Bits() int
}

// TypedNumeric is a numeric value depicted by type T.
type TypedNumeric[T constraints.Integer | constraints.Float] struct {
	Marshaler            Marshaler[T]
	NativeMin, NativeMax maybe.Value[T]
}

// Min returns the minimum permitted value as a literal.
func (s TypedNumeric[T]) Min() maybe.Value[Literal] {
	return mustMarshal(s.Marshaler, s.NativeMin)
}

// Max returns the maximum permitted value as a literal.
func (s TypedNumeric[T]) Max() maybe.Value[Literal] {
	return mustMarshal(s.Marshaler, s.NativeMax)
}

// Limits returns the range of permitted values.
//
// explicit is true if both the minimum and the maximum limits are specified
// by the application, or false if either of the limits is simply the limit
// of the underlying type.
func (s TypedNumeric[T]) Limits() (min, max Literal, explicit bool) {
	lower, upper := limits.Of[T]()
	explicit = true

	if v, ok := s.NativeMin.Get(); ok {
		lower = v
	} else {
		explicit = false
	}

	if v, ok := s.NativeMax.Get(); ok {
		lower = v
	} else {
		explicit = false
	}

	min, err := s.Marshaler.Marshal(lower)
	if err != nil {
		panic(err)
	}

	max, err = s.Marshaler.Marshal(upper)
	if err != nil {
		panic(err)
	}

	return min, max, explicit
}

// Bits is the number of bits used to store the number.
func (s TypedNumeric[T]) Bits() int {
	return int(unsafe.Sizeof(T(0))) * 8
}

// Type returns the type of the native value.
func (s TypedNumeric[T]) Type() reflect.Type {
	return typeOf[T]()
}

// Finalize prepares the schema for use.
//
// It returns an error if schema is invalid.
func (s TypedNumeric[T]) Finalize() error {
	if _, err := marshal(s.Marshaler, s.NativeMin); err != nil {
		return fmt.Errorf("minimum value: %w", err)
	}

	if _, err := marshal(s.Marshaler, s.NativeMax); err != nil {
		return fmt.Errorf("maximum value: %w", err)
	}

	return nil
}

// AcceptVisitor passes s to the appropriate method of v.
func (s TypedNumeric[T]) AcceptVisitor(v SchemaVisitor) {
	v.VisitNumeric(s)
}

// Marshal converts a value to its literal representation.
func (s TypedNumeric[T]) Marshal(v T) (Literal, error) {
	if err := s.validate(v); err != nil {
		return Literal{}, err
	}

	return s.Marshaler.Marshal(v)
}

// Unmarshal converts a literal value to it's native representation.
func (s TypedNumeric[T]) Unmarshal(v Literal) (T, error) {
	n, err := s.Marshaler.Unmarshal(v)
	if err != nil {
		return 0, err
	}

	return n, s.validate(n)
}

// Examples returns a (possibly empty) set of examples of valid values.
func (s TypedNumeric[T]) Examples(hasOtherExamples bool) []TypedExample[T] {
	if hasOtherExamples {
		return nil
	}

	// If there are no other examples we do a linear interpolation at 25% of the
	// (min, max) range in an attempt to provide _something_ that might be
	// useful.
	min, max := limits.Of[T]()

	return []TypedExample[T]{
		{
			Native:      T(float64(min) + float64(max-min)*0.25),
			Description: "randomly generated example",
		},
	}
}

// validate returns an error if v is invalid.
func (s TypedNumeric[T]) validate(v T) error {
	if min, ok := s.NativeMin.Get(); ok && v < min {
		return MinError{s}
	}

	if max, ok := s.NativeMax.Get(); ok && v > max {
		return MaxError{s}
	}

	return nil
}

// MinError indicates that a numeric value was less than the minimum permitted
// value.
type MinError struct {
	Numeric Numeric
}

var _ SchemaError = MinError{}

// Schema returns the schema that was violated.
func (e MinError) Schema() Schema {
	return e.Numeric
}

// AcceptVisitor passes the error to the appropriate method of v.
func (e MinError) AcceptVisitor(v SchemaErrorVisitor) {
	v.VisitMinError(e)
}

func (e MinError) Error() string {
	return fmt.Sprintf("too low, %s", explainRangeError(e.Numeric))
}

// MaxError indicates that a numeric value was greater than the maximum
// permitted value.
type MaxError struct {
	Numeric Numeric
}

var _ SchemaError = MaxError{}

// Schema returns the schema that was violated.
func (e MaxError) Schema() Schema {
	return e.Numeric
}

// AcceptVisitor passes the error to the appropriate method of v.
func (e MaxError) AcceptVisitor(v SchemaErrorVisitor) {
	v.VisitMaxError(e)
}

func (e MaxError) Error() string {
	return fmt.Sprintf("too high, %s", explainRangeError(e.Numeric))
}

func explainRangeError(s Numeric) string {
	min, hasMin := s.Min().Get()
	max, hasMax := s.Max().Get()

	if !hasMin {
		return fmt.Sprintf(
			"expected %s or less",
			max.Quote(),
		)
	}

	if !hasMax {
		return fmt.Sprintf(
			"expected %s or greater",
			min.Quote(),
		)
	}

	if min == max {
		return fmt.Sprintf(
			"expected exactly %s",
			min.Quote(),
		)
	}

	return fmt.Sprintf(
		"expected between %s and %s",
		min.Quote(),
		max.Quote(),
	)
}
