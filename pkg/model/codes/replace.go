package codes

import (
	"github.com/ty-bnn/myriad/pkg/model/vars"
)

type Replace struct {
	Kind   CodeKind
	RepVar vars.Var
}

func (r Replace) GetKind() CodeKind {
	return r.Kind
}
