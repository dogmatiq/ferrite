package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("func NetworkAddress", func() {
	var builder *NetworkAddressBuilder

	BeforeEach(func() {
		builder = NetworkAddress("FERRITE_NETWORK_ADDR", "<desc>")
	})

	AfterEach(func() {
		tearDown()
	})

	It("panics if the name is empty", func() {
		Expect(func() {
			NetworkAddress("", "<desc>").Optional()
		}).To(PanicWith("invalid specification: variable name must not be empty"))
	})

	It("panics if the description is empty", func() {
		Expect(func() {
			NetworkAddress("FERRITE_NETWORK_ADDR", "").Optional()
		}).To(PanicWith("specification for FERRITE_NETWORK_ADDR is invalid: variable description must not be empty"))
	})

	When("the variable is required", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				DescribeTable(
					"it returns the parsed address",
					func(input string, expected NetworkAddr) {
						os.Setenv("FERRITE_NETWORK_ADDR", input)

						v := builder.
							Required().
							Value()

						Expect(v).To(Equal(expected))
					},
					Entry(
						"IPv4 address",
						"192.168.0.1:8080",
						NetworkAddr{Host: "192.168.0.1", Port: "8080"},
					),
					Entry(
						"IPv6 address",
						"[::1]:8080",
						NetworkAddr{Host: "::1", Port: "8080"},
					),
					Entry(
						"hostname with numeric port",
						"host.example.org:8080",
						NetworkAddr{Host: "host.example.org", Port: "8080"},
					),
					Entry(
						"hostname with IANA service name",
						"host.example.org:https",
						NetworkAddr{Host: "host.example.org", Port: "https"},
					),
				)

				It("rejects an invalid address", func() {
					os.Setenv("FERRITE_NETWORK_ADDR", "not-an-address")

					Expect(func() {
						builder.
							Required().
							Value()
					}).To(Panic())
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							WithDefault("localhost:8080").
							Required().
							Value()

						Expect(v).To(Equal(NetworkAddr{Host: "localhost", Port: "8080"}))
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
							"FERRITE_NETWORK_ADDR is undefined and does not have a default value",
						))
					})
				})
			})
		})
	})

	When("the variable is optional", func() {
		When("the value is not empty", func() {
			Describe("func Value()", func() {
				It("returns the value", func() {
					os.Setenv("FERRITE_NETWORK_ADDR", "192.168.0.1:8080")

					v, ok := builder.
						Optional().
						Value()

					Expect(ok).To(BeTrue())
					Expect(v).To(Equal(NetworkAddr{Host: "192.168.0.1", Port: "8080"}))
				})
			})
		})

		When("the value is empty", func() {
			When("there is a default value", func() {
				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							WithDefault("localhost:8080").
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(NetworkAddr{Host: "localhost", Port: "8080"}))
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

	Describe("func NetworkAddr.String()", func() {
		It("returns the address in host:port form for IPv4", func() {
			addr := NetworkAddr{Host: "192.168.0.1", Port: "8080"}
			Expect(addr.String()).To(Equal("192.168.0.1:8080"))
		})

		It("returns the address in [host]:port form for IPv6", func() {
			addr := NetworkAddr{Host: "::1", Port: "8080"}
			Expect(addr.String()).To(Equal("[::1]:8080"))
		})
	})
})

func ExampleNetworkAddress_required() {
	defer example()()

	v := ferrite.
		NetworkAddress("FERRITE_NETWORK_ADDR", "example network address variable").
		Required()

	os.Setenv("FERRITE_NETWORK_ADDR", "host.example.org:https")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is host.example.org:https
}

func ExampleNetworkAddress_default() {
	defer example()()

	v := ferrite.
		NetworkAddress("FERRITE_NETWORK_ADDR", "example network address variable").
		WithDefault("localhost:8080").
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is localhost:8080
}

func ExampleNetworkAddress_optional() {
	defer example()()

	v := ferrite.
		NetworkAddress("FERRITE_NETWORK_ADDR", "example network address variable").
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

func ExampleNetworkAddress_deprecated() {
	defer example()()

	v := ferrite.
		NetworkAddress("FERRITE_NETWORK_ADDR", "example network address variable").
		Deprecated()

	os.Setenv("FERRITE_NETWORK_ADDR", "host.example.org:https")
	ferrite.Init()

	if x, ok := v.DeprecatedValue(); ok {
		fmt.Println("value is", x)
	} else {
		fmt.Println("value is undefined")
	}

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_NETWORK_ADDR  example network address variable  [ <string> ]  ⚠ deprecated variable set to host.example.org:https
	//
	// value is host.example.org:https
}
