package ferrite

import "github.com/dogmatiq/ferrite/variable"

// SeeAlso is an option that adds i to the "see also" section of the generated
// documentation.
func SeeAlso(i Input, options ...SeeAlsoOption) interface {
	DeprecatedOption
	RequiredOption
	OptionalOption
} {
	return option{
		Input: func(cfg *inputConfig) {
			for _, v := range i.variables() {
				seeAlso(cfg.Spec.Peek(), v.Spec(), options...)
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
		Sub:      from,
		RefersTo: to,
	}

	for _, opt := range options {
		opt.applyRefersToOption(&rel)
	}

	if err := variable.ApplyRelationship(rel); err != nil {
		panic(err.Error())
	}
}
