package vars

type Var interface {
	GetKind() VarKind
	GetName() string
}

type VarKind int

const (
	SINGLE VarKind = iota
	ARRAY
	ELEMENT
	LITERAL
)
