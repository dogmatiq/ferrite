package markdown_test

import (
	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable(
	"func Run()",
	tableTest(
		"relationship",
		WithoutExplanatoryText(),
		WithoutUsageExamples(),
	),
	Entry(
		"see also",
		"see-also.md",
		func(reg ferrite.Registry) {
			verbose := ferrite.
				Bool("VERBOSE", "enable verbose logging").
				Optional(ferrite.WithRegistry(reg))

			ferrite.
				Bool("DEBUG", "enable or disable debugging features").
				Optional(
					ferrite.WithRegistry(reg),
					ferrite.SeeAlso(verbose),
				)
		},
	),
	Entry(
		"deprecated + superseded",
		"deprecated-superseded.md",
		func(reg ferrite.Registry) {
			addr := ferrite.
				String(`BIND_HOST`, "listen host for the HTTP server").
				WithDefault("0.0.0.0").
				Required(ferrite.WithRegistry(reg))

			port := ferrite.
				NetworkPort("BIND_PORT", "listen port for the HTTP server").
				WithDefault("8080").
				Required(ferrite.WithRegistry(reg))

			version := ferrite.
				String("BIND_VERSION", "IP version for the HTTP server").
				WithDefault("4").
				Required(ferrite.WithRegistry(reg))

			ferrite.
				String("BIND_ADDRESS", "listen address for the HTTP server").
				WithDefault("0.0.0.0:8080").
				Deprecated(
					ferrite.WithRegistry(reg),
					ferrite.SupersededBy(addr),
					ferrite.SupersededBy(port),
					ferrite.SupersededBy(version),
				)
		},
	),
	Entry(
		"depends on + optional + default",
		"depends-on/with-default.md",
		func(reg ferrite.Registry) {
			widgetEnabled := ferrite.
				Bool("WIDGET_ENABLED", "enable the widget").
				Required(ferrite.WithRegistry(reg))

			ferrite.
				String("WIDGET_COLOR", "the color of the widget").
				WithDefault("turquoise").
				Optional(
					ferrite.WithRegistry(reg),
					ferrite.RelevantIf(widgetEnabled),
				)
		},
	),
	Entry(
		"depends on + optional",
		"depends-on/optional.md",
		func(reg ferrite.Registry) {
			widgetEnabled := ferrite.
				Bool("WIDGET_ENABLED", "enable the widget").
				Required(ferrite.WithRegistry(reg))

			ferrite.
				String("WIDGET_COLOR", "the color of the widget").
				Optional(
					ferrite.WithRegistry(reg),
					ferrite.RelevantIf(widgetEnabled),
				)
		},
	),
	Entry(
		"depends on + required + default",
		"depends-on/with-default.md",
		func(reg ferrite.Registry) {
			widgetEnabled := ferrite.
				Bool("WIDGET_ENABLED", "enable the widget").
				Required(ferrite.WithRegistry(reg))

			ferrite.
				String("WIDGET_COLOR", "the color of the widget").
				WithDefault("turquoise").
				Required(
					ferrite.WithRegistry(reg),
					ferrite.RelevantIf(widgetEnabled),
				)
		},
	),
	Entry(
		"depends on + required",
		"depends-on/required.md",
		func(reg ferrite.Registry) {
			widgetEnabled := ferrite.
				Bool("WIDGET_ENABLED", "enable the widget").
				Required(ferrite.WithRegistry(reg))

			ferrite.
				String("WIDGET_COLOR", "the color of the widget").
				Required(
					ferrite.WithRegistry(reg),
					ferrite.RelevantIf(widgetEnabled),
				)
		},
	),
	Entry(
		"depends on + required (string)",
		"depends-on/required-string.md",
		func(reg ferrite.Registry) {
			widgetType := ferrite.
				String("WIDGET_TYPE", "the type of widget").
				Required(ferrite.WithRegistry(reg))

			ferrite.
				String("WIDGET_COLOR", "the color of the widget").
				Required(
					ferrite.WithRegistry(reg),
					ferrite.RelevantIf(widgetType),
				)
		},
	),
	Entry(
		"depends on + required + constraint",
		"depends-on/required-with-constraint.md",
		func(reg ferrite.Registry) {
			widgetEnabled := ferrite.
				Bool("WIDGET_ENABLED", "enable the widget").
				Required(ferrite.WithRegistry(reg))

			ferrite.
				String("WIDGET_COLOR", "the color of the widget").
				WithConstraint(
					"must be a valid CSS color",
					func(s string) bool { return true },
				).
				Required(
					ferrite.WithRegistry(reg),
					ferrite.RelevantIf(widgetEnabled),
				)
		},
	),
	Entry(
		"depends on + deprecated",
		"depends-on/deprecated.md",
		func(reg ferrite.Registry) {
			widgetEnabled := ferrite.
				Bool("WIDGET_ENABLED", "enable the widget").
				Required(ferrite.WithRegistry(reg))

			ferrite.
				String("WIDGET_COLOR", "the color of the widget").
				Deprecated(
					ferrite.WithRegistry(reg),
					ferrite.RelevantIf(widgetEnabled),
				)
		},
	),
	Entry(
		"depends on + multiple (simple requirement)",
		"depends-on/multiple-simple.md",
		func(reg ferrite.Registry) {
			colorEnabled := ferrite.
				Bool("COLOR_ENABLED", "enable colors").
				Required(ferrite.WithRegistry(reg))

			widgetEnabled := ferrite.
				Bool("WIDGET_ENABLED", "enable the widget").
				Required(ferrite.WithRegistry(reg))

			ferrite.
				String("WIDGET_COLOR", "the color of the widget").
				Required(
					ferrite.WithRegistry(reg),
					ferrite.RelevantIf(colorEnabled),
					ferrite.RelevantIf(widgetEnabled),
				)
		},
	),
	Entry(
		"depends on + multiple (complex requirement)",
		"depends-on/multiple-complex.md",
		func(reg ferrite.Registry) {
			colorEnabled := ferrite.
				Bool("COLOR_ENABLED", "enable colors").
				Required(ferrite.WithRegistry(reg))

			widgetEnabled := ferrite.
				Enum("WIDGET_MODE", "set the widget mode").
				WithMembers("grayscale", "color").
				Required(ferrite.WithRegistry(reg))

			ferrite.
				String("WIDGET_COLOR", "the color of the widget").
				Required(
					ferrite.WithRegistry(reg),
					ferrite.RelevantIf(colorEnabled),
					ferrite.RelevantWhen(widgetEnabled, "color"),
				)
		},
	),
	Entry(
		"depends on + specific value",
		"depends-on/with-value.md",
		func(reg ferrite.Registry) {
			redisService := ferrite.
				KubernetesService("redis").
				Required(ferrite.WithRegistry(reg))

			ferrite.
				Bool("REDIS_DEBUG", "debug local redis connections").
				Required(
					ferrite.WithRegistry(reg),
					ferrite.RelevantWhen(
						redisService,
						ferrite.KubernetesAddress{
							Host: "localhost",
							Port: "6379",
						},
					),
				)
		},
	),
)
