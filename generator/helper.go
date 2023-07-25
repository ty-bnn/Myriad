package generator

import (
	"myriad/compiler"
)

func (g *Generator) generateCodeBlock(index int) (int, []string, error) {
	var codes []string
	var err error

	for index < len(g.MainCodes) {
		var code string
		var codeBlock []string
		// fmt.Println(index)

		switch g.MainCodes[index].GetKind() {
		case compiler.ROW:
			code = g.MainCodes[index].GetContent()
			codes = append(codes, code)
			index++
		case compiler.COMMAND:
			if g.command == "RUN" && g.MainCodes[index].GetContent() == "RUN" {
				// RUN命令の結合
				codes[len(codes)-1] = codes[len(codes)-1][:len(codes[len(codes)-1])-1] + " \\\n"
				code = "    "
				index++
			} else {
				code = g.MainCodes[index].GetContent()
				g.command = g.MainCodes[index].GetContent()
			}
			codes = append(codes, code)
			index++
		case compiler.IF:
			index, codeBlock, err = g.generateIfBlock(index)
			if err != nil {
				return index, codes, err
			}

			codes = append(codes, codeBlock...)
		case compiler.ENDIF:
			index++
			return index, codes, nil
		}
	}

	return index, codes, nil
}

func (g *Generator) generateIfBlock(index int) (int, []string, error) {
	var codes []string
	var codeBlock []string
	var err error

	// if節
	index, codeBlock, err = g.generateIf(index)
	if err != nil {
		return index, codes, err
	}
	// if節の条件がTrueならばそれ以降のifブロックは飛ばす
	if codeBlock != nil {
		codes = append(codes, codeBlock...)
		return index, codes, nil
	}

	// elif節
	for {
		if g.MainCodes[index].GetKind() != compiler.ELIF {
			break
		}

		index, codeBlock, err = g.generateIf(index)
		if err != nil {
			return index, codes, err
		}
		// elif節の条件がTrueならばそれ以降のifブロックは飛ばす
		if codeBlock != nil {
			codes = append(codes, codeBlock...)
			return index, codes, nil
		}
	}

	// else節
	if g.MainCodes[index].GetKind() != compiler.ELSE {
		return index, codes, nil
	}

	// else節はgenerateCodeBlockで対応可能（ENDIFを読み取ったら帰ってくる）
	index, codeBlock, err = g.generateCodeBlock(index + 1)
	if err != nil {
		return index, codes, err
	}
	codes = append(codes, codeBlock...)

	return index, codes, nil
}

// if節とelif節の処理は同じ
func (g *Generator) generateIf(index int) (int, []string, error) {
	var codes []string

	code, _ := g.MainCodes[index].(compiler.IfInterCode)
	nextOffset := code.IfContent.NextOffset
	endOffset := code.IfContent.EndOffset

	condition, err := getIfCondition(code.IfContent)
	if err != nil {
		return index, codes, err
	}

	// 条件式がfalseの場合、if節内の処理は返さずに次のindexを返す
	if !condition {
		return index + nextOffset + 1, codes, nil
	}

	_, codes, err = g.generateCodeBlock(index + 1)
	if err != nil {
		return index, codes, err
	}

	// 条件式がtrueの場合、if
	return index + endOffset + 1, codes, nil
}
