package ferrite_test

import (
	"os"
	"time"

	"github.com/dogmatiq/ferrite"
)

func ExampleInit_markdownUsage() {
	setUp()
	defer tearDown()

	// Tell ferrite to generate markdown documentation for the environment
	// variables.
	os.Setenv("FERRITE_MODE", "usage/markdown")

	ferrite.
		Duration("FERRITE_DURATION", "example duration").
		WithDefault(10 * time.Minute).
		Required()

	ferrite.Init()

	// Output:
	// # Environment Variables
	//
	// This document describes the environment variables used by `ferrite.test`. It is generated automatically by [dogmatiq/ferrite].
	//
	// <!-- references -->
	//
	// [dogmatiq/ferrite]: https://github.com/dogmatiq/ferrite
}
