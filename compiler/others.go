package compiler

func isCompiled(file string) bool {
	for _, compiledFile := range readFiles {
		if file == compiledFile {
			return true
		}
	}

	return false
}