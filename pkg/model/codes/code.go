package codes

type Code interface {
	GetKind() CodeKind
}

type CodeKind int

const (
	LITERAL CodeKind = iota
	COMMAND
	DEFINE
	ASSIGN
	REPLACE
	CALLPROC
	IF
	ELIF
	ELSE
	JUMP
	POP
)
