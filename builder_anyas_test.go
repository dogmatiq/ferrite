package ferrite_test

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// uppercased is a simple test type that wraps a string in uppercase form.
type uppercased struct {
	Value string
}

func parseUppercased(s string) (uppercased, error) {
	if s == "" {
		return uppercased{}, errors.New("must not be empty")
	}
	return uppercased{Value: strings.ToUpper(s)}, nil
}

func formatUppercased(v uppercased) (string, error) {
	if v.Value == "" {
		return "", errors.New("must not be empty")
	}
	return v.Value, nil
}

var _ = Describe("type AnyAsBuilder", func() {
	var builder *AnyAsBuilder[uppercased]

	BeforeEach(func() {
		builder = AnyAs("FERRITE_ANY_AS", "<desc>", parseUppercased, formatUppercased)
	})

	AfterEach(func() {
		tearDown()
	})

	It("panics if the name is empty", func() {
		Expect(func() {
			AnyAs("", "<desc>", parseUppercased, formatUppercased).Optional()
		}).To(PanicWith("invalid specification: variable name must not be empty"))
	})

	It("panics if the description is empty", func() {
		Expect(func() {
			AnyAs("FERRITE_ANY_AS", "", parseUppercased, formatUppercased).Optional()
		}).To(PanicWith("specification for FERRITE_ANY_AS is invalid: variable description must not be empty"))
	})

	When("the variable is required", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the parsed value", func() {
					os.Setenv("FERRITE_ANY_AS", "hello")

					v := builder.
						Required().
						Value()

					Expect(v).To(Equal(uppercased{Value: "HELLO"}))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							WithDefault(uppercased{Value: "DEFAULT"}).
							Required().
							Value()

						Expect(v).To(Equal(uppercased{Value: "DEFAULT"}))
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
							"FERRITE_ANY_AS is undefined and does not have a default value",
						))
					})
				})
			})
		})
	})

	When("the variable is optional", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the parsed value", func() {
					os.Setenv("FERRITE_ANY_AS", "hello")

					v, ok := builder.
						Optional().
						Value()

					Expect(ok).To(BeTrue())
					Expect(v).To(Equal(uppercased{Value: "HELLO"}))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							WithDefault(uppercased{Value: "DEFAULT"}).
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(uppercased{Value: "DEFAULT"}))
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

func ExampleAnyAs_required() {
	defer example()()

	v := ferrite.
		AnyAs("FERRITE_ANY_AS", "example parsed variable", parseUppercased, formatUppercased).
		Required()

	os.Setenv("FERRITE_ANY_AS", "hello")
	ferrite.Init()

	fmt.Println("value is", v.Value().Value)

	// Output:
	// value is HELLO
}

func ExampleAnyAs_default() {
	defer example()()

	v := ferrite.
		AnyAs("FERRITE_ANY_AS", "example parsed variable", parseUppercased, formatUppercased).
		WithDefault(uppercased{Value: "DEFAULT"}).
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value().Value)

	// Output:
	// value is DEFAULT
}

func ExampleAnyAs_optional() {
	defer example()()

	v := ferrite.
		AnyAs("FERRITE_ANY_AS", "example parsed variable", parseUppercased, formatUppercased).
		Optional()

	ferrite.Init()

	if x, ok := v.Value(); ok {
		fmt.Println("value is", x.Value)
	} else {
		fmt.Println("value is undefined")
	}

	// Output:
	// value is undefined
}

func ExampleAnyAs_constraint() {
	defer example()()

	ferrite.
		AnyAs("FERRITE_ANY_AS", "example constrained parsed variable", parseUppercased, formatUppercased).
		WithConstraint(
			"must contain X",
			func(v uppercased) bool {
				return strings.Contains(v.Value, "X")
			},
		).
		Required()

	os.Setenv("FERRITE_ANY_AS", "hello")
	ferrite.Init()

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_ANY_AS  example constrained parsed variable    <string>    ✗ set to hello, must contain X
	//
	// <process exited with error code 1>
}
