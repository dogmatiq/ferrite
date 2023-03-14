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
)
