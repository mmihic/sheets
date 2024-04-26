package formula

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

// TokenExpectedError returns an error when we receive a token other
// than what was expected.
func TokenExpectedError(tok lexer.Token, expected ...string) error {
	return ParseErrorf(tok.Pos, "expected one of [%s], found '%s'",
		strings.Join(expected, ", "), tok.Value)
}

// ParseErrorf returns a formula error.
func ParseErrorf(pos lexer.Position, msg string, args ...any) error {
	return WrapParseError(pos, fmt.Errorf(msg, args...))
}

// WrapParseError wraps an error in a parse error.
func WrapParseError(pos lexer.Position, err error) error {
	return fmt.Errorf("error at %s: %w", pos, err)
}
