package main

import(
	"fmt"
	"os"

	"dcc/helpers"
	"dcc/tokenizer"
	"dcc/parser"
	"dcc/compiler"
	"dcc/generator"
)

func main() {
	// Myriadファイルから全ての行を読む
	lines, err := helpers.ReadLinesFromFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// トークナイズ
	t := &tokenizer.Tokenizer{}
	err = t.Tokenize(lines)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// パース
	p := &parser.Parser{}
	err = p.Parse(t.Tokens)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	// コンパイル
	c := &compiler.Compiler{}
	err = c.Compile(t.Tokens)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// For debug
	// for k, v := range *c.FunctionInterCodeMap {
	// 	fmt.Println(k)
	// 	for _, c := range v {
	// 		fmt.Println(c)
	// 	}
	// }

	// コード生成
	g := &generator.Generator{}
	err = g.GenerateCode(c.FunctionInterCodeMap, c.FunctionVarMap)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = helpers.WriteFile(*g.Codes, os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
