package compiler

func (c Compiler) isCompiled(file string) bool {
	for _, compiledFile := range c.readFiles {
		if file == compiledFile {
			return true
		}
	}

	return false
}

/*
mainに入っている場合は引数（ARGUMENT）でも通常の変数として扱う
Why: 予めコンパイル時に引数として与えられるため既に値がわかっているから
*/
func (c Compiler) getVariableIndex(functionName string, variableName string) (int, bool) {
	for i, variable := range c.FunctionVarMap[functionName] {
		if variable.getName() == variableName && (functionName == "main" || variable.getKind() == VARIABLE) {
			return i, true
		}
	}

	return 0, false
}

func (c Compiler) getArgumentIndex(functionName string, argumentName string) (int, bool) {
	for i, argument := range c.FunctionVarMap[functionName] {
		if argument.getName() == argumentName && functionName != "main" && argument.getKind() == ARGUMENT {
			return i, true
		}
	}

	return 0, false
}