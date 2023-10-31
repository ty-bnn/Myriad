package values

type AddString struct {
	Kind   ValueKind
	Values []Value
}

func (a AddString) GetKind() ValueKind {
	return a.Kind
}

func (a AddString) GetName() string {
	return ""
}
