package main

import(
	"fmt"
	"os"

	"dcc/others"
	"dcc/tokenizer"
	"dcc/parser"
	"dcc/compiler"
	"dcc/generator"
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
	
	functionInterCodeMap, functionArgMap, err := compiler.Compile(tokens)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// For debug
	// for k, v := range functionInterCodeMap {
	// 	fmt.Println(k)
	// 	for _, c := range v {
	// 		fmt.Println(c)
	// 	}
	// }

	codes, err := generator.GenerateCode(functionInterCodeMap, functionArgMap)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = others.WriteFile(codes, os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
