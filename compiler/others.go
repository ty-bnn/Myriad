package compiler

func argumentExist(functionName string, variableName string) bool {
	for _, argument := range (*functionArgMap)[functionName] {
		if argument.Name == variableName {
			return true
		}
	}

	return false
}

func getArgumentValue(functionName string, variableName string) string {
	for _, argument := range (*functionArgMap)[functionName] {
		if argument.Name == variableName {
			return argument.Value
		}
	}

	return ""
}