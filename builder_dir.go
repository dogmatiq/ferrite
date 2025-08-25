package ferrite

import (
	"errors"
	"os"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// Dir configures an environment variable as a filesystem directory name.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Dir(name, desc string) *DirBuilder {
	b := &DirBuilder{}
	b.builder.Name(name)
	b.builder.Description(desc)
	b.builder.NonNormativeExample("/path/to/dir", "an absolute directory path")
	b.builder.NonNormativeExample("./path/to/dir", "a relative directory path")
	return b
}

// DirBuilder builds a specification for a string value that is the name of a
// filesystem directory.
type DirBuilder struct {
	schema  variable.TypedString[DirName]
	builder variable.TypedSpecBuilder[DirName]
}

var _ isBuilderOf[
	DirName,
	string,
	*DirBuilder,
]

// WithDefault sets the default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *DirBuilder) WithDefault(v string) *DirBuilder {
	b.builder.Default(variable.ConstDefault(DirName(v)))
	return b
}

// WithExample adds an example value to the variable's documentation.
func (b *DirBuilder) WithExample(v string, desc string) *DirBuilder {
	b.builder.NormativeExample(DirName(v), desc)
	return b
}

// WithMustExist adds a constraint that requires the value to refer to an
// existing directory.
func (b *DirBuilder) WithMustExist() *DirBuilder {
	b.builder.BuiltInConstraint(
		"**MUST** refer to a directory that already exists",
		func(ctx variable.ConstraintContext, v DirName) variable.ConstraintError {
			if ctx != variable.ConstraintContextFinal {
				return nil
			}

			info, err := os.Stat(string(v))
			if err != nil {
				if os.IsNotExist(err) {
					return errors.New("expected the directory to exist")
				}
				return err
			}

			if !info.IsDir() {
				return errors.New("the path refers to a file, expected a directory")
			}

			return nil
		},
	)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *DirBuilder) Required(options ...RequiredOption) Required[DirName] {
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *DirBuilder) Optional(options ...OptionalOption) Optional[DirName] {
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *DirBuilder) Deprecated(options ...DeprecatedOption) Deprecated[DirName] {
	return deprecated(b.schema, &b.builder, options...)
}

// DirName is the name of a directory.
type DirName string
