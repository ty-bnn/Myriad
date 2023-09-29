package codes

type Code interface {
	GetKind() CodeKind
}

type CodeKind int

// TODO: 全ての型にそれぞれの構造体を作る
const (
	LITERAL CodeKind = iota
	COMMAND
	DEFINE
	REPLACE
	CALLPROC
	IF
	ELIF
	ELSE
	JUMP
	POP
)
