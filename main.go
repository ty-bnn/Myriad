package main

import(
	"fmt"
	"os"

	"dcc/others"
	"dcc/tokenizer"
	// "dcc/parser"
	// "dcc/compiler"
	// "dcc/generator"
)

func main() {
	// Myriadファイルから全ての行を読む
	lines, err := others.ReadLinesFromFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// トークナイズ
	_, err = tokenizer.Tokenize(lines)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// // パース
	// err = parser.Parse(tokens)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	
	// // コンパイル
	// functionInterCodeMap, functionArgMap, err := compiler.Compile(tokens)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// // For debug
	// // for k, v := range functionInterCodeMap {
	// // 	fmt.Println(k)
	// // 	for _, c := range v {
	// // 		fmt.Println(c)
	// // 	}
	// // }

	// // コード生成
	// codes, err := generator.GenerateCode(functionInterCodeMap, functionArgMap)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// err = others.WriteFile(codes, os.Args[2])
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}
