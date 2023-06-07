package ferrite

import "github.com/dogmatiq/ferrite/internal/variable"

// SeeAlso is an option for a variable set that adds the variables in another
// set, s, to the "see also" section of the generated documentation.
func SeeAlso(s VariableSet, options ...SeeAlsoOption) interface {
	RequiredOption
	OptionalOption
	DeprecatedOption
} {
	return option{
		ApplyToSpec: func(b variable.SpecBuilder) {
			for _, v := range s.variables() {
				variable.EstablishRelationships(
					variable.RefersTo{
						Subject:  b.Peek(),
						RefersTo: v.Spec(),
					},
				)
			}
		},
	}
}

// SeeAlsoOption configures the behavior of the SeeAlso() variable set option.
type SeeAlsoOption interface {
	future()
}
