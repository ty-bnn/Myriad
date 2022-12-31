package generator

import (
	"dcc/types"
)

func getArgumentIndex(functionName string, variableName string) int {
	for i, argument := range functionArgMap[functionName] {
		if argument.Name == variableName {
			return i
		}
	}

	// TODO: argumentExitを合体させたい
	return 0
}

func getIfCondition(ifContent types.IfContent, functionName string, argValues []string) bool {
	var lValue, rValue string
	if ifContent.LFormula.Kind == types.SSTRING {
		lValue = ifContent.LFormula.Content
	} else if ifContent.LFormula.Kind == types.SIDENTIFIER {
		i := getArgumentIndex(functionName, ifContent.LFormula.Content)
		lValue = argValues[i]
	}

	if ifContent.RFormula.Kind == types.SSTRING {
		rValue = ifContent.RFormula.Content
	} else if ifContent.RFormula.Kind == types.SIDENTIFIER {
		i := getArgumentIndex(functionName, ifContent.RFormula.Content)
		rValue = argValues[i]
	}

	if ifContent.Operator == types.EQUAL && lValue == rValue {
		return true
	} else if ifContent.Operator == types.NOTEQUAL && lValue != rValue {
		return true
	}

	return false
}