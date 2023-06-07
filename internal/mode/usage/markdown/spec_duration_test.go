package markdown_test

import (
	"time"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	"github.com/dogmatiq/ferrite/internal/variable"
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable(
	"duration spec",
	tableTest(
		"spec/duration",
		WithoutExplanatoryText(),
		WithoutIndex(),
		WithoutUsageExamples(),
	),
	Entry(
		"deprecated",
		"deprecated.md",
		func(reg *variable.Registry) {
			ferrite.
				Duration("GRPC_TIMEOUT", "gRPC request timeout").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),

	Entry(
		"optional",
		"optional.md",
		func(reg *variable.Registry) {
			ferrite.
				Duration("GRPC_TIMEOUT", "gRPC request timeout").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required",
		"required.md",
		func(reg *variable.Registry) {
			ferrite.
				Duration("GRPC_TIMEOUT", "gRPC request timeout").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with default value",
		"with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Duration("GRPC_TIMEOUT", "gRPC request timeout").
				WithDefault(10 * time.Millisecond).
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with default value",
		"with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Duration("GRPC_TIMEOUT", "gRPC request timeout").
				WithDefault(10 * time.Millisecond).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"with maximum value",
		"with-max.md",
		func(reg *variable.Registry) {
			ferrite.
				Duration("GRPC_TIMEOUT", "gRPC request timeout").
				WithMaximum(24 * time.Hour).
				Required(ferrite.WithRegistry(reg))
		},
	),
)
