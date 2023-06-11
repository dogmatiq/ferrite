package markdown_test

import (
	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable(
	"bool spec",
	tableTest(
		"spec/bool",
		WithoutExplanatoryText(),
		WithoutIndex(),
		WithoutUsageExamples(),
	),
	Entry(
		"deprecated",
		"deprecated.md",
		func(reg ferrite.Registry) {
			ferrite.
				Bool("DEBUG", "enable or disable debugging features").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional",
		"optional.md",
		func(reg ferrite.Registry) {
			ferrite.
				Bool("DEBUG", "enable or disable debugging features").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required",
		"required.md",
		func(reg ferrite.Registry) {
			ferrite.
				Bool("DEBUG", "enable or disable debugging features").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with default value",
		"with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				Bool("DEBUG", "enable or disable debugging features").
				WithDefault(false).
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with default value",
		"with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				Bool("DEBUG", "enable or disable debugging features").
				WithDefault(false).
				Required(ferrite.WithRegistry(reg))
		},
	),
)
