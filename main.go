package main

import(
	"fmt"
	"os"
	"errors"

	"dcc/types"
	"dcc/others"
	"dcc/tokenizer"
	"dcc/parser"
	"dcc/compiler"
)

func main() {
	lines, err := others.ReadLinesFromFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tokens, err := tokenizer.Tokenize(lines)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	// For debug.
	// for _, token := range tokens {
	// 	fmt.Printf("%30s\t%10d\n", token.Content, token.Kind)
	// }

	err = parser.Parse(tokens)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	functionCodeMap := map[string][]types.Code{}
	functionArgMap := map[string][]types.Argument{}

	err = compiler.Compile(tokens, &functionArgMap, &functionCodeMap)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, ok := functionCodeMap["main"]; !ok {
		fmt.Println(errors.New(fmt.Sprintf("syntax error: cannot find main function")))
		os.Exit(1)
	}

	err = others.WriteFile(functionCodeMap["main"])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
