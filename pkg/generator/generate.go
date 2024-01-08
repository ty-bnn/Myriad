package generator

import (
	"errors"
	"fmt"
	"sort"

	"github.com/ty-bnn/myriad/pkg/utils"

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

	shapeRawCodes(g.RawCodes)

	fmt.Println("Generate Done.")

	utils.WriteStdOut(g.RawCodes)

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
		case codes.APPEND:
			appendCode := code.(codes.Append)
			elem, err := getLiteral(vTable, appendCode.Element)
			if err != nil {
				return nil, err
			}
			index, err := getIndex(vTable, appendCode.Array)
			if err != nil {
				return nil, err
			}
			if vTable[index].Value.GetKind() != values.LITERALS {
				return nil, errors.New(fmt.Sprintf("semantic error: %s is not array", appendCode.Array))
			}
			elements := vTable[index].Value.(values.Literals).Values
			elements = append(elements, elem)
			vTable[index].Value = values.Literals{Kind: values.LITERALS, Values: elements}
			g.index++
		case codes.SORT:
			appendCode := code.(codes.Sort)
			index, err := getIndex(vTable, appendCode.Array)
			if err != nil {
				return nil, err
			}
			if vTable[index].Value.GetKind() != values.LITERALS {
				return nil, errors.New(fmt.Sprintf("semantic error: %s is not array", appendCode.Array))
			}
			sort.Strings(vTable[index].Value.(values.Literals).Values)
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
			ifCodes, err := g.ifBlock(vTable)
			if err != nil {
				return nil, err
			}
			rawCodes = append(rawCodes, ifCodes...)
		case codes.FOR:
			forCodes, err := g.forBlock(vTable)
			if err != nil {
				return nil, err
			}
			rawCodes = append(rawCodes, forCodes...)
		case codes.OUTPUT:
			outCode := code.(codes.Output)
			outPath, err := getLiteral(vTable, outCode.FilePath)
			if err != nil {
				return nil, err
			}
			g.index++

			outCodes, err := g.codeBlock(vTable)
			if err != nil {
				return nil, err
			}

			// ENDコード
			g.index++

			shapeRawCodes(outCodes)

			err = utils.WriteFile(outCodes, outPath)
			if err != nil {
				return nil, err
			}
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
	ifCodePtr := g.index
	ok, err := evalCondition(vTable, ifCode.Condition)
	if err != nil {
		return nil, err
	}
	if !ok {
		g.index += ifCode.False
		return nil, nil
	}
	g.index++

	// コードブロック
	rowCodes, err := g.codeBlock(vTable)
	if err != nil {
		return nil, err
	}

	g.index = ifCodePtr + ifCode.True

	return rowCodes, nil
}

func (g *Generator) elifSection(vTable []vars.Var) ([]string, error) {
	funcCodes := g.funcToCodes[g.funcPtr]

	// ELIFコード
	elifCode := funcCodes[g.index].(codes.Elif)
	elifCodePtr := g.index
	ok, err := evalCondition(vTable, elifCode.Condition)
	if err != nil {
		return nil, err
	}
	if !ok {
		g.index += elifCode.False
		return nil, nil
	}
	g.index++

	// コードブロック
	rowCodes, err := g.codeBlock(vTable)
	if err != nil {
		return nil, err
	}

	g.index = elifCodePtr + elifCode.True

	return rowCodes, nil
}

// TODO: elseの到達判定を再考
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
