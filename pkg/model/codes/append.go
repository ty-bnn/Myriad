package codes

import "github.com/ty-bnn/myriad/pkg/model/values"

type Append struct {
	Kind    CodeKind
	Array   string
	Element values.Value
}

func (a Append) GetKind() CodeKind {
	return a.Kind
}
