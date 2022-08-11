package spec2

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Range limits a value to a specific range.
type Range interface {
	Constraint

	// Limits returns the minimum and maximum allowed values.
	Limits() (min, max Maybe[RValue])
}

type rangeOf[T constraints.Ordered] struct {
	Min       Maybe[T]
	Max       Maybe[T]
	Marshaler Marshaler[T]
}

func (c rangeOf[T]) Limits() (min, max Maybe[RValue]) {
	min = Map(c.Min, c.Marshaler.Marshal)
	max = Map(c.Max, c.Marshaler.Marshal)
	return min, max
}

func (c rangeOf[T]) Test(r RValue) error {
	n, err := c.Marshaler.Unmarshal(r)
	if err != nil {
		return err
	}

	if min, ok := c.Min.Get(); ok {
		if n.Native < min {
			return RangeError{r, c}
		}
	}

	if max, ok := c.Max.Get(); ok {
		if n.Native > max {
			return RangeError{r, c}
		}
	}

	return nil
}

func (c rangeOf[T]) AcceptVisitor(v ConstraintVisitor) {
	v.VisitRange(c)
}

// RangeError indicates that a Range constraint was violated.
type RangeError struct {
	Value RValue
	Range Range
}

func (e RangeError) Error() string {
	min, max := e.Range.Limits()

	if r, ok := min.Get(); ok {
		return fmt.Sprintf("must be %s or higher", r)
	}

	if r, ok := max.Get(); ok {
		return fmt.Sprintf("must be %s or lower", r)
	}

	return fmt.Sprintf(
		"must be between %s and %s",
		min.MustGet(),
		max.MustGet(),
	)
}
