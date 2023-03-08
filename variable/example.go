package variable

// ExampleSource is an enumeration describing the source of an example.
type ExampleSource int

const (
	// ExampleSourceUnknown indicates the source of an example is unknown.
	ExampleSourceUnknown ExampleSource = iota

	// ExampleSourceSchema indicates the example was taken from the schema, such
	// as a minimum or maximum value.
	ExampleSourceSchema

	// ExampleSourceSpecBuilder indicates that the examples was supplied to the
	// specification via a SpecBuilder.
	ExampleSourceSpecBuilder

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

// BestExample returns the heuristically "best" example to use for the given
// spec.
func BestExample(spec Spec) Example {
	var incumbent Example

	isBetter := func(a, b Example) bool {
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
	}

	for _, candidate := range spec.Examples() {
		// Always prefer the default value, as it's typically supplied by the
		// application author and "guaranteed" to work.
		if IsDefault(spec, candidate.Canonical) {
			return candidate
		}

		// Otherwise, keep going looking for a better candidate.
		if isBetter(candidate, incumbent) {
			incumbent = candidate
		}
	}

	return incumbent
}
