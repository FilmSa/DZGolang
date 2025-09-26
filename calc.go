package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type numStack []int

func (s *numStack) push(v int) {
	*s = append(*s, v)
}

func (s *numStack) pop() (int, error) {
	if len(*s) == 0 {
		return 0, errors.New("pop from empty num stack")
	}
	val := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return val, nil
}

type opStack []rune

func (s *opStack) push(op rune) {
	*s = append(*s, op)
}

func (s *opStack) pop() (rune, error) {
	if len(*s) == 0 {
		return 0, errors.New("pop from empty op stack")
	}
	val := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return val, nil
}

func (s *opStack) peek() (rune, bool) {
	if len(*s) == 0 {
		return 0, false
	}
	return (*s)[len(*s)-1], true
}

func priority(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	default:
		return 0
	}
}

func apply(a, b int, op rune) (int, error) {
	switch op {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	}
	return 0, fmt.Errorf("unknown operator: %q", op)
}

func eval(expr string) (int, error) {
	var nums numStack
	var ops opStack
	numStr := ""
	prevTokenIsOp := true
	pushNum := func() error {
		if numStr != "" {
			if numStr == "-" {
				return errors.New("dangling unary minus")
			}
			val, err := strconv.Atoi(numStr)
			if err != nil {
				return err
			}
			nums.push(val)
			numStr = ""
			prevTokenIsOp = false
		}
		return nil
	}

	for _, ch := range expr {
		if unicode.IsSpace(ch) {
			continue
		}
		if unicode.IsDigit(ch) {
			numStr += string(ch)
			prevTokenIsOp = false
			continue
		}

		if ch == '-' && prevTokenIsOp {
			numStr += string(ch)
			continue
		}

		if ch == '(' {
			if numStr == "-" {
				nums.push(0)
				ops.push('-')
				numStr = ""
			} else if err := pushNum(); err != nil {
				return 0, err
			}
			ops.push(ch)
			prevTokenIsOp = true
		} else if ch == ')' {
			if err := pushNum(); err != nil {
				return 0, err
			}
			for {
				top, ok := ops.peek()
				if !ok {
					return 0, errors.New("mismatched parentheses")
				}
				if top == '(' {
					_, _ = ops.pop()
					break
				}
				if err := calcOnce(&nums, &ops); err != nil {
					return 0, err
				}
			}
			prevTokenIsOp = false
		} else {
			if err := pushNum(); err != nil {
				return 0, err
			}
			for {
				top, ok := ops.peek()
				if !ok || priority(top) < priority(ch) {
					break
				}
				if err := calcOnce(&nums, &ops); err != nil {
					return 0, err
				}
			}
			ops.push(ch)
			prevTokenIsOp = true
		}
	}

	if err := pushNum(); err != nil {
		return 0, err
	}

	for len(ops) > 0 {
		if err := calcOnce(&nums, &ops); err != nil {
			return 0, err
		}
	}

	if len(nums) != 1 {
		return 0, errors.New("invalid expression")
	}
	return nums[0], nil
}

func calcOnce(nums *numStack, ops *opStack) error {
	b, err := nums.pop()
	if err != nil {
		return err
	}
	a, err := nums.pop()
	if err != nil {
		return err
	}
	op, err := ops.pop()
	if err != nil {
		return err
	}
	res, err := apply(a, b, op)
	if err != nil {
		return err
	}
	nums.push(res)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: no expression")
		return
	}
	expr := strings.TrimSpace(os.Args[1])
	res, err := eval(expr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(res)
}
