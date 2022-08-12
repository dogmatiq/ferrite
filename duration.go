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
		spec: variable.PendingSpecFor[time.Duration]{
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
	spec variable.PendingSpecFor[time.Duration]
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

	// Trim any trailing zero units, e.g. "10m0s" -> "10m".
	for i := len(runes) - 2; i >= 0; i-- {
		if runes[i+1] == '0' && !unicode.IsDigit(runes[i]) {
			runes = runes[:i+1]
			break
		}
	}

	return variable.Literal(runes), nil
}

func (durationMarshaler) Unmarshal(v variable.Literal) (time.Duration, error) {
	d, err := time.ParseDuration(string(v))
	if err == nil {
		return d, nil
	}

	m := err.Error()
	if !strings.Contains(m, "unit") {
		return 0, errors.New("must be a valid duration, e.g. 10m30s")
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
