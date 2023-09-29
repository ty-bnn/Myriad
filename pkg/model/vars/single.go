package vars

type Single struct {
	Kind  VarKind
	Name  string
	Value string
}

func (s Single) GetKind() VarKind {
	return s.Kind
}

func (s Single) GetName() string {
	return s.Name
}
