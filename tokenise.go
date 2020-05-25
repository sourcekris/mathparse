package mathparse

import (
	"fmt"
	"regexp"
	"strings"

	fmp "github.com/sourcekris/goflint"
)

// Parser encapsulates the parsed tokens and values.
type Parser struct {
	letterBuffer string
	numberBuffer string
	expression   string
	tokens       []Token
}

// Token encodes a single token and its children.
type Token struct {
	Type       TokenType
	Value      string
	ParseValue *fmp.Fmpz
	Children   []Token
}

// TokenType describes the type of token.
type TokenType uint

const (
	undefined TokenType = iota // 0
	space                      // 1
	literal                    // 2
	variable                   // 3
	operation                  // 4
	function                   // 5
	lparen                     // 6
	rparen                     // 7
	funcDelim                  // 8
)

// String returns a string representation of a token.
func (t *Token) String() string {
	res := fmt.Sprintf("Type: %d, Value: %s, ParseValue: %v, Num. Children: %d\n", t.Type, t.Value, t.ParseValue, len(t.Children))
	if t.Children != nil {
		for _, c := range t.Children {
			res = fmt.Sprintf("%s\tChild: %s", res, c.String())
		}
	}
	return res
}

// ReadExpression tokenises an expression given string str.
func (p *Parser) ReadExpression(str string) {
	p.expression = str
	p.tokens = []Token{}
	p.tokenise()
}

// ReadMultipartExpression tokenizes an expression expressed as multiple parts in a string slice
// str.
func (p *Parser) ReadMultipartExpression(str []string) {
	p.expression = strings.Join(str, " ")
	p.tokens = []Token{}
	p.tokenise()
}

func (p *Parser) tokenise() {
	dumpLetter := func(p *Parser) {
		for i, ch := range p.letterBuffer {
			p.tokens = append(p.tokens, newToken(variable, string(ch)))
			if i < len(p.letterBuffer)-1 {
				p.tokens = append(p.tokens, newToken(operation, "*"))
			}
		}
		p.letterBuffer = ""
	}
	dumpNumber := func(p *Parser) {
		if len(p.numberBuffer) > 0 {
			p.tokens = append(p.tokens, newToken(literal, p.numberBuffer))
			p.numberBuffer = ""
		}
	}
	for i, ch := range p.expression {
		switch getTokenType(ch) {
		case space:
			continue
		case literal:
			p.numberBuffer += string(ch)
		case variable:
			if len(p.numberBuffer) > 0 {
				dumpNumber(p)
				p.tokens = append(p.tokens, newToken(operation, "*"))
			}
			p.letterBuffer += string(ch)
		case operation:
			// A hack to support leading negative numbers, prepend a 0.
			if i == 0 && ch == '-' {
				p.tokens = append(p.tokens, newToken(literal, "0"))
			}
			dumpNumber(p)
			dumpLetter(p)
			p.tokens = append(p.tokens, newToken(operation, string(ch)))
		case lparen:
			if len(p.numberBuffer) > 0 {
				dumpNumber(p)
				p.tokens = append(p.tokens, newToken(operation, "*"))
			}
			if len(p.letterBuffer) > 0 {
				p.tokens = append(p.tokens, newToken(function, p.letterBuffer))
				p.letterBuffer = ""
			}
			p.tokens = append(p.tokens, newToken(lparen, "("))
		case rparen:
			dumpLetter(p)
			dumpNumber(p)
			p.tokens = append(p.tokens, newToken(rparen, ")"))
		case funcDelim:
			dumpNumber(p)
			dumpLetter(p)
			p.tokens = append(p.tokens, newToken(funcDelim, ","))
		}
	}

	if len(p.numberBuffer) > 0 {
		dumpNumber(p)
	}
	if len(p.letterBuffer) > 0 {
		dumpLetter(p)
	}

	p.tokens, _ = buildTree(p.tokens)
}

func buildTree(set []Token) ([]Token, int) {
	toks := []Token{}
	for i := 0; i < len(set); i++ {
		tok := set[i]
		switch tok.Type {
		case function:
			child, offset := buildTree(set[i+2:])
			tok.Children = child
			i += 2 + offset
			toks = append(toks, tok)
		case lparen:
			child, offset := buildTree(set[i+1:])
			tok.Children = child
			toks = append(toks, tok)
			i += 1 + offset
		case rparen:
			// toks = append(toks, tok)
			return toks, i
		default:
			toks = append(toks, tok)
		}
	}
	return toks, len(set)
}

func newToken(typ TokenType, value string) Token {
	tok := Token{
		Type:  typ,
		Value: value,
	}
	if typ == literal {
		tok.ParseValue, _ = new(fmp.Fmpz).SetString(value, 10)
	}
	return tok
}

func getTokenType(ch rune) TokenType {
	let := string(ch)
	if let == " " {
		return space
	} else if isDigit(let) {
		return literal
	} else if isLetter(let) {
		return variable
	} else if isOperator(let) {
		return operation
	} else if isOpenParen(let) {
		return lparen
	} else if isCloseParen(let) {
		return rparen
	} else if isComma(let) {
		return funcDelim
	}
	return undefined
}

func isComma(let string) bool {
	return let == ","
}

func isDigit(let string) bool {
	res, err := regexp.MatchString(`[0-9\.]`, let)
	if err != nil {
		fmt.Print(err)
	}
	return res
}

func isLetter(let string) bool {
	res, err := regexp.MatchString("[a-zA-Z]", let)
	if err != nil {
		fmt.Print(err)
	}
	return res
}

func isOperator(let string) bool {
	res, err := regexp.MatchString(`\*|\/|\+|\^|\-|%`, let)
	// "\x43|\x47|\x42|\x94|\x45"
	if err != nil {
		fmt.Print(err)
	}
	return res
}

func isOpenParen(let string) bool {
	return let == "("
}

func isCloseParen(let string) bool {
	return let == ")"
}
