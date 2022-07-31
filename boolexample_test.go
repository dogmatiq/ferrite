package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleBool_true() {
	os.Setenv("FERRITE_DEBUG", "true")

	debug := ferrite.
		Bool(
			"FERRITE_DEBUG",
			"enable debug logging",
		)

	ferrite.ResolveEnvironment()

	if debug.Value() {
		fmt.Println("DEBUG is on")
	} else {
		fmt.Println("DEBUG is off")
	}

	// Output:
	// DEBUG is on
}

func ExampleBool_false() {
	os.Setenv("FERRITE_DEBUG", "false")

	debug := ferrite.
		Bool(
			"FERRITE_DEBUG",
			"enable debug logging",
		)

	ferrite.ResolveEnvironment()

	if debug.Value() {
		fmt.Println("DEBUG is on")
	} else {
		fmt.Println("DEBUG is off")
	}

	// Output:
	// DEBUG is off
}

func ExampleBool_default() {
	os.Setenv("FERRITE_DEBUG", "")

	debug := ferrite.
		Bool(
			"FERRITE_DEBUG",
			"enable debug logging",
		).
		Default(true)

	ferrite.ResolveEnvironment()

	if debug.Value() {
		fmt.Println("DEBUG is on")
	} else {
		fmt.Println("DEBUG is off")
	}

	// Output:
	// DEBUG is on
}

func ExampleBool_optional() {
	os.Setenv("FERRITE_DEBUG", "")

	debug := ferrite.
		Bool(
			"FERRITE_DEBUG",
			"enable debug logging",
		).
		Optional()

	ferrite.ResolveEnvironment()

	if value, ok := debug.Value(); !ok {
		fmt.Println("DEBUG is empty/undefined")
	} else if value {
		fmt.Println("DEBUG is on")
	} else {
		fmt.Println("DEBUG is off")
	}

	// Output:
	// DEBUG is empty/undefined
}

func ExampleBool_required() {
	// Capture the error message from ResolveEnvironment() for testing, this
	// would not be done in production code.
	defer func() {
		fmt.Println(recover())
	}()

	debug := ferrite.
		Bool(
			"FERRITE_DEBUG",
			"enable debug logging",
		)

	ferrite.ResolveEnvironment()

	if debug.Value() {
		fmt.Println("DEBUG is on")
	} else {
		fmt.Println("DEBUG is off")
	}

	// Output:
	// ENVIRONMENT VARIABLES
	// 	✗ FERRITE_DEBUG [bool] (enable debug logging)
	// 		✗ must be set explicitly
	// 		✗ must be set to "true" or "false"
}
