package variable

import (
	"golang.org/x/exp/slices"
)

// ExampleSource is an enumeration describing the source of an example.
type ExampleSource int

const (
	// ExampleSourceUnknown indicates the source of an example is unknown.
	ExampleSourceUnknown ExampleSource = iota

	// ExampleSourceSchema indicates the example was taken from the schema, such
	// as a minimum or maximum value.
	ExampleSourceSchema

	// ExampleSourceSpecOption indicates that the examples was supplied to the
	// specification via a SpecOption.
	ExampleSourceSpecOption

	// ExampleSourceSpecDefault indicates that the example was generated from
	// the default value of the variable.
	ExampleSourceSpecDefault
)

// Example is an example value.
type Example struct {
	Canonical   Literal
	Description string
	IsNormative bool
	Source      ExampleSource
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

// BestExample returns the heuristically "best" example to use for the given
// spec.
func BestExample(spec Spec) Example {
	// Always prefer the default value, as it's always supplied by the
	// application author and "guaranteed" to work.
	if def, ok := spec.Default(); ok {
		for _, eg := range spec.Examples() {
			if eg.Canonical == def {
				return eg
			}
		}
	}

	// Otherwise, we rank the examples by their "quality" and return the
	// highest-ranked one.
	ranked := slices.Clone(
		spec.Examples(),
	)

	slices.SortFunc(
		ranked,
		func(a, b Example) bool {
			// Prefer normative examples over non-normative ones.
			if a.IsNormative != b.IsNormative {
				return a.IsNormative
			}

			// Prefer examples from higher-levels sources.
			if a.Source != b.Source {
				return a.Source > b.Source
			}

			// Prefer examples with (longer) descriptions than those without.
			if a.Description != b.Description {
				return len(a.Description) > len(b.Description)
			}

			// If both are schema-generated non-normative examples, assume that
			// the shorter value is less likely to be "weird".
			if a.Source == ExampleSourceSchema && !a.IsNormative {
				return len(a.Canonical.String) < len(b.Canonical.String)
			}

			// Otherwise, assume that the the longer value has better
			// illustrative properties.
			return len(a.Canonical.String) > len(b.Canonical.String)
		},
	)

	return ranked[0]
}

// buildExamples returns the complete set of examples to use for the given spec.
//
// It combines the examples provided by the user with the examples provided by
// the variable's schema.
func buildExamples[T any](
	spec TypedSpec[T],
	fromOptions []TypedExample[T],
) ([]Example, error) {
	var examples []Example

	// Add the examples provided as options to the spec.
	for _, eg := range fromOptions {
		lit, err := spec.Marshal(eg.Native)
		if err != nil {
			return nil, err
		}

		examples = appendUniqueExample(examples, Example{
			Canonical:   lit,
			Description: eg.Description,
			IsNormative: eg.IsNormative,
			Source:      ExampleSourceSpecOption,
		})
	}

	hasOtherExamples := len(fromOptions) > 0 || !spec.def.IsEmpty()

	// Generate examples from the schema and add each one only if it meets all
	// of the constraints.
	for _, eg := range spec.schema.Examples(hasOtherExamples) {
		if lit, err := spec.Marshal(eg.Native); err == nil {
			examples = appendUniqueExample(examples, Example{
				Canonical:   lit,
				Description: eg.Description,
				IsNormative: eg.IsNormative,
				Source:      ExampleSourceSchema,
			})
		}
	}

	// Add an example of the default value (if one is not already present).
	if def, ok := spec.def.Get(); ok {
		examples = prependUniqueExample(examples, Example{
			Canonical:   def.Canonical(),
			IsNormative: true,
			Source:      ExampleSourceSpecDefault,
		})
	}

	if len(examples) == 0 {
		panic("spec must contain at least one example")
	}

	return examples, nil
}

// appendUniqueExample appends eg to the examples only if there is no existing example
// with the same value.
func appendUniqueExample(examples []Example, eg Example) []Example {
	if containsExample(examples, eg.Canonical) {
		return examples
	}

	return append(examples, eg)
}

// prependUniqueExample prepends eg to the examples only if there is no existing
// example with the same value.
func prependUniqueExample(examples []Example, eg Example) []Example {
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
