package ferrite_test

import (
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleInit_markdownUsage() {
	defer example()()

	// Tell ferrite to generate markdown documentation for the environment
	// variables.
	os.Setenv("FERRITE_MODE", "usage/markdown")

	ferrite.Init()

	// In the interest of simplicity this example doesn't have any defined
	// environment variables, which is explained in the markdown output.

	// Output:
	// # Environment Variables
	//
	// This document describes the environment variables used by `ferrite.test`.
	//
	// **`ferrite.test` does not appear to use any environment variables.**
	//
	// > [!WARNING]
	// > This document only shows environment variables declared using [Ferrite].
	// > `ferrite.test` may consume other undocumented environment variables.
	//
	// <!-- references -->
	//
	// [ferrite]: https://github.com/dogmatiq/ferrite
	// <process exited successfully>
}
