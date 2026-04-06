package markdown

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable(
	"headingSlug",
	func(text, expected string) {
		Expect(headingSlug(text)).To(Equal(expected))
	},
	Entry(
		"uppercase environment variable name with backticks",
		"`COOKIE_SIGN_KEY`",
		"cookie_sign_key",
	),
	Entry(
		"uppercase environment variable name without backticks",
		"COOKIE_SIGN_KEY",
		"cookie_sign_key",
	),
	Entry(
		"already lowercase",
		"`read_dsn`",
		"read_dsn",
	),
	Entry(
		"spaces are converted to hyphens",
		"hello world",
		"hello-world",
	),
	Entry(
		"hyphens are preserved",
		"hello-world",
		"hello-world",
	),
	Entry(
		"punctuation is stripped",
		"hello, world!",
		"hello-world",
	),
	Entry(
		"unicode letters are preserved",
		"héllo",
		"héllo",
	),
	Entry(
		"digits are preserved",
		"v1.2.3",
		"v123",
	),
	Entry(
		"empty string",
		"",
		"",
	),
)
