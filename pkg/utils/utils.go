package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ReadLinesFromFile(samplePath string) (string, error) {
	// Open file.
	fp, err := os.Open(samplePath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("cannot open %s", samplePath))
	}
	defer func() {
		if err := fp.Close(); err != nil {
			fmt.Println("cannot close file", err)
		}
	}()

	// Read sample codes line by line.
	data, err := io.ReadAll(fp)
	if err != nil {
		return "", errors.New(fmt.Sprintf("cannot read %s", samplePath))
	}

	return string(data), nil
}

func WriteFile(codes []string, filePath string) error {
	// Create file.
	dirs := filepath.Dir(filePath)
	if _, err := os.Stat(dirs); os.IsNotExist(err) {
		err := os.MkdirAll(dirs, os.ModePerm)
		if err != nil {
			return err
		}
	}

	fp, err := os.Create(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot create Dockerfile"))
	}
	defer func() {
		if err := fp.Close(); err != nil {
			fmt.Println("cannot close file", err)
		}
	}()

	for _, code := range codes {
		_, err := fp.Write([]byte(code))
		if err != nil {
			return errors.New(fmt.Sprintf("cannot write byte data to the file"))
		}
	}

	return nil
}

func WriteStdOut(codes []string) {
	fmt.Println("==================================")
	for _, code := range codes {
		fmt.Print(code)
	}
}
