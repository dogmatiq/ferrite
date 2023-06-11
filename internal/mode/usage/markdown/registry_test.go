package markdown_test

import (
	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	"github.com/dogmatiq/ferrite/internal/variable"
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable(
	"func Run()",
	tableTest(
		"registry",
		WithoutExplanatoryText(),
		WithoutIndex(),
		WithoutUsageExamples(),
	),
	Entry(
		"with URL",
		"with-url.md",
		func(reg *variable.Registry) {
			*reg = *ferrite.NewRegistry(
				"3p",
				"Some Third-party Product",
				ferrite.WithDocumentationURL("https://example.org/docs/registry.html"),
			)

			ferrite.
				String("READ_DSN", "database connection string for read-models").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"without URL",
		"without-url.md",
		func(reg *variable.Registry) {
			*reg = *ferrite.NewRegistry(
				"3p",
				"Some Third-party Product",
			)

			ferrite.
				String("READ_DSN", "database connection string for read-models").
				Required(ferrite.WithRegistry(reg))
		},
	),
)
