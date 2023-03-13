package ferrite

import "github.com/dogmatiq/ferrite/variable"

// SupersededBy is a deprecation option that indicates i is a direct
// replacement for the deprecated variable(s).
func SupersededBy(i Input, options ...SupersededByOption) DeprecatedOption {
	return option{
		Deprecated: func(cfg *deprecatedConfig) {
			for _, v := range i.variables() {
				rel := variable.Supersedes{
					Subject:    v.Spec(),
					Supersedes: cfg.Spec.Peek(),
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
