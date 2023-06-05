package ferrite_test

import (
	"os"
	"strings"
	"time"

	"github.com/dogmatiq/ferrite"
)

func ExampleInit_validation() {
	defer example()()

	os.Setenv("FERRITE_BINARY", "PHZhbHVlPg==")
	ferrite.
		Binary("FERRITE_BINARY", "example binary").
		Required()

	os.Setenv("FERRITE_BINARY_SENSITIVE", "aHVudGVyMg==")
	ferrite.
		Binary("FERRITE_BINARY_SENSITIVE", "example sensitive binary").
		WithSensitiveContent().
		Required()

	os.Setenv("FERRITE_BOOL", "true")
	ferrite.
		Bool("FERRITE_BOOL", "example bool").
		Required()

	os.Setenv("FERRITE_DURATION", "3h20m")
	ferrite.
		Duration("FERRITE_DURATION", "example duration").
		Required()

	os.Setenv("FERRITE_ENUM", "foo")
	ferrite.
		Enum("FERRITE_ENUM", "example enum").
		WithMembers("foo", "bar", "baz").
		Required()

	os.Setenv("FERRITE_NETWORK_PORT", "8080")
	ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port").
		Required()

	os.Setenv("FERRITE_NUM_FLOAT", "-123.45")
	ferrite.
		Float[float32]("FERRITE_NUM_FLOAT", "example float-point").
		Required()

	os.Setenv("FERRITE_NUM_SIGNED", "-123")
	ferrite.
		Signed[int16]("FERRITE_NUM_SIGNED", "example signed integer").
		Required()

	os.Setenv("FERRITE_NUM_UNSIGNED", "456")
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
		WithSensitiveContent().
		Required()

	os.Setenv("FERRITE_SVC_SERVICE_HOST", "host.example.org")
	os.Setenv("FERRITE_SVC_SERVICE_PORT", "443")
	ferrite.
		KubernetesService("ferrite-svc").
		Required()

	os.Setenv("FERRITE_URL", "https://example.org")
	ferrite.
		URL("FERRITE_URL", "example URL").
		Required()

	ferrite.
		String("FERRITE_XTRIGGER", "trigger failure for example").
		Required()

	ferrite.Init()

	// Output:
	// Environment Variables:
	//
	//    FERRITE_BINARY            example binary                           <base64>           ✓ set to {12 bytes}
	//    FERRITE_BINARY_SENSITIVE  example sensitive binary                 <base64>           ✓ set to {12 bytes}
	//    FERRITE_BOOL              example bool                             true | false       ✓ set to true
	//    FERRITE_DURATION          example duration                         1ns ...            ✓ set to 3h20m
	//    FERRITE_ENUM              example enum                             foo | bar | baz    ✓ set to foo
	//    FERRITE_NETWORK_PORT      example network port                     <string>           ✓ set to 8080
	//    FERRITE_NUM_FLOAT         example float-point                      <float32>          ✓ set to -123.45
	//    FERRITE_NUM_SIGNED        example signed integer                   <int16>            ✓ set to -123
	//    FERRITE_NUM_UNSIGNED      example unsigned integer                 <uint16>           ✓ set to 456
	//    FERRITE_STRING            example string                           <string>           ✓ set to 'hello, world!'
	//    FERRITE_STRING_SENSITIVE  example sensitive string                 <string>           ✓ set to *******
	//    FERRITE_SVC_SERVICE_HOST  kubernetes "ferrite-svc" service host    <string>           ✓ set to host.example.org
	//    FERRITE_SVC_SERVICE_PORT  kubernetes "ferrite-svc" service port    <string>           ✓ set to 443
	//    FERRITE_URL               example URL                              <string>           ✓ set to https://example.org
	//  ❯ FERRITE_XTRIGGER          trigger failure for example              <string>           ✗ undefined
	//
	// <process exited with error code 1>
}

