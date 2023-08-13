package ferrite

import (
	"reflect"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// RelevantIf is an option that enables a variable set only if the value
// obtained from another set, s, is "truthy" (not the zero-value).
//
// A "irrelevant" variable set behaves as though its environment variables are
// undefined, irrespective of the actual values of the variables and any default
// values.
func RelevantIf(s VariableSet, _ ...RelevantIfOption) interface {
	RequiredOption
	OptionalOption
	DeprecatedOption
} {
	return option{
		ApplyToSpec: func(b variable.SpecBuilder) {
			for _, v := range s.variables() {
				variable.EstablishRelationships(
					variable.RefersTo{
						Subject:  b.Peek(),
						RefersTo: v.Spec(),
					},
					variable.DependsOn{
						Subject:   b.Peek(),
						DependsOn: v.Spec(),
					},
				)

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
	future()
}
