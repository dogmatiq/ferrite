package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("func NetworkPort", func() {
	var builder *NetworkPortBuilder

	BeforeEach(func() {
		builder = NetworkPort("FERRITE_NETWORK_PORT", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	It("panics if the name is empty", func() {
		Expect(func() {
			NetworkPort("", "<desc>").Optional()
		}).To(PanicWith("invalid specification: variable name must not be empty"))
	})

	It("panics if the description is empty", func() {
		Expect(func() {
			NetworkPort("FERRITE_NETWORK_PORT", "").Optional()
		}).To(PanicWith("specification for FERRITE_NETWORK_PORT is invalid: variable description must not be empty"))
	})

	When("the variable is required", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the value ", func() {
					os.Setenv("FERRITE_NETWORK_PORT", "12345")

					v := builder.
						Required().
						Value()

					Expect(v).To(Equal("12345"))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							WithDefault("https").
							Required().
							Value()

						Expect(v).To(Equal("https"))
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
							"FERRITE_NETWORK_PORT is undefined and does not have a default value",
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
					os.Setenv("FERRITE_NETWORK_PORT", "12345")

					v, ok := builder.
						Optional().
						Value()

					Expect(ok).To(BeTrue())
					Expect(v).To(Equal("12345"))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							WithDefault("12345").
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal("12345"))
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

func ExampleNetworkPort_required() {
	defer example()()

	v := ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port variable").
		Required()

	os.Setenv("FERRITE_NETWORK_PORT", "https")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is https
}

func ExampleNetworkPort_default() {
	defer example()()

	v := ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port variable").
		WithDefault("12345").
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is 12345
}

func ExampleNetworkPort_optional() {
	defer example()()

	v := ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port variable").
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

func ExampleNetworkPort_deprecated() {
	defer example()()

	os.Setenv("FERRITE_NETWORK_PORT", "https")
	v := ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port variable").
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
	//  ❯ FERRITE_NETWORK_PORT  example network port variable  [ <string> ]  ⚠ deprecated variable set to https
	//
	// value is https
}
