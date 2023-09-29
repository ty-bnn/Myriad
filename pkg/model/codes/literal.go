package codes

type Literal struct {
	Kind    CodeKind
	Content string
}

func (l Literal) GetKind() CodeKind {
	return l.Kind
}
