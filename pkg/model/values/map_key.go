package values

type MapKey struct {
	Kind ValueKind
	Name string
}

func (m MapKey) GetKind() ValueKind {
	return m.Kind
}

func (m MapKey) GetName() string {
	return m.Name
}