func ExampleInit_validationWithDefaultValues() {
	defer example()()

	ferrite.
		Binary("FERRITE_BINARY", "example binary").
		WithDefaultString("PHZhbHVlPg==").
		Required()

	ferrite.
		Binary("FERRITE_BINARY_SENSITIVE", "example sensitive binary").
		WithDefaultString("aHVudGVyMg==").
		WithSensitiveContent().
		Required()

	ferrite.
		Bool("FERRITE_BOOL", "example bool").
		WithDefault(true).
		Required()

	ferrite.
		Duration("FERRITE_DURATION", "example duration").
		WithDefault(10 * time.Second).
		Required()

	ferrite.
		Enum("FERRITE_ENUM", "example enum").
		WithMembers("foo", "bar", "baz").
		WithDefault("bar").
		Required()

	ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port").
		WithDefault("8080").
		Required()

	ferrite.
		Float[float32]("FERRITE_NUM_FLOAT", "example float-point").
		WithDefault(-123.45).
		Required()

	ferrite.
		Signed[int16]("FERRITE_NUM_SIGNED", "example signed integer").
		WithDefault(-123).
		Required()

	ferrite.
		Unsigned[uint16]("FERRITE_NUM_UNSIGNED", "example unsigned integer").
		WithDefault(123).
		Required()

	ferrite.
		String("FERRITE_STRING", "example string").
		WithDefault("hello, world!").
		Required()

	ferrite.
		String("FERRITE_STRING_SENSITIVE", "example sensitive string").
		WithDefault("hunter2").
		WithSensitiveContent().
		Required()

	ferrite.
		KubernetesService("ferrite-svc").
		WithDefault("host.example.org", "443").
		Required()

	ferrite.
		URL("FERRITE_URL", "example URL").
		WithDefault("https://example.org").
		Required()

	ferrite.
		String("FERRITE_XTRIGGER", "trigger failure for example").
		Required()

	ferrite.Init()

	// Output:
	// Environment Variables:
	//
	//    FERRITE_BINARY            example binary                         [ <base64> ] = {16 bytes}           ✓ using default value
	//    FERRITE_BINARY_SENSITIVE  example sensitive binary               [ <base64> ] = {16 bytes}           ✓ using default value
	//    FERRITE_BOOL              example bool                           [ true | false ] = true             ✓ using default value
	//    FERRITE_DURATION          example duration                       [ 1ns ... ] = 10s                   ✓ using default value
	//    FERRITE_ENUM              example enum                           [ foo | bar | baz ] = bar           ✓ using default value
	//    FERRITE_NETWORK_PORT      example network port                   [ <string> ] = 8080                 ✓ using default value
	//    FERRITE_NUM_FLOAT         example float-point                    [ <float32> ] = -123.45             ✓ using default value
	//    FERRITE_NUM_SIGNED        example signed integer                 [ <int16> ] = -123                  ✓ using default value
	//    FERRITE_NUM_UNSIGNED      example unsigned integer               [ <uint16> ] = 123                  ✓ using default value
	//    FERRITE_STRING            example string                         [ <string> ] = 'hello, world!'      ✓ using default value
	//    FERRITE_STRING_SENSITIVE  example sensitive string               [ <string> ] = *******              ✓ using default value
	//    FERRITE_SVC_SERVICE_HOST  kubernetes "ferrite-svc" service host  [ <string> ] = host.example.org     ✓ using default value
	//    FERRITE_SVC_SERVICE_PORT  kubernetes "ferrite-svc" service port  [ <string> ] = 443                  ✓ using default value
	//    FERRITE_URL               example URL                            [ <string> ] = https://example.org  ✓ using default value
	//  ❯ FERRITE_XTRIGGER          trigger failure for example              <string>                          ✗ undefined
	//
	// <process exited with error code 1>
}

