package codes

type Pop struct {
	Kind CodeKind
}

func (p Pop) GetKind() CodeKind {
	return p.Kind
}
