package generator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ty-bnn/myriad/pkg/model/vars"

	"github.com/ty-bnn/myriad/pkg/model/codes"
	"github.com/ty-bnn/myriad/pkg/model/values"
)

func (g *Generator) Generate() error {
	var err error

	fmt.Println("Generating...")

	//mainArgs := os.Args[3:]
	g.funcPtr = "main"

	g.RawCodes, err = g.callFunc(nil)
	if err != nil {
		return err
	}

	fmt.Println("Generate Done.")

	return nil
}

func (g *Generator) callFunc(args []values.Value) ([]string, error) {
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

		defCode := funcCodes[g.index].(codes.Define)
		vTable = append(vTable, vars.Var{Name: defCode.Key, Value: arg})

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
			// 末尾が'\', '\n'で終わっているか確認
			trimmed := strings.Replace(literal.Content, " ", "", -1)
			if 0 < len(trimmed) && trimmed[len(trimmed)-2:] == "\\\n" {
				rawCodes = append(rawCodes, whiteSpaces(g.commandPtr))
			}
			g.index++
		case codes.COMMAND:
			command := code.(codes.Command)
			rawCodes = append(rawCodes, command.Content)
			g.commandPtr = command.Content
			rawCodes = append(rawCodes, " ")
			g.index++
		case codes.DEFINE:
			define := code.(codes.Define)
			value, err := getValue(vTable, define.Value)
			if err != nil {
				return nil, err
			}
			vTable = append(vTable, vars.Var{Name: define.Key, Value: value})
			g.index++
		case codes.ASSIGN:
			assign := code.(codes.Assign)
			value, err := getValue(vTable, assign.Value)
			if err != nil {
				return nil, err
			}
			index, err := getIndex(vTable, assign.Key)
			if err != nil {
				return nil, err
			}
			vTable[index] = vars.Var{Name: assign.Key, Value: value}
			g.index++
		case codes.REPLACE:
			rep := code.(codes.Replace)
			value, err := getLiteral(vTable, rep.Value)
			if err != nil {
				return nil, err
			}
			rawCodes = append(rawCodes, value)
			g.index++
		case codes.CALLPROC:
			callProc := code.(codes.CallProc)

			var args []values.Value
			for _, arg := range callProc.Args {
				v, err := getValue(vTable, arg)
				if err != nil {
					return nil, err
				}
				args = append(args, v)
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
	for g.index < len(funcCodes) && funcCodes[g.index].GetKind() == codes.ELIF {
		var elifSecCodes []string
		elifSecCodes, err = g.elifSection(vTable)
		if err != nil {
			return nil, err
		}
		ifBCodes = append(ifBCodes, elifSecCodes...)
	}

	// ELSEコード
	if g.index < len(funcCodes) && funcCodes[g.index].GetKind() == codes.ELSE {
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
	ok, err := evalCondition(vTable, ifCode.Condition)
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
	ok, err := evalCondition(vTable, elifCode.Condition)
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

	literals, err := getLiterals(vTable, forCode.ArrayValue)
	if err != nil {
		return nil, err
	}

	start := g.index
	itrName := forCode.ItrName

	for _, literal := range literals {
		var rowCodes []string
		g.index = start

		vTable = append(vTable, vars.Var{Name: itrName, Value: values.Literal{Kind: values.LITERAL, Value: literal}})

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
