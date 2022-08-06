package ferrite_test

import (
	"errors"
	"os"

	. "github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/ferrite/schema"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type UnsignedSpec", func() {
	type customNumeric uint16

	var spec *UnsignedSpec[customNumeric]

	BeforeEach(func() {
		spec = Unsigned[customNumeric]("FERRITE_UNSIGNED", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	When("the environment variable is not empty", func() {
		BeforeEach(func() {
			os.Setenv("FERRITE_UNSIGNED", "123")
		})

		Describe("func Value()", func() {
			It("returns the numeric value", func() {
				Expect(spec.Value()).To(Equal(customNumeric(123)))
			})
		})

		Describe("func Validate()", func() {
			It("returns a successful result", func() {
				Expect(spec.Validate()).To(ConsistOf(
					ValidationResult{
						Name:        "FERRITE_UNSIGNED",
						Description: "<desc>",
						Schema: schema.Range{
							Min: "0",
							Max: "65535",
						},
						ExplicitValue: "123",
					},
				))
			})
		})
	})

	When("the environment variable is invalid", func() {
		Describe("func Value()", func() {
			DescribeTable(
				"it panics",
				func(value, expect string) {
					os.Setenv("FERRITE_UNSIGNED", value)
					Expect(func() {
						spec.Value()
					}).To(PanicWith(expect))
				},
				Entry(
					"underflow",
					"-1",
					`FERRITE_UNSIGNED is invalid: must be an integer between 0 and 65535`,
				),
				Entry(
					"overflow",
					"65537",
					`FERRITE_UNSIGNED is invalid: must be an integer between 0 and 65535`,
				),
				Entry(
					"decimal",
					"123.45",
					`FERRITE_UNSIGNED is invalid: must be an integer between 0 and 65535`,
				),
				Entry(
					"invalid characters",
					"123!",
					`FERRITE_UNSIGNED is invalid: must be an integer between 0 and 65535`,
				),
			)
		})

		Describe("func Validate()", func() {
			DescribeTable(
				"it returns a failure result",
				func(value, expect string) {
					os.Setenv("FERRITE_UNSIGNED", value)
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:        "FERRITE_UNSIGNED",
							Description: "<desc>",
							Schema: schema.Range{
								Min: "0",
								Max: "65535",
							},
							ExplicitValue: value,
							Error:         errors.New(expect),
						},
					))
				},
				Entry(
					"underflow",
					"-1",
					`must be an integer between 0 and 65535`,
				),
				Entry(
					"overflow",
					"65537",
					`must be an integer between 0 and 65535`,
				),
				Entry(
					"decimal",
					"123.45",
					`must be an integer between 0 and 65535`,
				),
				Entry(
					"invalid characters",
					"123!",
					`must be an integer between 0 and 65535`,
				),
			)
		})
	})

	When("the environment variable is empty", func() {
		When("there is a default value", func() {
			BeforeEach(func() {
				spec.WithDefault(123)
			})

			Describe("func Value()", func() {
				It("returns the default", func() {
					Expect(spec.Value()).To(Equal(customNumeric(123)))
				})
			})

			Describe("func Validate()", func() {
				It("returns a success result", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:        "FERRITE_UNSIGNED",
							Description: "<desc>",
							Schema: schema.Range{
								Min: "0",
								Max: "65535",
							},
							DefaultValue: "123",
							UsingDefault: true,
						},
					))
				})
			})
		})

		When("there is no default value", func() {
			Describe("func Value()", func() {
				It("panics", func() {
					Expect(func() {
						spec.Value()
					}).To(PanicWith("FERRITE_UNSIGNED is invalid: must not be empty"))
				})
			})

			Describe("func Validate()", func() {
				It("returns a failure result", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:        "FERRITE_UNSIGNED",
							Description: "<desc>",
							Schema: schema.Range{
								Min: "0",
								Max: "65535",
							},
							Error: errors.New(`must not be empty`),
						},
					))
				})
			})
		})
	})
})
