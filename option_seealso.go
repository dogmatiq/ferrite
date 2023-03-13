package ferrite

import "github.com/dogmatiq/ferrite/variable"

// SeeAlso is an option that the variables in s to the "see also" section of the
// generated documentation.
func SeeAlso(s VariableSet, options ...SeeAlsoOption) interface {
	DeprecatedOption
	RequiredOption
	OptionalOption
} {
	return option{
		ApplyToSpec: func(spec variable.SpecBuilder) {
			for _, v := range s.variables() {
				seeAlso(spec.Peek(), v.Spec(), options...)
			}
		},
	}
}

// SeeAlsoOption changes the behavior of the SeeAlso() option.
type SeeAlsoOption interface {
	applyRefersToOption(*variable.RefersTo)
}

func seeAlso(
	from, to variable.Spec,
	options ...SeeAlsoOption,
) {
	rel := variable.RefersTo{
		Subject:  from,
		RefersTo: to,
	}

	for _, opt := range options {
		opt.applyRefersToOption(&rel)
	}

	if err := variable.AddRelationship(rel); err != nil {
		panic(err.Error())
	}
}
