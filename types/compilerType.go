package types

type CodeKind int

const (
	ROW CodeKind = iota
	VAR
)

type ArgumentKind int

const (
	STRING ArgumentKind = iota
	ARRAY
)

type Argument struct {
	Name string
	Value string
	Kind ArgumentKind
}

type Code struct {
	Code string
	Kind CodeKind
}