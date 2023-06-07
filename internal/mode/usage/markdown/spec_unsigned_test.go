package markdown_test

import (
	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	"github.com/dogmatiq/ferrite/internal/variable"
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable(
	"unsigned spec",
	tableTest(
		"spec/unsigned",
		WithoutExplanatoryText(),
		WithoutIndex(),
		WithoutUsageExamples(),
	),
	Entry(
		"deprecated",
		"deprecated.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional",
		"optional.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required",
		"required.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with default value",
		"with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				WithDefault(900).
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with default value",
		"with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				WithDefault(900).
				Required(ferrite.WithRegistry(reg))
		},
	), Entry(
		"with minimum value",
		"with-min.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				WithMinimum(10).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"with maximum value",
		"with-max.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				WithMaximum(20).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"with minimum and maximum values",
		"with-minmax.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				WithMinimum(10).
				WithMaximum(20).
				Required(ferrite.WithRegistry(reg))
		},
	),
)
