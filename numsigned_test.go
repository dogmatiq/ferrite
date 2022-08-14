package ferrite_test

import (
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type userDefinedSigned int16

var _ = Describe("type SignedSpec", func() {
	var builder SignedBuilder[userDefinedSigned]

	BeforeEach(func() {
		builder = Signed[userDefinedSigned]("FERRITE_SIGNED", "<desc>")
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
						os.Setenv("FERRITE_SIGNED", value)

						v := builder.
							Required().
							Value()

						Expect(v).To(Equal(userDefinedSigned(expect)))
					},
					Entry("zero", "0", 0),
					Entry("positive", "+123", +123),
					Entry("negative", "-123", -123),
				)
			})
		})

		When("the value is invalid", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it panics",
					func(value, expect string) {
						os.Setenv("FERRITE_SIGNED", value)

						Expect(func() {
							builder.
								Required().
								Value()
						}).To(PanicWith(expect))
					},
					Entry(
						"underflow",
						"-32769",
						`value of FERRITE_SIGNED (-32769) is invalid: strconv.ParseInt: parsing "-32769": value out of range`,
					),
					Entry(
						"overflow",
						"32768",
						`value of FERRITE_SIGNED (32768) is invalid: strconv.ParseInt: parsing "32768": value out of range`,
					),
					Entry(
						"decimal",
						"123.45",
						`value of FERRITE_SIGNED (123.45) is invalid: strconv.ParseInt: parsing "123.45": invalid syntax`,
					),
					Entry(
						"invalid characters",
						"123!",
						`value of FERRITE_SIGNED ('123!') is invalid: strconv.ParseInt: parsing "123!": invalid syntax`,
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

						Expect(v).To(Equal(userDefinedSigned(-123)))
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
							"FERRITE_SIGNED is undefined and does not have a default value",
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
						os.Setenv("FERRITE_SIGNED", value)

						v, ok := builder.
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(userDefinedSigned(expect)))
					},
					Entry("zero", "0", 0),
					Entry("positive", "+123", +123),
					Entry("negative", "-123", -123),
				)
			})
		})

		When("the value is invalid", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it panics",
					func(value, expect string) {
						os.Setenv("FERRITE_SIGNED", value)

						Expect(func() {
							builder.
								Optional().
								Value()
						}).To(PanicWith(expect))
					},
					Entry(
						"underflow",
						"-32769",
						`value of FERRITE_SIGNED (-32769) is invalid: strconv.ParseInt: parsing "-32769": value out of range`,
					),
					Entry(
						"overflow",
						"32768",
						`value of FERRITE_SIGNED (32768) is invalid: strconv.ParseInt: parsing "32768": value out of range`,
					),
					Entry(
						"decimal",
						"123.45",
						`value of FERRITE_SIGNED (123.45) is invalid: strconv.ParseInt: parsing "123.45": invalid syntax`,
					),
					Entry(
						"invalid characters",
						"123!",
						`value of FERRITE_SIGNED ('123!') is invalid: strconv.ParseInt: parsing "123!": invalid syntax`,
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
						Expect(v).To(Equal(userDefinedSigned(-123)))
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
})
