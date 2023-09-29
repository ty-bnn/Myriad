package generator

import (
	"errors"
	"fmt"

	"github.com/ty-bnn/myriad/pkg/model/codes"
	"github.com/ty-bnn/myriad/pkg/model/vars"
)

func getConditionEval(vTable []vars.Var, condition codes.Condition) (bool, error) {
	left, err := getValue(vTable, condition.Left)
	if err != nil {
		return false, err
	}

	right, err := getValue(vTable, condition.Right)
	if err != nil {
		return false, err
	}

	if condition.Operator == codes.EQUAL && left == right {
		return true, nil
	} else if condition.Operator == codes.NOTEQUAL && left != right {
		return true, nil
	}

	return false, nil
}

// getValue returns the value of vars.
// 変数のスコープを実現するために、変数表の後ろから変数名を探索する
func getValue(vTable []vars.Var, target vars.Var) (string, error) {
	// 文字列が入っていた場合はそのまま値を返す
	if target.GetKind() == vars.LITERAL {
		return target.(vars.Literal).Value, nil
	}

	for i := len(vTable) - 1; i >= 0; i-- {
		if vTable[i].GetName() != target.GetName() {
			continue
		}

		switch target.GetKind() {
		case vars.SINGLE:
			if vTable[i].GetKind() != vars.SINGLE {
				return "", errors.New(fmt.Sprintf("semantic error: cannot use %s as type single", target.GetName()))
			}

			return vTable[i].(vars.Single).Value, nil
		case vars.ELEMENT:
			if vTable[i].GetKind() != vars.ARRAY {
				return "", errors.New(fmt.Sprintf("semantic error: cannot use %s as type array element", target.GetName()))
			}

			index := target.(vars.Element).Index
			if index < 0 || len(vTable[i].(vars.Array).Values) <= index {
				return "", errors.New(fmt.Sprintf("semantic error: out of index for %s", target.GetName()))
			}

			return vTable[i].(vars.Array).Values[index], nil
		}
	}

	return "", errors.New(fmt.Sprintf("semantic error: %s is not declared", target.GetName()))
}
