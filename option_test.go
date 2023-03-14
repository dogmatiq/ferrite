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

	os.Setenv("FERRITE_WIDGET_SPEED", "100")    // define the speed
	os.Setenv("FERRITE_WIDGET_ENABLED", "true") // and enable the widget
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

	widgetSpeed := ferrite.
		Unsigned[uint]("FERRITE_WIDGET_SPEED", "set the speed of the widget").
		Optional(ferrite.RelevantIf(widgetEnabled))

	os.Setenv("FERRITE_WIDGET_SPEED", "100")     // define the speed
	os.Setenv("FERRITE_WIDGET_ENABLED", "false") // but disable the widget
	ferrite.Init()

	if x, ok := widgetSpeed.Value(); ok {
		fmt.Println("value is", x)
	} else {
		fmt.Println("value is not relevant")
	}

	// Output:
	// value is not relevant
}
