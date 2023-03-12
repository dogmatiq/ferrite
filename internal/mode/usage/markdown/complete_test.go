package markdown_test

import (
	"github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/ferrite/variable"
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable(
	"func Run()",
	tableTest("complete"),
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
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"platform examples",
		"platform-examples.md",
		func(reg *variable.Registry) {
			ferrite.
				Bool("DEBUG", "enable or disable debugging features").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"platform examples use default values as examples when available",
		"platform-examples-use-defaults.md",
		func(reg *variable.Registry) {
			ferrite.
				NetworkPort("PORT", "an environment variable that has a default value").
				WithDefault("ftp").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"see also",
		"see-also.md",
		func(reg *variable.Registry) {
			verbose := ferrite.
				Bool("VERBOSE", "enable verbose logging").
				Optional(ferrite.WithRegistry(reg))

			ferrite.
				Bool("DEBUG", "enable or disable debugging features").
				SeeAlso(verbose).
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"deprecated + superseded",
		"deprecated-superseded.md",
		func(reg *variable.Registry) {
			addr := ferrite.
				String(`BIND_HOST`, "listen host for the HTTP server").
				WithDefault("0.0.0.0").
				Required(ferrite.WithRegistry(reg))

			port := ferrite.
				NetworkPort("BIND_PORT", "listen port for the HTTP server").
				WithDefault("8080").
				Required(ferrite.WithRegistry(reg))

			version := ferrite.
				String("BIND_VERSION", "IP version for the HTTP server").
				WithDefault("4").
				Required(ferrite.WithRegistry(reg))

			ferrite.
				String("BIND_ADDRESS", "listen address for the HTTP server").
				WithDefault("0.0.0.0:8080").
				Deprecated(
					ferrite.WithRegistry(reg),
					ferrite.SupersededBy(addr),
					ferrite.SupersededBy(port),
					ferrite.SupersededBy(version),
				)
		},
	),
)
