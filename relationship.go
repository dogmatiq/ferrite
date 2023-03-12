package ferrite

import "github.com/dogmatiq/ferrite/variable"

// SeeAlsoOption changes the behavior of a call the SeeAlso() method of the
// various builder implementations.
type SeeAlsoOption interface {
	applySeeAlsoOption(*variable.RefersTo)
}

func seeAlsoInput(
	from variable.SpecBuilder,
	to Input,
	options ...SeeAlsoOption,
) {
	for _, v := range to.variables() {
		seeAlso(from.Peek(), v.Spec(), options...)
	}
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
