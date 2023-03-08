package ferrite_test

import (
	"errors"
	"os"

	"github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/ferrite/variable"
)

func ExampleSensitive() {
	setUp()
	defer tearDown()

	os.Setenv("FERRITE_SENSITIVE", "hunter2")
	ferrite.
		String("FERRITE_SENSITIVE", "example sensitive variable").
		WithConstraintFunc(
			"always fail",
			func(s string) variable.ConstraintError {
				return errors.New("always fail")
			},
		).
		Required(ferrite.Sensitive())

	ferrite.Init()

	// Note that the variable's value is obscured in the console output.

	// Output:
	// Environment Variables:
	//
	//  ❯ FERRITE_SENSITIVE  example sensitive variable    <string>    ✗ set to *******, always fail
}
