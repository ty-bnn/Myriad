package others

import (
	"os"
	"fmt"
	"bufio"
	"errors"

	"dcc/types"
)

func ReadLinesFromFile(samplePath string) ([]string, error) {
	var lines []string

	// Open file.
	fp, err := os.Open(samplePath)
	if err != nil {
		fp.Close()
		return []string{}, errors.New(fmt.Sprintf("cannot open %s", samplePath))
	}
	defer fp.Close()

	// Read sample code line by line.
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return []string{}, errors.New(fmt.Sprintf("cannot read %s", samplePath))
	}

	return lines, nil
}

func WriteFile(contents []types.Code) (error) {
	// Creae file.
	fp, err := os.Create("Dockerfile")
	if err != nil {
		fp.Close()
		return errors.New(fmt.Sprintf("cannot create Dockerfile"))
	}
	defer fp.Close()

	for _, content := range contents {
		_, err := fp.Write([]byte(content.Code))
		if err != nil {
			return errors.New(fmt.Sprintf("cannot write byte data to the file"))
		}
	}

	return nil
}
