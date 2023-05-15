package markdown

import (
	"path/filepath"

	"github.com/dogmatiq/ferrite/internal/mode"
)

// Run generates environment variable usage instructions in markdown format.
func Run(cfg mode.Config, options ...Option) {
	r := renderer{
		App:    filepath.Base(cfg.Args[0]),
		Specs:  cfg.Registry.Specs(),
		Output: cfg.Out,
	}

	for _, opt := range options {
		opt(&r)
	}

	r.Render()
	cfg.Exit(0)
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

// WithoutExplanatoryText disables the inclusion of the informational paragraphs
// in the rendered output.
func WithoutExplanatoryText() Option {
	return func(r *renderer) {
		r.withoutExplanatoryText = true
	}
}

// WithoutIndex disables the inclusion of the index in the rendered output.
func WithoutIndex() Option {
	return func(r *renderer) {
		r.withoutIndex = true
	}
}
