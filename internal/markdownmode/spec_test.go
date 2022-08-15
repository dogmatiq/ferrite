package markdownmode_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"time"

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

		// BOOL

		Entry(
			"bool + optional + default",
			"bool/default.md",
			func(reg *variable.Registry) {
				ferrite.
					Bool("DEBUG", "enable or disable debugging features").
					WithDefault(false).
					Optional(variable.WithRegistry(reg))
			},
		),
		Entry(
			"bool + optional",
			"bool/optional.md",
			func(reg *variable.Registry) {
				ferrite.
					Bool("DEBUG", "enable or disable debugging features").
					Optional(variable.WithRegistry(reg))
			},
		),
		Entry(
			"bool + required + default",
			"bool/default.md",
			func(reg *variable.Registry) {
				ferrite.
					Bool("DEBUG", "enable or disable debugging features").
					WithDefault(false).
					Required(variable.WithRegistry(reg))
			},
		),
		Entry(
			"bool + required",
			"bool/required.md",
			func(reg *variable.Registry) {
				ferrite.
					Bool("DEBUG", "enable or disable debugging features").
					Required(variable.WithRegistry(reg))
			},
		),

		// ENUM

		Entry(
			"enum + optional + default",
			"enum/default.md",
			func(reg *variable.Registry) {
				ferrite.
					Enum("LOG_LEVEL", "the minimum log level to record").
					WithMember("debug", "show information for developers").
					WithMember("info", "standard log messages").
					WithMember("warn", "important, but don't need individual human review").
					WithMember("error", "a healthy application shouldn't produce any errors").
					WithMember("fatal", "the application cannot proceed").
					WithDefault("error").
					Optional(variable.WithRegistry(reg))
			},
		),
		Entry(
			"enum + optional",
			"enum/optional.md",
			func(reg *variable.Registry) {
				ferrite.
					Enum("LOG_LEVEL", "the minimum log level to record").
					WithMember("debug", "show information for developers").
					WithMember("info", "standard log messages").
					WithMember("warn", "important, but don't need individual human review").
					WithMember("error", "a healthy application shouldn't produce any errors").
					WithMember("fatal", "the application cannot proceed").
					Optional(variable.WithRegistry(reg))
			},
		),
		Entry(
			"enum + required + default",
			"enum/default.md",
			func(reg *variable.Registry) {
				ferrite.
					Enum("LOG_LEVEL", "the minimum log level to record").
					WithMember("debug", "show information for developers").
					WithMember("info", "standard log messages").
					WithMember("warn", "important, but don't need individual human review").
					WithMember("error", "a healthy application shouldn't produce any errors").
					WithMember("fatal", "the application cannot proceed").
					WithDefault("error").
					Required(variable.WithRegistry(reg))
			},
		),
		Entry(
			"enum + required",
			"enum/required.md",
			func(reg *variable.Registry) {
				ferrite.
					Enum("LOG_LEVEL", "the minimum log level to record").
					WithMember("debug", "show information for developers").
					WithMember("info", "standard log messages").
					WithMember("warn", "important, but don't need individual human review").
					WithMember("error", "a healthy application shouldn't produce any errors").
					WithMember("fatal", "the application cannot proceed").
					Required(variable.WithRegistry(reg))
			},
		),

		// DURATION

		Entry(
			"duration + optional + default",
			"duration/default.md",
			func(reg *variable.Registry) {
				ferrite.
					Duration("GRPC_TIMEOUT", "gRPC request timeout").
					WithDefault(10 * time.Millisecond).
					Optional(variable.WithRegistry(reg))
			},
		),
		Entry(
			"duration + optional",
			"duration/optional.md",
			func(reg *variable.Registry) {
				ferrite.
					Duration("GRPC_TIMEOUT", "gRPC request timeout").
					Optional(variable.WithRegistry(reg))
			},
		),
		Entry(
			"duration + required + default",
			"duration/default.md",
			func(reg *variable.Registry) {
				ferrite.
					Duration("GRPC_TIMEOUT", "gRPC request timeout").
					WithDefault(10 * time.Millisecond).
					Required(variable.WithRegistry(reg))
			},
		),
		Entry(
			"duration + required",
			"duration/required.md",
			func(reg *variable.Registry) {
				ferrite.
					Duration("GRPC_TIMEOUT", "gRPC request timeout").
					Required(variable.WithRegistry(reg))
			},
		),

		// STRING

		Entry(
			"string + optional + default",
			"string/default.md",
			func(reg *variable.Registry) {
				ferrite.
					String("READ_DSN", "database connection string for read-models").
					WithDefault("host=localhost dbname=readmodels user=projector").
					Optional(variable.WithRegistry(reg))
			},
		),
		Entry(
			"string + optional",
			"string/optional.md",
			func(reg *variable.Registry) {
				ferrite.
					String("READ_DSN", "database connection string for read-models").
					Optional(variable.WithRegistry(reg))
			},
		),
		Entry(
			"string + required + default",
			"string/default.md",
			func(reg *variable.Registry) {
				ferrite.
					String("READ_DSN", "database connection string for read-models").
					WithDefault("host=localhost dbname=readmodels user=projector").
					Required(variable.WithRegistry(reg))
			},
		),
		Entry(
			"string + required",
			"string/required.md",
			func(reg *variable.Registry) {
				ferrite.
					String("READ_DSN", "database connection string for read-models").
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
