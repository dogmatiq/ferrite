package markdownmode

import (
	"io"

	"github.com/dogmatiq/ferrite/variable"
)

// Run generates environment variable usage instructions in markdown format.
func Run(
	reg *variable.Registry,
	app string,
	w io.Writer,
	options ...Option,
) {
	r := renderer{
		App:    app,
		Specs:  reg.Specs(),
		Output: w,
	}

	for _, opt := range options {
		opt(&r)
	}

	r.Render()
}

// Option is a function that changes the behavior of a renderer.
type Option func(*renderer)

// WithoutUsageExamples disables the inclusion of usage examples in the rendered
// output.
func WithoutUsageExamples() Option {
	return func(r *renderer) {
		r.withoutUsageExamples = true
	}
}

// WithoutPreamble disables the inclusion of the preamble in the rendered
// output.
func WithoutPreamble() Option {
	return func(r *renderer) {
		r.withoutPreamble = true
	}
}

// WithoutIndex disables the inclusion of the index in the rendered output.
func WithoutIndex() Option {
	return func(r *renderer) {
		r.withoutIndex = true
	}
}
