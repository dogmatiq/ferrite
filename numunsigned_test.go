package ferrite_test

import (
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type userDefinedUnsigned uint16

var _ = Describe("type UnsignedBuilder", func() {

	var builder UnsignedBuilder[userDefinedUnsigned]

	BeforeEach(func() {
		builder = Unsigned[userDefinedUnsigned]("FERRITE_UNSIGNED", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	When("the variable is required", func() {
		When("the value is valid", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"returns the value",
					func(value string, expect int) {
						os.Setenv("FERRITE_UNSIGNED", value)

						v := builder.
							Required().
							Value()

						Expect(v).To(Equal(userDefinedUnsigned(expect)))
					},
					Entry("zero", "0", 0),
					Entry("positive", "123", +123),
				)
			})
		})

		When("the value is invalid", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it panics",
					func(value, expect string) {
						os.Setenv("FERRITE_UNSIGNED", value)

						Expect(func() {
							builder.
								Required().
								Value()
						}).To(PanicWith(expect))
					},
					Entry(
						"underflow",
						"-1",
						`value of FERRITE_UNSIGNED (-1) is invalid: strconv.ParseUint: parsing "-1": invalid syntax`,
					),
					Entry(
						"overflow",
						"65536",
						`value of FERRITE_UNSIGNED (65536) is invalid: strconv.ParseUint: parsing "65536": value out of range`,
					),
					Entry(
						"decimal",
						"123.45",
						`value of FERRITE_UNSIGNED (123.45) is invalid: strconv.ParseUint: parsing "123.45": invalid syntax`,
					),
					Entry(
						"invalid characters",
						"123!",
						`value of FERRITE_UNSIGNED ('123!') is invalid: strconv.ParseUint: parsing "123!": invalid syntax`,
					),
				)
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				When("there is a default value", func() {
					Describe("func Value()", func() {
						It("returns the default", func() {
							v := builder.
								WithDefault(123).
								Required().
								Value()

							Expect(v).To(Equal(userDefinedUnsigned(123)))
						})
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
							"FERRITE_UNSIGNED is undefined and does not have a default value",
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
					func(value string, expect int) {
						os.Setenv("FERRITE_UNSIGNED", value)

						v, ok := builder.
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(userDefinedUnsigned(expect)))
					},
					Entry("zero", "0", 0),
					Entry("positive", "123", 123),
				)
			})
		})

		When("the value is invalid", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it panics",
					func(value, expect string) {
						os.Setenv("FERRITE_UNSIGNED", value)

						Expect(func() {
							builder.
								Optional().
								Value()
						}).To(PanicWith(expect))
					},
					Entry(
						"underflow",
						"-1",
						`value of FERRITE_UNSIGNED (-1) is invalid: strconv.ParseUint: parsing "-1": invalid syntax`,
					),
					Entry(
						"overflow",
						"65536",
						`value of FERRITE_UNSIGNED (65536) is invalid: strconv.ParseUint: parsing "65536": value out of range`,
					),
					Entry(
						"decimal",
						"123.45",
						`value of FERRITE_UNSIGNED (123.45) is invalid: strconv.ParseUint: parsing "123.45": invalid syntax`,
					),
					Entry(
						"invalid characters",
						"123!",
						`value of FERRITE_UNSIGNED ('123!') is invalid: strconv.ParseUint: parsing "123!": invalid syntax`,
					),
				)
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							WithDefault(123).
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(userDefinedUnsigned(123)))
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
				os.Setenv("FERRITE_UNSIGNED", "1")

				builder.
					WithMinimum(5).
					Required().
					Value()
			}).To(PanicWith(
				`value of FERRITE_UNSIGNED (1) is invalid: too low, expected 5 or greater`,
			))
		})
	})

	When("the value is greater than the maximum limit", func() {
		It("panics", func() {
			Expect(func() {
				os.Setenv("FERRITE_UNSIGNED", "10")

				builder.
					WithMaximum(5).
					Required().
					Value()
			}).To(PanicWith(
				`value of FERRITE_UNSIGNED (10) is invalid: too high, expected 5 or less`,
			))
		})
	})
})
