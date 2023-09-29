package codes

import (
	"github.com/ty-bnn/myriad/pkg/model/vars"
)

type CallProc struct {
	Kind     CodeKind
	ProcName string
	Args     []vars.Var
}

func (i CallProc) GetKind() CodeKind {
	return i.Kind
}
