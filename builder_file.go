package ferrite

import (
	"errors"
	"io"
	"os"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// File configures an environment variable as a filename.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func File(name, desc string) *FileBuilder {
	b := &FileBuilder{}
	b.builder.Name(name)
	b.builder.Description(desc)
	b.builder.NonNormativeExample("/path/to/file", "an absolute file path")
	b.builder.NonNormativeExample("./path/to/file", "a relative file path")
	return b
}

// FileBuilder builds a specification for a string value that is the name of a
// file.
type FileBuilder struct {
	schema  variable.TypedString[FileName]
	builder variable.TypedSpecBuilder[FileName]
}

var _ isBuilderOf[
	FileName,
	string,
	*FileBuilder,
]

// WithDefault sets the default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *FileBuilder) WithDefault(v string) *FileBuilder {
	b.builder.Default(variable.ConstDefault(FileName(v)))
	return b
}

// WithExample adds an example value to the variable's documentation.
func (b *FileBuilder) WithExample(v string, desc string) *FileBuilder {
	b.builder.NormativeExample(FileName(v), desc)
	return b
}

// WithMustExist adds a constraint that requires the value to refer to an
// existing file.
func (b *FileBuilder) WithMustExist() *FileBuilder {
	b.builder.BuiltInConstraint(
		"**MUST** refer to a file that already exists",
		func(ctx variable.ConstraintContext, v FileName) variable.ConstraintError {
			if ctx != variable.ConstraintContextFinal {
				return nil
			}

			info, err := os.Stat(string(v))
			if err != nil {
				if os.IsNotExist(err) {
					return errors.New("expected the file to exist")
				}
				return err
			}

			if info.IsDir() {
				return errors.New("the path refers to a directory, expected a file")
			}

			return nil
		},
	)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *FileBuilder) Required(options ...RequiredOption) Required[FileName] {
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *FileBuilder) Optional(options ...OptionalOption) Optional[FileName] {
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *FileBuilder) Deprecated(options ...DeprecatedOption) Deprecated[FileName] {
	return deprecated(b.schema, &b.builder, options...)
}

// FileName is the name of a file.
type FileName string

// Reader returns a reader that produces the contents of the file.
func (n FileName) Reader() (io.ReadCloser, error) {
	return os.Open(string(n))
}

// ReadBytes returns the contents of the file as a byte slice.
func (n FileName) ReadBytes() ([]byte, error) {
	return os.ReadFile(string(n))
}

// ReadString returns the contents of the file as a string.
func (n FileName) ReadString() (string, error) {
	data, err := n.ReadBytes()
	return string(data), err
}
