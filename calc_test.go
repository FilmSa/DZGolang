package main

import (
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		expr string
		want int
	}{
		{"1+2", 3},
		{"4-2", 2},
		{"2*3", 6},
		{"8/2", 4},
		{"(1+2)*3", 9},
		{"(10-2)/2", 4},
		{"10-(2+3)", 5},
		{"(2+3)*(4-1)", 15},
		{"-2+3", 1},
		{"-(2+3)", -5},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			got, err := eval(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expr=%q got=%d want=%d", tt.expr, got, tt.want)
			}
		})
	}
}

func TestInvalidExpressions(t *testing.T) {
	tests := []string{
		"",     
		"1++2", 
		"(1+2", 
		"1/0",  
		"abc",  
		"()",  
	}

	for _, expr := range tests {
		t.Run(expr, func(t *testing.T) {
			_, err := eval(expr)
			if err == nil {
				t.Errorf("expected error for expr=%q but got none", expr)
			}
		})
	}
}
