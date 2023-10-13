package codes

import (
	"github.com/ty-bnn/myriad/pkg/model/values"
)

type Replace struct {
	Kind  CodeKind
	Value values.Value
}

func (r Replace) GetKind() CodeKind {
	return r.Kind
}
