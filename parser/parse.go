package parser

import (
	"dcc/tokenizer"
)

// type NodeKind int

// const (
// 	ND_PROGRAM NodeKind = iota // プログラム
// 	ND_IMPORT // インポート部
// 	ND_FUNCDEC // 関数宣言部
// 	ND_MAIN // メイン部
// 	ND_FUNC // 関数
// 	ND_ARGDEC // 引数宣言部
// 	ND_FUNCDESC // 関数記述部
// 	ND_DFILE // Dfile文
// 	ND_DFCMD // Df命令
// 	ND_DFARG // Df引数
// 	ND_FUNCCALL // 関数呼び出し文
// 	ND_FUNCNAME // 関数名
// 	ND_FILENAME // ファイル名
// 	ND_STRING // 文字列
// )

// type Node struct {
// 	Kind NodeKind
// 	Cnodes []Node
// 	Val string
// }

func Parse(tokens []tokenizer.Token) (error){
	err := program(tokens, 0)

	return err
}