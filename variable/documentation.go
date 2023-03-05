package variable

import (
	"fmt"
	"strings"
)

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

// DocumentationBuilder is a fluent interface for building a Documentation
// value.
type DocumentationBuilder[T any] struct {
	doc Documentation
}

// WithDocumentation returns an option that adds additional documentation about
// a variable.
//
// The T type parameter is not meaningful, but is required in order to produce a
// SpecOption of the correct type.
func WithDocumentation[T any]() DocumentationBuilder[T] {
	return DocumentationBuilder[T]{}
}

// Summary adds a summary to the documentation.
func (b DocumentationBuilder[T]) Summary(f string, v ...any) DocumentationBuilder[T] {
	b.doc.Summary = fmt.Sprintf(f, v...)
	return b
}

// Paragraph adds a paragraph to the documentation.
func (b DocumentationBuilder[T]) Paragraph(text ...string) ParagraphFormatter[T] {
	return ParagraphFormatter[T]{
		func(v ...any) DocumentationBuilder[T] {
			b.doc.Paragraphs = append(
				b.doc.Paragraphs,
				fmt.Sprintf(
					strings.Join(text, " "),
					v...,
				),
			)
			return b
		},
	}
}

// ParagraphFormatter is a fluent interface for applying values to a paragraph
// template.
type ParagraphFormatter[T any] struct {
	Format func(...any) DocumentationBuilder[T]
}

// Important marks the documentation as important.
func (b DocumentationBuilder[T]) Important() DocumentationBuilder[T] {
	b.doc.IsImportant = true
	return b
}

// Done returns an option that adds the documentation to the variable spec.
func (b DocumentationBuilder[T]) Done() SpecOption[T] {
	return func(opts *specOptions[T]) error {
		opts.Docs = append(opts.Docs, b.doc)
		return nil
	}
}
