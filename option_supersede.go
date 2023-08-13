package ferrite

import "github.com/dogmatiq/ferrite/internal/variable"

// SupersededBy is a option for a deprecated variable set that indicates the
// variables in another set, s, should be used instead.
func SupersededBy(s VariableSet, _ ...SupersededByOption) DeprecatedOption {
	return option{
		ApplyToSpecInDeprecatedSet: func(b variable.SpecBuilder) {
			for _, v := range s.variables() {
				variable.EstablishRelationships(
					variable.Supersedes{
						Subject:    v.Spec(),
						Supersedes: b.Peek(),
					},
				)
			}
		},
	}
}

// SupersededByOption changes the behavior of the SupersededBy() option.
type SupersededByOption interface {
	future()
}
