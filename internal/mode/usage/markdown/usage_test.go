package markdown_test

import (
	"github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/ferrite/variable"
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable(
	"func Run()",
	tableTest("usage"),
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
