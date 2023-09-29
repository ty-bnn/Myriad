package codes

type Command struct {
	Kind    CodeKind
	Content string
}

func (c Command) GetKind() CodeKind {
	return c.Kind
}
