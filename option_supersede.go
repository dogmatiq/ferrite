package ferrite

import "github.com/dogmatiq/ferrite/variable"

// SupersededBy is a option for a deprecated variable set that indicates the
// variables in another set, s, should be used instead.
func SupersededBy(s VariableSet, options ...SupersededByOption) DeprecatedOption {
	return option{
		ApplyToSpecInDeprecatedSet: func(b variable.SpecBuilder) {
			for _, v := range s.variables() {
				rel := variable.Supersedes{
					Subject:    v.Spec(),
					Supersedes: b.Peek(),
				}

				for _, opt := range options {
					opt.applySupersedesOption(&rel)
				}

				if err := variable.AddRelationship(rel); err != nil {
					panic(err.Error())
				}
			}
		},
	}
}

// SupersededByOption changes the behavior of the SupersededBy() option.
type SupersededByOption interface {
	applySupersedesOption(*variable.Supersedes)
}
