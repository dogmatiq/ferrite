package variable

// WithDocumentation returns an option that adds additional free-form
// information about a variable.
//
// The documentation string must be plain text (not Markdown).
//
// The T type parameter is not meaningful to the link itself, but is required in
// order to produce a SpecOption of the correct type.
func WithDocumentation[T any](doc string) SpecOption[T] {
	return func(opts *specOptions[T]) error {
		opts.Docs = append(opts.Docs, doc)
		return nil
	}
}
