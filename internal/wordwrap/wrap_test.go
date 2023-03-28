package wordwrap_test

import (
	"fmt"
	"strings"

	"github.com/dogmatiq/ferrite/internal/wordwrap"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Wrap()", func() {
	DescribeTable(
		"it wraps text to the specified number of columns",
		func(text string, expect []string) {
			actual := wordwrap.Wrap(text, 80)
			Expect(actual).To(Equal(expect))
		},
		Entry(
			"empty text",
			"",
			nil,
		),
		Entry(
			"ascii only, single line, byte-size less than max columns",
			"This variable is **sensitive**.",
			[]string{
				"This variable is **sensitive**.",
			},
		),
		Entry(
			"multi-byte character, single line, byte-size less than max columns",
			"⚠️ This variable is **sensitive**.",
			[]string{
				"⚠️ This variable is **sensitive**.",
			},
		),
		Entry(
			"multi-byte character, single line, byte-size greater than max columns",
			"⚠️ This variable is **sensitive**; its value may contain private information.",
			[]string{
				"⚠️ This variable is **sensitive**; its value may contain private information.",
			},
		),
		Entry(
			"multi-byte character, multiple lines",
			"⚠️ The application may consume other undocumented environment variables. This document only shows variables declared using [Ferrite].",
			[]string{
				"⚠️ The application may consume other undocumented environment variables. This",
				"document only shows variables declared using [Ferrite].",
			},
		),
		Entry(
			"single-line, single word that occupies entire line",
			strings.Repeat("X", 80),
			[]string{
				strings.Repeat("X", 80),
			},
		),
		Entry(
			"single-line, single word that does not fit on a line",
			strings.Repeat("X", 81),
			[]string{
				strings.Repeat("X", 81),
			},
		),
		Entry(
			"multi-line, single word that occupies entire line",
			fmt.Sprintf(
				"foo %s bar",
				strings.Repeat("X", 80),
			),
			[]string{
				"foo",
				strings.Repeat("X", 80),
				"bar",
			},
		),
		Entry(
			"multi-line, single word that does not fit on a line",
			fmt.Sprintf(
				"foo %s bar",
				strings.Repeat("X", 81),
			),
			[]string{
				"foo",
				strings.Repeat("X", 81),
				"bar",
			},
		),
		Entry(
			"multi-line, wrapping exactly at max columns",
			"The `DEBUG` variable **MAY** be left undefined. Otherwise, the value **MUST** be either `true` or `false`.",
			[]string{
				"The `DEBUG` variable **MAY** be left undefined. Otherwise, the value **MUST** be",
				"either `true` or `false`.",
			},
		),
		Entry(
			"explicit newlines",
			"foo\nbar",
			[]string{
				"foo",
				"bar",
			},
		),
		Entry(
			"explicit newlines with empty line",
			"foo\n\nbar",
			[]string{
				"foo",
				"",
				"bar",
			},
		),
		Entry(
			"first part of hyphenated word falls exactly on max columns",
			"⚠️ The `API_URL` variable is **deprecated**; its use is **NOT RECOMMENDED** as it may be removed in a future version. If defined, the value **MUST** be a fully-qualified URL.",
			[]string{
				"⚠️ The `API_URL` variable is **deprecated**; its use is **NOT RECOMMENDED** as",
				"it may be removed in a future version. If defined, the value **MUST** be a",
				"fully-qualified URL.",
			},
		),
		Entry(
			"second part of hyphenated word exceeds max columns",
			"⚠️ The `WEIGHT` variable is **deprecated**; its use is **NOT RECOMMENDED** as it may be removed in a future version. If defined, the value **MUST** be a non-negative whole number.",
			[]string{
				"⚠️ The `WEIGHT` variable is **deprecated**; its use is **NOT RECOMMENDED** as it",
				"may be removed in a future version. If defined, the value **MUST** be a non-",
				"negative whole number.",
			},
		),
	)
})
