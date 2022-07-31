package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleBool() {
	ferrite.DefaultRegistry.Reset()
	os.Setenv("FERRITE_BOOL", "true")
	defer os.Unsetenv("FERRITE_BOOL")

	debug := ferrite.
		Bool(
			"FERRITE_BOOL",
			"example boolean variable",
		)

	ferrite.ValidateEnvironment()

	if debug.Value() {
		fmt.Println("variable is true")
	} else {
		fmt.Println("variable is false")
	}

	// Output:
	// variable is true
}

func ExampleBool_default() {
	ferrite.DefaultRegistry.Reset()

	debug := ferrite.
		Bool(
			"FERRITE_BOOL",
			"example boolean variable",
		).
		Default(true)

	ferrite.ValidateEnvironment()

	if debug.Value() {
		fmt.Println("variable is true")
	} else {
		fmt.Println("variable is false")
	}

	// Output:
	// variable is true
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

	ferrite.ValidateEnvironment()

	if debug.Value() {
		fmt.Println("variable is true")
	} else {
		fmt.Println("variable is false")
	}

	// Output:
	// variable is true
}
