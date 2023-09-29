package generator

import (
	"fmt"
	"os"
)

func (g *Generator) Generate() error {
	fmt.Println("Generating...")

	mainArgs := os.Args[3:]

	err := g.callFunc("main", mainArgs)
	if err != nil {
		return err
	}

	fmt.Println("Generate Done.")

	return nil
}
