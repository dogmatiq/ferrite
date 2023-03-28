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
	// **There do not appear to be any environment variables.**
	//
	// ⚠️ `ferrite.test` may consume other undocumented environment variables. This
	// document only shows variables declared using [Ferrite].
	//
	// <!-- references -->
	//
	// [ferrite]: https://github.com/dogmatiq/ferrite
	// <process exited successfully>
}
