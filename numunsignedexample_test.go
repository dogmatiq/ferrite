package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleUnsigned_required() {
	defer example()()

	v := ferrite.
		Unsigned[uint]("FERRITE_UNSIGNED", "example unsigned integer variable").
		Required()

	os.Setenv("FERRITE_UNSIGNED", "123")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is 123
}

func ExampleUnsigned_default() {
	defer example()()

	v := ferrite.
		Unsigned[uint]("FERRITE_UNSIGNED", "example unsigned integer variable").
		WithDefault(123).
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is 123
}

func ExampleUnsigned_optional() {
	defer example()()

	v := ferrite.
		Unsigned[uint]("FERRITE_UNSIGNED", "example unsigned integer variable").
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

func ExampleUnsigned_limits() {
	defer example()()

	v := ferrite.
		Unsigned[uint]("FERRITE_UNSIGNED", "example unsigned integer variable").
		WithMinimum(5).
		WithMaximum(10).
		Required()

	os.Setenv("FERRITE_UNSIGNED", "7")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is 7
}
