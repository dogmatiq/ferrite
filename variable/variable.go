package variable

import (
	"regexp"
	"strings"
	"sync"

	"github.com/dogmatiq/ferrite/maybe"
)

// Name is the name of an environment variable.
type Name string

// String the string representation of a value, as used in the shell.
type String string

// needsQuotes is a pattern that matches values that require shell  escaping.
var needsQuotes = regexp.MustCompile(`[^\w@%+=:,./-]`)

// String returns an escaped string representation of the value that can be used
// directly in the shell.
func (s String) String() string {
	raw := string(s)

	if needsQuotes.MatchString(raw) {
		return `'` + strings.ReplaceAll(raw, `'`, `'"'"'`) + `'`
	}

	return raw
}

// Variable is an interface for an environment variable.
type Variable interface {
	// Name returns the name of the variable.
	Name() Name

	// Description returns a human-readable description of the variable.
	Description() string

	// Class returns the environment variable's class.
	Class() Class

	// Default returns the string representation of the default value.
	Default() maybe.Value[String]

	// Canonical returns the canonical string representation of the variable.
	Canonical() (maybe.Value[String], ValidationError)

	// Verbatim returns the string representation of the variable as it appears
	// in the environment.
	Verbatim() String

	// IsOptional returns true if the application can handle the absence of a
	// value for this variable.
	IsOptional() bool

	// IsDefault returns true if the value is the default, as opposed to being
	// set explicitly in the environment.
	IsDefault() bool
}

// TypedVariable is a variable depicted by type T.
type TypedVariable[T any] struct {
	spec Spec[T]
	env  Environment

	once      sync.Once
	value     maybe.Value[T]
	canonical maybe.Value[String]
	isDefault bool
	err       ValidationError
}

// Name returns the name of the variable.
func (v *TypedVariable[T]) Name() Name {
	return v.spec.Name
}

// Description returns a human-readable description of the variable.
func (v *TypedVariable[T]) Description() string {
	return v.spec.Description
}

// Class the environment varible's class.
func (v *TypedVariable[T]) Class() Class {
	return v.spec.Class
}

// Default returns the string representation of the default value.
func (v *TypedVariable[T]) Default() maybe.Value[String] {
	return maybe.Map(v.spec.Default, func(def T) String {
		return v.spec.Class.Marshal(def)
	})
}

// Value returns the native representation of the variable.
func (v *TypedVariable[T]) Value() (maybe.Value[T], ValidationError) {
	v.resolve()
	return v.value, v.err
}

// Canonical returns the canonical string representation of the variable.
func (v *TypedVariable[T]) Canonical() (maybe.Value[String], ValidationError) {
	v.resolve()
	return v.canonical, v.err
}

// Verbatim returns the string representation of the variable as it appears in
// the environment.
func (v *TypedVariable[T]) Verbatim() String {
	return v.env.Get(v.spec.Name)
}

// IsOptional returns true if the application can handle the absence of a value
// for this variable.
func (v *TypedVariable[T]) IsOptional() bool {
	return v.spec.IsOptional
}

// IsDefault returns true if the value is the default, as opposed to being set
// explicitly in the environment.
func (v *TypedVariable[T]) IsDefault() bool {
	v.resolve()
	return v.isDefault
}

func (v *TypedVariable[T]) resolve() {
	v.once.Do(func() {
		if s := v.Verbatim(); s != "" {
			n, c, err := v.spec.Class.Unmarshal(v.spec.Name, s)
			if err != nil {
				v.err = err
				return
			}

			v.value = maybe.Some(n)
			v.canonical = maybe.Some(c)
			return
		}

		if !v.spec.Default.IsEmpty() {
			v.value = v.spec.Default
			v.canonical = maybe.Map(v.value, v.spec.Class.Marshal)
			v.isDefault = true
		}
	})
}
