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
			Environment: variable.MemoryEnvironment{},
		}
	})

	DescribeTable(
		"it describes the environment variable",
		func(
			file string,
			setup func(*variable.Registry),
		) {
			setup(reg)

			expect, err := os.ReadFile(
				filepath.Join(
					"testdata",
					"markdown",
					"usage",
					file,
				),
			)
			Expect(err).ShouldNot(HaveOccurred())

			actual := Run("<app>", reg)
			ExpectWithOffset(1, actual).To(EqualX(string(expect)))
		},
		Entry(
			"usage",
			"usage.md",
			func(reg *variable.Registry) {
				ferrite.
					Bool("DEBUG", "enable or disable debugging features").
					Optional(variable.WithRegistry(reg))
			},
		),
		Entry(
			"usage shows the default value in examples if available",
			"usage-shows-default.md",
			func(reg *variable.Registry) {
				ferrite.
					NetworkPort("PORT", "an environment variable that has a default value").
					WithDefault("ftp").
					Required(variable.WithRegistry(reg))
			},
		),
	)
})
