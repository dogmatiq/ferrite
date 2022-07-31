package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleBool_true() {
	ferrite.DefaultRegistry.Reset()
	os.Setenv("FERRITE_BOOL", "true")
	defer os.Unsetenv("FERRITE_BOOL")

	debug := ferrite.
		Bool(
			"FERRITE_BOOL",
			"example boolean variable",
		)

	ferrite.ResolveEnvironment()

	if debug.Value() {
		fmt.Println("variable is true")
	} else {
		fmt.Println("variable is false")
	}

	// Output:
	// variable is true
}

func ExampleBool_false() {
	ferrite.DefaultRegistry.Reset()
	os.Setenv("FERRITE_BOOL", "false")
	defer os.Unsetenv("FERRITE_BOOL")

	debug := ferrite.
		Bool(
			"FERRITE_BOOL",
			"example boolean variable",
		)

	ferrite.ResolveEnvironment()

	if debug.Value() {
		fmt.Println("variable is true")
	} else {
		fmt.Println("variable is false")
	}

	// Output:
	// variable is false
}

func ExampleBool_default() {
	ferrite.DefaultRegistry.Reset()

	debug := ferrite.
		Bool(
			"FERRITE_BOOL",
			"example boolean variable",
		).
		Default(true)

	ferrite.ResolveEnvironment()

	if debug.Value() {
		fmt.Println("variable is true")
	} else {
		fmt.Println("variable is false")
	}

	// Output:
	// variable is true
}

func ExampleBool_undefined() {
	ferrite.DefaultRegistry.Reset()

	// Capture the error message from ResolveEnvironment() for testing, this
	// would not be done in production code.
	defer func() {
		fmt.Println(recover())
	}()

	ferrite.
		Bool(
			"FERRITE_BOOL",
			"example boolean variable",
		)

	ferrite.ResolveEnvironment()

	// Output:
	// ENVIRONMENT VARIABLES
	//  ✗ FERRITE_BOOL [bool] (example boolean variable)
	//    ✗ must be set explicitly
	//    ✗ must be either "true" or "false"
}

func ExampleBool_customLiterals() {
	ferrite.DefaultRegistry.Reset()
	os.Setenv("FERRITE_BOOL", "yes")
	defer os.Unsetenv("FERRITE_BOOL")

	debug := ferrite.
		Bool(
			"FERRITE_BOOL",
			"example boolean variable",
		).
		Literals("yes", "no")

	ferrite.ResolveEnvironment()

	if debug.Value() {
		fmt.Println("variable is true")
	} else {
		fmt.Println("variable is false")
	}

	// Output:
	// variable is true
}
