package markdown_test

import (
	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	"github.com/dogmatiq/ferrite/variable"
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable(
	"func Run()",
	tableTest(
		"relationship",
		WithoutPreamble(),
		WithoutUsageExamples(),
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
				Optional(
					ferrite.WithRegistry(reg),
					ferrite.SeeAlso(verbose),
				)
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
