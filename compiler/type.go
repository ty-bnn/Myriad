package compiler

import (
	"fmt"

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

type VariableKind int
const (
	VARIABLE VariableKind = iota
	ARGUMENT
)

type Variable struct {
	Name string
	Value string
	Kind VariableKind
}

func printInterCodes(functionName string) {
	intToString := map[CodeKind]string{
		ROW: "row",
		COMMAND: "command",
		VAR: "var",
		IF: "if",
		ELIF: "elif",
		ELSE: "else",
		ENDIF: "endif",
	}

	fmt.Printf("--------- inter codes in \"%s\"---------\n", functionName)
	for i, code := range functionInterCodeMap[functionName] {
		fmt.Println("{")
		fmt.Printf("  Content: %s\n", code.Content)
		fmt.Printf("     Kind: %s\n", intToString[code.Kind])
		fmt.Printf("    Index: %2d\n", i)
		fmt.Println("},")
	}
}