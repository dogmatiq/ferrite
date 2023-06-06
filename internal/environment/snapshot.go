package environment

// Snapshot is a snapshot of the environment.
type Snapshot struct {
	variables []variable
}

type variable struct {
	Name, Value string
}

// TakeSnapshot takes a snapshot of the variables within the environment.
func TakeSnapshot() *Snapshot {
	s := &Snapshot{}

	Range(func(name, value string) bool {
		s.variables = append(s.variables, variable{name, value})
		return true
	})

	return s
}

// RestoreSnapshot restores the environment to the state it was in when the
// given snapshot was taken.
func RestoreSnapshot(s *Snapshot) {
	Range(func(name, _ string) bool {
		Unset(name)
		return true
	})

	for _, v := range s.variables {
		Set(v.Name, v.Value)
	}
}
