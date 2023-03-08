package markdownmode_test

import (
	"os"
	"path/filepath"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite/internal/markdownmode"
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

			actual := Run(
				"<app>",
				reg,
				WithoutUsageExamples(),
			)
			ExpectWithOffset(1, actual).To(EqualX(string(expect)))
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
					Required(variable.UseRegisterOptionsWithBuilder(variable.WithRegistry(reg)))
			},
		),
	)
})
