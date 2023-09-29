package codes

import (
	"github.com/ty-bnn/myriad/pkg/model/vars"
)

type If struct {
	Kind      CodeKind
	Condition Condition
}

func (i If) GetKind() CodeKind {
	return i.Kind
}

type Elif struct {
	Kind      CodeKind
	Condition Condition
}

func (e Elif) GetKind() CodeKind {
	return e.Kind
}

type Else struct {
	Kind CodeKind
}

func (e Else) GetKind() CodeKind {
	return e.Kind
}

type OpeKind int

type Condition struct {
	Left, Right vars.Var
	Operator    OpeKind
}

const (
	EQUAL OpeKind = iota
	NOTEQUAL
)

/*
conditionの評価時にfalseになっても再帰的にif文内の実行を続ける
もしfalseであれば返り値のコードを出力しなければ良い
IF:
	evaluate condition
	...
	any tasks
	...
ENDIF
ELIF:
	evaluate condition
	...
	any tasks
	...
ENDIF
ELSE:
	evaluate condition
	...
	any tasks
	...
ENDELSE
*/
