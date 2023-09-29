package vars

type Element struct {
	Kind  VarKind
	Name  string
	Index int
}

func (e Element) GetKind() VarKind {
	return e.Kind
}

func (e Element) GetName() string {
	return e.Name
}