func ExampleInit_validationWithOptionalValues() {
	defer example()()

	ferrite.
		Binary("FERRITE_BINARY", "example binary").
		Optional()

	ferrite.
		Binary("FERRITE_BINARY_SENSITIVE", "example sensitive binary").
		WithSensitiveContent().
		Optional()

	ferrite.
		Bool("FERRITE_BOOL", "example bool").
		Optional()

	ferrite.
		Duration("FERRITE_DURATION", "example duration").
		Optional()

	ferrite.
		Enum("FERRITE_ENUM", "example enum").
		WithMembers("foo", "bar", "baz").
		Optional()

	ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port").
		Optional()

	ferrite.
		Float[float32]("FERRITE_NUM_FLOAT", "example float-point").
		Optional()

	ferrite.
		Signed[int16]("FERRITE_NUM_SIGNED", "example signed integer").
		Optional()

	ferrite.
		Unsigned[uint16]("FERRITE_NUM_UNSIGNED", "example unsigned integer").
		Optional()

	ferrite.
		String("FERRITE_STRING", "example string").
		Optional()

	ferrite.
		String("FERRITE_STRING_SENSITIVE", "example sensitive string").
		WithSensitiveContent().
		Optional()

	ferrite.
		KubernetesService("ferrite-svc").
		Optional()

	ferrite.
		URL("FERRITE_URL", "example URL").
		Optional()

	ferrite.
		String("FERRITE_XTRIGGER", "trigger failure for example").
		Required()

	ferrite.Init()

	// Output:
	// Environment Variables:
	//
	//    FERRITE_BINARY            example binary                         [ <base64> ]         • undefined
	//    FERRITE_BINARY_SENSITIVE  example sensitive binary               [ <base64> ]         • undefined
	//    FERRITE_BOOL              example bool                           [ true | false ]     • undefined
	//    FERRITE_DURATION          example duration                       [ 1ns ... ]          • undefined
	//    FERRITE_ENUM              example enum                           [ foo | bar | baz ]  • undefined
	//    FERRITE_NETWORK_PORT      example network port                   [ <string> ]         • undefined
	//    FERRITE_NUM_FLOAT         example float-point                    [ <float32> ]        • undefined
	//    FERRITE_NUM_SIGNED        example signed integer                 [ <int16> ]          • undefined
	//    FERRITE_NUM_UNSIGNED      example unsigned integer               [ <uint16> ]         • undefined
	//    FERRITE_STRING            example string                         [ <string> ]         • undefined
	//    FERRITE_STRING_SENSITIVE  example sensitive string               [ <string> ]         • undefined
	//    FERRITE_SVC_SERVICE_HOST  kubernetes "ferrite-svc" service host  [ <string> ]         • undefined
	//    FERRITE_SVC_SERVICE_PORT  kubernetes "ferrite-svc" service port  [ <string> ]         • undefined
	//    FERRITE_URL               example URL                            [ <string> ]         • undefined
	//  ❯ FERRITE_XTRIGGER          trigger failure for example              <string>           ✗ undefined
	//
	// <process exited with error code 1>
}

func ExampleInit_validationWithNonCanonicalValues() {
	defer example()()

	os.Setenv("FERRITE_DURATION", "3h 10m 0s")
	ferrite.
		Duration("FERRITE_DURATION", "example duration").
		Required()

	ferrite.
		String("FERRITE_XTRIGGER", "trigger failure for example").
		Required()

	ferrite.Init()

	// Output:
	// Environment Variables:
	//
	//    FERRITE_DURATION  example duration               1ns ...     ✓ set to '3h 10m 0s', equivalent to 3h10m
	//  ❯ FERRITE_XTRIGGER  trigger failure for example    <string>    ✗ undefined
	//
	// <process exited with error code 1>
}

