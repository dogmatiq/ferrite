package variable

// Example is an example value.
type Example struct {
	Canonical   Literal
	Description string
}

// TypedExample is an example value depicted by type T.
type TypedExample[T any] struct {
	Native      T
	Description string
}

// WithExample is a SpecOption that adds an example value.
func WithExample[T any](v T, desc string) SpecOption[T] {
	return func(opts *specOptions[T]) error {
		eg := TypedExample[T]{v, desc}
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

// containsExample returns true if lit is one of the example values.
func containsExample(examples []Example, lit Literal) bool {
	for _, eg := range examples {
		if eg.Canonical == lit {
			return true
		}
	}

	return false
}
