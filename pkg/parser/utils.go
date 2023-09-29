package parser

import (
	"errors"
	"fmt"

	"github.com/ty-bnn/myriad/pkg/model/codes"
)

func (p *Parser) isCompiled(file string) bool {
	for _, compiledFile := range p.compiledFiles {
		if file == compiledFile {
			return true
		}
	}

	return false
}

func (p *Parser) addFuncCodes(funcToCodes map[string][]codes.Code) error {
	for funcName, funcCodes := range funcToCodes {
		if _, has := p.FuncToCodes[funcName]; has {
			return errors.New(fmt.Sprintf("semantic error: %s is already declared", funcName))
		}

		p.FuncToCodes[funcName] = funcCodes
	}

	return nil
}
