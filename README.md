# mathparse

golang library for parsing maths expression strings into arbitrary precision integers using the 
FLINT library type Fmpz.

## Purpose

To simplify the expression of complex integer math in other libraries and programs I'm writing.

## Original

This library was forked from [github.com/Maldris/mathparse](https://github.com/Maldris/mathparse)
and was originally intended for the float64 type only. I removed float support as I do not need it
and modified it to support only the arbitrarty precision integer type Fmpz. I also added unit tests
to verify the library produces expected results.

## Limitations

Variables used in expressions can only be single letter, i.e.: `xy` is evaluated as `x*y`

## Usage
### Parser
The core functionality is provided by the Parser object, if you have the expression ready, you can 
create the parser, and process your expression at the same time with

```go
  expression := "2^1239-1"
  p := mathparse.NewParser(expression)
```
At this point the expression is tokenised and ready to parse.

This process is separated, and the Token class exported, to allow people to build other resolving
logic if they so wish.

To resolve the expression, call `Resolve`

```go
  p.Resolve()
```

When this is done, there are two possible results, either the expression has resolved down to a 
single integer value, and can be output as a Fmpz, or (due to variables in use, or an FmpzPoly) a 
potentially simplified expression string can be retreived.

To know which Option to use, check if the expression is a value with `FoundResult`, a return value
of true means an integer value can be retrieved, otherwise, a polynomial

```go
  if p.FoundResult() {
    var result *fmp.Fmpz
    result = p.GetValueResult()
    log.Print(result)
  } else {
    var expression string
    expression = p.GetExpressionResult()
    log.Print(expression)
  }
```

Here `GetValueResult` retreives the integer result of the expression, and `GetExpressionResult` will
return the expressions simplified form


Its worth noting that the parser object is reusable if need be, if after parsing one expression you
wish to parse another, simply load your next one with either `ReadExpression` or `ReadMultipartExpression`.

Each of which will read the expression in and tokenise it as `NewParser` does, but on an existing
parser.
`ReadMultipartExpression` exists so that if you have an expression already in multiple parts (i.e.
separated as function inputs to a text/template function call) the library is still simple to use.
It will simple concatenate the expression segments, and proceed to attempt to resolve the resultant
expression.

If multiple segment data is intended, but not desired as a string expression, the raw tokens can be
retreived via `GetTokens`, which will return the raw token tree from the parser.

Useful if you have funciton aruements separated by commas, and you want to then take the result and
pass into an external function.


### Tokens

The structure of each token is quite simple:

```go
  type Token struct {
    Type       TokenType
    Value      string
    ParseValue *fmp.Fmpz
    Children   []Token
  }
```

Type is the type of token, from an enum (see below), value, which is the raw string value of the
token (i.e. "3.8", "+", "a", etc), ParseValue, which is the integer of the value, if the type is a
literal. And lastly, Children, which will contain child tokens nested under this token, which is
only the case for functions, and Parenthesis.

TokenType may take on the following values:

```go
const (
  undefined TokenType = iota // 0 - unknown token character
  space                      // 1 - space character, ignored
  literal                    // 2 - a literal, a number
  variable                   // 3 - variables
  operation                  // 4 - any of the following mathematical operations: * / + - ^ %
  function                   // 5 - a function, it will have the expression for its arguements as Child tokens
  lparen                     // 6 - opening parenthesis, will have the enclosed expression as Child tokens
  rparen                     // 7 - closing parenthesis, used internally, stripped in tree creation, used to mark the end of the current function or parenthesis
  funcDelim                  // 8 - delimits function arguements, doesnt do anything, but prevents, adjacent expressions being evaluated together
)
```
