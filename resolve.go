package mathparse

import (
	"fmt"

	fmp "github.com/sourcekris/goflint"
)

// Resolve resolves the expression in Parser p into the resulting integer or expression string.
func (p *Parser) Resolve() {
	// parenthases
	// exponents/roots
	// multiplication/division
	// addition/subtraction
	// functions
	// repeat
	p.tokens = parseExpression(p.tokens)
}

// Eval continuously resolves the expression until we cannot resolve any further and returns the
// Value Result or nil if the expression never resolves.
func (p *Parser) Eval() *fmp.Fmpz {
	el := len(p.GetExpressionResult())
	fmt.Printf("el: %v\n", el)
	for {
		p.Resolve()
		if p.FoundResult() {
			break
		}

		if len(p.GetExpressionResult()) == el {
			return nil
		}

		el = len(p.GetExpressionResult())
	}

	return p.GetValueResult()
}

// Eval will create a new parser from an expressions and Evaluate the expression until it resolves
// or results in error. Returns an fmpz or an error.
func Eval(expression string) (*fmp.Fmpz, error) {
	p := NewParser(expression)

	if val := p.Eval(); val != nil {
		return val, nil
	}

	return nil, fmt.Errorf("eval could not resolve expression %q: %v", expression, p.GetExpressionResult())
}

func parseExpression(set []Token) []Token {
	mod := false
	if set[0].Type == function || set[0].Type == lparen {
		set[0].Children = parseExpression(set[0].Children)
	}
	for i := 1; i < len(set)-1; i++ {
		if set[i].Type == operation {
			if (set[i].Value == "^") && (set[i-1].Type == literal && set[i+1].Type == literal) {
				mod = true
				set[i-1].ParseValue = new(fmp.Fmpz).Exp(set[i-1].ParseValue, set[i+1].ParseValue, nil)
				set[i-1].Value = set[i-1].ParseValue.String()
				set = append(set[:i], set[i+2:]...)
				i--
			}
		}
	}
	for i := 1; i < len(set)-1; i++ {
		if set[i].Type == operation {
			if (set[i].Value == "*" || set[i].Value == "/" || set[i].Value == "%") && (set[i-1].Type == literal && set[i+1].Type == literal) {
				mod = true
				if set[i].Value == "*" {
					set[i-1].ParseValue = new(fmp.Fmpz).Mul(set[i-1].ParseValue, set[i+1].ParseValue)
				} else if set[i].Value == "/" {
					set[i-1].ParseValue = new(fmp.Fmpz).Div(set[i-1].ParseValue, set[i+1].ParseValue)
				} else if set[i].Value == "%" {
					set[i-1].ParseValue = new(fmp.Fmpz).Mod(set[i-1].ParseValue, set[i+1].ParseValue)
				}
				set[i-1].Value = set[i-1].ParseValue.String()
				set = append(set[:i], set[i+2:]...)
				i--
			}
		}
	}

	for i := 1; i < len(set)-1; i++ {
		if set[i].Type == operation {
			if (set[i].Value == "+" || set[i].Value == "-") && (set[i-1].Type == literal && set[i+1].Type == literal) {
				mod = true
				if set[i].Value == "+" {
					set[i-1].ParseValue = new(fmp.Fmpz).Add(set[i-1].ParseValue, set[i+1].ParseValue)
				} else if set[i].Value == "-" {
					set[i-1].ParseValue = new(fmp.Fmpz).Sub(set[i-1].ParseValue, set[i+1].ParseValue)
				}
				set[i-1].Value = set[i-1].ParseValue.String()
				set = append(set[:i], set[i+2:]...)
				i--
			}
		}
	}

	// functions
	for i := range set {
		if set[i].Type == lparen {
			if len(set[i].Children) == 1 {
				set[i] = set[i].Children[0]
			} else {
				set[i].Children = parseExpression(set[i].Children)
			}
		}
		if set[i].Type == function {
			mod = true
			switch set[i].Value {
			case "abs":
				set[i] = newToken(literal, new(fmp.Fmpz).Abs(set[i].Children[0].ParseValue).String())
			case "sqrt":
				set[i] = newToken(literal, new(fmp.Fmpz).Root(set[i].Children[0].ParseValue, 2).String())
			case "max":
				set[i] = newToken(literal, new(fmp.Fmpz).Max(set[i].Children[0].ParseValue, set[i].Children[2].ParseValue).String())
			case "min":
				set[i] = newToken(literal, new(fmp.Fmpz).Min(set[i].Children[0].ParseValue, set[i].Children[2].ParseValue).String())
			case "mod":
				fmt.Printf("set[i]: %v\n", set[i].String())
				set[i] = newToken(literal, new(fmp.Fmpz).Mod(set[i].Children[0].ParseValue, set[i].Children[2].ParseValue).String())
			case "pow":
				set[i] = newToken(literal, new(fmp.Fmpz).Exp(set[i].Children[0].ParseValue, set[i].Children[2].ParseValue, nil).String())
			case "invmod":
				set[i] = newToken(literal, new(fmp.Fmpz).ModInverse(set[i].Children[0].ParseValue, set[i].Children[2].ParseValue).String())
			}
		}
	}

	if len(set) > 1 && mod {
		set = parseExpression(set)
	}

	return set
}
