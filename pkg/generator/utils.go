package generator

import (
	"errors"
	"fmt"

	"github.com/ty-bnn/myriad/pkg/model/vars"

	"github.com/ty-bnn/myriad/pkg/model/codes"
	"github.com/ty-bnn/myriad/pkg/model/values"
)

func evalCondition(vTable []vars.Var, root codes.ConditionalNode) (bool, error) {
	if root.Operator == codes.EQUAL || root.Operator == codes.NOTEQUAL {
		eq, err := isEqual(vTable, root)
		if err != nil {
			return false, err
		}
		return eq, nil
	}

	lEq, err := evalCondition(vTable, *root.Left)
	if err != nil {
		return false, err
	}

	rEq, err := evalCondition(vTable, *root.Right)
	if err != nil {
		return false, err
	}

	if root.Operator == codes.OR {
		return lEq || rEq, nil
	}

	return lEq && rEq, nil
}

func isEqual(vTable []vars.Var, node codes.ConditionalNode) (bool, error) {
	left, err := getLiteral(vTable, node.Left.Var)
	if err != nil {
		return false, err
	}

	right, err := getLiteral(vTable, node.Right.Var)
	if err != nil {
		return false, err
	}

	if (node.Operator == codes.EQUAL && left == right) || (node.Operator == codes.NOTEQUAL && left != right) {
		return true, nil
	}

	return false, nil
}

// getLiteral returns a literal.
// 変数のスコープを実現するために、変数表の後ろから変数名を探索する
// literal, ident, element, map_valueに対応
func getLiteral(vTable []vars.Var, target values.Value) (string, error) {
	// 文字列が入っていた場合はそのまま値を返す
	if target.GetKind() == values.LITERAL {
		return target.(values.Literal).Value, nil
	}

	for i := len(vTable) - 1; i >= 0; i-- {
		if vTable[i].Name != target.GetName() {
			continue
		}

		switch target.GetKind() {
		case values.IDENT:
			if vTable[i].Value.GetKind() != values.LITERAL {
				return "", errors.New(fmt.Sprintf("semantic error: cannot use %s as type single", target.GetName()))
			}
			return vTable[i].Value.(values.Literal).Value, nil
		case values.ELEMENT:
			if vTable[i].Value.GetKind() != values.LITERALS {
				return "", errors.New(fmt.Sprintf("semantic error: cannot use %s as type array element", target.GetName()))
			}

			index := target.(values.Element).Index
			if index < 0 || len(vTable[i].Value.(values.Literals).Values) <= index {
				return "", errors.New(fmt.Sprintf("semantic error: out of index for %s", target.GetName()))
			}
			return vTable[i].Value.(values.Literals).Values[index], nil
		case values.MAPVALUE:
			if vTable[i].Value.GetKind() != values.MAP {
				return "", errors.New(fmt.Sprintf("semantic error: cannot use %s as type map", target.GetName()))
			}

			key := target.(values.MapValue).Key
			keyValue, err := getLiteral(vTable, key)
			if err != nil {
				return "", err
			}
			value, ok := vTable[i].Value.(values.Map).Value[keyValue]
			if !ok {
				return "", errors.New(fmt.Sprintf("semantic error: missing %s in %s as a key", keyValue, target.GetName()))
			}
			strValue, ok := value.(string)
			if !ok {
				return "", errors.New(fmt.Sprintf("semantic error: value is not type literal"))
			}
			return strValue, nil
		default:
			return "", errors.New(fmt.Sprintf("semantic error: value is not type literal"))
		}
	}

	return "", errors.New(fmt.Sprintf("semantic error: %s is not declared", target.GetName()))
}

