package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type userDefinedSigned int16

var _ = Describe("type SignedBuilder", func() {
	var builder *SignedBuilder[userDefinedSigned]

	BeforeEach(func() {
		builder = Signed[userDefinedSigned]("FERRITE_SIGNED", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	It("panics if the name is empty", func() {
		Expect(func() {
			Signed[userDefinedSigned]("", "<desc>").Optional()
		}).To(PanicWith("invalid specification: variable name must not be empty"))
	})

	It("panics if the description is empty", func() {
		Expect(func() {
			Signed[userDefinedSigned]("FERRITE_SIGNED", "").Optional()
		}).To(PanicWith("specification for FERRITE_SIGNED is invalid: variable description must not be empty"))
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
						`value of FERRITE_SIGNED (-32769) is invalid: too low, expected the smallest int16 value of -32768 or greater`,
					),
					Entry(
						"overflow",
						"32768",
						`value of FERRITE_SIGNED (32768) is invalid: too high, expected the largest int16 value of +32767 or less`,
					),
					Entry(
						"decimal",
						"123.45",
						`value of FERRITE_SIGNED (123.45) is invalid: unrecognized int16 syntax`,
					),
					Entry(
						"invalid characters",
						"123!",
						`value of FERRITE_SIGNED ('123!') is invalid: unrecognized int16 syntax`,
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
						`value of FERRITE_SIGNED (-32769) is invalid: too low, expected the smallest int16 value of -32768 or greater`,
					),
					Entry(
						"overflow",
						"32768",
						`value of FERRITE_SIGNED (32768) is invalid: too high, expected the largest int16 value of +32767 or less`,
					),
					Entry(
						"decimal",
						"123.45",
						`value of FERRITE_SIGNED (123.45) is invalid: unrecognized int16 syntax`,
					),
					Entry(
						"invalid characters",
						"123!",
						`value of FERRITE_SIGNED ('123!') is invalid: unrecognized int16 syntax`,
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

	When("the value is lower than the minimum limit", func() {
		It("panics", func() {
			Expect(func() {
				os.Setenv("FERRITE_SIGNED", "1")

				builder.
					WithMinimum(5).
					Required().
					Value()
			}).To(PanicWith(
				`value of FERRITE_SIGNED (1) is invalid: too low, expected +5 or greater`,
			))
		})
	})

	When("the value is greater than the maximum limit", func() {
		It("panics", func() {
			Expect(func() {
				os.Setenv("FERRITE_SIGNED", "10")

				builder.
					WithMaximum(5).
					Required().
					Value()
			}).To(PanicWith(
				`value of FERRITE_SIGNED (10) is invalid: too high, expected +5 or less`,
			))
		})
	})
})

func ExampleSigned_required() {
	defer example()()

	v := ferrite.
		Signed[int]("FERRITE_SIGNED", "example signed integer variable").
		Required()

	os.Setenv("FERRITE_SIGNED", "-123")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is -123
}

func ExampleSigned_default() {
	defer example()()

	v := ferrite.
		Signed[int]("FERRITE_SIGNED", "example signed integer variable").
		WithDefault(-123).
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is -123
}

func ExampleSigned_optional() {
	defer example()()

	v := ferrite.
		Signed[int]("FERRITE_SIGNED", "example signed integer variable").
		Optional()

	ferrite.Init()

	if x, ok := v.Value(); ok {
		fmt.Println("value is", x)
	} else {
		fmt.Println("value is undefined")
	}

	// Output:
	// value is undefined
}

func ExampleSigned_limits() {
	defer example()()

	v := ferrite.
		Signed[int]("FERRITE_SIGNED", "example signed integer variable").
		WithMinimum(-5).
		WithMaximum(10).
		Required()

	os.Setenv("FERRITE_SIGNED", "-2")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is -2
}

func ExampleSigned_deprecated() {
	defer example()()

	os.Setenv("FERRITE_SIGNED", "-123")
	v := ferrite.
		Signed[int]("FERRITE_SIGNED", "example signed integer variable").
		Deprecated()

	ferrite.Init()

	if x, ok := v.DeprecatedValue(); ok {
		fmt.Println("value is", x)
	} else {
		fmt.Println("value is undefined")
	}

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_SIGNED  example signed integer variable  [ <int> ]  ⚠ deprecated variable set to -123
	//
	// value is -123
}
