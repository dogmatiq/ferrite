package variable

import (
	"regexp"
	"strings"
	"sync"

	"github.com/dogmatiq/ferrite/maybe"
)

// Name is the name of an environment variable.
type Name string

// Literal the string representation of an environment variable value.
type Literal struct {
	String string
}

// needsQuotes is a pattern that matches values that require shell escaping.
var needsQuotes = regexp.MustCompile(`[^\w@%+=:,./-]`)

// Quote returns an escaped string representation of the value that can be used
// directly in the shell.
func (l Literal) Quote() string {
	if !needsQuotes.MatchString(l.String) {
		return l.String
	}

	return `'` + strings.ReplaceAll(l.String, `'`, `'"'"'`) + `'`
}

// Any is an interface for an environment variable of any type.
type Any interface {
	// Spec returns the variable's specification.
	Spec() Spec

	// IsValid returns true if the variable is valid.
	IsValid() bool

	// Value returns the variable's value.
	Value() (maybe.Value[Value], ValueError)
}

// OfType is an environment variable depicted by type T.
type OfType[T any] struct {
	spec TypedSpec[T]
	env  Environment

	once  sync.Once
	value maybe.Value[valueOf[T]]
	err   ValueError
}

// Spec returns the variable's specification.
func (v *OfType[T]) Spec() Spec {
	return v.spec
}

// IsValid returns true if the variable is valid.
func (v *OfType[T]) IsValid() bool {
	v.resolve()

	if v.err != nil {
		return false
	}

	if v.spec.required && v.value.IsEmpty() {
		return false
	}

	return true
}

// Value returns the variable's value.
func (v *OfType[T]) Value() (maybe.Value[Value], ValueError) {
	v.resolve()
	return maybe.Map(
		v.value,
		func(v valueOf[T]) Value {
			return v
		},
	), v.err
}

// NativeValue returns the variable's native value.
func (v *OfType[T]) NativeValue() (maybe.Value[T], ValueError) {
	v.resolve()
	return maybe.Map(
		v.value,
		func(v valueOf[T]) T {
			return v.native
		},
	), v.err
}

func (v *OfType[T]) resolve() {
	v.once.Do(func() {
		lit := v.env.Get(v.spec.name)

		if lit.String == "" {
			v.value = v.spec.def
			return
		}

		n, err := v.spec.schema.Unmarshal(lit)
		if err != nil {
			v.err = valueError{
				name:    v.spec.name,
				literal: lit,
				cause:   err,
			}
			return
		}

		for _, val := range v.spec.validators {
			if err := val.Validate(n); err != nil {
				v.err = valueError{
					name:    v.spec.name,
					literal: lit,
					cause:   err,
				}
				return
			}
		}

		c, err := v.spec.schema.Marshal(n)
		if err != nil {
			panic(err)
		}

		v.value = maybe.Some(valueOf[T]{
			verbatim:  lit,
			native:    n,
			canonical: c,
		})
	})
}
