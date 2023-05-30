package generator

import (
	"myriad/compiler"
)

type Generator struct {
	MainCodes *[]compiler.InterCode
	argsInMain *map[string]string
	command string
	Codes []string
}

func (g *Generator) GenerateCode(fInterCodeMap *map[string][]compiler.InterCode, fArgMap *map[string][]compiler.Variable) error {
	g.argsInMain = &map[string]string{}
	g.MainCodes = &[]compiler.InterCode{}
	*g.MainCodes = (*fInterCodeMap)["main"]
	g.command = ""

	_, codes, err := g.generateCodeBlock(0)
	if err != nil {
		return err
	}

	g.Codes = codes

	return nil
}
