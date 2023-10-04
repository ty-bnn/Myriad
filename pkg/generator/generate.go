package generator

import (
	"errors"
	"fmt"
	"os"

	"github.com/ty-bnn/myriad/pkg/model/codes"
	"github.com/ty-bnn/myriad/pkg/model/vars"
)

func (g *Generator) Generate() error {
	var err error

	fmt.Println("Generating...")

	mainArgs := os.Args[3:]
	g.funcPtr = "main"

	g.RawCodes, err = g.callFunc(mainArgs)
	if err != nil {
		return err
	}

	fmt.Println("Generate Done.")

	return nil
}

func (g *Generator) callFunc(args []string) ([]string, error) {
	funcCodes, ok := g.funcToCodes[g.funcPtr]
	if !ok {
		return nil, errors.New(fmt.Sprintf("semantic error: %s is not defined", g.funcPtr))
	}

	var index int
	var vTable []vars.Var

	// 引数を変数として定義
	for _, arg := range args {
		if len(funcCodes) <= index || funcCodes[index].GetKind() != codes.DEFINE {
			return nil, errors.New(fmt.Sprintf("semantic error: %s got too many args", g.funcPtr))
		}

		v := funcCodes[index].(codes.Define).Var.(vars.Single)
		v.Value = arg
		vTable = append(vTable, v)

		index++
	}

	// コードブロック
	rowCodes, _, err := g.codeBlock(index, vTable)
	if err != nil {
		return nil, err
	}

	return rowCodes, nil
}

func (g *Generator) codeBlock(index int, vTable []vars.Var) ([]string, int, error) {
	funcCodes := g.funcToCodes[g.funcPtr]
	var rawCodes []string

	for index < len(funcCodes) {
		code := funcCodes[index]

		switch code.GetKind() {
		case codes.LITERAL:
			literal := code.(codes.Literal)
			rawCodes = append(rawCodes, literal.Content)
			index++
		case codes.COMMAND:
			command := code.(codes.Command)
			rawCodes = append(rawCodes, command.Content)
			index++
			// TODO: RUNの連結はGenerateするときに行う
			//if g.command == "RUN" && command.Content == "RUN" {
			//	// RUN命令の結合
			//	str := g.Codes[len(g.Codes)-1]
			//	str = str[0:len(str)-1] + " \\\n"
			//	g.Codes[len(g.Codes)-1] = str
			//
			//	g.Codes = append(g.Codes, "   ")
			//} else {
			//	g.command = command.Content
			//	g.Codes = append(g.Codes, command.Content)
			//}
		case codes.DEFINE:
			define := code.(codes.Define)
			vTable = append(vTable, define.Var)
			index++
		case codes.ASSIGN:
			assign := code.(codes.Assign)
			if err := assignVar(vTable, assign.Var); err != nil {
				return nil, -1, err
			}
			index++
		case codes.REPLACE:
			rep := code.(codes.Replace)
			value, err := getValue(vTable, rep.RepVar)
			if err != nil {
				return nil, -1, err
			}
			rawCodes = append(rawCodes, value)
			index++
		case codes.CALLPROC:
			callProc := code.(codes.CallProc)
			var args []string
			for _, arg := range callProc.Args {
				value, err := getValue(vTable, arg)
				if err != nil {
					return nil, -1, err
				}

				args = append(args, value)
			}

			calledFrom := g.funcPtr
			g.funcPtr = callProc.ProcName

			funcRawCodes, err := g.callFunc(args)
			if err != nil {
				return nil, -1, err
			}

			g.funcPtr = calledFrom
			rawCodes = append(rawCodes, funcRawCodes...)
			index++
		case codes.IF:
			var ifCodes []string
			var err error
			ifCodes, index, err = g.ifBlock(index, vTable)
			if err != nil {
				return nil, -1, err
			}
			rawCodes = append(rawCodes, ifCodes...)
		case codes.FOR:
			var forCodes []string
			var err error
			forCodes, index, err = g.forBlock(index, vTable)
			if err != nil {
				return nil, -1, err
			}
			rawCodes = append(rawCodes, forCodes...)
		default:
			return rawCodes, index, nil
		}
	}

	return rawCodes, index, nil
}

func (g *Generator) ifBlock(index int, vTable []vars.Var) ([]string, int, error) {
	funcCodes := g.funcToCodes[g.funcPtr]
	var ifBCodes []string

	// IFコード
	ifSecCodes, index, err := g.ifSection(index, vTable)
	if err != nil {
		return nil, -1, err
	}
	ifBCodes = append(ifBCodes, ifSecCodes...)

	// ELIFコード
	for funcCodes[index].GetKind() == codes.ELIF {
		var elifSecCodes []string
		elifSecCodes, index, err = g.elifSection(index, vTable)
		if err != nil {
			return nil, -1, err
		}
		ifBCodes = append(ifBCodes, elifSecCodes...)
	}

	// ELSEコード
	if funcCodes[index].GetKind() == codes.ELSE {
		var elseSecCodes []string
		elseSecCodes, index, err = g.elseSection(index, vTable)
		if err != nil {
			return nil, -1, err
		}
		ifBCodes = append(ifBCodes, elseSecCodes...)
	}

	return ifBCodes, index, nil
}

func (g *Generator) ifSection(index int, vTable []vars.Var) ([]string, int, error) {
	funcCodes := g.funcToCodes[g.funcPtr]

	// IFコード
	ifCode := funcCodes[index].(codes.If)
	ok, err := getConditionEval(vTable, ifCode.Condition)
	if err != nil {
		return nil, -1, err
	}

	index++

	// コードブロック
	rowCodes, index, err := g.codeBlock(index, vTable)
	if err != nil {
		return nil, -1, err
	}

	// ENDコード
	index++

	if !ok {
		return nil, index, nil
	}

	return rowCodes, index, nil
}

func (g *Generator) elifSection(index int, vTable []vars.Var) ([]string, int, error) {
	funcCodes := g.funcToCodes[g.funcPtr]

	// ELIFコード
	elifCode := funcCodes[index].(codes.Elif)
	ok, err := getConditionEval(vTable, elifCode.Condition)
	if err != nil {
		return nil, -1, err
	}

	index++

	// コードブロック
	rowCodes, index, err := g.codeBlock(index, vTable)
	if err != nil {
		return nil, -1, err
	}

	// ENDコード
	index++

	if !ok {
		return nil, index, nil
	}

	return rowCodes, index, nil
}

func (g *Generator) elseSection(index int, vTable []vars.Var) ([]string, int, error) {
	// ELSEコード
	index++

	// コードブロック
	rowCodes, index, err := g.codeBlock(index, vTable)
	if err != nil {
		return nil, -1, err
	}

	// ENDブロック
	index++

	return rowCodes, index, nil
}

func (g *Generator) forBlock(index int, vTable []vars.Var) ([]string, int, error) {
	funcCodes := g.funcToCodes[g.funcPtr]

	var forCodes []string

	// FORコード
	forCode := funcCodes[index].(codes.For)
	index++

	array, err := getValues(vTable, forCode.Array)
	if err != nil {
		return nil, -1, err
	}

	start := index
	for _, value := range array {
		var rowCodes []string
		index = start

		forCode.Itr.Value = value
		vTable = append(vTable, forCode.Itr)

		// コードブロック
		rowCodes, index, err = g.codeBlock(index, vTable)
		if err != nil {
			return nil, -1, err
		}
		forCodes = append(forCodes, rowCodes...)

		// ENDコード
		index++

		// for文で定義したイテレータをPOP
		vTable = vTable[:len(vTable)-1]
	}

	return forCodes, index, nil
}
