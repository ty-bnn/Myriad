package values

type Element struct {
	Kind  ValueKind
	Name  string
	Index int
}

func (e Element) GetKind() ValueKind {
	return e.Kind
}

func (e Element) GetName() string {
	return e.Name
}
