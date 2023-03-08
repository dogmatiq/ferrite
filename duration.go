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
//
// Durations have a minimum value of 1 nanosecond by default.
func Duration(name, desc string) *DurationBuilder {
	b := &DurationBuilder{
		schema: variable.TypedNumeric[time.Duration]{
			Marshaler: durationMarshaler{},
			NativeMin: maybe.Some(1 * time.Nanosecond),
		},
	}

	b.spec.Name(name)
	b.spec.Description(desc)
	b.spec.Documentation().
		Summary("Duration syntax").
		Paragraph(
			"Durations are specified as a sequence of decimal numbers, each with an optional fraction and a unit suffix, such as `300ms`, `-1.5h` or `2h45m`.",
			"Supported time units are `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`.",
		).
		Format().
		Done()

	return b
}

// DurationBuilder builds a specification for a duration variable.
type DurationBuilder struct {
	schema variable.TypedNumeric[time.Duration]
	spec   variable.TypedSpecBuilder[time.Duration]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *DurationBuilder) WithDefault(v time.Duration) *DurationBuilder {
	b.spec.Default(v)
	return b
}

// WithMinimum sets the minimum acceptable value of the variable.
func (b *DurationBuilder) WithMinimum(v time.Duration) *DurationBuilder {
	b.schema.NativeMin = maybe.Some(v)
	return b
}

// WithMaximum sets the maximum acceptable value of the variable.
func (b *DurationBuilder) WithMaximum(v time.Duration) *DurationBuilder {
	b.schema.NativeMax = maybe.Some(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *DurationBuilder) Required(options ...Option) Required[time.Duration] {
	return req(b.schema, &b.spec, options)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *DurationBuilder) Optional(options ...Option) Optional[time.Duration] {
	return opt(b.schema, &b.spec, options)
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
