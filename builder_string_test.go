package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type userDefinedString string

var _ = Describe("type StringBuilder", func() {
	var builder *StringBuilder[userDefinedString]

	BeforeEach(func() {
		builder = StringAs[userDefinedString]("FERRITE_STRING", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	It("panics if the name is empty", func() {
		Expect(func() {
			StringAs[userDefinedString]("", "<desc>").Optional()
		}).To(PanicWith("invalid specification: variable name must not be empty"))
	})

	It("panics if the description is empty", func() {
		Expect(func() {
			StringAs[userDefinedString]("FERRITE_STRING", "").Optional()
		}).To(PanicWith("specification for FERRITE_STRING is invalid: variable description must not be empty"))
	})

	When("the variable is required", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the value ", func() {
					os.Setenv("FERRITE_STRING", "<value>")

					v := builder.
						Required().
						Value()

					Expect(v).To(Equal(userDefinedString("<value>")))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							WithDefault("<value>").
							Required().
							Value()

						Expect(v).To(Equal(userDefinedString("<value>")))
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
							"FERRITE_STRING is undefined and does not have a default value",
						))
					})
				})
			})
		})
	})

	When("the variable is optional", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the value ", func() {
					os.Setenv("FERRITE_STRING", "<value>")

					v, ok := builder.
						Optional().
						Value()

					Expect(ok).To(BeTrue())
					Expect(v).To(Equal(userDefinedString("<value>")))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							WithDefault("<value>").
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(userDefinedString("<value>")))
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

func ExampleString_required() {
	defer example()()

	v := ferrite.
		String("FERRITE_STRING", "example string variable").
		Required()

	os.Setenv("FERRITE_STRING", "<value>")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is <value>
}

func ExampleString_default() {
	defer example()()

	v := ferrite.
		String("FERRITE_STRING", "example string variable").
		WithDefault("<default>").
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is <default>
}

func ExampleString_optional() {
	defer example()()

	v := ferrite.
		String("FERRITE_STRING", "example string variable").
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

func ExampleString_sensitive() {
	defer example()()

	os.Setenv("FERRITE_STRING", "hunter2")
	ferrite.
		String("FERRITE_STRING", "example sensitive string variable").
		WithConstraint(
			"always fail",
			func(s string) bool {
				// Force the variable to be considered invalid so that the
				// variable table is rendered to the console.
				return false
			},
		).
		WithSensitiveContent().
		Required()

	ferrite.Init()

	// Note that the variable's value is obscured in the console output.

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_STRING  example sensitive string variable    <string>    ✗ set to *******, always fail
	//
	// <process exited with error code 1>
}

func ExampleString_deprecated() {
	defer example()()

	os.Setenv("FERRITE_STRING", "<value>")
	v := ferrite.
		String("FERRITE_STRING", "example string variable").
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
	//  ❯ FERRITE_STRING  example string variable  [ <string> ]  ⚠ deprecated variable set to '<value>'
	//
	// value is <value>
}
