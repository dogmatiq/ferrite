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
						ExplicitValue: "host.example.org",
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

	When("the host is invalid", func() {
		BeforeEach(func() {
			os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org.") // note trailing dot
			os.Setenv("FERRITE_SVC_SERVICE_PORT", "12345")
		})

		Describe("func Address()", func() {
			It("panics", func() {
				Expect(func() {
					spec.Address()
				}).To(PanicWith("FERRITE_SVC_SERVICE_HOST is invalid: hostname must not begin or end with a dot"))
			})
		})

		Describe("func Host()", func() {
			It("panics", func() {
				Expect(func() {
					spec.Host()
				}).To(PanicWith("FERRITE_SVC_SERVICE_HOST is invalid: hostname must not begin or end with a dot"))
			})
		})

		Describe("func Validate()", func() {
			It("returns a failure result for the port", func() {
				Expect(spec.Validate()).To(ConsistOf(
					ValidationResult{
						Name:          "FERRITE_SVC_SERVICE_HOST",
						Description:   `Hostname or IP address of the "ferrite-svc" service.`,
						ValidInput:    "[string]",
						ExplicitValue: "host.example.org.",
						Error:         errors.New(`hostname must not begin or end with a dot`),
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

	When("the port is invalid", func() {
		BeforeEach(func() {
			os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
			os.Setenv("FERRITE_SVC_SERVICE_PORT", "foo-") // note trailing hyphen
		})

		Describe("func Address()", func() {
			It("panics", func() {
				Expect(func() {
					spec.Address()
				}).To(PanicWith(`FERRITE_SVC_SERVICE_PORT is invalid: "foo-" is not a valid IANA service name (must not begin or end with a hyphen)`))
			})
		})

		Describe("func Port()", func() {
			It("panics", func() {
				Expect(func() {
					spec.Port()
				}).To(PanicWith(`FERRITE_SVC_SERVICE_PORT is invalid: "foo-" is not a valid IANA service name (must not begin or end with a hyphen)`))
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
						Name:          "FERRITE_SVC_SERVICE_PORT",
						Description:   `Network port of the "ferrite-svc" service.`,
						ValidInput:    "[string]|(1..65535)",
						ExplicitValue: "foo-",
						Error:         errors.New(`"foo-" is not a valid IANA service name (must not begin or end with a hyphen)`),
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
						ExplicitValue: "host.example.org",
					},
					ValidationResult{
						Name:          "FERRITE_SVC_SERVICE_PORT_NAMED_PORT",
						Description:   `Network port of the "ferrite-svc" service's "named-port" port.`,
						ValidInput:    "[string]|(1..65535)",
						ExplicitValue: "12345",
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

	Describe("func WithNamedPort()", func() {
		When("the port name is invalid", func() {
			DescribeTable(
				"it panics",
				func(port, expect string) {
					Expect(func() {
						spec.WithNamedPort(port)
					}).To(PanicWith(expect))
				},
				Entry(
					"empty",
					"",
					"kubernetes port name is invalid: must not be empty",
				),
				Entry(
					"starts with a hyphen",
					"-foo",
					"kubernetes port name is invalid: must not begin or end with a hyphen",
				),
				Entry(
					"ends with a hyphen",
					"foo-",
					"kubernetes port name is invalid: must not begin or end with a hyphen",
				),
				Entry(
					"contains an invalid character",
					"foo*bar",
					"kubernetes port name is invalid: must contain only lowercase ASCII letters, digits and hyphen",
				),
				Entry(
					"contains an uppercase character",
					"fooBar",
					"kubernetes port name is invalid: must contain only lowercase ASCII letters, digits and hyphen",
				),
			)
		})
	})

	Describe("func WithDefault()", func() {
		When("the host is valid", func() {
			DescribeTable(
				"it does not panic",
				func(host string) {
					Expect(func() {
						spec.WithDefault(host, "12345")
					}).NotTo(Panic())
				},
				Entry(
					"IPv4",
					"192.168.1.2",
				),
				Entry(
					"IPv6",
					"::1",
				),
				Entry(
					"unqualified DNS name",
					"svc-name",
				),
				Entry(
					"qualified DNS name",
					"svc-name.example.org",
				),
			)
		})

		When("the port is valid", func() {
			DescribeTable(
				"it does not panic",
				func(port string) {
					Expect(func() {
						spec.WithDefault("default.example.org", port)
					}).NotTo(Panic())
				},
				Entry(
					"numeric",
					"12345",
				),
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
					"hostname begins with a dot",
					".foo",
					"12345",
					"default value of FERRITE_SVC_SERVICE_HOST is invalid: hostname must not begin or end with a dot",
				),
				Entry(
					"hostname ends with a dot",
					"foo.",
					"12345",
					"default value of FERRITE_SVC_SERVICE_HOST is invalid: hostname must not begin or end with a dot",
				),
				Entry(
					"hostname contains whitespace",
					"foo .bar",
					"12345",
					"default value of FERRITE_SVC_SERVICE_HOST is invalid: hostname must not contain whitespace",
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
					"IANA service name is too long",
					"host.example.org",
					"this-name-is-very-long",
					`default value of FERRITE_SVC_SERVICE_PORT is invalid: "this-name-is-very-long" is not a valid IANA service name (must be between 1 and 15 characters)`,
				),
				Entry(
					"IANA service name does not contain any letters",
					"host.example.org",
					"100-200",
					`default value of FERRITE_SVC_SERVICE_PORT is invalid: "100-200" is not a valid IANA service name (must contain at least one letter)`,
				),
				Entry(
					"IANA service name starts with a hyphen",
					"host.example.org",
					"-foo",
					`default value of FERRITE_SVC_SERVICE_PORT is invalid: "-foo" is not a valid IANA service name (must not begin or end with a hyphen)`,
				),
				Entry(
					"IANA service name ends with a hyphen",
					"host.example.org",
					"foo-",
					`default value of FERRITE_SVC_SERVICE_PORT is invalid: "foo-" is not a valid IANA service name (must not begin or end with a hyphen)`,
				),
				Entry(
					"IANA service name contains adjacent hyphens",
					"host.example.org",
					"foo--bar",
					`default value of FERRITE_SVC_SERVICE_PORT is invalid: "foo--bar" is not a valid IANA service name (must not contain adjacent hyphens)`,
				),
				Entry(
					"IANA service name contains an invalid character",
					"host.example.org",
					"foo*bar",
					`default value of FERRITE_SVC_SERVICE_PORT is invalid: "foo*bar" is not a valid IANA service name (must contain only ASCII letters, digits and hyphen)`,
				),
			)
		})
	})
})
