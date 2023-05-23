package compiler

func (c Compiler) isCompiled(file string) bool {
	for _, compiledFile := range *c.readFiles {
		if file == compiledFile {
			return true
		}
	}

	return false
}

func (c Compiler) getVariableIndex(functionName string, variableName string) (int, bool) {
	for i, variable := range (*c.FunctionVarMap)[functionName] {
		if variable.Name == variableName && (functionName == "main" || variable.Kind == VARIABLE) {
			return i, true
		}
	}

	return 0, false
}

func (c Compiler) getArgumentIndex(functionName string, argumentName string) (int, bool) {
	for i, argument := range (*c.FunctionVarMap)[functionName] {
		if argument.Name == argumentName && functionName != "main" && argument.Kind == ARGUMENT {
			return i, true
		}
	}

	return 0, false
}