package ferrite_test

import (
	"os"

	. "github.com/dogmatiq/ferrite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type KubernetesServiceBuilder", func() {
	var builder *KubernetesServiceBuilder

	BeforeEach(func() {
		builder = KubernetesService("ferrite-svc")
	})

	AfterEach(func() {
		tearDown()
	})

	When("the variables are required", func() {
		When("the the host and port are valid", func() {
			BeforeEach(func() {
				os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
				os.Setenv("FERRITE_SVC_SERVICE_PORT", "12345")
			})

			It("returns the value", func() {
				v := builder.
					Required().
					Value()

				Expect(v).To(Equal(
					KubernetesAddress{
						Host: "host.example.org",
						Port: "12345",
					},
				))
			})
		})

		When("the host is invalid", func() {
			Describe("func Value()", func() {
				invalidHostTable(
					"it panics",
					func(value, expect string) {
						os.Setenv("FERRITE_SVC_SERVICE_HOST", value) // note trailing dot
						os.Setenv("FERRITE_SVC_SERVICE_PORT", "12345")

						Expect(func() {
							builder.
								Required().
								Value()
						}).To(PanicWith(expect))
					},
				)
			})
		})

		When("the port is invalid", func() {
			Describe("func Value()", func() {
				invalidPortTable(
					"it panics",
					func(value, expect string) {
						Expect(func() {
							os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
							os.Setenv("FERRITE_SVC_SERVICE_PORT", value)

							builder.
								Required().
								Value()
						}).To(PanicWith(expect))
					},
				)
			})
		})

		When("both values are empty", func() {
			When("there is a default value", func() {
				BeforeEach(func() {
					builder = builder.WithDefault("default.example.org", "54321")
				})

				Describe("func Value()", func() {
					It("returns the default", func() {
						v := builder.
							Required().
							Value()

						Expect(v).To(Equal(
							KubernetesAddress{
								Host: "default.example.org",
								Port: "54321",
							},
						))
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
							"FERRITE_SVC_SERVICE_HOST is undefined and does not have a default value",
						))
					})
				})
			})
		})

		When("the host is empty", func() {
			BeforeEach(func() {
				os.Setenv("FERRITE_SVC_SERVICE_PORT", "12345")
			})

			When("there is a default value", func() {
				BeforeEach(func() {
					builder = builder.WithDefault("default.example.org", "54321")
				})

				Describe("func Value()", func() {
					It("returns the value with a default host", func() {
						v := builder.
							Required().
							Value()

						Expect(v).To(Equal(
							KubernetesAddress{
								Host: "default.example.org",
								Port: "12345",
							},
						))
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
							"FERRITE_SVC_SERVICE_HOST is undefined and does not have a default value",
						))
					})
				})
			})
		})

		When("the port is empty", func() {
			BeforeEach(func() {
				os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
			})

			When("there is a default value", func() {
				BeforeEach(func() {
					builder = builder.WithDefault("default.example.org", "54321")
				})

				Describe("func Value()", func() {
					It("returns the value with a default port", func() {
						v := builder.
							Required().
							Value()

						Expect(v).To(Equal(
							KubernetesAddress{
								Host: "host.example.org",
								Port: "54321",
							},
						))
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
							"FERRITE_SVC_SERVICE_PORT is undefined and does not have a default value",
						))
					})
				})
			})
		})
	})

	When("the variables are optional", func() {
		When("the the host and port are valid", func() {
			BeforeEach(func() {
				os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
				os.Setenv("FERRITE_SVC_SERVICE_PORT", "12345")
			})

			It("returns the value", func() {
				v, ok := builder.
					Optional().
					Value()

				Expect(ok).To(BeTrue())
				Expect(v).To(Equal(
					KubernetesAddress{
						Host: "host.example.org",
						Port: "12345",
					},
				))
			})
		})

		When("the host is invalid", func() {
			Describe("func Value()", func() {
				invalidHostTable(
					"it panics",
					func(value, expect string) {
						os.Setenv("FERRITE_SVC_SERVICE_HOST", value) // note trailing dot
						os.Setenv("FERRITE_SVC_SERVICE_PORT", "12345")

						Expect(func() {
							builder.
								Optional().
								Value()
						}).To(PanicWith(expect))
					},
				)
			})
		})

		When("the port is invalid", func() {
			Describe("func Value()", func() {
				invalidPortTable(
					"it panics",
					func(value, expect string) {
						Expect(func() {
							os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
							os.Setenv("FERRITE_SVC_SERVICE_PORT", value)

							builder.
								Optional().
								Value()
						}).To(PanicWith(expect))
					},
				)
			})
		})

		When("both values are empty", func() {
			When("there is a default value", func() {
				BeforeEach(func() {
					builder = builder.WithDefault("default.example.org", "54321")
				})

				Describe("func Value()", func() {
					It("returns the default", func() {
						v, ok := builder.
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(
							KubernetesAddress{
								Host: "default.example.org",
								Port: "54321",
							},
						))
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

		When("the host is empty", func() {
			BeforeEach(func() {
				os.Setenv("FERRITE_SVC_SERVICE_PORT", "12345")
			})

			When("there is a default value", func() {
				BeforeEach(func() {
					builder = builder.WithDefault("default.example.org", "54321")
				})

				Describe("func Value()", func() {
					It("returns the value with a default host", func() {
						v, ok := builder.
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(
							KubernetesAddress{
								Host: "default.example.org",
								Port: "12345",
							},
						))
					})
				})
			})

			When("there is no default value", func() {
				Describe("func Value()", func() {
					It("panics", func() {
						Expect(func() {
							builder.
								Optional().
								Value()
						}).To(PanicWith(
							`FERRITE_SVC_SERVICE_PORT is defined but FERRITE_SVC_SERVICE_HOST is not, define both or neither`,
						))
					})
				})
			})
		})

		When("the port is empty", func() {
			BeforeEach(func() {
				os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
			})

			When("there is a default value", func() {
				BeforeEach(func() {
					builder = builder.WithDefault("default.example.org", "54321")
				})

				Describe("func Value()", func() {
					It("returns the value with a default port", func() {
						v, ok := builder.
							Optional().
							Value()

						Expect(ok).To(BeTrue())
						Expect(v).To(Equal(
							KubernetesAddress{
								Host: "host.example.org",
								Port: "54321",
							},
						))
					})
				})
			})

			When("there is no default value", func() {
				Describe("func Value()", func() {
					It("panics", func() {
						Expect(func() {
							builder.
								Optional().
								Value()
						}).To(PanicWith(
							`FERRITE_SVC_SERVICE_HOST is defined but FERRITE_SVC_SERVICE_PORT is not, define both or neither`,
						))
					})
				})
			})
		})
	})

	When("using a named port", func() {
		BeforeEach(func() {
			builder = builder.WithNamedPort("named-port")

			os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
			os.Setenv("FERRITE_SVC_SERVICE_PORT_NAMED_PORT", "12345")
		})

		Describe("func Value()", func() {
			It("returns the value", func() {
				v := builder.
					Required().
					Value()

				Expect(v).To(Equal(
					KubernetesAddress{
						Host: "host.example.org",
						Port: "12345",
					},
				))
			})
		})

		When("the port name is invalid", func() {
			Describe("func WithNamedPort()", func() {
				DescribeTable(
					"it panics",
					func(port, expect string) {
						Expect(func() {
							builder.WithNamedPort(port)
						}).To(PanicWith(expect))
					},
					Entry(
						"empty",
						"",
						`specification of kubernetes "ferrite-svc" service is invalid: invalid named port: name must not be empty`,
					),
					Entry(
						"starts with a hyphen",
						"-foo",
						`specification of kubernetes "ferrite-svc" service is invalid: invalid named port: name must not begin or end with a hyphen`,
					),
					Entry(
						"ends with a hyphen",
						"foo-",
						`specification of kubernetes "ferrite-svc" service is invalid: invalid named port: name must not begin or end with a hyphen`,
					),
					Entry(
						"contains an invalid character",
						"foo*bar",
						`specification of kubernetes "ferrite-svc" service is invalid: invalid named port: name must contain only lowercase ASCII letters, digits and hyphen`,
					),
					Entry(
						"contains an uppercase character",
						"fooBar",
						`specification of kubernetes "ferrite-svc" service is invalid: invalid named port: name must contain only lowercase ASCII letters, digits and hyphen`,
					),
				)
			})
		})
	})

	Describe("func WithDefault()", func() {
		When("the host is valid", func() {
			DescribeTable(
				"it does not panic",
				func(host string) {
					Expect(func() {
						builder.WithDefault(host, "12345")
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
						builder.WithDefault("default.example.org", port)
					}).NotTo(Panic())
				},
				Entry(
					"numeric",
					"12345",
				),
				Entry(
					"IANA service name in lowercase",
					"https",
				),
				Entry(
					"IANA service name in uppercase",
					"HTTPS",
				),
				Entry(
					"IANA service name with mixedcase",
					"HttpS",
				),
				Entry(
					"IANA service name contains hyphens",
					"foo-bar-spam",
				),
				Entry(
					"IANA service name contains digits",
					"0foo1bar2",
				),
			)
		})

		When("the default is invalid", func() {
			DescribeTable(
				"it panics",
				func(host, port, expect string) {
					Expect(func() {
						builder.
							WithDefault(host, port).
							Required()
					}).To(PanicWith(expect))
				},
				Entry(
					"empty host",
					"",
					"12345",
					`specification for FERRITE_SVC_SERVICE_HOST is invalid: default value: host must not be empty`,
				),
				Entry(
					"hostname begins with a dot",
					".foo",
					"12345",
					`specification for FERRITE_SVC_SERVICE_HOST is invalid: default value: host must not begin or end with a dot`,
				),
				Entry(
					"hostname ends with a dot",
					"foo.",
					"12345",
					`specification for FERRITE_SVC_SERVICE_HOST is invalid: default value: host must not begin or end with a dot`,
				),
				Entry(
					"hostname contains whitespace",
					"foo .bar",
					"12345",
					`specification for FERRITE_SVC_SERVICE_HOST is invalid: default value: host must not contain whitespace`,
				),
				Entry(
					"empty port",
					"host.example.org",
					"",
					`specification for FERRITE_SVC_SERVICE_PORT is invalid: default value: port must not be empty`,
				),
				Entry(
					"numeric port too low",
					"host.example.org",
					"0",
					`specification for FERRITE_SVC_SERVICE_PORT is invalid: default value: numeric ports must be between 1 and 65535`,
				),
				Entry(
					"numeric port too high",
					"host.example.org",
					"65536",
					`specification for FERRITE_SVC_SERVICE_PORT is invalid: default value: numeric ports must be between 1 and 65535`,
				),
				Entry(
					"IANA service name is too long",
					"host.example.org",
					"this-name-is-very-long",
					`specification for FERRITE_SVC_SERVICE_PORT is invalid: default value: IANA service name must be between 1 and 15 characters`,
				),
				Entry(
					"IANA service name does not contain any letters",
					"host.example.org",
					"100-200",
					`specification for FERRITE_SVC_SERVICE_PORT is invalid: default value: IANA service name must contain at least one letter`,
				),
				Entry(
					"IANA service name starts with a hyphen",
					"host.example.org",
					"-foo",
					`specification for FERRITE_SVC_SERVICE_PORT is invalid: default value: IANA service name must not begin or end with a hyphen`,
				),
				Entry(
					"IANA service name ends with a hyphen",
					"host.example.org",
					"foo-",
					`specification for FERRITE_SVC_SERVICE_PORT is invalid: default value: IANA service name must not begin or end with a hyphen`,
				),
				Entry(
					"IANA service name contains adjacent hyphens",
					"host.example.org",
					"foo--bar",
					`specification for FERRITE_SVC_SERVICE_PORT is invalid: default value: IANA service name must not contain adjacent hyphens`,
				),
				Entry(
					"IANA service name contains an invalid character",
					"host.example.org",
					"foo*bar",
					`specification for FERRITE_SVC_SERVICE_PORT is invalid: default value: IANA service name must contain only ASCII letters, digits and hyphen`,
				),
			)
		})
	})
})

func invalidHostTable(
	desc string,
	fn func(value, expect string),
) {
	DescribeTable(
		desc,
		fn,
		Entry(
			"leading dot",
			".host.example.org",
			`value of FERRITE_SVC_SERVICE_HOST (.host.example.org) is invalid: host must not begin or end with a dot`,
		),
		Entry(
			"trailing dot",
			"host.example.org.",
			`value of FERRITE_SVC_SERVICE_HOST (host.example.org.) is invalid: host must not begin or end with a dot`,
		),
		Entry(
			"whitespace",
			"host.examp le.org",
			`value of FERRITE_SVC_SERVICE_HOST ('host.examp le.org') is invalid: host must not contain whitespace`,
		),
	)
}

func invalidPortTable(
	desc string,
	fn func(value, expect string),
) {
	DescribeTable(
		desc,
		fn,
		Entry(
			"numeric port too low",
			"0",
			`value of FERRITE_SVC_SERVICE_PORT (0) is invalid: numeric ports must be between 1 and 65535`,
		),
		Entry(
			"numeric port too high",
			"65536",
			`value of FERRITE_SVC_SERVICE_PORT (65536) is invalid: numeric ports must be between 1 and 65535`,
		),
		Entry(
			"IANA service name is too long",
			"this-name-is-very-long",
			`value of FERRITE_SVC_SERVICE_PORT (this-name-is-very-long) is invalid: IANA service name must be between 1 and 15 characters`,
		),
		Entry(
			"IANA service name does not contain any letters",
			"100-200",
			`value of FERRITE_SVC_SERVICE_PORT (100-200) is invalid: IANA service name must contain at least one letter`,
		),
		Entry(
			"IANA service name starts with a hyphen",
			"-foo",
			`value of FERRITE_SVC_SERVICE_PORT (-foo) is invalid: IANA service name must not begin or end with a hyphen`,
		),
		Entry(
			"IANA service name ends with a hyphen",
			"foo-",
			`value of FERRITE_SVC_SERVICE_PORT (foo-) is invalid: IANA service name must not begin or end with a hyphen`,
		),
		Entry(
			"IANA service name contains adjacent hyphens",
			"foo--bar",
			`value of FERRITE_SVC_SERVICE_PORT (foo--bar) is invalid: IANA service name must not contain adjacent hyphens`,
		),
		Entry(
			"IANA service name contains an invalid character",
			"foo*bar",
			`value of FERRITE_SVC_SERVICE_PORT ('foo*bar') is invalid: IANA service name must contain only ASCII letters, digits and hyphen`,
		),
	)
}
