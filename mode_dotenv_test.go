package ferrite_test

import (
	"os"
	"time"

	"github.com/dogmatiq/ferrite"
)

func ExampleInit_exportDotEnvFile() {
	defer example()()

	os.Setenv("FERRITE_BOOL", "true")
	ferrite.
		Bool("FERRITE_BOOL", "example bool").
		Required()

	os.Setenv("FERRITE_DURATION", "620s")
	ferrite.
		Duration("FERRITE_DURATION", "example duration").
		WithDefault(1 * time.Hour).
		Required()

	ferrite.
		Enum("FERRITE_ENUM", "example enum").
		WithMembers("foo", "bar", "baz").
		WithDefault("bar").
		Required()

	os.Setenv("FERRITE_NETWORK_PORT", "8080")
	ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port").
		Optional()

	ferrite.
		Float[float32]("FERRITE_NUM_FLOAT", "example floating-point").
		Required()

	ferrite.
		Signed[int16]("FERRITE_NUM_SIGNED", "example signed integer").
		Required()

	ferrite.
		Unsigned[uint16]("FERRITE_NUM_UNSIGNED", "example unsigned integer").
		Required()

	os.Setenv("FERRITE_STRING", "hello, world!")
	ferrite.
		String("FERRITE_STRING", "example string").
		Required()

	os.Setenv("FERRITE_STRING_SENSITIVE", "hunter2")
	ferrite.
		String("FERRITE_STRING_SENSITIVE", "example sensitive string").
		WithDefault("password").
		WithSensitiveContent().
		Required()

	os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
	os.Setenv("FERRITE_SVC_SERVICE_PORT", "443")
	ferrite.
		KubernetesService("ferrite-svc").
		Deprecated()

	os.Setenv("FERRITE_URL", "https//example.org")
	ferrite.
		URL("FERRITE_URL", "example URL").
		Required()

	// Tell ferrite to export an env file containing the environment variables.
	os.Setenv("FERRITE_MODE", "export/dotenv")

	ferrite.Init()

	// Output:
	// # example bool (required)
	// export FERRITE_BOOL=true
	//
	// # example duration (default: 1h)
	// export FERRITE_DURATION=620s # equivalent to 10m20s
	//
	// # example enum (default: bar)
	// export FERRITE_ENUM=
	//
	// # example network port (optional)
	// export FERRITE_NETWORK_PORT=8080
	//
	// # example floating-point (required)
	// export FERRITE_NUM_FLOAT=
	//
	// # example signed integer (required)
	// export FERRITE_NUM_SIGNED=
	//
	// # example unsigned integer (required)
	// export FERRITE_NUM_UNSIGNED=
	//
	// # example string (required)
	// export FERRITE_STRING='hello, world!'
	//
	// # example sensitive string (default: ********, sensitive)
	// export FERRITE_STRING_SENSITIVE=hunter2
	//
	// # kubernetes "ferrite-svc" service host (deprecated)
	// export FERRITE_SVC_SERVICE_HOST=host.example.org
	//
	// # kubernetes "ferrite-svc" service port (deprecated)
	// export FERRITE_SVC_SERVICE_PORT=443
	//
	// # example URL (required)
	// export FERRITE_URL= # https//example.org is invalid: URL must have a scheme
	// <process exited successfully>
}
