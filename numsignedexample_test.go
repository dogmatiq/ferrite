package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleSigned_required() {
	defer example()()

	v := ferrite.
		Signed[int]("FERRITE_SIGNED", "example signed integer variable").
		Required()

	os.Setenv("FERRITE_SIGNED", "-123")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is -123
}

func ExampleSigned_default() {
	defer example()()

	v := ferrite.
		Signed[int]("FERRITE_SIGNED", "example signed integer variable").
		WithDefault(-123).
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is -123
}

func ExampleSigned_optional() {
	defer example()()

	v := ferrite.
		Signed[int]("FERRITE_SIGNED", "example signed integer variable").
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

func ExampleSigned_limits() {
	defer example()()

	v := ferrite.
		Signed[int]("FERRITE_SIGNED", "example signed integer variable").
		WithMinimum(-5).
		WithMaximum(10).
		Required()

	os.Setenv("FERRITE_SIGNED", "-2")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is -2
}

func ExampleSigned_deprecated() {
	defer example()()

	os.Setenv("FERRITE_SIGNED", "-123")
	v := ferrite.
		Signed[int]("FERRITE_SIGNED", "example signed integer variable").
		Deprecated()

	ferrite.Init()

	if x, ok := v.DeprecatedValue(); ok {
		fmt.Println("value is", x)
	} else {
		fmt.Println("value is undefined")
	}

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_SIGNED  example signed integer variable  [ <int> ]  ⚠ deprecated variable set to -123
	//
	// value is -123
}
