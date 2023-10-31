package values

type Value interface {
	GetKind() ValueKind
	GetName() string
}

type ValueKind int

const (
	LITERAL ValueKind = iota
	LITERALS
	ELEMENT
	MAP
	MAPKEY
	MAPVALUE
	IDENT
	ADDSTRING
)
