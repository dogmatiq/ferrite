package markdown_test

import (
	"github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable(
	"func Run()",
	tableTest("complete"),
	Entry(
		"no variables",
		"empty.md",
		func(reg ferrite.Registry) {},
	),
	Entry(
		"non-normative examples",
		"non-normative.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("READ_DSN", "database connection string for read-models").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"user-defined normative examples",
		"normative.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("READ_DSN", "database connection string for read-models").
				WithExample("postgres://user:pass@localhost:5432/dbname", "user-defined example").
				Required(ferrite.WithRegistry(reg))
		},
	),
)
