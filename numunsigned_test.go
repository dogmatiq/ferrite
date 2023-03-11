package ferrite_test

import (
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type userDefinedUnsigned uint16

var _ = Describe("type UnsignedBuilder", func() {
	var builder *UnsignedBuilder[userDefinedUnsigned]

	BeforeEach(func() {
		builder = Unsigned[userDefinedUnsigned]("FERRITE_UNSIGNED", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	It("panics if the name is empty", func() {
		Expect(func() {
			Unsigned[userDefinedUnsigned]("", "<desc>").Optional()
		}).To(PanicWith("invalid specification: variable name must not be empty"))
	})

	It("panics if the description is empty", func() {
		Expect(func() {
			Unsigned[userDefinedUnsigned]("FERRITE_UNSIGNED", "").Optional()
		}).To(PanicWith("specification for FERRITE_UNSIGNED is invalid: variable description must not be empty"))
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
						`value of FERRITE_UNSIGNED (-1) is invalid: unrecognized uint16 syntax`,
					),
					Entry(
						"overflow",
						"65536",
						`value of FERRITE_UNSIGNED (65536) is invalid: too high, expected the largest uint16 value of 65535 or less`,
					),
					Entry(
						"decimal",
						"123.45",
						`value of FERRITE_UNSIGNED (123.45) is invalid: unrecognized uint16 syntax`,
					),
					Entry(
						"invalid characters",
						"123!",
						`value of FERRITE_UNSIGNED ('123!') is invalid: unrecognized uint16 syntax`,
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
						`value of FERRITE_UNSIGNED (-1) is invalid: unrecognized uint16 syntax`,
					),
					Entry(
						"overflow",
						"65536",
						`value of FERRITE_UNSIGNED (65536) is invalid: too high, expected the largest uint16 value of 65535 or less`,
					),
					Entry(
						"decimal",
						"123.45",
						`value of FERRITE_UNSIGNED (123.45) is invalid: unrecognized uint16 syntax`,
					),
					Entry(
						"invalid characters",
						"123!",
						`value of FERRITE_UNSIGNED ('123!') is invalid: unrecognized uint16 syntax`,
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
