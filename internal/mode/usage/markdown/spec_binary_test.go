package markdown_test

import (
	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	"github.com/dogmatiq/ferrite/variable"
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable(
	"binary spec",
	tableTest(
		"spec/binary",
		WithoutExplanatoryText(),
		WithoutIndex(),
		WithoutUsageExamples(),
	),
	Entry(
		"deprecated",
		"deprecated.md",
		func(reg *variable.Registry) {
			ferrite.
				Binary("FAVICON", "the content of the favicon.png file").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional",
		"optional.md",
		func(reg *variable.Registry) {
			ferrite.
				Binary("FAVICON", "the content of the favicon.png file").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required",
		"required.md",
		func(reg *variable.Registry) {
			ferrite.
				Binary("FAVICON", "the content of the favicon.png file").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with default value",
		"with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Binary("FAVICON", "the content of the favicon.png file").
				WithDefault([]byte("<favicon content>")).
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with default value",
		"with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Binary("FAVICON", "the content of the favicon.png file").
				WithDefault([]byte("<favicon content>")).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with sensitive content",
		"with-sensitive-optional.md",
		func(reg *variable.Registry) {
			ferrite.
				Binary("SECRET_KEY", "a very secret machine-readable key").
				WithSensitiveContent().
				Optional(
					ferrite.WithRegistry(reg),
				)
		},
	),
	Entry(
		"required with sensitive content",
		"with-sensitive-required.md",
		func(reg *variable.Registry) {
			ferrite.
				Binary("SECRET_KEY", "a very secret machine-readable key").
				WithSensitiveContent().
				Required(
					ferrite.WithRegistry(reg),
				)
		},
	),
	Entry(
		"optional with sensitive content and default value",
		"with-sensitive-with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Binary("SECRET_KEY", "a very secret machine-readable key").
				WithDefault([]byte("hunter2")).
				WithSensitiveContent().
				Optional(
					ferrite.WithRegistry(reg),
				)
		},
	),
	Entry(
		"required with sensitive content and default value",
		"with-sensitive-with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Binary("SECRET_KEY", "a very secret machine-readable key").
				WithDefault([]byte("hunter2")).
				WithSensitiveContent().
				Required(
					ferrite.WithRegistry(reg),
				)
		},
	),
)
