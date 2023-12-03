package codes

import "github.com/ty-bnn/myriad/pkg/model/values"

type (
	If struct {
		Kind      CodeKind
		Condition ConditionalNode
		Jump
	}

	Elif struct {
		Kind      CodeKind
		Condition ConditionalNode
		Jump
	}

	Else struct {
		Kind CodeKind
	}

	NodeKind        int
	OperatorKind    int
	ConditionalNode struct {
		Operator    OperatorKind
		Var         values.Value
		Left, Right *ConditionalNode
	}
	Jump struct {
		True  int
		False int
	}
)

func (i If) GetKind() CodeKind {
	return i.Kind
}

func (e Elif) GetKind() CodeKind {
	return e.Kind
}

func (e Else) GetKind() CodeKind {
	return e.Kind
}

const (
	AND OperatorKind = iota
	OR
	EQUAL
	NOTEQUAL
	STARTWITH
	ENDWITH
)

var CompOperator = map[OperatorKind]bool{
	EQUAL:     true,
	NOTEQUAL:  true,
	STARTWITH: true,
	ENDWITH:   true,
}
