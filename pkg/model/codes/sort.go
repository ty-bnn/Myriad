package codes

type Sort struct {
	Kind  CodeKind
	Array string
}

func (s Sort) GetKind() CodeKind {
	return s.Kind
}
