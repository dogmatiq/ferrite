package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleBool() {
	setUp()
	defer tearDown()

	value := ferrite.
		Bool("FERRITE_BOOL", "example boolean variable")

	os.Setenv("FERRITE_BOOL", "true")
	ferrite.ValidateEnvironment()

	fmt.Println("value is", value.Value())

	// Output:
	// value is true
}

func ExampleBool_default() {
	setUp()
	defer tearDown()

	value := ferrite.
		Bool("FERRITE_BOOL", "example boolean variable").
		WithDefault(true)

	ferrite.ValidateEnvironment()

	fmt.Println("value is", value.Value())

	// Output:
	// value is true
}

func ExampleBool_customLiterals() {
	setUp()
	defer tearDown()

	value := ferrite.
		Bool("FERRITE_BOOL", "example boolean variable").
		WithLiterals("yes", "no")

	os.Setenv("FERRITE_BOOL", "yes")
	ferrite.ValidateEnvironment()

	fmt.Println("value is", value.Value())

	// Output:
	// value is true
}
