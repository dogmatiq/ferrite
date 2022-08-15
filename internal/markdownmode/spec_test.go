package markdownmode_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

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

			expect, err := os.ReadFile(filepath.Join("testdata", "markdown", file))
			Expect(err).ShouldNot(HaveOccurred())

			actual := Run("<app>", reg, false)
			ExpectWithOffset(1, actual).To(EqualX(string(expect)))
		},
		Entry(
			nil,
			"empty.md",
			func(reg *variable.Registry) {},
		),
		Entry(
			nil,
			"bool-optional-default.md",
			func(reg *variable.Registry) {
				ferrite.
					Bool("DEBUG", "enable or disable debugging features").
					WithDefault(false).
					Optional(variable.WithRegistry(reg))
			},
		),
		Entry(
			nil,
			"bool-optional.md",
			func(reg *variable.Registry) {
				ferrite.
					Bool("DEBUG", "enable or disable debugging features").
					Optional(variable.WithRegistry(reg))
			},
		),
		Entry(
			nil,
			"bool-required-default.md",
			func(reg *variable.Registry) {
				ferrite.
					Bool("DEBUG", "enable or disable debugging features").
					WithDefault(false).
					Required(variable.WithRegistry(reg))
			},
		),
		Entry(
			nil,
			"bool-required.md",
			func(reg *variable.Registry) {
				ferrite.
					Bool("DEBUG", "enable or disable debugging features").
					Required(variable.WithRegistry(reg))
			},
		),
	)
})

// expectLines verifies that text consists of the given lines.
func expectLines(buf *bytes.Buffer, lines ...string) {
	actual := buf.String()
	expect := strings.Join(lines, "\n") + "\n"
	ExpectWithOffset(1, actual).To(EqualX(expect))
}
