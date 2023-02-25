package parser

import (
	"dcc/types"
)

func Parse(tokens []types.Token) (error) {
	err := program(tokens, 0)
	if err != nil {
		return err
	}

	return nil
}
