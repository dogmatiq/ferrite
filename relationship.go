package ferrite

import "github.com/dogmatiq/ferrite/variable"

// SeeAlsoOption changes the behavior of a call the SeeAlso() method of the
// various builder implementations.
type SeeAlsoOption interface {
	applySeeAlsoOption(*variable.SeeAlso)
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
	r := variable.SeeAlso{
		From: from,
		To:   to,
	}

	for _, opt := range options {
		opt.applySeeAlsoOption(&r)
	}

	from.AddRelationship(r)
	to.AddRelationship(r)
}
