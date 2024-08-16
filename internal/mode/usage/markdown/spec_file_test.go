package markdown_test

import (
	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable(
	"file spec",
	tableTest(
		"spec/file",
		WithoutExplanatoryText(),
		WithoutIndex(),
		WithoutUsageExamples(),
	),
	Entry(
		"deprecated",
		"deprecated.md",
		func(reg ferrite.Registry) {
			ferrite.
				File("PRIVATE_KEY", "path to the private key file").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional",
		"optional.md",
		func(reg ferrite.Registry) {
			ferrite.
				File("PRIVATE_KEY", "path to the private key file").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required",
		"required.md",
		func(reg ferrite.Registry) {
			ferrite.
				File("PRIVATE_KEY", "path to the private key file").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with default value",
		"with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				File("PRIVATE_KEY", "path to the private key file").
				WithDefault("/etc/ssh/id_rsa").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with default value",
		"with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				File("PRIVATE_KEY", "path to the private key file").
				WithDefault("/etc/ssh/id_rsa").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"deprecated with must-exist requirement",
		"with-must-exist-deprecated.md",
		func(reg ferrite.Registry) {
			ferrite.
				File("PRIVATE_KEY", "path to the private key file").
				WithMustExist().
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with must-exist requirement",
		"with-must-exist-optional.md",
		func(reg ferrite.Registry) {
			ferrite.
				File("PRIVATE_KEY", "path to the private key file").
				WithMustExist().
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with must-exist requirement",
		"with-must-exist-required.md",
		func(reg ferrite.Registry) {
			ferrite.
				File("PRIVATE_KEY", "path to the private key file").
				WithMustExist().
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional with must-exist requirement and default value",
		"with-must-exist-with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				File("PRIVATE_KEY", "path to the private key file").
				WithDefault("/etc/ssh/id_rsa").
				WithMustExist().
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with must-exist requirement and default value",
		"with-must-exist-with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				File("PRIVATE_KEY", "path to the private key file").
				WithDefault("/etc/ssh/id_rsa").
				WithMustExist().
				Required(ferrite.WithRegistry(reg))
		},
	),
)
