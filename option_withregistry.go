package ferrite

// WithRegistry is an option that sets the variable registries to use.
func WithRegistry(reg Registry) interface {
	InitOption
	RequiredOption
	OptionalOption
	DeprecatedOption
} {
	if reg == nil {
		panic("registry must not be nil")
	}

	return option{
		ApplyToInitConfig: func(cfg *initConfig) {
			cfg.ModeConfig.Registries.Add(reg)
		},
		ApplyToSetConfig: func(cfg *variableSetConfig) {
			cfg.Registries = append(cfg.Registries, reg)
		},
	}
}
