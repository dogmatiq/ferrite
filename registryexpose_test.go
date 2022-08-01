package ferrite

import (
	"io"
)

func SetExitBehavior(w io.Writer, fn func(code int)) {
	stderr = w
	exit = fn
}
