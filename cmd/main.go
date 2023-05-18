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
	// Myriadファイルから全ての行を読む
	lines, err := others.ReadLinesFromFile(os.Args[1])
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
	functionInterCodeMap, functionVarMap, err := compiler.Compile(t.Tokens)
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

	// コード生成
	codes, err := generator.GenerateCode(functionInterCodeMap, functionVarMap)
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
