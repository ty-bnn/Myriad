package generator

import (
	"dcc/compiler"
)

func getIfCondition(ifContent compiler.IfContent, functionName string) (bool, error) {
	var lValue, rValue string
	lValue = ifContent.LFormula.Content
	rValue = ifContent.RFormula.Content

	if ifContent.Operator == compiler.EQUAL && lValue == rValue {
		return true, nil
	} else if ifContent.Operator == compiler.NOTEQUAL && lValue != rValue {
		return true, nil
	}

	return false, nil
}
