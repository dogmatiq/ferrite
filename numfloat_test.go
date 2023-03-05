package ferrite_test

import (
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type userDefinedFloat float32

var _ = Describe("type FloatBuilder", func() {
	var builder FloatBuilder[userDefinedFloat]

	BeforeEach(func() {
		builder = Float[userDefinedFloat]("FERRITE_FLOAT", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	When("the variable is required", func() {
		When("the value is valid", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"returns the value",
					func(value string, expect float64) {
						os.Setenv("FERRITE_FLOAT", value)

						v := builder.
							Required().
							Value()

						Expect(v).To(Equal(userDefinedFloat(expect)))
					},
					Entry("zero", "0", 0.0),
					Entry("positive", "+123.45", +123.45),
					Entry("negative", "-123.45", -123.45),
				)
			})
		})

		When("the value is invalid", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it panics",
					func(value, expect string) {
						os.Setenv("FERRITE_FLOAT", value)

						Expect(func() {
							builder.
								Required().
								Value()
						}).To(PanicWith(expect))
					},
					Entry(
						"underflow",
						"-3.40282346638528859811704183484516925440e+39",
						`value of FERRITE_FLOAT (-3.40282346638528859811704183484516925440e+39) is invalid: strconv.ParseFloat: parsing "-3.40282346638528859811704183484516925440e+39": value out of range`,
					),
					Entry(
						"overflow",
						"3.40282346638528859811704183484516925440e+39",
						`value of FERRITE_FLOAT (3.40282346638528859811704183484516925440e+39) is invalid: strconv.ParseFloat: parsing "3.40282346638528859811704183484516925440e+39": value out of range`,
					),
					Entry(
						"invalid characters",
						"123!",
						`value of FERRITE_FLOAT ('123!') is invalid: strconv.ParseFloat: parsing "123!": invalid syntax`,
					),
				)
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							WithDefault(-123).
							Required().
							Value()

						Expect(v).To(Equal(userDefinedFloat(-123)))
					})
				})
			})

			When("there is no default value", func() {
				Describe("func Value()", func() {
					It("panics", func() {
						Expect(func() {
							builder.
								Required().
								Value()
						}).To(PanicWith(
							"FERRITE_FLOAT is undefined and does not have a default value",
						))
					})
				})
			})
		})
	})

	When("the variable is optional", func() {
		When("the value is valid", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"returns the value",
					func(value string, expect float64) {
						os.Setenv("FERRITE_FLOAT", value)

						v, ok := builder.
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(userDefinedFloat(expect)))
					},
					Entry("zero", "0", 0.0),
					Entry("positive", "+123.45", +123.45),
					Entry("negative", "-123.45", -123.45),
				)
			})
		})

		When("the value is invalid", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it panics",
					func(value, expect string) {
						os.Setenv("FERRITE_FLOAT", value)

						Expect(func() {
							builder.
								Optional().
								Value()
						}).To(PanicWith(expect))
					},
					Entry(
						"underflow",
						"-3.40282346638528859811704183484516925440e+39",
						`value of FERRITE_FLOAT (-3.40282346638528859811704183484516925440e+39) is invalid: strconv.ParseFloat: parsing "-3.40282346638528859811704183484516925440e+39": value out of range`,
					),
					Entry(
						"overflow",
						"3.40282346638528859811704183484516925440e+39",
						`value of FERRITE_FLOAT (3.40282346638528859811704183484516925440e+39) is invalid: strconv.ParseFloat: parsing "3.40282346638528859811704183484516925440e+39": value out of range`,
					),
					Entry(
						"invalid characters",
						"123!",
						`value of FERRITE_FLOAT ('123!') is invalid: strconv.ParseFloat: parsing "123!": invalid syntax`,
					),
				)
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							WithDefault(-123).
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(userDefinedFloat(-123)))
					})
				})
			})

			When("there is no default value", func() {
				Describe("func Value()", func() {
					It("returns with ok == false", func() {
						_, ok := builder.
							Optional().
							Value()

						Expect(ok).To(BeFalse())
					})
				})
			})
		})
	})

	When("the value is lower than the minimum limit", func() {
		It("panics", func() {
			Expect(func() {
				os.Setenv("FERRITE_FLOAT", "-1.1")

				builder.
					WithMinimum(5.5).
					Required().
					Value()
			}).To(PanicWith(
				`value of FERRITE_FLOAT (-1.1) is invalid: too low, expected +5.5 or greater`,
			))
		})
	})

	When("the value is greater than the maximum limit", func() {
		It("panics", func() {
			Expect(func() {
				os.Setenv("FERRITE_FLOAT", "10.1")

				builder.
					WithMaximum(5.5).
					Required().
					Value()
			}).To(PanicWith(
				`value of FERRITE_FLOAT (10.1) is invalid: too high, expected +5.5 or less`,
			))
		})
	})
})
