package markdown_test

import (
	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable(
	"dir spec",
	tableTest(
		"spec/dir",
		WithoutExplanatoryText(),
		WithoutIndex(),
		WithoutUsageExamples(),
	),
	Entry(
		"deprecated",
		"deprecated.md",
		func(reg ferrite.Registry) {
			ferrite.
				Dir("CACHE_DIR", "path to the cache directory").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional",
		"optional.md",
		func(reg ferrite.Registry) {
			ferrite.
				Dir("CACHE_DIR", "path to the cache directory").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required",
		"required.md",
		func(reg ferrite.Registry) {
			ferrite.
				Dir("CACHE_DIR", "path to the cache directory").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with default value",
		"with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				Dir("CACHE_DIR", "path to the cache directory").
				WithDefault("/var/run/cache").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with default value",
		"with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				Dir("CACHE_DIR", "path to the cache directory").
				WithDefault("/var/run/cache").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"deprecated with must-exist requirement",
		"with-must-exist-deprecated.md",
		func(reg ferrite.Registry) {
			ferrite.
				Dir("CACHE_DIR", "path to the cache directory").
				WithMustExist().
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with must-exist requirement",
		"with-must-exist-optional.md",
		func(reg ferrite.Registry) {
			ferrite.
				Dir("CACHE_DIR", "path to the cache directory").
				WithMustExist().
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with must-exist requirement",
		"with-must-exist-required.md",
		func(reg ferrite.Registry) {
			ferrite.
				Dir("CACHE_DIR", "path to the cache directory").
				WithMustExist().
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with must-exist requirement and default value",
		"with-must-exist-with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				Dir("CACHE_DIR", "path to the cache directory").
				WithDefault("/var/run/cache").
				WithMustExist().
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with must-exist requirement and default value",
		"with-must-exist-with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				Dir("CACHE_DIR", "path to the cache directory").
				WithDefault("/var/run/cache").
				WithMustExist().
				Required(ferrite.WithRegistry(reg))
		},
	),
)
