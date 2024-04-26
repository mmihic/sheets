// Package backtrack contains a lexer that allows backtracking
package backtrack

import "github.com/alecthomas/participle/v2/lexer"

// A Lexer is a Lexer that can back-track, allowing for
// tokens to be pushed "back" onto the Lexer.
type Lexer interface {
	lexer.Lexer
	Push(tokens ...lexer.Token)
}

// EnableBacktracking takes a Lexer and turns it into a Lexer.
func EnableBacktracking(lex lexer.Lexer) Lexer {
	return &backTrackingLexer{
		lex: lex,
	}
}

type backTrackingLexer struct {
	lex    lexer.Lexer
	buffer []lexer.Token
}

func (l *backTrackingLexer) Push(tokens ...lexer.Token) {
	l.buffer = append(l.buffer, tokens...)
}

func (l *backTrackingLexer) Next() (lexer.Token, error) {
	if len(l.buffer) > 0 {
		tok := l.buffer[len(l.buffer)-1]
		l.buffer = l.buffer[:len(l.buffer)-1]
		return tok, nil
	}

	return l.lex.Next()
}
