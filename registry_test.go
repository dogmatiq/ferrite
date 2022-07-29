package ferrite

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup the default registry such that it panics rather than calling
	// os.Exit() so that tests can recover.
	//
	// In production setups we don't use a panic because we don't want to see a
	// stack trace.
	DefaultRegistry.fatal = func(err error) {
		panic(err)
	}

	os.Exit(m.Run())
}
