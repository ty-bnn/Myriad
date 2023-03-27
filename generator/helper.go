package generator

import (
	"fmt"

	"dcc/compiler"
)

func generateCodeBlock(index int) (int, []string, error) {
	fmt.Println("---------- generating code ---------")

	var codes []string
	var err error
	
	for index < len(mainCodes) {
		var code string
		var codeBlock []string
		// fmt.Println(index)

		switch mainCodes[index].Kind {
			case compiler.ROW:
				code = mainCodes[index].Content
				codes = append(codes, code)
				index++
			case compiler.COMMAND:
				if command == "RUN" && mainCodes[index].Content == "RUN" {
					// RUN命令の結合
					codes[len(codes) - 1] = codes[len(codes) - 1][:len(codes[len(codes) - 1]) - 1] + " \\\n"
					code = "    "
					index++
				} else {
					code = mainCodes[index].Content
					command = mainCodes[index].Content
				}
				codes = append(codes, code)
				index++
			case compiler.IF:
				index, codeBlock, err = generateIfBlock(index)
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

func generateIfBlock(index int) (int, []string, error) {
	var codes []string
	var codeBlock []string
	var err error

	// if節
	index, codeBlock, err = generateIf(index)
	if err != nil {
		return index, codes, err
	}
	// if節の条件がTrueならばそれ以降のifブロックは飛ばす
	if codeBlock != nil {
		codes = append(codes, codeBlock...)
		return index, codes, nil
	}

	// elif節
	for ;; {
		if mainCodes[index].Kind != compiler.ELIF {
			break
		}

		index, codeBlock, err = generateIf(index)
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
	if mainCodes[index].Kind != compiler.ELSE {
		return index, codes, nil
	}
	
	// else節はgenerateCodeBlockで対応可能（ENDIFを読み取ったら帰ってくる）
	index, codeBlock, err = generateCodeBlock(index + 1)
	if err != nil {
		return index, codes, err
	}
	codes = append(codes, codeBlock...)

	return index, codes, nil
}

// if節とelif節の処理は同じ
func generateIf(index int) (int, []string, error) {
	var codes []string

	nextOffset := mainCodes[index].IfContent.NextOffset
	endOffset := mainCodes[index].IfContent.EndOffset

	condition, err := getIfCondition(mainCodes[index].IfContent)
	if err != nil {
		return index, codes, err
	}

	// 条件式がflaseの場合、if節内の処理は返さずに次のindexを返す
	if !condition {
		return index + nextOffset + 1, codes, nil
	}

	_, codes, err = generateCodeBlock(index + 1)
	if err != nil {
		return index, codes, err
	}

	// 条件式がtrueの場合、if
	return index + endOffset + 1, codes, nil
}