func ExampleInit_validationWithInvalidValues() {
	defer example()()

	os.Setenv("FERRITE_BINARY", "<invalid base64>")
	ferrite.
		Binary("FERRITE_BINARY", "example binary").
		Required()

	os.Setenv("FERRITE_BINARY_SENSITIVE", "<invalid base64>")
	ferrite.
		Binary("FERRITE_BINARY_SENSITIVE", "example sensitive binary").
		WithSensitiveContent().
		Required()

	os.Setenv("FERRITE_BOOL", "yes")
	ferrite.
		Bool("FERRITE_BOOL", "example bool").
		Required()

	os.Setenv("FERRITE_DURATION", "-+10s")
	ferrite.
		Duration("FERRITE_DURATION", "example duration").
		Required()

	os.Setenv("FERRITE_ENUM", "qux")
	ferrite.
		Enum("FERRITE_ENUM", "example enum").
		WithMembers("foo", "bar", "baz").
		Required()

	os.Setenv("FERRITE_NETWORK_PORT", "<invalid port>")
	ferrite.
		NetworkPort("FERRITE_NETWORK_PORT", "example network port").
		Required()

	os.Setenv("FERRITE_NUM_FLOAT", "-123w45")
	ferrite.
		Float[float32]("FERRITE_NUM_FLOAT", "example float-point").
		Required()

	os.Setenv("FERRITE_NUM_SIGNED", "123.3")
	ferrite.
		Signed[int16]("FERRITE_NUM_SIGNED", "example signed integer").
		Required()

	os.Setenv("FERRITE_NUM_UNSIGNED", "-123")
	ferrite.
		Unsigned[uint16]("FERRITE_NUM_UNSIGNED", "example unsigned integer").
		Required()

	os.Setenv("FERRITE_STRING", "foo bar")
	ferrite.
		String("FERRITE_STRING", "example string").
		WithConstraint(
			"must not contain whitespace",
			func(s string) bool {
				return !strings.ContainsRune(s, ' ')
			},
		).
		Required()

	os.Setenv("FERRITE_STRING_SENSITIVE", "foo bar")
	ferrite.
		String("FERRITE_STRING_SENSITIVE", "example sensitive string").
		WithConstraint(
			"must not contain whitespace",
			func(s string) bool {
				return !strings.ContainsRune(s, ' ')
			},
		).
		WithSensitiveContent().
		Required()

	os.Setenv("FERRITE_SVC_SERVICE_HOST", ".local")
	os.Setenv("FERRITE_SVC_SERVICE_PORT", "https-")
	ferrite.
		KubernetesService("ferrite-svc").
		Required()

	os.Setenv("FERRITE_URL", "/relative/path")
	ferrite.
		URL("FERRITE_URL", "example URL").
		Required()

	ferrite.Init()

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_BINARY            example binary                           <base64>           ✗ set to {16 bytes}, illegal base64 data at input byte 0
	//  ❯ FERRITE_BINARY_SENSITIVE  example sensitive binary                 <base64>           ✗ set to {16 bytes}, illegal base64 data at input byte 0
	//  ❯ FERRITE_BOOL              example bool                             true | false       ✗ set to yes, expected either true or false
	//  ❯ FERRITE_DURATION          example duration                         1ns ...            ✗ set to -+10s, expected duration
	//  ❯ FERRITE_ENUM              example enum                             foo | bar | baz    ✗ set to qux, expected foo, bar or baz
	//  ❯ FERRITE_NETWORK_PORT      example network port                     <string>           ✗ set to '<invalid port>', IANA service name must contain only ASCII letters, digits and hyphen
	//  ❯ FERRITE_NUM_FLOAT         example float-point                      <float32>          ✗ set to -123w45, expected float32
	//  ❯ FERRITE_NUM_SIGNED        example signed integer                   <int16>            ✗ set to 123.3, expected integer between -32768 and +32767
	//  ❯ FERRITE_NUM_UNSIGNED      example unsigned integer                 <uint16>           ✗ set to -123, expected integer between 0 and 65535
	//  ❯ FERRITE_STRING            example string                           <string>           ✗ set to 'foo bar', must not contain whitespace
	//  ❯ FERRITE_STRING_SENSITIVE  example sensitive string                 <string>           ✗ set to *******, must not contain whitespace
	//  ❯ FERRITE_SVC_SERVICE_HOST  kubernetes "ferrite-svc" service host    <string>           ✗ set to .local, host must not begin or end with a dot
	//  ❯ FERRITE_SVC_SERVICE_PORT  kubernetes "ferrite-svc" service port    <string>           ✗ set to https-, IANA service name must not begin or end with a hyphen
	//  ❯ FERRITE_URL               example URL                              <string>           ✗ set to /relative/path, URL must have a scheme
	//
	// <process exited with error code 1>
}
