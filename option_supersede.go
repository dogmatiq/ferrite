package ferrite

import "github.com/dogmatiq/ferrite/variable"

// SupersededBy is a deprecation option that indicates i is a direct
// replacement for the deprecated variable(s).
func SupersededBy(i Input, options ...SupersededByOption) DeprecatedOption {
	return option{
		Deprecated: func(cfg *deprecatedConfig) {
			for _, v := range i.variables() {
				rel := variable.SupersededBy{
					Spec:         cfg.Spec.Peek(),
					SupersededBy: v.Spec(),
				}

				for _, opt := range options {
					opt.applySupersededByOption(&rel)
				}

				variable.ApplyRelationship(rel)
			}
		},
	}
}

// SupersededByOption changes the behavior of the SupersededBy() option.
type SupersededByOption interface {
	applySupersededByOption(*variable.SupersededBy)
}
