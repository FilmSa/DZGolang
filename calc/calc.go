package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func priority(op rune) int {
	if op == '+' || op == '-' {
		return 1
	}
	if op == '*' || op == '/' {
		return 2
	}
	return 0
}

func apply(a, b int, op rune) int {
	switch op {
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	case '/':
		return a / b
	}
	return 0
}

func eval(expr string) int {
	nums := []int{}
	ops := []rune{}
	numStr := ""

	pushNum := func() {
		if numStr != "" {
			val, _ := strconv.Atoi(numStr)
			nums = append(nums, val)
			numStr = ""
		}
	}

	for _, ch := range expr {
		if unicode.IsSpace(ch) {
			continue
		}
		if unicode.IsDigit(ch) {
			numStr += string(ch)
			continue
		}
		pushNum()
		if ch == '(' {
			ops = append(ops, ch)
		} else if ch == ')' {
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				b := nums[len(nums)-1]
				a := nums[len(nums)-2]
				nums = nums[:len(nums)-2]
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]
				nums = append(nums, apply(a, b, op))
			}
			ops = ops[:len(ops)-1] // убираем '('
		} else {
			for len(ops) > 0 && priority(ops[len(ops)-1]) >= priority(ch) {
				b := nums[len(nums)-1]
				a := nums[len(nums)-2]
				nums = nums[:len(nums)-2]
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]
				nums = append(nums, apply(a, b, op))
			}
			ops = append(ops, ch)
		}
	}

	pushNum()

	for len(ops) > 0 {
		b := nums[len(nums)-1]
		a := nums[len(nums)-2]
		nums = nums[:len(nums)-2]
		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]
		nums = append(nums, apply(a, b, op))
	}

	if len(nums) == 0 {
		return 0
	}
	return nums[0]
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error")
		return
	}
	expr := os.Args[1]
	expr = strings.TrimSpace(expr)
	res := eval(expr)
	fmt.Println(res)
}
