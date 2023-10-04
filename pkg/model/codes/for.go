package codes

import "github.com/ty-bnn/myriad/pkg/model/vars"

type For struct {
	Kind  CodeKind
	Itr   vars.Single
	Array vars.Array
}

func (f For) GetKind() CodeKind {
	return f.Kind
}
