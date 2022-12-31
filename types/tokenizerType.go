package types

type TokenKind int
const (
	SIMPORT TokenKind = iota
	SFROM
	SMAIN
	SIF
	SELIF
	SELSE
	SLPAREN
	SRPAREN
	SCOMMA
	SARRANGE
	SLBRACE
	SRBRACE
	SEQUAL
	SNOTEQUAL
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