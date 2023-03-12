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

// SeeAlsoOption changes the behavior of a call the SeeAlso() method of the
// various builder implementations.
type SeeAlsoOption interface {
	applySeeAlsoOption(*variable.RefersTo)
}

func seeAlso(
	from, to variable.Spec,
	options ...SeeAlsoOption,
) {
	rel := variable.RefersTo{
		Spec:     from,
		RefersTo: to,
	}

	for _, opt := range options {
		opt.applySeeAlsoOption(&rel)
	}

	variable.ApplyRelationship(rel)
}
