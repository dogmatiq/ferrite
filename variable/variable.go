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
type Literal string

// needsQuotes is a pattern that matches values that require shell escaping.
var needsQuotes = regexp.MustCompile(`[^\w@%+=:,./-]`)

// String returns an escaped string representation of the value that can be used
// directly in the shell.
func (l Literal) String() string {
	s := string(l)

	if needsQuotes.MatchString(s) {
		return `'` + strings.ReplaceAll(s, `'`, `'"'"'`) + `'`
	}

	return s
}

// Any is an interface for an environment variable of any type.
type Any interface {
	// Spec returns the variable's specification.
	Spec() Spec

	// IsValid returns true if the variable is valid.
	IsValid() bool

	// Value returns the variable's value.
	Value() (maybe.Value[Value], ValidationError)
}

// OfType is an environment variable depicted by type T.
type OfType[T any] struct {
	spec SpecFor[T]
	env  Environment

	once  sync.Once
	value maybe.Value[valueOf[T]]
	err   ValidationError
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

	if v.value.IsEmpty() {
		return v.spec.isOptional
	}

	return true
}

// Value returns the variable's value.
func (v *OfType[T]) Value() (maybe.Value[Value], ValidationError) {
	v.resolve()
	return maybe.Map(
		v.value,
		func(v valueOf[T]) Value {
			return v
		},
	), v.err
}

// NativeValue returns the variable's native value.
func (v *OfType[T]) NativeValue() (maybe.Value[T], ValidationError) {
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

		if lit == "" {
			if n, ok := v.spec.def.Get(); ok {
				v.value = maybe.Some(valueOf[T]{
					native:    n,
					canonical: v.spec.class.Marshal(n),
					isDefault: true,
				})
			}

			return
		}

		n, c, err := v.spec.class.Unmarshal(v.spec.name, lit)
		if err != nil {
			v.err = err
			return
		}

		v.value = maybe.Some(valueOf[T]{
			verbatim:  lit,
			native:    n,
			canonical: c,
		})
	})
}
