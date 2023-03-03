package markdownmode

import (
	"github.com/dogmatiq/ferrite/variable"
)

// Run generates environment variable usage instructions in markdown format.
func Run(
	app string,
	reg *variable.Registry,
	options ...Option,
) string {
	r := renderer{
		App:   app,
		Specs: reg.Specs(),
	}

	for _, opt := range options {
		opt(&r)
	}

	return r.Render()
}

// Option is a function that changes the behavior of a renderer.
type Option func(*renderer)

// WithoutUsageExamples disables the inclusion of usage examples in the rendered
// output.
func WithoutUsageExamples() Option {
	return func(r *renderer) {
		r.hideUsageExamples = true
	}
}
