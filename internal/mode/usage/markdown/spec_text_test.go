package markdown_test

import (
	"fmt"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	. "github.com/onsi/ginkgo/v2"
)

// textSpecValue is a test type that implements encoding.TextMarshaler and
// encoding.TextUnmarshaler via its pointer receiver.
type textSpecValue struct {
	Data string
}

func (v textSpecValue) MarshalText() ([]byte, error) {
	return []byte("text:" + v.Data), nil
}

func (v *textSpecValue) UnmarshalText(data []byte) error {
	if len(data) < 5 || string(data[:5]) != "text:" {
		return fmt.Errorf("expected text: prefix")
	}
	v.Data = string(data[5:])
	return nil
}

var _ = DescribeTable(
	"text spec",
	tableTest(
		"spec/text",
		WithoutExplanatoryText(),
		WithoutIndex(),
		WithoutUsageExamples(),
	),
	Entry(
		"deprecated",
		"deprecated.md",
		func(reg ferrite.Registry) {
			ferrite.
				TextAs[textSpecValue]("FERRITE_TEXT", "example text-marshaled value").
				WithExample(textSpecValue{Data: "hello"}, "a greeting").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"optional",
		"optional.md",
		func(reg ferrite.Registry) {
			ferrite.
				TextAs[textSpecValue]("FERRITE_TEXT", "example text-marshaled value").
				WithExample(textSpecValue{Data: "hello"}, "a greeting").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required",
		"required.md",
		func(reg ferrite.Registry) {
			ferrite.
				TextAs[textSpecValue]("FERRITE_TEXT", "example text-marshaled value").
				WithExample(textSpecValue{Data: "hello"}, "a greeting").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"required with default value",
		"with-default.md",
		func(reg ferrite.Registry) {
			ferrite.
				TextAs[textSpecValue]("FERRITE_TEXT", "example text-marshaled value").
				WithDefault(textSpecValue{Data: "fallback"}).
				Required(ferrite.WithRegistry(reg))
		},
	),
)
