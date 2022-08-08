package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleBool_required() {
	setUp()
	defer tearDown()

	v := ferrite.
		Bool("FERRITE_BOOL", "example boolean variable").
		Required()

	os.Setenv("FERRITE_BOOL", "true")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is true
}

func ExampleBool_default() {
	setUp()
	defer tearDown()

	v := ferrite.
		Bool("FERRITE_BOOL", "example boolean variable").
		WithDefault(true).
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is true
}

func ExampleBool_optional() {
	setUp()
	defer tearDown()

	v := ferrite.
		Bool("FERRITE_BOOL", "example boolean variable").
		Optional()

	ferrite.Init()

	if x, ok := v.Value(); ok {
		fmt.Println("value is", x)
	} else {
		fmt.Println("value is undefined")
	}

	// Output:
	// value is undefined
}

func ExampleBool_customLiterals() {
	setUp()
	defer tearDown()

	v := ferrite.
		Bool("FERRITE_BOOL", "example boolean variable").
		WithLiterals("yes", "no").
		Required()

	os.Setenv("FERRITE_BOOL", "yes")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is true
}
