package ferrite

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/dogmatiq/ferrite/maybe"
	"github.com/dogmatiq/ferrite/variable"
)

// File configures an environment variable as a filename.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func File(name, desc string) *FileBuilder {
	return &FileBuilder{
		name: name,
		desc: desc,
	}
}

// FileBuilder builds a specification for a boolean value.
type FileBuilder struct {
	name, desc string
	def        maybe.Value[FileName]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *FileBuilder) WithDefault(v string) *FileBuilder {
	b.def = maybe.Some(FileName(v))
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *FileBuilder) Required(options ...variable.RegisterOption) Required[FileName] {
	v := variable.Register(b.spec(true), options)
	return requiredVar[FileName]{v}
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *FileBuilder) Optional(options ...variable.RegisterOption) Optional[FileName] {
	v := variable.Register(b.spec(false), options)
	return optionalVar[FileName]{v}

}

func (b *FileBuilder) spec(req bool) variable.TypedSpec[FileName] {
	s, err := variable.NewSpec(
		b.name,
		b.desc,
		b.def,
		req,
		variable.TypedString[FileName]{},
	)
	if err != nil {
		panic(err.Error())
	}

	return s
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
