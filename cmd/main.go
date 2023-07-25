package main

import (
	"fmt"
	"os"

	"myriad/compiler"
	"myriad/generator"
	"myriad/helpers"
	"myriad/parser"
	"myriad/tokenizer"
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

	// コード生成
	g := &generator.Generator{}
	err = g.GenerateCode(c.FunctionInterCodeMap)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// ファイル出力
	err = helpers.WriteFile(g.Codes, os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
