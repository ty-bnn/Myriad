package values

type Literals struct {
	Kind   ValueKind
	Values []string
}

func (l Literals) GetKind() ValueKind {
	return l.Kind
}

func (l Literals) GetName() string {
	return ""
}
