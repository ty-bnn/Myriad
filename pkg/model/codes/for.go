package codes

import "github.com/ty-bnn/myriad/pkg/model/values"

type For struct {
	Kind       CodeKind
	ItrName    string
	ArrayValue values.Value
}

func (f For) GetKind() CodeKind {
	return f.Kind
}
