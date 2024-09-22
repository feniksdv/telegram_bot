package services

import (
	"fmt"
	"github.com/username/expression"
)

func Calc(text string) string {
	result, err := expression.Eval(text)
	if err != nil {
		fmt.Println("Error evaluating expression:", err)
		return
	}
	fmt.Printf("Result: %.4f\n", result)
	return text
}
