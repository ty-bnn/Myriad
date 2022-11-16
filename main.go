package main

import(
	"fmt"
	"os"
	"bufio"
	"dcc/tokenizer"
)

func main() {
	samplePath := "sample/sample01.ty"

	lines := readLinesFromSample(samplePath)

	tokens, err := tokenizer.Tokenize(lines)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	// For debug.
	for _, token := range tokens {
		fmt.Printf("%30s\t%10s\n", token.Content, token.Kind)
	}
}

func readLinesFromSample(samplePath string) []string {
	var lines []string

	// Open file.
	fp, err := os.Open(samplePath)
	if err != nil {
		fmt.Println("cannot open", samplePath)
	}
	defer fp.Close()

	// Read sample code line by line.
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		fmt.Println("cannot read lines from", samplePath)
		os.Exit(0)
	}

	return lines
}