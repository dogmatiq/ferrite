package ferrite

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/dogmatiq/ferrite/maybe"
	"github.com/dogmatiq/ferrite/variable"
)

// Duration configures an environment variable as a duration.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Duration(name, desc string) DurationBuilder {
	return DurationBuilder{
		name: name,
		desc: desc,
	}
}

// DurationBuilder builds a specification for a duration variable.
type DurationBuilder struct {
	name, desc string
	def        maybe.Value[time.Duration]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b DurationBuilder) WithDefault(v time.Duration) DurationBuilder {
	b.def = maybe.Some(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b DurationBuilder) Required(options ...variable.RegisterOption) Required[time.Duration] {
	return registerRequired(b.spec(true), options)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b DurationBuilder) Optional(options ...variable.RegisterOption) Optional[time.Duration] {
	return registerOptional(b.spec(false), options)
}

func (b DurationBuilder) spec(req bool) variable.TypedSpec[time.Duration] {
	s, err := variable.NewSpec(
		b.name,
		b.desc,
		b.def,
		req,
		variable.TypedNumeric[time.Duration]{
			Marshaler: durationMarshaler{},
			NativeMin: maybe.Some(1 * time.Nanosecond),
		},
	)
	if err != nil {
		panic(err.Error())
	}

	return s
}

type durationMarshaler struct{}

func (durationMarshaler) Marshal(v time.Duration) (variable.Literal, error) {
	runes := []rune(v.String())
	i := len(runes) - 1

	// Skip over the last units.
	for !unicode.IsDigit(runes[i]) {
		i--
	}

	// Look for the second-to-last units, if any non-zero digit is encountered
	// then we need to keep the last units.
	for unicode.IsDigit(runes[i]) {
		if runes[i] != '0' {
			return variable.Literal{
				String: string(runes),
			}, nil
		}

		i--
	}

	// Otherwise the last units have a zero value and we omit them.
	return variable.Literal{
		String: string(runes[:i+1]),
	}, nil
}

func (durationMarshaler) Unmarshal(v variable.Literal) (time.Duration, error) {
	d, err := time.ParseDuration(strings.ReplaceAll(v.String, " ", ""))
	if err == nil {
		return d, nil
	}

	m := err.Error()
	if !strings.Contains(m, "unit") {
		return 0, err
	}

	return 0, errors.New(
		strings.Replace(
			strings.TrimPrefix(m, "time: "),
			fmt.Sprintf(` in duration %q`, v.String),
			"",
			1,
		),
	)
}
