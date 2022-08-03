package ferrite_test

import (
	"errors"
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type KubeServiceSpec", func() {
	var spec *KubeServiceSpec

	BeforeEach(func() {
		// FERRITE_SVC_SERVICE_HOST=host.example.org
		// FERRITE_SVC_SERVICE_PORT=12345
		// FERRITE_SVC_SERVICE_PORT_NAMED_PORT=12345
		// FERRITE_SVC_PORT=tcp://host.example.org:12345

		// FERRITE_SVC_PORT_12345_TCP=tcp://host.example.org:12345
		// FERRITE_SVC_PORT_12345_TCP_PROTO=tcp
		// FERRITE_SVC_PORT_12345_TCP_PORT=12345
		// FERRITE_SVC_PORT_12345_TCP_ADDR=host.example.org

		spec = KubeService("ferrite-svc")
	})

	AfterEach(func() {
		tearDown()
	})

	When("both the host and port environment variables are set", func() {
		BeforeEach(func() {
			os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
			os.Setenv("FERRITE_SVC_SERVICE_PORT", "12345")
		})

		Describe("func Address(), Host() and Port()", func() {
			It("return the network address", func() {
				Expect(spec.Address()).To(Equal("host.example.org:12345"))
				Expect(spec.Host()).To(Equal("host.example.org"))
				Expect(spec.Port()).To(Equal("12345"))
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

		Describe("func Validate()", func() {
			It("returns success results", func() {
				Expect(spec.Validate()).To(ConsistOf(
					ValidationResult{
						Name:          "FERRITE_SVC_SERVICE_HOST",
						Description:   `Hostname or IP address of the "ferrite-svc" service.`,
						ValidInput:    "[string]",
						DefaultValue:  "",
						ExplicitValue: "host.example.org",
						Error:         nil,
					},
					ValidationResult{
						Name:          "FERRITE_SVC_SERVICE_PORT",
						Description:   `Network port of the "ferrite-svc" service.`,
						ValidInput:    "[string]|(1..65535)",
						DefaultValue:  "",
						ExplicitValue: "12345",
						Error:         nil,
					},
				))
			})
		})
	})

	When("using a named port", func() {
		BeforeEach(func() {
			spec.WithNamedPort("named-port")

			os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
			os.Setenv("FERRITE_SVC_SERVICE_PORT_NAMED_PORT", "12345")
		})

		Describe("func Address(), Host() and Port()", func() {
			It("return the network address", func() {
				Expect(spec.Address()).To(Equal("host.example.org:12345"))
				Expect(spec.Host()).To(Equal("host.example.org"))
				Expect(spec.Port()).To(Equal("12345"))
			})
		})

		Describe("func Validate()", func() {
			It("returns success results", func() {
				Expect(spec.Validate()).To(ConsistOf(
					ValidationResult{
						Name:          "FERRITE_SVC_SERVICE_HOST",
						Description:   `Hostname or IP address of the "ferrite-svc" service.`,
						ValidInput:    "[string]",
						DefaultValue:  "",
						ExplicitValue: "host.example.org",
						Error:         nil,
					},
					ValidationResult{
						Name:          "FERRITE_SVC_SERVICE_PORT_NAMED_PORT",
						Description:   `Network port of the "ferrite-svc" service's "named-port" port.`,
						ValidInput:    "[string]|(1..65535)",
						DefaultValue:  "",
						ExplicitValue: "12345",
						Error:         nil,
					},
				))
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

			Describe("func Address(), Host() and Port()", func() {
				It("return the network address with the default host", func() {
					Expect(spec.Address()).To(Equal("default.example.org:12345"))
					Expect(spec.Host()).To(Equal("default.example.org"))
					Expect(spec.Port()).To(Equal("12345"))
				})
			})

			Describe("func Validate()", func() {
				It("returns sucess results", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:         "FERRITE_SVC_SERVICE_HOST",
							Description:  `Hostname or IP address of the "ferrite-svc" service.`,
							ValidInput:   "[string]",
							DefaultValue: "default.example.org",
							UsingDefault: true,
						},
						ValidationResult{
							Name:          "FERRITE_SVC_SERVICE_PORT",
							Description:   `Network port of the "ferrite-svc" service.`,
							ValidInput:    "[string]|(1..65535)",
							DefaultValue:  "54321",
							ExplicitValue: "12345",
						},
					))
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

			Describe("func Host()", func() {
				It("panics", func() {
					Expect(func() {
						spec.Host()
					}).To(PanicWith("FERRITE_SVC_SERVICE_HOST is invalid: must not be empty"))
				})
			})

			Describe("func Validate()", func() {
				It("returns a failure result for the host", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:        "FERRITE_SVC_SERVICE_HOST",
							Description: `Hostname or IP address of the "ferrite-svc" service.`,
							ValidInput:  "[string]",
							Error:       errors.New(`must not be empty`),
						},
						ValidationResult{
							Name:          "FERRITE_SVC_SERVICE_PORT",
							Description:   `Network port of the "ferrite-svc" service.`,
							ValidInput:    "[string]|(1..65535)",
							ExplicitValue: "12345",
						},
					))
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

			Describe("func Address(), Host() and Port()", func() {
				It("return the network address with the default port", func() {
					Expect(spec.Address()).To(Equal("host.example.org:54321"))
					Expect(spec.Host()).To(Equal("host.example.org"))
					Expect(spec.Port()).To(Equal("54321"))
				})
			})

			Describe("func Validate()", func() {
				It("returns sucess results", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:          "FERRITE_SVC_SERVICE_HOST",
							Description:   `Hostname or IP address of the "ferrite-svc" service.`,
							ValidInput:    "[string]",
							DefaultValue:  "default.example.org",
							ExplicitValue: "host.example.org",
						},
						ValidationResult{
							Name:         "FERRITE_SVC_SERVICE_PORT",
							Description:  `Network port of the "ferrite-svc" service.`,
							ValidInput:   "[string]|(1..65535)",
							DefaultValue: "54321",
							UsingDefault: true,
						},
					))
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

			Describe("func Port()", func() {
				It("panics", func() {
					Expect(func() {
						spec.Port()
					}).To(PanicWith("FERRITE_SVC_SERVICE_PORT is invalid: must not be empty"))
				})
			})

			Describe("func Validate()", func() {
				It("returns a failure result for the port", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:          "FERRITE_SVC_SERVICE_HOST",
							Description:   `Hostname or IP address of the "ferrite-svc" service.`,
							ValidInput:    "[string]",
							ExplicitValue: "host.example.org",
						},
						ValidationResult{
							Name:        "FERRITE_SVC_SERVICE_PORT",
							Description: `Network port of the "ferrite-svc" service.`,
							ValidInput:  "[string]|(1..65535)",
							Error:       errors.New(`must not be empty`),
						},
					))
				})
			})
		})
	})

	When("both environment variables are empty", func() {
		When("there is a default value", func() {
			BeforeEach(func() {
				spec.WithDefault("default.example.org", "54321")
			})

			Describe("func Address(), Host() and Port()", func() {
				It("return the default network address", func() {
					Expect(spec.Address()).To(Equal("default.example.org:54321"))
					Expect(spec.Host()).To(Equal("default.example.org"))
					Expect(spec.Port()).To(Equal("54321"))
				})
			})

			Describe("func Validate()", func() {
				It("returns success results", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:         "FERRITE_SVC_SERVICE_HOST",
							Description:  `Hostname or IP address of the "ferrite-svc" service.`,
							ValidInput:   "[string]",
							DefaultValue: "default.example.org",
							UsingDefault: true,
						},
						ValidationResult{
							Name:         "FERRITE_SVC_SERVICE_PORT",
							Description:  `Network port of the "ferrite-svc" service.`,
							ValidInput:   "[string]|(1..65535)",
							DefaultValue: "54321",
							UsingDefault: true,
						},
					))
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

			Describe("func Host()", func() {
				It("panics", func() {
					Expect(func() {
						spec.Host()
					}).To(PanicWith("FERRITE_SVC_SERVICE_HOST is invalid: must not be empty"))
				})
			})

			Describe("func Port()", func() {
				It("panics", func() {
					Expect(func() {
						spec.Port()
					}).To(PanicWith("FERRITE_SVC_SERVICE_PORT is invalid: must not be empty"))
				})
			})

			Describe("func Validate()", func() {
				It("returns failure results", func() {
					Expect(spec.Validate()).To(ConsistOf(
						ValidationResult{
							Name:        "FERRITE_SVC_SERVICE_HOST",
							Description: `Hostname or IP address of the "ferrite-svc" service.`,
							ValidInput:  "[string]",
							Error:       errors.New(`must not be empty`),
						},
						ValidationResult{
							Name:        "FERRITE_SVC_SERVICE_PORT",
							Description: `Network port of the "ferrite-svc" service.`,
							ValidInput:  "[string]|(1..65535)",
							Error:       errors.New(`must not be empty`),
						},
					))
				})
			})
		})
	})

	Describe("func WithDefault()", func() {
		When("the port is a valid IANA port name", func() {
			DescribeTable(
				"it does not panics",
				func(port string) {
					Expect(func() {
						spec.WithDefault("default.example.org", port)
					}).NotTo(Panic())
				},
				Entry(
					"lowercase",
					"https",
				),
				Entry(
					"uppercase",
					"HTTPS",
				),
				Entry(
					"mixed",
					"HttpS",
				),
				Entry(
					"hypenated",
					"foo-bar-spam",
				),
				Entry(
					"contains digits",
					"0foo1bar2",
				),
			)
		})
		When("the default is invalid", func() {
			DescribeTable(
				"it panics",
				func(host, port, expect string) {
					Expect(func() {
						spec.WithDefault(host, port)
					}).To(PanicWith(expect))
				},
				Entry(
					"empty host",
					"",
					"12345",
					"default value of FERRITE_SVC_SERVICE_HOST is invalid: must not be empty",
				),
				Entry(
					"empty port",
					"host.example.org",
					"",
					"default value of FERRITE_SVC_SERVICE_PORT is invalid: must not be empty",
				),
				Entry(
					"numeric port too low",
					"host.example.org",
					"0",
					"default value of FERRITE_SVC_SERVICE_PORT is invalid: numeric ports must be between 1 and 65535",
				),
				Entry(
					"numeric port too high",
					"host.example.org",
					"65536",
					"default value of FERRITE_SVC_SERVICE_PORT is invalid: numeric ports must be between 1 and 65535",
				),
				Entry(
					"IANA registered port does not contain any letters",
					"host.example.org",
					"100-200",
					`default value of FERRITE_SVC_SERVICE_PORT is invalid: "100-200" is not a valid numeric port or well-formed IANA service name`,
				),
				Entry(
					"IANA registered port starts with a hyphen",
					"host.example.org",
					"-https",
					`default value of FERRITE_SVC_SERVICE_PORT is invalid: "-https" is not a valid numeric port or well-formed IANA service name`,
				),
			)
		})
	})
})
