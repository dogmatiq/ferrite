package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
	"github.com/dogmatiq/ferrite/is"
)

func ExampleSeeAlso() {
	defer example()()

	verbose := ferrite.
		Bool("FERRITE_VERBOSE", "enable verbose logging").
		Optional()

	ferrite.
		Bool("FERRITE_DEBUG", "enable or disable debugging features").
		Optional(ferrite.SeeAlso(verbose))

	ferrite.Init()

	// Output:
}

func ExampleSupersededBy() {
	defer example()()

	verbose := ferrite.
		Bool("FERRITE_VERBOSE", "enable verbose logging").
		Optional()

	ferrite.
		Bool("FERRITE_DEBUG", "enable debug logging").
		Deprecated(ferrite.SupersededBy(verbose))

	ferrite.Init()

	// Output:
}

func ExampleIgnoreIf() {
	defer example()()

	enabled := ferrite.
		Bool("FERRITE_WIDGET_ENABLED", "enable the widget").
		Required()

	speed := ferrite.
		Unsigned[uint]("FERRITE_WIDGET_SPEED", "set the speed of the widget").
		Optional(ferrite.IgnoreIf(enabled, is.Equal(false)))

	os.Setenv("FERRITE_WIDGET_SPEED", "100")     // define the speed
	os.Setenv("FERRITE_WIDGET_ENABLED", "false") // but disable the widget
	ferrite.Init()

	if x, ok := speed.Value(); ok {
		fmt.Println("value is", x)
	} else {
		fmt.Println("value is ignored")
	}

	// Output:
	// value is ignored
}
