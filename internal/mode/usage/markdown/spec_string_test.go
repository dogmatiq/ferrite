package markdown_test

import (
	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable(
	"string spec",
	tableTest(
		"spec/string",
		WithoutExplanatoryText(),
		WithoutIndex(),
		WithoutUsageExamples(),
	),
	Entry(
		"deprecated",
		"deprecated.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("READ_DSN", "database connection string for read-models").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional",
		"optional.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("READ_DSN", "database connection string for read-models").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required",
		"required.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("READ_DSN", "database connection string for read-models").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with default value",
		"with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("READ_DSN", "database connection string for read-models").
				WithDefault("host=localhost dbname=readmodels user=projector").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with default value",
		"with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("READ_DSN", "database connection string for read-models").
				WithDefault("host=localhost dbname=readmodels user=projector").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"deprecated with exact length limit",
		"with-exact-length-deprecated.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("SEED", "the seed for the random-number generator").
				WithLength(5).
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with exact length limit",
		"with-exact-length-optional.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("SEED", "the seed for the random-number generator").
				WithLength(5).
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with exact length limit",
		"with-exact-length-required.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("SEED", "the seed for the random-number generator").
				WithLength(5).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"deprecated with maximum length limit",
		"with-max-length-deprecated.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("SEED", "the seed for the random-number generator").
				WithMaximumLength(10).
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with maximum length limit",
		"with-max-length-optional.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("SEED", "the seed for the random-number generator").
				WithMaximumLength(10).
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with maximum length limit",
		"with-max-length-required.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("SEED", "the seed for the random-number generator").
				WithMaximumLength(10).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"deprecated with minimum length limit",
		"with-min-length-deprecated.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("SEED", "the seed for the random-number generator").
				WithMinimumLength(5).
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with minimum length limit",
		"with-min-length-optional.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("SEED", "the seed for the random-number generator").
				WithMinimumLength(5).
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with minimum length limit",
		"with-min-length-required.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("SEED", "the seed for the random-number generator").
				WithMinimumLength(5).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"deprecated with minimum and maximum length limit",
		"with-minmax-length-deprecated.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("SEED", "the seed for the random-number generator").
				WithMinimumLength(5).
				WithMaximumLength(10).
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with minimum and maximum length limit",
		"with-minmax-length-optional.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("SEED", "the seed for the random-number generator").
				WithMinimumLength(5).
				WithMaximumLength(10).
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with minimum and maximum length limit",
		"with-minmax-length-required.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("SEED", "the seed for the random-number generator").
				WithMinimumLength(5).
				WithMaximumLength(10).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with sensitive content",
		"with-sensitive-optional.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("PASSWORD", "a very secret password").
				WithSensitiveContent().
				Optional(
					ferrite.WithRegistry(reg),
				)
		},
	),
	Entry(
		"required with sensitive content",
		"with-sensitive-required.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("PASSWORD", "a very secret password").
				WithSensitiveContent().
				Required(
					ferrite.WithRegistry(reg),
				)
		},
	),
	Entry(
		"optional with sensitive content and default value",
		"with-sensitive-with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("PASSWORD", "a very secret password").
				WithDefault("hunter2").
				WithSensitiveContent().
				Optional(
					ferrite.WithRegistry(reg),
				)
		},
	),
	Entry(
		"required with sensitive content and default value",
		"with-sensitive-with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				String("PASSWORD", "a very secret password").
				WithDefault("hunter2").
				WithSensitiveContent().
				Required(
					ferrite.WithRegistry(reg),
				)
		},
	),
)
