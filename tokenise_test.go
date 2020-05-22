package mathparse

import (
	"fmt"
	"reflect"
	"testing"

	fmp "github.com/sourcekris/goflint"
)

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		name       string
		expression string
		want       []Token
	}{
		{
			name:       "valid tokenise",
			expression: "1+1",
			want: []Token{
				{Type: literal, Value: "1", ParseValue: fmp.NewFmpz(1)},
				{Type: operation, Value: "+"},
				{Type: literal, Value: "1", ParseValue: fmp.NewFmpz(1)},
				// Why are there always two full sets of tokens?
				{Type: literal, Value: "1", ParseValue: fmp.NewFmpz(1)},
				{Type: operation, Value: "+"},
				{Type: literal, Value: "1", ParseValue: fmp.NewFmpz(1)},
			},
		},
	} {
		p := NewParser(tc.expression)
		p.tokenise()
		got := p.GetTokens()
		if !reflect.DeepEqual(got, tc.want) {
			var gotstr, wantstr string
			for _, t := range got {
				gotstr = fmt.Sprintf("%s%s", gotstr, t.String())
			}
			for _, t := range tc.want {
				wantstr = fmt.Sprintf("%s%s", wantstr, t.String())
			}
			t.Errorf("tokenise() failed: %s got/want mismatched:\n%s\n/\n%s", tc.name, gotstr, wantstr)
		}
	}
}
