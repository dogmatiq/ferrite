package markdownmode_test

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/ferrite/internal/mode"
	. "github.com/dogmatiq/ferrite/internal/mode/markdownmode"
	"github.com/dogmatiq/ferrite/variable"
	. "github.com/jmalloc/gomegax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Run()", func() {
	var reg *variable.Registry

	BeforeEach(func() {
		reg = &variable.Registry{
			Environment: &variable.MemoryEnvironment{},
		}
	})

	DescribeTable(
		"it generates the correct preamble",
		func(
			file string,
			setup func(*variable.Registry),
		) {
			setup(reg)

			expect, err := os.ReadFile(
				filepath.Join(
					"testdata",
					"markdown",
					"preamble",
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
				WithoutUsageExamples(),
			)
			ExpectWithOffset(1, actual.String()).To(EqualX(string(expect)))
			Expect(exited).To(BeTrue())
		},
		Entry(
			"no variables",
			"empty.md",
			func(reg *variable.Registry) {},
		),
		Entry(
			"non-normative examples",
			"non-normative.md",
			func(reg *variable.Registry) {
				ferrite.
					String("READ_DSN", "database connection string for read-models").
					Required(variable.WithRegistry(reg))
			},
		),
	)
})
