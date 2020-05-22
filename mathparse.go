package mathparse

import fmp "github.com/sourcekris/goflint"

// NewParser constructs a new parser from expression.
func NewParser(expression string) Parser {
	parse := Parser{}
	parse.ReadExpression(expression)
	return parse
}

// FoundResult returns true if the expression has a result.
func (p *Parser) FoundResult() bool {
	return len(p.tokens) <= 1 && p.tokens[0].Type == literal
}

// GetValueResult returns the integer value from the expression.
func (p *Parser) GetValueResult() *fmp.Fmpz {
	return p.tokens[0].ParseValue
}

// GetExpressionResult returns the string representation of the expression.
func (p *Parser) GetExpressionResult() string {
	return getStringExpression(p.tokens)
}

func getStringExpression(set []Token) string {
	str := ""
	for _, tok := range set {
		switch tok.Type {
		case space:
			str += " "
		case literal:
			str += tok.Value
		case variable:
			str += tok.Value
		case operation:
			str += tok.Value
		case function:
			str += tok.Value + "(" + getStringExpression(tok.Children) + ")"
		case lparen:
			str += "(" + getStringExpression(tok.Children) + ")"
		case funcDelim:
			str += ","
		}
	}
	return str
}

// GetTokens returns the tokens from the parser.
func (p *Parser) GetTokens() []Token {
	return p.tokens
}
