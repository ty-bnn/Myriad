package vars

type Literal struct {
	Kind  VarKind
	Value string
}

func (l Literal) GetKind() VarKind {
	return l.Kind
}

func (l Literal) GetName() string {
	return ""
}
