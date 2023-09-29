package codes

import (
	"github.com/ty-bnn/myriad/pkg/model/vars"
)

type Define struct {
	Kind CodeKind
	Var  vars.Var
}

func (d Define) GetKind() CodeKind {
	return d.Kind
}
