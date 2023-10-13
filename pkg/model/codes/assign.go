package codes

import "github.com/ty-bnn/myriad/pkg/model/values"

type Assign struct {
	Kind  CodeKind
	Key   string
	Value values.Value
}

func (a Assign) GetKind() CodeKind {
	return a.Kind
}
