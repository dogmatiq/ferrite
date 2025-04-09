package ferrite

import (
	"fmt"
	"reflect"

	"github.com/dogmatiq/ferrite/internal/maybe"
	"github.com/dogmatiq/ferrite/internal/variable"
)

// RelevantIf is an option that enables a variable set only if the value
// obtained from another set, s, is "truthy" (not the zero-value).
//
// An "irrelevant" variable set behaves as though its environment variables are
// undefined, irrespective of the actual values of the variables and any default
// values.
func RelevantIf[T any](s VariableSet[T], _ ...RelevantIfOption) interface {
	RequiredOption
	OptionalOption
	DeprecatedOption
} {
	return option{
		ApplyToSpec: func(b variable.SpecBuilder) {
			b.Precondition(
				func() bool {
					x, ok := s.native()
					return ok && !reflect.ValueOf(x).IsZero()
				},
			)

			for _, vari := range s.variables() {
				dependsOn := variable.DependsOn{
					Subject:   b.Peek(),
					DependsOn: vari.Spec(),
				}

				variable.EstablishRelationships(
					variable.RefersTo{
						Subject:  b.Peek(),
						RefersTo: vari.Spec(),
					},
					dependsOn,
				)
			}
		},
	}
}

// RelevantWhen is an option that enables a variable set only if the value
// obtained from another set, s, produce the value v.
//
// An "irrelevant" variable set behaves as though its environment variables are
// undefined, irrespective of the actual values of the variables and any default
// values.
func RelevantWhen[T comparable](s VariableSet[T], v T, _ ...RelevantWhenOption) interface {
	RequiredOption
	OptionalOption
	DeprecatedOption
} {
	literals, err := s.literals(v)
	if err != nil {
		panic(fmt.Sprintf(
			"cannot use value as precondition: %s",
			err,
		))
	}

	return option{
		ApplyToSpec: func(b variable.SpecBuilder) {
			b.Precondition(
				func() bool {
					x, ok := s.native()
					return ok && x == v
				},
			)

			for i, vari := range s.variables() {
				dependsOn := variable.DependsOn{
					Subject:   b.Peek(),
					DependsOn: vari.Spec(),
					Value:     maybe.Some(literals[i]),
				}

				variable.EstablishRelationships(
					variable.RefersTo{
						Subject:  b.Peek(),
						RefersTo: vari.Spec(),
					},
					dependsOn,
				)
			}
		},
	}
}

// RelevantIfOption changes the behavior of the [RelevantIf] options.
type RelevantIfOption interface {
	future()
}

// RelevantWhenOption changes the behavior of the [RelevantWhen] option.
type RelevantWhenOption interface {
	future()
}
