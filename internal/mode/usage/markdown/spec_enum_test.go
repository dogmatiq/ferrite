package markdown_test

import (
	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	"github.com/dogmatiq/ferrite/variable"
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable(
	"enum spec",
	tableTest(
		"spec/enum",
		WithoutExplanatoryText(),
		WithoutIndex(),
		WithoutUsageExamples(),
	),
	Entry(
		"deprecated",
		"deprecated.md",
		func(reg *variable.Registry) {
			ferrite.
				Enum("LOG_LEVEL", "the minimum log level to record").
				WithMember("debug", "show information for developers").
				WithMember("info", "standard log messages").
				WithMember("warn", "important, but don't need individual human review").
				WithMember("error", "a healthy application shouldn't produce any errors").
				WithMember("fatal", "the application cannot proceed").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional",
		"optional.md",
		func(reg *variable.Registry) {
			ferrite.
				Enum("LOG_LEVEL", "the minimum log level to record").
				WithMember("debug", "show information for developers").
				WithMember("info", "standard log messages").
				WithMember("warn", "important, but don't need individual human review").
				WithMember("error", "a healthy application shouldn't produce any errors").
				WithMember("fatal", "the application cannot proceed").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required",
		"required.md",
		func(reg *variable.Registry) {
			ferrite.
				Enum("LOG_LEVEL", "the minimum log level to record").
				WithMember("debug", "show information for developers").
				WithMember("info", "standard log messages").
				WithMember("warn", "important, but don't need individual human review").
				WithMember("error", "a healthy application shouldn't produce any errors").
				WithMember("fatal", "the application cannot proceed").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with default value",
		"with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Enum("LOG_LEVEL", "the minimum log level to record").
				WithMember("debug", "show information for developers").
				WithMember("info", "standard log messages").
				WithMember("warn", "important, but don't need individual human review").
				WithMember("error", "a healthy application shouldn't produce any errors").
				WithMember("fatal", "the application cannot proceed").
				WithDefault("error").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with default value",
		"with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Enum("LOG_LEVEL", "the minimum log level to record").
				WithMember("debug", "show information for developers").
				WithMember("info", "standard log messages").
				WithMember("warn", "important, but don't need individual human review").
				WithMember("error", "a healthy application shouldn't produce any errors").
				WithMember("fatal", "the application cannot proceed").
				WithDefault("error").
				Required(ferrite.WithRegistry(reg))
		},
	),
)
