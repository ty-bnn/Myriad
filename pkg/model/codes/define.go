package codes

import (
	"github.com/ty-bnn/myriad/pkg/model/values"
)

type Define struct {
	Kind  CodeKind
	Key   string
	Value values.Value
}

func (d Define) GetKind() CodeKind {
	return d.Kind
}
