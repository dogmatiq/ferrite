package markdown_test

import (
	"bytes"
	"strings"
	"time"

	"github.com/dogmatiq/ferrite"
	. "github.com/dogmatiq/ferrite/internal/mode/usage/markdown"
	"github.com/dogmatiq/ferrite/variable"
	. "github.com/jmalloc/gomegax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable(
	"func Run()",
	tableTest(
		"spec",
		WithoutPreamble(),
		WithoutIndex(),
		WithoutUsageExamples(),
	),

	// BOOL

	Entry(
		"bool + optional + default",
		"bool/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Bool("DEBUG", "enable or disable debugging features").
				WithDefault(false).
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"bool + optional",
		"bool/optional.md",
		func(reg *variable.Registry) {
			ferrite.
				Bool("DEBUG", "enable or disable debugging features").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"bool + required + default",
		"bool/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Bool("DEBUG", "enable or disable debugging features").
				WithDefault(false).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"bool + required",
		"bool/required.md",
		func(reg *variable.Registry) {
			ferrite.
				Bool("DEBUG", "enable or disable debugging features").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"bool + deprecated + default",
		"bool/deprecated-with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Bool("DEBUG", "enable or disable debugging features").
				WithDefault(false).
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"bool + deprecated",
		"bool/deprecated.md",
		func(reg *variable.Registry) {
			ferrite.
				Bool("DEBUG", "enable or disable debugging features").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),

	// ENUM

	Entry(
		"enum + optional + default",
		"enum/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Enum("LOG_LEVEL", "the minimum log level to record").
				WithMember("debug", "show information for developers").
				WithMember("info", "standard log messages").
				WithMember("warn", "important, but don't need individual human review").
				WithMember("error", "a healthy application shouldn't produce any errors").
				WithMember("fatal", "the application cannot proceed").
				WithDefault("error").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"enum + optional",
		"enum/optional.md",
		func(reg *variable.Registry) {
			ferrite.
				Enum("LOG_LEVEL", "the minimum log level to record").
				WithMember("debug", "show information for developers").
				WithMember("info", "standard log messages").
				WithMember("warn", "important, but don't need individual human review").
				WithMember("error", "a healthy application shouldn't produce any errors").
				WithMember("fatal", "the application cannot proceed").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"enum + required + default",
		"enum/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Enum("LOG_LEVEL", "the minimum log level to record").
				WithMember("debug", "show information for developers").
				WithMember("info", "standard log messages").
				WithMember("warn", "important, but don't need individual human review").
				WithMember("error", "a healthy application shouldn't produce any errors").
				WithMember("fatal", "the application cannot proceed").
				WithDefault("error").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"enum + required",
		"enum/required.md",
		func(reg *variable.Registry) {
			ferrite.
				Enum("LOG_LEVEL", "the minimum log level to record").
				WithMember("debug", "show information for developers").
				WithMember("info", "standard log messages").
				WithMember("warn", "important, but don't need individual human review").
				WithMember("error", "a healthy application shouldn't produce any errors").
				WithMember("fatal", "the application cannot proceed").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"enum + deprecated + default",
		"enum/deprecated-with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Enum("LOG_LEVEL", "the minimum log level to record").
				WithMember("debug", "show information for developers").
				WithMember("info", "standard log messages").
				WithMember("warn", "important, but don't need individual human review").
				WithMember("error", "a healthy application shouldn't produce any errors").
				WithMember("fatal", "the application cannot proceed").
				WithDefault("error").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"enum + deprecated",
		"enum/deprecated.md",
		func(reg *variable.Registry) {
			ferrite.
				Enum("LOG_LEVEL", "the minimum log level to record").
				WithMember("debug", "show information for developers").
				WithMember("info", "standard log messages").
				WithMember("warn", "important, but don't need individual human review").
				WithMember("error", "a healthy application shouldn't produce any errors").
				WithMember("fatal", "the application cannot proceed").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),

	// DURATION

	Entry(
		"duration + optional + default",
		"duration/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Duration("GRPC_TIMEOUT", "gRPC request timeout").
				WithDefault(10 * time.Millisecond).
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"duration + optional",
		"duration/optional.md",
		func(reg *variable.Registry) {
			ferrite.
				Duration("GRPC_TIMEOUT", "gRPC request timeout").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"duration + required + default",
		"duration/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Duration("GRPC_TIMEOUT", "gRPC request timeout").
				WithDefault(10 * time.Millisecond).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"duration + required",
		"duration/required.md",
		func(reg *variable.Registry) {
			ferrite.
				Duration("GRPC_TIMEOUT", "gRPC request timeout").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"duration + maximum",
		"duration/with-max.md",
		func(reg *variable.Registry) {
			ferrite.
				Duration("GRPC_TIMEOUT", "gRPC request timeout").
				WithMaximum(24 * time.Hour).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"duration + deprecated + default",
		"duration/deprecated-with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Duration("GRPC_TIMEOUT", "gRPC request timeout").
				WithDefault(10 * time.Millisecond).
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"duration + deprecated",
		"duration/deprecated.md",
		func(reg *variable.Registry) {
			ferrite.
				Duration("GRPC_TIMEOUT", "gRPC request timeout").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),

	// KUBERNETES SERVICE

	Entry(
		"k8s service + optional + default",
		"k8s-service/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				KubernetesService("redis").
				WithDefault("redis.example.org", "6379").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"k8s service + optional",
		"k8s-service/optional.md",
		func(reg *variable.Registry) {
			ferrite.
				KubernetesService("redis").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"k8s service + required + default",
		"k8s-service/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				KubernetesService("redis").
				WithDefault("redis.example.org", "6379").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"k8s service + required",
		"k8s-service/required.md",
		func(reg *variable.Registry) {
			ferrite.
				KubernetesService("redis").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"k8s service + deprecated + default",
		"k8s-service/deprecated-with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				KubernetesService("redis").
				WithDefault("redis.example.org", "6379").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"k8s service + deprecated",
		"k8s-service/deprecated.md",
		func(reg *variable.Registry) {
			ferrite.
				KubernetesService("redis").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	// NETWORK PORT

	Entry(
		"network port + optional + default",
		"network-port/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				NetworkPort("PORT", "listen port for the HTTP server").
				WithDefault("8080").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"network port + optional",
		"network-port/optional.md",
		func(reg *variable.Registry) {
			ferrite.
				NetworkPort("PORT", "listen port for the HTTP server").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"network port + required + default",
		"network-port/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				NetworkPort("PORT", "listen port for the HTTP server").
				WithDefault("8080").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"network port + required",
		"network-port/required.md",
		func(reg *variable.Registry) {
			ferrite.
				NetworkPort("PORT", "listen port for the HTTP server").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"network port + deprecated + default",
		"network-port/deprecated-with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				NetworkPort("PORT", "listen port for the HTTP server").
				WithDefault("8080").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"network port + deprecated",
		"network-port/deprecated.md",
		func(reg *variable.Registry) {
			ferrite.
				NetworkPort("PORT", "listen port for the HTTP server").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),

	// NUMBER - FLOAT

	Entry(
		"float + optional + default",
		"float/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Float[float32]("WEIGHT", "weighting for this node").
				WithDefault(123.5).
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"float + optional",
		"float/optional.md",
		func(reg *variable.Registry) {
			ferrite.
				Float[float32]("WEIGHT", "weighting for this node").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"float + required + default",
		"float/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Float[float32]("WEIGHT", "weighting for this node").
				WithDefault(123.5).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"float + required",
		"float/required.md",
		func(reg *variable.Registry) {
			ferrite.
				Float[float32]("WEIGHT", "weighting for this node").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"float + min",
		"float/with-min.md",
		func(reg *variable.Registry) {
			ferrite.
				Float[float32]("WEIGHT", "weighting for this node").
				WithMinimum(-10.5).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"float + max",
		"float/with-max.md",
		func(reg *variable.Registry) {
			ferrite.
				Float[float32]("WEIGHT", "weighting for this node").
				WithMaximum(+20.5).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"float + min/max",
		"float/with-minmax.md",
		func(reg *variable.Registry) {
			ferrite.
				Float[float32]("WEIGHT", "weighting for this node").
				WithMinimum(-10.5).
				WithMaximum(+20.5).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"float + deprecated + default",
		"float/deprecated-with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Float[float32]("WEIGHT", "weighting for this node").
				WithDefault(123.5).
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"float + deprecated",
		"float/deprecated.md",
		func(reg *variable.Registry) {
			ferrite.
				Float[float32]("WEIGHT", "weighting for this node").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),

	// NUMBER - SIGNED

	Entry(
		"signed + optional + default",
		"signed/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Signed[int8]("WEIGHT", "weighting for this node").
				WithDefault(100).
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"signed + optional",
		"signed/optional.md",
		func(reg *variable.Registry) {
			ferrite.
				Signed[int8]("WEIGHT", "weighting for this node").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"signed + required + default",
		"signed/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Signed[int8]("WEIGHT", "weighting for this node").
				WithDefault(100).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"signed + required",
		"signed/required.md",
		func(reg *variable.Registry) {
			ferrite.
				Signed[int8]("WEIGHT", "weighting for this node").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"signed + min",
		"signed/with-min.md",
		func(reg *variable.Registry) {
			ferrite.
				Signed[int8]("WEIGHT", "weighting for this node").
				WithMinimum(-10).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"signed + max",
		"signed/with-max.md",
		func(reg *variable.Registry) {
			ferrite.
				Signed[int8]("WEIGHT", "weighting for this node").
				WithMaximum(+20).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"signed + min/max",
		"signed/with-minmax.md",
		func(reg *variable.Registry) {
			ferrite.
				Signed[int8]("WEIGHT", "weighting for this node").
				WithMinimum(-10).
				WithMaximum(+20).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"signed + deprecated + default",
		"signed/deprecated-with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Signed[int8]("WEIGHT", "weighting for this node").
				WithDefault(100).
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"signed + deprecated",
		"signed/deprecated.md",
		func(reg *variable.Registry) {
			ferrite.
				Signed[int8]("WEIGHT", "weighting for this node").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),

	// NUMBER - UNSIGNED

	Entry(
		"unsigned + optional + default",
		"unsigned/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				WithDefault(900).
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"unsigned + optional",
		"unsigned/optional.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"unsigned + required + default",
		"unsigned/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				WithDefault(900).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"unsigned + required",
		"unsigned/required.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"unsigned + min",
		"unsigned/with-min.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				WithMinimum(10).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"unsigned + max",
		"unsigned/with-max.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				WithMaximum(20).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"unsigned + min/max",
		"unsigned/with-minmax.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				WithMinimum(10).
				WithMaximum(20).
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"unsigned + deprecated + default",
		"unsigned/deprecated-with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				WithDefault(900).
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"unsigned + deprecated",
		"unsigned/deprecated.md",
		func(reg *variable.Registry) {
			ferrite.
				Unsigned[uint16]("WEIGHT", "weighting for this node").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),

	// STRING

	Entry(
		"string + optional + default",
		"string/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				String("READ_DSN", "database connection string for read-models").
				WithDefault("host=localhost dbname=readmodels user=projector").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"string + optional",
		"string/optional.md",
		func(reg *variable.Registry) {
			ferrite.
				String("READ_DSN", "database connection string for read-models").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"string + required + default",
		"string/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				String("READ_DSN", "database connection string for read-models").
				WithDefault("host=localhost dbname=readmodels user=projector").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"string + required",
		"string/required.md",
		func(reg *variable.Registry) {
			ferrite.
				String("READ_DSN", "database connection string for read-models").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"string + deprecated + default",
		"string/deprecated-with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				String("READ_DSN", "database connection string for read-models").
				WithDefault("host=localhost dbname=readmodels user=projector").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"string + deprecated",
		"string/deprecated.md",
		func(reg *variable.Registry) {
			ferrite.
				String("READ_DSN", "database connection string for read-models").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"string + sensitive + optional + default",
		"string/sensitive-with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				String("PASSWORD", "a very secret password").
				WithDefault("hunter2").
				WithSensitiveContent().
				Optional(
					ferrite.WithRegistry(reg),
				)
		},
	),
	Entry(
		"string + sensitive + optional",
		"string/sensitive-optional.md",
		func(reg *variable.Registry) {
			ferrite.
				String("PASSWORD", "a very secret password").
				WithSensitiveContent().
				Optional(
					ferrite.WithRegistry(reg),
				)
		},
	),
	Entry(
		"string + sensitive + required + default",
		"string/sensitive-with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				String("PASSWORD", "a very secret password").
				WithDefault("hunter2").
				WithSensitiveContent().
				Required(
					ferrite.WithRegistry(reg),
				)
		},
	),
	Entry(
		"string + sensitive + required",
		"string/sensitive-required.md",
		func(reg *variable.Registry) {
			ferrite.
				String("PASSWORD", "a very secret password").
				WithSensitiveContent().
				Required(
					ferrite.WithRegistry(reg),
				)
		},
	),

	// URL

	Entry(
		"url + optional + default",
		"url/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				URL("API_URL", "the URL of the REST API").
				WithDefault("http://localhost:8080").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"url + optional",
		"url/optional.md",
		func(reg *variable.Registry) {
			ferrite.
				URL("API_URL", "the URL of the REST API").
				Optional(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"url + required + default",
		"url/with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				URL("API_URL", "the URL of the REST API").
				WithDefault("http://localhost:8080").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"url + required",
		"url/required.md",
		func(reg *variable.Registry) {
			ferrite.
				URL("API_URL", "the URL of the REST API").
				Required(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"url + deprecated + default",
		"url/deprecated-with-default.md",
		func(reg *variable.Registry) {
			ferrite.
				URL("API_URL", "the URL of the REST API").
				WithDefault("http://localhost:8080").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
	Entry(
		"url + deprecated",
		"url/deprecated.md",
		func(reg *variable.Registry) {
			ferrite.
				URL("API_URL", "the URL of the REST API").
				Deprecated(ferrite.WithRegistry(reg))
		},
	),
)

// expectLines verifies that text consists of the given lines.
func expectLines(buf *bytes.Buffer, lines ...string) {
	actual := buf.String()
	expect := strings.Join(lines, "\n") + "\n"
	ExpectWithOffset(1, actual).To(EqualX(expect))
}
