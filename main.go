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
	for i, token := range tokens {
		if token.Content == "\n" {
			fmt.Printf("%2d\t%30s\t%10d\n", i, "\\n",token.Kind)
		} else {
			fmt.Printf("%2d\t%30s\t%10d\n", i, token.Content, token.Kind)
		}
	}

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

	codes := generator.GenerateCode(functionInterCodeMap, functionArgMap)

	err = others.WriteFile(codes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
