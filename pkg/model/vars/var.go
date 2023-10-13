package vars

import "github.com/ty-bnn/myriad/pkg/model/values"

// Var will be used in variable table.
type Var struct {
	Name  string
	Value values.Value
}
