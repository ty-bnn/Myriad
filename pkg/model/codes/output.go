package codes

import "github.com/ty-bnn/myriad/pkg/model/values"

type Output struct {
	Kind     CodeKind
	FilePath values.Value
}

func (o Output) GetKind() CodeKind {
	return o.Kind
}
