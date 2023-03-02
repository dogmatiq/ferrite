package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleURL_required() {
	setUp()
	defer tearDown()

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
	setUp()
	defer tearDown()

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
	setUp()
	defer tearDown()

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
