package ferrite

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/dogmatiq/ferrite/variable"
)

// File configures an environment variable as a filename.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func File(name, desc string) *FileBuilder {
	b := &FileBuilder{}
	b.spec.Name(name)
	b.spec.Description(desc)
	return b
}

// FileBuilder builds a specification for a boolean value.
type FileBuilder struct {
	schema variable.TypedString[FileName]
	spec   variable.TypedSpecBuilder[FileName]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *FileBuilder) WithDefault(v string) *FileBuilder {
	b.spec.Default(FileName(v))
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *FileBuilder) Required(options ...VariableOption) Required[FileName] {
	return req(b.schema, &b.spec, options)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *FileBuilder) Optional(options ...VariableOption) Optional[FileName] {
	return opt(b.schema, &b.spec, options)
}

// FileName is the name of a file.
type FileName string

// Reader returns a reader that produces the contents of the file.
func (n FileName) Reader() (io.ReadCloser, error) {
	return os.Open(string(n))
}

// ReadBytes returns the contents of the file as a byte slice.
func (n FileName) ReadBytes() ([]byte, error) {
	return ioutil.ReadFile(string(n))
}

// ReadString returns the contents of the file as a string.
func (n FileName) ReadString() (string, error) {
	data, err := n.ReadBytes()
	return string(data), err
}
