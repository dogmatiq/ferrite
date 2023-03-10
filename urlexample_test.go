package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleURL_required() {
	defer example()()

	v := ferrite.
		URL("FERRITE_URL", "example URL variable").
		Required()

	os.Setenv("FERRITE_URL", "https://example.org/path")
	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is https://example.org/path
}

func ExampleURL_default() {
	defer example()()

	v := ferrite.
		URL("FERRITE_URL", "example URL variable").
		WithDefault("https://example.org/default").
		Required()

	ferrite.Init()

	fmt.Println("value is", v.Value())

	// Output:
	// value is https://example.org/default
}

func ExampleURL_optional() {
	defer example()()

	v := ferrite.
		URL("FERRITE_URL", "example URL variable").
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
func ExampleURL_deprecated() {
	defer example()()

	os.Setenv("FERRITE_URL", "https://example.org/path")
	ferrite.
		URL("FERRITE_URL", "example URL variable").
		Deprecated()

	ferrite.Init()

	fmt.Println("<execution continues>")

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_URL  example URL variable  [ <string> ]  ⚠ deprecated variable set to https://example.org/path
	//
	// <execution continues>
}
