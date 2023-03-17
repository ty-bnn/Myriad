package compiler

func isCompiled(file string) bool {
	for _, compiledFile := range readFiles {
		if file == compiledFile {
			return true
		}
	}

	return false
}

func getVariableIndex(functionName string, variableName string) (int, bool) {
	for i, variable := range functionVarMap[functionName] {
		if variable.Name == variableName && (functionName == "main" || variable.Kind == VARIABLE) {
			return i, true
		}
	}

	return 0, false
}

func getArgumentIndex(functionName string, argumentName string) (int, bool) {
	for i, argument := range functionVarMap[functionName] {
		if argument.Name == argumentName && functionName != "main" && argument.Kind == ARGUMENT {
			return i, true
		}
	}

	return 0, false
}