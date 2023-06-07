package markdown_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/dogmatiq/ferrite/internal/mode"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	"github.com/dogmatiq/ferrite/internal/variable"
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
			mode.Config{
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

		// Split strings into lines which producers a more human-friendly diff
		// in case of a failure.
		actualLines := strings.Split(actual.String(), "\n")
		expectLines := strings.Split(string(expect), "\n")

		ExpectWithOffset(1, actualLines).To(EqualX(expectLines))
		Expect(exited).To(BeTrue())
	}
}
