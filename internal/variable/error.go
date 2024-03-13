package variable

import "fmt"

// Error is an error that indicates a problem parsing or validating an
// environment variable.
type Error interface {
	error

	// Name returns the name of the environment variable.
	Name() string
}

// MinLengthError indicates that a value was shorter than the minimum permitted
// length.
type MinLengthError struct {
	ViolatedSchema LengthLimited
}

var _ SchemaError = MinLengthError{}

// Schema returns the schema that was violated.
func (e MinLengthError) Schema() Schema {
	return e.ViolatedSchema
}

// AcceptVisitor passes the error to the appropriate method of v.
func (e MinLengthError) AcceptVisitor(v SchemaErrorVisitor) {
	v.VisitMinLengthError(e)
}

func (e MinLengthError) Error() string {
	return fmt.Sprintf("too short, %s", e.ViolatedSchema.ExplainLengthError())
}

// MaxLengthError indicates that a value was greater than the maximum permitted
// length.
type MaxLengthError struct {
	ViolatedSchema LengthLimited
}

var _ SchemaError = MaxLengthError{}

// Schema returns the schema that was violated.
func (e MaxLengthError) Schema() Schema {
	return e.ViolatedSchema
}

// AcceptVisitor passes the error to the appropriate method of v.
func (e MaxLengthError) AcceptVisitor(v SchemaErrorVisitor) {
	v.VisitMaxLengthError(e)
}

func (e MaxLengthError) Error() string {
	return fmt.Sprintf("too long, %s", e.ViolatedSchema.ExplainLengthError())
}

func explainLengthError(s LengthLimited, lengthDesc string) string {
	min, hasMin := s.MinLength()
	max, hasMax := s.MaxLength()

	if !hasMin {
		return fmt.Sprintf("expected %s to be %d bytes or fewer", lengthDesc, max)
	}

	if !hasMax {
		return fmt.Sprintf("expected %s to be %d bytes or more", lengthDesc, min)
	}

	if min == max {
		return fmt.Sprintf("expected %s to be exactly %d bytes", lengthDesc, min)
	}

	return fmt.Sprintf("expected %s to be between %d and %d bytes", lengthDesc, min, max)
}
