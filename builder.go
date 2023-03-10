package ferrite

// isBuilderOf makes a static assertion that B implements the standard builder
// interface for a variable of type T.
type isBuilderOf[
	T any,
	B interface {
		Required(options ...VariableOption) Required[T]
		Optional(options ...VariableOption) Optional[T]
	},
] struct{}
