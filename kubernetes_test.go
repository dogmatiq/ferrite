package ferrite_test

import (
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type KubeServiceSpec", func() {
	var spec *KubeServiceSpec

	BeforeEach(func() {
		spec = KubeService("ferrite-svc", "<desc>")
	})

	AfterEach(func() {
		Teardown()
	})

	When("both the host and port environment variables are set", func() {
		BeforeEach(func() {
			os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
			os.Setenv("FERRITE_SVC_SERVICE_PORT", "12345")
		})

		Describe("func Address()", func() {
			It("returns the network address", func() {
				Expect(spec.Address()).To(Equal("host.example.org:12345"))
			})
		})

		When("the host is an IPv6 address", func() {
			BeforeEach(func() {
				os.Setenv("FERRITE_SVC_SERVICE_HOST", "::1")
			})

			Describe("func Address()", func() {
				It("properly escapes the IP address", func() {
					Expect(spec.Address()).To(Equal("[::1]:12345"))
				})
			})
		})
	})

	When("using a named port", func() {
		BeforeEach(func() {
			spec.WithNamedPort("named-port")

			os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
			os.Setenv("FERRITE_SVC_SERVICE_PORT_NAMED_PORT", "12345")
		})

		Describe("func Address()", func() {
			It("returns the network address", func() {
				Expect(spec.Address()).To(Equal("host.example.org:12345"))
			})
		})
	})

	When("the host environment variable is empty", func() {
		BeforeEach(func() {
			os.Setenv("FERRITE_SVC_SERVICE_HOST", "")
			os.Setenv("FERRITE_SVC_SERVICE_PORT", "12345")
		})

		When("there is a default value", func() {
			BeforeEach(func() {
				spec.WithDefault("default.example.org", "54321")
			})

			Describe("func Address()", func() {
				It("returns the default host with the explicit port", func() {
					Expect(spec.Address()).To(Equal("default.example.org:12345"))
				})
			})
		})

		When("there is no default value", func() {
			Describe("func Address()", func() {
				It("panics", func() {
					Expect(func() {
						spec.Address()
					}).To(PanicWith("FERRITE_SVC_SERVICE_HOST is invalid: must not be empty"))
				})
			})
		})
	})

	When("the port environment variable is empty", func() {
		BeforeEach(func() {
			os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
			os.Setenv("FERRITE_SVC_SERVICE_PORT", "")
		})

		When("there is a default value", func() {
			BeforeEach(func() {
				spec.WithDefault("default.example.org", "54321")
			})

			Describe("func Address()", func() {
				It("returns the explicit host with the default port", func() {
					Expect(spec.Address()).To(Equal("host.example.org:54321"))
				})
			})
		})

		When("there is no default value", func() {
			Describe("func Address()", func() {
				It("panics", func() {
					Expect(func() {
						spec.Address()
					}).To(PanicWith("FERRITE_SVC_SERVICE_PORT is invalid: must not be empty"))
				})
			})
		})
	})

	When("both environment variables are empty", func() {
		When("there is a default value", func() {
			BeforeEach(func() {
				spec.WithDefault("default.example.org", "54321")
			})

			Describe("func Address()", func() {
				It("returns the default", func() {
					Expect(spec.Address()).To(Equal("default.example.org:54321"))
				})
			})
		})

		When("there is no default value", func() {
			Describe("func Address()", func() {
				It("panics", func() {
					Expect(func() {
						spec.Address()
					}).To(PanicWith("FERRITE_SVC_SERVICE_HOST is invalid: must not be empty"))
				})
			})
		})
	})
})
