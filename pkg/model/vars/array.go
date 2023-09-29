package vars

type Array struct {
	Kind   VarKind
	Name   string
	Values []string
}

func (a Array) GetKind() VarKind {
	return a.Kind
}

func (a Array) GetName() string {
	return a.Name
}
