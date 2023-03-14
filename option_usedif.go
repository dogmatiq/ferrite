package ferrite

import (
	"reflect"

	"github.com/dogmatiq/ferrite/variable"
)

// RelevantIf is an option that "enables" a variable set only if the value
// obtained from another set, s, is non-zero.
//
// An "irrelevant" or "disabled" variable set behaves as though its environment
// variables are undefined, irrespective of the actual values of the variables
// and any default values.
func RelevantIf(s VariableSet, options ...RelevantIfOption) interface {
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
					opt.applyRelevantIfOption(&rel)
				}

				if err := variable.AddRelationship(rel); err != nil {
					panic(err.Error())
				}

				b.Precondition(
					func() bool {
						return !reflect.ValueOf(
							s.value(),
						).IsZero()
					},
				)
			}
		},
	}
}

// RelevantIfOption changes the behavior of the RelevantIf() option.
type RelevantIfOption interface {
	applyRelevantIfOption(*variable.DependsOn)
}
