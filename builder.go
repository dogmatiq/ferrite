package ferrite

// isBuilderOf makes a static assertion that B has the common methods required
// for all "builder" types.
type isBuilderOf[Result, Literal any, Builder interface {
	WithDefault(Literal) Builder
	// WithDefaultFrom(VariableSet[Result]) Builder
	WithExample(Literal, string) Builder

	Required(...RequiredOption) Required[Result]
	Optional(...OptionalOption) Optional[Result]
	Deprecated(...DeprecatedOption) Deprecated[Result]
}] struct{}

// isBuilderOf makes a static assertion that B has the common methods required
// for all "builder" types.
type isBuilderOfMinimal[Result any, Builder interface {
	Required(...RequiredOption) Required[Result]
	Optional(...OptionalOption) Optional[Result]
	Deprecated(...DeprecatedOption) Deprecated[Result]
}] struct{}
