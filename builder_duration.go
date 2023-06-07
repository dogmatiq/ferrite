package ferrite

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dogmatiq/ferrite/internal/maybe"
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

	b.builder.Name(name)
	b.builder.Description(desc)
	b.builder.Documentation().
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
	schema  variable.TypedNumeric[time.Duration]
	builder variable.TypedSpecBuilder[time.Duration]
}

var _ isBuilderOf[time.Duration, *DurationBuilder]

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *DurationBuilder) WithDefault(v time.Duration) *DurationBuilder {
	b.builder.Default(v)
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
func (b *DurationBuilder) Required(options ...RequiredOption) Required[time.Duration] {
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *DurationBuilder) Optional(options ...OptionalOption) Optional[time.Duration] {
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *DurationBuilder) Deprecated(options ...DeprecatedOption) Deprecated[time.Duration] {
	return deprecated(b.schema, &b.builder, options...)
}

type durationMarshaler struct{}

func (durationMarshaler) Marshal(v time.Duration) (variable.Literal, error) {
	runes := []rune(v.String())
	zeroes := false

loop:
	for end := len(runes); end > 0; end-- {
		switch runes[end-1] {
		case '0':
			// Keep track of when we're in a run of zeroes.
			zeroes = true
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			// If we find any non-zero digits then we need to keep the current
			// units.
			break loop
		default:
			// If we were traversing a run of zeroes and have found another unit
			// string (s, ms, us, etc) then we can discard everything after this
			// point.
			if zeroes {
				zeroes = false
				runes = runes[:end]
			}
		}
	}

	return variable.Literal{
		String: string(runes),
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
