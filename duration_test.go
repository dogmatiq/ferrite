package ferrite_test

import (
	"errors"
	"os"
	"time"

	. "github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/ferrite/schema"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type DurationSpec", func() {
	type customNumeric int16

	var spec *DurationSpec

	BeforeEach(func() {
		spec = Duration("FERRITE_DURATION", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	Describe("func Describe()", func() {
		It("describes the variable", func() {
			Expect(spec.Describe()).To(ConsistOf(
				VariableXXX{
					Name:        "FERRITE_DURATION",
					Description: "<desc>",
					Schema: schema.Range{
						Min: "1ns",
					},
				},
			))
		})
	})

	When("the environment variable is not empty", func() {
		BeforeEach(func() {
			os.Setenv("FERRITE_DURATION", "630s")
		})

		Describe("func Value()", func() {
			It("returns the numeric value", func() {
				Expect(spec.Value()).To(Equal(10*time.Minute + 30*time.Second))
			})
		})

		Describe("func Validate()", func() {
			It("returns a successful result", func() {
				Expect(spec.Validate()).To(ConsistOf(
					ValidationResult{
						Name:  "FERRITE_DURATION",
						Value: "10m30s",
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
					os.Setenv("FERRITE_DURATION", value)
					Expect(func() {
						spec.Value()
					}).To(PanicWith(expect))
				},
				Entry(
					"zero",
					"0s",
					`FERRITE_DURATION is invalid: must be a positive duration`,
				),
				Entry(
					"negative",
					"-1s",
					`FERRITE_DURATION is invalid: must be a positive duration`,
				),
			)
		})

		Describe("func Validate()", func() {
			DescribeTable(
				"it returns a failure result",
				func(value, expect string) {
					os.Setenv("FERRITE_DURATION", value)
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:  "FERRITE_DURATION",
							Error: errors.New(expect),
						},
					))
				},
				Entry(
					"zero",
					"0s",
					`must be a positive duration`,
				),
				Entry(
					"negative",
					"-1s",
					`must be a positive duration`,
				),
			)
		})
	})

	When("the environment variable is empty", func() {
		When("there is a default value", func() {
			BeforeEach(func() {
				spec.WithDefault(630 * time.Second)
			})

			Describe("func Value()", func() {
				It("returns the default", func() {
					Expect(spec.Value()).To(Equal(630 * time.Second))
				})
			})

			Describe("func Describe()", func() {
				It("describes the variable", func() {
					Expect(spec.Describe()).To(ConsistOf(
						VariableXXX{
							Name:        "FERRITE_DURATION",
							Description: "<desc>",
							Schema: schema.Range{
								Min: "1ns",
							},
							Default: "10m30s",
						},
					))
				})
			})

			Describe("func Validate()", func() {
				It("returns a success result", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:        "FERRITE_DURATION",
							Value:       "10m30s",
							UsedDefault: true,
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
					}).To(PanicWith("FERRITE_DURATION is invalid: must not be empty"))
				})
			})

			Describe("func Validate()", func() {
				It("returns a failure result", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:  "FERRITE_DURATION",
							Error: errors.New(`must not be empty`),
						},
					))
				})
			})
		})
	})
})
