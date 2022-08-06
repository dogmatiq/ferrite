package ferrite_test

import (
	"errors"
	"os"

	. "github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/ferrite/schema"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type SignedSpec", func() {
	type customNumeric int16

	var spec *SignedSpec[customNumeric]

	BeforeEach(func() {
		spec = Signed[customNumeric]("FERRITE_SIGNED", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	When("the environment variable is not empty", func() {
		BeforeEach(func() {
			os.Setenv("FERRITE_SIGNED", "-123")
		})

		Describe("func Value()", func() {
			It("returns the numeric value", func() {
				Expect(spec.Value()).To(Equal(customNumeric(-123)))
			})
		})

		Describe("func Validate()", func() {
			It("returns a successful result", func() {
				Expect(spec.Validate()).To(ConsistOf(
					ValidationResult{
						Name:        "FERRITE_SIGNED",
						Description: "<desc>",
						Schema: schema.Range{
							Min: "-32768",
							Max: "+32767",
						},
						ExplicitValue: "-123",
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
					os.Setenv("FERRITE_SIGNED", value)
					Expect(func() {
						spec.Value()
					}).To(PanicWith(expect))
				},
				Entry(
					"underflow",
					"-32769",
					`FERRITE_SIGNED is invalid: must be an integer between -32768 and +32767`,
				),
				Entry(
					"overflow",
					"32768",
					`FERRITE_SIGNED is invalid: must be an integer between -32768 and +32767`,
				),
				Entry(
					"decimal",
					"123.45",
					`FERRITE_SIGNED is invalid: must be an integer between -32768 and +32767`,
				),
				Entry(
					"invalid characters",
					"123!",
					`FERRITE_SIGNED is invalid: must be an integer between -32768 and +32767`,
				),
			)
		})

		Describe("func Validate()", func() {
			DescribeTable(
				"it returns a failure result",
				func(value, expect string) {
					os.Setenv("FERRITE_SIGNED", value)
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:        "FERRITE_SIGNED",
							Description: "<desc>",
							Schema: schema.Range{
								Min: "-32768",
								Max: "+32767",
							},
							ExplicitValue: value,
							Error:         errors.New(expect),
						},
					))
				},
				Entry(
					"underflow",
					"-32769",
					`must be an integer between -32768 and +32767`,
				),
				Entry(
					"overflow",
					"32768",
					`must be an integer between -32768 and +32767`,
				),
				Entry(
					"decimal",
					"123.45",
					`must be an integer between -32768 and +32767`,
				),
				Entry(
					"invalid characters",
					"123!",
					`must be an integer between -32768 and +32767`,
				),
			)
		})
	})

	When("the environment variable is empty", func() {
		When("there is a default value", func() {
			BeforeEach(func() {
				spec.WithDefault(-123)
			})

			Describe("func Value()", func() {
				It("returns the default", func() {
					Expect(spec.Value()).To(Equal(customNumeric(-123)))
				})
			})

			Describe("func Validate()", func() {
				It("returns a success result", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:        "FERRITE_SIGNED",
							Description: "<desc>",
							Schema: schema.Range{
								Min: "-32768",
								Max: "+32767",
							},
							DefaultValue: "-123",
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
					}).To(PanicWith("FERRITE_SIGNED is invalid: must not be empty"))
				})
			})

			Describe("func Validate()", func() {
				It("returns a failure result", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:        "FERRITE_SIGNED",
							Description: "<desc>",
							Schema: schema.Range{
								Min: "-32768",
								Max: "+32767",
							},
							Error: errors.New(`must not be empty`),
						},
					))
				})
			})
		})
	})
})
