package variable

// Example is an example value.
type Example struct {
	Canonical   Literal
	Description string
	IsNormative bool
}

// TypedExample is an example value depicted by type T.
type TypedExample[T any] struct {
	Native      T
	Description string
	IsNormative bool
}

// WithNormativeExample is a SpecOption that adds a "normative" example value.
//
// A normative example is one that is meaningful in the context of the
// variable's use.
func WithNormativeExample[T any](v T, desc string) SpecOption[T] {
	return func(opts *specOptions[T]) error {
		eg := TypedExample[T]{v, desc, true}
		opts.Examples = append(opts.Examples, eg)
		return nil
	}
}

// WithNonNormativeExample is a SpecOption that adds a "non-normative" example
// value.
//
// A non-normative example is one that may not be meaningful in the context of
// the variable's use, but is included for illustrative purposes.
//
// For example, a variable that represents a URL may have a non-normative
// example of "https://example.org/path", even if the actual use-case for the
// variable requires an "ftp" URL.
func WithNonNormativeExample[T any](v T, desc string) SpecOption[T] {
	return func(opts *specOptions[T]) error {
		eg := TypedExample[T]{v, desc, false}
		opts.Examples = append(opts.Examples, eg)
		return nil
	}
}

// appendExample appends eg to the examples only if there is no existing example
// with the same value.
func appendExample(examples []Example, eg Example) []Example {
	if containsExample(examples, eg.Canonical) {
		return examples
	}

	return append(examples, eg)
}

// prependExample prepends eg to the examples only if there is no existing
// example with the same value.
func prependExample(examples []Example, eg Example) []Example {
	if containsExample(examples, eg.Canonical) {
		return examples
	}

	return append([]Example{eg}, examples...)
}

// containsExample returns true if lit is one of the example values.
func containsExample(examples []Example, lit Literal) bool {
	for _, eg := range examples {
		if eg.Canonical == lit {
			return true
		}
	}

	return false
}
