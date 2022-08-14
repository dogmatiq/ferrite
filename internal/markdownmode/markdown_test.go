package markdownmode_test

// import (
// 	"bytes"
// 	"os"
// 	"path/filepath"
// 	"strings"

// 	. "github.com/dogmatiq/ferrite"
// 	. "github.com/jmalloc/gomegax"
// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// )

// var _ = Describe("func generateMarkdownUsage()", func() {
// 	buffer := &bytes.Buffer{}

// 	BeforeEach(func() {
// 		buffer.Reset()

// 		SetExitBehavior(
// 			buffer,
// 			func(code int) {
// 				Expect(code).To(Equal(0))
// 			},
// 		)

// 		os.Setenv("FERRITE_MODE", "usage/markdown")
// 	})

// 	AfterEach(func() {
// 		tearDown()
// 	})

// 	DescribeTable(
// 		"it describes the environment variable",
// 		func(file string, setup func(), lines ...string) {
// 			expect, err := os.ReadFile(filepath.Join("testdata", "markdown", file))
// 			Expect(err).ShouldNot(HaveOccurred())

// 			setup()
// 			Init()

// 			actual := buffer.String()
// 			ExpectWithOffset(1, actual).To(EqualX(string(expect)))
// 		},
// 		Entry(
// 			nil,
// 			"empty.md",
// 			func() {},
// 		),
// 		Entry(
// 			nil,
// 			"bool-required-default.md",
// 			func() {
// 				Bool("DEBUG", "enable debug mode").
// 					WithDefault(false).
// 					Required()
// 			},
// 		),
// 	)
// })

// // expectLines verifies that text consists of the given lines.
// func expectLines(buf *bytes.Buffer, lines ...string) {
// 	actual := buf.String()
// 	expect := strings.Join(lines, "\n") + "\n"
// 	ExpectWithOffset(1, actual).To(EqualX(expect))
// }
