package ferrite_test

import (
	"os"
	"time"

	"github.com/dogmatiq/ferrite"
)

func ExampleValidateEnvironment_undefined() {
	setUp()
	defer tearDown()

	ferrite.
		Bool("FERRITE_BOOL", "example bool").
		Required()

	ferrite.
		Duration("FERRITE_DURATION", "example duration").
		Required()

	ferrite.
		Enum("FERRITE_ENUM", "example enum").
		WithMembers("foo", "bar").
		Required()

	ferrite.
		KubernetesService("ferrite-svc").
		Required()

	ferrite.
		Signed[int16]("FERRITE_NUM_SIGNED", "example signed integer").
		Required()

	ferrite.
		Unsigned[uint16]("FERRITE_NUM_UNSIGNED", "example unsigned integer").
		Required()

	ferrite.
		String("FERRITE_STRING", "example string").
		Required()

	ferrite.ValidateEnvironment()

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_BOOL              example bool                  true | false             ✗ undefined
	//  ❯ FERRITE_DURATION          example duration              1ns ...                  ✗ undefined
	//  ❯ FERRITE_ENUM              example enum                  foo | bar                ✗ undefined
	//  ❯ FERRITE_NUM_SIGNED        example signed integer        -32768 .. +32767         ✗ undefined
	//  ❯ FERRITE_NUM_UNSIGNED      example unsigned integer      0 .. 65535               ✗ undefined
	//  ❯ FERRITE_STRING            example string                <string>                 ✗ undefined
	//  ❯ FERRITE_SVC_SERVICE_HOST  "ferrite-svc" service host    <string>                 ✗ undefined
	//  ❯ FERRITE_SVC_SERVICE_PORT  "ferrite-svc" service port    <string> | 1 .. 65535    ✗ undefined
}

func ExampleValidateEnvironment_defaultValues() {
	setUp()
	defer tearDown()

	ferrite.
		Bool("FERRITE_BOOL", "example bool").
		WithDefault(true).
		Required()

	ferrite.
		Duration("FERRITE_DURATION", "example duration").
		WithDefault(10 * time.Second).
		Optional()

	ferrite.
		Enum("FERRITE_ENUM", "example enum").
		WithMembers("foo", "bar").
		WithDefault("bar").
		Required()

	ferrite.
		KubernetesService("ferrite-svc").
		WithDefault("host.example.org", "443").
		Optional()

	ferrite.
		Signed[int16]("FERRITE_NUM_SIGNED", "example signed integer").
		WithDefault(-123).
		Required()

	ferrite.
		Unsigned[uint16]("FERRITE_NUM_UNSIGNED", "example unsigned integer").
		WithDefault(123).
		Optional()

	ferrite.
		String("FERRITE_STRING", "example string").
		WithDefault("hello, world!").
		Required()

	ferrite.
		String("FERRITE_XTRIGGER", "trigger failure").
		Required()

	ferrite.ValidateEnvironment()

	// Output:
	// Environment Variables:
	//
	//    FERRITE_BOOL              example bool                [ true | false ] = true          ✓ using default value
	//    FERRITE_DURATION          example duration            [ 1ns ... ] = 10s                ✓ using default value
	//    FERRITE_ENUM              example enum                [ foo | bar ] = bar              ✓ using default value
	//    FERRITE_NUM_SIGNED        example signed integer      [ -32768 .. +32767 ] = -123      ✓ using default value
	//    FERRITE_NUM_UNSIGNED      example unsigned integer    [ 0 .. 65535 ] = 123             ✓ using default value
	//    FERRITE_STRING            example string              [ <string> ] = 'hello, world!'   ✓ using default value
	//    FERRITE_SVC_SERVICE_HOST  "ferrite-svc" service host  [ <string> ] = host.example.org  ✓ using default value
	//    FERRITE_SVC_SERVICE_PORT  "ferrite-svc" service port  [ <string> | 1 .. 65535 ] = 443  ✓ using default value
	//  ❯ FERRITE_XTRIGGER          trigger failure               <string>                       ✗ undefined
}

