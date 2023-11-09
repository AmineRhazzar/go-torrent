package stackops

import "fmt"

func Push(stack *[]string, el string) {
	*stack = append(*stack, el);
}

func Pop(stack *[]string) (string, error) {
	len := len(*stack)
	if len == 0 {
		var t string
		return t, fmt.Errorf("pop error on stack: stack is empty")
	}
	poppedElement := (*stack)[len-1]
	(*stack) = (*stack)[:len-1]
	return poppedElement, nil
}

func Front(stack *[]string) (string, error) {
	len := len(*stack)
	if len == 0 {
		var t string
		return t, fmt.Errorf("pop error on stack: stack is empty")
	}
	frontElement := (*stack)[len-1]
	return frontElement, nil
}