// getLiterals returns literals.
// literals, ident, map_key, map_valueに対応
func getLiterals(vTable []vars.Var, target values.Value) ([]string, error) {
	if target.GetKind() == values.LITERALS {
		return target.(values.Literals).Values, nil
	}

	for i := len(vTable) - 1; i >= 0; i-- {
		if vTable[i].Name != target.GetName() {
			continue
		}

		switch target.GetKind() {
		case values.IDENT:
			if vTable[i].Value.GetKind() != values.LITERALS {
				return nil, errors.New(fmt.Sprintf("semantic error: cannot use %s as type array", target.GetName()))
			}

			return vTable[i].Value.(values.Literals).Values, nil
		case values.MAPKEY:
			if vTable[i].Value.GetKind() != values.MAP {
				return nil, errors.New(fmt.Sprintf("semantic error: cannot use %s as type map", target.GetName()))
			}

			var keys []string
			for key := range vTable[i].Value.(values.Map).Value {
				keys = append(keys, key)
			}

			return keys, nil
		case values.MAPVALUE:
			if vTable[i].Value.GetKind() != values.MAP {
				return nil, errors.New(fmt.Sprintf("semantic error: cannot use %s as type map", target.GetName()))
			}

			key := target.(values.MapValue).Key
			keyValue, err := getLiteral(vTable, key)
			if err != nil {
				return nil, err
			}
			value, ok := vTable[i].Value.(values.Map).Value[keyValue]
			if !ok {
				return nil, errors.New(fmt.Sprintf("semantic error: missing %s in %s as a key", keyValue, target.GetName()))
			}
			arrayValue, ok := assertionToStringSlice(value)
			if !ok {
				return nil, errors.New(fmt.Sprintf("semantic error: value is not type literals"))
			}
			return arrayValue, nil
		default:
			return nil, errors.New(fmt.Sprintf("semantic error: value is not type literals"))
		}
	}

	return nil, errors.New(fmt.Sprintf("semantic error: %s is not declared", target.GetName()))
}

// getMap returns map_value.
// map, map_valueに対応
func getMap(vTable []vars.Var, target values.Value) (map[string]interface{}, error) {
	if target.GetKind() == values.MAP {
		return target.(values.Map).Value, nil
	}

	for i := len(vTable) - 1; i >= 0; i-- {
		if vTable[i].Name != target.GetName() {
			continue
		}

		switch target.GetKind() {
		case values.MAPVALUE:
			if vTable[i].Value.GetKind() != values.MAP {
				return nil, errors.New(fmt.Sprintf("semantic error: cannot use %s as type map", target.GetName()))
			}

			key := target.(values.MapValue).Key
			keyValue, err := getLiteral(vTable, key)
			if err != nil {
				return nil, err
			}
			value, ok := vTable[i].Value.(values.Map).Value[keyValue]
			if !ok {
				return nil, errors.New(fmt.Sprintf("semantic error: missing %s in %s as a key", keyValue, target.GetName()))
			}
			mapValue, ok := value.(map[string]interface{})
			if !ok {
				return nil, errors.New(fmt.Sprintf("semantic error: value is not type map"))
			}
			return mapValue, nil
		default:
			return nil, errors.New(fmt.Sprintf("semantic error: value is not type literals"))
		}
	}

	return nil, errors.New(fmt.Sprintf("semantic error: %s is not declared", target.GetName()))
}

func makeVar(vTable []vars.Var, value values.Value, vName string) (vars.Var, error) {
	// getLiteral, getLiterals, getMapを順に回していき、適切なvalueを探す
	literal, err := getLiteral(vTable, value)
	if err == nil {
		return vars.Var{Name: vName, Value: values.Literal{Kind: values.LITERAL, Value: literal}}, nil
	}

	literals, err := getLiterals(vTable, value)
	if err == nil {
		return vars.Var{Name: vName, Value: values.Literals{Kind: values.LITERALS, Values: literals}}, nil
	}

	mapLiteral, err := getMap(vTable, value)
	if err == nil {
		return vars.Var{Name: vName, Value: values.Map{Kind: values.MAP, Value: mapLiteral}}, nil
	}

	return vars.Var{}, err
}

func getIndex(vTable []vars.Var, vName string) (int, error) {
	for i := len(vTable) - 1; i >= 0; i-- {
		if vTable[i].Name == vName {
			return i, nil
		}
	}

	return -1, errors.New(fmt.Sprintf("semantic error: %s is not declared", vName))
}

func assertionToStringSlice(target interface{}) ([]string, bool) {
	slices, ok := target.([]interface{})
	if !ok {
		return nil, false
	}

	var strings []string
	for _, v := range slices {
		str, ok := v.(string)
		if !ok {
			return nil, false
		}

		strings = append(strings, str)
	}

	return strings, true
}

func whiteSpaces(word string) string {
	var spaces string
	for i := 0; i < len(word)+1; i++ {
		spaces += " "
	}

	return spaces
}
