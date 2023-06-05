package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type userDefinedByte byte
type userDefinedBinary []userDefinedByte

var _ = Describe("type BinaryBuilder", func() {
	var builder *BinaryBuilder[userDefinedBinary, userDefinedByte]

	BeforeEach(func() {
		builder = BinaryAs[userDefinedBinary]("FERRITE_BINARY", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	It("panics if the name is empty", func() {
		Expect(func() {
			BinaryAs[userDefinedBinary]("", "<desc>").Optional()
		}).To(PanicWith("invalid specification: variable name must not be empty"))
	})

	It("panics if the description is empty", func() {
		Expect(func() {
			BinaryAs[userDefinedBinary]("FERRITE_BINARY", "").Optional()
		}).To(PanicWith("specification for FERRITE_BINARY is invalid: variable description must not be empty"))
	})

	When("the variable is required", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the value ", func() {
					os.Setenv("FERRITE_BINARY", "PHZhbHVlPg==")

					v := builder.
						Required().
						Value()

					Expect(v).To(Equal(userDefinedBinary("<value>")))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							WithDefault(userDefinedBinary("<value>")).
							Required().
							Value()

						Expect(v).To(Equal(userDefinedBinary("<value>")))
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
							"FERRITE_BINARY is undefined and does not have a default value",
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
					os.Setenv("FERRITE_BINARY", "PHZhbHVlPg==")

					v, ok := builder.
						Optional().
						Value()

					Expect(ok).To(BeTrue())
					Expect(v).To(Equal(userDefinedBinary("<value>")))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							WithDefault(userDefinedBinary("<value>")).
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(userDefinedBinary("<value>")))
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

func ExampleBinary_required() {
	defer example()()

	v := ferrite.
		Binary("FERRITE_BINARY", "example binary variable").
		Required()

	os.Setenv("FERRITE_BINARY", "PHZhbHVlPg==")
	ferrite.Init()

	fmt.Println("value is", string(v.Value()))

	// Output:
	// value is <value>
}

func ExampleBinary_default() {
	defer example()()

	v := ferrite.
		Binary("FERRITE_BINARY", "example binary variable").
		WithDefault([]byte("<default>")).
		Required()

	ferrite.Init()

	fmt.Println("value is", string(v.Value()))

	// Output:
	// value is <default>
}

func ExampleBinary_optional() {
	defer example()()

	v := ferrite.
		Binary("FERRITE_BINARY", "example binary variable").
		Optional()

	ferrite.Init()

	if x, ok := v.Value(); ok {
		fmt.Println("value is", string(x))
	} else {
		fmt.Println("value is undefined")
	}

	// Output:
	// value is undefined
}

func ExampleBinary_sensitive() {
	defer example()()

	os.Setenv("FERRITE_BINARY", "aHVudGVyMg==")
	ferrite.
		Binary("FERRITE_BINARY", "example sensitive binary variable").
		WithConstraint(
			"always fail",
			func(v []byte) bool {
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
	//  ❯ FERRITE_BINARY  example sensitive binary variable    <base64>    ✗ set to 12-byte value, always fail
	//
	// <process exited with error code 1>
}

func ExampleBinary_deprecated() {
	defer example()()

	os.Setenv("FERRITE_BINARY", "PHZhbHVlPg==")
	v := ferrite.
		Binary("FERRITE_BINARY", "example binary variable").
		Deprecated()

	ferrite.Init()

	if x, ok := v.DeprecatedValue(); ok {
		fmt.Println("value is", string(x))
	} else {
		fmt.Println("value is undefined")
	}

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_BINARY  example binary variable  [ <base64> ]  ⚠ deprecated variable set to 12-byte value
	//
	// value is <value>
}