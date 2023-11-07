package token

type TokenKind int

const (
	IMPORT TokenKind = iota
	FROM
	MAIN
	IF
	ELIF
	ELSE
	FOR
	IN
	KEYS
	JSONUNMARSHAL
	LPAREN
	RPAREN
	COMMA
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET
	DOT
	DEFINE
	ASSIGN
	EQUAL
	NOTEQUAL
	AND
	OR
	LDOUBLEBRA
	RDOUBLEBRA
	DFBEGIN
	DFEND
	PLUS
	DOUBLELESS
	STRING
	DFCOMMAND
	DFARG
	IDENTIFIER
	NUMBER
)

var ReservedKeywords = map[string]TokenKind{
	"import":        IMPORT,
	"from":          FROM,
	"main":          MAIN,
	"if":            IF,
	"else if":       ELIF,
	"else":          ELSE,
	"for":           FOR,
	"in":            IN,
	"keys":          KEYS,
	"JsonUnmarshal": JSONUNMARSHAL,
}

var DockerfileCommands = map[string]bool{
	"ADD":         true,
	"ARG":         true,
	"CMD":         true,
	"COPY":        true,
	"ENTRYPOINT":  true,
	"ENV":         true,
	"EXPOSE":      true,
	"FROM":        true,
	"HEALTHCHECK": true,
	"LABEL":       true,
	"MAINTAINER":  true,
	"ONBUILD":     true,
	"RUN":         true,
	"SHELL":       true,
	"STOPSIGNAL":  true,
	"USER":        true,
	"VOLUME":      true,
	"WORKDIR":     true,
}

type Token struct {
	Content string
	Kind    TokenKind
	Line    int
}
