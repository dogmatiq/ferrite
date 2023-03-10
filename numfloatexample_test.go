package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleFloat_required() {
	defer example()()

	v := ferrite.
		Float[float64]("FERRITE_FLOAT", "example floating-point variable").
		Required()

	os.Setenv("FERRITE_FLOAT", "-123.45")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is -123.45
}

func ExampleFloat_default() {
	defer example()()

	v := ferrite.
		Float[float64]("FERRITE_FLOAT", "example floating-point variable").
		WithDefault(-123.45).
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is -123.45
}

func ExampleFloat_optional() {
	defer example()()

	v := ferrite.
		Float[float64]("FERRITE_FLOAT", "example floating-point variable").
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
	defer example()()

	v := ferrite.
		Float[float64]("FERRITE_FLOAT", "example floating-point variable").
		WithMinimum(-5).
		WithMaximum(10).
		Required()

	os.Setenv("FERRITE_FLOAT", "-2")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is -2
}

func ExampleFloat_deprecated() {
	defer example()()

	os.Setenv("FERRITE_FLOAT", "-123.45")
	ferrite.
		Float[float64]("FERRITE_FLOAT", "example floating-point variable").
		Deprecated()

	ferrite.Init()

	fmt.Println("<execution continues>")

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_FLOAT  example floating-point variable  [ <float64> ]  ⚠ deprecated variable set to -123.45
	//
	// <execution continues>
}
