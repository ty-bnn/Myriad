package generator

import (
	"fmt"
	"myriad/compiler"
)

type Generator struct {
	MainCodes []compiler.InterCode
	command   string
	Codes     []string
}

func (g *Generator) GenerateCode(fInterCodeMap map[string][]compiler.InterCode) error {
	fmt.Println("Generating...")
	g.MainCodes = fInterCodeMap["main"]
	g.command = ""

	_, codes, err := g.generateCodeBlock(0)
	if err != nil {
		return err
	}

	g.Codes = codes

	fmt.Println("Generate Done.")

	return nil
}
