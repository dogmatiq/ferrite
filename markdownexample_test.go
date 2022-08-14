package ferrite_test

import (
	"os"

	"github.com/dogmatiq/ferrite"
)

func ExampleInit_markdownUsage() {
	setUp()
	defer tearDown()

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
	// The application may consume other undocumented environment variables; this
	// document only shows those variables defined using [dogmatiq/ferrite].
	//
	// <!-- references -->
	//
	// [dogmatiq/ferrite]: https://github.com/dogmatiq/ferrite
}
