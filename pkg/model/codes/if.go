package codes

import "github.com/ty-bnn/myriad/pkg/model/values"

type (
	If struct {
		Kind      CodeKind
		Condition ConditionalNode
	}

	Elif struct {
		Kind      CodeKind
		Condition ConditionalNode
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