func ExampleValidateEnvironment_optional() {
	setUp()
	defer tearDown()

	ferrite.
		Bool("FERRITE_BOOL", "example bool").
		Optional()

	ferrite.
		Duration("FERRITE_DURATION", "example duration").
		Optional()

	ferrite.
		Enum("FERRITE_ENUM", "example enum").
		WithMembers("foo", "bar").
		Optional()

	ferrite.
		KubernetesService("ferrite-svc").
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
		String("FERRITE_XTRIGGER", "trigger failure").
		Required()

	ferrite.ValidateEnvironment()

	// Output:
	// Environment Variables:
	//
	//    FERRITE_BOOL              example bool                [ true | false ]           • undefined
	//    FERRITE_DURATION          example duration            [ 1ns ... ]                • undefined
	//    FERRITE_ENUM              example enum                [ foo | bar ]              • undefined
	//    FERRITE_NUM_SIGNED        example signed integer      [ -32768 .. +32767 ]       • undefined
	//    FERRITE_NUM_UNSIGNED      example unsigned integer    [ 0 .. 65535 ]             • undefined
	//    FERRITE_STRING            example string              [ <string> ]               • undefined
	//    FERRITE_SVC_SERVICE_HOST  "ferrite-svc" service host  [ <string> ]               • undefined
	//    FERRITE_SVC_SERVICE_PORT  "ferrite-svc" service port  [ <string> | 1 .. 65535 ]  • undefined
	//  ❯ FERRITE_XTRIGGER          trigger failure               <string>                 ✗ undefined
}

func ExampleValidateEnvironment_invalid() {
	setUp()
	defer tearDown()

	os.Setenv("FERRITE_BOOL", "yes")
	ferrite.
		Bool("FERRITE_BOOL", "example bool").
		Required()

	os.Setenv("FERRITE_DURATION", "-+10s")
	ferrite.
		Duration("FERRITE_DURATION", "example duration").
		Required()

	os.Setenv("FERRITE_ENUM", "baz")
	ferrite.
		Enum("FERRITE_ENUM", "example enum").
		WithMembers("foo", "bar").
		Required()

	os.Setenv("FERRITE_NUM_SIGNED", "123.3")
	ferrite.
		Signed[int16]("FERRITE_NUM_SIGNED", "example signed integer").
		Required()

	os.Setenv("FERRITE_NUM_UNSIGNED", "-123")
	ferrite.
		Unsigned[uint16]("FERRITE_NUM_UNSIGNED", "example unsigned integer").
		Required()

	// There is currently no way to make a string defined but invalid.
	ferrite.
		String("FERRITE_STRING", "example string").
		Required()

	os.Setenv("FERRITE_SVC_SERVICE_HOST", ".local")
	os.Setenv("FERRITE_SVC_SERVICE_PORT", "https-")
	ferrite.
		KubernetesService("ferrite-svc").
		Required()

	ferrite.ValidateEnvironment()

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_BOOL              example bool                  true | false             ✗ set to yes, must be either true or false
	//  ❯ FERRITE_DURATION          example duration              1ns ...                  ✗ set to -+10s, must be a valid duration, e.g. 10m30s
	//  ❯ FERRITE_ENUM              example enum                  foo | bar                ✗ set to baz, must be one of the enum members
	//  ❯ FERRITE_NUM_SIGNED        example signed integer        -32768 .. +32767         ✗ set to 123.3, must be an integer between -32768 and +32767
	//  ❯ FERRITE_NUM_UNSIGNED      example unsigned integer      0 .. 65535               ✗ set to -123, must be an integer between 0 and 65535
	//  ❯ FERRITE_STRING            example string                <string>                 ✗ undefined
	//  ❯ FERRITE_SVC_SERVICE_HOST  "ferrite-svc" service host    <string>                 ✗ set to .local, host must not begin or end with a dot
	//  ❯ FERRITE_SVC_SERVICE_PORT  "ferrite-svc" service port    <string> | 1 .. 65535    ✗ set to https-, IANA service name must not begin or end with a hyphen
}
