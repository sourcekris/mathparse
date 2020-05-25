package mathparse

import (
	"testing"
)

func TestGetExpressionResults(t *testing.T) {
	for _, tc := range []struct {
		name       string
		expression string
		want       string
	}{
		{
			name:       "valid mersenne prime calculation",
			expression: "2^607-1",
			want:       "531137992816767098689588206552468627329593117727031923199444138200403559860852242739162502265229285668889329486246501015346579337652707239409519978766587351943831270835393219031728127",
		},
		{
			name:       "valid rsa encryption calculation - c = m^e % n",
			expression: "1289^3%25777",
			want:       "18524",
		},
		{
			name:       "valid algebraic equation",
			expression: "5x+1+1",
			want:       "5*x+2",
		},
	} {
		p := NewParser(tc.expression)
		p.Resolve()
		got := p.GetExpressionResult()
		if got != tc.want {
			t.Errorf("GetExpressionResults() failed: %s got/want mismatched:\n%s\n/\n%s", tc.name, got, tc.want)
		}
	}
}

func TestGetValueResult(t *testing.T) {
	for _, tc := range []struct {
		expression string
		wantErr    bool
		want       int
	}{
		// Copied the unit tests from a similar library.
		{"2^2", false, 4},
		{"1+1", false, 2},
		{"-1+2", false, 1},
		{"2-1", false, 1},
		{"1-10", false, -9},
		{"1+2*3", false, 7},
		{"2*3+1", false, 7},
		{"2*3/2", false, 3},
		{"2/2*3", false, 3},
		// Testing precedence.
		{"1+2*3/2", false, 4},
		{"-3+3*2+5-2*2", false, 4},
		{"4+3-2+1", false, 6},
		{"2-3+4-2", false, 1},
		{"24*3+15*2-31*4-1+2", false, -21},
		// Testing brackets
		{"(1+2)*3", false, 9},
		{"3*(1+2)", false, 9},
		{"3*(1+2)*4", false, 36},
		// Embedded expressions.
		{"(3*4+(4*7)+55)+(5*5)", false, 120},
		// Functions.
		{"mod(300, 40)", false, 20},
	} {
		//t.Errorf("%v", new(fmp.Fmpz).Add(fmp.NewFmpz(-1), fmp.NewFmpz(2)))
		p := NewParser(tc.expression)
		v := p.Eval()
		var got int
		if v != nil {
			got = v.GetInt()
		} else {
			t.Fatalf("GetValueResult() failed: %s got nil Fmpz", tc.expression)
		}

		if got != tc.want {
			t.Errorf("GetValueResult() %q failed: got/want mismatched: %d / %d", tc.expression, got, tc.want)
		}
	}
}
