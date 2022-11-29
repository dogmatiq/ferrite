package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleFloat_required() {
	setUp()
	defer tearDown()

	v := ferrite.
		Float[float64]("FERRITE_FLOAT", "example signed floating-point variable").
		Required()

	os.Setenv("FERRITE_FLOAT", "-123.45")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is -123.45
}

func ExampleFloat_default() {
	setUp()
	defer tearDown()

	v := ferrite.
		Float[float64]("FERRITE_FLOAT", "example signed floating-point variable").
		WithDefault(-123.45).
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is -123.45
}

func ExampleFloat_optional() {
	setUp()
	defer tearDown()

	v := ferrite.
		Float[float64]("FERRITE_FLOAT", "example signed floating-point variable").
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

func ExampleFloat_limits() {
	setUp()
	defer tearDown()

	v := ferrite.
		Float[float64]("FERRITE_FLOAT", "example signed floating-point variable").
		WithMinimum(-5).
		WithMaximum(10).
		Required()

	os.Setenv("FERRITE_FLOAT", "-2")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is -2
}
