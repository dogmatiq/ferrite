package ferrite

import (
	"github.com/dogmatiq/ferrite/is"
	"github.com/dogmatiq/ferrite/variable"
)

// IgnoreIf is an option that causes the variable set to be ignored if the
// value obtained from another set, s, matched a predicate.
//
// Ignored variables are treated as though they are undefined, regardless of the
// actual value of the environment variables or the default values.
func IgnoreIf[Set TypedVariableSet[T], T any](
	s Set,
	fn is.Predicate[T],
	options ...IgnoreIfOption[T],
) interface {
	RequiredOption
	OptionalOption
	DeprecatedOption
} {
	return option{
		ApplyToSpec: func(b variable.SpecBuilder) {
			for _, v := range s.variables() {
				rel := variable.DependsOn{
					Subject:   b.Peek(),
					DependsOn: v.Spec(),
				}

				for _, opt := range options {
					opt.applyDependsOnOption(&rel)
				}

				if err := variable.AddRelationship(rel); err != nil {
					panic(err.Error())
				}

				b.Precondition(func() bool {
					return !fn(s.value())
				})
			}
		},
	}
}

// IgnoreIfOption changes the behavior of the IgnoreIf() and IgnoreUnless option.
type IgnoreIfOption[T any] interface {
	applyDependsOnOption(*variable.DependsOn)
}
