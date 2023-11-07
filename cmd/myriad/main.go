package main

import (
	"fmt"
	"os"

	"github.com/ty-bnn/myriad/pkg/generator"
	"github.com/ty-bnn/myriad/pkg/parser"
	"github.com/ty-bnn/myriad/pkg/tokenizer"
	"github.com/ty-bnn/myriad/pkg/utils"
)

func main() {
	// Myriadファイルから全ての行を読む
	data, err := utils.ReadLinesFromFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// トークナイズ
	t := tokenizer.New(data)
	err = t.Tokenize()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// パース
	p := parser.New(t.Tokens)
	err = p.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Dockerfile生成
	g := generator.New(p.FuncToCodes)
	err = g.Generate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
