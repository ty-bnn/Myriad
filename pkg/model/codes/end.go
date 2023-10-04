package codes

type End struct {
	Kind CodeKind
}

func (e End) GetKind() CodeKind {
	return e.Kind
}
