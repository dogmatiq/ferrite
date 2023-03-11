package ferrite_test

import (
	"errors"
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/ferrite/variable"
)

func ExampleString_required() {
	defer example()()

	v := ferrite.
		String("FERRITE_STRING", "example string variable").
		Required()

	os.Setenv("FERRITE_STRING", "<value>")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is <value>
}

func ExampleString_default() {
	defer example()()

	v := ferrite.
		String("FERRITE_STRING", "example string variable").
		WithDefault("<default>").
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is <default>
}

func ExampleString_optional() {
	defer example()()

	v := ferrite.
		String("FERRITE_STRING", "example string variable").
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

func ExampleString_sensitive() {
	defer example()()

	os.Setenv("FERRITE_STRING", "hunter2")
	ferrite.
		String("FERRITE_STRING", "example sensitive string variable").
		WithConstraintFunc(
			"always fail",
			func(s string) variable.ConstraintError {
				return errors.New("always fail")
			},
		).
		WithSensitiveContent().
		Required()

	ferrite.Init()

	// Note that the variable's value is obscured in the console output.

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_STRING  example sensitive string variable    <string>    ✗ set to *******, always fail
	//
	// <process exited with error code 1>
}

func ExampleString_deprecated() {
	defer example()()

	os.Setenv("FERRITE_STRING", "<value>")
	v := ferrite.
		String("FERRITE_STRING", "example string variable").
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
	//  ❯ FERRITE_STRING  example string variable  [ <string> ]  ⚠ deprecated variable set to '<value>'
	//
	// value is <value>
}
