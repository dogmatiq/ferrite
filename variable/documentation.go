package variable

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
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

// DocumentationBuilderX is a fluent interface for building a Documentation
// value.
// Deprecated:
type DocumentationBuilderX[T any] struct {
	doc Documentation
}

// WithDocumentation returns documentation builder that adds additional
// documentation about a variable.
//
// The T type parameter is not meaningful, but is required in order to produce a
// TypedSpecOption of the correct type.
func WithDocumentation[T any]() DocumentationBuilderX[T] {
	return DocumentationBuilderX[T]{}
}

// Summary adds a summary to the documentation.
func (b DocumentationBuilderX[T]) Summary(f string, v ...any) DocumentationBuilderX[T] {
	b.doc.Summary = fmt.Sprintf(f, v...)
	return b
}

// Paragraph adds a paragraph to the documentation.
//
// text is concatenated together with a space to form the paragraph text.
// The entire paragraph is a Printf() style format specifier.
func (b DocumentationBuilderX[T]) Paragraph(text ...string) ParagraphFormatterX[T] {
	return ParagraphFormatterX[T]{
		func(v ...any) DocumentationBuilderX[T] {
			b.doc.Paragraphs = append(
				slices.Clone(b.doc.Paragraphs),
				fmt.Sprintf(
					strings.Join(text, " "),
					v...,
				),
			)
			return b
		},
	}
}

// ParagraphFormatterX is a fluent interface for applying values to a paragraph
// template.
// Deprecated:
type ParagraphFormatterX[T any] struct {
	// Format applies the given values to the paragraph template.
	Format func(...any) DocumentationBuilderX[T]
}

// Important marks the documentation as important.
func (b DocumentationBuilderX[T]) Important() DocumentationBuilderX[T] {
	b.doc.IsImportant = true
	return b
}

// Done returns an option that adds the documentation to the variable spec.
func (b DocumentationBuilderX[T]) Done() TypedSpecOption[T] {
	return func(opts *specOptions[T]) error {
		opts.Docs = append(opts.Docs, b.doc)
		return nil
	}
}

// DocumentationBuilder is a fluent interface for building a documentation.
type DocumentationBuilder struct {
	docs *[]Documentation
	doc  Documentation
}

// Paragraph adds a paragraph to the documentation.
//
// text is concatenated together with a space to form the paragraph text.
// The entire paragraph is a Printf() style format specifier.
func (b DocumentationBuilder) Paragraph(text ...string) ParagraphFormatter {
	return ParagraphFormatter{
		func(v ...any) DocumentationBuilder {
			b.doc.Paragraphs = append(
				slices.Clone(b.doc.Paragraphs),
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
type ParagraphFormatter struct {
	// Format applies the given values to the paragraph template.
	Format func(...any) DocumentationBuilder
}

// Important marks the documentation as important.
func (b DocumentationBuilder) Important() DocumentationBuilder {
	b.doc.IsImportant = true
	return b
}

// Done returns an option that adds the documentation to the variable spec.
func (b DocumentationBuilder) Done() {
	*b.docs = append(*b.docs, b.doc)
}
