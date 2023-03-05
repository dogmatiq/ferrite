package variable

// Documentation is free-form documentation about a variable.
type Documentation struct {
	// Summary is a short summary of the documentation.
	//
	// It may be empty. If provided it must be plain text.
	Summary string

	// Paragraphs is a list of paragraphs to show.
	//
	// Simple inline Markdown formatting is allowed, but the value used in
	// contexts where the plain text is shown directly to the user.
	Paragraphs []string

	// IsImportant indicates that the documentation is important and should be
	// made obvious to the user.
	IsImportant bool
}

// WithDocumentation returns an option that adds additional documentation about
// a variable.
//
// The T type parameter is not meaningful, but is required in order to produce a
// SpecOption of the correct type.
func WithDocumentation[T any](docs ...Documentation) SpecOption[T] {
	return func(opts *specOptions[T]) error {
		for _, d := range docs {
			opts.Docs = append(opts.Docs, d)
		}

		return nil
	}
}
