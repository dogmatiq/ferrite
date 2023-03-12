package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type userDefinedBool bool

func (v userDefinedBool) String() string {
	if v {
		return "yes"
	}
	return "no"
}

var _ = Describe("type BoolBuilder", func() {
	var builder *BoolBuilder[userDefinedBool]

	BeforeEach(func() {
		builder = BoolAs[userDefinedBool]("FERRITE_BOOL", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	It("panics if the name is empty", func() {
		Expect(func() {
			BoolAs[userDefinedBool]("", "<desc>").Optional()
		}).To(PanicWith("invalid specification: variable name must not be empty"))
	})

	It("panics if the description is empty", func() {
		Expect(func() {
			BoolAs[userDefinedBool]("FERRITE_BOOL", "").Optional()
		}).To(PanicWith("specification for FERRITE_BOOL is invalid: variable description must not be empty"))
	})

	When("the variable is required", func() {
		When("the value is one of the accepted literals", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it returns the value associated with the literal",
					func(value string, expect userDefinedBool) {
						os.Setenv("FERRITE_BOOL", value)

						v := builder.
							Required().
							Value()

						Expect(v).To(Equal(expect))
					},
					Entry("true", "yes", userDefinedBool(true)),
					Entry("false", "no", userDefinedBool(false)),
				)
			})
		})

		When("the value is invalid", func() {
			BeforeEach(func() {
				// we don't accept true/false for the userDefinedBool type
				os.Setenv("FERRITE_BOOL", "true")
			})

			Describe("func Value()", func() {
				It("panics", func() {
					Expect(func() {
						builder.
							Required().
							Value()
					}).To(PanicWith(
						`value of FERRITE_BOOL (true) is invalid: expected either yes or no`,
					))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					DescribeTable(
						"it returns the default",
						func(expect userDefinedBool) {
							v := builder.
								WithDefault(expect).
								Required().
								Value()

							Expect(v).To(Equal(expect))
						},
						Entry("true", userDefinedBool(true)),
						Entry("false", userDefinedBool(false)),
					)
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
							"FERRITE_BOOL is undefined and does not have a default value",
						))
					})
				})
			})
		})
	})

	When("the variable is optional", func() {
		When("the value is invalid", func() {
			BeforeEach(func() {
				// we don't accept true/false for the userDefinedBool type
				os.Setenv("FERRITE_BOOL", "true")
			})

			Describe("func Value()", func() {
				It("panics", func() {
					Expect(func() {
						builder.
							Optional().
							Value()
					}).To(PanicWith(
						`value of FERRITE_BOOL (true) is invalid: expected either yes or no`,
					))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					DescribeTable(
						"it returns the default",
						func(expect userDefinedBool) {
							v, ok := builder.
								WithDefault(expect).
								Optional().
								Value()

							Expect(ok).To(BeTrue())
							Expect(v).To(Equal(expect))
						},
						Entry("true", userDefinedBool(true)),
						Entry("false", userDefinedBool(false)),
					)
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

	Describe("func WithLiterals()", func() {
		DescribeTable(
			"it overrides the default literals",
			func(value string, expect userDefinedBool) {
				os.Setenv("FERRITE_BOOL", value)

				v := builder.
					WithLiterals("true", "false").
					Required().
					Value()

				Expect(v).To(Equal(expect))
			},
			Entry("true", "true", userDefinedBool(true)),
			Entry("false", "false", userDefinedBool(false)),
		)

		When("the true literal is empty", func() {
			It("panics", func() {
				Expect(func() {
					builder.
						WithLiterals("", "no").
						Optional()
				}).To(PanicWith(
					"specification for FERRITE_BOOL is invalid: literals can not be an empty string",
				))
			})
		})

		When("the true literal is empty", func() {
			It("panics", func() {
				Expect(func() {
					builder.
						WithLiterals("yes", "").
						Optional()

				}).To(PanicWith(
					"specification for FERRITE_BOOL is invalid: literals can not be an empty string",
				))
			})
		})
	})
})

func ExampleBool_required() {
	defer example()()

	v := ferrite.
		Bool("FERRITE_BOOL", "example boolean variable").
		Required()

	os.Setenv("FERRITE_BOOL", "true")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is true
}

func ExampleBool_default() {
	defer example()()

	v := ferrite.
		Bool("FERRITE_BOOL", "example boolean variable").
		WithDefault(true).
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is true
}

func ExampleBool_optional() {
	defer example()()

	v := ferrite.
		Bool("FERRITE_BOOL", "example boolean variable").
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

func ExampleBool_customLiterals() {
	defer example()()

	v := ferrite.
		Bool("FERRITE_BOOL", "example boolean variable").
		WithLiterals("yes", "no").
		Required()

	os.Setenv("FERRITE_BOOL", "yes")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is true
}

func ExampleBool_deprecated() {
	defer example()()

	os.Setenv("FERRITE_BOOL", "true")
	v := ferrite.
		Bool("FERRITE_BOOL", "example boolean variable").
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
	//  ❯ FERRITE_BOOL  example boolean variable  [ true | false ]  ⚠ deprecated variable set to true
	//
	// value is true
}
