package markdownmode_test

import (
	"os"
	"path/filepath"

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
			Environment: variable.MemoryEnvironment{},
		}
	})

	DescribeTable(
		"it generates a useful preamble",
		func(
			file string,
			setup func(*variable.Registry),
		) {
			setup(reg)

			expect, err := os.ReadFile(filepath.Join("testdata", "markdown", file))
			Expect(err).ShouldNot(HaveOccurred())

			actual := Run(
				"<app>",
				reg,
				WithoutUsageExamples(),
			)
			ExpectWithOffset(1, actual).To(EqualX(string(expect)))
		},
		Entry(
			nil,
			"empty.md",
			func(reg *variable.Registry) {},
		),
	)
})
