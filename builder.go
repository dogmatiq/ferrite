package ferrite

// isBuilderOf makes a static assertion that B has the common methods required
// for all "builder" types.
type isBuilderOf[T any, B interface {
	Required(options ...RequiredOption) Required[T]
	Optional(options ...OptionalOption) Optional[T]
	Deprecated(options ...DeprecatedOption) Deprecated[T]
}] struct{}
