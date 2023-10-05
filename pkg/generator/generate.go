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

	var vTable []vars.Var

	// 引数を変数として定義
	for _, arg := range args {
		if len(funcCodes) <= g.index || funcCodes[g.index].GetKind() != codes.DEFINE {
			return nil, errors.New(fmt.Sprintf("semantic error: %s got too many args", g.funcPtr))
		}

		v := funcCodes[g.index].(codes.Define).Var.(vars.Single)
		v.Value = arg
		vTable = append(vTable, v)

		g.index++
	}

	// コードブロック
	rowCodes, err := g.codeBlock(vTable)
	if err != nil {
		return nil, err
	}

	return rowCodes, nil
}

func (g *Generator) codeBlock(vTable []vars.Var) ([]string, error) {
	funcCodes := g.funcToCodes[g.funcPtr]
	var rawCodes []string

	for g.index < len(funcCodes) {
		code := funcCodes[g.index]

		switch code.GetKind() {
		case codes.LITERAL:
			literal := code.(codes.Literal)
			rawCodes = append(rawCodes, literal.Content)
			g.index++
		case codes.COMMAND:
			command := code.(codes.Command)
			rawCodes = append(rawCodes, command.Content)
			g.index++
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
			g.index++
		case codes.ASSIGN:
			assign := code.(codes.Assign)
			if err := assignVar(vTable, assign.Var); err != nil {
				return nil, err
			}
			g.index++
		case codes.REPLACE:
			rep := code.(codes.Replace)
			value, err := getValue(vTable, rep.RepVar)
			if err != nil {
				return nil, err
			}
			rawCodes = append(rawCodes, value)
			g.index++
		case codes.CALLPROC:
			callProc := code.(codes.CallProc)
			var args []string
			for _, arg := range callProc.Args {
				value, err := getValue(vTable, arg)
				if err != nil {
					return nil, err
				}

				args = append(args, value)
			}

			funcStack := g.funcPtr
			g.funcPtr = callProc.ProcName
			indexStack := g.index
			g.index = 0

			funcRawCodes, err := g.callFunc(args)
			if err != nil {
				return nil, err
			}

			g.funcPtr = funcStack
			g.index = indexStack

			rawCodes = append(rawCodes, funcRawCodes...)
			g.index++
		case codes.IF:
			var ifCodes []string
			var err error
			ifCodes, err = g.ifBlock(vTable)
			if err != nil {
				return nil, err
			}
			rawCodes = append(rawCodes, ifCodes...)
		case codes.FOR:
			var forCodes []string
			var err error
			forCodes, err = g.forBlock(vTable)
			if err != nil {
				return nil, err
			}
			rawCodes = append(rawCodes, forCodes...)
		default:
			return rawCodes, nil
		}
	}

	return rawCodes, nil
}

func (g *Generator) ifBlock(vTable []vars.Var) ([]string, error) {
	funcCodes := g.funcToCodes[g.funcPtr]
	var ifBCodes []string

	// IFコード
	ifSecCodes, err := g.ifSection(vTable)
	if err != nil {
		return nil, err
	}
	ifBCodes = append(ifBCodes, ifSecCodes...)

	// ELIFコード
	for funcCodes[g.index].GetKind() == codes.ELIF {
		var elifSecCodes []string
		elifSecCodes, err = g.elifSection(vTable)
		if err != nil {
			return nil, err
		}
		ifBCodes = append(ifBCodes, elifSecCodes...)
	}

	// ELSEコード
	if funcCodes[g.index].GetKind() == codes.ELSE {
		var elseSecCodes []string
		elseSecCodes, err = g.elseSection(vTable)
		if err != nil {
			return nil, err
		}
		ifBCodes = append(ifBCodes, elseSecCodes...)
	}

	return ifBCodes, nil
}

func (g *Generator) ifSection(vTable []vars.Var) ([]string, error) {
	funcCodes := g.funcToCodes[g.funcPtr]

	// IFコード
	ifCode := funcCodes[g.index].(codes.If)
	ok, err := getConditionEval(vTable, ifCode.Condition)
	if err != nil {
		return nil, err
	}

	g.index++

	// コードブロック
	rowCodes, err := g.codeBlock(vTable)
	if err != nil {
		return nil, err
	}

	// ENDコード
	g.index++

	if !ok {
		return nil, nil
	}

	return rowCodes, nil
}

func (g *Generator) elifSection(vTable []vars.Var) ([]string, error) {
	funcCodes := g.funcToCodes[g.funcPtr]

	// ELIFコード
	elifCode := funcCodes[g.index].(codes.Elif)
	ok, err := getConditionEval(vTable, elifCode.Condition)
	if err != nil {
		return nil, err
	}

	g.index++

	// コードブロック
	rowCodes, err := g.codeBlock(vTable)
	if err != nil {
		return nil, err
	}

	// ENDコード
	g.index++

	if !ok {
		return nil, nil
	}

	return rowCodes, nil
}

func (g *Generator) elseSection(vTable []vars.Var) ([]string, error) {
	// ELSEコード
	g.index++

	// コードブロック
	rowCodes, err := g.codeBlock(vTable)
	if err != nil {
		return nil, err
	}

	// ENDブロック
	g.index++

	return rowCodes, nil
}

func (g *Generator) forBlock(vTable []vars.Var) ([]string, error) {
	funcCodes := g.funcToCodes[g.funcPtr]

	var forCodes []string

	// FORコード
	forCode := funcCodes[g.index].(codes.For)
	g.index++

	array, err := getValues(vTable, forCode.Array)
	if err != nil {
		return nil, err
	}

	start := g.index
	for _, value := range array {
		var rowCodes []string
		g.index = start

		forCode.Itr.Value = value
		vTable = append(vTable, forCode.Itr)

		// コードブロック
		rowCodes, err = g.codeBlock(vTable)
		if err != nil {
			return nil, err
		}
		forCodes = append(forCodes, rowCodes...)

		// ENDコード
		g.index++

		// for文で定義したイテレータをPOP
		vTable = vTable[:len(vTable)-1]
	}

	return forCodes, nil
}
