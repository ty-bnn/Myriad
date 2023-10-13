package values

type Literal struct {
	Kind  ValueKind
	Value string
}

func (l Literal) GetKind() ValueKind {
	return l.Kind
}

func (l Literal) GetName() string {
	return ""
}
