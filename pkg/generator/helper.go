package generator

import (
	"errors"
	"fmt"

	"github.com/ty-bnn/myriad/pkg/model/codes"
	"github.com/ty-bnn/myriad/pkg/model/vars"
)

func (g *Generator) callFunc(funcName string, args []string) error {
	var vTable []vars.Var

	funcCodes, ok := g.FuncToCodes[funcName]
	if !ok {
		return errors.New(fmt.Sprintf("semantic error: %s is not defined", funcName))
	}

	var index int

	for _, arg := range args {
		code := funcCodes[index]
		v := code.(codes.Define).Var.(vars.Single)
		v.Value = arg
		vTable = append(vTable, v)

		index++
	}

	for index < len(funcCodes) {
		code := funcCodes[index]

		switch code.GetKind() {
		case codes.LITERAL:
			literal := code.(codes.Literal)
			g.Codes = append(g.Codes, literal.Content)
			index++
		case codes.COMMAND:
			command := code.(codes.Command)
			if g.command == "RUN" && command.Content == "RUN" {
				// RUN命令の結合
				str := g.Codes[len(g.Codes)-1]
				str = str[0:len(str)-1] + " \\\n"
				g.Codes[len(g.Codes)-1] = str

				g.Codes = append(g.Codes, "   ")
			} else {
				g.command = command.Content
				g.Codes = append(g.Codes, g.command)
			}
			index++
		case codes.DEFINE:
			define := code.(codes.Define)
			vTable = append(vTable, define.Var)
			index++
		case codes.ASSIGN:
			assign := code.(codes.Assign)
			if err := assignVar(vTable, assign.Var); err != nil {
				return err
			}

			index++
		case codes.REPLACE:
			rep := code.(codes.Replace)
			value, err := getValue(vTable, rep.RepVar)
			if err != nil {
				return err
			}

			g.Codes = append(g.Codes, value)
			index++
		case codes.CALLPROC:
			callProc := code.(codes.CallProc)
			var args []string
			for _, arg := range callProc.Args {
				value, err := getValue(vTable, arg)
				if err != nil {
					return err
				}

				args = append(args, value)
			}

			err := g.callFunc(callProc.ProcName, args)
			if err != nil {
				return err
			}

			index++
		case codes.IF:
			ifCode := code.(codes.If)
			ok, err := getConditionEval(vTable, ifCode.Condition)
			if err != nil {
				return err
			}

			if ok {
				index += 2
			} else {
				index++
			}
		case codes.ELIF:
			elifCode := code.(codes.Elif)
			ok, err := getConditionEval(vTable, elifCode.Condition)
			if err != nil {
				return err
			}

			if ok {
				index += 2
			} else {
				index++
			}
		case codes.ELSE:
			index++
		case codes.JUMP:
			jump := code.(codes.Jump)
			index += jump.NextOffset
		case codes.POP:
			vTable = vTable[:len(vTable)-1]
			index++
		}
	}

	return nil
}
