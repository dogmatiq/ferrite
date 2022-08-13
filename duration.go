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
	b := DurationBuilder{
		spec: variable.PendingSpec[time.Duration]{
			Name:        variable.Name(name),
			Description: desc,
		},
	}

	s, err := variable.NewNumber(
		durationMarshaler{},
		maybe.Some(1*time.Nanosecond),
		maybe.None[time.Duration](),
	)
	if err != nil {
		b.spec.InvalidErr(err)
	}

	b.spec.Schema = s
	return b
}

// DurationBuilder builds a specification for a duration variable.
type DurationBuilder struct {
	spec variable.PendingSpec[time.Duration]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b DurationBuilder) WithDefault(v time.Duration) DurationBuilder {
	b.spec.Default = maybe.Some(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b DurationBuilder) Required(options ...variable.RegisterOption) Required[time.Duration] {
	return req[time.Duration]{variable.Register(b.spec, options)}
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b DurationBuilder) Optional(options ...variable.RegisterOption) Optional[time.Duration] {
	b.spec.IsOptional = true
	return opt[time.Duration]{variable.Register(b.spec, options)}
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
			return variable.Literal(runes), nil
		}

		i--
	}

	// Otherwise the last units have a zero value and we omit them.
	return variable.Literal(runes[:i+1]), nil
}

func (durationMarshaler) Unmarshal(v variable.Literal) (time.Duration, error) {
	d, err := time.ParseDuration(strings.ReplaceAll(string(v), " ", ""))
	if err == nil {
		return d, nil
	}

	m := err.Error()
	if !strings.Contains(m, "unit") {
		return 0, errors.New("expected a valid duration, e.g. 10m30s")
	}

	return 0, errors.New(
		strings.Replace(
			strings.TrimPrefix(m, "time: "),
			fmt.Sprintf(` in duration %q`, string(v)),
			"",
			1,
		),
	)
}
