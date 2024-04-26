// Package formula are internal helpers for dealing with formula.
package formula

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/mmihic/sheets/src/pkg/participlex/backtrack"
)

// Various token types.
const (
	TokenTypeString    = "String"
	TokenTypeIdent     = "Ident"
	TokenTypeCellRange = "CellRange"
	TokenTypeTrue      = "True"
	TokenTypeFalse     = "False"
	TokenTypeNumber    = "Number"
)

// A Token is a lexical token.
type Token struct {
	Type     string
	Value    string
	Position lexer.Position

	eof bool
}

// EOF returns true if this token is the EOF.
func (t Token) EOF() bool {
	return t.eof
}

// LexString returns a new lexer for the given string.
func LexString(text string) (*Lexer, error) {
	l, err := lex.LexString("", text)
	if err != nil {
		return nil, err
	}

	return &Lexer{
		lex: backtrack.EnableBacktracking(l),
	}, nil
}

// Lexer is a lexer for formulas.
type Lexer struct {
	lex  backtrack.Lexer
	next []Token
}

// Next returns the next token from the lexer.
func (l *Lexer) Next() (Token, error) {
	if len(l.next) != 0 {
		tok := l.next[len(l.next)-1]
		l.next = l.next[:len(l.next)-1]
		return tok, nil
	}

	lexTok, err := l.lex.Next()
	if err != nil {
		return Token{}, err
	}

	if isStartOfString(lexTok) {
		l.lex.Push(lexTok)
		text, err := consumeString(l.lex)
		if err != nil {
			return Token{}, err
		}

		return Token{
			Type:     TokenTypeString,
			Value:    text,
			Position: lexTok.Pos,
		}, nil
	}

	tokenType, ok := symbols[lexTok.Type]
	if !ok {
		for symbol, typ := range lex.Symbols() {
			if typ == lexTok.Type {
				return Token{}, fmt.Errorf("unknown token type: '%s'", symbol)
			}
		}

		return Token{}, fmt.Errorf("unknown token type (%d) '%s'", lexTok.Type, lexTok.Value)
	}

	return Token{
		Type:     tokenType,
		Value:    lexTok.Value,
		Position: lexTok.Pos,
		eof:      lexTok.EOF(),
	}, nil
}

// Push pushes a set of tokens back onto the Lexer.
func (l *Lexer) Push(tokens ...Token) {
	l.next = append(l.next, tokens...)
}

func isStartOfString(tok lexer.Token) bool {
	switch symbols[tok.Type] {
	case "DoubleQuotes", "SingleQuotes", "TickQuotes":
		return true
	default:
		return false
	}
}

func consumeString(lex backtrack.Lexer) (string, error) {
	tok, err := lex.Next()
	if err != nil {
		return "", err
	}

	switch symbols[tok.Type] {
	case "DoubleQuotes":
		return consumeStringUsing(lex, "DoubleQuotedStringChars", "DoubleQuotes")
	case "SingleQuotes":
		return consumeStringUsing(lex, "SingleQuotedStringChars", "SingleQuotes")
	case "TickQuotes":
		return consumeStringUsing(lex, "TickQuotedStringChars", "TickQuotes")
	default:
		return "", TokenExpectedError(tok, "String")
	}
}

func consumeStringUsing(lex backtrack.Lexer, charsToken, stopToken string) (string, error) {
	var elts []string
	for {
		tok, err := lex.Next()
		if err != nil {
			return "", nil
		}

		switch symbols[tok.Type] {
		case charsToken, "Char":
			elts = append(elts, tok.Value)
		case stopToken:
			return strings.Join(elts, ""), nil
		default:
			return "", ParseErrorf(tok.Pos, "expected end of string: '%s'", tok.Value)
		}
	}
}

var (
	lex = lexer.MustStateful(lexer.Rules{
		"Root": {
			{`SingleQuotes`, `'`, lexer.Push("SingleQuotedString")},
			{`DoubleQuotes`, `"`, lexer.Push("DoubleQuotedString")},
			{"TickQuotes", "`", lexer.Push("TickQuotedString")},
			{"True", `[Tt][Rr][Uu][Ee]`, nil},
			{"False", `[Ff][Aa][Ll][Ss][Ee]`, nil},
			{"CellRange", `([A-Za-z]{1,3})?(\d+)?\s*:\s*([A-Za-z]{1,3})?(\d+)?`, nil},
			{"Ident", `[A-Za-z_][A-Za-z0-9_]*`, nil},
			{"!", "!", nil},
			{":", ":", nil},
			{",", ",", nil},
			{"(", `\(`, nil},
			{")", `\)`, nil},
			{"Number", `[0-9]+(\.[0-9]+)?`, nil},
			{"whitespace", `[\s]+`, nil},
		},
		"SingleQuotedString": {
			{"backslash", `\\`, lexer.Push("EscapedChar")},
			{"SingleQuotes", `'`, lexer.Pop()},
			{"SingleQuotedStringChars", `[^'\\]+`, nil},
		},
		"DoubleQuotedString": {
			{"backslash", `\\`, lexer.Push("EscapedChar")},
			{"DoubleQuotes", `"`, lexer.Pop()},
			{"DoubleQuotedStringChars", `[^"\\]+`, nil},
		},
		"TickQuotedString": {
			{"backslash", `\\`, lexer.Push("EscapedChar")},
			{"TickQuotes", "`", lexer.Pop()},
			{"TickQuotedStringChars", "[^`\\\\]+", nil},
		},
		"EscapedChar": {
			{"Char", `.`, lexer.Pop()},
		},
	})

	symbols = lexer.SymbolsByRune(lex)
)
