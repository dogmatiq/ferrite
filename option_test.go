package ferrite_test

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite"
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

func ExampleRelevantIf_whenRelevant() {
	defer example()()

	widgetEnabled := ferrite.
		Bool("FERRITE_WIDGET_ENABLED", "enable the widget").
		Required()

	widgetSpeed := ferrite.
		Unsigned[uint]("FERRITE_WIDGET_SPEED", "set the speed of the widget").
		Optional(ferrite.RelevantIf(widgetEnabled))

	os.Setenv("FERRITE_WIDGET_SPEED", "100")
	os.Setenv("FERRITE_WIDGET_ENABLED", "true")
	ferrite.Init()

	if x, ok := widgetSpeed.Value(); ok {
		fmt.Println("value is", x)
	} else {
		fmt.Println("value is not relevant")
	}

	// Output:
	// value is 100
}

func ExampleRelevantIf_whenNotRelevant() {
	defer example()()

	widgetEnabled := ferrite.
		Bool("FERRITE_WIDGET_ENABLED", "enable the widget").
		Required()

	ferrite.
		Unsigned[uint]("FERRITE_WIDGET_SPEED", "set the speed of the widget").
		Required(ferrite.RelevantIf(widgetEnabled))

	// FERRITE_WIDGET_SPEED is "required" but we can leave it undefined when
	// FERRITE_WIDGET_ENABLED is "false" without causing a validation failure.
	os.Setenv("FERRITE_WIDGET_ENABLED", "false")
	ferrite.Init()

	// Output:
}
