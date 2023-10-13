package codes

import (
	"github.com/ty-bnn/myriad/pkg/model/values"
)

type CallProc struct {
	Kind     CodeKind
	ProcName string
	Args     []values.Value
}

func (i CallProc) GetKind() CodeKind {
	return i.Kind
}
