package markdown_test

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/dogmatiq/ferrite/internal/mode"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	"github.com/dogmatiq/ferrite/variable"
	. "github.com/jmalloc/gomegax"
	. "github.com/onsi/gomega"
)

func tableTest(
	path string,
	options ...Option,
) func(
	file string,
	setup func(*variable.Registry),
) {
	return func(
		file string,
		setup func(*variable.Registry),
	) {
		reg := &variable.Registry{
			Environment: &variable.MemoryEnvironment{},
		}

		setup(reg)

		expect, err := os.ReadFile(
			filepath.Join(
				"testdata",
				path,
				file,
			),
		)
		Expect(err).ShouldNot(HaveOccurred())

		actual := &bytes.Buffer{}
		exited := false

		Run(
			mode.Options{
				Registry: reg,
				Args:     []string{"<app>"},
				Out:      actual,
				Exit: func(code int) {
					exited = true
					Expect(code).To(Equal(0))
				},
			},
			options...,
		)
		ExpectWithOffset(1, actual.String()).To(EqualX(string(expect)))
		Expect(exited).To(BeTrue())
	}
}
