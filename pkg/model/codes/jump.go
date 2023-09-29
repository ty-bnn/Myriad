package codes

type Jump struct {
	Kind       CodeKind
	NextOffset int
}

func (j Jump) GetKind() CodeKind {
	return j.Kind
}
