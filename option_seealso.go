package ferrite

import "github.com/dogmatiq/ferrite/variable"

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
				seeAlso(b.Peek(), v.Spec(), options...)
			}
		},
	}
}

// SeeAlsoOption configures the behavior of the SeeAlso() variable set option.
type SeeAlsoOption interface {
	applyRefersToOption(*variable.RefersTo)
}

func seeAlso(
	subject, refersTo variable.Spec,
	options ...SeeAlsoOption,
) {
	rel := variable.RefersTo{
		Subject:  subject,
		RefersTo: refersTo,
	}

	for _, opt := range options {
		opt.applyRefersToOption(&rel)
	}

	if err := variable.AddRelationship(rel); err != nil {
		panic(err.Error())
	}
}
