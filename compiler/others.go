package compiler

func argumentExist(functionName string, variableName string) bool {
	for _, argument := range functionArgMap[functionName] {
		if argument.Name == variableName {
			return true
		}
	}

	return false
}
