package codes

import "github.com/ty-bnn/myriad/pkg/model/vars"

type Assign struct {
	Kind CodeKind
	Var  vars.Var
}

func (a Assign) GetKind() CodeKind {
	return a.Kind
}
