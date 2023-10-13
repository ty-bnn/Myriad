package values

type Map struct {
	Kind  ValueKind
	Value map[string]interface{}
}

func (m Map) GetKind() ValueKind {
	return m.Kind
}

func (m Map) GetName() string {
	return ""
}
