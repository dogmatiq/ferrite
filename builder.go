package ferrite

import (
	"fmt"

	"github.com/dogmatiq/ferrite/variable"
)

// undefinedError returns an error that indicates that a variable is undefined.
func undefinedError(v variable.Any) error {
	return fmt.Errorf(
		"%s is undefined and does not have a default value",
		v.Spec().Name(),
	)
}

// isBuilderOf makes a static assertion that B meats
type isBuilderOf[T any, B builderOf[T]] struct{}

// builderOf is an interface and type constriant common to all builders that
// produce a value of type T.
type builderOf[T any] interface {
	Required(options ...RequiredOption) Required[T]
	Optional(options ...OptionalOption) Optional[T]
	Deprecated(options ...DeprecatedOption) Deprecated[T]
}
