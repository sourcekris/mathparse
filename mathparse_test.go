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
			name:       "valid rsa encryption calculation",
			expression: "2^607-1",
			want:       "531137992816767098689588206552468627329593117727031923199444138200403559860852242739162502265229285668889329486246501015346579337652707239409519978766587351943831270835393219031728127",
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
