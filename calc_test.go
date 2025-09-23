package main

import "testing"

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
	}

	for _, tt := range tests {
		got := eval(tt.expr)
		if got != tt.want {
			t.Errorf("FAIL expr=%q got=%d want=%d", tt.expr, got, tt.want)
		} else {
			t.Logf("PASS expr=%q = %d", tt.expr, got)
		}
	}
}
