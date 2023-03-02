package generator

import (
	"fmt"
	"errors"

	"dcc/tokenizer"
	"dcc/compiler"
)

func getArgumentIndex(functionName string, variableName string) (int, bool) {
	for i, argument := range functionArgMap[functionName] {
		if argument.Name == variableName {
			return i, true
		}
	}

	return 0, false
}

func getIfCondition(ifContent compiler.IfContent, functionName string, argValues []string) (bool, error) {
	var lValue, rValue string
	if ifContent.LFormula.Kind == tokenizer.SSTRING {
		lValue = ifContent.LFormula.Content
	} else if ifContent.LFormula.Kind == tokenizer.SIDENTIFIER {
		if i, ok := getArgumentIndex(functionName, ifContent.LFormula.Content); ok {
			lValue = argValues[i]
		} else if val, ok := argsInMain[ifContent.LFormula.Content]; ok {
			lValue = val
		} else {
			return false, errors.New(fmt.Sprintf("semantic error: variable is not declared"))
		}
	}

	if ifContent.RFormula.Kind == tokenizer.SSTRING {
		rValue = ifContent.RFormula.Content
	} else if ifContent.RFormula.Kind == tokenizer.SIDENTIFIER {
		if i, ok := getArgumentIndex(functionName, ifContent.LFormula.Content); ok {
			rValue = argValues[i]
		} else if val, ok := argsInMain[ifContent.LFormula.Content]; ok {
			rValue = val
		} else {
			return false, errors.New(fmt.Sprintf("semantic error: variable is not declared"))
		}
	}

	if ifContent.Operator == compiler.EQUAL && lValue == rValue {
		return true, nil
	} else if ifContent.Operator == compiler.NOTEQUAL && lValue != rValue {
		return true, nil
	}

	return false, nil
}
