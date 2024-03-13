package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type enumMember int

const (
	member0 enumMember = iota
	member1
	member2
)

// String returns the value to use as the enum's key.
func (m enumMember) String() string {
	return fmt.Sprintf("<member-%d>", m)
}

var _ = Describe("type EnumBuilder", func() {
	var builder *EnumBuilder[enumMember]

	BeforeEach(func() {
		builder = EnumAs[enumMember]("FERRITE_ENUM", "<desc>").
			WithMember(member0, "<desc-of-0>").
			WithMember(member1, "<desc-of-1>").
			WithMember(member2, "<desc-of-2>")
	})

	AfterEach(func() {
		tearDown()
	})

	It("panics if the name is empty", func() {
		Expect(func() {
			EnumAs[enumMember]("", "<desc>").Optional()
		}).To(PanicWith("invalid specification: variable name must not be empty"))
	})

	It("panics if the description is empty", func() {
		Expect(func() {
			EnumAs[enumMember]("FERRITE_ENUM", "").Optional()
		}).To(PanicWith("specification for FERRITE_ENUM is invalid: variable description must not be empty"))
	})

	When("the variable is required", func() {
		When("the value is one of the accepted literals", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it returns the value associated with the literal",
					func(value string, expect enumMember) {
						os.Setenv("FERRITE_ENUM", value)

						v := builder.
							Required().
							Value()

						Expect(v).To(Equal(expect))
					},
					Entry("member 0", "<member-0>", member0),
					Entry("member 1", "<member-1>", member1),
					Entry("member 2", "<member-2>", member2),
				)
			})
		})

		When("the value is invalid", func() {
			BeforeEach(func() {
				os.Setenv("FERRITE_ENUM", "<non-member>")
			})

			Describe("func Value()", func() {
				It("panics", func() {
					Expect(func() {
						builder.
							Required().
							Value()
					}).To(PanicWith(
						`value of FERRITE_ENUM ('<non-member>') is invalid: expected '<member-0>', '<member-1>' or '<member-2>'`,
					))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							WithDefault(member1).
							Required().
							Value()

						Expect(v).To(Equal(member1))
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
							"FERRITE_ENUM is undefined and does not have a default value",
						))
					})
				})
			})
		})
	})

	When("the variable is optional", func() {
		When("the value is one of the accepted literals", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it returns the value associated with the literal",
					func(value string, expect enumMember) {
						os.Setenv("FERRITE_ENUM", value)

						v, ok := builder.
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(expect))
					},
					Entry("member 0", "<member-0>", member0),
					Entry("member 1", "<member-1>", member1),
					Entry("member 2", "<member-2>", member2),
				)
			})
		})

		When("the value is invalid", func() {
			BeforeEach(func() {
				os.Setenv("FERRITE_ENUM", "<non-member>")
			})

			Describe("func Value()", func() {
				It("panics", func() {
					Expect(func() {
						builder.
							Optional().
							Value()
					}).To(PanicWith(
						`value of FERRITE_ENUM ('<non-member>') is invalid: expected '<member-0>', '<member-1>' or '<member-2>'`,
					))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							WithDefault(member1).
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(member1))
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

	When("multiple members have the same literal representation", func() {
		It("panics", func() {
			Expect(func() {
				builder.
					WithMember(member1, "").
					WithMember(member1, "").
					Required()
			}).To(PanicWith(
				`specification for FERRITE_ENUM is invalid: literals must be unique but multiple values are represented as "<member-1>"`,
			))
		})
	})

	When("there are no members", func() {
		It("panics", func() {
			Expect(func() {
				ferrite.
					EnumAs[string]("FERRITE_ENUM", "<desc>").
					Required()
			}).To(PanicWith(
				`specification for FERRITE_ENUM is invalid: must allow at least two distinct values`,
			))
		})
	})

	When("the default value is not a member of the enum", func() {
		It("panics", func() {
			Expect(func() {
				builder.
					WithDefault(enumMember(100)).
					Required()
			}).To(PanicWith(
				`specification for FERRITE_ENUM is invalid: default value: expected '<member-0>', '<member-1>' or '<member-2>'`,
			))
		})
	})
})

func ExampleEnum_required() {
	defer example()()

	v := ferrite.
		Enum("FERRITE_ENUM", "example enum variable").
		WithMembers("red", "green", "blue").
		Required()

	os.Setenv("FERRITE_ENUM", "red")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is red
}

func ExampleEnum_default() {
	defer example()()

	v := ferrite.
		Enum("FERRITE_ENUM", "example enum variable").
		WithMembers("red", "green", "blue").
		WithDefault("green").
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is green
}

func ExampleEnum_optional() {
	defer example()()

	v := ferrite.
		Enum("FERRITE_ENUM", "example enum variable").
		WithMembers("red", "green", "blue").
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

func ExampleEnum_descriptions() {
	defer example()()

	v := ferrite.
		Enum("FERRITE_ENUM", "example enum variable").
		WithMember("red", "the color red").
		WithMember("green", "the color green").
		WithMember("blue", "the color blue").
		Required()

	os.Setenv("FERRITE_ENUM", "red")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is red
}

func ExampleEnum_deprecated() {
	defer example()()

	v := ferrite.
		Enum("FERRITE_ENUM", "example enum variable").
		WithMember("red", "the color red").
		WithMember("green", "the color green").
		WithMember("blue", "the color blue").
		Deprecated()

	os.Setenv("FERRITE_ENUM", "red")
	ferrite.Init()

	if x, ok := v.DeprecatedValue(); ok {
		fmt.Println("value is", x)
	} else {
		fmt.Println("value is undefined")
	}

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_ENUM  example enum variable  [ red | green | blue ]  ⚠ deprecated variable set to red
	//
	// value is red
}
