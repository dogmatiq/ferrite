package variable

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/dogmatiq/ferrite/internal/limits"
	"github.com/dogmatiq/ferrite/internal/maybe"
	"github.com/dogmatiq/ferrite/internal/reflectx"
	"golang.org/x/exp/constraints"
)

// Numeric is a schema that allows numeric values.
type Numeric interface {
	Schema

	// Min returns the minimum permitted value as a literal.
	Min() (Literal, bool)

	// Max returns the maximum permitted value as a literal.
	Max() (Literal, bool)

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
func (s TypedNumeric[T]) Min() (Literal, bool) {
	return mustMarshal(s.Marshaler, s.NativeMin).Get()
}

// Max returns the maximum permitted value as a literal.
func (s TypedNumeric[T]) Max() (Literal, bool) {
	return mustMarshal(s.Marshaler, s.NativeMax).Get()
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
	return reflectx.TypeOf[T]()
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
func (s TypedNumeric[T]) Examples(conservative bool) []TypedExample[T] {
	var examples []TypedExample[T]

	min, max := limits.Of[T]()
	if v, ok := s.NativeMin.Get(); ok {
		min = v
		examples = append(
			examples,
			TypedExample[T]{
				Native:      min,
				Description: "the minimum accepted value",
			},
		)
	}

	if v, ok := s.NativeMax.Get(); ok {
		max = v
		examples = append(
			examples,
			TypedExample[T]{
				Native:      max,
				Description: "the maximum accepted value",
			},
		)
	}

	if !conservative {
		// If there are no other examples we do a linear interpolation to find
		// some values within the (min, max) range in an attempt to provide
		// _something_ that might be useful.
		lerp := func(distance float64) T {
			a := float64(min) * (1 - distance)
			b := float64(max) * distance
			return T(a + b)
		}

		examples = append(
			examples,
			TypedExample[T]{
				Native: lerp(0.45),
			},
			TypedExample[T]{
				Native: lerp(0.60),
			},
		)
	}

	return examples
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
	min, hasMin := s.Min()
	max, hasMax := s.Max()

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

// UnwrapNumericParseError returns a more user-friendly error message for
// errors returned by strconv.ParseInt(), ParseUint(), and ParseFloat().
func UnwrapNumericParseError[T constraints.Integer | constraints.Float](
	err error,
	format func(T) string,
) error {
	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		min, max := limits.Of[T]()
		kind := reflectx.KindOf[T]()

		switch numErr.Err {
		case strconv.ErrSyntax:
			return fmt.Errorf("unrecognized %s syntax", kind)
		case strconv.ErrRange:
			if strings.TrimSpace(numErr.Num)[0] == '-' {
				return fmt.Errorf(
					"too low, expected the smallest %s value of %s or greater",
					kind,
					format(min),
				)
			}

			return fmt.Errorf(
				"too high, expected the largest %s value of %s or less",
				kind,
				format(max),
			)
		}
	}

	return err
}
