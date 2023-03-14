package compiler

import (
	"dcc/tokenizer"
)

// 中間言語
type InterCode struct {
	Content string
	Kind CodeKind
	ArgValues []string // for 関数呼び出し
	IfContent IfContent // for if文
}

type CodeKind int
const (
	ROW CodeKind = iota
	COMMAND
	VAR
	CALLFUNC
	IF
	ELIF
	ELSE
	ENDIF
)

type IfContent struct {
	LFormula, RFormula Formula
	Operator OperaterKind
	EndIndex int
}

type Formula struct {
	Content string
	Kind tokenizer.TokenKind
}

type OperaterKind int
const (
	EQUAL OperaterKind = iota
	NOTEQUAL
)

type ArgumentKind int
const (
	STRING ArgumentKind = iota
	ARRAY
)

type Argument struct {
	Name string
	Kind ArgumentKind
}
