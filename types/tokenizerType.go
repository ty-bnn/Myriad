package types

type TokenKind int

const (
	SIMPORT TokenKind = iota
	SFROM
	SMAIN
	SLPAREN
	SRPAREN
	SCOMMA
	SARRANGE
	SLBRACE
	SRBRACE
	SSTRING
	SDFCOMMAND
	SDFARG
	SIDENTIFIER
	SASSIGNVARIABLE
)

type Token struct {
	Content string
	Kind TokenKind
	Line int
}